package login

import (
	"context"
	"net/http"
	"time"
)

type Middleware interface {
	Handler(h http.Handler) http.Handler
}

type Token struct {
	Access  string
	Refresh string
	Expires time.Time
}
type key int

const (
	tokenKey key = iota
	errorKey
)

func WithToken(parent context.Context, token *Token) context.Context {
	return context.WithValue(parent, tokenKey, token)
}

func WithError(parent context.Context, err error) context.Context {
	return context.WithValue(parent, errorKey, err)
}

func TokenFrom(ctx context.Context) *Token {
	token, _ := ctx.Value(tokenKey).(*Token)
	return token
}

func ErrorFrom(ctx context.Context) error {
	err, _ := ctx.Value(errorKey).(error)
	return err
}
