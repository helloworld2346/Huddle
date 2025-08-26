package file

import (
	"time"

	"huddle/internal/user"
)

// File represents a file stored in the system
type File struct {
	ID             uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         uint       `json:"user_id" gorm:"not null"`
	ConversationID *uint      `json:"conversation_id"`
	MessageID      *uint      `json:"message_id"`
	
	// File Information
	FileName      string `json:"file_name" gorm:"not null"`
	OriginalName  string `json:"original_name" gorm:"not null"`
	FileSize      int64  `json:"file_size" gorm:"not null"`
	MimeType      string `json:"mime_type" gorm:"not null"`
	FileExtension string `json:"file_extension"`
	
	// Storage Information
	BucketName string `json:"bucket_name" gorm:"not null;default:'huddle-files'"`
	ObjectKey  string `json:"object_key" gorm:"not null"`
	StoragePath string `json:"storage_path" gorm:"not null"`
	
	// File Processing
	IsProcessed  bool   `json:"is_processed" gorm:"default:false"`
	ThumbnailURL string `json:"thumbnail_url"`
	PreviewURL   string `json:"preview_url"`
	
	// Security & Access
	IsPublic    bool       `json:"is_public" gorm:"default:false"`
	AccessToken string     `json:"access_token"`
	ExpiresAt   *time.Time `json:"expires_at"`
	
	// Metadata
	Width    *int `json:"width"`    // For images/videos
	Height   *int `json:"height"`   // For images/videos
	Duration *int `json:"duration"` // For videos/audio (seconds)
	
	// Timestamps
	CreatedAt time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:now()"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	// Relations
	User         user.User `json:"user" gorm:"foreignKey:UserID"`
	Conversation *FileConversation `json:"conversation" gorm:"foreignKey:ConversationID"`
	Message      *FileMessage      `json:"message" gorm:"foreignKey:MessageID"`
	Shares       []FileShare       `json:"shares" gorm:"foreignKey:FileID"`
}

// FileConversation represents conversation info for file context
type FileConversation struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// FileMessage represents message info for file context
type FileMessage struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
}

// FileShare represents a file share
type FileShare struct {
	ID             uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	FileID         uint       `json:"file_id" gorm:"not null"`
	SharedBy       uint       `json:"shared_by" gorm:"not null"`
	SharedWith     *uint      `json:"shared_with"`
	ConversationID *uint      `json:"conversation_id"`
	
	// Share Settings
	CanDownload bool       `json:"can_download" gorm:"default:true"`
	CanEdit     bool       `json:"can_edit" gorm:"default:false"`
	ExpiresAt   *time.Time `json:"expires_at"`
	
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`

	// Relations
	File         File       `json:"file" gorm:"foreignKey:FileID"`
	SharedByUser user.User  `json:"shared_by_user" gorm:"foreignKey:SharedBy"`
	SharedWithUser *user.User `json:"shared_with_user" gorm:"foreignKey:SharedWith"`
}

// File Type Constants
const (
	FileTypeImage   = "image"
	FileTypeVideo   = "video"
	FileTypeAudio   = "audio"
	FileTypeDocument = "document"
	FileTypeArchive = "archive"
	FileTypeOther   = "other"
)

// DTOs for API requests/responses

// UploadFileRequest represents request to upload a file
type UploadFileRequest struct {
	ConversationID *uint `json:"conversation_id" form:"conversation_id"`
	MessageID      *uint `json:"message_id" form:"message_id"`
	IsPublic       bool  `json:"is_public" form:"is_public"`
}

// UpdateFileRequest represents request to update file metadata
type UpdateFileRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	IsPublic   bool   `json:"is_public"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// ShareFileRequest represents request to share a file
type ShareFileRequest struct {
	FileID         uint       `json:"file_id" binding:"required"`
	SharedWith     *uint      `json:"shared_with"`
	ConversationID *uint      `json:"conversation_id"`
	CanDownload    bool       `json:"can_download"`
	CanEdit        bool       `json:"can_edit"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

// FileResponse represents file response
type FileResponse struct {
	ID             uint       `json:"id"`
	UserID         uint       `json:"user_id"`
	ConversationID *uint      `json:"conversation_id"`
	MessageID      *uint      `json:"message_id"`
	
	// File Information
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	FileSize      int64  `json:"file_size"`
	MimeType      string `json:"mime_type"`
	FileExtension string `json:"file_extension"`
	FileType      string `json:"file_type"`
	
	// URLs
	DownloadURL   string `json:"download_url"`
	ThumbnailURL  string `json:"thumbnail_url,omitempty"`
	PreviewURL    string `json:"preview_url,omitempty"`
	
	// Security & Access
	IsPublic    bool       `json:"is_public"`
	AccessToken string     `json:"access_token,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	
	// Metadata
	Width    *int `json:"width,omitempty"`
	Height   *int `json:"height,omitempty"`
	Duration *int `json:"duration,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	User         user.UserResponse `json:"user"`
	Conversation *FileConversation `json:"conversation,omitempty"`
	Message      *FileMessage      `json:"message,omitempty"`
}

// FileListResponse represents list of files response
type FileListResponse struct {
	Files    []FileResponse `json:"files"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// FileShareResponse represents file share response
type FileShareResponse struct {
	ID             uint       `json:"id"`
	FileID         uint       `json:"file_id"`
	SharedBy       uint       `json:"shared_by"`
	SharedWith     *uint      `json:"shared_with"`
	ConversationID *uint      `json:"conversation_id"`
	
	// Share Settings
	CanDownload bool       `json:"can_download"`
	CanEdit     bool       `json:"can_edit"`
	ExpiresAt   *time.Time `json:"expires_at"`
	
	CreatedAt time.Time `json:"created_at"`
	
	// Relations
	File         FileResponse `json:"file"`
	SharedByUser user.UserResponse `json:"shared_by_user"`
	SharedWithUser *user.UserResponse `json:"shared_with_user,omitempty"`
}

// FileSearchRequest represents file search request
type FileSearchRequest struct {
	Query          string `json:"query" form:"query"`
	FileType       string `json:"file_type" form:"file_type"`
	ConversationID *uint  `json:"conversation_id" form:"conversation_id"`
	UserID         *uint  `json:"user_id" form:"user_id"`
	Page           int    `json:"page" form:"page" binding:"min=1"`
	PageSize       int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}
