package error_code

import (
	"fmt"
	"net/http"
)

// Error represents a standardized API error
type Error struct {
	Code    int      `json:"code"`              // Error code
	Msg     string   `json:"message"`           // Error message
	Details []string `json:"details,omitempty"` // Optional error details
	HTTP    int      `json:"-"`                 // HTTP status code (not exposed in JSON)
	DocRef  string   `json:"doc_ref,omitempty"` // Reference to documentation
}

var codes = map[int]string{}

// Basic error code
const (
	SuccessCode = 0

	ServerErrorCode     = 10000
	InvalidParamsCode   = 10001
	NotFoundCode        = 10002
	TooManyRequestsCode = 10003

	UnauthorizedAuthNotExistErrorCode  = 20001
	UnauthorizedTokenErrorCode         = 20002
	UnauthorizedTokenTimeoutErrorCode  = 20003
	UnauthorizedTokenGenerateErrorCode = 20004

	CopyErrorErrorCode = 30001
	JSONErrorErrorCode = 30002

	AccountExistErrorCode  = 40001
	UserNameExistErrorCode = 40002
)

// API error code
var (
	Success         = NewError(SuccessCode, "success")
	ServerError     = NewError(ServerErrorCode, "server internal error")
	InvalidParams   = NewError(InvalidParamsCode, "invalid params")
	NotFound        = NewError(NotFoundCode, "record not found")
	TooManyRequests = NewError(TooManyRequestsCode, "too many requests")
)

// Auth error code
var (
	UnauthorizedAuthNotExist  = NewError(UnauthorizedAuthNotExistErrorCode, "unauthorized, auth not exists")
	UnauthorizedTokenError    = NewError(UnauthorizedTokenErrorCode, "unauthorized, token invalid")
	UnauthorizedTokenTimeout  = NewError(UnauthorizedTokenTimeoutErrorCode, "unauthorized, token timeout")
	UnauthorizedTokenGenerate = NewError(UnauthorizedTokenGenerateErrorCode, "unauthorized, token generate failed")
)

// Internal error code
var (
	CopyError = NewError(CopyErrorErrorCode, "copy obj error")
	JSONError = NewError(JSONErrorErrorCode, "json marshal/unmarshal error")
)

// Business error code
var (
	AccountExist  = NewError(AccountExistErrorCode, "account already exists")
	UserNameExist = NewError(UserNameExistErrorCode, "username already exists")
)

// NewError creates a new Error instance with the specified code and message
func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("error code %d already exists, please replace one", code))
	}

	codes[code] = msg
	return &Error{
		Code: code,
		Msg:  msg,
		HTTP: determineHTTPStatusCode(code), // Determine default HTTP status code
	}
}

// NewErrorWithStatus creates a new Error instance with the specified code, message, and HTTP status
func NewErrorWithStatus(code int, msg string, status int) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("error code %d already exists, please replace one", code))
	}

	codes[code] = msg
	return &Error{
		Code: code,
		Msg:  msg,
		HTTP: status,
	}
}

// Error implements the error interface
func (e Error) Error() string {
	return fmt.Sprintf("err_code: %d, err_msg: %s, details: %v", e.Code, e.Msg, e.Details)
}

// Msgf formats the error message with the provided arguments
func (e Error) Msgf(args []any) string {
	return fmt.Sprintf(e.Msg, args...)
}

// WithDetails adds error details
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.Details = []string{}
	newError.Details = append(newError.Details, details...)

	return &newError
}

// WithDocRef adds a documentation reference
func (e *Error) WithDocRef(docRef string) *Error {
	newError := *e
	newError.DocRef = docRef
	return &newError
}

// WithMessage customizes the error message
func (e *Error) WithMessage(format string, args ...interface{}) *Error {
	newError := *e
	newError.Msg = fmt.Sprintf(format, args...)
	return &newError
}

// Is implements error comparison
func (e *Error) Is(tgt error) bool {
	target, ok := tgt.(*Error)
	if !ok {
		return false
	}
	return target.Code == e.Code
}

// StatusCode returns the HTTP status code
func (e *Error) StatusCode() int {
	if e.HTTP != 0 {
		return e.HTTP
	}
	return determineHTTPStatusCode(e.Code)
}

// determineHTTPStatusCode maps error codes to HTTP status codes
func determineHTTPStatusCode(code int) int {
	switch code {
	case SuccessCode:
		return http.StatusOK
	case ServerErrorCode:
		return http.StatusInternalServerError
	case InvalidParamsCode:
		return http.StatusBadRequest
	case NotFoundCode:
		return http.StatusNotFound
	case UnauthorizedAuthNotExistErrorCode,
		UnauthorizedTokenErrorCode,
		UnauthorizedTokenGenerateErrorCode,
		UnauthorizedTokenTimeoutErrorCode:
		return http.StatusUnauthorized
	case TooManyRequestsCode:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
