package template

import (
	"icon_exchange/internal/asset/model"
	"icon_exchange/internal/asset/repository"
	"icon_exchange/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repository
}

func NewAssetHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// RegisterRoutes registers all template endpoints.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/user", h.GetByUserID)
}

// GetByUserID returns all markets.
func (h *Handler) GetByUserID(ctx *gin.Context) {
	// Import "icon_exchange/internal/user" nếu chưa được import, hoặc parse trực tiếp
	// Ở đây chúng ta sẽ lấy trực tiếp thông tin từ context đã được set bởi middleware
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id type"})
		return
	}

	asset, err := h.repo.GetByUserID(nil, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, shared.Response[[]*model.Asset]{
		Data: asset,
	})
}
