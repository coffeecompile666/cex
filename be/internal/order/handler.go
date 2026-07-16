package template

import (
	"icon_exchange/internal/order/service"
)

type Handler struct {
	service *service.Service
}

func NewOrderHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
