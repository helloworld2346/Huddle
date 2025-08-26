package file

import (
	"context"
	"mime/multipart"
)

type Repository interface {
	Create(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id uint) (*File, error)
	GetByObjectKey(ctx context.Context, objectKey string) (*File, error)
	Update(ctx context.Context, file *File) error
	Delete(ctx context.Context, id uint) error
	SoftDelete(ctx context.Context, id uint) error
	
	// File queries
	ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]File, int64, error)
	ListByConversation(ctx context.Context, conversationID uint, page, pageSize int) ([]File, int64, error)
	SearchFiles(ctx context.Context, req *FileSearchRequest) ([]File, int64, error)
	
	// File shares
	CreateShare(ctx context.Context, share *FileShare) error
	GetShareByID(ctx context.Context, id uint) (*FileShare, error)
	GetSharesByFile(ctx context.Context, fileID uint) ([]FileShare, error)
	GetSharesByUser(ctx context.Context, userID uint) ([]FileShare, error)
	DeleteShare(ctx context.Context, id uint) error
	
	// Access control
	CheckFileAccess(ctx context.Context, fileID, userID uint) (bool, error)
	GetUserFiles(ctx context.Context, userID uint, page, pageSize int) ([]File, int64, error)
}

type Service interface {
	// File operations
	UploadFile(ctx context.Context, userID uint, file multipart.File, header *multipart.FileHeader, req *UploadFileRequest) (*FileResponse, error)
	GetFile(ctx context.Context, fileID uint) (*FileResponse, error)
	GetFileByID(ctx context.Context, fileID, userID uint) (*FileResponse, error)
	UpdateFile(ctx context.Context, fileID, userID uint, req *UpdateFileRequest) (*FileResponse, error)
	DeleteFile(ctx context.Context, fileID, userID uint) error
	
	// File listing
	ListUserFiles(ctx context.Context, userID uint, page, pageSize int) (*FileListResponse, error)
	ListConversationFiles(ctx context.Context, conversationID, userID uint, page, pageSize int) (*FileListResponse, error)
	SearchFiles(ctx context.Context, userID uint, req *FileSearchRequest) (*FileListResponse, error)
	
	// File sharing
	ShareFile(ctx context.Context, userID uint, req *ShareFileRequest) (*FileShareResponse, error)
	GetFileShares(ctx context.Context, fileID, userID uint) ([]FileShareResponse, error)
	DeleteFileShare(ctx context.Context, shareID, userID uint) error
	
	// File access
	GetDownloadURL(ctx context.Context, fileID, userID uint) (string, error)
	CheckFileAccess(ctx context.Context, fileID, userID uint) (bool, error)
	
	// File processing
	GenerateThumbnail(ctx context.Context, fileID uint) (string, error)
	ProcessFileMetadata(ctx context.Context, fileID uint) error
	
	// Validation
	ValidateFile(header *multipart.FileHeader) error
	ValidateFileAccess(ctx context.Context, fileID, userID uint) error
}
