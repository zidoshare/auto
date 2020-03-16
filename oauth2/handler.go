package oauth2

import (
	"auto/login"
	"errors"
	"time"

	"github.com/kataras/iris/v12"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	Conf *Config
}

func (h *Handler) Handle(ctx iris.Context) {
	if err := ctx.FormValue("error"); err != "" {
		logrus.Errorf("oauth:认证错误: %s", err)
		login.WithError(ctx, errors.New(err))
		ctx.Next()
		return
	}

	code := ctx.FormValue("code")
	if len(code) == 0 {
		state := createState(ctx)
		ctx.Redirect(h.Conf.authorizeRedirect(state), 303)
		return
	}
	state := ctx.FormValue("state")
	deleteState(ctx)
	if err := validateState(ctx, state); err != nil {
		logrus.Errorln("oauth: state缺少或已经失效")
		login.WithError(ctx, err)
		ctx.Next()
		return
	}
	source, err := h.Conf.exchange(code, state)
	if err != nil {
		logrus.Errorf("oauth: 无法交换code: %s: %s", code, err)
		login.WithError(ctx, err)
		ctx.Next()
		return
	}

	login.WithToken(ctx, &login.Token{
		Access:  source.AccessToken,
		Refresh: source.RefreshToken,
		Expires: time.Now().UTC().Add(
			time.Duration(source.Expires) * time.Second,
		),
	})

	ctx.Next()
}
