# 🚀 Huddle - Real-time Chat Application

Huddle là một ứng dụng chat realtime hiện đại, lấy cảm hứng từ Messenger, được xây dựng bằng Golang và PostgreSQL. Ứng dụng hỗ trợ chat 1-1, nhóm chat, kết bạn, chia sẻ file và nhiều tính năng khác.

## 📋 Mục lục

- [Tính năng](#-tính-năng)
- [Kiến trúc hệ thống](#-kiến-trúc-hệ-thống)
- [Cấu trúc dự án](#-cấu-trúc-dự-án)
- [Flow hoạt động](#-flow-hoạt-động)
- [Công nghệ sử dụng](#-công-nghệ-sử-dụng)
- [Roadmap phát triển](#-roadmap-phát-triển)
- [Cài đặt và chạy](#-cài-đặt-và-chạy)
- [API Documentation](#-api-documentation)

## ✨ Tính năng

### 🔐 Authentication & User Management

- Đăng ký, đăng nhập, đăng xuất
- JWT token authentication
- Quản lý profile người dùng
- Upload avatar
- Tìm kiếm người dùng

### 👥 Friend System

- Gửi lời mời kết bạn
- Chấp nhận/từ chối lời mời
- Danh sách bạn bè
- Chặn/bỏ chặn người dùng
- Quản lý lời mời kết bạn

### 💬 Chat Features

- **Chat 1-1**: Tin nhắn riêng tư giữa 2 người
- **Group Chat**: Chat nhóm với nhiều thành viên
- **Real-time messaging**: WebSocket cho tin nhắn tức thì
- **Message history**: Lưu trữ và tìm kiếm tin nhắn
- **Message reactions**: Like, heart, emoji reactions
- **Read receipts**: Hiển thị trạng thái đã đọc
- **Typing indicators**: Hiển thị đang gõ

### 📁 File Sharing

- Upload và chia sẻ file
- Hỗ trợ nhiều định dạng file
- Lưu trữ file trên MinIO
- Preview hình ảnh
- Download file

### 🏢 Group Management

- Tạo nhóm chat
- Thêm/xóa thành viên
- Phân quyền admin/member
- Quản lý thông tin nhóm
- Avatar nhóm

### 🔔 Notifications

- Push notifications
- Thông báo tin nhắn mới
- Thông báo lời mời kết bạn
- Online/offline status

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
│   │       └── main.go                 # Entry point
│   ├── internal/
│   │   ├── auth/                       # Authentication module
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── user/                       # User management module
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── friend/                     # Friend system module
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── chat/                       # Chat module
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── websocket.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── group/                      # Group management module
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── file/                       # File handling module
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── minio.go
│   │   │   ├── model.go
│   │   │   └── interface.go
│   │   ├── database/                   # Database connection
│   │   │   ├── connection.go
│   │   │   ├── migrations.go
│   │   │   └── models.go
│   │   ├── middleware/                 # HTTP middleware
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── rate_limit.go
│   │   └── routes/                     # Route aggregation
│   │       └── routes.go
│   ├── pkg/
│   │   ├── utils/                      # Utility functions
│   │   │   ├── jwt.go
│   │   │   ├── password.go
│   │   │   ├── response.go
│   │   │   └── validator.go
│   │   └── config/                     # Configuration
│   │       └── config.go
│   ├── migrations/                     # Database migrations
│   │   ├── 001_initial_schema.sql
│   │   ├── 002_add_friendship.sql
│   │   └── 003_add_file_storage.sql
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── components/                 # React components
│   │   ├── pages/                      # Page components
│   │   ├── services/                   # API services
│   │   ├── store/                      # State management
│   │   └── utils/                      # Frontend utilities
│   ├── public/
│   ├── package.json
│   └── README.md
├── docker-compose.yml                  # Docker setup
├── Dockerfile.backend
├── Dockerfile.frontend
└── README.md
```

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

### Backend

- **Golang** (1.21+) - Ngôn ngữ lập trình chính
- **Gin** - HTTP web framework
- **Gorilla WebSocket** - Real-time communication
- **GORM** - ORM cho database
- **PostgreSQL** - Relational database
- **Redis** - Cache và session storage
- **MinIO** - Object storage cho file
- **JWT** - Authentication tokens
- **bcrypt** - Password hashing

### Frontend

- **React** (18+) - UI framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling framework
- **Socket.io-client** - WebSocket client
- **React Query** - Data fetching
- **React Router** - Client-side routing
- **Zustand** - State management

### DevOps & Tools

- **Docker** - Containerization
- **Docker Compose** - Multi-container setup
- **Git** - Version control
- **Postman** - API testing
- **pgAdmin** - Database management

## 🗺️ Roadmap phát triển

### Phase 1: Foundation (2-3 tuần)

**Mục tiêu**: Setup cơ bản và authentication

- [ ] Setup project structure
- [ ] Database schema và migrations
- [ ] User authentication (register/login)
- [ ] Basic WebSocket connection
- [ ] Simple chat interface
- [ ] Docker setup

### Phase 2: Core Features (2-3 tuần)

**Mục tiêu**: Tính năng chat cơ bản

- [ ] Friend system (request, accept, reject)
- [ ] Direct messaging
- [ ] Message history
- [ ] Online status
- [ ] Basic UI/UX
- [ ] File upload với MinIO

### Phase 3: Advanced Features (2-3 tuần)

**Mục tiêu**: Tính năng nâng cao

- [ ] Group creation và management
- [ ] Group messaging
- [ ] Message reactions
- [ ] Push notifications
- [ ] Search functionality
- [ ] User profiles và avatars

### Phase 4: Polish & Optimization (1-2 tuần)

**Mục tiêu**: Hoàn thiện và tối ưu

- [ ] Error handling
- [ ] Performance optimization
- [ ] Security improvements
- [ ] Testing (unit, integration)
- [ ] Documentation
- [ ] Deployment setup

### Phase 5: Enhancement (1-2 tuần)

**Mục tiêu**: Tính năng bổ sung

- [ ] Voice messages
- [ ] Video calls (future)
- [ ] Message encryption
- [ ] Advanced search
- [ ] Message forwarding
- [ ] Mobile responsive

## 🚀 Cài đặt và chạy

### Yêu cầu hệ thống

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Quick Start với Docker

```bash
# Clone repository
git clone https://github.com/your-username/huddle.git
cd huddle

# Chạy với Docker Compose
docker-compose up -d

# Truy cập ứng dụng
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# MinIO Console: http://localhost:9001
```

### Development Setup

```bash
# Backend
cd backend
go mod download
go run cmd/server/main.go

# Frontend
cd frontend
npm install
npm start
```

## 📚 API Documentation

### Authentication Endpoints

- `POST /api/auth/register` - Đăng ký
- `POST /api/auth/login` - Đăng nhập
- `POST /api/auth/logout` - Đăng xuất
- `GET /api/auth/me` - Lấy thông tin user hiện tại

### User Endpoints

- `GET /api/users` - Lấy danh sách users
- `GET /api/users/:id` - Lấy thông tin user
- `PUT /api/users/profile` - Cập nhật profile
- `POST /api/users/avatar` - Upload avatar

### Friend Endpoints

- `GET /api/friends` - Lấy danh sách bạn bè
- `POST /api/friends/request/:user_id` - Gửi lời mời kết bạn
- `PUT /api/friends/request/:request_id` - Phản hồi lời mời
- `GET /api/friends/requests` - Lấy danh sách lời mời

### Chat Endpoints

- `GET /api/messages/direct/:user_id` - Lấy tin nhắn 1-1
- `POST /api/messages/direct/:user_id` - Gửi tin nhắn 1-1
- `GET /api/groups/:id/messages` - Lấy tin nhắn nhóm
- `POST /api/groups/:id/messages` - Gửi tin nhắn nhóm

### Group Endpoints

- `GET /api/groups` - Lấy danh sách nhóm
- `POST /api/groups` - Tạo nhóm mới
- `POST /api/groups/:id/members` - Thêm thành viên
- `DELETE /api/groups/:id/members/:user_id` - Xóa thành viên

### File Endpoints

- `POST /api/files/upload` - Upload file
- `GET /api/files/:filename` - Download file
- `DELETE /api/files/:filename` - Xóa file

### WebSocket

- `WS /ws` - WebSocket connection cho real-time chat

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
