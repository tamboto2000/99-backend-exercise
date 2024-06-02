package aggregator

import (
	"context"
	"fmt"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/external/listing"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/external/user"
	"github.com/tamboto2000/99-backend-exercise/internal/common/errors"
	"github.com/tamboto2000/99-backend-exercise/pkg/logger"
)

type AggregatorService interface {
	CreateUser(ctx context.Context, name string) (user.User, error)
	CreateListing(ctx context.Context, userId int64, lt string, price int64) (listing.Listing, error)
	GetAllListing(ctx context.Context, pageNum, pageSize int) ([]Listing, error)
}

type aggrSvc struct {
	userSvc    user.UserService
	listingSvc listing.ListingService
}

func NewAggregatorService(userSvc user.UserService, listingSvc listing.ListingService) AggregatorService {
	return &aggrSvc{
		userSvc:    userSvc,
		listingSvc: listingSvc,
	}
}

func (aggr *aggrSvc) CreateUser(ctx context.Context, name string) (user.User, error) {
	u, err := aggr.userSvc.CreateUser(ctx, name)
	if err != nil {
		logger.Error(err.Error())
	}

	return u, err
}

func (aggr *aggrSvc) CreateListing(ctx context.Context, userId int64, lt string, price int64) (listing.Listing, error) {
	// validate if user with userId is exist
	_, err := aggr.userSvc.GetUserDetail(ctx, userId)
	if err != nil {
		var errNotExists errors.ErrNotExists
		if errors.As(err, &errNotExists) {
			return listing.Listing{}, errors.New("user not found", errors.CodeValidation)
		}

		logger.Error(err.Error())
		return listing.Listing{}, err
	}

	// create listing
	ls, err := aggr.listingSvc.CreateListing(ctx, userId, lt, price)
	if err != nil {
		logger.Error(err.Error())
		return ls, err
	}

	return ls, nil
}

func (aggr *aggrSvc) GetAllListing(ctx context.Context, pageNum, pageSize int) ([]Listing, error) {
	ls, err := aggr.listingSvc.GetAllListing(ctx, pageNum, pageSize)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var lsu []Listing
	for _, l := range ls {
		u, err := aggr.userSvc.GetUserDetail(ctx, l.UserID)
		if err != nil {
			logger.Error(fmt.Sprintf("error get user for listing: %v", err))
			return nil, err
		}

		lsu = append(lsu, Listing{
			Listing: l,
			User:    u,
		})
	}

	return lsu, nil
}
