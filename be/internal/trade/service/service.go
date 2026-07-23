package service

import (
	"icon_exchange/internal/template/repository"
)

type ITradeService interface {
}

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}
