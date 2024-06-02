package handlers

import (
	"net/http"
	"strconv"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/entities"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/rest/messages"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/services/user"
	"github.com/tamboto2000/99-backend-exercise/internal/common/errors"
	"github.com/tamboto2000/99-backend-exercise/internal/common/response"
	"github.com/tamboto2000/99-backend-exercise/pkg/mux"
)

type createUserResponse struct {
	response.ResultResponse
	User messages.User `json:"user"`
}

type getUserDetailResponse struct {
	response.ResultResponse
	User messages.User `json:"user"`
}

type getAllUsersPaginateResponse struct {
	response.ResultResponse
	Users []messages.User `json:"users"`
}

func CreateUser(svc user.UserService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		name := ctx.FormVal("name")
		u, err := user.NewUser(name)
		if err != nil {
			return response.GiveResultErrResponse(ctx, err)
		}

		err = svc.CreateUser(ctx.RequestContext(), u)
		if err != nil {
			return response.GiveResultErrResponse(ctx, err)
		}

		ue := u.User()
		resp := createUserResponse{
			ResultResponse: response.ResultResponse{Result: true},
			User:           userEntityToMessage(ue),
		}

		return ctx.WriteJSON(http.StatusOK, resp)
	}
}

func GetUserDetail(svc user.UserService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		idStr := ctx.PathVal("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			fields := make(errors.Fields)
			fields.Add("id", "must be integer")

			return response.GiveResultErrResponse(ctx, errors.NewErrValidation(response.MsgBadRequest, fields))
		}

		ue, err := svc.GetUserDetail(ctx.RequestContext(), id)
		if err != nil {
			return response.GiveResultErrResponse(ctx, err)
		}

		resp := getUserDetailResponse{
			ResultResponse: response.ResultResponse{Result: true},
			User:           userEntityToMessage(ue),
		}

		return ctx.WriteJSON(http.StatusOK, resp)
	}
}

func GetAllUsers(svc user.UserService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		pageNumStr := ctx.QueryVal("page_num")
		pageSizeStr := ctx.QueryVal("page_size")

		fields := make(errors.Fields)
		pageNum := validatePageNum(fields, pageNumStr)
		pageSize := validatePageSize(fields, pageSizeStr)

		if fields.NotEmpty() {
			return response.GiveResultErrResponse(ctx, errors.NewErrValidation(response.MsgBadRequest, fields))
		}

		us, err := svc.GetAllUsersPaginate(ctx.RequestContext(), pageNum, pageSize)
		if err != nil {
			return response.GiveResultErrResponse(ctx, err)
		}

		resp := getAllUsersPaginateResponse{
			ResultResponse: response.ResultResponse{Result: true},
			Users: func(us []entities.User) []messages.User {
				var userResps []messages.User
				for _, u := range us {
					userResps = append(userResps, userEntityToMessage(u))
				}

				return userResps
			}(us),
		}

		return ctx.WriteJSON(http.StatusOK, resp)
	}
}

func validatePageNum(fields errors.Fields, pageNumStr string) int {
	if pageNumStr == "" {
		return 0
	}

	field := "page_num"

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		fields.Add(field, "must be integer")
	}

	return pageNum
}

func validatePageSize(fields errors.Fields, pageSizeStr string) int {
	if pageSizeStr == "" {
		return 0
	}

	field := "page_size"

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		fields.Add(field, "must be integer")
	}

	return pageSize
}

func userEntityToMessage(user entities.User) messages.User {
	return messages.User{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.UnixNano(),
		UpdatedAt: user.UpdatedAt.UnixNano(),
	}
}
