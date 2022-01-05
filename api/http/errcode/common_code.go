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
    UnauthorizedAuthNotExist  = NewError(10000003, "鉴权失败，找不到对应的AppKey和AppSecret")
    UnauthorizedTokenError    = NewError(10000004, "鉴权失败，Token错误")
    UnauthorizedTokenTimeout  = NewError(10000005, "鉴权失败，Token超时")
    UnauthorizedTokenGenerate = NewError(10000006, "鉴权失败，Token生成失败")
    TooManyRequests           = NewError(10000007, "too many requests")
)
