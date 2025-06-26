package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

type Uploader struct {
	path    string
	queries *sqlc.Queries
	auth    *auth.Service
}

func NewUploader(p string, auth *auth.Service, queries *sqlc.Queries) *Uploader {
	return &Uploader{
		queries: queries,
		path:    path.Join(p, "uploads"),
		auth:    auth,
	}
}

func (u *Uploader) Routes() (http.Handler, error) {
	// Create a new FileStore instance which is responsible for
	// storing the uploaded file on disk in the specified directory.
	// This path _must_ exist before tusd will store uploads in it.
	// If you want to save them on a different medium, for example
	// a remote FTP server, you can implement your own storage backend
	// by implementing the tusd.DataStore interface.
	store := filestore.New(u.path)

	// A locking mechanism helps preventing data loss or corruption from
	// parallel requests to a upload resource. A good match for the disk-based
	// storage is the filelocker package which uses disk-based file lock for
	// coordinating access.
	// More information is available at https://tus.github.io/tusd/advanced-topics/locks/.
	locker := filelocker.New(u.path)

	// A storage backend for tusd may consist of multiple different parts which
	// handle upload creation, locking, termination and so on. The composer is a
	// place where all those separated pieces are joined together. In this example
	// we only use the file store but you may plug in multiple.
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

			ext := path.Ext(filename)
			prefix := strings.TrimSuffix(filename, ext)
			uniq := database.GenerateID("")

			hook.Upload.Storage["Path"] = path.Join(id, fmt.Sprintf("%s_%s%s", prefix, uniq, ext))

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

type InfoFile struct {
	ID             string `json:"ID"`
	Size           int    `json:"Size"`
	SizeIsDeferred bool   `json:"SizeIsDeferred"`
	Offset         int    `json:"Offset"`
	MetaData       struct {
		Filename     string `json:"filename"`
		Filetype     string `json:"filetype"`
		Name         string `json:"name"`
		RelativePath string `json:"relativePath"`
		Type         string `json:"type"`
	} `json:"MetaData"`
	IsPartial      bool        `json:"IsPartial"`
	IsFinal        bool        `json:"IsFinal"`
	PartialUploads interface{} `json:"PartialUploads"`
	Storage        struct {
		InfoPath string `json:"InfoPath"`
		Path     string `json:"Path"`
		Type     string `json:"Type"`
	} `json:"Storage"`
}

func (u *Uploader) File(id string, blob string) (*os.File, string, int64, error) {
	infoFilePath := path.Join(u.path, id+".info")

	infoFile, err := os.Open(infoFilePath)
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

	filePath := path.Join(u.path, id, blob)

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, "", 0, fmt.Errorf("file %s does not exist", filePath)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	return f, infoFileData.MetaData.Filetype, info.Size(), nil
}

func (u *Uploader) DeleteFile(id string) error {
	return errors.Join(
		os.RemoveAll(path.Join(u.path, id)),
		os.RemoveAll(id+".info"),
	)
}
