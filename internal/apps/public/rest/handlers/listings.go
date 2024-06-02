package handlers

import (
	"net/http"
	"strconv"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/rest/messages"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/services/aggregator"
	"github.com/tamboto2000/99-backend-exercise/internal/common/errors"
	"github.com/tamboto2000/99-backend-exercise/internal/common/response"
	"github.com/tamboto2000/99-backend-exercise/pkg/logger"
	"github.com/tamboto2000/99-backend-exercise/pkg/mux"
)

func CreateListing(svc aggregator.AggregatorService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		var reqBody messages.CreateListingRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			err := errors.New(response.MsgBadRequest, errors.CodeValidation)
			return response.GiveErrResponse(ctx, err)
		}

		ls, err := svc.CreateListing(ctx.RequestContext(), reqBody.UserID, reqBody.ListingType, reqBody.Price)
		if err != nil {
			// DELET
			logger.Error(err.Error())

			return response.GiveErrResponse(ctx, err)
		}

		respBody := messages.CreateListingResponse{
			Listing: ls,
		}

		return ctx.WriteJSON(http.StatusOK, respBody)
	}
}

func GetAllListing(svc aggregator.AggregatorService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		pageNumStr := ctx.QueryVal("page_num")
		pageSizeStr := ctx.QueryVal("page_size")

		fields := make(errors.Fields)
		pageNum := validatePageNum(fields, pageNumStr)
		pageSize := validatePageSize(fields, pageSizeStr)

		if fields.NotEmpty() {
			return response.GiveResultErrResponse(ctx, errors.NewErrValidation(response.MsgBadRequest, fields))
		}

		ls, err := svc.GetAllListing(ctx.RequestContext(), pageNum, pageSize)
		if err != nil {
			return response.GiveResultErrResponse(ctx, err)
		}

		respBody := messages.GetAllListingResponse{
			ResultResponse: response.ResultResponse{
				Result: true,
			},
			Listings: ls,
		}

		return ctx.WriteJSON(http.StatusOK, respBody)
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
