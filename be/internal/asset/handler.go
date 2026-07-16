package template

import service2 "icon_exchange/internal/template/service"

type Handler struct {
	service *service2.Service
}

func NewHandler(service *service2.Service) *Handler {
	return &Handler{service: service}
}
