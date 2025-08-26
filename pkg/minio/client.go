package minio

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"huddle/internal/config"
	"huddle/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type Client struct {
	client     *minio.Client
	bucketName string
	endpoint   string
}

var minioClient *Client

// NewClient creates a new MinIO client
func NewClient() (*Client, error) {
	cfg := config.GetConfig()
	
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKeyID, cfg.MinIO.SecretAccessKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	client := &Client{
		client:     minioClient,
		bucketName: cfg.MinIO.BucketName,
		endpoint:   cfg.MinIO.Endpoint,
	}

	// Ensure bucket exists
	if err := client.ensureBucketExists(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	logger.Info("MinIO client initialized successfully",
		zap.String("endpoint", cfg.MinIO.Endpoint),
		zap.String("bucket", cfg.MinIO.BucketName))

	return client, nil
}

// GetClient returns the singleton MinIO client
func GetClient() *Client {
	return minioClient
}

// SetClient sets the singleton MinIO client
func SetClient(client *Client) {
	minioClient = client
}

// ensureBucketExists ensures the bucket exists, creates it if it doesn't
func (c *Client) ensureBucketExists() error {
	ctx := context.Background()
	
	exists, err := c.client.BucketExists(ctx, c.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = c.client.MakeBucket(ctx, c.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		
		logger.Info("Created MinIO bucket", zap.String("bucket", c.bucketName))
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (c *Client) UploadFile(ctx context.Context, objectKey string, file multipart.File, header *multipart.FileHeader) error {
	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	// Upload file
	_, err := c.client.PutObject(ctx, c.bucketName, objectKey, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
		UserMetadata: map[string]string{
			"original-name": header.Filename,
			"upload-time":   time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	logger.Info("File uploaded to MinIO successfully",
		zap.String("object_key", objectKey),
		zap.String("filename", header.Filename),
		zap.Int64("size", header.Size))

	return nil
}

// DownloadFile downloads a file from MinIO
func (c *Client) DownloadFile(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	obj, err := c.client.GetObject(ctx, c.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}

	return obj, nil
}

// DeleteFile deletes a file from MinIO
func (c *Client) DeleteFile(ctx context.Context, objectKey string) error {
	err := c.client.RemoveObject(ctx, c.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	logger.Info("File deleted from MinIO", zap.String("object_key", objectKey))
	return nil
}

// GetPresignedURL generates a presigned URL for file access
func (c *Client) GetPresignedURL(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, c.bucketName, objectKey, expires, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// GetFileInfo gets file information from MinIO
func (c *Client) GetFileInfo(ctx context.Context, objectKey string) (*minio.ObjectInfo, error) {
	info, err := c.client.StatObject(ctx, c.bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info from MinIO: %w", err)
	}

	return &info, nil
}

// GenerateObjectKey generates a unique object key for file storage
func (c *Client) GenerateObjectKey(userID uint, originalName string) string {
	// Get file extension
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".bin"
	}
	
	// Generate unique filename
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%d%s", userID, timestamp, ext)
	
	// Create object key with user folder structure
	objectKey := fmt.Sprintf("uploads/%d/%s", userID, filename)
	
	return objectKey
}

// GetFileType determines file type from MIME type
func (c *Client) GetFileType(mimeType string) string {
	mimeType = strings.ToLower(mimeType)
	
	// Check image types
	imageTypes := []string{
		"image/jpeg", "image/jpg", "image/png", "image/gif", 
		"image/webp", "image/bmp", "image/svg+xml",
	}
	
	// Check video types
	videoTypes := []string{
		"video/mp4", "video/avi", "video/mov", "video/wmv",
		"video/flv", "video/webm", "video/mkv",
	}
	
	// Check audio types
	audioTypes := []string{
		"audio/mpeg", "audio/mp3", "audio/wav", "audio/aac",
		"audio/ogg", "audio/flac", "audio/m4a",
	}
	
	// Check document types
	documentTypes := []string{
		"application/pdf", "application/msword", 
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/plain", "text/csv", "text/html",
	}
	
	// Check archive types
	archiveTypes := []string{
		"application/zip", "application/x-rar-compressed",
		"application/x-7z-compressed", "application/gzip",
		"application/x-tar",
	}
	
	// Check image types
	for _, imgType := range imageTypes {
		if mimeType == imgType {
			return "image"
		}
	}
	
	// Check video types
	for _, vidType := range videoTypes {
		if mimeType == vidType {
			return "video"
		}
	}
	
	// Check audio types
	for _, audType := range audioTypes {
		if mimeType == audType {
			return "audio"
		}
	}
	
	// Check document types
	for _, docType := range documentTypes {
		if mimeType == docType {
			return "document"
		}
	}
	
	// Check archive types
	for _, arcType := range archiveTypes {
		if mimeType == arcType {
			return "archive"
		}
	}
	
	return "other"
}

// ValidateFileSize checks if file size is within limits
func (c *Client) ValidateFileSize(size int64) error {
	maxSize := int64(50 * 1024 * 1024) // 50MB
	
	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", size, maxSize)
	}
	
	return nil
}

// ValidateFileType checks if file type is allowed
func (c *Client) ValidateFileType(mimeType string) error {
	allowedTypes := []string{
		"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp", "image/svg+xml",
		"video/mp4", "video/avi", "video/mov", "video/wmv", "video/flv", "video/webm", "video/mkv",
		"audio/mpeg", "audio/mp3", "audio/wav", "audio/aac", "audio/ogg", "audio/flac", "audio/m4a",
		"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/plain", "text/csv", "text/html",
		"application/zip", "application/x-rar-compressed", "application/x-7z-compressed", "application/gzip", "application/x-tar",
	}
	
	mimeType = strings.ToLower(mimeType)
	for _, allowedType := range allowedTypes {
		if mimeType == allowedType {
			return nil
		}
	}
	
	return fmt.Errorf("file type %s is not allowed", mimeType)
}

