package test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/SecurityBrewery/catalyst"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/pointer"
	ctime "github.com/SecurityBrewery/catalyst/time"
)

type testClock struct{}

func (testClock) Now() time.Time {
	return time.Date(2021, 12, 12, 12, 12, 12, 12, time.UTC)
}

func TestServer(t *testing.T) {
	ctime.DefaultClock = testClock{}

	for _, tt := range api.Tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, _, _, _, _, db, _, server, cleanup, err := Server(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			if err := SetupTestData(ctx, db); err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			// setup request
			var req *http.Request
			if tt.Args.Data != nil {
				b, err := json.Marshal(tt.Args.Data)
				if err != nil {
					t.Fatal(err)
				}

				req = httptest.NewRequest(strings.ToUpper(tt.Args.Method), tt.Args.URL, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(strings.ToUpper(tt.Args.Method), tt.Args.URL, nil)
			}

			// run request
			server.ServeHTTP(w, req)

			result := w.Result()

			// assert results
			if result.StatusCode != tt.Want.Status {
				msg, _ := io.ReadAll(result.Body)

				t.Fatalf("Status got = %v, want %v: %s", result.Status, tt.Want.Status, msg)
			}
			if tt.Want.Status != http.StatusNoContent {
				jsonEqual(t, result.Body, tt.Want.Body)
			}
		})
	}
}

func TestService(t *testing.T) {
	type args struct {
		method string
		url    string
		data   interface{}
	}
	type want struct {
		status int
		body   interface{}
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "GetUser not existing", args: args{method: http.MethodGet, url: "/users/123"}, want: want{status: http.StatusNotFound, body: map[string]string{"error": "document not found"}}},
		{name: "ListUsers", args: args{method: http.MethodGet, url: "/users"}, want: want{status: http.StatusOK}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, _, _, server, cleanup, err := Server(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			// server.ConfigureRoutes()
			w := httptest.NewRecorder()

			// setup request
			var req *http.Request
			if tt.args.data != nil {
				b, err := json.Marshal(tt.args.data)
				if err != nil {
					t.Fatal(err)
				}

				req = httptest.NewRequest(tt.args.method, tt.args.url, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.args.method, tt.args.url, nil)
			}

			// run request
			server.ServeHTTP(w, req)

			result := w.Result()

			// assert results
			if result.StatusCode != tt.want.status {
				t.Fatalf("Status got = %v, want %v", result.Status, tt.want.status)
			}
			if tt.want.status != http.StatusNoContent {
				jsonEqual(t, result.Body, tt.want.body)
			}
		})
	}
}

func TestBackupAndRestore(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	type want struct {
		status int
	}
	tests := []struct {
		name string
		want want
	}{
		{name: "Backup", want: want{status: http.StatusOK}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, server, err := Catalyst(t)
			if err != nil {
				t.Fatal(err)
			}

			if err := SetupTestData(ctx, server.DB); err != nil {
				t.Fatal(err)
			}

			createFile(ctx, server)

			zipB := assertBackup(t, server)

			assertZipFile(t, readZipFile(t, zipB))

			clearAllDatabases(server)
			_, err = server.DB.UserCreateSetupAPIKey(ctx, "test")
			if err != nil {
				log.Fatal(err)
			}

			deleteAllBuckets(t, server)

			assertRestore(t, zipB, server)

			assertTicketExists(t, server)

			assertFileExists(t, server)
		})
	}
}

func assertBackup(t *testing.T, server *catalyst.Server) []byte {
	// setup request
	req := httptest.NewRequest(http.MethodGet, "/api/backup/create", nil)
	req.Header.Set("PRIVATE-TOKEN", "test")

	// run request
	backupRequestRecorder := httptest.NewRecorder()
	server.Server.ServeHTTP(backupRequestRecorder, req)
	backupResult := backupRequestRecorder.Result()

	// assert results
	assert.Equal(t, http.StatusOK, backupResult.StatusCode)

	zipBuf := &bytes.Buffer{}
	if _, err := io.Copy(zipBuf, backupResult.Body); err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, backupResult.Body.Close())

	return zipBuf.Bytes()
}

func assertZipFile(t *testing.T, r *zip.Reader) {
	var names []string
	for _, f := range r.File {
		names = append(names, f.Name)
	}

	if !includes(t, names, "minio/catalyst-8125/test.txt") {
		t.Error("Minio file missing")
	}

	for _, p := range []string{
		"arango/ENCRYPTION", "arango/automations_.*.data.json.gz", "arango/automations_.*.structure.json", "arango/dump.json", "arango/jobs_.*.data.json.gz", "arango/jobs_.*.structure.json", "arango/logs_.*.data.json.gz", "arango/logs_.*.structure.json", "arango/migrations_.*.data.json.gz", "arango/migrations_.*.structure.json", "arango/playbooks_.*.data.json.gz", "arango/playbooks_.*.structure.json", "arango/related_.*.data.json.gz", "arango/related_.*.structure.json", "arango/templates_.*.data.json.gz", "arango/templates_.*.structure.json", "arango/tickets_.*.data.json.gz", "arango/tickets_.*.structure.json", "arango/tickettypes_.*.data.json.gz", "arango/tickettypes_.*.structure.json", "arango/userdata_.*.data.json.gz", "arango/userdata_.*.structure.json", "arango/users_.*.data.json.gz", "arango/users_.*.structure.json",
	} {
		if !includes(t, names, p) {
			t.Errorf("Arango file missing: %s", p)
		}
	}
}

func clearAllDatabases(server *catalyst.Server) {
	server.DB.Truncate(context.Background())
}

func deleteAllBuckets(t *testing.T, server *catalyst.Server) {
	buckets, err := server.Storage.S3().ListBuckets(&s3.ListBucketsInput{})
	for _, bucket := range buckets.Buckets {
		server.Storage.S3().DeleteBucket(&s3.DeleteBucketInput{
			Bucket: bucket.Name,
		})
	}

	if err != nil {
		t.Fatal(err)
	}
}

func assertRestore(t *testing.T, zipB []byte, server *catalyst.Server) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("backup", "backup.zip")
	if err != nil {
		log.Fatal(err)
	}

	_, err = fileWriter.Write(zipB)
	if err != nil {
		log.Fatal(err)
	}

	assert.NoError(t, bodyWriter.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/backup/restore", bodyBuf)
	req.Header.Set("PRIVATE-TOKEN", "test")
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	// run request
	restoreRequestRecorder := httptest.NewRecorder()
	server.Server.ServeHTTP(restoreRequestRecorder, req)
	restoreResult := restoreRequestRecorder.Result()

	if !assert.Equal(t, http.StatusOK, restoreResult.StatusCode) {
		b, _ := io.ReadAll(restoreResult.Body)
		log.Println(string(b))
		t.FailNow()
	}
}

func createFile(ctx context.Context, server *catalyst.Server) {
	buf := bytes.NewBufferString("test text")

	server.Storage.S3().CreateBucket(&s3.CreateBucketInput{Bucket: pointer.String("catalyst-8125")})

	if _, err := server.Storage.Uploader().Upload(&s3manager.UploadInput{Body: buf, Bucket: pointer.String("catalyst-8125"), Key: pointer.String("test.txt")}); err != nil {
		log.Fatal(err)
	}

	if _, err := server.DB.LinkFiles(ctx, 8125, []*model.File{{Key: "test.txt", Name: "test.txt"}}); err != nil {
		log.Fatal(err)
	}
}

func assertTicketExists(t *testing.T, server *catalyst.Server) {
	req := httptest.NewRequest(http.MethodGet, "/api/tickets/8125", nil)
	req.Header.Set("PRIVATE-TOKEN", "test")

	// run request
	backupRequestRecorder := httptest.NewRecorder()
	server.Server.ServeHTTP(backupRequestRecorder, req)
	backupResult := backupRequestRecorder.Result()

	// assert results
	assert.Equal(t, http.StatusOK, backupResult.StatusCode)

	zipBuf := &bytes.Buffer{}
	if _, err := io.Copy(zipBuf, backupResult.Body); err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, backupResult.Body.Close())

	var ticket model.Ticket
	assert.NoError(t, json.Unmarshal(zipBuf.Bytes(), &ticket))

	assert.Equal(t, "phishing from selenafadel@von.com detected", ticket.Name)
}

func assertFileExists(t *testing.T, server *catalyst.Server) {
	obj, err := server.Storage.S3().GetObject(&s3.GetObjectInput{
		Bucket: aws.String("catalyst-8125"),
		Key:    aws.String("test.txt"),
	})
	assert.NoError(t, err)

	b, err := io.ReadAll(obj.Body)
	assert.NoError(t, err)

	assert.Equal(t, "test text", string(b))
}

func includes(t *testing.T, names []string, s string) bool {
	for _, name := range names {
		match, err := regexp.MatchString(s, name)
		if err != nil {
			t.Fatal(err)
		}

		if match {
			return true
		}
	}
	return false
}

func readZipFile(t *testing.T, b []byte) *zip.Reader {
	buf := bytes.NewReader(b)

	zr, err := zip.NewReader(buf, int64(buf.Len()))
	if err != nil {
		t.Fatal(string(b), err)
	}

	return zr
}

func jsonEqual(t *testing.T, got io.Reader, want interface{}) {
	var gotObject, wantObject interface{}

	// load bytes
	wantBytes, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	gotBytes, err := io.ReadAll(got)
	if err != nil {
		t.Fatal(err)
	}

	fields := []string{"secret"}
	for _, field := range fields {
		gField := gjson.GetBytes(wantBytes, field)
		if gField.Exists() && gjson.GetBytes(gotBytes, field).Exists() {
			gotBytes, err = sjson.SetBytes(gotBytes, field, gField.Value())
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	// normalize bytes
	if err = json.Unmarshal(wantBytes, &wantObject); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(gotBytes, &gotObject); err != nil {
		t.Fatal(string(gotBytes), err)
	}

	// compare
	assert.Equal(t, wantObject, gotObject)
}
