package storage

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"hrm-app/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	repo   StorageRepository
	bucket string
}

func NewHandler(repo StorageRepository, bucket string) *Handler {
	return &Handler{
		repo:   repo,
		bucket: bucket,
	}
}

func (h *Handler) UploadFile(c *gin.Context) {
	// 1. Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 2. Parse multipart form (max 10MB)
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to parse form: "+err.Error())
		return
	}

	// 3. Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	// 4. Generate unique file key
	// Format: users/{userID}/{uuid}{ext}
	ext := filepath.Ext(header.Filename)
	if ext == "" && len(header.Filename) > 0 {
		if idx := strings.LastIndex(header.Filename, "."); idx != -1 {
			ext = header.Filename[idx:]
		}
	}

	// Fallback/Default extension if still empty
	if ext == "" {
		ext = ".bin"
	}

	key := fmt.Sprintf("users/%d/%s%s", userID.(uint), uuid.New().String(), ext)

	// 5. Get content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 6. Upload to storage
	ctx := c.Request.Context()
	err = h.repo.Upload(ctx, h.bucket, key, file, contentType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload file: "+err.Error())
		return
	}

	// 7. Get Public URL
	url := h.repo.GetURL(h.bucket, key)

	response.Success(c, gin.H{
		"url": url,
	})
}
