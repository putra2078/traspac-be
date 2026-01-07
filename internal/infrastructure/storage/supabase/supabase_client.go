package supabase

import (
	"context"
	"fmt"
	"hrm-app/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// SupabaseS3Client menyimpan instance S3 client
type SupabaseS3Client struct {
	Client *s3.S3
}

// NewSupabaseS3Client inisialisasi S3 client Supabase
func NewSupabaseS3Client(cfg *config.Config) (*SupabaseS3Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Supabase.S3.Region),
		Endpoint:         aws.String(cfg.Supabase.S3.Endpoint),
		Credentials:      credentials.NewStaticCredentials(cfg.Supabase.S3.AccessKeyID, cfg.Supabase.S3.SecretAccessKey, ""),
		S3ForcePathStyle: aws.Bool(true), // Wajib untuk Supabase S3
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := s3.New(sess)
	return &SupabaseS3Client{
		Client: client,
	}, nil
}

// ListBuckets untuk test connection
func (c *SupabaseS3Client) ListBuckets(ctx context.Context) error {
	output, err := c.Client.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("failed to list buckets: %w", err)
	}

	for _, b := range output.Buckets {
		fmt.Println("Bucket:", aws.StringValue(b.Name))
	}
	return nil
}
