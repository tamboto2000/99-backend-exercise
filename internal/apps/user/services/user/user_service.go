package user

import (
	"context"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/entities"
)

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserDetail(ctx context.Context, id int64) (entities.User, error)
	GetAllUsersPaginate(ctx context.Context, page, size int) ([]entities.User, error)
}

type userSvc struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userSvc{
		userRepo: userRepo,
	}
}

func (u *userSvc) CreateUser(ctx context.Context, user *User) error {
	return u.userRepo.Create(ctx, user.User())
}

func (u *userSvc) GetUserDetail(ctx context.Context, id int64) (entities.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *userSvc) GetAllUsersPaginate(ctx context.Context, page, size int) ([]entities.User, error) {
	if page < 1 {
		page = 1
	}

	if size < 1 {
		size = 10
	}

	offset := size * (page - 1)

	return u.userRepo.GetAll(ctx, offset, size)
}
