package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/afidyoga/cinema-api/internal/model"
	"github.com/afidyoga/cinema-api/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	user, err := h.authSvc.Register(&req)
	if err != nil {
		c.JSON(http.StatusConflict, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{Success: true, Message: "registration successful", Data: user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	resp, err := h.authSvc.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Success: true, Message: "login successful", Data: resp})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	c.JSON(http.StatusOK, model.APIResponse{Success: true, Data: gin.H{"user_id": userID}})
}
