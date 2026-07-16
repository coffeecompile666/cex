package user

import (
	"errors"
	"net/http"

	"icon_exchange/internal/shared"
	"icon_exchange/internal/user/repository"
	"icon_exchange/internal/user/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	authService *service.AuthService
}

func NewHandler(authService *service.AuthService) *Handler {
	return &Handler{authService: authService}
}

// RegisterRoutes registers all auth endpoints.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.Signup)
		auth.POST("/signup/verify", h.VerifyOTP)
		auth.POST("/signup/resend", h.ResendOTP)
		auth.POST("/signin", h.Signin)
		auth.POST("/refresh-token", h.RefreshToken)
	}
}

// ─── Request DTOs ──────────────────────────────────────────────────────────────

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SigninRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ─── Response helpers ─────────────────────────────────────────────────────────

// respondError checks if err is a shared.Error and uses its Status + Code.
// Falls back to 500 for unknown errors.
func respondError(c *gin.Context, err error) {
	var appErr shared.Error
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    shared.ErrInternalServerError.Code,
		"message": shared.ErrInternalServerError.Message,
	})
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

// Signup handles POST /auth/signup
func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, shared.ErrBadRequest)
		return
	}

	if err := h.authService.Signup(req.Email, req.Password); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email. Please verify to complete registration."})
}

// VerifyOTP handles POST /auth/signup/verify
func (h *Handler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, shared.ErrBadRequest)
		return
	}

	tokenPair, err := h.authService.VerifyOTP(req.Email, req.Code)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// ResendOTP handles POST /auth/signup/resend
func (h *Handler) ResendOTP(c *gin.Context) {
	var req ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, shared.ErrBadRequest)
		return
	}

	if err := h.authService.ResendOTP(req.Email); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "A new OTP has been sent to your email."})
}

// Signin handles POST /auth/signin
func (h *Handler) Signin(c *gin.Context) {
	var req SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, shared.ErrBadRequest)
		return
	}

	tokenPair, err := h.authService.Signin(req.Email, req.Password)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// RefreshToken handles POST /auth/refresh-token
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, shared.ErrBadRequest)
		return
	}

	tokenPair, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// ─── Backward compat ─────────────────────────────────────────────────────────

func NewRepository(db *gorm.DB) *repository.UserRepo {
	return repository.NewUserRepo(db)
}
