package mux

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w: w,
		r: r,
	}
}

func (c *Context) BindJSON(dest any) error {
	err := json.NewDecoder(c.r.Body).Decode(dest)
	return err
}

func (c *Context) AddHeader(key, val string) {
	c.r.Header.Add(key, val)
}

func (c *Context) SetHeader(key, val string) {
	c.r.Header.Set(key, val)
}

func (c *Context) GetHeader(key string) string {
	return c.r.Header.Get(key)
}

func (c *Context) DeleteHeader(key string) {
	c.r.Header.Del(key)
}

func (c *Context) SetResponseHeader(header http.Header) {
	c.r.Header = header
}

func (c *Context) WriteJSON(status int, data any) error {
	c.w.Header().Set("Content-type", "application/json")
	c.w.WriteHeader(status)
	err := json.NewEncoder(c.w).Encode(data)

	return err
}

func (c *Context) Write(status int, data []byte) error {
	c.w.WriteHeader(status)
	_, err := c.w.Write(data)

	return err
}

func (c *Context) RequestContext() context.Context {
	return c.r.Context()
}

func (c *Context) SetRequestContext(ctx context.Context) {
	c.r = c.r.Clone(ctx)
}

func (c *Context) Path() string {
	return c.r.URL.Path
}

func (c *Context) PathVal(key string) string {
	return c.r.PathValue(key)
}

func (c *Context) QueryVal(key string) string {
	return c.r.URL.Query().Get(key)
}

func (c *Context) Query() url.Values {
	return c.r.URL.Query()
}

func (c *Context) Method() string {
	return c.r.Method
}

func (c *Context) FormVal(key string) string {
	if c.r.Form == nil {
		c.r.ParseForm()
	}

	return c.r.FormValue(key)
}
