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

type ImageFinder interface {
	FindImagesByUserId(userId string) ([]ImageRecord, error)
}

type ImageRecord struct {
	Id string
}

type StorageDeleter interface {
	Delete(imageId string) error
}

type UserService struct {
	userRepo  *repository.UserRepository
	imageFinder ImageFinder
	storage     StorageDeleter
	jwtSecret string
	jwtExpiry time.Duration
}

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

func (s *UserService) CreateUser(requestId string, req request.CreateUserRequest) (response.CreateUserResponse, *constant.AppError) {
	if err := validateEmail(req.Email); err != nil {
		logger.Warn(requestId, "email validation failed", "email", req.Email)
		return response.CreateUserResponse{}, err
	}

	if len(req.Email) > 255 {
		logger.Warn(requestId, "email exceeds 255 characters")
		return response.CreateUserResponse{}, &constant.ErrValidationFailed
	}

	if len(req.Password) < 8 || len(req.Password) > 72 {
		logger.Warn(requestId, "password length invalid")
		return response.CreateUserResponse{}, &constant.ErrValidationFailed
	}

	user, err := s.userRepo.Insert(req.Email, req.Password, req.Name)
	if err != nil {
		if isDuplicateKeyError(err) {
			logger.Warn(requestId, "duplicate email", "email", req.Email)
			return response.CreateUserResponse{}, &constant.ErrDuplicateResource
		}
		logger.Error(requestId, "failed to insert user", err)
		return response.CreateUserResponse{}, &constant.ErrDatabaseFailed
	}

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

func (s *UserService) Login(requestId string, req request.LoginRequest) (response.LoginResponse, *constant.AppError) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		logger.Warn(requestId, "login failed - user not found", "email", req.Email)
		return response.LoginResponse{}, &constant.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn(requestId, "login failed - password mismatch", "userId", user.Id)
		return response.LoginResponse{}, &constant.ErrUnauthorized
	}

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

func (s *UserService) DeleteUser(requestId string, userId string) (response.DeleteUserResponse, *constant.AppError) {
	if s.imageFinder != nil && s.storage != nil {
		images, err := s.imageFinder.FindImagesByUserId(userId)
		if err != nil {
			logger.Error(requestId, "failed to find user images for cleanup", err)
		} else {
			for _, img := range images {
				if err := s.storage.Delete(img.Id); err != nil {
					logger.Warn(requestId, "failed to delete image file", "imageId", img.Id, "error", err.Error())
				}
			}
			logger.Info(requestId, "filesystem cleanup completed", "imageCount", len(images))
		}
	}

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

func (s *UserService) generateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(s.jwtExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func validateEmail(email string) *constant.AppError {
	if email == "" {
		return &constant.ErrValidationFailed
	}

	atIndex := strings.IndexByte(email, '@')
	if atIndex < 0 {
		return &constant.ErrValidationFailed
	}

	if strings.Count(email, "@") != 1 {
		return &constant.ErrValidationFailed
	}

	local := email[:atIndex]
	domain := email[atIndex+1:]

	if local == "" || domain == "" {
		return &constant.ErrValidationFailed
	}

	return nil
}

func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "unique constraint")
}
