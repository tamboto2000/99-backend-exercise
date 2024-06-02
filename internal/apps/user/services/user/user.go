package user

import (
	"time"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/entities"
	"github.com/tamboto2000/99-backend-exercise/internal/common/errors"
	"github.com/tamboto2000/99-backend-exercise/pkg/snowid"
)

type User struct {
	user entities.User
}

func NewUser(name string) (*User, error) {
	id := snowid.Generate()
	fields := make(errors.Fields)

	if name == "" {
		fields.Add("name", "can not be empty")
		return nil, errors.NewErrValidation("invalid input", fields)
	}

	if len(name) > 100 {
		fields.Add("name", "maximum length is 100 characters")
		return nil, errors.NewErrValidation("invalid input", nil)
	}

	now := time.Now()
	user := entities.User{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return &User{user}, nil
}

func (u User) User() entities.User {
	return u.user
}
