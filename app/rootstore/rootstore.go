package rootstore

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tus/tusd/v2/pkg/handler"
)

var defaultFilePerm = os.FileMode(0664)
var defaultDirectoryPerm = os.FileMode(0754)

const (
	// StorageKeyPath is the key of the path of uploaded file in handler.FileInfo.Storage
	StorageKeyPath = "Path"
	// StorageKeyInfoPath is the key of the path of .info file in handler.FileInfo.Storage
	StorageKeyInfoPath = "InfoPath"
)

// RootStore is a file system based data store for tusd.
type RootStore struct {
	root *os.Root
}

func New(root *os.Root) RootStore {
	return RootStore{root: root}
}

// UseIn sets this store as the core data store in the passed composer and adds
// all possible extension to it.
func (store RootStore) UseIn(composer *handler.StoreComposer) {
	composer.UseCore(store)
	composer.UseTerminater(store)
	composer.UseConcater(store)
	composer.UseLengthDeferrer(store)
	composer.UseContentServer(store)
}

func (store RootStore) NewUpload(ctx context.Context, info handler.FileInfo) (handler.Upload, error) {
	if info.ID == "" {
		info.ID = rand.Text()
	}

	// The .info file's location can directly be deduced from the upload ID
	infoPath := store.infoPath(info.ID)
	// The binary file's location might be modified by the pre-create hook.
	var binPath string
	if info.Storage != nil && info.Storage[StorageKeyPath] != "" {
		binPath = info.Storage[StorageKeyPath]
	} else {
		binPath = store.defaultBinPath(info.ID)
	}

	info.Storage = map[string]string{
		"Type":             "rootstore",
		StorageKeyPath:     binPath,
		StorageKeyInfoPath: infoPath,
	}

	_ = store.root.MkdirAll(filepath.Dir(binPath), defaultDirectoryPerm)

	// Create binary file with no content
	if err := store.root.WriteFile(binPath, nil, defaultFilePerm); err != nil {
		return nil, err
	}

	upload := &fileUpload{
		root:     store.root,
		info:     info,
		infoPath: infoPath,
		binPath:  binPath,
	}

	// writeInfo creates the file by itself if necessary
	if err := upload.writeInfo(); err != nil {
		return nil, err
	}

	return upload, nil
}

func (store RootStore) GetUpload(ctx context.Context, id string) (handler.Upload, error) {
	infoPath := store.infoPath(id)
	data, err := fs.ReadFile(store.root.FS(), filepath.ToSlash(infoPath))
	if err != nil {
		if os.IsNotExist(err) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = handler.ErrNotFound
		}
		return nil, err
	}
	var info handler.FileInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	// If the info file contains a custom path to the binary file, we use that. If not, we
	// fall back to the default value (although the Path property should always be set in recent
	// tusd versions).
	var binPath string
	if info.Storage != nil && info.Storage[StorageKeyPath] != "" {
		// No filepath.Join here because the joining already happened in NewUpload. Duplicate joining
		// with relative paths lead to incorrect paths
		binPath = info.Storage[StorageKeyPath]
	} else {
		binPath = store.defaultBinPath(info.ID)
	}

	stat, err := store.root.Stat(binPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = handler.ErrNotFound
		}
		return nil, err
	}

	info.Offset = stat.Size()

	return &fileUpload{
		root:     store.root,
		info:     info,
		binPath:  binPath,
		infoPath: infoPath,
	}, nil
}

func (store RootStore) AsTerminatableUpload(upload handler.Upload) handler.TerminatableUpload {
	return upload.(*fileUpload)
}

func (store RootStore) AsLengthDeclarableUpload(upload handler.Upload) handler.LengthDeclarableUpload {
	return upload.(*fileUpload)
}

func (store RootStore) AsConcatableUpload(upload handler.Upload) handler.ConcatableUpload {
	return upload.(*fileUpload)
}

func (store RootStore) AsServableUpload(upload handler.Upload) handler.ServableUpload {
	return upload.(*fileUpload)
}

// defaultBinPath returns the path to the file storing the binary data, if it is
// not customized using the pre-create hook.
func (store RootStore) defaultBinPath(id string) string {
	return id
}

// infoPath returns the path to the .info file storing the file's info.
func (store RootStore) infoPath(id string) string {
	return id + ".info"
}

type fileUpload struct {
	root *os.Root

	// info stores the current information about the upload
	info handler.FileInfo
	// infoPath is the path to the .info file
	infoPath string
	// binPath is the path to the binary file (which has no extension)
	binPath string
}

func (upload *fileUpload) GetInfo(ctx context.Context) (handler.FileInfo, error) {
	return upload.info, nil
}

func (upload *fileUpload) WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error) {
	file, err := upload.root.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return 0, err
	}
	// Avoid the use of defer file.Close() here to ensure no errors are lost
	// See https://github.com/tus/tusd/issues/698.

	n, err := io.Copy(file, src)
	upload.info.Offset += n
	if err != nil {
		file.Close()
		return n, err
	}

	return n, file.Close()
}

func (upload *fileUpload) GetReader(ctx context.Context) (io.ReadCloser, error) {
	return upload.root.Open(upload.binPath)
}

func (upload *fileUpload) Terminate(ctx context.Context) error {
	// We ignore errors indicating that the files cannot be found because we want
	// to delete them anyways. The files might be removed by a cron job for cleaning up
	// or some file might have been removed when tusd crashed during the termination.
	err := upload.root.Remove(upload.binPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = upload.root.Remove(upload.infoPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func (upload *fileUpload) ConcatUploads(ctx context.Context, uploads []handler.Upload) (err error) {
	file, err := upload.root.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return err
	}
	defer func() {
		// Ensure that close error is propagated, if it occurs.
		// See https://github.com/tus/tusd/issues/698.
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	for _, partialUpload := range uploads {
		if err := partialUpload.(*fileUpload).appendTo(file); err != nil {
			return err
		}
	}

	return
}

func (upload *fileUpload) appendTo(file *os.File) error {
	src, err := upload.root.Open(upload.binPath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(file, src); err != nil {
		src.Close()
		return err
	}

	return src.Close()
}

func (upload *fileUpload) DeclareLength(ctx context.Context, length int64) error {
	upload.info.Size = length
	upload.info.SizeIsDeferred = false
	return upload.writeInfo()
}

// writeInfo updates the entire information. Everything will be overwritten.
func (upload *fileUpload) writeInfo() error {
	data, err := json.Marshal(upload.info)
	if err != nil {
		return err
	}

	_ = upload.root.MkdirAll(filepath.Dir(upload.infoPath), defaultDirectoryPerm)

	return upload.root.WriteFile(upload.infoPath, data, defaultFilePerm)
}

func (upload *fileUpload) FinishUpload(ctx context.Context) error {
	return nil
}

func (upload *fileUpload) ServeContent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	http.ServeFileFS(w, r, upload.root.FS(), filepath.ToSlash(upload.binPath))

	return nil
}
