# Authentication Feature Implementation Plan

## 🎯 **Feature Overview**
Implement email/password authentication with JWT tokens stored in HTTP-only cookies for a multi-tenant B2B SaaS platform.

### **Requirements Summary:**
- Email/password signup/login (local provider)
- JWT access token with 3-day expiration
- HTTP-only cookies with SameSite strict
- Password validation (min 6 chars, uppercase, number, special char)
- bcrypt hashing with cost 10
- No email verification required
- Separation of auth and user domains

---

## 🏗️ **Implementation Steps**

### **Phase 1: Foundation & Configuration**

#### **Step 1: Environment Configuration**
- [x] Update `pkg/setting/section.go` - Add JWT and Cookie config structs
- [x] Update `configs/local.yaml` - Add JWT and cookie configuration
- [ ] Update `global/global.go` - Add global variables if needed

```yaml
# configs/local.yaml additions needed:
jwt:
  secret: "your-jwt-secret-key"
  expiration_time: "72h"  # 3 days

cookie:
  domain: ""  # Empty for localhost, set via env
  secure: false  # true for production
  http_only: true
  same_site: "Strict"
```

#### **Step 2: Constants Definition**
- [x] Create `internal/common/auth_constants.go` - Auth-related constants
- [x] Create `internal/common/user_constants.go` - User status/role constants  
- [x] Create `internal/common/validation_constants.go` - Validation rules
- [x] Update `internal/common/httpError.go` - Add auth error constants

#### **Step 3: Utility Functions**
- [x] Create `pkg/utils/password.go` - Password hashing/validation utilities (moved from internal/common)
- [x] Create `pkg/utils/jwt.go` - JWT generation/validation utilities (moved from internal/common)
- [x] Removed `internal/common/validation.go` - Using validator library with DTOs instead
- [x] Write comprehensive unit tests for utilities

---

### **Phase 2: Data Layer**

#### **Step 4: Repository Interfaces**
- [x] Create `internal/repo/interfaces.go` - Define all repository interfaces
- [x] Define `UserRepositoryInterface` with required methods
- [x] Plan for future repository interfaces

#### **Step 5: User Repository Implementation**
- [x] Create `internal/repo/user_repository.go` - Implement UserRepository
- [x] Implement `CreateUserWithAuth()` method with transaction
- [x] Implement `GetUserByEmail()` method
- [x] Implement `GetUserByID()` method  
- [x] Implement `GetUserWithWorkspaces()` method with relations
- [x] Write comprehensive repository unit tests with SQLite in-memory database

---

### **Phase 3: Business Logic Layer**

#### **Step 6: DTOs Definition**
- [x] Create `internal/dto/auth_dto.go` - Auth request/response DTOs
  - [x] `SignupRequest` struct with validation tags
  - [x] `LoginRequest` struct with validation tags
  - [x] `AuthResponse` struct for API responses
- [x] Create `internal/dto/user_dto.go` - User-related DTOs
  - [x] `UserWithWorkspaces` struct
  - [x] `WorkspaceMembershipInfo` struct
  - [x] `UserResponse` struct

#### **Step 7: Auth Service Implementation**
- [x] Create `internal/services/auth_service.go` - Authentication business logic
- [x] Implement `AuthServiceInterface` interface
- [x] Implement `Signup()` method with transaction handling
  - [x] Password strength validation
  - [x] Email uniqueness check
  - [x] User + UserAuthProvider creation
  - [x] UserProfile creation
- [x] Implement `Login()` method
  - [x] Credential validation
  - [x] JWT token generation
- [x] Implement `ValidateToken()` method
- [ ] Write service unit tests

#### **Step 8: User Service Implementation**
- [x] Create `internal/services/user_service.go` - User management business logic
- [x] Implement `UserServiceInterface` interface
- [x] Implement `GetUserWithWorkspaces()` method
- [x] Implement `GetUserProfile()` method
- [ ] Write service unit tests

---

### **Phase 4: API Layer**

#### **Step 9: Auth Controller Implementation**
- [x] Create `internal/controllers/auth_controller.go` - Auth endpoints
- [x] Implement `NewAuthController()` constructor
- [x] Implement `Signup()` endpoint handler
  - [x] Request parsing and validation
  - [x] Service layer integration
  - [x] Success response formatting
- [x] Implement `Login()` endpoint handler
  - [x] Credential validation
  - [x] JWT cookie setting
  - [x] Success response formatting
- [x] Implement `Logout()` endpoint handler
  - [x] Cookie clearing

#### **Step 10: User Controller Implementation**
- [x] Create `internal/controllers/user_controller.go` - User endpoints
- [x] Implement `NewUserController()` constructor
- [x] Implement `GetCurrentUser()` endpoint handler
  - [x] JWT token extraction from context
  - [x] User info retrieval with workspaces
  - [x] Response formatting (direct model return)

#### **Step 11: Authentication Middleware**
- [x] Create `internal/middlewares/auth_middleware.go` - JWT validation middleware
- [x] Implement JWT cookie extraction
- [x] Implement token validation
- [x] Implement user ID context injection
- [x] Handle authentication errors
- [x] Additional middleware variants (optional auth, direct token auth)

---

### **Phase 5: Routing & Integration**

#### **Step 12: Auth Routes Setup**
- [x] Create `internal/router/routes/auth_routes.go` - Auth route definitions
- [x] Implement `AuthRoutes` struct and interface
- [x] Define route prefix `/auth`
- [x] Setup routes:
  - [x] `POST /auth/signup`
  - [x] `POST /auth/login`
  - [x] `POST /auth/logout`
- [x] Integrated dependency injection

#### **Step 13: User Routes Setup**
- [x] Create `internal/router/routes/user_routes.go` - User route definitions
- [x] Implement `UserRoutes` struct and interface
- [x] Define route prefix `/users`
- [x] Setup routes:
  - [x] `GET /users/me` (with auth middleware)
- [x] Integrated dependency injection and middleware

#### **Step 14: Route Registration**
- [x] Update `internal/router/router.go` - Register new route modules
- [x] Add `AuthRoutes` to route manager
- [x] Add `UserRoutes` to route manager
- [x] Successful compilation test completed

---