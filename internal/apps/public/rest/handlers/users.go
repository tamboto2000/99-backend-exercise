package handlers

import (
	"net/http"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/rest/messages"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/services/aggregator"
	"github.com/tamboto2000/99-backend-exercise/internal/common/errors"
	"github.com/tamboto2000/99-backend-exercise/internal/common/response"
	"github.com/tamboto2000/99-backend-exercise/pkg/mux"
)

func CreateUser(svc aggregator.AggregatorService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		var reqBody messages.CreateUserRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			err := errors.New(response.MsgBadRequest, errors.CodeValidation)
			return response.GiveErrResponse(ctx, err)
		}

		u, err := svc.CreateUser(ctx.RequestContext(), reqBody.Name)
		if err != nil {
			return response.GiveErrResponse(ctx, err)
		}

		resp := messages.CreateUserResponse{User: u}

		return ctx.WriteJSON(http.StatusOK, resp)
	}
}
