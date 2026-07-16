package template

import (
	"icon_exchange/internal/market/model"
	"icon_exchange/internal/market/repository"
	"icon_exchange/internal/market/service"
	"icon_exchange/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service    *service.Service
	repository *repository.Repository
}

func NewMarketHandler(service *service.Service, repository *repository.Repository) *Handler {
	return &Handler{service: service, repository: repository}
}

// RegisterRoutes registers all template endpoints.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/markets", h.GetAllMarkets)
}

// GetAllMarkets returns all markets.
func (h *Handler) GetAllMarkets(ctx *gin.Context) {
	markets, err := h.repository.GetAll()
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, shared.Response[[]model.Market]{
		Data: markets,
	})
}
