package test_server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/ryanfaerman/dispatch"
)

type testServer struct {
	*httptest.Server
	t *testing.T
}

func (s *testServer) Request(method, path string, payload io.Reader) (*http.Request, *http.Response, []byte) {
	req, err := http.NewRequest(method, s.URL+path, payload)
	if err != nil {
		s.t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.t.Fatal(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return req, res, body
}

func New(t *testing.T, options ...func(*dispatch.Dispatch)) *testServer {
	r := dispatch.New()

	for _, option := range options {
		option(r)
	}

	s := &testServer{httptest.NewServer(r), t}
	runtime.SetFinalizer(s, func(s *testServer) { s.Server.Close() })
	return s
}
