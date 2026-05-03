package middleware

import (
	"net/http"
	"strings"
	"ticktick-ai/internal/domain"
	"ticktick-ai/pkg/jwt"
	"ticktick-ai/pkg/response"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	jwtManager *jwt.Manager
}

func NewJWTMiddleware(jwtManager *jwt.Manager) *JWTMiddleware {
	return &JWTMiddleware{
		jwtManager: jwtManager,
	}
}

// RequireAuth validates JWT token and sets user info in context
func (m *JWTMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "missing authorization header", domain.ErrInvalidToken)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorResponse(c, http.StatusUnauthorized, "invalid authorization format", domain.ErrInvalidToken)
			c.Abort()
			return
		}

		claims, err := m.jwtManager.ValidateToken(parts[1])
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token", err)
			c.Abort()
			return
		}

		c.Set(domain.CtxKeyUserID, int64(claims.UserID))
		c.Set(domain.CtxKeyRole, claims.Role)

		c.Next()
	}
}

// OptionalAuth validates JWT token if present, but doesn't require it
func (m *JWTMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := m.jwtManager.ValidateToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set(domain.CtxKeyUserID, int64(claims.UserID))
		c.Set(domain.CtxKeyRole, claims.Role)

		c.Next()
	}
}

// RequireRole validates that the user has one of the required roles
func (m *JWTMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(domain.CtxKeyRole)
		if !exists {
			response.ErrorResponse(c, http.StatusForbidden, "access denied", domain.ErrForbidden)
			c.Abort()
			return
		}

		userRole := roleValue.(string)

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.ErrorResponse(c, http.StatusForbidden, "insufficient permissions", domain.ErrForbidden)
		c.Abort()
	}
}
