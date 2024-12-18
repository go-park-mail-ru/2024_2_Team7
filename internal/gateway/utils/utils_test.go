package utils

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFilename(t *testing.T) {
	t.Run("Valid filename generation", func(t *testing.T) {
		header := &multipart.FileHeader{
			Filename: "test_image.jpg",
		}

		err := GenerateFilename(header)

		assert.NoError(t, err)
		assert.Contains(t, header.Filename, "_")

		assert.Equal(t, "jpg", getFileExtension(header.Filename))
	})

	t.Run("Invalid image format", func(t *testing.T) {
		header := &multipart.FileHeader{
			Filename: "test_document.pdf",
		}

		err := GenerateFilename(header)

		assert.Error(t, err)
		assert.Equal(t, models.ErrInvalidImageFormat, err)
	})
}

func TestGetPaginationParams(t *testing.T) {
	t.Run("Default pagination params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/items?page=0&limit=30", nil)

		params := GetPaginationParams(req)

		assert.Equal(t, 0, params.Offset)
		assert.Equal(t, 30, params.Limit)
	})

	t.Run("Custom pagination params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/items?page=1&limit=10", nil)

		params := GetPaginationParams(req)

		assert.Equal(t, 10, params.Offset)
		assert.Equal(t, 10, params.Limit)
	})
}
