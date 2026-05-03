package row

import (
	"ticktick-ai/internal/domain"
	"time"
)

// User represents user data in database
type User struct {
	ID           int
	Name         string
	PhoneNumber  string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ToUser converts database row to domain user
func ToUser(row User) domain.User {
	return domain.User{
		ID:           row.ID,
		Name:         row.Name,
		PhoneNumber:  row.PhoneNumber,
		PasswordHash: row.PasswordHash,
		Role:         domain.Role(row.Role),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

// FromUser converts domain user to database row
func FromUser(user domain.User) User {
	return User{
		ID:           user.ID,
		Name:         user.Name,
		PhoneNumber:  user.PhoneNumber,
		PasswordHash: user.PasswordHash,
		Role:         string(user.Role),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
