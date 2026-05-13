package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/afidyoga/cinema-api/internal/model"
	"github.com/afidyoga/cinema-api/internal/service"
)

type ScheduleHandler struct {
	scheduleSvc *service.ScheduleService
}

func NewScheduleHandler(scheduleSvc *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleSvc: scheduleSvc}
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	var req model.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	schedule, err := h.scheduleSvc.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{Success: true, Message: "schedule created", Data: schedule})
}

func (h *ScheduleHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	schedules, total, err := h.scheduleSvc.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: gin.H{
			"schedules": schedules,
			"total":     total,
			"page":      page,
			"limit":     limit,
		},
	})
}

func (h *ScheduleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	schedule, err := h.scheduleSvc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Success: true, Data: schedule})
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	schedule, err := h.scheduleSvc.Update(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "schedule not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Success: true, Message: "schedule updated", Data: schedule})
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.scheduleSvc.Delete(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "schedule not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Success: true, Message: "schedule deleted"})
}
