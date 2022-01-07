package catalyst

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-chi/chi"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"

	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/storage"
)

func upload(client *s3.S3, external string) http.HandlerFunc {
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
			BasePath:      external + "/api/files/" + ticketID + "/upload/",
			StoreComposer: composer,
		})
		if err != nil {
			api.JSONErrorStatus(w, http.StatusBadRequest, fmt.Errorf("could not create tusd handler: %w", err))
			return
		}

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
