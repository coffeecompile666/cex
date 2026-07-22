package template

import (
	"icon_exchange/internal/market/model"
	"icon_exchange/internal/market/repository"
	"icon_exchange/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository *repository.Repository
}

func NewMarketHandler(repository *repository.Repository) *Handler {
	return &Handler{repository: repository}
}

// RegisterRoutes registers all template endpoints.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", h.GetAllMarkets)
	router.GET("/base")
}

// GetAllMarkets returns all markets.
func (h *Handler) GetAllMarkets(ctx *gin.Context) {
	markets, err := h.repository.GetMarkets()
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, shared.Response[[]model.Market]{
		Data: markets,
	})
}

// GetBaseCurrency returns usd currency config
func (h *Handler) GetBaseCurrency(ctx *gin.Context) {
	base, err := h.repository.GetBaseCurrency()
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, shared.Response[*model.Market]{
		Data: base,
	})
}
