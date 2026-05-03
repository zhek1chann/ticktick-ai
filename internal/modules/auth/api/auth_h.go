package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"ticktick-ai/internal/domain"
	"ticktick-ai/internal/modules/auth/api/dto"
	"ticktick-ai/pkg/response"
	"ticktick-ai/pkg/validation"

	"github.com/gin-gonic/gin"
)

type authService interface {
	Register(ctx context.Context, name, phoneNumber, password string) (domain.AuthResult, error)
	Login(ctx context.Context, phoneNumber, password string) (domain.AuthResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (domain.TokenPair, error)
	UserByID(ctx context.Context, userID int) (domain.User, error)
	UpdateUserRole(ctx context.Context, userID int, newRole domain.Role) error
	ListUsers(ctx context.Context) ([]domain.User, error)
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account and returns JWT tokens
// @Tags         user.auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Register request"
// @Success      200 {object} dto.RegisterSwaggerResponse "Successful registration with tokens"
// @Failure      400 {object} response.ErrorSwagger "Invalid request"
// @Failure      409 {object} response.ErrorSwagger "User already exists"
// @Failure      500 {object} response.ErrorSwagger "Internal error"
// @Router       /api/auth/register [post]
func (h *Handlerauth) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !validation.BindAndValidate(c, h.validator, &req) {
		return
	}

	authResult, err := h.svcAuth.Register(c.Request.Context(), req.Name, req.PhoneNumber, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			response.ErrorResponse(c, http.StatusConflict, "user already exists", err)
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, "internal error", err)
		return
	}

	accessToken, refreshToken := dto.FromTokenPair(authResult.Tokens)
	response.SuccessResponse(c, http.StatusOK, "registration successful", dto.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.FromUser(authResult.User),
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates user and returns JWT tokens
// @Tags         user.auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login request"
// @Success      200 {object} dto.LoginSwaggerResponse "Successful login with tokens"
// @Failure      400 {object} response.ErrorSwagger "Invalid request"
// @Failure      401 {object} response.ErrorSwagger "Invalid credentials"
// @Failure      500 {object} response.ErrorSwagger "Internal error"
// @Router       /api/auth/login [post]
func (h *Handlerauth) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !validation.BindAndValidate(c, h.validator, &req) {
		return
	}

	authResult, err := h.svcAuth.Login(c.Request.Context(), req.PhoneNumber, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.ErrorResponse(c, http.StatusUnauthorized, "invalid credentials", err)
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, "internal error", err)
		return
	}

	accessToken, refreshToken := dto.FromTokenPair(authResult.Tokens)
	response.SuccessResponse(c, http.StatusOK, "login successful", dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.FromUser(authResult.User),
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Generates new access and refresh tokens using a valid refresh token
// @Tags         user.auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh token request"
// @Success      200 {object} dto.RefreshTokenSwaggerResponse "Successful token refresh"
// @Failure      400 {object} response.ErrorSwagger "Invalid request"
// @Failure      401 {object} response.ErrorSwagger "Invalid or expired token"
// @Router       /api/auth/refresh [post]
func (h *Handlerauth) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if !validation.BindAndValidate(c, h.validator, &req) {
		return
	}

	tokens, err := h.svcAuth.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token", err)
		return
	}

	accessToken, refreshToken := dto.FromTokenPair(tokens)
	response.SuccessResponse(c, http.StatusOK, "token refresh successful", dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Me godoc
// @Summary      Get current user
// @Description  Returns current user information from JWT token
// @Tags         user.auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.GetMeSwaggerResponse "Current user information"
// @Failure      401 {object} response.ErrorSwagger "Unauthorized"
// @Failure      500 {object} response.ErrorSwagger "Internal error"
// @Router       /api/auth/me [get]
func (h *Handlerauth) Me(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := domain.GetUserIDFromContext(c)
	if userID == 0 {
		response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized", domain.ErrInvalidToken)
		return
	}

	user, err := h.svcAuth.UserByID(c.Request.Context(), int(userID))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "internal error", err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "user information retrieved", dto.GetMeResponse{
		User: dto.FromUser(user),
	})
}

// UpdateUserRole godoc
// @Summary      Update user role (super_admin only)
// @Description  Updates a user's role. Only accessible by super_admin. Cannot assign super_admin role.
// @Tags         admin.auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id path int true "User ID"
// @Param        request body dto.UpdateUserRoleRequest true "New role (user or shop_owner only)"
// @Success      200 {object} dto.UpdateUserRoleSwaggerResponse "Role updated successfully"
// @Failure      400 {object} response.ErrorSwagger "Invalid request or cannot set super_admin"
// @Failure      401 {object} response.ErrorSwagger "Unauthorized"
// @Failure      403 {object} response.ErrorSwagger "Forbidden - not super_admin"
// @Failure      404 {object} response.ErrorSwagger "User not found"
// @Failure      500 {object} response.ErrorSwagger "Internal error"
// @Router       /api/auth/users/{user_id}/role [put]
func (h *Handlerauth) UpdateUserRole(c *gin.Context) {
	// Get authenticated user's role from context
	userRole, exists := c.Get(domain.CtxKeyRole)
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", domain.ErrInvalidToken)
		return
	}

	// Only super_admin can update roles
	if userRole.(string) != domain.RoleSuperAdmin.String() {
		response.ErrorResponse(c, http.StatusForbidden, "Only super_admin can update user roles", domain.ErrForbidden)
		return
	}

	// Get target user ID from path
	userIDParam := c.Param("user_id")
	userID := 0
	if _, err := fmt.Sscanf(userIDParam, "%d", &userID); err != nil || userID == 0 {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Parse request body
	var req dto.UpdateUserRoleRequest
	if !validation.BindAndValidate(c, h.validator, &req) {
		return
	}

	// Convert string role to domain.Role
	var newRole domain.Role
	switch req.Role {
	case "user":
		newRole = domain.RoleUser
	case "shop_owner":
		newRole = domain.RoleShopOwner
	case "super_admin":
		newRole = domain.RoleSuperAdmin
	default:
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid role", domain.ErrInvalidInput)
		return
	}

	// Update role
	err := h.svcAuth.UpdateUserRole(c.Request.Context(), userID, newRole)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}

		if errors.Is(err, domain.ErrCannotSetSuperAdmin) {
			response.ErrorResponse(c, http.StatusBadRequest, "Cannot assign super_admin role", err)
			return
		}

		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user role", err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "User role updated successfully", nil)
}

// ListUsers godoc
// @Summary      List all users (super_admin only)
// @Description  Returns a list of all users in the system. Only accessible by super_admin.
// @Tags         admin.auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.ListUsersSwaggerResponse "List of all users"
// @Failure      401 {object} response.ErrorSwagger "Unauthorized"
// @Failure      403 {object} response.ErrorSwagger "Forbidden - not super_admin"
// @Failure      500 {object} response.ErrorSwagger "Internal error"
// @Router       /api/auth/users [get]
func (h *Handlerauth) ListUsers(c *gin.Context) {
	// Get authenticated user's role from context
	userRole, exists := c.Get(domain.CtxKeyRole)
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", domain.ErrInvalidToken)
		return
	}

	// Only super_admin can list users
	if userRole.(string) != domain.RoleSuperAdmin.String() {
		response.ErrorResponse(c, http.StatusForbidden, "Only super_admin can list users", domain.ErrForbidden)
		return
	}

	users, err := h.svcAuth.ListUsers(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to list users", err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", dto.ListUsersResponse{
		Users: dto.FromUsers(users),
	})
}
