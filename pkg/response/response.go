package response

import (
	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Operation completed successfully"`
	Data    T      `json:"data"`
}

type ErrorDetail struct {
	Err string `json:"error" example:"error description"`
}

type ErrorSwagger struct {
	Status  string      `json:"status" example:"error"`
	Message string      `json:"message" example:"Error message"`
	Data    ErrorDetail `json:"data"`
}

func ErrorResponse(c *gin.Context, statusCode int, msg string, err error) {
	errDetail := ErrorDetail{
		Err: err.Error(),
	}
	c.JSON(statusCode, Response[any]{
		Status:  "error",
		Message: msg,
		Data:    errDetail,
	})
}

func SuccessResponse(c *gin.Context, statusCode int, msg string, data interface{}) {
	c.JSON(statusCode, Response[any]{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}
