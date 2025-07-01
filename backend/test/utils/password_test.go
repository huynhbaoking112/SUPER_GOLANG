package utils_test

import (
	"go-backend-v2/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PasswordTestSuite struct {
	suite.Suite
}

func (suite *PasswordTestSuite) TestHashPassword_Success() {
	password := "TestPassword123!"

	hashedPassword, err := utils.HashPassword(password)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
	assert.NotEqual(suite.T(), password, hashedPassword)
	assert.Greater(suite.T(), len(hashedPassword), 50) // bcrypt hashes are typically 60+ chars
}

func (suite *PasswordTestSuite) TestHashPassword_EmptyPassword() {
	password := ""

	hashedPassword, err := utils.HashPassword(password)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
}

func (suite *PasswordTestSuite) TestHashPassword_DifferentPasswords() {
	password1 := "TestPassword123!"
	password2 := "DifferentPassword456@"

	hash1, err1 := utils.HashPassword(password1)
	hash2, err2 := utils.HashPassword(password2)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NotEqual(suite.T(), hash1, hash2)
}

func (suite *PasswordTestSuite) TestHashPassword_SamePasswordDifferentHashes() {
	password := "TestPassword123!"

	hash1, err1 := utils.HashPassword(password)
	hash2, err2 := utils.HashPassword(password)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	// bcrypt generates different salts, so same password should have different hashes
	assert.NotEqual(suite.T(), hash1, hash2)
}

func (suite *PasswordTestSuite) TestCheckPassword_ValidPassword() {
	password := "TestPassword123!"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(suite.T(), err)

	isValid := utils.CheckPassword(password, hashedPassword)

	assert.True(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestCheckPassword_InvalidPassword() {
	password := "TestPassword123!"
	wrongPassword := "WrongPassword456@"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(suite.T(), err)

	isValid := utils.CheckPassword(wrongPassword, hashedPassword)

	assert.False(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestCheckPassword_EmptyPassword() {
	password := "TestPassword123!"
	emptyPassword := ""

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(suite.T(), err)

	isValid := utils.CheckPassword(emptyPassword, hashedPassword)

	assert.False(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestCheckPassword_EmptyHash() {
	password := "TestPassword123!"
	emptyHash := ""

	isValid := utils.CheckPassword(password, emptyHash)

	assert.False(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestCheckPassword_InvalidHash() {
	password := "TestPassword123!"
	invalidHash := "invalid-hash-string"

	isValid := utils.CheckPassword(password, invalidHash)

	assert.False(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestValidatePassword_ValidPassword() {
	validPasswords := []string{
		"Password123!",
		"MySecure@Pass1",
		"ComplexP@ssw0rd",
		"Test123#",
		"Abc123$def",
	}

	for _, password := range validPasswords {
		err := utils.ValidatePassword(password)
		assert.NoError(suite.T(), err, "Password should be valid: %s", password)
	}
}

func (suite *PasswordTestSuite) TestValidatePassword_TooShort() {
	shortPasswords := []string{
		"Abc1!", // 5 chars
		"Ab1!",  // 4 chars
		"A1!",   // 3 chars
		"",      // empty
	}

	for _, password := range shortPasswords {
		err := utils.ValidatePassword(password)
		assert.Error(suite.T(), err, "Password should be invalid (too short): %s", password)
		assert.Contains(suite.T(), err.Error(), "PASSWORD_TOO_SHORT")
	}
}

func (suite *PasswordTestSuite) TestValidatePassword_TooLong() {
	// Create a password longer than 128 characters
	longPassword := "A1!" + string(make([]byte, 130))

	err := utils.ValidatePassword(longPassword)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "PASSWORD_TOO_LONG")
}

func (suite *PasswordTestSuite) TestValidatePassword_NoUppercase() {
	passwordsWithoutUpper := []string{
		"password123!",
		"lowercase@123",
		"nouppercase1#",
	}

	for _, password := range passwordsWithoutUpper {
		err := utils.ValidatePassword(password)
		assert.Error(suite.T(), err, "Password should be invalid (no uppercase): %s", password)
		assert.Contains(suite.T(), err.Error(), "PASSWORD_NO_UPPER")
	}
}

func (suite *PasswordTestSuite) TestValidatePassword_NoLowercase() {
	passwordsWithoutLower := []string{
		"PASSWORD123!",
		"UPPERCASE@123",
		"NOLOWERCASE1#",
	}

	for _, password := range passwordsWithoutLower {
		err := utils.ValidatePassword(password)
		assert.Error(suite.T(), err, "Password should be invalid (no lowercase): %s", password)
		assert.Contains(suite.T(), err.Error(), "PASSWORD_NO_LOWER")
	}
}

func (suite *PasswordTestSuite) TestValidatePassword_NoDigit() {
	passwordsWithoutDigit := []string{
		"Password!",
		"NoNumbers@",
		"OnlyLetters#",
	}

	for _, password := range passwordsWithoutDigit {
		err := utils.ValidatePassword(password)
		assert.Error(suite.T(), err, "Password should be invalid (no digit): %s", password)
		assert.Contains(suite.T(), err.Error(), "PASSWORD_NO_DIGIT")
	}
}

func (suite *PasswordTestSuite) TestValidatePassword_NoSpecialCharacter() {
	passwordsWithoutSpecial := []string{
		"Password123",
		"NoSpecialChars1",
		"OnlyAlphaNum123",
	}

	for _, password := range passwordsWithoutSpecial {
		err := utils.ValidatePassword(password)
		assert.Error(suite.T(), err, "Password should be invalid (no special char): %s", password)
		assert.Contains(suite.T(), err.Error(), "PASSWORD_NO_SPECIAL")
	}
}

func (suite *PasswordTestSuite) TestValidatePassword_EdgeCases() {
	// Test minimum valid password (exactly 6 chars with all requirements)
	minValidPassword := "Aa1!"
	err := utils.ValidatePassword(minValidPassword)
	assert.Error(suite.T(), err) // Should still be too short (4 chars)

	// Test minimum valid password (exactly 6 chars with all requirements)
	minValidPassword = "Aa1!bb"
	err = utils.ValidatePassword(minValidPassword)
	assert.NoError(suite.T(), err) // Should be valid

	// Test various special characters
	specialChars := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "-", "=", "[", "]", "{", "}", ";", "'", ":", "\"", "\\", "|", ",", ".", "<", ">", "/", "?", "~", "`"}

	for _, char := range specialChars {
		password := "Password123" + char
		err := utils.ValidatePassword(password)
		assert.NoError(suite.T(), err, "Password with special char '%s' should be valid", char)
	}
}

func (suite *PasswordTestSuite) TestBcryptCost() {
	// Test that the bcrypt cost is set correctly
	assert.Equal(suite.T(), 10, utils.BcryptCost)
}

func TestPasswordSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}
