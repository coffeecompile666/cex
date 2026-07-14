package repository

import (
	"time"

	"icon_exchange/internal/module_user/model"

	"gorm.io/gorm"
)

type SessionRepo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(session *model.Session, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(session).Error
}

// GetByHashToken finds an active (non-revoked, non-expired) session by its hashed refresh token.
func (r *SessionRepo) GetByHashToken(hash string) (*model.Session, error) {
	var session model.Session
	err := r.db.
		Where("hash_refresh_token = ? AND revoked_at IS NULL AND expired_at > ?", hash, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// RevokeByID marks a session as revoked.
func (r *SessionRepo) RevokeByID(sessionID uint) error {
	now := time.Now()
	return r.db.Model(&model.Session{}).
		Where("id = ?", sessionID).
		Update("revoked_at", now).Error
}
