package template

import (
	"icon_exchange/internal/order/model"
	"icon_exchange/internal/order/repository"
	"icon_exchange/internal/shared"
	"icon_exchange/internal/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repository
}

func NewOrderHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRouter(group *gin.RouterGroup) {
	group.GET("/me", h.getOrderByUserID)
}

// getOrderByUserID get all orders by user id
func (h *Handler) getOrderByUserID(c *gin.Context) {
	userID, ok := user.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, shared.ErrInternalServerError)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, shared.ErrInvalidOffset)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, shared.ErrInvalidLimit)
		return
	}

	orders, err := h.repo.GetByUserID(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, shared.Response[[]*model.Order]{
		Data: orders,
	})
}
