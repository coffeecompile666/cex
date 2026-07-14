package module_wallet

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the wallet API endpoints
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	wallets := router.Group("/wallets")
	{
		wallets.GET("/:userId", h.GetBalance)
	}
}

// GetBalance handles GET /wallets/:userId
func (h *Handler) GetBalance(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	wallet, err := h.service.GetBalance(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}
