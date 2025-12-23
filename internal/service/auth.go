package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/whauzan/todo-api/internal/domain"
	"github.com/whauzan/todo-api/internal/pkg/apperror"
	"github.com/whauzan/todo-api/internal/pkg/jwt"
	"github.com/whauzan/todo-api/internal/pkg/password"
	"github.com/whauzan/todo-api/internal/repository"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo     repository.UserRepository
	tokenManager *jwt.TokenManager
	hasher       *password.Hasher
	logger       *slog.Logger
}

// NewAuthService creates a new AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	tokenManager *jwt.TokenManager,
	hasher *password.Hasher,
	logger *slog.Logger,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		hasher:       hasher,
		logger:       logger,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.UserInfo, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to check existing user", "error", err)
		return nil, apperror.ErrInternal
	}

	if existingUser != nil {
		return nil, apperror.ErrUserExists
	}

	// Hash password
	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to hash password", "error", err)
		return nil, apperror.ErrInternal
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "failed to create user", "error", err)
		return nil, apperror.ErrInternal
	}

	s.logger.InfoContext(ctx, "user registered successfully", "user_id", user.ID, "email", user.Email)

	return user.ToUserInfo(), nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get user by email", "error", err)
		return nil, apperror.ErrInternal
	}

	if user == nil {
		return nil, apperror.ErrInvalidCredentials
	}

	// Verify password
	if err := s.hasher.Verify(req.Password, user.PasswordHash); err != nil {
		if errors.Is(err, password.ErrMismatchedHashAndPassword) {
			return nil, apperror.ErrInvalidCredentials
		}
		s.logger.ErrorContext(ctx, "failed to verify password", "error", err)
		return nil, apperror.ErrInternal
	}

	// Generate JWT token
	tokenResp, err := s.tokenManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to generate token", "error", err)
		return nil, apperror.ErrInternal
	}

	s.logger.InfoContext(ctx, "user logged in successfully", "user_id", user.ID, "email", user.Email)

	return &domain.LoginResponse{
		Token:     tokenResp.Token,
		ExpiresAt: tokenResp.ExpiresAt,
		User:      user.ToUserInfo(),
	}, nil
}

// Refresh refreshes an existing JWT token
func (s *AuthService) Refresh(ctx context.Context, tokenString string) (*domain.LoginResponse, error) {
	// Refresh the token using the token manager
	tokenResp, err := s.tokenManager.RefreshToken(tokenString)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to refresh token", "error", err)
		return nil, apperror.NewAppError(
			apperror.CodeUnauthorized,
			"Invalid or expired token",
			401,
			err,
		)
	}

	// Validate the token to get user claims
	claims, err := s.tokenManager.ValidateToken(tokenResp.Token)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to validate refreshed token", "error", err)
		return nil, apperror.ErrInternal
	}

	// Get user info
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get user by ID", "error", err, "user_id", claims.UserID)
		return nil, apperror.ErrInternal
	}

	if user == nil {
		return nil, apperror.NewAppError(
			apperror.CodeNotFound,
			"User not found",
			404,
			fmt.Errorf("user with ID %s not found", claims.UserID),
		)
	}

	s.logger.InfoContext(ctx, "token refreshed successfully", "user_id", user.ID, "email", user.Email)

	return &domain.LoginResponse{
		Token:     tokenResp.Token,
		ExpiresAt: tokenResp.ExpiresAt,
		User:      user.ToUserInfo(),
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get user by ID", "error", err, "user_id", userID)
		return nil, apperror.ErrInternal
	}

	if user == nil {
		return nil, apperror.NewAppError(
			apperror.CodeNotFound,
			"User not found",
			404,
			fmt.Errorf("user with ID %s not found", userID),
		)
	}

	return user, nil
}
