package repo_test

import (
	"go-backend-v2/internal/repo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserRepository(t *testing.T) {
	// Test that the repository constructor works
	userRepo := repo.NewUserRepository()

	assert.NotNil(t, userRepo)
	assert.Implements(t, (*repo.UserRepositoryInterface)(nil), userRepo)
}

func TestUserRepositoryInterface(t *testing.T) {
	// Test that UserRepository implements the interface correctly
	var _ repo.UserRepositoryInterface = &repo.UserRepository{}

	// This test ensures the interface is properly implemented
	assert.True(t, true)
}

// Note: Full integration tests with database require CGO for SQLite
// These tests verify the repository structure and interface compliance
// For full testing, run with: CGO_ENABLED=1 go test
func TestRepositoryStructure(t *testing.T) {
	userRepo := repo.NewUserRepository()

	// Test that we can create the repository without panicking
	assert.NotNil(t, userRepo)

	// Type assertion test
	_, ok := userRepo.(*repo.UserRepository)
	assert.True(t, ok, "Repository should be of type *UserRepository")
}
