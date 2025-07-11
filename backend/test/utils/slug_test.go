package utils_test

import (
	"fmt"
	"go-backend-v2/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple workspace name",
			input:    "My Workspace",
			expected: "my-workspace",
		},
		{
			name:     "Name with special characters",
			input:    "My@Awesome#Workspace!",
			expected: "myawesomeworkspace",
		},
		{
			name:     "Name with multiple spaces",
			input:    "My    Awesome    Workspace",
			expected: "my-awesome-workspace",
		},
		{
			name:     "Name with leading and trailing spaces",
			input:    "  My Workspace  ",
			expected: "my-workspace",
		},
		{
			name:     "Name with numbers",
			input:    "Workspace 123",
			expected: "workspace-123",
		},
		{
			name:     "Name with hyphens already",
			input:    "My-Awesome-Workspace",
			expected: "my-awesome-workspace",
		},
		{
			name:     "Name with multiple consecutive hyphens",
			input:    "My---Workspace",
			expected: "my-workspace",
		},
		{
			name:     "Single word",
			input:    "Workspace",
			expected: "workspace",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "workspace",
		},
		{
			name:     "Only special characters",
			input:    "@#$%^&*()",
			expected: "workspace",
		},
		{
			name:     "Mixed case with underscores",
			input:    "MyAwesomeWorkspace",
			expected: "myawesomeworkspace",
		},
		{
			name:     "Name with dots",
			input:    "My.Awesome.Workspace",
			expected: "myawesomeworkspace",
		},
		{
			name:     "Very long name",
			input:    "This Is A Very Long Workspace Name That Contains Many Words",
			expected: "this-is-a-very-long-workspace-name-that-contains-many-words",
		},
		{
			name:     "Name with accented characters",
			input:    "Café Workspace",
			expected: "caf-workspace",
		},
		{
			name:     "Name with only hyphens",
			input:    "---",
			expected: "workspace",
		},
		{
			name:     "Vietnamese workspace name",
			input:    "Không gian làm việc",
			expected: "khng-gian-lm-vic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GenerateSlug(tt.input)
			assert.Equal(t, tt.expected, result, "GenerateSlug(%q) = %q, want %q", tt.input, result, tt.expected)
		})
	}
}

func TestGenerateUniqueSlug_NoConflict(t *testing.T) {
	checkExists := func(slug string) (bool, error) {
		return false, nil
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple name without conflict",
			input:    "My Workspace",
			expected: "my-workspace",
		},
		{
			name:     "Complex name without conflict",
			input:    "My@Awesome#Workspace!",
			expected: "myawesomeworkspace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.GenerateUniqueSlug(tt.input, checkExists)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateUniqueSlug_WithConflicts(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		conflictingSlugs []string
		expected         string
	}{
		{
			name:             "Single conflict",
			input:            "My Workspace",
			conflictingSlugs: []string{"my-workspace"},
			expected:         "my-workspace-1",
		},
		{
			name:             "Multiple conflicts",
			input:            "My Workspace",
			conflictingSlugs: []string{"my-workspace", "my-workspace-1", "my-workspace-2"},
			expected:         "my-workspace-3",
		},
		{
			name:             "Conflicts with gaps",
			input:            "My Workspace",
			conflictingSlugs: []string{"my-workspace", "my-workspace-1", "my-workspace-3"},
			expected:         "my-workspace-2",
		},
		{
			name:             "Many conflicts",
			input:            "Popular Name",
			conflictingSlugs: generateSlugSequence("popular-name", 10),
			expected:         "popular-name-10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a map for quick lookup
			conflicts := make(map[string]bool)
			for _, slug := range tt.conflictingSlugs {
				conflicts[slug] = true
			}

			checkExists := func(slug string) (bool, error) {
				return conflicts[slug], nil
			}

			result, err := utils.GenerateUniqueSlug(tt.input, checkExists)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateUniqueSlug_CheckExistsError(t *testing.T) {
	// Mock function that returns an error
	checkExists := func(slug string) (bool, error) {
		return false, fmt.Errorf("database connection failed")
	}

	result, err := utils.GenerateUniqueSlug("My Workspace", checkExists)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "failed to check slug existence")
}

func TestGenerateUniqueSlug_TooManyConflicts(t *testing.T) {
	// Mock function that always returns true (everything conflicts)
	checkExists := func(slug string) (bool, error) {
		return true, nil
	}

	result, err := utils.GenerateUniqueSlug("My Workspace", checkExists)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "unable to generate unique slug")
}

func TestGenerateUniqueSlug_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string input",
			input:    "",
			expected: "workspace",
		},
		{
			name:     "Only special characters",
			input:    "@#$%",
			expected: "workspace",
		},
	}

	// Mock function that never conflicts
	checkExists := func(slug string) (bool, error) {
		return false, nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.GenerateUniqueSlug(tt.input, checkExists)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests
func BenchmarkGenerateSlug(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.GenerateSlug("My Awesome Workspace With Many Words")
	}
}

// Helper functions
func generateSlugSequence(baseSlug string, count int) []string {
	result := []string{baseSlug}
	for i := 1; i < count; i++ {
		result = append(result, fmt.Sprintf("%s-%d", baseSlug, i))
	}
	return result
}

// Table-driven test for parallel execution
func TestGenerateSlug_Parallel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Parallel test 1", "Workspace One", "workspace-one"},
		{"Parallel test 2", "Workspace Two", "workspace-two"},
		{"Parallel test 3", "Workspace Three", "workspace-three"},
		{"Parallel test 4", "Workspace Four", "workspace-four"},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := utils.GenerateSlug(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
