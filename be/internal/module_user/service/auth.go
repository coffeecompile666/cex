package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"icon_exchange/internal/module_user/model"
	"icon_exchange/internal/module_user/repository"
	"icon_exchange/internal/shared"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ─── TokenPair ─────────────────────────────────────────────────────────────────

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"` // raw token returned to client
}

// ─── MailerInterface ──────────────────────────────────────────────────────────

type MailerInterface interface {
	SendOTPEmail(toEmail string, otpCode string) error
}

// ─── AuthService ───────────────────────────────────────────────────────────────

type AuthService struct {
	userRepo    *repository.UserRepo
	otpRepo     *repository.OTPRepo
	sessionRepo *repository.SessionRepo
	mailer      MailerInterface
	db          *gorm.DB
}

func NewAuthService(
	db *gorm.DB,
	userRepo *repository.UserRepo,
	otpRepo *repository.OTPRepo,
	sessionRepo *repository.SessionRepo,
	mailer MailerInterface,
) *AuthService {
	return &AuthService{
		db:          db,
		userRepo:    userRepo,
		otpRepo:     otpRepo,
		sessionRepo: sessionRepo,
		mailer:      mailer,
	}
}

// ─── Signup ────────────────────────────────────────────────────────────────────

// Signup registers a new user (unverified) and sends an OTP email.
// If the email already exists but is unverified, it revokes the old OTP and resends a new one.
func (s *AuthService) Signup(email, password string) error {
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return shared.ErrInternalServerError
	}

	var userID uint

	if existingUser != nil {
		if existingUser.IsVerified {
			return shared.ErrEmailAlreadyRegistered
		}
		// Unverified: allow resending OTP
		userID = existingUser.ID
	} else {
		hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return shared.ErrInternalServerError
		}

		user := &model.User{
			Mail:         email,
			HashPassword: string(hashPwd),
			IsVerified:   false,
		}
		if err := s.userRepo.Create(user, nil); err != nil {
			return shared.ErrInternalServerError
		}
		userID = user.ID
	}

	return s.createAndSendOTP(userID, email)
}

// ─── VerifyOTP ─────────────────────────────────────────────────────────────────

// VerifyOTP validates the OTP for signup and returns a token pair on success.
func (s *AuthService) VerifyOTP(email, code string) (*TokenPair, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, shared.ErrInternalServerError
	}

	if user.IsVerified {
		return nil, shared.ErrConflict
	}

	otp, err := s.otpRepo.GetLatestActiveByUserID(user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrOTPNotFound
		}
		return nil, shared.ErrInternalServerError
	}

	if otp.FailAttempts >= s.getOTPMaxAttempts() {
		return nil, shared.ErrOTPMaxAttempts
	}

	if err := bcrypt.CompareHashAndPassword([]byte(otp.HashCode), []byte(code)); err != nil {
		_ = s.otpRepo.IncrementFailAttempts(otp.ID)
		return nil, shared.ErrOTPInvalid
	}

	var tokenPair *TokenPair
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepo.MarkVerified(user.ID, tx); err != nil {
			return err
		}
		if err := s.otpRepo.RevokeAllByUserID(user.ID, tx); err != nil {
			return err
		}
		pair, err := s.createSession(user.ID, tx)
		if err != nil {
			return err
		}
		tokenPair = pair
		return nil
	})
	if txErr != nil {
		return nil, shared.ErrInternalServerError
	}

	return tokenPair, nil
}

// ─── ResendOTP ─────────────────────────────────────────────────────────────────

// ResendOTP revokes all existing OTPs and sends a fresh one.
func (s *AuthService) ResendOTP(email string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared.ErrNotFound
		}
		return shared.ErrInternalServerError
	}
	if user.IsVerified {
		return shared.ErrConflict
	}
	return s.createAndSendOTP(user.ID, email)
}

// ─── Signin ────────────────────────────────────────────────────────────────────

// Signin authenticates an existing verified user and returns a token pair.
func (s *AuthService) Signin(email, password string) (*TokenPair, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// Always return ErrInvalidCredentials to avoid user enumeration
		return nil, shared.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return nil, shared.ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, shared.ErrUserNotVerified
	}

	return s.createSession(user.ID, nil)
}

// ─── RefreshToken ──────────────────────────────────────────────────────────────

// RefreshToken validates a raw refresh token, rotates it, and issues a new token pair.
func (s *AuthService) RefreshToken(rawToken string) (*TokenPair, error) {
	hash := s.hashRefreshToken(rawToken)

	session, err := s.sessionRepo.GetByHashToken(hash)
	if err != nil {
		return nil, shared.ErrInvalidToken
	}

	var tokenPair *TokenPair
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.sessionRepo.RevokeByID(session.ID); err != nil {
			return err
		}
		pair, err := s.createSession(session.UserID, tx)
		if err != nil {
			return err
		}
		tokenPair = pair
		return nil
	})
	if txErr != nil {
		return nil, shared.ErrInternalServerError
	}

	return tokenPair, nil
}

// ─── Helpers ───────────────────────────────────────────────────────────────────

func (s *AuthService) createAndSendOTP(userID uint, email string) error {
	code, err := generateOTPCode()
	if err != nil {
		return shared.ErrInternalServerError
	}

	hashCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return shared.ErrInternalServerError
	}

	expiryMinutes := s.getOTPExpiryMinutes()

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.otpRepo.RevokeAllByUserID(userID, tx); err != nil {
			return err
		}
		otp := &model.OTP{
			HashCode:     string(hashCode),
			FailAttempts: 0,
			ExpiredAt:    time.Now().Add(time.Duration(expiryMinutes) * time.Minute),
			UserID:       userID,
		}
		if err := s.otpRepo.Create(otp, tx); err != nil {
			return err
		}
		return s.mailer.SendOTPEmail(email, code)
	})
	if txErr != nil {
		return fmt.Errorf("%w: %v", shared.ErrInternalServerError, txErr)
	}
	return nil
}

func (s *AuthService) createSession(userID uint, tx *gorm.DB) (*TokenPair, error) {
	rawRefreshToken, err := generateSecureToken()
	if err != nil {
		return nil, shared.ErrInternalServerError
	}

	session := &model.Session{
		HashRefreshToken: s.hashRefreshToken(rawRefreshToken),
		ExpiredAt:        time.Now().Add(30 * 24 * time.Hour),
		UserID:           userID,
	}
	if err := s.sessionRepo.Create(session, tx); err != nil {
		return nil, shared.ErrInternalServerError
	}

	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, shared.ErrInternalServerError
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
	}, nil
}

func (s *AuthService) generateAccessToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_MINUTES"))
	if expiry == 0 {
		expiry = 15
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(expiry) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func (s *AuthService) hashRefreshToken(raw string) string {
	masterKey := os.Getenv("REFRESH_TOKEN_MASTER_KEY")
	mac := hmac.New(sha256.New, []byte(masterKey))
	mac.Write([]byte(raw))
	return hex.EncodeToString(mac.Sum(nil))
}

func generateOTPCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *AuthService) getOTPExpiryMinutes() int {
	v, _ := strconv.Atoi(os.Getenv("OTP_EXPIRY_MINUTES"))
	if v == 0 {
		return 5
	}
	return v
}

func (s *AuthService) getOTPMaxAttempts() int {
	v, _ := strconv.Atoi(os.Getenv("OTP_MAX_ATTEMPTS"))
	if v == 0 {
		return 5
	}
	return v
}
