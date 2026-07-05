// Package utils berisi helper lintas-layer: format response, error
// terstruktur, dan JWT.
package utils

import "github.com/gin-gonic/gin"

// APIResponse — envelope sukses yang konsisten untuk semua endpoint,
// padanan dari ResponseInterceptor di NestJS:
// { statusCode, message, data }.
type APIResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

// Success mengirim response sukses dengan envelope standar. `message`
// dipasok per-handler, padanan dari decorator @ResponseMessage; gunakan
// "Success" untuk default.
func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}
