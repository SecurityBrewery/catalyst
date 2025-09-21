package rootstore

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tus/tusd/v2/pkg/handler"
)

// Test interface implementation of FSStore.
var (
	_ handler.DataStore               = RootStore{}
	_ handler.TerminaterDataStore     = RootStore{}
	_ handler.ConcaterDataStore       = RootStore{}
	_ handler.LengthDeferrerDataStore = RootStore{}
)

func TestFSStore(t *testing.T) {
	t.Parallel()

	root, err := os.OpenRoot(t.TempDir())
	require.NoError(t, err)

	t.Cleanup(func() { root.Close() })

	store := New(root)
	ctx := t.Context()

	// Create new upload
	upload, err := store.NewUpload(ctx, handler.FileInfo{
		Size: 42,
		MetaData: map[string]string{
			"hello": "world",
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, upload)

	// Check info without writing
	info, err := upload.GetInfo(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 42, info.Size)
	assert.EqualValues(t, 0, info.Offset)
	assert.Equal(t, handler.MetaData{"hello": "world"}, info.MetaData)
	assert.Len(t, info.Storage, 3)
	assert.Equal(t, "rootstore", info.Storage["Type"])
	assert.Equal(t, info.ID, info.Storage["Path"])
	assert.Equal(t, info.ID+".info", info.Storage["InfoPath"])

	// Write data to upload
	bytesWritten, err := upload.WriteChunk(ctx, 0, strings.NewReader("hello world"))
	require.NoError(t, err)
	assert.EqualValues(t, len("hello world"), bytesWritten)

	// Check new offset
	info, err = upload.GetInfo(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 42, info.Size)
	assert.EqualValues(t, 11, info.Offset)

	// Read content
	reader, err := upload.GetReader(ctx)
	require.NoError(t, err)

	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(content))
	reader.Close()

	// Serve content
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Range", "bytes=0-4")

	err = store.AsServableUpload(upload).ServeContent(t.Context(), w, r)
	require.NoError(t, err)

	assert.Equal(t, http.StatusPartialContent, w.Code)
	assert.Equal(t, "5", w.Header().Get("Content-Length"))
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "bytes 0-4/11", w.Header().Get("Content-Range"))
	assert.NotEmpty(t, w.Header().Get("Last-Modified"))
	assert.Equal(t, "hello", w.Body.String())

	// Terminate upload
	require.NoError(t, store.AsTerminatableUpload(upload).Terminate(ctx))

	// Test if upload is deleted
	upload, err = store.GetUpload(ctx, info.ID)
	assert.Nil(t, upload)
	assert.Equal(t, handler.ErrNotFound, err)
}

// TestCreateDirectories tests whether an upload with a slash in its ID causes
// the correct directories to be created.
func TestFSStoreCreateDirectories(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	root, err := os.OpenRoot(tmp)
	require.NoError(t, err)

	t.Cleanup(func() { root.Close() })

	store := New(root)
	ctx := t.Context()

	// Create new upload
	upload, err := store.NewUpload(ctx, handler.FileInfo{
		ID:   "hello/world/123",
		Size: 42,
		MetaData: map[string]string{
			"hello": "world",
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, upload)

	// Check info without writing
	info, err := upload.GetInfo(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 42, info.Size)
	assert.EqualValues(t, 0, info.Offset)
	assert.Equal(t, handler.MetaData{"hello": "world"}, info.MetaData)
	assert.Len(t, info.Storage, 3)
	assert.Equal(t, "rootstore", info.Storage["Type"])
	assert.Equal(t, filepath.FromSlash(info.ID), info.Storage["Path"])
	assert.Equal(t, filepath.FromSlash(info.ID+".info"), info.Storage["InfoPath"])

	// Write data to upload
	bytesWritten, err := upload.WriteChunk(ctx, 0, strings.NewReader("hello world"))
	require.NoError(t, err)
	assert.EqualValues(t, len("hello world"), bytesWritten)

	// Check new offset
	info, err = upload.GetInfo(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 42, info.Size)
	assert.EqualValues(t, 11, info.Offset)

	// Read content
	reader, err := upload.GetReader(ctx)
	require.NoError(t, err)

	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(content))
	reader.Close()

	// Check that the file and directory exists on disk
	statInfo, err := os.Stat(filepath.Join(tmp, "hello/world/123"))
	require.NoError(t, err)
	assert.True(t, statInfo.Mode().IsRegular())
	assert.EqualValues(t, 11, statInfo.Size())
	statInfo, err = os.Stat(filepath.Join(tmp, "hello/world/"))
	require.NoError(t, err)
	assert.True(t, statInfo.Mode().IsDir())

	// Terminate upload
	require.NoError(t, store.AsTerminatableUpload(upload).Terminate(ctx))

	// Test if upload is deleted
	upload, err = store.GetUpload(ctx, info.ID)
	assert.Nil(t, upload)
	assert.Equal(t, handler.ErrNotFound, err)
}

func TestFSStoreNotFound(t *testing.T) {
	t.Parallel()

	root, err := os.OpenRoot(t.TempDir())
	require.NoError(t, err)

	t.Cleanup(func() { root.Close() })

	store := New(root)
	ctx := t.Context()

	upload, err := store.GetUpload(ctx, "upload-that-does-not-exist")
	require.Error(t, err)
	assert.Equal(t, handler.ErrNotFound, err)
	assert.Nil(t, upload)
}

func TestFSStoreConcatUploads(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	root, err := os.OpenRoot(tmp)
	require.NoError(t, err)

	t.Cleanup(func() { root.Close() })

	store := New(root)
	ctx := t.Context()

	// Create new upload to hold concatenated upload
	finUpload, err := store.NewUpload(ctx, handler.FileInfo{Size: 9})
	require.NoError(t, err)
	assert.NotNil(t, finUpload)

	finInfo, err := finUpload.GetInfo(ctx)
	require.NoError(t, err)

	finID := finInfo.ID

	// Create three uploads for concatenating
	partialUploads := make([]handler.Upload, 3)
	contents := []string{
		"abc",
		"def",
		"ghi",
	}

	for i := range 3 {
		upload, err := store.NewUpload(ctx, handler.FileInfo{Size: 3})
		require.NoError(t, err)

		n, err := upload.WriteChunk(ctx, 0, strings.NewReader(contents[i]))
		require.NoError(t, err)
		assert.EqualValues(t, 3, n)

		partialUploads[i] = upload
	}

	err = store.AsConcatableUpload(finUpload).ConcatUploads(ctx, partialUploads)
	require.NoError(t, err)

	// Check offset
	finUpload, err = store.GetUpload(ctx, finID)
	require.NoError(t, err)

	info, err := finUpload.GetInfo(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 9, info.Size)
	assert.EqualValues(t, 9, info.Offset)

	// Read content
	reader, err := finUpload.GetReader(ctx)
	require.NoError(t, err)

	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "abcdefghi", string(content))
	reader.Close()
}

func TestFSStoreDeclareLength(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	root, err := os.OpenRoot(tmp)
	require.NoError(t, err)

	t.Cleanup(func() { root.Close() })

	store := New(root)
	ctx := t.Context()

	upload, err := store.NewUpload(ctx, handler.FileInfo{
		Size:           0,
		SizeIsDeferred: true,
	})
	require.NoError(t, err)
	assert.NotNil(t, upload)

	info, err := upload.GetInfo(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 0, info.Size)
	assert.True(t, info.SizeIsDeferred)

	err = store.AsLengthDeclarableUpload(upload).DeclareLength(ctx, 100)
	require.NoError(t, err)

	updatedInfo, err := upload.GetInfo(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 100, updatedInfo.Size)
	assert.False(t, updatedInfo.SizeIsDeferred)
}

// TestCustomRelativePath tests whether the upload's destination can be customized
// relative to the storage directory.
func TestFSStoreCustomRelativePath(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	root, err := os.OpenRoot(tmp)
	require.NoError(t, err)

	t.Cleanup(func() { root.Close() })

	store := New(root)
	ctx := t.Context()

	// Create new upload
	upload, err := store.NewUpload(ctx, handler.FileInfo{
		ID:   "folder1/info",
		Size: 42,
		Storage: map[string]string{
			"Path": "./folder2/bin",
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, upload)

	// Check info without writing
	info, err := upload.GetInfo(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 42, info.Size)
	assert.EqualValues(t, 0, info.Offset)
	assert.Len(t, info.Storage, 3)
	assert.Equal(t, "rootstore", info.Storage["Type"])
	assert.Equal(t, filepath.FromSlash("./folder2/bin"), info.Storage["Path"])
	assert.Equal(t, filepath.FromSlash("folder1/info.info"), info.Storage["InfoPath"])

	// Write data to upload
	bytesWritten, err := upload.WriteChunk(ctx, 0, strings.NewReader("hello world"))
	require.NoError(t, err)
	assert.EqualValues(t, len("hello world"), bytesWritten)

	// Check new offset
	info, err = upload.GetInfo(ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 42, info.Size)
	assert.EqualValues(t, 11, info.Offset)

	// Read content
	reader, err := upload.GetReader(ctx)
	require.NoError(t, err)

	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(content))
	reader.Close()

	// Check that the output file and info file exist on disk
	statInfo, err := os.Stat(filepath.Join(tmp, "folder2/bin"))
	require.NoError(t, err)
	assert.True(t, statInfo.Mode().IsRegular())
	assert.EqualValues(t, 11, statInfo.Size())
	statInfo, err = os.Stat(filepath.Join(tmp, "folder1/info.info"))
	require.NoError(t, err)
	assert.True(t, statInfo.Mode().IsRegular())

	// Terminate upload
	require.NoError(t, store.AsTerminatableUpload(upload).Terminate(ctx))

	// Test if upload is deleted
	upload, err = store.GetUpload(ctx, info.ID)
	assert.Nil(t, upload)
	assert.Equal(t, handler.ErrNotFound, err)
}

// TestCustomAbsolutePath tests whether the upload's destination can be customized
// using an absolute path to the storage directory.
func TestFSStoreCustomAbsolutePath(t *testing.T) {
	t.Parallel()

	root, err := os.OpenRoot(t.TempDir())
	require.NoError(t, err)

	t.Cleanup(func() { root.Close() })

	store := New(root)

	// Create new upload, but the Path property points to a directory
	// outside of the directory given to FSStore
	binPath := filepath.Join(t.TempDir(), "dir/my-upload.bin")
	_, err = store.NewUpload(t.Context(), handler.FileInfo{
		ID:   "my-upload",
		Size: 42,
		Storage: map[string]string{
			"Path": binPath,
		},
	})
	require.Error(t, err)

	_, err = os.Stat(binPath)
	require.Error(t, err)
}
