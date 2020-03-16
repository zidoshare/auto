package login

import (
	"time"

	"github.com/kataras/iris/v12"
)

//type Middleware interface {
//	Handler(h http.Handler) http.Handler
//}

type Token struct {
	Access  string
	Refresh string
	Expires time.Time
}
type key int

const (
	tokenKey = "KEY_TOKEN"
	errorKey = "KEY_ERROR"
)

//TODO handle login
func HandleLogin(ctx iris.Context) {

}

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
