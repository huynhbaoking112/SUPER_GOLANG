# Go Backend Project

Dự án backend Go này tuân thủ theo Clean Architecture và [Standard Go Project Layout](https://github.com/golang-standards/project-layout) để đảm bảo tính nhất quán, khả năng mở rộng và dễ bảo trì.

## Cấu trúc thư mục

```
backend/
├── cmd/                    # Application entrypoints
│   ├── server/            # Main API server
│   ├── cli/               # Command line interface
│   └── cronjob/           # Background job runner
├── internal/               # Private application code
│   ├── controllers/       # HTTP handlers & controllers
│   ├── services/          # Business logic & use cases
│   ├── repo/              # Data access layer & repositories
│   ├── models/            # Domain entities & data models
│   ├── router/            # Route definitions & middleware setup
│   ├── middlewares/       # Custom middleware functions
│   └── initialize/        # Application initialization
├── pkg/                   # Shared utility packages
│   ├── logger/            # Logging utilities
│   ├── setting/           # Configuration management
│   └── utils/             # Common helper functions
├── global/                # Global variables & singletons
├── response/              # HTTP response utilities & status codes
├── configs/               # Configuration files & templates
├── docs/                  # Documentation & API specs
├── deployments/           # Deployment configurations
├── scripts/               # Build & deployment scripts
├── test/                  # Integration & E2E tests
└── tools/                 # Development tools & utilities
```

## Giải thích chi tiết

### 🏗️ Thư mục Go cốt lõi

#### `/cmd`
**Mục đích**: Chứa các điểm vào (entry points) của ứng dụng.

- **`server/`**: Main API server application - điểm vào chính cho REST/gRPC API
- **`cli/`**: Command line interface tools - các công cụ dòng lệnh
- **`cronjob/`**: Background job runner - xử lý các tác vụ nền và scheduled jobs

**Nguyên tắc**:
- Code trong đây nên được giữ ở mức tối thiểu
- Chỉ chịu trách nhiệm khởi tạo và "wire" các components
- Mỗi thư mục con tương ứng với một executable riêng biệt

#### `/internal`
**Mục đích**: Chứa toàn bộ business logic và core application code.

**Clean Architecture Layers**:

- **`controllers/`** (Presentation Layer):
  - HTTP handlers và controllers
  - Request/Response transformation
  - Input validation và sanitization
  - Authentication & authorization logic

- **`services/`** (Use Case Layer):
  - Business logic và use cases
  - Application-specific rules
  - Orchestration between different domains
  - Transaction management

- **`repo/`** (Data Access Layer):
  - Repository interfaces và implementations
  - Database queries và operations
  - External API integrations
  - Data mapping và transformation

- **`models/`** (Domain Layer):
  - Domain entities và value objects
  - Business rules và domain logic
  - Core data structures
  - Database schemas

- **`router/`** (Infrastructure):
  - Route definitions và grouping
  - Middleware setup và configuration
  - HTTP server configuration

- **`middlewares/`** (Cross-cutting Concerns):
  - Authentication middleware
  - Logging middleware
  - CORS handling
  - Rate limiting
  - Error handling

- **`initialize/`** (Dependency Injection):
  - Application bootstrap code
  - Dependency injection setup
  - Database connections
  - External service initialization

#### `/pkg`
**Mục đích**: Shared utilities có thể được sử dụng bởi nhiều services.

- **`logger/`**: Structured logging utilities với multiple output formats
- **`setting/`**: Configuration loading và validation từ files/environment
- **`utils/`**: Common helper functions và shared utilities

**Nguyên tắc**: Chỉ đặt code ở đây khi chắc chắn có thể tái sử dụng.

### 🚀 Thư mục ứng dụng đặc thù

#### `/global`
**Mục đích**: Global variables, singletons và shared state.

- Global database connections
- Shared configuration instances
- Application-wide constants
- Singleton instances

**Lưu ý**: Sử dụng cẩn thận để tránh tight coupling.

#### `/response`
**Mục đích**: HTTP response utilities và standardized response formats.

- HTTP status code constants
- Standardized API response structures
- Error response formatting
- Success response helpers

### 📚 Thư mục hỗ trợ

#### `/configs`
**Mục đích**: Configuration files và templates.

- Environment-specific configurations (dev, staging, prod)
- Database connection configurations
- External service configurations
- Feature flags và runtime settings

#### `/docs`
**Mục đích**: Project documentation.

- API documentation (OpenAPI/Swagger specs)
- Architecture decision records (ADRs)
- Development guides
- Deployment instructions

#### `/deployments`
**Mục đích**: Deployment configurations cho các môi trường khác nhau.

- Docker configurations (Dockerfile, docker-compose.yml)
- Kubernetes manifests
- CI/CD pipeline definitions
- Infrastructure as Code (Terraform, CloudFormation)

#### `/scripts`
**Mục đích**: Automation scripts.

- Build scripts
- Database migration scripts
- Deployment automation
- Development setup scripts

#### `/test`
**Mục đích**: Integration tests, E2E tests và test utilities.

- Integration test suites
- End-to-end test scenarios
- Test data fixtures
- Test helper functions

**Lưu ý**: Unit tests nên được đặt cùng với code được test (`*_test.go` files).

#### `/tools`
**Mục đích**: Development tools và code generation.

- Mock generation tools
- Code generation utilities
- Database migration tools
- Development helper scripts

## Architecture Principles

### 🎯 Clean Architecture

1. **Dependency Inversion**: Dependencies point inward, từ frameworks về domain
2. **Interface Segregation**: Sử dụng interfaces để decouple components
3. **Single Responsibility**: Mỗi layer có một trách nhiệm rõ ràng  
4. **Separation of Concerns**: Business logic tách biệt khỏi infrastructure

### 🔧 Go Best Practices

1. **Interface-Driven Development**: Define interfaces trong consumer packages
2. **Error Handling**: Explicit error handling với wrapped errors
3. **Context Propagation**: Sử dụng context.Context cho request lifecycle
4. **Goroutine Safety**: Thread-safe access to shared resources
5. **Resource Management**: Proper cleanup với defer statements

## Development Workflow

### 🚀 Getting Started

1. **Clone và setup**:
   ```bash
   git clone <repository>
   cd backend
   go mod download
   ```

2. **Run development server**:
   ```bash
   go run cmd/server/main.go
   ```

3. **Run CLI tools**:
   ```bash
   go run cmd/cli/main.go [command]
   ```

4. **Run background jobs**:
   ```bash
   go run cmd/cronjob/main.go
   ```

### 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test ./test/...

# Run specific package tests
go test ./internal/services/...
```

### 🏗️ Building

```bash
# Build all applications
./scripts/build.sh

# Build specific application
go build -o bin/server cmd/server/main.go
go build -o bin/cli cmd/cli/main.go
go build -o bin/cronjob cmd/cronjob/main.go
```

## Code Organization Guidelines

### 📁 Package Naming
- Sử dụng singular nouns (user, not users)
- Short, clear và descriptive names
- Avoid abbreviations unless commonly understood
- Consistent với Go naming conventions

### 🔗 Dependencies
- **Controllers** depend on **Services**
- **Services** depend on **Repository interfaces**
- **Repositories** implement interfaces defined in **Services**
- **Models** should be dependency-free domain objects

### 📝 File Organization
```
internal/services/user/
├── service.go          # Service interface definition
├── user_service.go     # Service implementation
├── user_service_test.go # Unit tests
└── errors.go           # Service-specific errors
```

## Contributing

1. **Code Style**: Follow Go conventions và use `gofmt`
2. **Testing**: Write tests for new functionality
3. **Documentation**: Update relevant documentation
4. **Error Handling**: Implement proper error handling với context
5. **Logging**: Add appropriate logging cho debugging

## Security Best Practices

- **Input Validation**: Validate all external inputs
- **Authentication**: Implement proper JWT/session handling
- **Authorization**: Role-based access control
- **Rate Limiting**: Protect endpoints from abuse
- **Secure Defaults**: Use secure configuration defaults
- **Error Messages**: Don't leak sensitive information in errors

## Performance Considerations

- **Database**: Use connection pooling và query optimization
- **Caching**: Implement appropriate caching strategies
- **Concurrency**: Use goroutines safely với proper synchronization
- **Memory**: Avoid memory leaks với proper resource cleanup
- **Monitoring**: Implement health checks và metrics

## Resources

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) 