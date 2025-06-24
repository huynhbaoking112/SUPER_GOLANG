# Go Backend Project

Dự án backend Go này tuân thủ theo [Standard Go Project Layout](https://github.com/golang-standards/project-layout) để đảm bảo tính nhất quán, khả năng mở rộng và dễ bảo trì.

## Cấu trúc thư mục

```
backend/
├── api/                    # Định nghĩa API
├── build/                  # Packaging và CI/CD
│   ├── ci/                # CI/CD configurations
│   └── package/           # Package configurations
├── cmd/                    # Application entrypoints
│   └── server/            # Main server application
├── configs/                # Configuration files
├── deployments/            # Deployment configurations
├── docs/                   # Documentation
├── examples/               # Examples và demos
├── githooks/               # Git hooks
├── internal/               # Private application code
│   ├── app/               # Use cases/Services
│   ├── domain/            # Domain entities
│   └── repository/        # Data access layer
├── pkg/                    # Public library code
├── scripts/                # Automation scripts
├── test/                   # Additional tests
├── tools/                  # Supporting tools
```

## Giải thích chi tiết

### 🏗️ Thư mục Go cốt lõi

#### `/cmd`
**Mục đích**: Chứa các điểm vào (entry points) của ứng dụng.

- Mỗi thư mục con tương ứng với một executable
- Tên thư mục con nên trùng với tên file thực thi
- Code trong đây nên được giữ ở mức tối thiểu
- Chỉ chịu trách nhiệm khởi tạo và "wire" các components

**Ví dụ**:
```
cmd/
├── server/         # API server
├── worker/         # Background worker
└── migrate/        # Database migration tool
```

#### `/internal`
**Mục đích**: Chứa code riêng tư của ứng dụng.

- Code trong đây không thể được import bởi các project khác
- Đây là quy tắc được Go compiler thực thi
- Chứa toàn bộ business logic chính của ứng dụng

**Cấu trúc con**:
- `internal/app/`: Use cases, services, business logic
- `internal/domain/`: Domain entities, models, business rules
- `internal/repository/`: Data access layer, database interfaces

#### `/pkg`
**Mục đích**: Chứa code có thể được sử dụng bởi các ứng dụng khác.

- Code công khai (public) có thể được import
- Chỉ đặt code ở đây khi chắc chắn muốn chia sẻ
- Quy tắc: Bắt đầu với `/internal`, chỉ chuyển sang `/pkg` khi cần thiết

**Ví dụ**:
```
pkg/
├── logger/         # Logging utilities
├── validator/      # Validation helpers
└── httpclient/     # HTTP client wrapper
```


### 🚀 Thư mục ứng dụng

#### `/api`
**Mục đích**: Chứa các file định nghĩa API contract.

- OpenAPI/Swagger specifications cho REST API
- Protocol Buffer files (`.proto`) cho gRPC
- API documentation và schemas

**Ví dụ**:
```
api/
├── openapi/
│   └── api.yaml    # OpenAPI specification
├── proto/
│   └── user.proto  # gRPC definitions
└── docs/           # API documentation
```

#### `/configs`
**Mục đích**: Chứa configuration files và templates.

- Configuration templates hoặc default configs
- Environment-specific configurations
- Tách biệt cấu hình khỏi code

**Ví dụ**:
```
configs/
├── config.yaml     # Default configuration
├── dev.yaml        # Development config
├── prod.yaml       # Production config
└── docker.yaml     # Docker environment config
```

#### `/build`
**Mục đích**: Packaging và Continuous Integration.

- `build/package/`: Docker, OS package configurations
- `build/ci/`: CI/CD configurations và scripts

**Ví dụ**:
```
build/
├── package/
│   ├── Dockerfile
│   └── docker-compose.yml
└── ci/
    ├── .github/
    └── jenkins/
```

#### `/deployments`
**Mục đích**: IaaS, PaaS, system và container orchestration deployment configurations.

- Kubernetes manifests
- Terraform configurations
- Docker Compose files
- Helm charts

**Ví dụ**:
```
deployments/
├── kubernetes/
├── terraform/
├── helm/
└── docker-compose/
```

### 📚 Thư mục hỗ trợ

#### `/docs`
**Mục đích**: Design và user documents.

- Architecture documentation
- API documentation
- User guides
- Development guides

#### `/scripts`
**Mục đích**: Scripts để thực hiện các operations khác nhau.

- Build scripts
- Installation scripts
- Database migration scripts
- Analysis scripts

**Ví dụ**:
```
scripts/
├── build.sh        # Build application
├── migrate.sh      # Database migration
├── test.sh         # Run tests
└── deploy.sh       # Deployment script
```

#### `/test`
**Mục đích**: Additional external test apps và test data.

- Integration tests
- End-to-end tests
- Test data và fixtures
- Performance tests

**Lưu ý**: Unit tests nên được đặt trong files `*_test.go` bên cạnh code được test.

#### `/tools`
**Mục đích**: Supporting tools cho project.

- Code generation tools
- Development utilities
- Build tools
- Analysis tools

**Ví dụ**:
```
tools/
├── mockgen/        # Mock generation
├── swagger/        # API doc generation
└── migrate/        # Database migration tool
```

#### `/examples`
**Mục đích**: Examples cho applications và/hoặc public libraries.

- Usage examples
- Demo applications
- Sample configurations
- Tutorials

#### `/githooks`
**Mục đích**: Git hooks.

- Pre-commit hooks
- Pre-push hooks
- Commit message validation
- Code quality checks

## Best Practices

### 🎯 Nguyên tắc tổ chức code

1. **Separation of Concerns**: Mỗi thư mục có một trách nhiệm rõ ràng
2. **Dependency Direction**: Dependencies nên point inward (từ ngoài vào trong)
3. **Interface Segregation**: Sử dụng interfaces để decouple components
4. **Single Responsibility**: Mỗi package nên có một lý do duy nhất để thay đổi

### 🔧 Development Workflow

1. **Bắt đầu với `/internal`**: Đặt tất cả code mới trong `/internal` trước
2. **Extract to `/pkg`**: Chỉ chuyển sang `/pkg` khi cần chia sẻ
3. **Keep `/cmd` minimal**: Chỉ chứa initialization code
4. **Document everything**: Maintain clear documentation trong `/docs`

### 📝 Naming Conventions

- Sử dụng tên thư mục ngắn gọn, rõ ràng
- Tránh viết tắt không rõ nghĩa
- Consistent với Go naming conventions
- Sử dụng singular forms cho package names

## Getting Started

1. **Initialize Go module**:
   ```bash
   go mod init your-project-name
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

4. **Run tests**:
   ```bash
   go test ./...
   ```

## Contributing

1. Tuân thủ cấu trúc thư mục hiện tại
2. Viết tests cho code mới
3. Update documentation khi cần thiết
4. Follow Go best practices và coding standards

## Resources

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Clean Architecture in Go](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) 