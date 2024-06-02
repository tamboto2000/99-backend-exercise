package rest

import (
	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/modules"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/rest/handlers"
	"github.com/tamboto2000/99-backend-exercise/pkg/mux"
)

func RegisterRoutes(r *mux.Router, svcs modules.Services) {
	r.Post("/users", handlers.CreateUser(svcs.UserService))
	r.Get("/users/{id}", handlers.GetUserDetail(svcs.UserService))
	r.Get("/users", handlers.GetAllUsers(svcs.UserService))
}
