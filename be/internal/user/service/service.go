package service

import (
	"icon_exchange/internal/user/model"
	"icon_exchange/internal/user/repository"
)

type IUserService interface {
	GetUserByID(id uint) (*model.User, error)
}

type UserService struct {
	repo repository.UserRepo
}

func (u UserService) GetUserByID(id uint) (*model.User, error) {
	return u.repo.GetByID(id)
}

func NewUserService(repo repository.UserRepo) IUserService {
	return &UserService{repo: repo}
}
