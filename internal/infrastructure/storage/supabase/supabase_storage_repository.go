package supabase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"hrm-app/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type SupabaseStorageRepository struct {
	client *SupabaseS3Client
	config *config.Config
}

func NewSupabaseStorageRepository(cfg *config.Config) (*SupabaseStorageRepository, error) {
	// Initialize S3 client using SupabaseS3Client
	client, err := NewSupabaseS3Client(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase S3 client: %w", err)
	}

	return &SupabaseStorageRepository{
		client: client,
		config: cfg,
	}, nil
}

func (s *SupabaseStorageRepository) Upload(ctx context.Context, bucket string, key string, data io.Reader, contentType string) error {
	// Read all data from the reader
	body, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	// Upload to S3 using the client
	_, err = s.client.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"), // Make file publicly accessible
	})

	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}

func (s *SupabaseStorageRepository) Delete(ctx context.Context, bucket string, key string) error {
	_, err := s.client.Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

func (s *SupabaseStorageRepository) GetURL(bucket string, key string) string {
	// Construct the public URL for Supabase storage
	// Standard Supabase API format: https://PROJECT_ID.supabase.co/storage/v1/object/public/bucket/key
	// S3 Config Endpoint format: https://PROJECT_ID.storage.supabase.co/storage/v1/s3

	u, err := url.Parse(s.config.Supabase.S3.Endpoint)
	if err != nil {
		return "" // or handle error gracefully
	}

	// 1. Clean up Host: remove ".storage" if present to switch from S3 endpoint to main API
	host := strings.Replace(u.Host, ".storage.supabase.co", ".supabase.co", 1)

	// 2. Construct the standard API URL
	// Note: The path should be /storage/v1/object/public/<bucket>/<key>
	return fmt.Sprintf("%s://%s/storage/v1/object/public/%s/%s", u.Scheme, host, bucket, key)
}
