package error_code

// basic error code
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
