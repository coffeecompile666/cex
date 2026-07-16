package repository

import (
	"icon_exchange/internal/user/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *model.User, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(user).Error
}

func (r *UserRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("mail = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) MarkVerified(userID uint, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(&model.User{}).Where("id = ?", userID).Update("is_verified", true).Error
}

func (r *UserRepo) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
