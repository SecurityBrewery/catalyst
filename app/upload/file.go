package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tus/tusd/v2/pkg/filelocker"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/rootstore"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

type Uploader struct {
	queries *sqlc.Queries
	auth    *auth.Service
	root    *os.Root
}

func NewUploader(dir string, auth *auth.Service, queries *sqlc.Queries) (*Uploader, error) {
	uploadsDir := path.Join(dir, "uploads")

	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create uploads directory: %w", err)
	}

	root, err := os.OpenRoot(uploadsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open uploads directory: %w", err)
	}

	return &Uploader{
		queries: queries,
		root:    root,
		auth:    auth,
	}, nil
}

func (u *Uploader) Routes() (http.Handler, error) {
	store := rootstore.New(u.root)
	locker := filelocker.New(u.root.Name())
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	// Create a new HTTP handler for the tusd server by providing a configuration.
	// The StoreComposer property must be set to allow the handler to function.
	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
		PreUploadCreateCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
			// This hook is called before an upload is created. You can use it to
			// modify the upload information, for example to set a custom ID or
			// storage path.
			id := database.GenerateID("")

			if hook.Upload.Storage == nil {
				hook.Upload.Storage = make(map[string]string)
			}

			filename, ok := hook.Upload.MetaData["filename"]
			if !ok || filename == "" {
				filename = id
			}

			_, filePath := u.paths(id, filepath.Base(filename))

			hook.Upload.Storage["Path"] = filePath

			return tusd.HTTPResponse{}, tusd.FileInfoChanges{
				ID:      id,
				Storage: hook.Upload.Storage,
			}, nil
		},
		PreFinishResponseCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
			filename, ok := hook.Upload.MetaData["filename"]
			if !ok || filename == "" {
				filename = hook.Upload.ID
			}

			_, err := u.queries.InsertFile(hook.Context, sqlc.InsertFileParams{
				ID:      hook.Upload.ID,
				Name:    filename,
				Blob:    path.Base(hook.Upload.Storage["Path"]),
				Size:    float64(hook.Upload.Size),
				Ticket:  hook.HTTPRequest.Header.Get("X-Ticket-ID"),
				Created: time.Now().UTC(),
				Updated: time.Now().UTC(),
			})

			return tusd.HTTPResponse{}, err
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tusd handler: %w", err)
	}

	// Start another goroutine for receiving events from the handler whenever
	// an upload is completed. The event will contains details about the upload
	// itself and the relevant HTTP request.
	go func() {
		for {
			event := <-handler.CompleteUploads
			slog.Info("Upload %s finished", "id", event.Upload.ID)
		}
	}()

	return chi.Chain(u.auth.Middleware, auth.ValidateFileScopes).Handler(handler), nil
}

type InfoFileMetaData struct {
	Filename     string `json:"filename"`
	Filetype     string `json:"filetype"`
	Name         string `json:"name"`
	RelativePath string `json:"relativePath"`
	Type         string `json:"type"`
}

type InfoFileStorage struct {
	InfoPath string `json:"InfoPath"`
	Path     string `json:"Path"`
	Type     string `json:"Type"`
}
type InfoFile struct {
	ID             string           `json:"ID"`
	Size           int              `json:"Size"`
	SizeIsDeferred bool             `json:"SizeIsDeferred"`
	Offset         int              `json:"Offset"`
	MetaData       InfoFileMetaData `json:"MetaData"`
	IsPartial      bool             `json:"IsPartial"`
	IsFinal        bool             `json:"IsFinal"`
	PartialUploads interface{}      `json:"PartialUploads"`
	Storage        InfoFileStorage  `json:"Storage"`
}

func (u *Uploader) CreateFile(id string, filename string, blob []byte) (string, error) {
	filename = filepath.Base(filename)

	infoFilePath, filePath := u.paths(id, filename)

	fileType := http.DetectContentType(blob)

	infoFileData := InfoFile{
		ID:             id,
		Size:           len(blob),
		SizeIsDeferred: true,
		Offset:         0,
		MetaData: InfoFileMetaData{
			Filename:     filename,
			Filetype:     fileType,
			Name:         filename,
			RelativePath: "null",
			Type:         fileType,
		},
		IsPartial:      false,
		IsFinal:        false,
		PartialUploads: nil,
		Storage: InfoFileStorage{
			InfoPath: infoFilePath,
			Path:     filePath,
			Type:     "filestore",
		},
	}

	if err := u.root.Mkdir(id, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory for file %s: %w", id, err)
	}

	file, err := u.root.Create(infoFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file info %s: %w", infoFilePath, err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(infoFileData); err != nil {
		return "", fmt.Errorf("failed to encode file info %s: %w", infoFilePath, err)
	}

	file, err = u.root.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	if _, err := file.Write(blob); err != nil {
		return "", fmt.Errorf("failed to write blob to file %s: %w", filePath, err)
	}

	return path.Base(filePath), nil
}

func (u *Uploader) File(id, name string) (*os.File, string, int64, error) {
	infoFilePath := id + ".info"

	infoFile, err := u.root.Open(infoFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", 0, fmt.Errorf("file info %s does not exist", infoFilePath)
		}

		return nil, "", 0, fmt.Errorf("failed to open file info %s: %w", infoFilePath, err)
	}
	defer infoFile.Close()

	var infoFileData InfoFile
	if err := json.NewDecoder(infoFile).Decode(&infoFileData); err != nil {
		return nil, "", 0, fmt.Errorf("failed to decode file info %s: %w", infoFilePath, err)
	}

	filePath := path.Join(id, name)

	info, err := u.root.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, "", 0, fmt.Errorf("file %s does not exist", filePath)
	}

	f, err := u.root.Open(filePath)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	return f, infoFileData.MetaData.Filetype, info.Size(), nil
}

func (u *Uploader) DeleteFile(id, name string) error {
	return errors.Join(
		u.root.Remove(path.Join(id, name)),
		u.root.Remove(id),
		u.root.Remove(id+".info"),
	)
}

func (u *Uploader) paths(id string, filename string) (infoFilePath, filePath string) {
	infoFilePath = id + ".info"
	ext := path.Ext(filename)
	prefix := strings.TrimSuffix(filename, ext)
	uniq := database.GenerateID("")
	filePath = path.Join(id, fmt.Sprintf("%s_%s%s", prefix, uniq, ext))

	return infoFilePath, filePath
}
