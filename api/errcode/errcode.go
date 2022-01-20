package errcode

import (
    "fmt"
    "net/http"
)

/**
 * @author Rancho
 * @date 2022/1/5
 */

type Error struct {
    Code    int      `json:"code"`
    Msg     string   `json:"msg"`
    Details []string `json:"details"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
    if _, ok := codes[code]; ok {
        panic(fmt.Sprintf("error code %d already exists, please replace one", code))
    }

    codes[code] = msg
    return &Error{
        Code: code,
        Msg:  msg,
    }
}

func (e *Error) Error() string {
    return fmt.Sprintf("err_code: %d, err_msg: %s", e.Code, e.Msg)
}

func (e *Error) Msgf(args []interface{}) string {
    return fmt.Sprintf(e.Msg, args...)
}

func (e *Error) WithDetails(details ...string) *Error {
    newError := *e
    newError.Details = []string{}

    for _, d := range details {
        newError.Details = append(newError.Details, d)
    }

    return &newError
}

func (e *Error) StatusCode() int {
    switch e.Code {
    case Success.Code:
        return http.StatusOK
    case ServerError.Code:
        return http.StatusInternalServerError
    case InvalidParams.Code:
        return http.StatusBadRequest
    case NotFound.Code:
        return http.StatusNotFound
    case UnauthorizedAuthNotExist.Code:
        fallthrough
    case UnauthorizedTokenError.Code:
        fallthrough
    case UnauthorizedTokenGenerate.Code:
        fallthrough
    case UnauthorizedTokenTimeout.Code:
        return http.StatusUnauthorized
    case TooManyRequests.Code:
        return http.StatusTooManyRequests
    }

    return http.StatusInternalServerError
}
