package repository

import (
	"time"

	"icon_exchange/internal/user/model"

	"gorm.io/gorm"
)

type OTPRepo struct {
	db *gorm.DB
}

func NewOTPRepo(db *gorm.DB) *OTPRepo {
	return &OTPRepo{db: db}
}

func (r *OTPRepo) Create(otp *model.OTP, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(otp).Error
}

// RevokeAllByUserID sets RevokedAt = now() for all active OTPs of a user.
func (r *OTPRepo) RevokeAllByUserID(userID uint, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	now := time.Now()
	return db.Model(&model.OTP{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

// GetLatestActiveByUserID returns the most recent non-revoked, non-expired OTP for a user.
func (r *OTPRepo) GetLatestActiveByUserID(userID uint) (*model.OTP, error) {
	var otp model.OTP
	err := r.db.
		Where("user_id = ? AND revoked_at IS NULL AND expired_at > ?", userID, time.Now()).
		Order("created_at DESC").
		First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// IncrementFailAttempts increments the fail counter for a given OTP.
func (r *OTPRepo) IncrementFailAttempts(otpID uint) error {
	return r.db.Model(&model.OTP{}).
		Where("id = ?", otpID).
		UpdateColumn("fail_attempts", gorm.Expr("fail_attempts + 1")).Error
}
