package error_code

/**
 * @author Rancho
 * @date 2022/1/5
 */

var (
    Success                   = NewError(0, "success")
    ServerError               = NewError(1000, "server internal error")
    InvalidParams             = NewError(10001, "invalid params")
    NotFound                  = NewError(10002, "record not found")
    UnauthorizedAuthNotExist  = NewError(10003, "unauthorized, auth not exists")
    UnauthorizedTokenError    = NewError(10004, "unauthorized, token invalid")
    UnauthorizedTokenTimeout  = NewError(10005, "unauthorized, token timeout")
    UnauthorizedTokenGenerate = NewError(10006, "unauthorized，token generate failed")
    TooManyRequests           = NewError(10007, "too many requests")
)
