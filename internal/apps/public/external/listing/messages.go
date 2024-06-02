package listing

import "github.com/tamboto2000/99-backend-exercise/internal/common/response"

type Listing struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	ListingType string `json:"listing_type"`
	Price       int64  `json:"price"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type createListingResponse struct {
	response.ResultResponse
	Listing Listing  `json:"listing"`
	Errors  []string `json:"errors"`
}

type getAllListingResponse struct {
	response.ResultResponse
	Listings []Listing `json:"listings"`
}
