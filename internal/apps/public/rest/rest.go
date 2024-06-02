package rest

import (
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/modules"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/rest/handlers"
	"github.com/tamboto2000/99-backend-exercise/pkg/mux"
)

func RegisterRoutes(r *mux.Router, svcs modules.Services) {
	pb := r.SubRouter("/public-api")

	// create user
	pb.Post("/users", handlers.CreateUser(svcs.AggregatorService))

	// create listing
	pb.Post("/listings", handlers.CreateListing(svcs.AggregatorService))

	// get all listing
	pb.Get("/listings", handlers.GetAllListing(svcs.AggregatorService))
}
