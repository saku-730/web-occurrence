package handler

import (
	"github.com/saku-730/web-occurrence/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	masterService service.MasterService
}

func NewMasterHandler(s service.MasterService) *MasterHandler {
	return &MasterHandler{masterService: s}
}

// GetMasterData はマスターデータをJSONで返すのだ
func (h *MasterHandler) GetMasterData(c *gin.Context) {
	data, err := h.masterService.GetMasterData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "マスターデータの取得に失敗: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
