package dispatch

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nbio/httpcontext"
)

type DispatchRequest struct{ Request *http.Request }

func Request(req *http.Request) *DispatchRequest {
	return &DispatchRequest{Request: req}
}

func (dc *DispatchRequest) Params() httprouter.Params {
	if value, ok := dc.GetOk(paramsKey); ok {
		if params, ok := value.(httprouter.Params); ok {
			return params
		}
	}
	return httprouter.Params{}
}

func (dc *DispatchRequest) Set(key interface{}, value interface{}) {
	httpcontext.Set(dc.Request, key, value)
}

func (dc *DispatchRequest) Get(key interface{}) (val interface{}) {
	return httpcontext.Get(dc.Request, key)
}

func (dc *DispatchRequest) GetOk(key interface{}) (val interface{}, ok bool) {
	return httpcontext.GetOk(dc.Request, key)
}

// Params is a shortcut method for when you don't want the entire upgraded request
func Params(req *http.Request) httprouter.Params {
	return Request(req).Params()
}

func (d *DispatchRequest) ContentType() string {
	content := d.Request.Header.Get("Content-Type")
	for i, a := range content {
		if a == ' ' || a == ';' {
			return content[:i]
		}
	}
	return content
}
