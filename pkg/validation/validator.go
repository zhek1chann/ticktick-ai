package validation

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"ticktick-ai/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidate binds JSON and validates the object.
// On error writes HTTP 400 and returns false.
func BindAndValidate(c *gin.Context, v *validator.Validate, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON", err)
		return false
	}

	if err := v.Struct(obj); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			var sb strings.Builder
			for _, e := range ve {
				sb.WriteString(fmt.Sprintf("%s failed on '%s'; ", e.Field(), e.Tag()))
			}
			response.ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors.New(strings.TrimSuffix(sb.String(), "; ")))
			return false
		}

		// if not belong to ValidationErrors
		response.ErrorResponse(c, http.StatusBadRequest, "Validation error", err)
		return false
	}

	return true
}
