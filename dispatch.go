package dispatch

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nbio/httpcontext"
)

const (
	paramsKey = iota
)

type Dispatch struct {
	*httprouter.Router

	middleware []func(http.Handler) http.HandlerFunc
	pathPrefix string
}

func New() *Dispatch {
	return &Dispatch{Router: httprouter.New()}
}

func (d *Dispatch) Handle(method, path string, handler http.HandlerFunc) {
	d.Router.Handle(method, d.pathPrefix+path, d.Wrap(handler))
}

func (d *Dispatch) Wrap(handler http.Handler) httprouter.Handle {
	for _, middleware := range d.middleware {
		handler = middleware(handler)
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		httpcontext.Set(r, paramsKey, p)
		handler.ServeHTTP(w, r)
	}
}

func (d *Dispatch) Group(path string, middleware ...func(http.Handler) http.HandlerFunc) *Dispatch {
	return &Dispatch{
		Router:     d.Router,
		pathPrefix: path,
		middleware: append(d.middleware, middleware...),
	}
}

func (d *Dispatch) Use(middleware ...func(http.Handler) http.HandlerFunc) {
	d.middleware = append(d.middleware, middleware...)
}

func (d *Dispatch) GET(path string, handler http.HandlerFunc)     { d.Handle("GET", path, handler) }
func (d *Dispatch) PUT(path string, handler http.HandlerFunc)     { d.Handle("PUT", path, handler) }
func (d *Dispatch) POST(path string, handler http.HandlerFunc)    { d.Handle("POST", path, handler) }
func (d *Dispatch) PATCH(path string, handler http.HandlerFunc)   { d.Handle("PATCH", path, handler) }
func (d *Dispatch) DELETE(path string, handler http.HandlerFunc)  { d.Handle("DELETE", path, handler) }
func (d *Dispatch) OPTIONS(path string, handler http.HandlerFunc) { d.Handle("OPTIONS", path, handler) }
