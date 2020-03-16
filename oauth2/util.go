package oauth2

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/kataras/iris/v12"
)

const cookieName = "_oauth_state"

func createState(ctx iris.Context) string {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  random(),
		MaxAge: 1800,
	}
	ctx.SetCookie(cookie)
	return cookie.Value
}

func validateState(ctx iris.Context, state string) error {
	cookie := ctx.GetCookie(cookieName)
	if state != cookie {
		return ErrState
	}
	return nil
}

func deleteState(ctx iris.Context) {
	ctx.RemoveCookie(cookieName)
}

func random() string {
	return fmt.Sprintf("%x", rand.Uint64())
}
