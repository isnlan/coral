package errors

import (
	"github.com/pkg/errors"
)

var (
	New         = errors.New
	Wrap        = errors.Wrap
	Wrapf       = errors.Wrapf
	Errorf      = errors.Errorf
	WithStack   = errors.WithStack
	WithMessage = errors.WithMessage
	Cause       = errors.Cause
)

type CodeError interface {
	error
	Code() int
}

const (
	InternalErrorCode = 500 // 内部错误
)

var InternalError = &SvrError{code: InternalErrorCode, info: "internal error"}

type SvrError struct {
	code int
	info string
}

func (e *SvrError) Error() string {
	return e.info
}

func (e *SvrError) Code() int {
	return e.code
}

func NewWithCode(err error, code int) CodeError {
	if e, ok := err.(CodeError); ok {
		return e
	}

	return &SvrError{
		code: code,
		info: err.Error(),
	}
}

func NewWithInfo(info string, code int) CodeError {
	return &SvrError{
		code: code,
		info: info,
	}
}
