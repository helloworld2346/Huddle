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

- [ ] Đăng ký, đăng nhập, đăng xuất
- [ ] JWT token authentication
- [ ] Quản lý profile người dùng
- [ ] Upload avatar
- [ ] Tìm kiếm người dùng

### 👥 Friend System

- [ ] Gửi lời mời kết bạn
- [ ] Chấp nhận/từ chối lời mời
- [ ] Danh sách bạn bè
- [ ] Chặn/bỏ chặn người dùng
- [ ] Quản lý lời mời kết bạn

### 💬 Chat Features

- [ ] **Chat 1-1**: Tin nhắn riêng tư giữa 2 người
- [ ] **Group Chat**: Chat nhóm với nhiều thành viên
- [ ] **Real-time messaging**: WebSocket cho tin nhắn tức thì
- [ ] **Message history**: Lưu trữ và tìm kiếm tin nhắn
- [ ] **Message reactions**: Like, heart, emoji reactions
- [ ] **Read receipts**: Hiển thị trạng thái đã đọc
- [ ] **Typing indicators**: Hiển thị đang gõ

### 📁 File Sharing

- [ ] Upload và chia sẻ file
- [ ] Hỗ trợ nhiều định dạng file
- [ ] Lưu trữ file trên MinIO
- [ ] Preview hình ảnh
- [ ] Download file

### 🏢 Group Management

- [ ] Tạo nhóm chat
- [ ] Thêm/xóa thành viên
- [ ] Phân quyền admin/member
- [ ] Quản lý thông tin nhóm
- [ ] Avatar nhóm

### 🔔 Notifications

- [ ] Push notifications
- [ ] Thông báo tin nhắn mới
- [ ] Thông báo lời mời kết bạn
- [ ] Online/offline status

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

- **Users**: Thông tin người dùng
- **Friend Requests**: Lời mời kết bạn
- **Friendships**: Mối quan hệ bạn bè
- **Groups**: Thông tin nhóm
- **Group Members**: Thành viên nhóm
- **Direct Messages**: Tin nhắn 1-1
- **Group Messages**: Tin nhắn nhóm
- **Message Reactions**: Phản ứng tin nhắn
- **User Sessions**: Phiên đăng nhập
- **Blocked Users**: Người dùng bị chặn

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
│   │   ├── auth/                       # Authentication module ⏳
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── user/                       # User management module ⏳
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── friend/                     # Friend system module ⏳
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── chat/                       # Chat module ⏳
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── websocket.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── group/                      # Group management module ⏳
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── file/                       # File handling module ⏳
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── minio.go
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
│   │   │   └── password.go
│   │   ├── logger/                     # Structured logging ✅
│   │   │   └── logger.go
│   │   ├── utils/                      # Common utilities ✅
│   │   │   └── response.go
│   │   └── validation/                 # Request validation ✅
│   │       └── validator.go
│   ├── migrations/                     # Database migrations ✅
│   │   ├── 001_initial_schema.sql
│   │   └── 002_auth_schema.sql
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
├── .gitignore                         # Git ignore ✅
├── SETUP.md                           # Setup guide ✅
└── README.md
```

## 🎯 Tiến độ phát triển

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

### ⏳ **Đang thực hiện (Phase 2 - Core Features)**

#### **Authentication System:**

- [x] **JWT Utilities** - Token generation, validation, refresh
- [x] **Password Utilities** - Hashing, validation, strength check
- [x] **Authentication Middleware** - JWT validation, user context
- [x] **Request Validation** - Input validation utilities
- [x] **Database Schema** - Auth tables và indexes
- [ ] User registration và login endpoints
- [ ] Session management với Redis

#### **User Management:**

- [ ] User CRUD operations
- [ ] Profile management
- [ ] Avatar upload
- [ ] User search

### 📋 **Còn lại (Phase 3-5)**

#### **Friend System:**

- [ ] Friend requests
- [ ] Friend list management
- [ ] User blocking

#### **Chat Features:**

- [ ] WebSocket setup
- [ ] Direct messaging
- [ ] Group messaging
- [ ] Message history

#### **File Sharing:**

- [ ] MinIO integration
- [ ] File upload/download
- [ ] Image preview

#### **Frontend:**

- [ ] React/Vue.js setup
- [ ] UI components
- [ ] Real-time chat interface

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
User → Upload file → MinIO → Get URL → Save to database → Send message with file
```

### 6. WebSocket Connection Flow

```
Client → Connect WebSocket → Authenticate → Join user room → Listen for messages
```

## 🛠️ Công nghệ sử dụng

### Backend ✅

- **Golang** (1.24.6) - Ngôn ngữ lập trình chính
- **Gin** (v1.10.1) - HTTP web framework
- **Gorilla WebSocket** - Real-time communication ⏳
- **GORM** (v1.30.1) - ORM cho database
- **PostgreSQL** (15-alpine) - Relational database
- **Redis** (7-alpine) - Cache và session storage
- **MinIO** - Object storage cho file ⏳
- **JWT** - Authentication tokens ⏳
- **bcrypt** - Password hashing ⏳
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

### ⏳ Phase 2: Core Features (Đang thực hiện)

**Mục tiêu**: Authentication và user management

- [ ] User authentication (register/login)
- [ ] JWT token management
- [ ] User profile management
- [ ] Basic WebSocket connection
- [ ] File upload với MinIO

### 📋 Phase 3: Advanced Features (Chưa bắt đầu)

**Mục tiêu**: Chat và friend system

- [ ] Friend system (request, accept, reject)
- [ ] Direct messaging
- [ ] Group creation và management
- [ ] Message history
- [ ] Message reactions
- [ ] Online status

### 📋 Phase 4: Polish & Optimization (Chưa bắt đầu)

**Mục tiêu**: Hoàn thiện và tối ưu

- [ ] Push notifications
- [ ] Search functionality
- [ ] Performance optimization
- [ ] Security improvements
- [ ] Testing (unit, integration)
- [ ] Frontend development

### 📋 Phase 5: Enhancement (Chưa bắt đầu)

**Mục tiêu**: Tính năng bổ sung

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

### Quick Start với Docker

```bash
# Clone repository
git clone https://github.com/your-username/huddle.git
cd huddle

# Chạy với Docker Compose
make docker-up

# Download dependencies
make deps

# Chạy ứng dụng
make run

# Truy cập ứng dụng
# Backend API: http://localhost:8080
# Health Check: http://localhost:8080/api/health
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

- `POST /api/auth/register` - Đăng ký
- `POST /api/auth/login` - Đăng nhập
- `POST /api/auth/logout` - Đăng xuất
- `GET /api/auth/me` - Lấy thông tin user hiện tại

#### User Endpoints

- `GET /api/users` - Lấy danh sách users
- `GET /api/users/:id` - Lấy thông tin user
- `PUT /api/users/profile` - Cập nhật profile
- `POST /api/users/avatar` - Upload avatar

#### Friend Endpoints

- `GET /api/friends` - Lấy danh sách bạn bè
- `POST /api/friends/request/:user_id` - Gửi lời mời kết bạn
- `PUT /api/friends/request/:request_id` - Phản hồi lời mời
- `GET /api/friends/requests` - Lấy danh sách lời mời

#### Chat Endpoints

- `GET /api/messages/direct/:user_id` - Lấy tin nhắn 1-1
- `POST /api/messages/direct/:user_id` - Gửi tin nhắn 1-1
- `GET /api/groups/:id/messages` - Lấy tin nhắn nhóm
- `POST /api/groups/:id/messages` - Gửi tin nhắn nhóm

#### WebSocket

- `WS /ws` - WebSocket connection cho real-time chat

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
```

## 📊 Performance Metrics

### Current Performance

- **Response Time**: ~1ms cho health check
- **Database Connection**: Pool size 10-100 connections
- **Redis Connection**: Pool size 10 connections
- **Memory Usage**: ~28MB cho binary
- **Logging**: Structured JSON với Zap

### Monitoring

- **Health Check**: Real-time service status
- **Request Logging**: Method, path, status, latency
- **Error Logging**: Structured error tracking
- **Database Logging**: Query performance
- **Redis Logging**: Operation tracking

## 🤝 Đóng góp

1. Fork dự án
2. Tạo feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Mở Pull Request

## 📄 License

Dự án này được phát hành dưới MIT License - xem file [LICENSE](LICENSE) để biết thêm chi tiết.

## 👨‍💻 Tác giả

**Your Name** - [your-email@example.com](mailto:your-email@example.com)

---

⭐ Nếu dự án này hữu ích, hãy cho một star!
