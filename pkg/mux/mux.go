package mux

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tamboto2000/99-backend-exercise/pkg/logger"
)

type Middleware func(HandleFunc) HandleFunc

type Router struct {
	pfx    string
	srvMux *http.ServeMux
	mddls  []Middleware
	server *http.Server
}

func NewRouter() *Router {
	r := &Router{
		srvMux: http.NewServeMux(),
	}

	return r
}

func (r *Router) SubRouter(pfx string) *Router {
	var mddls []Middleware
	copy(mddls, r.mddls)

	subR := &Router{
		pfx:    r.pfx + pfx,
		srvMux: r.srvMux,
		mddls:  mddls,
	}

	return subR
}

func (r *Router) Use(mddls ...Middleware) {
	r.mddls = append(r.mddls, mddls...)
}

func (r *Router) Post(path string, h HandleFunc, mddls ...Middleware) {
	r.Handle(http.MethodPost, path, h, mddls...)
}

func (r *Router) Get(path string, h HandleFunc, mddls ...Middleware) {
	r.Handle(http.MethodGet, path, h, mddls...)
}

func (r *Router) Patch(path string, h HandleFunc, mddls ...Middleware) {
	r.Handle(http.MethodPatch, path, h, mddls...)
}

func (r *Router) Put(path string, h HandleFunc, mddls ...Middleware) {
	r.Handle(http.MethodPut, path, h, mddls...)
}

func (r *Router) Delete(path string, h HandleFunc, mddls ...Middleware) {
	r.Handle(http.MethodDelete, path, h, mddls...)
}

func (r *Router) Run(addr string) error {
	if r.server == nil {
		r.server = new(http.Server)
	}

	r.server.Addr = addr
	r.server.Handler = r.srvMux

	return r.server.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}

func (r *Router) Handle(method, path string, h HandleFunc, mddls ...Middleware) {
	path = constructPath(method, r.pfx, path)

	mddls = append(r.mddls, mddls...)
	for i := len(mddls) - 1; i > -1; i-- {
		mddl := mddls[i]
		h = mddl(h)
	}

	hh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		if err := h(ctx); err != nil {
			logger.Error(err.Error())
		}
	})

	r.srvMux.Handle(path, hh)
}

type HandleFunc func(ctx *Context) error

func constructPath(method, pref, path string) string {
	f := "%s %s%s"
	return fmt.Sprintf(f, method, pref, path)
}
