package response

import "github.com/snlansky/coral/pkg/errors"

const SuccessCode = 0

const SuccessMsg = "ok"

type JsonResponse struct {
	ErrorCode   int         `json:"errcode"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

func New(data interface{}) *JsonResponse {
	resp := JsonResponse{
		ErrorCode:   SuccessCode,
		Description: SuccessMsg,
		Data:        data,
	}
	return &resp
}

func Fail(code int, desc string) *JsonResponse {
	return &JsonResponse{
		ErrorCode:   code,
		Description: desc,
		Data:        nil,
	}
}

func Err(err error) *JsonResponse {
	if e, ok := err.(errors.CodeError); ok {
		return &JsonResponse{
			ErrorCode:   e.Code(),
			Description: e.Error(),
			Data:        nil,
		}
	}

	return &JsonResponse{
		ErrorCode:   errors.InternalErrorCode,
		Description: err.Error(),
	}

}
