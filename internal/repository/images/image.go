package images

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"kudago/internal/models"

	"github.com/pkg/errors"

	"golang.org/x/net/context"
)

type ImageDB struct {
	UploadPath string
}

type ImageConfig struct {
	Path string
}

func NewDB(config ImageConfig) *ImageDB {
	return &ImageDB{UploadPath: config.Path}
}

func (r *ImageDB) SaveImage(ctx context.Context, header multipart.FileHeader, file multipart.File) (string, error) {
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", errors.Wrap(err, models.LevelDB)
	}

	fileType := http.DetectContentType(buffer)
	if !isSupportedImageType(fileType) {
		return "", errors.Wrap(models.ErrUnsupportedFile, models.LevelDB)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", errors.Wrap(err, models.LevelDB)
	}

	newPath := filepath.Join(r.UploadPath, header.Filename)
	if err := os.MkdirAll(r.UploadPath, os.ModePerm); err != nil {
		return "", errors.Wrap(err, models.LevelDB)
	}

	newFile, err := os.Create(newPath)
	if err != nil {
		return "", errors.Wrap(err, models.LevelDB)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return "", errors.Wrap(err, models.LevelDB)
	}
	return newPath, nil
}

func (r *ImageDB) DeleteImage(ctx context.Context, imagePath string) error {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return errors.Wrap(err, models.LevelDB)
	}

	if err := os.Remove(imagePath); err != nil {
		return errors.Wrap(err, models.LevelDB)
	}

	return nil
}

func isSupportedImageType(fileType string) bool {
	supportedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}

	return supportedTypes[fileType]
}
