package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/config"
	"github.com/tamboto2000/99-backend-exercise/internal/common/errors"
	"github.com/tamboto2000/99-backend-exercise/pkg/gocb"
	"github.com/tamboto2000/99-backend-exercise/pkg/httpx"
)

const (
	userSvcCfgKey        = "user_svc"
	createUserEndpKey    = "create_user"
	getUserDetailEndpKey = "get_user_detail"
)

type UserService interface {
	CreateUser(ctx context.Context, name string) (User, error)
	GetUserDetail(ctx context.Context, id int64) (User, error)
}

type userSvc struct {
	cl       httpx.Client
	host     string
	endpsCfg config.Endpoints
}

func NewUserService(cfg config.Config) UserService {
	svcCfg := cfg.Services.Get(userSvcCfgKey)

	endpsCfg := svcCfg.Endpoints

	// TODO: configuration for cb and http client
	// can be put into config
	cb := gocb.NewCircuitBreaker("user-svc", gocb.Settings{
		ErrThreshold: 10,
		ErrInterval:  1 * time.Minute,
		Timeout:      1 * time.Minute,
		Retry:        3,
	})

	cl := http.DefaultClient
	cl.Timeout = 10 * time.Second

	return &userSvc{
		cl:       httpx.NewClientWithCB(cl, cb),
		host:     svcCfg.Host,
		endpsCfg: endpsCfg,
	}
}

func (u *userSvc) CreateUser(ctx context.Context, name string) (User, error) {
	endp := fmt.Sprintf("%s%s", u.host, u.endpsCfg.Get(createUserEndpKey))

	var respBody userResponse
	var user User

	data := make(url.Values)
	data.Set("name", name)
	reqBody := bytes.NewReader([]byte(data.Encode()))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endp, reqBody)
	if err != nil {
		return user, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := u.cl.Do(req)
	if err != nil {
		return user, err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return user, err
	}

	user = respBody.User

	if resp.StatusCode != http.StatusOK {
		errResp := respBody.Error
		if errResp.Code == errors.CodeValidation {
			return user, errors.NewErrValidation(errResp.Msg, errResp.Fields)
		}

		return user, errors.New(errResp.Msg, errResp.Code)
	}

	return user, nil
}

func (u *userSvc) GetUserDetail(ctx context.Context, id int64) (User, error) {
	endp := fmt.Sprintf("%s%s", u.host, u.endpsCfg.Get(getUserDetailEndpKey))
	endp = fmt.Sprintf(endp, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endp, nil)
	if err != nil {
		return User{}, err
	}

	resp, err := u.cl.Do(req)
	if err != nil {
		return User{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusInternalServerError {
		return User{}, errors.New("internal server error", errors.CodeInternal)
	}

	var respData userResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return User{}, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return User{}, errors.NewErrValidation(respData.Error.Msg, respData.Error.Fields)
		}

		if resp.StatusCode == http.StatusNotFound {
			return User{}, errors.NewErrNotExists(respData.Error.Msg)
		}
	}

	return respData.User, nil
}
