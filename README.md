# 🚀 Huddle - Real-time Chat Application

Huddle là một ứng dụng chat realtime hiện đại, lấy cảm hứng từ Messenger, được xây dựng bằng Golang và PostgreSQL. Ứng dụng hỗ trợ chat 1-1, nhóm chat, kết bạn, chia sẻ file và nhiều tính năng khác.

## 📋 Mục lục

- [Tính năng](#-tính-năng)
- [Kiến trúc hệ thống](#-kiến-trúc-hệ-thống)
- [Cấu trúc dự án](#-cấu-trúc-dự-án)
- [Tiến độ phát triển](#-tiến-độ-phát-triển)
- [Flow hoạt động](#-flow-hoạt-động)
- [Công nghệ sử dụng](#-công-nghệ-sử-dụng)
- [Roadmap phát triển](#-roadmap-phát-triển)
- [Cài đặt và chạy](#-cài-đặt-và-chạy)
- [API Documentation](#-api-documentation)

## ✨ Tính năng

### 🔐 Authentication & User Management

- [x] **Đăng ký, đăng nhập, đăng xuất** - Complete auth flow
- [x] **JWT token authentication** - Access & refresh tokens
- [x] **Token blacklisting** - Redis-based token revocation
- [x] **Session management** - Database & Redis sessions
- [x] **Password reset** - Forgot/reset password flow
- [x] **Activity logging** - User activity tracking
- [x] **Rate limiting** - Login attempt protection
- [x] **Quản lý profile người dùng** - CRUD operations
- [x] **Tìm kiếm người dùng** - Search by username/email
- [ ] Upload avatar (MinIO integration pending)

### 👥 Friend System

- [x] **Gửi lời mời kết bạn** - Send friend requests
- [x] **Chấp nhận/từ chối lời mời** - Accept/reject friend requests
- [x] **Danh sách bạn bè** - Friends list management
- [x] **Chặn/bỏ chặn người dùng** - Block/unblock users
- [x] **Quản lý lời mời kết bạn** - Friend request management
- [x] **Bidirectional friendships** - Two-way friend relationships
- [x] **Validation & error handling** - Proper business logic

### 💬 Chat Features

- [x] **Chat 1-1**: Tin nhắn riêng tư giữa 2 người ✅
- [x] **Group Chat**: Chat nhóm với nhiều thành viên ✅
- [x] **Message history**: Lưu trữ và tìm kiếm tin nhắn ✅
- [x] **Message reactions**: Like, heart, emoji reactions ✅
- [x] **Read receipts**: Hiển thị trạng thái đã đọc (conversation level) ✅
- [x] **Real-time messaging**: WebSocket cho tin nhắn tức thì ✅
- [x] **Typing indicators**: Hiển thị đang gõ ✅
- [x] **Online/Offline status**: Track user presence ✅

### 📁 File Sharing

- [x] **Upload và chia sẻ file** - Complete file upload system ✅
- [x] **Hỗ trợ nhiều định dạng file** - Images, videos, documents, archives ✅
- [x] **Lưu trữ file trên MinIO** - Object storage integration ✅
- [x] **File metadata management** - Database storage with MinIO ✅
- [x] **File sharing system** - Share files with users/conversations ✅
- [x] **Access control** - Public/private files, ownership validation ✅
- [x] **File search** - Search by name, type, conversation ✅
- [x] **Presigned URLs** - Secure file access ✅
- [x] **File validation** - Size limits, type restrictions ✅
- [ ] Preview hình ảnh (future enhancement)
- [ ] Thumbnail generation (future enhancement)

### 🏢 Group Management

- [x] **Tạo nhóm chat** - Create group conversations ✅
- [x] **Thêm/xóa thành viên** - Add/remove participants ✅
- [x] **Phân quyền admin/member** - Role-based permissions ✅
- [x] **Quản lý thông tin nhóm** - Update conversation details ✅
- [ ] Avatar nhóm (MinIO integration pending)

### 🔔 Notifications

- [x] **Real-time notifications** - WebSocket-based notifications ✅
- [x] **Thông báo tin nhắn mới** - New message notifications ✅
- [x] **Thông báo lời mời kết bạn** - Friend request notifications ✅
- [x] **Online/offline status** - User presence notifications ✅
- [ ] Push notifications (Mobile app pending)

## 🏗️ Kiến trúc hệ thống

### Backend Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   WebSocket     │    │   PostgreSQL    │
│   (React/Vue)   │◄──►│   Connection    │◄──►│   Database      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP API      │    │   Redis Cache   │    │   MinIO Storage │
│   (Gin)         │    │   (Sessions)    │    │   (Files)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Database Schema

#### ✅ **Core Tables**

- **Users**: Thông tin người dùng, profiles, settings
- **Sessions**: Phiên đăng nhập và token management
- **Password Resets**: Token reset mật khẩu
- **User Activities**: Log hoạt động người dùng

#### ✅ **Friend System Tables**

- **Friend Requests**: Lời mời kết bạn (pending, accepted, rejected, cancelled)
- **Friendships**: Mối quan hệ bạn bè (bidirectional)
- **Blocked Users**: Người dùng bị chặn

#### ✅ **Chat System Tables**

- **Conversations**: Cuộc hội thoại (direct, group)
- **Conversation Participants**: Thành viên conversation với roles (admin, member)
- **Messages**: Tin nhắn (text, image, file, system)
- **Message Reactions**: Phản ứng tin nhắn (like, love, haha, wow, sad, angry)
- **Message Reads**: Trạng thái đã đọc tin nhắn (future enhancement)

#### ✅ **File System Tables**

- **Files**: File metadata và thông tin lưu trữ
- **File Shares**: Chia sẻ file với users/conversations

#### ✅ **Database Schema Details**

**Users Table:**

- `id`, `username`, `email`, `display_name`, `bio`, `avatar`
- `is_public`, `last_login`, `login_attempts`, `locked_until`
- `created_at`, `updated_at`, `deleted_at`

**Sessions Table:**

- `id`, `user_id`, `token`, `expires_at`, `created_at`

**Friend Requests Table:**

- `id`, `sender_id`, `receiver_id`, `status`, `created_at`, `updated_at`

**Friendships Table:**

- `id`, `user_id`, `friend_id`, `created_at`

**Blocked Users Table:**

- `id`, `blocker_id`, `blocked_id`, `created_at`

**Conversations Table:**

- `id`, `name`, `type`, `created_by`, `created_at`, `updated_at`

**Conversation Participants Table:**

- `id`, `conversation_id`, `user_id`, `role`, `joined_at`, `last_read_at`

**Messages Table:**

- `id`, `conversation_id`, `sender_id`, `content`, `message_type`
- `file_url`, `file_name`, `file_size`, `reply_to_id`
- `is_edited`, `edited_at`, `created_at`, `updated_at`

**Message Reactions Table:**

- `id`, `message_id`, `user_id`, `reaction_type`, `created_at`
- Unique constraint: `(message_id, user_id, reaction_type)`

**Message Reads Table:**

- `id`, `message_id`, `user_id`, `read_at` (future enhancement)

**Files Table:**

- `id`, `user_id`, `conversation_id`, `message_id`
- `file_name`, `original_name`, `file_size`, `mime_type`, `file_extension`
- `bucket_name`, `object_key`, `storage_path`
- `is_processed`, `thumbnail_url`, `preview_url`
- `is_public`, `access_token`, `expires_at`
- `width`, `height`, `duration` (for media files)
- `created_at`, `updated_at`, `deleted_at`

**File Shares Table:**

- `id`, `file_id`, `shared_by`, `shared_with`, `conversation_id`
- `can_download`, `can_edit`, `expires_at`, `created_at`

#### ⏳ **Future Tables**

- **Groups**: Thông tin nhóm (separate from conversations)
- **Group Members**: Thành viên nhóm

## 📁 Cấu trúc dự án

```
huddle/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                 # Entry point ✅
│   ├── internal/
│   │   ├── app/                        # App setup ✅
│   │   │   └── app.go
│   │   ├── auth/                       # Authentication module ✅
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── user/                       # User management module ✅
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── friend/                     # Friend system module ✅
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── conversation/               # Conversation management ✅
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── message/                    # Message system ✅
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── websocket/                  # WebSocket hub ✅
│   │   │   ├── hub.go
│   │   │   ├── client.go
│   │   │   ├── service.go
│   │   │   ├── handler.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── file/                       # File handling module ✅
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── health/                     # Health check ✅
│   │   │   ├── handler.go
│   │   │   └── routes.go
│   │   ├── database/                   # Database connection ✅
│   │   │   └── connection.go
│   │   ├── middleware/                 # HTTP middleware ✅
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   └── config/                     # Configuration ✅
│   │       ├── config.go
│   │       ├── redis.go
│   │       └── app.env
│   ├── pkg/
│   │   ├── auth/                       # Authentication utilities ✅
│   │   │   ├── jwt.go
│   │   │   ├── password.go
│   │   │   └── redis.go
│   │   ├── minio/                      # MinIO client utilities ✅
│   │   │   └── client.go
│   │   ├── logger/                     # Structured logging ✅
│   │   │   └── logger.go
│   │   ├── utils/                      # Common utilities ✅
│   │   │   └── response.go
│   │   └── validation/                 # Request validation ✅
│   │       └── validator.go
│   ├── migrations/                     # Database migrations ✅
│   │   ├── 001_initial_schema.sql
│   │   ├── 002_auth_schema.sql
│   │   ├── 003_update_user_schema.sql
│   │   ├── 004_auth_tables.sql
│   │   ├── 005_friend_system.sql
│   │   ├── 006_chat_system.sql
│   │   └── 007_file_system.sql
│   ├── go.mod
│   └── go.sum
├── frontend/                           # ⏳ Chưa implement
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   ├── store/
│   │   └── utils/
│   ├── public/
│   ├── package.json
│   └── README.md
├── docker-compose.yml                  # Docker setup ✅
├── Dockerfile.backend                  # ⏳ Chưa tạo
├── Dockerfile.frontend                 # ⏳ Chưa tạo
├── Makefile                           # Development commands ✅
├── test_websocket.html                # WebSocket testing tool ✅
├── .gitignore                         # Git ignore ✅
├── SETUP.md                           # Setup guide ✅
└── README.md
```

## 🎯 Tiến độ phát triển

### 🚀 **Current Status (August 2025)**

**✅ Phase 1, 2, 3 & 4 COMPLETED** - Core infrastructure, authentication system, chat system, và real-time messaging đã hoàn thành 100%

**🎯 Next Target**: Frontend Development (Phase 6)

**📊 Progress**: 99% of total project (Core features + Friend System + Conversation System + Message System + WebSocket Hub + Online/Offline System + File Sharing System ready)

### ✅ **Đã hoàn thành (Phase 1 - Foundation)**

#### **Core Infrastructure:**

- [x] **Project Structure** - Cấu trúc thư mục modular
- [x] **Go Modules** - Dependencies management
- [x] **Configuration** - Environment variables với godotenv
- [x] **Docker Setup** - PostgreSQL và Redis containers
- [x] **Database Schema** - Initial migration với users và sessions
- [x] **Database Connection** - PostgreSQL với GORM
- [x] **Redis Connection** - Cache và session storage

#### **Web Framework:**

- [x] **Gin Server** - HTTP web framework
- [x] **CORS Middleware** - Cross-origin support
- [x] **Structured Logging** - JSON logging với Zap
- [x] **Request Logging** - Performance metrics
- [x] **Error Handling** - Graceful shutdown
- [x] **Health Check API** - `/api/health`, `/api/health/ping`

#### **Development Tools:**

- [x] **Makefile** - Development commands
- [x] **Git Setup** - Version control
- [x] **Documentation** - README và SETUP guides

### ✅ **Đã hoàn thành (Phase 2 - Core Features)**

#### **Authentication System:**

- [x] **JWT Utilities** - Token generation, validation, refresh
- [x] **Password Utilities** - Hashing, validation, strength check
- [x] **Redis Utilities** - Token blacklisting, session storage
- [x] **Authentication Middleware** - JWT validation, user context, blacklist check
- [x] **Request Validation** - Input validation utilities
- [x] **Database Schema** - Auth tables và indexes
- [x] **User Registration/Login** - Complete auth endpoints
- [x] **Session Management** - Database & Redis sessions
- [x] **Token Blacklisting** - Immediate token revocation
- [x] **Password Reset** - Forgot/reset password flow
- [x] **Activity Logging** - User activity tracking
- [x] **Rate Limiting** - Login attempt protection

#### **User Management:**

- [x] **User CRUD Operations** - Create, read, update, delete users
- [x] **Profile Management** - Update user profile
- [x] **Password Management** - Change password
- [x] **User Search** - Search by username/email
- [x] **Current User** - Get authenticated user info
- [ ] Avatar upload (MinIO integration pending)

### ✅ **Đã hoàn thành (Phase 2 - Friend System)**

- [x] **Friend requests** - Send, accept, reject, cancel
- [x] **Friend list management** - Get friends, remove friends
- [x] **User blocking** - Block/unblock users
- [x] **Validation & error handling** - Complete business logic

### ✅ **Đã hoàn thành (Phase 3 - Conversation System)**

- [x] **Conversation Management** - Create, list, update, delete conversations
- [x] **Participant Management** - Add, remove, leave conversations
- [x] **Smart Admin Transfer** - Hybrid admin leave logic with auto-promote
- [x] **Database Schema** - conversations, conversation_participants, messages, message_reactions, message_reads
- [x] **API Endpoints** - Complete conversation system APIs
- [x] **Business Logic** - Admin validation, access control, auto-promotion
- [x] **Testing** - All success and error cases tested

### ✅ **Đã hoàn thành (Phase 3 - Message System)**

- [x] **Message CRUD** - Create, read, update, delete messages
- [x] **Message Reactions** - Add/remove reactions (like, love, haha, wow, sad, angry)
- [x] **Message Search** - Search messages by content
- [x] **Message History** - Retrieve chat history with pagination
- [x] **Access Control** - Only conversation participants can access messages
- [x] **Message Validation** - Sender validation, content validation
- [x] **Database Schema** - messages, message_reactions, message_reads tables
- [x] **API Endpoints** - Complete message system APIs
- [x] **Testing** - All message features tested successfully

### ✅ **Đã hoàn thành (Phase 4 - Real-time Features)**

- [x] **WebSocket Hub** - Real-time communication ✅
- [x] **Real-time messaging** - Instant message delivery ✅
- [x] **Online/offline status** - User presence tracking ✅
- [x] **Typing indicators** - Real-time typing status ✅
- [x] **Connection health checker** - Automatic offline detection ✅
- [x] **Room-based messaging** - Conversation-specific broadcasting ✅
- [x] **JWT authentication** - Secure WebSocket connections ✅

### ✅ **Đã hoàn thành (Phase 5 - File Sharing)**

- [x] **MinIO Integration** - Object storage setup ✅
- [x] **File Upload/Download** - Complete file management ✅
- [x] **File Metadata Management** - Database storage ✅
- [x] **File Sharing System** - Share with users/conversations ✅
- [x] **Access Control** - Public/private files ✅
- [x] **File Validation** - Size limits, type restrictions ✅
- [x] **Presigned URLs** - Secure file access ✅
- [x] **File Search** - Search by name, type ✅
- [x] **Conversation File Isolation** - Files separated by conversation ✅
- [x] **Error Handling** - Complete error cases tested ✅

### 📋 **Còn lại (Phase 6)**

#### **Frontend Development:**

- [ ] React/Vue.js setup
- [ ] UI components
- [ ] Real-time chat interface
- [ ] File upload interface
- [ ] User management interface

#### **Advanced Features:**

- [ ] Push notifications
- [ ] Voice messages
- [ ] Message encryption
- [ ] Advanced search
- [ ] Image preview
- [ ] Avatar upload

#### **Frontend:**

- [ ] React/Vue.js setup
- [ ] UI components
- [ ] Real-time chat interface

#### **Advanced Features:**

- [ ] Push notifications
- [ ] Voice messages
- [ ] Message encryption
- [ ] Advanced search

## 🔄 Flow hoạt động

### 1. Authentication Flow

```
User → Register/Login → JWT Token → WebSocket Connection → Start Chatting
```

### 2. Friend Request Flow

```
User A → Send friend request → User B → Accept/Reject → Friendship created
```

### 3. Direct Message Flow

```
User A → Send DM → WebSocket → Server → Database → Broadcast to User B
```

### 4. Group Chat Flow

```
User → Send group message → WebSocket → Server → Database → Broadcast to all group members
```

### 5. File Upload Flow

```
User → Upload file → MinIO → Get URL → Save to database → Share with users/conversations
```

### 6. WebSocket Connection Flow

```
Client → Connect WebSocket → Authenticate → Join user room → Listen for messages
```

## 🛠️ Công nghệ sử dụng

### Backend ✅

- **Golang** (1.24.6) - Ngôn ngữ lập trình chính
- **Gin** (v1.10.1) - HTTP web framework
- **Gorilla WebSocket** - Real-time communication ✅
- **GORM** (v1.30.1) - ORM cho database
- **PostgreSQL** (15-alpine) - Relational database
- **Redis** (7-alpine) - Cache, session storage, token blacklisting ✅
- **MinIO** - Object storage cho file ✅
- **JWT** - Authentication tokens ✅
- **bcrypt** - Password hashing ✅
- **Zap** (v1.27.0) - Structured logging

### Frontend ⏳

- **React** (18+) - UI framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling framework
- **Socket.io-client** - WebSocket client
- **React Query** - Data fetching
- **React Router** - Client-side routing
- **Zustand** - State management

### DevOps & Tools ✅

- **Docker** - Containerization
- **Docker Compose** - Multi-container setup
- **Git** - Version control
- **Postman** - API testing
- **TablePlus** - Database management

## 🗺️ Roadmap phát triển

### ✅ Phase 1: Foundation (Đã hoàn thành)

**Mục tiêu**: Setup cơ bản và infrastructure

- [x] Setup project structure
- [x] Database schema và migrations
- [x] Gin server với middleware
- [x] Structured logging với Zap
- [x] Health check API
- [x] Docker setup cho PostgreSQL và Redis
- [x] Configuration management

### ✅ Phase 2: Core Features (Đã hoàn thành)

**Mục tiêu**: Authentication và user management

- [x] User authentication (register/login/logout)
- [x] JWT token management (access/refresh tokens)
- [x] Token blacklisting với Redis
- [x] Session management (database & Redis)
- [x] Password reset functionality
- [x] User profile management
- [x] User search và CRUD operations
- [x] Activity logging và rate limiting
- [ ] File upload với MinIO (pending)

### ✅ Phase 3: Advanced Features (Đã hoàn thành)

**Mục tiêu**: Chat và friend system

- [x] Friend system (request, accept, reject)
- [x] Direct messaging
- [x] Group creation và management
- [x] Message history
- [x] Message reactions
- [x] Online status

### ✅ Phase 4: Real-time Features (Đã hoàn thành)

**Mục tiêu**: WebSocket và real-time messaging

- [x] WebSocket hub
- [x] Real-time messaging
- [x] Online/offline status
- [x] Typing indicators
- [x] Connection health checker

### ✅ Phase 5: File Sharing (Đã hoàn thành)

**Mục tiêu**: File management và sharing

- [x] MinIO integration
- [x] File upload/download
- [x] File sharing system
- [x] Access control
- [x] File validation

### 📋 Phase 6: Frontend Development (Chưa bắt đầu)

**Mục tiêu**: Frontend interface

- [ ] React/Vue.js setup
- [ ] UI components
- [ ] Real-time chat interface
- [ ] File upload interface
- [ ] User management interface

### 📋 Phase 7: Enhancement (Chưa bắt đầu)

**Mục tiêu**: Tính năng bổ sung

- [ ] Push notifications
- [ ] Voice messages
- [ ] Video calls (future)
- [ ] Message encryption
- [ ] Advanced search
- [ ] Message forwarding
- [ ] Mobile responsive

## 🚀 Cài đặt và chạy

### Yêu cầu hệ thống

- Go 1.24.6+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### 🧪 Testing WebSocket

Để test real-time messaging, sử dụng file `test_websocket.html`:

```bash
# Mở file test trong browser
open test_websocket.html

# Hoặc truy cập trực tiếp
# file:///path/to/huddle/test_websocket.html
```

**Test Steps:**

1. Login với 2 users khác nhau (testuser1, testuser2)
2. Connect WebSocket cho cả 2 users
3. Join conversation 10
4. Gửi messages qua API - sẽ thấy real-time broadcasting
5. Test typing indicators
6. Check online users

### Quick Start với Docker

```bash
# Clone repository
git clone https://github.com/your-username/huddle.git
cd huddle

# Chạy với Docker Compose (PostgreSQL, Redis, MinIO)
make docker-up

# Download dependencies
make deps

# Chạy ứng dụng
make run

# Truy cập ứng dụng
# Backend API: http://localhost:8080
# Health Check: http://localhost:8080/api/health
# MinIO Console: http://localhost:9001 (minioadmin/minioadmin)
```

### Development Setup

```bash
# Backend
cd backend
go mod download
go run cmd/server/main.go

# Hoặc sử dụng Makefile
make dev  # docker-up + deps + run
```

## 📚 API Documentation

### ✅ Available Endpoints

#### Health Check Endpoints

- `GET /` - Welcome page
- `GET /api/health` - Health check với database và Redis status
- `GET /api/health/ping` - Simple ping endpoint

#### Example Responses

**Health Check:**

```json
{
  "status": "healthy",
  "timestamp": "2025-08-26T14:26:08.897371+07:00",
  "services": {
    "database": "up",
    "redis": "up"
  },
  "version": "1.0.0"
}
```

**Welcome Page:**

```json
{
  "message": "Welcome to Huddle API",
  "version": "1.0.0",
  "docs": "/api/health"
}
```

### ⏳ Planned Endpoints

#### Authentication Endpoints

- `POST /api/auth/register` - Đăng ký user mới ✅
- `POST /api/auth/login` - Đăng nhập ✅
- `POST /api/auth/logout` - Đăng xuất (blacklist tokens) ✅
- `POST /api/auth/refresh` - Refresh access token ✅
- `POST /api/auth/forgot-password` - Gửi email reset password ✅
- `POST /api/auth/reset-password` - Reset password với token ✅
- `GET /api/auth/stats` - Thống kê auth (protected) ✅

#### User Endpoints ✅

- `GET /api/users` - Lấy danh sách users ✅
- `GET /api/users/search` - Tìm kiếm users ✅
- `GET /api/users/:id` - Lấy thông tin user theo ID ✅
- `GET /api/users/username/:username` - Lấy user theo username ✅
- `GET /api/users/me` - Lấy thông tin user hiện tại (protected) ✅
- `PUT /api/users/me` - Cập nhật profile (protected) ✅
- `DELETE /api/users/me` - Xóa user (protected) ✅
- `PUT /api/users/me/password` - Đổi mật khẩu (protected) ✅
- `PUT /api/users/me/avatar` - Upload avatar (protected) ✅

#### Friend Endpoints ✅

- `GET /api/friends` - Lấy danh sách bạn bè ✅
- `POST /api/friends/requests` - Gửi lời mời kết bạn ✅
- `PUT /api/friends/requests/respond` - Phản hồi lời mời ✅
- `GET /api/friends/requests` - Lấy danh sách lời mời ✅
- `GET /api/friends/requests/sent` - Lấy lời mời đã gửi ✅
- `DELETE /api/friends/requests/:request_id` - Hủy lời mời ✅
- `DELETE /api/friends/:friend_id` - Xóa bạn bè ✅
- `POST /api/friends/block` - Chặn người dùng ✅
- `DELETE /api/friends/block/:user_id` - Bỏ chặn ✅
- `GET /api/friends/blocked` - Danh sách người bị chặn ✅

#### Conversation Endpoints ✅

- `POST /api/conversations` - Tạo conversation ✅
- `GET /api/conversations` - Lấy danh sách conversations ✅
- `GET /api/conversations/:id` - Lấy conversation chi tiết ✅
- `PUT /api/conversations/:id` - Cập nhật conversation ✅
- `DELETE /api/conversations/:id` - Xóa conversation ✅
- `POST /api/conversations/:id/participants` - Thêm thành viên ✅
- `DELETE /api/conversations/:id/participants` - Xóa thành viên ✅
- `POST /api/conversations/:id/leave` - Rời conversation ✅

#### Message Endpoints ✅

- `POST /api/conversations/:id/messages` - Gửi tin nhắn ✅
- `GET /api/conversations/:id/messages` - Lấy tin nhắn ✅
- `GET /api/conversations/:id/messages/before` - Lấy tin nhắn trước ID ✅
- `GET /api/conversations/:id/messages/search` - Tìm kiếm tin nhắn ✅
- `GET /api/conversations/:id/messages/:message_id` - Lấy tin nhắn chi tiết ✅
- `PUT /api/conversations/:id/messages/:message_id` - Cập nhật tin nhắn ✅
- `DELETE /api/conversations/:id/messages/:message_id` - Xóa tin nhắn ✅
- `POST /api/conversations/:id/messages/:message_id/reactions` - Thêm reaction ✅
- `DELETE /api/conversations/:id/messages/:message_id/reactions/:reaction_type` - Xóa reaction ✅

#### File Endpoints ✅

- `POST /api/files/upload` - Upload file ✅
- `GET /api/files/my` - Lấy files của user ✅
- `GET /api/files/search` - Tìm kiếm files ✅
- `GET /api/files/:id` - Lấy file (public) ✅
- `GET /api/files/:id/details` - Lấy file chi tiết (auth) ✅
- `PUT /api/files/:id` - Cập nhật file ✅
- `DELETE /api/files/:id` - Xóa file ✅
- `GET /api/files/:id/download` - Download file ✅
- `POST /api/files/share` - Chia sẻ file ✅
- `GET /api/files/:id/shares` - Lấy danh sách shares ✅
- `DELETE /api/files/shares/:id` - Xóa share ✅
- `GET /api/conversations/:id/files` - Lấy files trong conversation ✅

#### WebSocket ✅

- `WS /api/ws/connect` - WebSocket connection cho real-time chat ✅
- `GET /api/ws/users/online` - Lấy danh sách users online ✅
- `GET /api/ws/users/:user_id/status` - Lấy trạng thái user ✅

## 🛠️ Development Commands

```bash
# Xem tất cả lệnh có sẵn
make help

# Khởi động services
make docker-up

# Dừng services
make docker-down

# Xem logs
make docker-logs

# Download dependencies
make deps

# Build ứng dụng
make build

# Chạy ứng dụng
make run

# Clean build artifacts
make clean

# Development mode (docker-up + deps + run)
make dev

# Restart services
make restart

# Database migrations
make migrate
```

## 📊 Performance Metrics

### Current Performance

- **Response Time**: ~1ms cho health check, ~10ms cho auth operations
- **Database Connection**: Pool size 10-100 connections
- **Redis Connection**: Pool size 10 connections
- **Memory Usage**: ~28MB cho binary
- **Logging**: Structured JSON với Zap
- **Token Blacklisting**: Immediate revocation (< 1ms)

### Monitoring

- **Health Check**: Real-time service status
- **Request Logging**: Method, path, status, latency
- **Error Logging**: Structured error tracking
- **Database Logging**: Query performance
- **Redis Logging**: Operation tracking
- **Auth Logging**: Login attempts, token operations, activity tracking

## 🤝 Đóng góp

1. Fork dự án
2. Tạo feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Mở Pull Request

## 📚 API Documentation

### 🔐 Authentication APIs

#### Authentication Flow

```bash
# 1. Register new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "StrongPass123!",
    "display_name": "Test User"
  }'

# 2. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "StrongPass123!"
  }'

# 3. Use access token for protected routes
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/users/me

# 4. Logout (blacklists tokens)
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "YOUR_REFRESH_TOKEN"}'
```

### 👥 User Management APIs

```bash
# Search users
curl "http://localhost:8080/api/users/search?q=test"

# Update profile
curl -X PUT http://localhost:8080/api/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "display_name": "Updated Name",
    "bio": "New bio"
  }'

# Change password
curl -X PUT http://localhost:8080/api/users/me/password \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "StrongPass123!",
    "new_password": "NewStrongPass456!"
  }'
```

### 👥 Friend System APIs

```bash
# Send friend request
curl -X POST http://localhost:8080/api/friends/requests \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 123}'

# Accept friend request
curl -X PUT http://localhost:8080/api/friends/requests/456/accept \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Get friends list
curl -X GET http://localhost:8080/api/friends \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Block user
curl -X POST http://localhost:8080/api/friends/block \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 789}'
```

### 💬 Conversation APIs

```bash
# Create conversation
curl -X POST http://localhost:8080/api/conversations \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Group Chat",
    "type": "group",
    "participant_ids": [123, 456, 789]
  }'

# Get conversations
curl -X GET http://localhost:8080/api/conversations \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Get specific conversation
curl -X GET http://localhost:8080/api/conversations/10 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Add participant
curl -X POST http://localhost:8080/api/conversations/10/participants \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 999, "role": "member"}'

# Leave conversation
curl -X POST http://localhost:8080/api/conversations/10/leave \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 💬 Message APIs

```bash
# Create message
curl -X POST http://localhost:8080/api/conversations/10/messages \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello everyone!",
    "message_type": "text"
  }'

# Get messages
curl -X GET http://localhost:8080/api/conversations/10/messages \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Get messages before ID
curl -X GET "http://localhost:8080/api/conversations/10/messages/before?before_id=50&limit=20" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Search messages
curl -X GET "http://localhost:8080/api/conversations/10/messages/search?q=hello" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Update message
curl -X PUT http://localhost:8080/api/conversations/10/messages/123 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "Updated message"}'

# Delete message
curl -X DELETE http://localhost:8080/api/conversations/10/messages/123 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Add reaction
curl -X POST http://localhost:8080/api/conversations/10/messages/123/reactions \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"reaction_type": "like"}'

# Remove reaction
curl -X DELETE http://localhost:8080/api/conversations/10/messages/123/reactions/like \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 📁 File APIs

````bash
# Upload file
curl -X POST http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "file=@document.pdf" \
  -F "conversation_id=10" \
  -F "is_public=false"

# Get user files
curl -X GET http://localhost:8080/api/files/my \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Search files
curl -X GET "http://localhost:8080/api/files/search?query=document&page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Get file details
curl -X GET http://localhost:8080/api/files/123/details \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Get public file
curl -X GET http://localhost:8080/api/files/123

# Share file
curl -X POST http://localhost:8080/api/files/share \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "file_id": 123,
    "shared_with": 456,
    "can_download": true,
    "can_edit": false
  }'

# Get conversation files
curl -X GET http://localhost:8080/api/conversations/10/files \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

### 🔌 WebSocket APIs

```bash
# Connect to WebSocket (with JWT token in query parameter)
wscat -c "ws://localhost:8080/api/ws/connect?token=YOUR_ACCESS_TOKEN"

# Get online users
curl -X GET http://localhost:8080/api/ws/users/online \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Get user status
curl -X GET http://localhost:8080/api/ws/users/123/status \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
````

#### WebSocket Events

**Client to Server:**

```json
// Join conversation
{
  "type": "join_conversation",
  "data": {
    "conversation_id": 10
  },
  "timestamp": "2025-08-26T14:00:00.000Z"
}

// Send typing indicator
{
  "type": "typing",
  "data": {
    "conversation_id": 10
  },
  "timestamp": "2025-08-26T14:00:00.000Z"
}

// Stop typing
{
  "type": "stop_typing",
  "data": {
    "conversation_id": 10
  },
  "timestamp": "2025-08-26T14:00:00.000Z"
}
```

**Server to Client:**

```json
// New message received
{
  "type": "new_message",
  "data": {
    "id": 123,
    "conversation_id": 10,
    "sender_id": 456,
    "sender_name": "testuser1",
    "content": "Hello everyone!",
    "message_type": "text",
    "created_at": "2025-08-26T14:00:00.000Z"
  },
  "timestamp": "2025-08-26T14:00:00.000Z"
}

// User typing indicator
{
  "type": "user_typing",
  "data": {
    "conversation_id": 10
  },
  "timestamp": "2025-08-26T14:00:00.000Z",
  "user_id": 456,
  "username": "testuser1"
}

// User joined conversation
{
  "type": "user_joined",
  "data": {
    "conversation_id": 10,
    "user_id": 456,
    "username": "testuser1"
  },
  "timestamp": "2025-08-26T14:00:00.000Z"
}

// User online status
{
  "type": "user_online",
  "data": {
    "user_id": 456,
    "username": "testuser1",
    "is_online": true,
    "timestamp": "2025-08-26T14:00:00.000Z"
  },
  "timestamp": "2025-08-26T14:00:00.000Z"
}

// User offline status
{
  "type": "user_offline",
  "data": {
    "user_id": 456,
    "username": "testuser1",
    "is_online": false,
    "timestamp": "2025-08-26T14:00:00.000Z"
  },
  "timestamp": "2025-08-26T14:00:00.000Z"
}

## 🧪 Testing & API Examples

### Security Features

- ✅ **Token Blacklisting**: Immediate revocation after logout
- ✅ **Rate Limiting**: Login attempt protection
- ✅ **Password Strength**: Validation và hashing
- ✅ **Session Management**: Database & Redis sessions
- ✅ **Activity Logging**: Complete audit trail
- ✅ **File Access Control**: Public/private files, ownership validation
- ✅ **File Validation**: Size limits, type restrictions

### 🧪 Testing Results

#### ✅ **Authentication System**

- User registration, login, logout tested
- JWT token generation and validation working
- Token blacklisting functional
- Password reset flow tested

#### ✅ **User Management**

- User CRUD operations tested
- Profile updates working
- User search functionality tested
- Password change tested

#### ✅ **Friend System**

- Friend requests (send, accept, reject, cancel) tested
- Friendships creation and management tested
- User blocking/unblocking tested
- All validation and error cases tested

#### ✅ **Conversation System**

- Conversation creation (direct/group) tested
- Participant management (add/remove/leave) tested
- Admin transfer logic tested (hybrid approach)
- Access control and validation tested

#### ✅ **Message System**

- Message CRUD operations tested
- Message reactions (add/remove) tested
- Message search functionality tested
- Message history and pagination tested
- Access control (only participants can access) tested
- Message sender validation tested
- New participants can see old messages tested
- Multiple users can add same reaction type tested
- One user per reaction type (Facebook-style) tested
- Reaction persistence in database tested

#### ✅ **WebSocket System**

- WebSocket connection and authentication tested
- Real-time message broadcasting tested
- Typing indicators tested
- Online/offline status tracking tested
- Room-based messaging tested
- Client/hub management tested
- JWT token authentication via query parameter tested
- Connection health checker tested (automatic offline detection)
- Real-time status broadcasting tested (online/offline events)

#### ✅ **File System**

- File upload/download tested
- File metadata management tested
- File sharing system tested
- Access control (public/private) tested
- File validation (size, type) tested
- Presigned URLs tested
- File search functionality tested
- Conversation file isolation tested
- Error handling tested
- MinIO integration tested

## 📄 License

Dự án này được phát hành dưới MIT License - xem file [LICENSE](LICENSE) để biết thêm chi tiết.

## 👨‍💻 Tác giả

**Your Name** - [your-email@example.com](mailto:your-email@example.com)

---

⭐ Nếu dự án này hữu ích, hãy cho một star!
```
