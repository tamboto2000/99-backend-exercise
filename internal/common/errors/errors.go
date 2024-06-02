package errors

import "errors"

const (
	ValidationErr = iota
)

type Err struct {
	msg  string
	code int
}

func New(msg string, code int) error {
	return Err{
		msg:  msg,
		code: code,
	}
}

func (e Err) Error() string {
	return e.msg
}

func (e Err) Code() int {
	return e.code
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
