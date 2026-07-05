package utils

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AppError — error HTTP terstruktur, padanan dari HttpException NestJS
// (NotFoundException, BadRequestException, dst.). Message bertipe `any`
// karena bisa berupa string biasa atau array field error hasil validasi.
type AppError struct {
	StatusCode int
	Message    any
}

func (e *AppError) Error() string {
	if s, ok := e.Message.(string); ok {
		return s
	}
	return http.StatusText(e.StatusCode)
}

func NewBadRequest(message any) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Message: message}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{StatusCode: http.StatusUnauthorized, Message: message}
}

func NewForbidden(message string) *AppError {
	return &AppError{StatusCode: http.StatusForbidden, Message: message}
}

func NewNotFound(message string) *AppError {
	return &AppError{StatusCode: http.StatusNotFound, Message: message}
}

// errorResponse — bentuk body error, meniru exception layer NestJS:
// { statusCode, message, error } dengan `error` berisi reason phrase HTTP.
type errorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    any    `json:"message"`
	Error      string `json:"error"`
}

// RespondError menerjemahkan error apa pun menjadi response error yang
// konsisten. AppError memakai status & message-nya sendiri; error tak
// terduga menjadi 500 tanpa membocorkan detail internal ke client.
func RespondError(c *gin.Context, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		log.Printf("unexpected error: %v", err)
		appErr = &AppError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	c.AbortWithStatusJSON(appErr.StatusCode, errorResponse{
		StatusCode: appErr.StatusCode,
		Message:    appErr.Message,
		Error:      http.StatusText(appErr.StatusCode),
	})
}

// FieldError — satu entri kesalahan validasi per field, meniru
// exceptionFactory di main.ts: { property, message }.
type FieldError struct {
	Property string `json:"property"`
	Message  string `json:"message"`
}

// RespondValidationError mengubah error hasil ShouldBindJSON menjadi 400
// dengan message berupa array { property, message } — satu entri per field
// yang tidak valid, sama seperti ValidationPipe di NestJS.
func RespondValidationError(c *gin.Context, err error) {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		// Body bukan JSON valid / tipe field salah (mis. tanggal bukan RFC3339).
		RespondError(c, NewBadRequest("Invalid request body: "+err.Error()))
		return
	}

	fields := make([]FieldError, 0, len(verrs))
	for _, fe := range verrs {
		fields = append(fields, FieldError{
			Property: fe.Field(),
			Message:  validationMessage(fe),
		})
	}
	RespondError(c, NewBadRequest(fields))
}

// validationMessage menyusun pesan per constraint, meniru gaya pesan
// class-validator.
func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " should not be empty"
	case "email":
		return fe.Field() + " must be an email"
	case "min":
		return fe.Field() + " must be longer than or equal to " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be shorter than or equal to " + fe.Param() + " characters"
	case "future":
		return fe.Field() + " must be a future date"
	default:
		return fe.Field() + " is invalid"
	}
}
