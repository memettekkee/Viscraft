package service

import (
	"strings"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/request"
	"viscraft-backend/model/response"
	"viscraft-backend/pkg/logger"
	"viscraft-backend/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ImageFinder defines the minimal interface for finding images by user ID.
// Used during user deletion to locate files for filesystem cleanup.
type ImageFinder interface {
	FindImagesByUserId(userId string) ([]ImageRecord, error)
}

// ImageRecord holds the minimal fields needed for cleanup during user deletion.
type ImageRecord struct {
	Id string
}

// StorageDeleter defines the minimal interface for deleting image files from storage.
type StorageDeleter interface {
	Delete(imageId string) error
}

// UserService handles user-related business logic including registration,
// authentication, and account deletion.
type UserService struct {
	userRepo  *repository.UserRepository
	imageFinder ImageFinder
	storage     StorageDeleter
	jwtSecret string
	jwtExpiry time.Duration
}

// NewUserService creates a new UserService with the required dependencies.
// imageFinder and storage can be nil if filesystem cleanup is not needed (they are
// checked before use during deletion).
func NewUserService(
	userRepo *repository.UserRepository,
	imageFinder ImageFinder,
	storage StorageDeleter,
	jwtSecret string,
	jwtExpiry time.Duration,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		imageFinder: imageFinder,
		storage:     storage,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
	}
}

// CreateUser validates the input, creates a new user, and returns a JWT token.
func (s *UserService) CreateUser(requestId string, req request.CreateUserRequest) (response.CreateUserResponse, *constant.AppError) {
	// Validate email format: must contain exactly one "@" with non-empty local and domain parts
	if err := validateEmail(req.Email); err != nil {
		logger.Warn(requestId, "email validation failed", "email", req.Email)
		return response.CreateUserResponse{}, err
	}

	// Validate email length
	if len(req.Email) > 255 {
		logger.Warn(requestId, "email exceeds 255 characters")
		return response.CreateUserResponse{}, &constant.ErrValidationFailed
	}

	// Validate password length: 8-72 characters
	if len(req.Password) < 8 || len(req.Password) > 72 {
		logger.Warn(requestId, "password length invalid")
		return response.CreateUserResponse{}, &constant.ErrValidationFailed
	}

	// Insert user via repository (handles bcrypt hashing)
	user, err := s.userRepo.Insert(req.Email, req.Password, req.Name)
	if err != nil {
		// Check for duplicate email (PostgreSQL unique violation)
		if isDuplicateKeyError(err) {
			logger.Warn(requestId, "duplicate email", "email", req.Email)
			return response.CreateUserResponse{}, &constant.ErrDuplicateResource
		}
		logger.Error(requestId, "failed to insert user", err)
		return response.CreateUserResponse{}, &constant.ErrDatabaseFailed
	}

	// Generate JWT token
	token, tokenErr := s.generateJWT(user.Id)
	if tokenErr != nil {
		logger.Error(requestId, "failed to generate JWT", tokenErr)
		return response.CreateUserResponse{}, &constant.ErrInternalServer
	}

	logger.Info(requestId, "user created", "userId", user.Id)

	return response.CreateUserResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "User created successfully",
		},
		Token: token,
		Data: &response.UserData{
			Id:        user.Id,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// Login authenticates a user by email and password, returning a JWT token on success.
func (s *UserService) Login(requestId string, req request.LoginRequest) (response.LoginResponse, *constant.AppError) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Don't reveal whether email or password was wrong (requirement 2.2)
		logger.Warn(requestId, "login failed - user not found", "email", req.Email)
		return response.LoginResponse{}, &constant.ErrUnauthorized
	}

	// Compare bcrypt hash with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// Don't reveal whether email or password was wrong (requirement 2.2)
		logger.Warn(requestId, "login failed - password mismatch", "userId", user.Id)
		return response.LoginResponse{}, &constant.ErrUnauthorized
	}

	// Generate JWT token
	token, tokenErr := s.generateJWT(user.Id)
	if tokenErr != nil {
		logger.Error(requestId, "failed to generate JWT", tokenErr)
		return response.LoginResponse{}, &constant.ErrInternalServer
	}

	logger.Info(requestId, "user logged in", "userId", user.Id)

	return response.LoginResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Login successful",
		},
		Token: token,
		Data: &response.UserData{
			Id:        user.Id,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// GetUser retrieves the authenticated user's profile data by their ID.
func (s *UserService) GetUser(requestId string, userId string) (response.GetUserResponse, *constant.AppError) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		logger.Warn(requestId, "user not found", "userId", userId)
		return response.GetUserResponse{}, &constant.ErrUnauthorized
	}

	logger.Info(requestId, "user retrieved", "userId", user.Id)

	return response.GetUserResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "User retrieved successfully",
		},
		Data: &response.UserData{
			Id:        user.Id,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// DeleteUser removes a user account and triggers filesystem cleanup for all
// associated images. The userId parameter comes from the JWT context (self-deletion only).
func (s *UserService) DeleteUser(requestId string, userId string) (response.DeleteUserResponse, *constant.AppError) {
	// Perform filesystem cleanup if dependencies are available
	if s.imageFinder != nil && s.storage != nil {
		images, err := s.imageFinder.FindImagesByUserId(userId)
		if err != nil {
			logger.Error(requestId, "failed to find user images for cleanup", err)
			// Continue with deletion even if image lookup fails
		} else {
			for _, img := range images {
				if err := s.storage.Delete(img.Id); err != nil {
					// Log failure and continue (requirement 3.4)
					logger.Warn(requestId, "failed to delete image file", "imageId", img.Id, "error", err.Error())
				}
			}
			logger.Info(requestId, "filesystem cleanup completed", "imageCount", len(images))
		}
	}

	// Delete user from database (cascade handles DB cleanup)
	if err := s.userRepo.Delete(userId); err != nil {
		logger.Error(requestId, "failed to delete user", err)
		return response.DeleteUserResponse{}, &constant.ErrDatabaseFailed
	}

	logger.Info(requestId, "user deleted", "userId", userId)

	return response.DeleteUserResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "User deleted successfully",
		},
	}, nil
}

// generateJWT creates a signed JWT token with userId and expiry claims.
func (s *UserService) generateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(s.jwtExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// validateEmail checks that the email contains exactly one "@" with non-empty
// local and domain parts.
func validateEmail(email string) *constant.AppError {
	if email == "" {
		return &constant.ErrValidationFailed
	}

	atIndex := strings.IndexByte(email, '@')
	if atIndex < 0 {
		return &constant.ErrValidationFailed
	}

	// Check for exactly one "@"
	if strings.Count(email, "@") != 1 {
		return &constant.ErrValidationFailed
	}

	local := email[:atIndex]
	domain := email[atIndex+1:]

	// Both local and domain parts must be non-empty
	if local == "" || domain == "" {
		return &constant.ErrValidationFailed
	}

	return nil
}

// isDuplicateKeyError checks if a database error is a PostgreSQL unique violation.
func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "unique constraint")
}
