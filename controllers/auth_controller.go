package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hackaton-management-app/dto"
	"hackaton-management-app/services"
	"hackaton-management-app/utils"
)

// AuthController — pengganti route /api/auth/* milik Better Auth,
// versi tradisional: register + login berbasis JWT.
type AuthController struct {
	Service *services.AuthService
}

// Register — POST /auth/register (publik).
func (ctrl *AuthController) Register(c *gin.Context) {
	var input dto.RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	result, err := ctrl.Service.Register(input)
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "User registered successfully", result)
}

// Login — POST /auth/login (publik).
func (ctrl *AuthController) Login(c *gin.Context) {
	var input dto.LoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	result, err := ctrl.Service.Login(input)
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Login successful", result)
}
