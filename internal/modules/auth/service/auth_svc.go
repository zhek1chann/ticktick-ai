package service

import (
	"context"
	"log/slog"
	"ticktick-ai/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type userRepo interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	UserByPhoneNumber(ctx context.Context, phoneNumber string) (domain.User, error)
	UserByID(ctx context.Context, id int) (domain.User, error)
	UpdateUserRole(ctx context.Context, userID int, role string) error
	ListUsers(ctx context.Context) ([]domain.User, error)
}

func (s *Service) Register(ctx context.Context, name, phoneNumber, password string) (domain.AuthResult, error) {
	var result domain.AuthResult

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Check if user already exists
		existingUser, errTx := s.repo.UserByPhoneNumber(ctx, phoneNumber)
		if errTx == nil && existingUser.ID != 0 {
			return domain.ErrAlreadyExists
		}

		// Hash password
		hashedPassword, errTx := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if errTx != nil {
			return errTx
		}

		// Create user
		user := domain.User{
			Name:         name,
			PhoneNumber:  phoneNumber,
			PasswordHash: string(hashedPassword),
			Role:         domain.RoleUser,
		}

		userID, errTx := s.repo.CreateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		// GetCart created user
		createdUser, errTx := s.repo.UserByID(ctx, userID)
		if errTx != nil {
			return errTx
		}

		// Generate tokens
		accessToken, errTx := s.jwtManager.GenerateAccessToken(userID, domain.RoleUser.String())
		if errTx != nil {
			return errTx
		}

		refreshToken, errTx := s.jwtManager.GenerateRefreshToken(userID, domain.RoleUser.String())
		if errTx != nil {
			return errTx
		}

		result = domain.AuthResult{
			Tokens: domain.TokenPair{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
			User: createdUser,
		}

		return nil
	})

	if err != nil {
		slog.ErrorContext(ctx, "Register", "error", err)
		return domain.AuthResult{}, err
	}

	return result, nil
}

func (s *Service) Login(ctx context.Context, phoneNumber, password string) (domain.AuthResult, error) {
	var result domain.AuthResult

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// GetCart user by phone number
		user, errTx := s.repo.UserByPhoneNumber(ctx, phoneNumber)
		if errTx != nil {
			return domain.ErrInvalidCredentials
		}

		// Verify password
		errTx = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if errTx != nil {
			return domain.ErrInvalidCredentials
		}

		// Generate tokens
		accessToken, errTx := s.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
		if errTx != nil {
			return errTx
		}

		refreshToken, errTx := s.jwtManager.GenerateRefreshToken(user.ID, string(user.Role))
		if errTx != nil {
			return errTx
		}

		result = domain.AuthResult{
			Tokens: domain.TokenPair{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
			User: user,
		}

		return nil
	})

	if err != nil {
		slog.ErrorContext(ctx, "Login", "error", err)
		return domain.AuthResult{}, err
	}

	return result, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	var tokens domain.TokenPair

	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		slog.ErrorContext(ctx, "RefreshToken", "error", err)
		return domain.TokenPair{}, err
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Verify user still exists
		user, errTx := s.repo.UserByID(ctx, claims.UserID)
		if errTx != nil {
			return domain.ErrNotFound
		}

		// Generate new tokens
		accessToken, errTx := s.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
		if errTx != nil {
			return errTx
		}

		newRefreshToken, errTx := s.jwtManager.GenerateRefreshToken(user.ID, string(user.Role))
		if errTx != nil {
			return errTx
		}

		tokens = domain.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
		}

		return nil
	})

	if err != nil {
		slog.ErrorContext(ctx, "RefreshToken", "error", err)
		return domain.TokenPair{}, err
	}

	return tokens, nil
}

func (s *Service) UserByID(ctx context.Context, userID int) (domain.User, error) {
	var user domain.User

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		user, errTx = s.repo.UserByID(ctx, userID)
		if errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		slog.ErrorContext(ctx, "UserByID", "error", err)
		return domain.User{}, err
	}

	return user, nil
}

// ValidateUserIsShopOwner checks if user exists and has shop_owner role
func (s *Service) ValidateUserIsShopOwner(ctx context.Context, userID int) (bool, error) {
	var isShopOwner bool

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		// GetCart user
		user, errTx := s.repo.UserByID(ctx, userID)
		if errTx != nil {
			return errTx
		}

		// Check if user has shop_owner role
		isShopOwner = user.Role == domain.RoleShopOwner

		return nil
	})

	if err != nil {
		slog.ErrorContext(ctx, "ValidateUserIsShopOwner", "error", err, "user_id", userID)
		return false, err
	}

	return isShopOwner, nil
}

// UpdateUserRole updates a user's role (admin only)
// Cannot assign super_admin role - returns ErrCannotSetSuperAdmin
func (s *Service) UpdateUserRole(ctx context.Context, userID int, newRole domain.Role) error {
	// Prevent assigning super_admin role
	if newRole == domain.RoleSuperAdmin {
		return domain.ErrCannotSetSuperAdmin
	}

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Verify user exists
		_, errTx := s.repo.UserByID(ctx, userID)
		if errTx != nil {
			return errTx
		}

		// Update role
		errTx = s.repo.UpdateUserRole(ctx, userID, newRole.String())
		return errTx
	})

	if err != nil {
		slog.ErrorContext(ctx, "UpdateUserRole", "error", err, "user_id", userID, "new_role", newRole)
		return err
	}

	slog.InfoContext(ctx, "User role updated", "user_id", userID, "new_role", newRole)
	return nil
}

// ListUsers returns all users (admin only)
func (s *Service) ListUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		users, errTx = s.repo.ListUsers(ctx)
		return errTx
	})

	if err != nil {
		slog.ErrorContext(ctx, "ListUsers", "error", err)
		return nil, err
	}

	return users, nil
}
