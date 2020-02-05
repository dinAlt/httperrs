package httperrs

import (
	"context"
	"net/http"
)

type key string

const (
	keyError = key("error")
)

func SetError(r *http.Request, err error) *http.Request {
	return r.WithContext(
		context.WithValue(r.Context(), keyError, err))
}

func GetError(r *http.Request) error {
	err, _ := r.Context().Value(keyError).(error)

	if err == nil {
		return nil
	}

	return err
}
