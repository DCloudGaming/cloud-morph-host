package errors

import (
	"errors"
	"net/http"
)

var (
	UserNotFound = errors.New("Invalid User")
	BadRequestMethod = errors.New(http.StatusText(http.StatusMethodNotAllowed))
	InternalError    = errors.New(http.StatusText(http.StatusInternalServerError))

	NoJSONBody   = errors.New("Unable to decode JSON")
	BadCSRF           = errors.New("Missing CSRF Header")
	BadOrigin         = errors.New("Invalid Origin Header")
	RouteUnauthorized = errors.New("You don't have permission to view this resource")
	RouteNotFound     = errors.New("Route not found")
	ExpiredToken      = errors.New("Your access token expired")
	InvalidToken      = errors.New("Your access token is invalid")
)

// codeMap returns a map of errors to http status codes
func codeMap() map[error]int {
	return map[error]int{
		BadRequestMethod: http.StatusMethodNotAllowed,
		InternalError:    http.StatusInternalServerError,

		NoJSONBody:   http.StatusBadRequest,
		UserNotFound:         http.StatusNotFound,
		BadCSRF:           http.StatusUnauthorized,
		BadOrigin:         http.StatusUnauthorized,
		RouteUnauthorized: http.StatusUnauthorized,
		RouteNotFound:     http.StatusNotFound,
		ExpiredToken:      http.StatusUnauthorized,
		InvalidToken:      http.StatusUnauthorized,
	}
}

// GetCode is a helper to get the relevant code for an error, or just return 500
func GetCode(e error) (bool, int) {
	if code, ok := codeMap()[e]; ok {
		return true, code
	}
	return false, http.StatusInternalServerError
}
