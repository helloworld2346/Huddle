package file

import (
	"context"
	"fmt"

	"huddle/internal/database"
	"huddle/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// Create creates a new file record
func (r *repository) Create(ctx context.Context, file *File) error {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		logger.Error("Failed to create file", zap.Error(err))
		return fmt.Errorf("failed to create file: %w", err)
	}
	return nil
}

// GetByID gets a file by ID
func (r *repository) GetByID(ctx context.Context, id uint) (*File, error) {
	var file File
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Conversation").
		Preload("Message").
		Preload("Shares").
		First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("file not found")
		}
		logger.Error("Failed to get file by ID", zap.Error(err), zap.Uint("file_id", id))
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return &file, nil
}

// GetByObjectKey gets a file by object key
func (r *repository) GetByObjectKey(ctx context.Context, objectKey string) (*File, error) {
	var file File
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Conversation").
		Preload("Message").
		Where("object_key = ?", objectKey).
		First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("file not found")
		}
		logger.Error("Failed to get file by object key", zap.Error(err), zap.String("object_key", objectKey))
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return &file, nil
}

// Update updates a file record
func (r *repository) Update(ctx context.Context, file *File) error {
	if err := r.db.WithContext(ctx).Save(file).Error; err != nil {
		logger.Error("Failed to update file", zap.Error(err), zap.Uint("file_id", file.ID))
		return fmt.Errorf("failed to update file: %w", err)
	}
	return nil
}

// Delete permanently deletes a file record
func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&File{}, id).Error; err != nil {
		logger.Error("Failed to delete file", zap.Error(err), zap.Uint("file_id", id))
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// SoftDelete soft deletes a file record
func (r *repository) SoftDelete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&File{}, id).Error; err != nil {
		logger.Error("Failed to soft delete file", zap.Error(err), zap.Uint("file_id", id))
		return fmt.Errorf("failed to soft delete file: %w", err)
	}
	return nil
}

// ListByUser gets files by user ID with pagination
func (r *repository) ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]File, int64, error) {
	var files []File
	var total int64

	offset := (page - 1) * pageSize

	// Get total count
	if err := r.db.WithContext(ctx).Model(&File{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&total).Error; err != nil {
		logger.Error("Failed to count user files", zap.Error(err), zap.Uint("user_id", userID))
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}

	// Get files
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&files).Error; err != nil {
		logger.Error("Failed to list user files", zap.Error(err), zap.Uint("user_id", userID))
		return nil, 0, fmt.Errorf("failed to list files: %w", err)
	}

	return files, total, nil
}

// ListByConversation gets files by conversation ID with pagination
func (r *repository) ListByConversation(ctx context.Context, conversationID uint, page, pageSize int) ([]File, int64, error) {
	var files []File
	var total int64

	offset := (page - 1) * pageSize

	// Get total count
	if err := r.db.WithContext(ctx).Model(&File{}).Where("conversation_id = ? AND deleted_at IS NULL", conversationID).Count(&total).Error; err != nil {
		logger.Error("Failed to count conversation files", zap.Error(err), zap.Uint("conversation_id", conversationID))
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}

	// Get files
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("conversation_id = ? AND deleted_at IS NULL", conversationID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&files).Error; err != nil {
		logger.Error("Failed to list conversation files", zap.Error(err), zap.Uint("conversation_id", conversationID))
		return nil, 0, fmt.Errorf("failed to list files: %w", err)
	}

	return files, total, nil
}

// SearchFiles searches files with filters
func (r *repository) SearchFiles(ctx context.Context, req *FileSearchRequest) ([]File, int64, error) {
	var files []File
	var total int64

	query := r.db.WithContext(ctx).Model(&File{}).Where("deleted_at IS NULL")

	// Apply filters
	if req.Query != "" {
		query = query.Where("file_name ILIKE ? OR original_name ILIKE ?", "%"+req.Query+"%", "%"+req.Query+"%")
	}

	if req.FileType != "" {
		// This would need to be implemented based on MIME type mapping
		// For now, we'll skip this filter
	}

	if req.ConversationID != nil {
		query = query.Where("conversation_id = ?", *req.ConversationID)
	}

	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}

	offset := (req.Page - 1) * req.PageSize

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		logger.Error("Failed to count search results", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}

	// Get files
	if err := query.
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&files).Error; err != nil {
		logger.Error("Failed to search files", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to search files: %w", err)
	}

	return files, total, nil
}

// CreateShare creates a new file share
func (r *repository) CreateShare(ctx context.Context, share *FileShare) error {
	if err := r.db.WithContext(ctx).Create(share).Error; err != nil {
		logger.Error("Failed to create file share", zap.Error(err))
		return fmt.Errorf("failed to create file share: %w", err)
	}
	return nil
}

// GetShareByID gets a file share by ID
func (r *repository) GetShareByID(ctx context.Context, id uint) (*FileShare, error) {
	var share FileShare
	if err := r.db.WithContext(ctx).
		Preload("File").
		Preload("SharedByUser").
		Preload("SharedWithUser").
		First(&share, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("file share not found")
		}
		logger.Error("Failed to get file share by ID", zap.Error(err), zap.Uint("share_id", id))
		return nil, fmt.Errorf("failed to get file share: %w", err)
	}
	return &share, nil
}

// GetSharesByFile gets all shares for a file
func (r *repository) GetSharesByFile(ctx context.Context, fileID uint) ([]FileShare, error) {
	var shares []FileShare
	if err := r.db.WithContext(ctx).
		Preload("File").
		Preload("SharedByUser").
		Preload("SharedWithUser").
		Where("file_id = ?", fileID).
		Find(&shares).Error; err != nil {
		logger.Error("Failed to get file shares", zap.Error(err), zap.Uint("file_id", fileID))
		return nil, fmt.Errorf("failed to get file shares: %w", err)
	}
	return shares, nil
}

// GetSharesByUser gets all shares by a user
func (r *repository) GetSharesByUser(ctx context.Context, userID uint) ([]FileShare, error) {
	var shares []FileShare
	if err := r.db.WithContext(ctx).
		Preload("File").
		Preload("SharedByUser").
		Preload("SharedWithUser").
		Where("shared_by = ? OR shared_with = ?", userID, userID).
		Find(&shares).Error; err != nil {
		logger.Error("Failed to get user shares", zap.Error(err), zap.Uint("user_id", userID))
		return nil, fmt.Errorf("failed to get user shares: %w", err)
	}
	return shares, nil
}

// DeleteShare deletes a file share
func (r *repository) DeleteShare(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&FileShare{}, id).Error; err != nil {
		logger.Error("Failed to delete file share", zap.Error(err), zap.Uint("share_id", id))
		return fmt.Errorf("failed to delete file share: %w", err)
	}
	return nil
}

// CheckFileAccess checks if a user has access to a file
func (r *repository) CheckFileAccess(ctx context.Context, fileID, userID uint) (bool, error) {
	var count int64
	
	// Check if user owns the file
	if err := r.db.WithContext(ctx).Model(&File{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", fileID, userID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check file ownership", zap.Error(err))
		return false, fmt.Errorf("failed to check file access: %w", err)
	}
	
	if count > 0 {
		return true, nil
	}
	
	// Check if file is shared with user
	if err := r.db.WithContext(ctx).Model(&FileShare{}).
		Where("file_id = ? AND shared_with = ?", fileID, userID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check file shares", zap.Error(err))
		return false, fmt.Errorf("failed to check file access: %w", err)
	}
	
	return count > 0, nil
}

// GetUserFiles gets files accessible by a user
func (r *repository) GetUserFiles(ctx context.Context, userID uint, page, pageSize int) ([]File, int64, error) {
	var files []File
	var total int64

	offset := (page - 1) * pageSize

	// Get total count - owned files + shared files
	var ownedCount, sharedCount int64
	
	if err := r.db.WithContext(ctx).Model(&File{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&ownedCount).Error; err != nil {
		logger.Error("Failed to count owned files", zap.Error(err), zap.Uint("user_id", userID))
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&FileShare{}).
		Where("shared_with = ?", userID).
		Count(&sharedCount).Error; err != nil {
		logger.Error("Failed to count shared files", zap.Error(err), zap.Uint("user_id", userID))
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}
	
	total = ownedCount + sharedCount

	// Get owned files
	var ownedFiles []File
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Conversation").
		Preload("Message").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&ownedFiles).Error; err != nil {
		logger.Error("Failed to list owned files", zap.Error(err), zap.Uint("user_id", userID))
		return nil, 0, fmt.Errorf("failed to list files: %w", err)
	}
	
	// Get shared files
	var sharedFiles []File
	sharedOffset := offset - int(ownedCount)
	if sharedOffset < 0 {
		sharedOffset = 0
	}
	
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Conversation").
		Preload("Message").
		Joins("JOIN file_shares ON files.id = file_shares.file_id").
		Where("file_shares.shared_with = ? AND files.deleted_at IS NULL", userID).
		Order("files.created_at DESC").
		Offset(sharedOffset).
		Limit(pageSize - len(ownedFiles)).
		Find(&sharedFiles).Error; err != nil {
		logger.Error("Failed to list shared files", zap.Error(err), zap.Uint("user_id", userID))
		return nil, 0, fmt.Errorf("failed to list files: %w", err)
	}
	
	// Combine results
	files = append(ownedFiles, sharedFiles...)

	return files, total, nil
}
