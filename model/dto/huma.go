package dto

import (
	"net/http"

	"github.com/rfancn/airpress/util/xerr"
)

// HumaOut 是 huma handler 的响应包装。
// huma 约定：Output 结构体的 Status 字段作为 HTTP 状态码，Body 字段作为响应体序列化。
// 这里 Body 统一为 BaseDTO 信封，保持与迁移前一致的响应格式。
type HumaOut[T any] struct {
	Status int        `json:"-" status:"200"`
	Body   BaseDTO[T] `json:"-"`
}

// HumaOK 把业务数据包成成功响应（HTTP 200，BaseDTO 信封）。
func HumaOK[T any](data T) (*HumaOut[T], error) {
	return &HumaOut[T]{
		Status: http.StatusOK,
		Body:   BaseDTO[T]{Status: http.StatusOK, Data: data, Message: "OK"},
	}, nil
}

// HumaErr 把业务错误包成失败响应（HTTP 状态码取自 xerr，BaseDTO 信封）。
// 注意：handler 返回 nil 的 error，状态码由 HumaOut.Status 控制，避免 huma 输出 RFC9457 格式。
func HumaErr[T any](err error) (*HumaOut[T], error) {
	status := xerr.GetHTTPStatus(err)
	return &HumaOut[T]{
		Status: status,
		Body:   BaseDTO[T]{Status: status, Message: xerr.GetMessage(err)},
	}, nil
}
