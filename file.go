package catalyst

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/arangodb/go-driver"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	maut "github.com/cugu/maut/auth"
	"github.com/go-chi/chi/v5"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/storage"
)

func tusdUpload(db *database.Database, catalystBus *bus.Bus, client *s3.S3, external string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticketID := chi.URLParam(r, "ticketID")
		if ticketID == "" {
			api.JSONErrorStatus(w, http.StatusBadRequest, errors.New("ticketID not given"))

			return
		}

		if err := storage.CreateBucket(client, ticketID); err != nil {
			api.JSONErrorStatus(w, http.StatusBadRequest, fmt.Errorf("could not create bucket: %w", err))

			return
		}

		store := s3store.New("catalyst-"+ticketID, client)

		composer := tusd.NewStoreComposer()
		store.UseIn(composer)

		handler, err := tusd.NewUnroutedHandler(tusd.Config{
			BasePath:              external + "/api/files/" + ticketID + "/tusd/",
			StoreComposer:         composer,
			NotifyCompleteUploads: true,
		})
		if err != nil {
			api.JSONErrorStatus(w, http.StatusBadRequest, fmt.Errorf("could not create tusd handler: %w", err))

			return
		}

		userID := "unknown"
		user, _, ok := maut.UserFromContext(r.Context())
		if ok {
			userID = user.ID
		}

		go func() {
			event := <-handler.CompleteUploads

			id, err := strconv.ParseInt(ticketID, 10, 64)
			if err != nil {
				return
			}

			file := &model.File{Key: event.Upload.Storage["Key"], Name: event.Upload.MetaData["filename"]}

			ctx := context.Background()
			doc, err := db.AddFile(ctx, id, file)
			if err != nil {
				log.Println(err)

				return
			}

			catalystBus.RequestChannel.Publish(&bus.RequestMsg{
				User:     userID,
				Function: "LinkFiles",
				IDs:      []driver.DocumentID{driver.DocumentID(fmt.Sprintf("tickets/%d", doc.ID))},
			})
		}()

		switch r.Method {
		case http.MethodHead:
			handler.HeadFile(w, r)
		case http.MethodPost:
			handler.PostFile(w, r)
		case http.MethodPatch:
			handler.PatchFile(w, r)
		default:
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("unknown method"))
		}
	}
}

func upload(db *database.Database, client *s3.S3, uploader *s3manager.Uploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticketID := chi.URLParam(r, "ticketID")
		if ticketID == "" {
			api.JSONErrorStatus(w, http.StatusBadRequest, errors.New("ticketID not given"))

			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			api.JSONErrorStatus(w, http.StatusBadRequest, err)

			return
		}
		defer file.Close()

		if err := storage.CreateBucket(client, ticketID); err != nil {
			api.JSONErrorStatus(w, http.StatusBadRequest, fmt.Errorf("could not create bucket: %w", err))

			return
		}

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("catalyst-" + ticketID),
			Key:    aws.String(header.Filename),
			Body:   file,
		})
		if err != nil {
			api.JSONErrorStatus(w, http.StatusBadRequest, err)

			return
		}

		id, err := strconv.ParseInt(ticketID, 10, 64)
		if err != nil {
			api.JSONErrorStatus(w, http.StatusBadRequest, err)

			return
		}

		_, err = db.AddFile(r.Context(), id, &model.File{
			Key:  header.Filename,
			Name: header.Filename,
		})
		if err != nil {
			api.JSONErrorStatus(w, http.StatusBadRequest, err)

			return
		}
	}
}

func download(downloader *s3manager.Downloader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticketID := chi.URLParam(r, "ticketID")
		if ticketID == "" {
			api.JSONErrorStatus(w, http.StatusBadRequest, errors.New("ticketID not given"))

			return
		}

		key := chi.URLParam(r, "key")
		if key == "" {
			api.JSONErrorStatus(w, http.StatusBadRequest, errors.New("key not given"))

			return
		}

		buf := sequentialWriter{w}

		downloader.Concurrency = 1
		_, err := downloader.Download(buf, &s3.GetObjectInput{
			Bucket: aws.String("catalyst-" + ticketID),
			Key:    aws.String(key),
		})
		if err != nil {
			api.JSONErrorStatus(w, http.StatusInternalServerError, err)
		}
	}
}

type sequentialWriter struct {
	w io.Writer
}

func (fw sequentialWriter) WriteAt(p []byte, _ int64) (n int, err error) {
	return fw.w.Write(p)
}
