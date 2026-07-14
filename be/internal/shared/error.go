package shared

import "net/http"

// Error is the base application error type.
// It carries an HTTP status, an internal error code, and a human-readable message.
type Error struct {
	Message string
	Code    int
	Status  int
}

// Error implements the built-in error interface.
func (e Error) Error() string {
	return e.Message
}

// ─── Common / Base Errors (1xxxxx) ─────────────────────────────────────────────

var (
	ErrInternalServerError Error = Error{Code: 100000, Message: "internal server error", Status: http.StatusInternalServerError}
	ErrNotFound            Error = Error{Code: 100001, Message: "not found", Status: http.StatusNotFound}
	ErrBadRequest          Error = Error{Code: 100002, Message: "bad request", Status: http.StatusBadRequest}
	ErrUnauthenticated     Error = Error{Code: 100003, Message: "unauthenticated", Status: http.StatusUnauthorized}
	ErrForbidden           Error = Error{Code: 100004, Message: "forbidden", Status: http.StatusForbidden}
	ErrConflict            Error = Error{Code: 100005, Message: "conflict", Status: http.StatusConflict}
	ErrTooManyRequests     Error = Error{Code: 100006, Message: "too many requests", Status: http.StatusTooManyRequests}
	ErrGone                Error = Error{Code: 100007, Message: "gone", Status: http.StatusGone}
)

// ─── Auth Errors (2xxxxx) ───────────────────────────────────────────────────────

var (
	// Signup
	ErrEmailAlreadyRegistered Error = Error{Code: 200001, Message: "email already registered", Status: http.StatusConflict}
	ErrUserNotVerified        Error = Error{Code: 200002, Message: "user not verified, please check your email for OTP", Status: http.StatusForbidden}

	// Credentials
	ErrInvalidCredentials Error = Error{Code: 200003, Message: "invalid email or password", Status: http.StatusUnauthorized}

	// OTP
	ErrOTPNotFound    Error = Error{Code: 200010, Message: "no active OTP found", Status: http.StatusGone}
	ErrOTPExpired     Error = Error{Code: 200011, Message: "OTP has expired", Status: http.StatusGone}
	ErrOTPRevoked     Error = Error{Code: 200012, Message: "OTP has been revoked", Status: http.StatusGone}
	ErrOTPMaxAttempts Error = Error{Code: 200013, Message: "OTP locked: too many failed attempts", Status: http.StatusTooManyRequests}
	ErrOTPInvalid     Error = Error{Code: 200014, Message: "invalid OTP code", Status: http.StatusBadRequest}

	// Token
	ErrInvalidToken Error = Error{Code: 200020, Message: "invalid or expired refresh token", Status: http.StatusUnauthorized}
)
