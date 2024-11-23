package images

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"kudago/internal/models"
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

func (r *ImageDB) UploadImage(ctx context.Context, media models.MediaFile) (string, error) {
	defer media.File.Close()

	buffer := make([]byte, 512)
	if _, err := media.File.Read(buffer); err != nil {
		return "", fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	fileType := http.DetectContentType(buffer)
	if !isSupportedImageType(fileType) {
		return "", fmt.Errorf("%s: %w", models.LevelDB, models.ErrUnsupportedFile)
	}

	if _, err := media.File.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	newPath := filepath.Join(r.UploadPath, media.Filename)
	if err := os.MkdirAll(r.UploadPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("%s: %w", models.LevelDB, models.ErrUnsupportedFile)
	}

	newFile, err := os.Create(newPath)
	if err != nil {
		return "", fmt.Errorf("%s: %w", models.LevelDB, models.ErrUnsupportedFile)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, media.File)
	if err != nil {
		return "", fmt.Errorf("%s: %w", models.LevelDB, models.ErrUnsupportedFile)
	}
	return newPath, nil
}

func (r *ImageDB) DeleteImage(ctx context.Context, imagePath string) error {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("%s: %w", models.LevelDB, models.ErrNotFound)
	}

	if err := os.Remove(imagePath); err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, models.ErrUnsupportedFile)
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
