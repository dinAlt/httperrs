package httperrs

import (
	"context"
	"net/http"
)

type key string

const (
	keyOnError = key("onerror")
)

type ErrorHandler func(*HTTPError)

type HTTPError struct {
	HTTPRequest *http.Request
	Inner       error
}

func (e *HTTPError) Error() string {
	return e.Inner.Error() +
		" [host=" + e.HTTPRequest.URL.Hostname() + "] " +
		"[path=" + e.HTTPRequest.URL.Path + "] " +
		"[query=" + e.HTTPRequest.URL.RawQuery + "] " +
		"[method=" + e.HTTPRequest.Method + "] "
}


type Middleware struct {
	ErrorHandler
	Next http.Handler
}

func Push(r *http.Request, err error) {
	onErr, _ := r.Context().Value(keyOnError).(ErrorHandler)
	if onErr == nil {
		panic("httperrs: no error handler in context")
	}

	onErr(&HTTPError{HTTPRequest: r, Inner: err})
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Next.ServeHTTP(w,
		r.WithContext(
			context.WithValue(
				r.Context(),
				keyOnError,
				m.ErrorHandler),
		),
	)
}
