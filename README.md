# Go Commerce - E-commerce API

E-commerce REST API built with Go Fiber framework following Clean Architecture principles. Features include user authentication, product management, transaction processing with atomic operations, comprehensive API documentation, and full Docker containerization.

## Features

### Core Features
- **Authentication & Authorization** - JWT-based auth with role-based access control
- **User Management** - Registration, profile management, password changes
- **Store Management** - Auto store creation with "toko-username" format, store profiles
- **Product Management** - CRUD operations, file upload, pagination, filtering
- **Address Management** - Indonesia region API integration, province/city validation
- **Category Management** - Admin-only category management
- **Transaction System** - Atomic transactions with product logging
- **API Documentation** - Interactive Swagger UI
- **Docker Support** - Full containerization with docker-compose

### Technical Features
- **Clean Architecture** - Separation of concerns with proper layering
- **Database Transactions** - ACID compliance with rollback mechanisms
- **Goroutine Logging** - Asynchronous product snapshot logging
- **File Upload System** - Local file storage with validation
- **Comprehensive Testing** - 50+ unit tests with mock repositories
- **Input Validation** - Request validation with proper error handling
- **Pagination & Filtering** - Efficient data retrieval
- **Docker Containerization** - Production-ready deployment
- **Region Integration** - Indonesia provinces and cities API integration

## Architecture

```
go-commerce/
├── cmd/app/                    # Application entry point
├── internal/
│   ├── domain/                 # Business entities & interfaces
│   ├── usecase/                # Business logic layer
│   ├── repository/mysql/       # Data access layer
│   ├── handler/http/           # HTTP handlers & middleware
│   └── service/                # External service integrations
├── pkg/                        # Shared utilities
│   ├── config/                 # Configuration management
│   ├── database/               # Database connection
│   └── jwt/                    # JWT token management
├── docs/                       # Auto-generated API documentation
├── migrations/                 # Database migration files
├── uploads/                    # File upload storage
├── Dockerfile                  # Docker build configuration
├── docker-compose.yml          # Multi-container orchestration
└── .dockerignore              # Docker build exclusions
```

## Tech Stack

- **Framework**: [Go Fiber v2](https://gofiber.io/) - Fast HTTP web framework
- **Database**: MySQL with [GORM](https://gorm.io/) ORM
- **Authentication**: JWT tokens with bcrypt password hashing
- **Documentation**: [Swagger/OpenAPI 3.0](https://swagger.io/)
- **Testing**: [Testify](https://github.com/stretchr/testify) with mock repositories
- **Validation**: [Go Playground Validator](https://github.com/go-playground/validator)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Containerization**: Docker & Docker Compose

## Prerequisites

### Option 1: Docker (Recommended)
- Docker & Docker Compose
- Git

### Option 2: Local Development
- Go 1.25 or higher
- MySQL 8.0 or higher
- Git

## Quick Start

### 🐳 Docker Setup (Recommended)

```bash
# 1. Clone Repository
git clone <repository-url>
cd go-commerce

# 2. Start with Docker
docker-compose up -d

# 3. Access Application
# API: http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html
# Health: http://localhost:8080/health
```

### 🔧 Local Development Setup

```bash
# 1. Clone Repository
git clone <repository-url>
cd go-commerce

# 2. Install Dependencies
go mod download

# 3. Environment Configuration
cp .env.example .env
# Edit .env with your database credentials

# 4. Database Setup
mysql -u root -p -e "CREATE DATABASE go_commerce;"

# 5. Run Migrations
migrate -path migrations -database "mysql://user:pass@tcp(localhost:3306)/go_commerce" up

# 6. Run Application
go run cmd/app/main.go
```

## Docker Commands

```bash
# Start containers
docker-compose up -d

# View logs
docker-compose logs app
docker-compose logs mysql

# Stop containers
docker-compose down

# Rebuild after code changes
docker-compose up --build -d

# View container status
docker-compose ps
```

## Environment Variables

```env
# Database Configuration
DB_HOST=mysql                    # Use 'mysql' for Docker, 'localhost' for local
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=go_commerce

# Application Configuration
APP_PORT=8080
APP_ENV=development

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRE_HOURS=24
JWT_REFRESH_EXPIRE_HOURS=168

# Upload Configuration
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=5242880
```

## API Documentation

### Swagger UI
Access interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Health Check
```
GET http://localhost:8080/health
```

### Key Endpoints

#### Authentication
- `POST /api/v1/auth/register` - User registration (auto creates "toko-username" store)
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/users/my` - Get profile (protected)

#### Stores
- `GET /api/v1/stores` - Get all active stores (public)
- `GET /api/v1/stores/my` - Get my store (protected)
- `PUT /api/v1/stores/my` - Update my store (protected)

#### Products
- `GET /api/v1/products` - Get all products (public)
- `POST /api/v1/products` - Create product (protected)
- `GET /api/v1/products/{id}` - Get product by ID

#### Addresses
- `GET /api/v1/addresses` - Get my addresses (protected)
- `POST /api/v1/addresses` - Create address with region validation (protected)
- `PUT /api/v1/addresses/{id}` - Update address (protected)
- `GET /api/v1/regions/provinces` - Get all provinces (public)
- `GET /api/v1/regions/provinces/{id}/cities` - Get cities by province (public)

#### Transactions
- `POST /api/v1/transactions` - Create transaction (protected)
- `GET /api/v1/transactions/my` - Get my transactions (protected)

## New Features

### Auto Store Creation
When users register, a store is automatically created with:
- **Store Name Format**: `toko-username` (e.g., "toko-johndoe")
- **Status**: Active (no pending approval required)
- **Description**: Welcome message

### Address Management with Indonesia Region API
Full integration with Indonesia region data:
- **Province & City Validation**: Real-time validation using Indonesia API
- **Auto-populate Names**: Province and city names automatically filled
- **Required Fields**: `province_id` and `city_id` are mandatory
- **String Format**: Uses string IDs to match API format (e.g., "31", "3171")

#### Address Request Example:
```json
{
  "judul_alamat": "Rumah Utama",
  "nama_penerima": "Ahmad Rizki",
  "notelp": "081234567890",
  "detail_alamat": "Jl. Sudirman No. 100, Senayan",
  "province_id": "31",
  "city_id": "3171",
  "kode_pos": "10270",
  "is_default": true
}
```

#### Address Response Example:
```json
{
  "id": 1,
  "judul_alamat": "Rumah Utama",
  "nama_penerima": "Ahmad Rizki",
  "detail_alamat": "Jl. Sudirman No. 100, Senayan",
  "province_id": "31",
  "city_id": "3171",
  "province_name": "DKI Jakarta",
  "city_name": "Jakarta Pusat",
  "kode_pos": "10270",
  "is_default": true
}
```

## Testing

### Run All Tests
```bash
go test ./... -v
```

### Run Specific Package Tests
```bash
go test ./internal/usecase -v
go test ./pkg/jwt -v
```

### Test Coverage
```bash
go test ./... -cover
```

**Current Test Coverage**: 50+ tests covering all critical business logic

## Development

### Generate Swagger Documentation
```bash
swag init -g cmd/app/main.go
```

### Database Migrations
Migrations are handled automatically by GORM AutoMigrate on application startup.

### Project Structure Guidelines
- **Domain Layer**: Business entities and interfaces (no dependencies)
- **Usecase Layer**: Business logic and orchestration
- **Repository Layer**: Data access and persistence
- **Handler Layer**: HTTP request/response handling
- **Service Layer**: External service integrations

## Security Features

- **JWT Authentication** with secure token generation
- **Password Hashing** using bcrypt
- **Input Validation** on all endpoints
- **Ownership Validation** for user-specific resources
- **Role-Based Access Control** for admin operations
- **File Upload Validation** with type and size checks
- **Region Data Validation** using external Indonesia API
- **Auto Store Activation** (no manual approval required)

## Deployment

### Docker Production Deployment

```bash
# Build production image
docker build -t go-commerce:latest .

# Run with production environment
docker run -d \
  --name go-commerce-app \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-secure-password \
  -e JWT_SECRET=your-production-jwt-secret \
  go-commerce:latest
```

### Docker Compose Production

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=mysql
      - DB_PASSWORD=secure-production-password
      - JWT_SECRET=super-secure-production-jwt-key
    depends_on:
      - mysql
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=secure-production-password
      - MYSQL_DATABASE=go_commerce
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

### Traditional Build

```bash
# Build for current platform
go build -o go-commerce cmd/app/main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o go-commerce-linux cmd/app/main.go
```

## Development Workflow

### With Docker (Recommended)

```bash
# 1. Make code changes
# 2. Rebuild and restart
docker-compose up --build -d

# 3. View logs
docker-compose logs -f app

# 4. Test changes
curl http://localhost:8080/health
```

### Local Development

```bash
# 1. Make code changes
# 2. Run tests
go test ./... -v

# 3. Start application
go run cmd/app/main.go

# 4. Generate Swagger docs (if API changes)
swag init -g cmd/app/main.go
```

## Performance Features

- **Database Connection Pooling** - Optimized connection management
- **Goroutine Logging** - Asynchronous operations for better performance
- **Pagination** - Efficient data retrieval for large datasets
- **Database Indexing** - Optimized query performance including region fields
- **Static File Serving** - Efficient file delivery
- **External API Caching** - Optimized region data retrieval
- **Auto Store Creation** - Streamlined user onboarding process

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Authors

- **Your Name** - *Initial work* - [YourGitHub](https://github.com/yourusername)

## Acknowledgments

- [Go Fiber](https://gofiber.io/) - Amazing web framework
- [GORM](https://gorm.io/) - Fantastic ORM for Go
- [Swagger](https://swagger.io/) - API documentation standard
- Clean Architecture principles by Robert C. Martin

---

**Built with Go and Clean Architecture**