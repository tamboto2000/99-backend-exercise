package messages

import (
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/external/listing"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/external/user"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/services/aggregator"
	"github.com/tamboto2000/99-backend-exercise/internal/common/response"
)

type CreateListingRequest struct {
	UserID      int64  `json:"user_id"`
	ListingType string `json:"listing_type"`
	Price       int64  `json:"price"`
}

type CreateListingResponse struct {
	Listing listing.Listing `json:"listing"`
}

type ListingWithUser struct {
	listing.Listing
	User user.User `json:"user"`
}

type GetAllListingResponse struct {
	response.ResultResponse
	Listings []aggregator.Listing `json:"listings"`
}
