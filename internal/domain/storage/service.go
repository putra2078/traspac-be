package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, bucket string, folder string) (string, error)
}

type service struct {
	repo StorageRepository
}

func NewService(repo StorageRepository) Service {
	return &service{repo: repo}
}

func (s *service) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, bucket string, folder string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Determine extension
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" && len(fileHeader.Filename) > 0 {
		if idx := strings.LastIndex(fileHeader.Filename, "."); idx != -1 {
			ext = fileHeader.Filename[idx:]
		}
	}
	if ext == "" {
		ext = ".jpg" // Default
	}

	// Generate key: {folder}/{uuid}{ext}
	// Ensure folder doesn't have trailing slash for clean join, or just use Sprintf
	folder = strings.TrimSuffix(folder, "/")
	key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Content Type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	if err := s.repo.Upload(ctx, bucket, key, file, contentType); err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return s.repo.GetURL(bucket, key), nil
}
