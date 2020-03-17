package core

import (
	"time"

	"github.com/kataras/iris/v12"
)

const (
	tokenKey = "KEY_TOKEN"
	errorKey = "KEY_ERROR"
)

func WithToken(ctx iris.Context, token *Token) {
	ctx.Values().Set(tokenKey, token)
}

func WithError(ctx iris.Context, err error) {
	ctx.Values().Set(errorKey, err)
}

func TokenFrom(ctx iris.Context) *Token {
	token, _ := ctx.Values().Get(tokenKey).(*Token)
	return token
}

func ErrorFrom(ctx iris.Context) error {
	err, _ := ctx.Values().Get(errorKey).(error)
	return err
}

type Token struct {
	Access  string
	Refresh string
	Expires time.Time
}
