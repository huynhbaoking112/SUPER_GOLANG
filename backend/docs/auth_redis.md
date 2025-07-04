# Redis Token Caching Implementation Plan

## 🎯 **Feature Design Overview**

### **Problem Statement**
- IAM service hiện tại sẽ trở thành bottleneck khi có nhiều microservices cần validate tokens
- Mỗi request đến services khác đều phải gọi qua IAM để check RBAC
- Cần giảm load cho IAM service và tăng performance validation

### **Solution Architecture**
- **Redis Caching Layer**: Cache token data để avoid repeated IAM calls
- **Encrypted Token Approach**: Encrypt JWT tokens làm session identifier
- **Direct Redis Access**: Services khác gọi trực tiếp Redis thay vì IAM

### **Core Flow Design**

#### **Login Process:**
1. User login → IAM validate credentials
2. Generate JWT token → Encrypt JWT → `encrypted_token`
3. Build RBAC data từ database (global_role + workspace_memberships)
4. Store Redis: `auth:token:{user_id}:{encrypted_token}` → RBAC data
5. Set dual cookies: `access_token` (JWT) + `encrypted_token`
6. Return full user data to client

#### **Service Request Validation:**
1. Extract `encrypted_token` từ cookie + `user_id` từ params
2. Direct Redis lookup: `auth:token:{user_id}:{encrypted_token}`
3. If found → check permissions → allow/deny
4. If not found → 401 Unauthorized

#### **Logout Process:**
1. Delete specific Redis key: `auth:token:{user_id}:{encrypted_token}`
2. Clear cookies

### **Redis Structure**
```
Key: "auth:token:{user_id}:{encrypted_token}"
Value: {
  "global_role": "customer",
  "workspace_memberships": [
    {
      "workspace_id": "ws_abc123",
      "role_name": "admin",
      "permissions": ["user:*", "document:*", "workspace:*"],
      "status": "active"
    }
  ]
}
TTL: 72 hours
```

### **Benefits**
- **Performance**: Giảm latency validation từ ~50ms xuống ~2ms
- **Scalability**: IAM service không bị bottleneck
- **Security**: Token encryption + automatic expiration
- **Simplicity**: Services chỉ cần Redis lookup

---

## 📋 **Phase 1: Foundation Setup** ✅ **COMPLETED**

### **Step 1: Update Configuration**
- [x] Add encryption key to `pkg/setting/section.go` JWT struct
- [x] Update `configs/local.yaml` with encryption key (32 bytes for AES-256)
- [x] Verify Redis connection configuration in `internal/initialize/redis.go`

### **Step 2: Create Encryption Utilities**
- [x] Create `pkg/utils/encryption.go`
- [x] Implement AES-256-GCM encrypt function
- [x] Implement AES-256-GCM decrypt function
- [x] Add error handling for encryption failures
- [x] Add unit tests for encryption functions

### **Step 3: Create Token DTOs**
- [x] Create `internal/dto/token_dto.go`
- [x] Define `UserTokenData` struct (global_role + workspace_memberships)
- [x] Define `WorkspaceMembershipTokenData` struct (workspace_id, role_name, permissions, status)
- [x] Add JSON marshal/unmarshal tags
- [x] Add helper methods for permission checking

---

## 📋 **Phase 2: Auth Service Enhancement** ✅ **COMPLETED**

### **Step 4: Enhance Auth Service với Redis Operations**
- [x] Add Redis token storage method to `internal/services/auth_service.go`
- [x] Add Redis token retrieval method
- [x] Add Redis token deletion method
- [x] Add build RBAC data method (query workspace memberships + permissions)
- [x] Implement Redis pipeline for better performance

### **Step 5: Update Login Logic**
- [x] Modify login method to generate encrypted token
- [x] Build RBAC data after successful authentication
- [x] Store token data in Redis with proper TTL
- [x] Update response to include encrypted_token cookie
- [x] Maintain backward compatibility with existing JWT cookie

### **Step 6: Update Logout Logic**
- [x] Modify logout method to delete Redis token
- [x] Clear both access_token and encrypted_token cookies
- [x] Handle logout errors gracefully

---
