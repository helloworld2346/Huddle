# 🚀 Quick Setup Guide - Huddle

Hướng dẫn setup nhanh để chạy PostgreSQL và Redis cho dự án Huddle.

## 📋 Yêu cầu

- Docker & Docker Compose
- Go 1.24.6+
- TablePlus (hoặc tool quản lý database khác)

## ⚡ Quick Start

### 1. Khởi động services

```bash
# Khởi động PostgreSQL và Redis
make docker-up

# Hoặc sử dụng docker-compose trực tiếp
docker-compose up -d
```

### 2. Download dependencies

```bash
make deps
```

### 3. Chạy ứng dụng

```bash
make run
```

## 🔧 Kết nối Database

### PostgreSQL

- **Host**: `localhost`
- **Port**: `5432`
- **Database**: `huddle`
- **Username**: `huddle_user`
- **Password**: `huddle_password`

### Redis

- **Host**: `localhost`
- **Port**: `6379`
- **Password**: (không có)
- **Database**: `0`

## 📊 Kết nối với TablePlus

### PostgreSQL Connection

1. Mở TablePlus
2. Tạo connection mới
3. Chọn PostgreSQL
4. Điền thông tin:
   - Host: `localhost`
   - Port: `5432`
   - Database: `huddle`
   - User: `huddle_user`
   - Password: `huddle_password`

### Redis Connection (nếu cần)

1. Tạo connection mới
2. Chọn Redis
3. Điền thông tin:
   - Host: `localhost`
   - Port: `6379`

## 🛠️ Các lệnh hữu ích

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

# Chạy tests
make test

# Clean build artifacts
make clean

# Development mode (docker-up + deps + run)
make dev

# Restart services
make restart
```

## 🔍 Kiểm tra services

### Kiểm tra PostgreSQL

```bash
# Kiểm tra container
docker ps | grep postgres

# Kiểm tra logs
docker logs huddle_postgres

# Kết nối trực tiếp
docker exec -it huddle_postgres psql -U huddle_user -d huddle
```

### Kiểm tra Redis

```bash
# Kiểm tra container
docker ps | grep redis

# Kiểm tra logs
docker logs huddle_redis

# Kết nối trực tiếp
docker exec -it huddle_redis redis-cli
```

## 🚨 Troubleshooting

### PostgreSQL không khởi động

```bash
# Kiểm tra port có bị chiếm không
lsof -i :5432

# Restart container
docker restart huddle_postgres
```

### Redis không khởi động

```bash
# Kiểm tra port có bị chiếm không
lsof -i :6379

# Restart container
docker restart huddle_redis
```

### Go dependencies lỗi

```bash
# Clean và download lại
go clean -modcache
make deps
```

## 📝 Notes

- Database sẽ được tự động tạo khi container khởi động
- Migrations sẽ chạy tự động từ thư mục `migrations/`
- Data được lưu trong Docker volumes nên không mất khi restart
- Có thể xóa volumes để reset data: `docker-compose down -v`

## 🎯 Next Steps

Sau khi setup thành công:

1. ✅ PostgreSQL và Redis đã sẵn sàng
2. 🔄 Tiếp theo sẽ setup Gin server
3. 🔐 Implement authentication
4. 💬 Setup WebSocket cho chat
5. 👥 Implement friend system
