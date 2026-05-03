package auth

import (
	"ticktick-ai/internal/middleware"
	"ticktick-ai/internal/modules/auth/api"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *api.Handlerauth, jwtMiddleware *middleware.JWTMiddleware) {
	authRoutes := router.Group("/auth")
	{
		// Public routes (no authentication required)
		authRoutes.POST("/register", h.Register)
		authRoutes.POST("/login", h.Login)
		authRoutes.POST("/refresh", h.RefreshToken)

		// Protected routes (authentication required)
		authRoutes.GET("/me", jwtMiddleware.RequireAuth(), h.Me)

		// Admin routes (super_admin only)
		authRoutes.GET("/users", jwtMiddleware.RequireAuth(), h.ListUsers)
		authRoutes.PUT("/users/:user_id/role", jwtMiddleware.RequireAuth(), h.UpdateUserRole)
	}
}
