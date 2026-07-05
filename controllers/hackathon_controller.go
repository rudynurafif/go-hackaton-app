package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hackaton-management-app/dto"
	"hackaton-management-app/middleware"
	"hackaton-management-app/services"
	"hackaton-management-app/utils"
)

// HackathonController — padanan hackaton.controller.ts.
// Read terbuka untuk publik; write khusus ADMIN; join khusus PARTICIPANT
// (diatur di routes).
type HackathonController struct {
	Service *services.HackathonService
}

// Create — POST /hackaton (khusus ADMIN).
func (ctrl *HackathonController) Create(c *gin.Context) {
	var input dto.CreateHackathonDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	// Author adalah admin yang login, diambil dari token — bukan dari client.
	user := middleware.CurrentUser(c)
	hackathon, err := ctrl.Service.Create(input, user.ID)
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "Hackathon created successfully", hackathon)
}

// FindAll — GET /hackaton (publik).
func (ctrl *HackathonController) FindAll(c *gin.Context) {
	hackathons, err := ctrl.Service.FindAll()
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Success", hackathons)
}

// FindOne — GET /hackaton/:id (publik).
func (ctrl *HackathonController) FindOne(c *gin.Context) {
	hackathon, err := ctrl.Service.FindOne(c.Param("id"))
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Success", hackathon)
}

// Update — PATCH /hackaton/:id (khusus ADMIN).
func (ctrl *HackathonController) Update(c *gin.Context) {
	var input dto.UpdateHackathonDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	hackathon, err := ctrl.Service.Update(c.Param("id"), input)
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Hackathon updated successfully", hackathon)
}

// Remove — DELETE /hackaton/:id (khusus ADMIN).
func (ctrl *HackathonController) Remove(c *gin.Context) {
	hackathon, err := ctrl.Service.Remove(c.Param("id"))
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Hackathon deleted successfully", hackathon)
}

// Join — POST /hackaton/:id/join (khusus PARTICIPANT).
func (ctrl *HackathonController) Join(c *gin.Context) {
	user := middleware.CurrentUser(c)
	participant, err := ctrl.Service.Join(c.Param("id"), user.ID)
	if err != nil {
		utils.RespondError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "Success joined hackathon", participant)
}
