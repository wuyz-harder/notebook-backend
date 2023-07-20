package common

import "net/http"

// 错误处理的结构体
type ApiError struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

// 错误处理的结构体
type ApiSuccess struct {
	StatusCode int         `json:"-"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

var (
	ServerError = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	// ParamsError = NewError(http.StatusPaymentRequired, 200400, "参数错误")
	NotFound = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
)

func OtherError(message string) *ApiError {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *ApiError) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}

func NewResp(statusCode, Code int, mes string, Data map[string]interface{}) *ApiSuccess {
	return &ApiSuccess{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        mes,
		Data:       Data,
	}
}
