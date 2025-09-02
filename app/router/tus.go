package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tus/tusd/v2/pkg/filelocker"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/rootstore"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/upload"
)

func tusRoutes(queries *sqlc.Queries, u *upload.Uploader) (http.Handler, error) {
	store := rootstore.New(u.Root)
	locker := filelocker.New(u.Root.Name())
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

			_, filePath := u.Paths(id, filepath.Base(filename))

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

			_, err := queries.InsertFile(hook.Context, sqlc.InsertFileParams{
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

	return chi.Chain(auth.Middleware(queries), auth.ValidateFileScopes).Handler(handler), nil
}
