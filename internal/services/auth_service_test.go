package services

import (
	"errors"
	"testing"
	"time"

	"github.com/baltop/commet/internal/config"
	"github.com/baltop/commet/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

// Test helpers
func newTestJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		Secret:      "test-secret-key-for-testing",
		ExpiryHours: 24,
	}
}

func newTestAuthService(mockRepo *MockUserRepository) *AuthService {
	return NewAuthService(mockRepo, newTestJWTConfig())
}

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

// =============================================================================
// Register Tests
// =============================================================================

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	// Setup expectations
	mockRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Execute
	user, err := authService.Register(req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Name, user.Name)
	assert.NotEmpty(t, user.PasswordHash)
	// Password should be hashed, not plain text
	assert.NotEqual(t, req.Password, user.PasswordHash)

	mockRepo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	// Email already exists
	mockRepo.On("ExistsByEmail", req.Email).Return(true, nil)

	// Execute
	user, err := authService.Register(req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrUserExists, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

func TestRegister_ExistsByEmailError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	dbError := errors.New("database connection error")
	mockRepo.On("ExistsByEmail", req.Email).Return(false, dbError)

	// Execute
	user, err := authService.Register(req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

func TestRegister_CreateError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	createError := errors.New("failed to create user")
	mockRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(createError)

	// Execute
	user, err := authService.Register(req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, createError, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Login Tests
// =============================================================================

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	password := "password123"
	existingUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hashPassword(password),
		Name:         "Test User",
	}

	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	mockRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

	// Execute
	user, token, err := authService.Login(req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, existingUser.Email, user.Email)
	assert.Equal(t, existingUser.ID, user.ID)

	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	req := &models.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", req.Email).Return(nil, errors.New("record not found"))

	// Execute
	user, token, err := authService.Login(req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	existingUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hashPassword("correctpassword"),
		Name:         "Test User",
	}

	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

	// Execute
	user, token, err := authService.Login(req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

// =============================================================================
// ValidateToken Tests
// =============================================================================

func TestValidateToken_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	// First, generate a valid token by logging in
	password := "password123"
	existingUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hashPassword(password),
		Name:         "Test User",
	}

	mockRepo.On("FindByEmail", existingUser.Email).Return(existingUser, nil)

	_, token, _ := authService.Login(&models.LoginRequest{
		Email:    existingUser.Email,
		Password: password,
	})

	// Now validate the token
	claims, err := authService.ValidateToken(token)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, existingUser.ID, claims.UserID)
	assert.Equal(t, existingUser.Email, claims.Email)
	assert.Equal(t, existingUser.Name, claims.Name)
}

func TestValidateToken_Expired(t *testing.T) {
	mockRepo := new(MockUserRepository)
	// Create auth service with very short expiry (negative to simulate expired)
	jwtConfig := config.JWTConfig{
		Secret:      "test-secret-key-for-testing",
		ExpiryHours: -1, // Already expired
	}
	authService := NewAuthService(mockRepo, jwtConfig)

	password := "password123"
	existingUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hashPassword(password),
		Name:         "Test User",
	}

	mockRepo.On("FindByEmail", existingUser.Email).Return(existingUser, nil)

	_, token, _ := authService.Login(&models.LoginRequest{
		Email:    existingUser.Email,
		Password: password,
	})

	// Validate expired token
	claims, err := authService.ValidateToken(token)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestValidateToken_Invalid(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	// Test with completely invalid token
	invalidToken := "invalid.token.string"

	claims, err := authService.ValidateToken(invalidToken)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	// Generate token with one secret
	password := "password123"
	existingUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hashPassword(password),
		Name:         "Test User",
	}

	mockRepo.On("FindByEmail", existingUser.Email).Return(existingUser, nil)

	_, token, _ := authService.Login(&models.LoginRequest{
		Email:    existingUser.Email,
		Password: password,
	})

	// Create new auth service with different secret
	differentSecretService := NewAuthService(mockRepo, config.JWTConfig{
		Secret:      "different-secret-key",
		ExpiryHours: 24,
	})

	// Try to validate with different secret
	claims, err := differentSecretService.ValidateToken(token)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	claims, err := authService.ValidateToken("")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// =============================================================================
// GetUserByID Tests
// =============================================================================

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	expectedUser := &models.User{
		ID:    1,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedUser, nil)

	// Execute
	user, err := authService.GetUserByID(1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("record not found"))

	// Execute
	user, err := authService.GetUserByID(999)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Token Generation Tests
// =============================================================================

func TestGenerateToken_ContainsCorrectClaims(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	password := "password123"
	existingUser := &models.User{
		ID:           42,
		Email:        "test@example.com",
		PasswordHash: hashPassword(password),
		Name:         "Test User",
	}

	mockRepo.On("FindByEmail", existingUser.Email).Return(existingUser, nil)

	_, token, err := authService.Login(&models.LoginRequest{
		Email:    existingUser.Email,
		Password: password,
	})

	assert.NoError(t, err)

	// Validate and check claims
	claims, err := authService.ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
	// Check that expiry is set correctly (should be in the future)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

// =============================================================================
// Password Hashing Tests
// =============================================================================

func TestRegister_PasswordIsHashed(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := newTestAuthService(mockRepo)

	var capturedUser *models.User

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "myplainpassword",
		Name:     "Test User",
	}

	mockRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		capturedUser = args.Get(0).(*models.User)
	}).Return(nil)

	user, err := authService.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Verify password is actually hashed using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(capturedUser.PasswordHash), []byte(req.Password))
	assert.NoError(t, err, "Password should be properly hashed with bcrypt")

	// Verify hash is different from plain password
	assert.NotEqual(t, req.Password, capturedUser.PasswordHash)
}
