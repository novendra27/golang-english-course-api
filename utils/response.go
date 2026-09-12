// Package utils provides common utility helpers including standardized JSON responses and i18n support.
package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APIResponse represents the standard JSON API response structure.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse writes a standardized JSON success response to the Gin context.
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse writes a standardized JSON error response to the Gin context.
func ErrorResponse(c *gin.Context, statusCode int, message string, errDetails interface{}) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errDetails,
	})
}

// ValidationErrorResponse translates validator.ValidationErrors into field-level localized error messages.
func ValidationErrorResponse(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errMap := make(map[string]string)
		for _, fe := range ve {
			errMap[fe.Field()] = TranslateValidationError(c, fe)
		}
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Success: false,
			Message: Translate(c, "system.validation_failed", nil),
			Errors:  errMap,
		})
		return
	}

	// Fallback for non-validator errors (e.g. malformed JSON syntax)
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: Translate(c, "system.invalid_json", nil),
		Errors:  err.Error(),
	})
}


