package httpx

import (
	"errors"
	"net/http"

	"github.com/tamboto2000/99-backend-exercise/pkg/gocb"
)

var ErrServerDown = errors.New("server is down")

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type clientWithCb struct {
	cl *http.Client
	cb gocb.Circuit
}

func NewClientWithCB(cl *http.Client, cb gocb.Circuit) Client {
	return &clientWithCb{
		cl: cl,
		cb: cb,
	}
}

func (cwc *clientWithCb) Do(req *http.Request) (resp *http.Response, err error) {
	err = cwc.cb.Run(func() error {
		resp, err = cwc.cl.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode >= http.StatusInternalServerError {
			return ErrServerDown
		}

		return err
	})

	if err == gocb.ErrCircuitOpen {
		return nil, ErrServerDown
	}

	return
}
