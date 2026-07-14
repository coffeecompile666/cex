package shared

import "net/http"

type Error struct {
	Message string
	Code    int
	Status  int
}

var (
	ErrInternalServerError Error = Error{Code: 100000, Message: "internal server error", Status: http.StatusInternalServerError}
	ErrNotFound            Error = Error{Code: 100001, Message: "not found", Status: http.StatusNotFound}
	ErrBadRequest          Error = Error{Code: 100002, Message: "bad request", Status: http.StatusBadRequest}
	ErrUnauthenticated     Error = Error{Code: 100003, Message: "unauthenticated", Status: http.StatusUnauthorized}
	ErrForbidden           Error = Error{Code: 100004, Message: "forbidden", Status: http.StatusForbidden}
)
