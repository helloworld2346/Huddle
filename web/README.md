# 🚀 Huddle Frontend - Real-time Chat Application

Frontend cho ứng dụng chat realtime Huddle, được xây dựng bằng Next.js 15, TypeScript, và Tailwind CSS.

## 📋 Mục lục

- [Tính năng](#-tính-năng)
- [Công nghệ sử dụng](#-công-nghệ-sử-dụng)
- [Cài đặt và chạy](#-cài-đặt-và-chạy)
- [Development](#-development)

## ✨ Tính năng

### 🔐 Authentication & User Management

- [x] **Đăng ký, đăng nhập, đăng xuất** - Complete auth flow
- [x] **JWT token authentication** - Access & refresh tokens
- [x] **Remember me functionality** - Persistent login
- [x] **Password reset** - Forgot/reset password flow

### 💬 Chat Features

- [x] **Chat 1-1 & Group Chat** - Real-time messaging
- [x] **Message reactions** - Like, heart, emoji reactions
- [x] **Typing indicators** - Real-time typing status
- [x] **Online/Offline status** - User presence tracking
- [x] **Message search** - Search functionality
- [x] **File sharing** - Upload and share files

### 👥 Friend System

- [x] **Friend requests** - Send, accept, reject
- [x] **User search** - Find and add friends
- [x] **Blocking system** - Block/unblock users

### 📁 File Management

- [x] **File upload** - Drag & drop interface
- [x] **File preview** - Image and document preview
- [x] **File sharing** - Share with users/conversations

## 🛠️ Công nghệ sử dụng

### Frontend Core

- **Next.js 15.5.0** - React framework với App Router
- **React 19.1.0** - UI library
- **TypeScript 5.x** - Type safety
- **Tailwind CSS 4.x** - Utility-first CSS framework

### State Management & Data Fetching

- **Zustand 4.5.0** - Lightweight state management
- **TanStack Query 5.0.0** - Data fetching & caching
- **React Hook Form 7.50.0** - Form handling & validation

### Real-time Communication

- **Socket.io Client 4.7.0** - WebSocket client
- **WebSocket API** - Real-time messaging

### UI/UX

- **Lucide React 0.400.0** - Icon library
- **React Hot Toast 2.4.0** - Toast notifications
- **Date-fns 3.6.0** - Date formatting

## 🚀 Cài đặt và chạy

### Yêu cầu hệ thống

- Node.js 18.17.0+
- npm hoặc yarn
- Backend server đang chạy (localhost:8080)

### Quick Start

```bash
# Clone repository (nếu chưa có)
git clone https://github.com/your-username/huddle.git
cd huddle/web

# Install dependencies
npm install

# Setup environment variables
cp .env.example .env.local

# Start development server
npm run dev

# Truy cập ứng dụng
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

### Environment Variables

Tạo file `.env.local`:

```env
# Backend API
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# App Configuration
NEXT_PUBLIC_APP_NAME=Huddle
NEXT_PUBLIC_APP_VERSION=1.0.0
```

## 🛠️ Development

### Available Scripts

```bash
# Development
npm run dev          # Start development server
npm run build        # Build for production
npm run start        # Start production server
npm run lint         # Run ESLint
npm run type-check   # Run TypeScript check
```

### Development Workflow

1. **Start Backend First**

   ```bash
   # Terminal 1 - Backend
   cd ../
   make docker-up
   make run
   ```

2. **Start Frontend**

   ```bash
   # Terminal 2 - Frontend
   cd web
   npm run dev
   ```

3. **Access Application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - MinIO Console: http://localhost:9001

## 📁 Project Structure

```
web/src/
├── app/                          # Next.js App Router
│   ├── (auth)/                   # Auth routes group
│   │   ├── login/page.tsx        # Login page
│   │   ├── register/page.tsx     # Register page
│   │   └── forgot-password/page.tsx # Forgot password
│   ├── (dashboard)/              # Protected routes group
│   │   ├── layout.tsx            # Dashboard layout
│   │   ├── page.tsx              # Main chat interface
│   │   ├── conversations/        # Conversation pages
│   │   ├── friends/              # Friend pages
│   │   ├── profile/              # Profile pages
│   │   └── files/                # File pages
│   ├── layout.tsx                # Root layout
│   └── globals.css               # Global styles
├── components/                   # Reusable components
│   ├── ui/                       # Base UI components
│   ├── auth/                     # Auth components
│   ├── chat/                     # Chat components
│   ├── friends/                  # Friend components
│   ├── files/                    # File components
│   └── layout/                   # Layout components
├── lib/                          # Utilities & configs
├── hooks/                        # Custom hooks
├── stores/                       # Zustand stores
└── types/                        # TypeScript types
```

## 🔌 API Integration

### API Client Setup

```typescript
// lib/api.ts
import axios from "axios";

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  timeout: 10000,
});

// Request interceptor - Add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

### Key API Endpoints

- **Authentication**: `/api/auth/*`
- **Conversations**: `/api/conversations/*`
- **Messages**: `/api/conversations/:id/messages/*`
- **Friends**: `/api/friends/*`
- **Files**: `/api/files/*`
- **WebSocket**: `/api/ws/connect`

## 🗃️ State Management

### Zustand Stores

```typescript
// stores/auth-store.ts
interface AuthStore {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
}

// stores/chat-store.ts
interface ChatStore {
  conversations: Conversation[];
  currentConversation: Conversation | null;
  messages: Record<string, Message[]>;
  onlineUsers: User[];
}
```

## 🎨 Design System

### Color Palette

```css
:root {
  --primary: #3b82f6; /* Blue */
  --primary-dark: #2563eb;
  --secondary: #64748b; /* Slate */
  --success: #10b981; /* Emerald */
  --warning: #f59e0b; /* Amber */
  --error: #ef4444; /* Red */
  --background: #ffffff; /* White */
  --surface: #f8fafc; /* Slate 50 */
  --text: #1e293b; /* Slate 800 */
}
```

## 📱 Responsive Design

### Breakpoints

- **Mobile**: < 768px
- **Tablet**: 768px - 1024px
- **Desktop**: > 1024px

### Layout Adaptations

- **Mobile**: Single column, bottom navigation
- **Tablet**: Sidebar + main content
- **Desktop**: Full sidebar + main content + details panel

## 🚀 Deployment

### Production Build

```bash
# Build for production
npm run build

# Start production server
npm run start
```

### Environment Variables (Production)

```env
# Production
NEXT_PUBLIC_API_URL=https://api.huddle.com
NEXT_PUBLIC_WS_URL=wss://api.huddle.com
NEXT_PUBLIC_APP_NAME=Huddle
NEXT_PUBLIC_APP_VERSION=1.0.0
```

### Deployment Platforms

- **Vercel** (Recommended)
- **Netlify**
- **AWS Amplify**
- **Docker**

## 📚 Additional Resources

### Documentation

- [Next.js Documentation](https://nextjs.org/docs)
- [React Documentation](https://react.dev)
- [TypeScript Documentation](https://www.typescriptlang.org/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Backend API Documentation](../README.md#api-documentation)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

---

⭐ If this project is helpful, please give it a star!
