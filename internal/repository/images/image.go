package images

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

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
	newPath := filepath.Join(r.UploadPath, header.Filename)
	if err := os.MkdirAll(r.UploadPath, os.ModePerm); err != nil {
		return "", err
	}

	newFile, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return "", err
	}
	return newPath, nil
}

func (r *ImageDB) DeleteImage(imagePath string) error {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return err
	}

	if err := os.Remove(imagePath); err != nil {
		return err
	}

	return nil
}
