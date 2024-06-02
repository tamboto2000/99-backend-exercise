package aggregator

import (
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/external/listing"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/external/user"
)

type Listing struct {
	listing.Listing
	User user.User `json:"user"`
}
