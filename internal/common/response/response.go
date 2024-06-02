package response

import (
	"net/http"

	"github.com/tamboto2000/99-backend-exercise/internal/common/errors"
	"github.com/tamboto2000/99-backend-exercise/pkg/mux"
)

const (
	StatusError   = "error"
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

const (
	MsgBadRequest  = "bad request"
	MsgInternalErr = "internal server error"
)

type Error struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Fields any    `json:"fields,omitempty"`
}

type ResultResponse struct {
	Result bool   `json:"result"`
	Error  *Error `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

func GiveResultErrResponse(ctx *mux.Context, err error) error {
	errMsg, status := BuildErrMessage(err)
	resp := ResultResponse{
		Result: false,
		Error:  &errMsg,
	}

	return ctx.WriteJSON(status, resp)
}

func GiveErrResponse(ctx *mux.Context, err error) error {
	errMsg, status := BuildErrMessage(err)
	errResp := ErrorResponse{
		Error: errMsg,
	}
	return ctx.WriteJSON(status, errResp)
}

func BuildErrMessage(err error) (errMsg Error, status int) {
	errMsg = errToErrResp(err)
	status = http.StatusInternalServerError

	switch errMsg.Code {
	case errors.CodeValidation:
		status = http.StatusBadRequest

	case errors.CodeAlreadyExists:
		status = http.StatusConflict

	case errors.CodeNotExists:
		status = http.StatusNotFound

	case errors.CodeInvalidAuth:
		status = http.StatusUnauthorized

	case errors.CodeLimitExceeded:
		status = http.StatusTooManyRequests

	default:
		errMsg.Msg = MsgInternalErr
	}

	return
}

func errToErrResp(err error) Error {
	var errVld errors.ErrValidation
	if errors.As(err, &errVld) {
		return Error{
			Msg:    errVld.Error(),
			Code:   errVld.Code(),
			Fields: errVld.Fields(),
		}
	}

	var errE errors.Err
	if errors.As(err, &errE) {
		return Error{
			Msg:  errE.Error(),
			Code: errE.Code(),
		}
	}

	return Error{
		Msg:  err.Error(),
		Code: errors.CodeInternal,
	}
}
