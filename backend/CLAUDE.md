
## Environment Variables and Global Configuration Guidelines

When working with environment variables and configuration in this codebase, follow these patterns:

## 1. Config Structure Pattern

### Defining Config Struct
All config structs must be defined in `pkg/setting/section.go`:

```go
// Add new struct to section.go
type NewService struct {
Host string `mapstructure:"host"`
Port int `mapstructure:"port"`
Username string `mapstructure:"username"`
Password string `mapstructure:"password"`
}

// Update Config struct to include new service
type Config struct {
Mysql Mysql `mapstructure:"mysql"`
Redis Redis `mapstructure:"redis"`
Server Server `mapstructure:"server"`
NewService NewService `mapstructure:"newservice"`
}
```

### Naming rules for Config Struct:
- Use PascalCase for struct names
- Use mapstructure tags with lowercase or camelCase
- Fields must have clear types (string, int, bool, time.Duration)

## 2. Global Variables Pattern

### Declaring Global Variables
All global variables must be declared in `global/global.go`:

```go
package global

import (
"go-backend-v2/pkg/setting"
"gorm.io/gorm"
"github.com/redis/go-redis/v9"
)

var (
Config *setting.Config
DB *gorm.DB // Database connection
RedisClient *redis.Client // Redis connection
Logger *zap.Logger // Logger instance
// Add new global variables here
NewServiceClient *NewServiceClient
)
```

### Rules for Global Variables:

- Use PascalCase for variable names

- Always use pointer types for complex objects

- Group related variables together

- Add comments to explain the purpose of the variable

## 3. Config Loading Pattern

### Viper Configuration
The `internal/initialize/loadconfig.go` file handles config loading:

```go
func LoadConfig() {
viper := viper.New()
viper.AddConfigPath("configs")
viper.SetConfigName("local") // Or based on ENV
viper.SetConfigType("yaml")

// Support environment variables
viper.AutomaticEnv()
viper.SetEnvPrefix("APP") // APP_MYSQL_HOST, APP_REDIS_PORT, etc.

if err := viper.ReadInConfig(); err != nil { 
panic(fmt.Errorf("ERROR READING CONFIG FILE: %v", err)) 
} 

if err := viper.Unmarshal(&global.Config); err != nil { 
panic(fmt.Errorf("ERROR UNMARSHALLING CONFIG: %v", err)) 
}
}
```

## 4. Service Initialization Pattern

### Create Init Functions
Each service needs an init function in `internal/initialize/`:

```go
// File: internal/initialize/newservice.go
package initialize

import ( 
"go-backend-v2/global" 
"fmt"
)

func InitNewService() { 
cfg := global.Config.NewService 

client := NewServiceClient{ 
Host: cfg.Host, 
Port: cfg.Port, 
Username: cfg.Username, 
Password: cfg.Password, 
} 

if err := client.Connect(); err != nil { 
panic(fmt.Errorf("failed to connect to NewService: %v", err)) 
} 

global.NewServiceClient = &client
}
```

### Update Run Function
Add init function to `internal/initialize/run.go`:

```go
func Run() { 
LoadConfig() 
InitMysql() 
InitRedis() 
InitLogger() 
InitNewService() // Add new init function
}
```

## 5. Configuration File Pattern

### YAML Configuration Structure
File `configs/local.yaml` must follow structure:

```yaml
server: 
port: 8080 
host: "localhost"

mysql: 
host: "localhost" 
port: 3306 
user: "root" 
password: "root" 
dbName: "IAM"
maxIdleConns: 10
maxOpenConns: 100
connMaxLifetime: "10m"
connMaxIdleTime: "10m"

redis:
host: "localhost"
port: 6379
password: "root"

newservice:
host: "localhost"
port: 9090
username: "admin"
password: "secret"
```

When you need to add a new configuration, follow the correct sequence:

1. Add struct to `pkg/setting/section.go`
2. Add field to Config struct
3. Add values to `configs/local.yaml`
4. Create init function in `internal/initialize/`
5. Add global variable to `global/global.go`
6. Update `Run()` function

# API Development Rules & Conventions

## Architecture Overview

Our API follows **Clean Architecture** principles with modular router design:

```
backend/
├── internal/
│   ├── controllers/           # Request handlers (presentation layer)
│   ├── services/             # Business logic (use case layer)  
│   ├── repo/                 # Data access (repository layer)
|   |── common/               # common data ex: constant message error
│   ├── models/               # Domain entities
│   ├── middlewares/          # Cross-cutting concerns
│   └── router/
│       ├── router.go         # Main router orchestrator
│       └── routes/           # Modular route definitions
│           ├── route_interface.go
│           ├── auth_routes.go
│           ├── user_routes.go
│           ├── admin_routes.go
│           ├── member_routes.go
│           └── public_routes.go
└── cmd/server/main.go        # Application entry point
```

## Router Module Rules

### 1. **Route Module Structure**

Every route module MUST implement the `RouteModule` interface:

```go
type RouteModule interface {
    SetupRoutes(router fiber.Router)
    GetPrefix() string
}
```

### 2. **File Naming Convention**

- Route files: `{domain}_routes.go` (e.g., `user_routes.go`, `admin_routes.go`)
- Controller files: `{domain}_controller.go` 
- Service files: `{domain}_service.go`
- Repository files: `{domain}_repository.go`

### 3. **Route Module Template**

```go
package routes

import (
    "go-backend-v2/internal/controllers"
    "go-backend-v2/internal/middlewares"
    "github.com/gofiber/fiber/v2"
)

type {Domain}Routes struct {
    controller *controllers.{Domain}Controller
}

func New{Domain}Routes() *{Domain}Routes {
    return &{Domain}Routes{
        controller: controllers.New{Domain}Controller(),
    }
}

func (r *{Domain}Routes) GetPrefix() string {
    return "/{domain_prefix}"
}

func (r *{Domain}Routes) SetupRoutes(router fiber.Router) {
    group := router.Group(r.GetPrefix())
    
    // Apply middleware if needed
    group.Use(middlewares.AuthMiddleware())
    
    // Define routes
    group.Get("/", r.controller.List)
    group.Post("/", r.controller.Create)
    group.Get("/:id", r.controller.GetByID)
    group.Put("/:id", r.controller.Update)
    group.Delete("/:id", r.controller.Delete)
}
```


## Controller Rules

### 1. **Controller Structure**

```go
type {Domain}Controller struct {
    service *services.{Domain}Service
}

func New{Domain}Controller() *{Domain}Controller {
    return &{Domain}Controller{
        service: services.New{Domain}Service(),
    }
}
```

### 2. **Handler Method Signature**

```go
func (c *{Domain}Controller) MethodName(ctx *fiber.Ctx) error {
    // 1. Input validation
    // 2. Extract parameters
    // 3. Call service layer
    // 4. Format response
    // 5. Error handling
}
```

### 3. **Standard CRUD Operations**

| Method | Route | Handler | Purpose |
|--------|-------|---------|---------|
| GET | `/{resource}` | `List()` | Get all resources |
| POST | `/{resource}` | `Create()` | Create new resource |
| GET | `/{resource}/:id` | `GetByID()` | Get specific resource |
| PUT | `/{resource}/:id` | `Update()` | Update resource |
| DELETE | `/{resource}/:id` | `Delete()` | Delete resource |

### 4. **Response Format Standards**

#### Success Response:
```go
return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "message": "Operation successful",
    "data":    result,
    "meta": fiber.Map{
        "timestamp": time.Now(),
        "path":      c.Path(),
    },
})
```

#### Error Response:
```go
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "error":   "Detailed error message",
    "code":    "ERROR_CODE",
    "path":    c.Path(),
    "details": validationErrors, // optional
})
```

#### Paginated Response:
```go
return c.JSON(fiber.Map{
    "data": items,
    "pagination": fiber.Map{
        "page":       page,
        "limit":      limit,
        "total":      total,
        "totalPages": totalPages,
    },
})
```

## Route Registration Rules

### 1. **Main Router Setup**

```go
// router/router.go
func SetupRoutes(app *fiber.App) {
    manager := NewRouteManager()
    manager.SetupRoutes(app)
}

type RouteManager struct {
    config       *routes.RouteConfig
    routeModules []routes.RouteModule
}
```

### 2. **Adding New Route Module**

1. Create controller: `internal/controllers/{domain}_controller.go`
2. Create routes: `internal/router/routes/{domain}_routes.go`
3. Register in RouteManager:

```go
func NewRouteManager() *RouteManager {
    routeModules := []routes.RouteModule{
        routes.NewPublicRoutes(),
        routes.NewAuthRoutes(),
        routes.NewUserRoutes(),
        routes.NewAdminRoutes(),
        routes.NewMemberRoutes(),
        routes.New{NewDomain}Routes(), // Add here
    }
    // ...
}
```

## API Versioning

### 1. **Version Structure**

```
/api/v1/{resource}
/api/v2/{resource}
```

### 2. **Version Configuration**

```go
type RouteConfig struct {
    APIVersion string // "v1", "v2", etc.  
    BaseURL    string // "/api"
}
```

### 2. **Global Error Handler**

```go
app := fiber.New(fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        code := fiber.StatusInternalServerError
        if e, ok := err.(*fiber.Error); ok {
            code = e.Code
        }
        
        return c.Status(code).JSON(fiber.Map{
            "error":  err.Error(),
            "code":   code,
            "path":   c.Path(),
            "method": c.Method(),
        })
    },
})
```