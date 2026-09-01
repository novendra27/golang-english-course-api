package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse merepresentasikan format JSON response standar aplikasi
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse mengirimkan format response sukses dengan status code dan data
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirimkan format response error terstandarisasi
func ErrorResponse(c *gin.Context, statusCode int, message string, errDetails interface{}) {
	// Jika status code tidak dispesifikasikan, gunakan 500 Internal Server Error
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errDetails,
	})
}
