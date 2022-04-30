package base

// @Author: Derek
// @Description: Error
// @Date: 2022/4/30 10:42
// @Version 1.0

import (
	"fmt"

	"github.com/pkg/errors"
)

func NewError(code int, message string) Error {
	return &baseError{
		ErrNo:  code,
		ErrMsg: message,
	}
}

type (
	Error interface {
		error
		SetErrMsg(errMsg string, v ...interface{})
		Sprintf(v ...interface{}) Error
		Equal(e error) bool
		Wrap(core error) error
		WrapPrintf(core error, format string, message ...interface{}) error
	}

	baseError struct {
		ErrNo  int
		ErrMsg string
	}
)

func (err baseError) Error() string {
	return err.ErrMsg
}

func (err *baseError) SetErrMsg(format string, v ...interface{}) {
	err.ErrMsg = fmt.Sprintf(format, v...)
}

func (err *baseError) Sprintf(v ...interface{}) Error {
	err.ErrMsg = fmt.Sprintf(err.ErrMsg, v...)
	return err
}

func (err *baseError) Equal(e error) bool {
	b, ok := e.(*baseError)
	if !ok {
		return false
	}

	return b.ErrNo == err.ErrNo // && b.ErrMsg == err.ErrMsg
}

func (err *baseError) WrapPrintf(core error, format string, message ...interface{}) error {
	if core == nil {
		return nil
	}

	err.ErrMsg = core.Error()
	return errors.Wrap(err, fmt.Sprintf(format, message...))
}

func (err *baseError) Wrap(core error) error {
	if core == nil {
		return nil
	}

	msg := err.ErrMsg
	err.ErrMsg = core.Error()
	return errors.Wrap(err, msg)
}
