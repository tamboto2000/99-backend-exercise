package modules

import (
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/config"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/external/listing"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/external/user"
)

type ExternalServices struct {
	UserService    user.UserService
	ListingService listing.ListingService
}

func NewExternalServices(cfg config.Config) ExternalServices {
	// user service
	userSvc := user.NewUserService(cfg)

	// listing service
	listingSvc := listing.NewListingSvc(cfg)

	return ExternalServices{
		UserService:    userSvc,
		ListingService: listingSvc,
	}
}
