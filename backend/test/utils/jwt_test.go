package utils_test

import (
	"go-backend-v2/global"
	"go-backend-v2/pkg/setting"
	"go-backend-v2/pkg/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JWTTestSuite struct {
	suite.Suite
	originalConfig *setting.Config
}

func (suite *JWTTestSuite) SetupSuite() {
	// Backup original config
	suite.originalConfig = global.Config

	// Set test config
	global.Config = &setting.Config{
		JWT: setting.JWT{
			Secret:         "test-jwt-secret-key-for-testing",
			ExpirationTime: 24 * time.Hour,
		},
	}
}

func (suite *JWTTestSuite) TearDownSuite() {
	// Restore original config
	global.Config = suite.originalConfig
}

func (suite *JWTTestSuite) TestGenerateToken_Success() {
	userID := "test-user-123"

	token, err := utils.GenerateToken(userID)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	assert.Contains(suite.T(), token, ".")
}

func (suite *JWTTestSuite) TestGenerateToken_EmptyUserID() {
	userID := ""

	token, err := utils.GenerateToken(userID)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
}

func (suite *JWTTestSuite) TestValidateToken_ValidToken() {
	userID := "test-user-123"

	// Generate a token
	token, err := utils.GenerateToken(userID)
	assert.NoError(suite.T(), err)

	// Validate the token
	extractedUserID, err := utils.ValidateToken(token)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)
}

func (suite *JWTTestSuite) TestValidateToken_InvalidToken() {
	invalidToken := "invalid.token.string"

	extractedUserID, err := utils.ValidateToken(invalidToken)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), extractedUserID)
	assert.Contains(suite.T(), err.Error(), "failed to parse token")
}

func (suite *JWTTestSuite) TestValidateToken_ExpiredToken() {
	// Set very short expiration for test
	global.Config.JWT.ExpirationTime = -1 * time.Hour // Already expired

	userID := "test-user-123"
	token, err := utils.GenerateToken(userID)
	assert.NoError(suite.T(), err)

	// Reset expiration to normal
	global.Config.JWT.ExpirationTime = 24 * time.Hour

	// Validate expired token
	extractedUserID, err := utils.ValidateToken(token)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), extractedUserID)
}

func (suite *JWTTestSuite) TestValidateToken_WrongSecret() {
	userID := "test-user-123"

	// Generate token with current secret
	token, err := utils.GenerateToken(userID)
	assert.NoError(suite.T(), err)

	// Change secret
	originalSecret := global.Config.JWT.Secret
	global.Config.JWT.Secret = "wrong-secret-key"

	// Try to validate with wrong secret
	extractedUserID, err := utils.ValidateToken(token)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), extractedUserID)

	// Restore original secret
	global.Config.JWT.Secret = originalSecret
}

func (suite *JWTTestSuite) TestExtractUserIDFromToken_ValidToken() {
	userID := "test-user-123"

	// Generate a token
	token, err := utils.GenerateToken(userID)
	assert.NoError(suite.T(), err)

	// Extract user ID without validation
	extractedUserID := utils.ExtractUserIDFromToken(token)

	assert.Equal(suite.T(), userID, extractedUserID)
}

func (suite *JWTTestSuite) TestExtractUserIDFromToken_InvalidToken() {
	invalidToken := "invalid.token.string"

	extractedUserID := utils.ExtractUserIDFromToken(invalidToken)

	assert.Empty(suite.T(), extractedUserID)
}

func (suite *JWTTestSuite) TestExtractUserIDFromToken_EmptyToken() {
	extractedUserID := utils.ExtractUserIDFromToken("")

	assert.Empty(suite.T(), extractedUserID)
}

func (suite *JWTTestSuite) TestIsTokenExpired_ValidToken() {
	userID := "test-user-123"

	// Generate a valid token
	token, err := utils.GenerateToken(userID)
	assert.NoError(suite.T(), err)

	// Check if token is expired
	isExpired := utils.IsTokenExpired(token)

	assert.False(suite.T(), isExpired)
}

func (suite *JWTTestSuite) TestIsTokenExpired_ExpiredToken() {
	// Set very short expiration for test
	global.Config.JWT.ExpirationTime = -1 * time.Hour // Already expired

	userID := "test-user-123"
	token, err := utils.GenerateToken(userID)
	assert.NoError(suite.T(), err)

	// Reset expiration to normal
	global.Config.JWT.ExpirationTime = 24 * time.Hour

	// Check if token is expired
	isExpired := utils.IsTokenExpired(token)

	assert.True(suite.T(), isExpired)
}

func (suite *JWTTestSuite) TestIsTokenExpired_InvalidToken() {
	invalidToken := "invalid.token.string"

	isExpired := utils.IsTokenExpired(invalidToken)

	assert.True(suite.T(), isExpired)
}

func (suite *JWTTestSuite) TestJWTClaims_Structure() {
	userID := "test-user-123"

	// Generate a token
	token, err := utils.GenerateToken(userID)
	assert.NoError(suite.T(), err)

	// Parse token to check claims structure
	parsedToken, err := jwt.ParseWithClaims(token, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.Secret), nil
	})

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*utils.JWTClaims)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), userID, claims.UserID)
	assert.Equal(suite.T(), "go-backend-v2", claims.Issuer)
	assert.Equal(suite.T(), userID, claims.Subject)
	assert.NotNil(suite.T(), claims.ExpiresAt)
	assert.NotNil(suite.T(), claims.IssuedAt)
	assert.NotNil(suite.T(), claims.NotBefore)
}

func TestJWTSuite(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}
