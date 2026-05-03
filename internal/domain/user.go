package domain

import "time"

type User struct {
	ID           int
	Name         string
	PhoneNumber  string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthResult struct {
	Tokens TokenPair
	User   User
}

type TokenClaims struct {
	UserID int
	Role   Role
	Exp    time.Time
}
