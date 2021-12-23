package catalyst

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/storage"
)

func BackupHandler(catalystStorage *storage.Storage, c *database.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Content-Disposition", "attachment; filename=backup.zip")
		context.Header("Content-Type", "application/zip")
		err := Backup(catalystStorage, c, context.Writer)
		if err != nil {
			log.Println(err)
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

type WriterAtBuffer struct {
	bytes.Buffer
}

func (fw WriterAtBuffer) WriteAt(p []byte, offset int64) (n int, err error) {
	return fw.Write(p)
}

func Backup(catalystStorage *storage.Storage, c *database.Config, writer io.Writer) error {
	archive := zip.NewWriter(writer)
	defer archive.Close()

	archive.SetComment(GetVersion())

	// S3
	if err := backupS3(catalystStorage, archive); err != nil {
		return err
	}

	// Arango
	return backupArango(c, archive)
}

func backupS3(catalystStorage *storage.Storage, archive *zip.Writer) error {
	buckets, err := catalystStorage.S3().ListBuckets(nil)
	if err != nil {
		return err
	}

	for _, bucket := range buckets.Buckets {
		objects, err := catalystStorage.S3().ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: bucket.Name,
		})
		if err != nil {
			return err
		}

		for _, content := range objects.Contents {
			rbuf := &WriterAtBuffer{}
			_, err := catalystStorage.Downloader().Download(rbuf, &s3.GetObjectInput{
				Bucket: bucket.Name,
				Key:    content.Key,
			})
			if err != nil {
				return err
			}

			a, err := archive.Create(path.Join("minio", *bucket.Name, *content.Key))
			if err != nil {
				return err
			}

			if _, err := io.Copy(a, rbuf); err != nil {
				return err
			}
		}
	}
	return nil
}

func backupArango(c *database.Config, archive *zip.Writer) error {
	dir, err := os.MkdirTemp("", "catalyst-backup")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err := arangodump(dir, c); err != nil {
		return err
	}

	return zipDump(dir, archive)
}

func zipDump(dir string, archive *zip.Writer) error {
	fsys := os.DirFS(dir)
	return fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		a, err := archive.Create(path.Join("arango", p))
		if err != nil {
			return err
		}

		f, err := fsys.Open(p)
		if err != nil {
			return err
		}

		if _, err := io.Copy(a, f); err != nil {
			return err
		}
		return nil
	})
}

func arangodump(dir string, config *database.Config) error {
	host := strings.Replace(config.Host, "http", "tcp", 1)

	name := config.Name
	if config.Name == "" {
		name = database.Name
	}
	args := []string{
		"--output-directory", dir, "--server.endpoint", host,
		"--server.username", config.User, "--server.password", config.Password,
		"--server.database", name,
	}
	cmd := exec.Command("arangodump", args...)
	return cmd.Run()
}
