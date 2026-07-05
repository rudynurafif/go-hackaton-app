package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hackaton-management-app/services"
	"hackaton-management-app/utils"
)

// UserController — padanan user.controller.ts. Pembatasan akses
// (login + role) dipasang di routes, bukan di sini.
type UserController struct {
	Service *services.UserService
}

// FindAll — GET /user/all (khusus ADMIN).
func (ctrl *UserController) FindAll(c *gin.Context) {
	users, err := ctrl.Service.FindAll()
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "All users retrieved", users)
}

// FindOne — GET /user/:id (semua user yang login).
func (ctrl *UserController) FindOne(c *gin.Context) {
	user, err := ctrl.Service.FindOne(c.Param("id"))
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Success", user)
}
