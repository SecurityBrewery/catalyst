package catalyst

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"

	"github.com/SecurityBrewery/catalyst/storage"
)

func upload(client *s3.S3, external string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ticketID, exists := ctx.Params.Get("ticketID")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ticketID not given"})
			return
		}

		if err := storage.CreateBucket(client, ticketID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("could not create bucket: %w", err)})
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
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("could not create tusd handler: %w", err)})
			return
		}

		switch ctx.Request.Method {
		case http.MethodHead:
			gin.WrapF(handler.HeadFile)(ctx)
		case http.MethodPost:
			gin.WrapF(handler.PostFile)(ctx)
		case http.MethodPatch:
			gin.WrapF(handler.PatchFile)(ctx)
		default:
			log.Println(errors.New("unknown method"))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unknown method"})
		}
	}
}

func download(downloader *s3manager.Downloader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ticketID, exists := ctx.Params.Get("ticketID")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ticketID not given"})
			return
		}

		key, exists := ctx.Params.Get("key")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "key not given"})
			return
		}

		buf := sequentialWriter{ctx.Writer}

		downloader.Concurrency = 1
		_, err := downloader.Download(buf, &s3.GetObjectInput{
			Bucket: aws.String("catalyst-" + ticketID),
			Key:    aws.String(key),
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}
}

type sequentialWriter struct {
	w io.Writer
}

func (fw sequentialWriter) WriteAt(p []byte, _ int64) (n int, err error) {
	return fw.w.Write(p)
}
