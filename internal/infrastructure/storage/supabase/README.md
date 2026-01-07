# Supabase S3 Storage Implementation

This package provides an implementation of the `StorageRepository` interface using Supabase's S3-compatible storage.

## Configuration

Add the following configuration to your `config/config.yaml`:

```yaml
supabase:
  s3:
    endpoint: "https://your-project-ref.supabase.co/storage/v1/s3"
    region: "us-east-1"
    access_key_id: "your_supabase_access_key_id"
    secret_access_key: "your_supabase_secret_access_key"
    bucket: "your_bucket_name"
```

### Getting Supabase S3 Credentials

1. Go to your Supabase project dashboard
2. Navigate to **Settings** → **API**
3. Find your **S3 Access Keys** section
4. Copy the `access_key_id` and `secret_access_key`
5. The endpoint format is: `https://<project-ref>.supabase.co/storage/v1/s3`

## Usage

### Initialize the Storage Client

```go
package main

import (
    "hrm-app/config"
    "hrm-app/internal/infrastructure/storage/supabase"
)

func main() {
    cfg := config.LoadConfig()
    
    storage, err := supabase.NewSupabaseStorage(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use storage...
}
```

### Upload a File

```go
import (
    "context"
    "os"
)

func uploadExample(storage *supabase.SupabaseStorage) error {
    file, err := os.Open("path/to/file.jpg")
    if err != nil {
        return err
    }
    defer file.Close()
    
    ctx := context.Background()
    bucket := "avatars"
    key := "user/123/profile.jpg"
    contentType := "image/jpeg"
    
    err = storage.Upload(ctx, bucket, key, file, contentType)
    if err != nil {
        return err
    }
    
    // Get the public URL
    url := storage.GetURL(bucket, key)
    fmt.Println("File uploaded to:", url)
    
    return nil
}
```

### Delete a File

```go
func deleteExample(storage *supabase.SupabaseStorage) error {
    ctx := context.Background()
    bucket := "avatars"
    key := "user/123/profile.jpg"
    
    err := storage.Delete(ctx, bucket, key)
    if err != nil {
        return err
    }
    
    fmt.Println("File deleted successfully")
    return nil
}
```

### Get Public URL

```go
func getURLExample(storage *supabase.SupabaseStorage) {
    bucket := "avatars"
    key := "user/123/profile.jpg"
    
    url := storage.GetURL(bucket, key)
    fmt.Println("Public URL:", url)
}
```

## Features

- **Upload**: Upload files with custom content types and public read access
- **Delete**: Remove files from storage
- **GetURL**: Generate public URLs for accessing uploaded files

## Notes

- Files are uploaded with `public-read` ACL by default
- The implementation uses AWS SDK v1 (note: AWS SDK Go v1 is deprecated, consider migrating to v2 in the future)
- All operations support context for timeout and cancellation
