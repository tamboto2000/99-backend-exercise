package messages

import "github.com/tamboto2000/99-backend-exercise/internal/apps/public/external/user"

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	User user.User `json:"user"`
}
