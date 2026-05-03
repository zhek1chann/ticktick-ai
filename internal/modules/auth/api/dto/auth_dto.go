package dto

import (
	"ticktick-ai/internal/domain"
	"ticktick-ai/pkg/response"
)

// ===== Swagger type shit ======

type RegisterSwaggerResponse struct {
	response.Response[RegisterResponse]
}

type LoginSwaggerResponse struct {
	response.Response[LoginResponse]
}

type RefreshTokenSwaggerResponse struct {
	response.Response[RefreshTokenResponse]
}

type GetMeSwaggerResponse struct {
	response.Response[GetMeResponse]
}

type UpdateUserRoleSwaggerResponse struct {
	response.Response[any]
}

type ListUsersSwaggerResponse struct {
	response.Response[ListUsersResponse]
}

// ===== Normal req/resp =====

type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GetMeResponse struct {
	User User `json:"user"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user shop_owner"`
}

type ListUsersResponse struct {
	Users []User `json:"users"`
}

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

// FromUser converts domain user to DTO user
func FromUser(user domain.User) User {
	return User{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Role:        string(user.Role),
	}
}

// FromUsers converts domain users to DTO users
func FromUsers(users []domain.User) []User {
	result := make([]User, len(users))
	for i, user := range users {
		result[i] = FromUser(user)
	}
	return result
}

// FromTokenPair converts domain token pair to DTO
func FromTokenPair(tokens domain.TokenPair) (string, string) {
	return tokens.AccessToken, tokens.RefreshToken
}
