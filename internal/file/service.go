package file

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"huddle/internal/user"
	"huddle/pkg/logger"
	"huddle/pkg/minio"

	"go.uber.org/zap"
)

type service struct {
	repo      Repository
	minioClient *minio.Client
}

func NewService(repo Repository) Service {
	minioClient := minio.GetClient()
	if minioClient == nil {
		logger.Info("MinIO client not found, creating new one...")
		// Create new client if not exists
		newClient, err := minio.NewClient()
		if err != nil {
			logger.Error("Failed to create MinIO client", zap.Error(err))
			// Return service without MinIO client for now
			return &service{
				repo:      repo,
				minioClient: nil,
			}
		} else {
			minio.SetClient(newClient)
			minioClient = newClient
		}
	}

	return &service{
		repo:      repo,
		minioClient: minioClient,
	}
}

// UploadFile uploads a file to MinIO and creates a file record
func (s *service) UploadFile(ctx context.Context, userID uint, file multipart.File, header *multipart.FileHeader, req *UploadFileRequest) (*FileResponse, error) {
	// Check if MinIO client is available
	if s.minioClient == nil {
		return nil, fmt.Errorf("MinIO client not available")
	}

	// Validate file
	if err := s.ValidateFile(header); err != nil {
		return nil, err
	}

	// Generate unique object key
	objectKey := s.minioClient.GenerateObjectKey(userID, header.Filename)
	
	// Upload file to MinIO
	if err := s.minioClient.UploadFile(ctx, objectKey, file, header); err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Get file type
	fileType := s.minioClient.GetFileType(header.Header.Get("Content-Type"))
	
	// Get file extension
	fileExtension := filepath.Ext(header.Filename)

	// Create file record
	fileRecord := &File{
		UserID:         userID,
		ConversationID: req.ConversationID,
		MessageID:      req.MessageID,
		FileName:       filepath.Base(objectKey),
		OriginalName:   header.Filename,
		FileSize:       header.Size,
		MimeType:       header.Header.Get("Content-Type"),
		FileExtension:  fileExtension,
		BucketName:     "huddle-files",
		ObjectKey:      objectKey,
		StoragePath:    objectKey,
		IsProcessed:    true,
		IsPublic:       req.IsPublic,
	}

	// Save to database
	if err := s.repo.Create(ctx, fileRecord); err != nil {
		// Clean up MinIO file if database save fails
		s.minioClient.DeleteFile(ctx, objectKey)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	// Generate download URL
	downloadURL, err := s.minioClient.GetPresignedURL(ctx, objectKey, 24*time.Hour)
	if err != nil {
		logger.Error("Failed to generate download URL", zap.Error(err))
		downloadURL = "" // Continue without download URL
	}

	// Build response
	response := &FileResponse{
		ID:             fileRecord.ID,
		UserID:         fileRecord.UserID,
		ConversationID: fileRecord.ConversationID,
		MessageID:      fileRecord.MessageID,
		FileName:       fileRecord.FileName,
		OriginalName:   fileRecord.OriginalName,
		FileSize:       fileRecord.FileSize,
		MimeType:       fileRecord.MimeType,
		FileExtension:  fileRecord.FileExtension,
		FileType:       fileType,
		DownloadURL:    downloadURL,
		ThumbnailURL:   fileRecord.ThumbnailURL,
		PreviewURL:     fileRecord.PreviewURL,
		IsPublic:       fileRecord.IsPublic,
		AccessToken:    fileRecord.AccessToken,
		ExpiresAt:      fileRecord.ExpiresAt,
		Width:          fileRecord.Width,
		Height:         fileRecord.Height,
		Duration:       fileRecord.Duration,
		CreatedAt:      fileRecord.CreatedAt,
		UpdatedAt:      fileRecord.UpdatedAt,
	}

	logger.Info("File uploaded successfully",
		zap.Uint("file_id", fileRecord.ID),
		zap.Uint("user_id", userID),
		zap.String("filename", header.Filename),
		zap.Int64("size", header.Size))

	return response, nil
}

// GetFile gets a file by ID (public access)
func (s *service) GetFile(ctx context.Context, fileID uint) (*FileResponse, error) {
	file, err := s.repo.GetByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// Check if file is public
	if !file.IsPublic {
		return nil, fmt.Errorf("file is not public")
	}

	return s.buildFileResponse(ctx, file), nil
}

// GetFileByID gets a file by ID with access control
func (s *service) GetFileByID(ctx context.Context, fileID, userID uint) (*FileResponse, error) {
	file, err := s.repo.GetByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// Check access
	if err := s.ValidateFileAccess(ctx, fileID, userID); err != nil {
		return nil, err
	}

	return s.buildFileResponse(ctx, file), nil
}

// UpdateFile updates file metadata
func (s *service) UpdateFile(ctx context.Context, fileID, userID uint, req *UpdateFileRequest) (*FileResponse, error) {
	// Check access
	if err := s.ValidateFileAccess(ctx, fileID, userID); err != nil {
		return nil, err
	}

	// Get file
	file, err := s.repo.GetByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// Update fields
	file.FileName = req.FileName
	file.IsPublic = req.IsPublic
	file.ExpiresAt = req.ExpiresAt
	file.UpdatedAt = time.Now()

	// Save to database
	if err := s.repo.Update(ctx, file); err != nil {
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	return s.buildFileResponse(ctx, file), nil
}

// DeleteFile deletes a file
func (s *service) DeleteFile(ctx context.Context, fileID, userID uint) error {
	// Check access
	if err := s.ValidateFileAccess(ctx, fileID, userID); err != nil {
		return err
	}

	// Get file
	file, err := s.repo.GetByID(ctx, fileID)
	if err != nil {
		return err
	}

	// Delete from MinIO
	if err := s.minioClient.DeleteFile(ctx, file.ObjectKey); err != nil {
		logger.Error("Failed to delete file from MinIO", zap.Error(err))
		// Continue with database deletion even if MinIO fails
	}

	// Delete from database
	if err := s.repo.Delete(ctx, fileID); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	logger.Info("File deleted successfully",
		zap.Uint("file_id", fileID),
		zap.Uint("user_id", userID))

	return nil
}

// ListUserFiles lists files owned by a user
func (s *service) ListUserFiles(ctx context.Context, userID uint, page, pageSize int) (*FileListResponse, error) {
	files, total, err := s.repo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]FileResponse, len(files))
	for i, file := range files {
		responses[i] = *s.buildFileResponse(ctx, &file)
	}

	return &FileListResponse{
		Files:    responses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// ListConversationFiles lists files in a conversation
func (s *service) ListConversationFiles(ctx context.Context, conversationID, userID uint, page, pageSize int) (*FileListResponse, error) {
	files, total, err := s.repo.ListByConversation(ctx, conversationID, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]FileResponse, len(files))
	for i, file := range files {
		responses[i] = *s.buildFileResponse(ctx, &file)
	}

	return &FileListResponse{
		Files:    responses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// SearchFiles searches files with filters
func (s *service) SearchFiles(ctx context.Context, userID uint, req *FileSearchRequest) (*FileListResponse, error) {
	files, total, err := s.repo.SearchFiles(ctx, req)
	if err != nil {
		return nil, err
	}

	responses := make([]FileResponse, len(files))
	for i, file := range files {
		responses[i] = *s.buildFileResponse(ctx, &file)
	}

	return &FileListResponse{
		Files:    responses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ShareFile shares a file with another user or conversation
func (s *service) ShareFile(ctx context.Context, userID uint, req *ShareFileRequest) (*FileShareResponse, error) {
	// Check if user owns the file
	file, err := s.repo.GetByID(ctx, req.FileID)
	if err != nil {
		return nil, err
	}

	if file.UserID != userID {
		return nil, fmt.Errorf("access denied: you don't own this file")
	}

	// Create share record
	share := &FileShare{
		FileID:         req.FileID,
		SharedBy:       userID,
		SharedWith:     req.SharedWith,
		ConversationID: req.ConversationID,
		CanDownload:    req.CanDownload,
		CanEdit:        req.CanEdit,
		ExpiresAt:      req.ExpiresAt,
	}

	if err := s.repo.CreateShare(ctx, share); err != nil {
		return nil, fmt.Errorf("failed to create file share: %w", err)
	}

	// Get share with relations
	shareWithRelations, err := s.repo.GetShareByID(ctx, share.ID)
	if err != nil {
		return nil, err
	}

	return s.buildFileShareResponse(ctx, shareWithRelations), nil
}

// GetFileShares gets all shares for a file
func (s *service) GetFileShares(ctx context.Context, fileID, userID uint) ([]FileShareResponse, error) {
	// Check if user owns the file
	file, err := s.repo.GetByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	if file.UserID != userID {
		return nil, fmt.Errorf("access denied: you don't own this file")
	}

	shares, err := s.repo.GetSharesByFile(ctx, fileID)
	if err != nil {
		return nil, err
	}

	responses := make([]FileShareResponse, len(shares))
	for i, share := range shares {
		responses[i] = *s.buildFileShareResponse(ctx, &share)
	}

	return responses, nil
}

// DeleteFileShare deletes a file share
func (s *service) DeleteFileShare(ctx context.Context, shareID, userID uint) error {
	// Get share
	share, err := s.repo.GetShareByID(ctx, shareID)
	if err != nil {
		return err
	}

	// Check if user owns the file
	file, err := s.repo.GetByID(ctx, share.FileID)
	if err != nil {
		return err
	}

	if file.UserID != userID {
		return fmt.Errorf("access denied: you don't own this file")
	}

	return s.repo.DeleteShare(ctx, shareID)
}

// GetDownloadURL generates a download URL for a file
func (s *service) GetDownloadURL(ctx context.Context, fileID, userID uint) (string, error) {
	// Check access
	if err := s.ValidateFileAccess(ctx, fileID, userID); err != nil {
		return "", err
	}

	// Get file
	file, err := s.repo.GetByID(ctx, fileID)
	if err != nil {
		return "", err
	}

	// Generate presigned URL
	url, err := s.minioClient.GetPresignedURL(ctx, file.ObjectKey, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}

	return url, nil
}

// CheckFileAccess checks if a user has access to a file
func (s *service) CheckFileAccess(ctx context.Context, fileID, userID uint) (bool, error) {
	return s.repo.CheckFileAccess(ctx, fileID, userID)
}

// GenerateThumbnail generates a thumbnail for an image file
func (s *service) GenerateThumbnail(ctx context.Context, fileID uint) (string, error) {
	// This would implement thumbnail generation
	// For now, return empty string
	return "", nil
}

// ProcessFileMetadata processes file metadata (dimensions, duration, etc.)
func (s *service) ProcessFileMetadata(ctx context.Context, fileID uint) error {
	// This would implement metadata processing
	// For now, return nil
	return nil
}

// ValidateFile validates file size and type
func (s *service) ValidateFile(header *multipart.FileHeader) error {
	// Validate file size
	if err := s.minioClient.ValidateFileSize(header.Size); err != nil {
		return err
	}

	// Validate file type
	mimeType := header.Header.Get("Content-Type")
	if err := s.minioClient.ValidateFileType(mimeType); err != nil {
		return err
	}

	return nil
}

// ValidateFileAccess validates if a user has access to a file
func (s *service) ValidateFileAccess(ctx context.Context, fileID, userID uint) error {
	hasAccess, err := s.repo.CheckFileAccess(ctx, fileID, userID)
	if err != nil {
		return fmt.Errorf("failed to check file access: %w", err)
	}

	if !hasAccess {
		return fmt.Errorf("access denied: you don't have permission to access this file")
	}

	return nil
}

// buildFileResponse builds a FileResponse from a File
func (s *service) buildFileResponse(ctx context.Context, file *File) *FileResponse {
	// Generate download URL
	downloadURL, err := s.minioClient.GetPresignedURL(ctx, file.ObjectKey, 24*time.Hour)
	if err != nil {
		logger.Error("Failed to generate download URL", zap.Error(err))
		downloadURL = ""
	}

	// Get file type
	fileType := s.minioClient.GetFileType(file.MimeType)

	// Build user response
	userResponse := file.User.ToResponse()

	// Build conversation response
	var conversationResponse *FileConversation
	if file.Conversation != nil {
		conversationResponse = file.Conversation
	}

	// Build message response
	var messageResponse *FileMessage
	if file.Message != nil {
		messageResponse = file.Message
	}

	return &FileResponse{
		ID:             file.ID,
		UserID:         file.UserID,
		ConversationID: file.ConversationID,
		MessageID:      file.MessageID,
		FileName:       file.FileName,
		OriginalName:   file.OriginalName,
		FileSize:       file.FileSize,
		MimeType:       file.MimeType,
		FileExtension:  file.FileExtension,
		FileType:       fileType,
		DownloadURL:    downloadURL,
		ThumbnailURL:   file.ThumbnailURL,
		PreviewURL:     file.PreviewURL,
		IsPublic:       file.IsPublic,
		AccessToken:    file.AccessToken,
		ExpiresAt:      file.ExpiresAt,
		Width:          file.Width,
		Height:         file.Height,
		Duration:       file.Duration,
		CreatedAt:      file.CreatedAt,
		UpdatedAt:      file.UpdatedAt,
		User:           userResponse,
		Conversation:   conversationResponse,
		Message:        messageResponse,
	}
}

// buildFileShareResponse builds a FileShareResponse from a FileShare
func (s *service) buildFileShareResponse(ctx context.Context, share *FileShare) *FileShareResponse {
	// Build file response
	fileResponse := s.buildFileResponse(ctx, &share.File)

	// Build user responses
	sharedByUserResponse := share.SharedByUser.ToResponse()
	var sharedWithUserResponse *user.UserResponse
	if share.SharedWithUser != nil {
		response := share.SharedWithUser.ToResponse()
		sharedWithUserResponse = &response
	}

	return &FileShareResponse{
		ID:             share.ID,
		FileID:         share.FileID,
		SharedBy:       share.SharedBy,
		SharedWith:     share.SharedWith,
		ConversationID: share.ConversationID,
		CanDownload:    share.CanDownload,
		CanEdit:        share.CanEdit,
		ExpiresAt:      share.ExpiresAt,
		CreatedAt:      share.CreatedAt,
		File:           *fileResponse,
		SharedByUser:   sharedByUserResponse,
		SharedWithUser: sharedWithUserResponse,
	}
}
