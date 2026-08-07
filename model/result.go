package model

import (
	"net/http"
	"time"
)

// ResultJson 返回结果
type ResultJson struct {
	Code    int       `json:"code"`
	Status  bool      `json:"status"`
	Message string    `json:"message"`
	Result  any       `json:"result"`
	Date    time.Time `json:"date"`
}

// ResultJsonSuccess 返回成功结果
func ResultJsonSuccess() ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "succeess",
		Status:  true,
		Result:  nil,
		Date:    time.Now(),
	}
}

// ResultJsonSuccessWithData 返回成功结果
func ResultJsonSuccessWithData(data any) ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "success",
		Status:  true,
		Result:  data,
		Date:    time.Now(),
	}
}

// ResultJsonError 返回错误结果
func ResultJsonError(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusInternalServerError,
		Message: message,
		Status:  true,
		Result:  nil,
		Date:    time.Now(),
	}
}

// ResultJsonBadRequest 返回错误结果
func ResultJsonBadRequestq(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusBadRequest,
		Message: message,
		Status:  false,
		Result:  nil,
		Date:    time.Now(),
	}
}

// ResultJsonUnauthorized 返回错误结果
func ResultJsonUnauthorized(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusUnauthorized,
		Message: message,
		Status:  false,
		Result:  nil,
		Date:    time.Now(),
	}
}
