package error_code

/**
 * @author Rancho
 * @date 2022/1/5
 */

// basic error code
var (
    Success         = NewError(0, "success")
    ServerError     = NewError(1000, "server internal error")
    InvalidParams   = NewError(10001, "invalid params")
    NotFound        = NewError(10002, "record not found")
    TooManyRequests = NewError(10003, "too many requests")
)

// auth error code
var (
    UnauthorizedAuthNotExist  = NewError(20001, "unauthorized, auth not exists")
    UnauthorizedTokenError    = NewError(20002, "unauthorized, token invalid")
    UnauthorizedTokenTimeout  = NewError(20003, "unauthorized, token timeout")
    UnauthorizedTokenGenerate = NewError(20004, "unauthorized，token generate failed")
)

// internal error code
var (
    CopyError = NewError(30001, "copy obj error")
    JSONError = NewError(30002, "json marshal/unmarshal error")
)
