// internal/controller/http/handlers/file_handler.go
package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/response"
	"github.com/google/uuid"
)

type FileHandler struct {
	storage    ports.FileStorage
	bucketName string
}

func NewFileHandler(storage ports.FileStorage, bucketName string) *FileHandler {
	return &FileHandler{
		storage:    storage,
		bucketName: bucketName,
	}
}

// POST /api/uploads/dishes
// form-data: file=<image>
func (h *FileHandler) UploadDishPhoto(w http.ResponseWriter, r *http.Request) {
	const maxSize = 5 << 20 // 5 MB

	if err := r.ParseMultipartForm(maxSize); err != nil {
		response.BadRequest(w, "failed to parse multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.BadRequest(w, "file is required")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Разрешаем только картинки
	if !strings.HasPrefix(contentType, "image/") {
		response.BadRequest(w, "only image uploads are allowed")
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	// Имя файла: dishes/<uuid>.<ext>
	filename := fmt.Sprintf("dishes/%s%s", uuid.New().String(), ext)

	url, err := h.storage.Upload(
		r.Context(),
		h.bucketName,
		filename,
		file,
		header.Size,
		contentType,
	)
	if err != nil {
		response.InternalError(w, "failed to upload file")
		return
	}

	// Возвращаем URL, который потом положим в dish.photo_url
	response.Created(w, map[string]string{
		"url": url,
	})
}
