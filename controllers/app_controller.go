// Package controllers berisi HTTP handler tipis: bind + validasi request,
// panggil service, kirim response lewat envelope standar.
// Padanan dari file *.controller.ts di NestJS.
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hackaton-management-app/middleware"
	"hackaton-management-app/utils"
)

// AppController — padanan app.controller.ts.
type AppController struct{}

// Hello — GET / (publik).
func (ctrl *AppController) Hello(c *gin.Context) {
	utils.Success(c, http.StatusOK, "Success", "Hello World!")
}

// Me — GET /me (butuh login). Mengembalikan user yang sedang login,
// padanan { user: session.user }.
func (ctrl *AppController) Me(c *gin.Context) {
	user := middleware.CurrentUser(c)
	utils.Success(c, http.StatusOK, "Success", gin.H{"user": user})
}
