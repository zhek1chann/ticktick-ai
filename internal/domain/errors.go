package domain

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Generic errors - used across all modules
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrForbidden     = errors.New("forbidden")
	ErrInvalidInput  = errors.New("invalid input")
	ErrNoChange      = errors.New("no change")
)

// PostgreSQL error codes
const (
	PgUniqueViolation     = "23505" // unique_violation
	PgForeignKeyViolation = "23503" // foreign_key_violation
	PgNotNullViolation    = "23502" // not_null_violation
)

// MapError converts database errors to domain errors.
// Use this in repositories to return consistent domain errors.
func MapError(err error) error {
	if err == nil {
		return nil
	}

	// Check for no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	// Check for postgres-specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PgUniqueViolation:
			return ErrAlreadyExists
		case PgForeignKeyViolation:
			return ErrInvalidInput // Referenced entity doesn't exist
		}
	}

	return err
}

// IsUniqueViolation checks if error is a unique constraint violation
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == PgUniqueViolation
	}
	return false
}

// IsForeignKeyViolation checks if error is a foreign key violation
func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == PgForeignKeyViolation
	}
	return false
}

// Auth errors
var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidToken        = errors.New("invalid token")
	ErrExpiredToken        = errors.New("token expired")
	ErrCannotSetSuperAdmin = errors.New("cannot assign super_admin role")
)

// Shop errors
var (
	ErrCannotOpenShop = errors.New("shop cannot be opened: settings and delivery house required")
	ErrNotShopOwner   = errors.New("user is not shop owner")
)

// Product errors
var (
	ErrCategoryHasProducts = errors.New("category has products and cannot be deleted")
)

// Product Request errors
var (
	ErrNotFilled   = errors.New("product request is not filled")
	ErrNoDataFound = errors.New("no data found from any source")
)

// Cart Errors
var (
	ErrOutOfStock = errors.New("out of stock")
)

// Order Errors
var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrEmptyOrder              = errors.New("order must have at least one item")
)
