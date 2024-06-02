package user

import "github.com/tamboto2000/99-backend-exercise/internal/common/errors"

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Error struct {
	Code   int           `json:"code"`
	Msg    string        `json:"msg"`
	Fields errors.Fields `json:"fields,omitempty"`
}

type userResponse struct {
	Result bool  `json:"result"`
	User   User  `json:"user"`
	Error  Error `json:"error"`
}
