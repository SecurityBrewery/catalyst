package upload

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Uploader struct {
	Root *os.Root
}

func New(dir string) (*Uploader, error) {
	uploadsDir := path.Join(dir, "uploads")

	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create uploads directory: %w", err)
	}

	root, err := os.OpenRoot(uploadsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open uploads directory: %w", err)
	}

	return &Uploader{
		Root: root,
	}, nil
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

	infoFilePath, filePath := u.Paths(id, filename)

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

	if err := u.Root.Mkdir(id, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory for file %s: %w", id, err)
	}

	file, err := u.Root.Create(infoFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file info %s: %w", infoFilePath, err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(infoFileData); err != nil {
		return "", fmt.Errorf("failed to encode file info %s: %w", infoFilePath, err)
	}

	file, err = u.Root.Create(filePath)
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

	infoFile, err := u.Root.Open(infoFilePath)
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

	info, err := u.Root.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, "", 0, fmt.Errorf("file %s does not exist", filePath)
	}

	f, err := u.Root.Open(filePath)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	return f, infoFileData.MetaData.Filetype, info.Size(), nil
}

func (u *Uploader) DeleteFile(id, name string) error {
	return errors.Join(
		u.Root.Remove(path.Join(id, name)),
		u.Root.Remove(id),
		u.Root.Remove(id+".info"),
	)
}

func (u *Uploader) Paths(id string, filename string) (infoFilePath, filePath string) {
	infoFilePath = id + ".info"
	ext := path.Ext(filename)
	prefix := strings.TrimSuffix(filename, ext)
	uniq := rand.Text()
	filePath = path.Join(id, fmt.Sprintf("%s_%s%s", prefix, uniq, ext))

	return infoFilePath, filePath
}
