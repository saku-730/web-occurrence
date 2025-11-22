package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/service"
)

type WorkstationHandler struct {
	wsService service.WorkstationService
}

func NewWorkstationHandler(s service.WorkstationService) *WorkstationHandler {
	return &WorkstationHandler{wsService: s}
}

func (h *WorkstationHandler) Create(c *gin.Context) {
	var req model.CreateWorkstationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ws, err := h.wsService.CreateWorkstation(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ws)
}
