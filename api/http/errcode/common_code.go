package errcode

/**
 * @author Rancho
 * @date 2022/1/5
 */

var (
    Success                   = NewError(0, "success")
    ServerError               = NewError(10000000, "server internal error")
    InvalidParams             = NewError(10000001, "invalid params")
    NotFound                  = NewError(10000002, "record not found")
    UnauthorizedAuthNotExist  = NewError(10000003, "unauthorized, auth not exists")
    UnauthorizedTokenError    = NewError(10000004, "unauthorized, token invalid")
    UnauthorizedTokenTimeout  = NewError(10000005, "unauthorized, token timeout")
    UnauthorizedTokenGenerate = NewError(10000006, "unauthorized，token generate failed")
    TooManyRequests           = NewError(10000007, "too many requests")
)
