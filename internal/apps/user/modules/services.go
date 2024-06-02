package modules

import "github.com/tamboto2000/99-backend-exercise/internal/apps/user/services/user"

type Services struct {
	UserService user.UserService
}

func NewServices(inf *Infra) Services {
	// user service
	userRepo := user.NewUserRepository(inf.DB)
	userSvc := user.NewUserService(userRepo)

	svcs := Services{
		UserService: userSvc,
	}

	return svcs
}
