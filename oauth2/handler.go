package oauth2

import (
	"auto/login"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Handler(h http.Handler, c *Config) http.Handler {
	return &handler{next: h, conf: c}
}

type handler struct {
	conf *Config
	next http.Handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.FormValue("error"); err != "" {
		logrus.Errorf("oauth:认证错误: %s", err)
		ctx = login.WithError(ctx, errors.New(err))
		h.next.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	code := r.FormValue("code")
	if len(code) == 0 {
		state := createState(w)
		http.Redirect(w, r, h.conf.authorizeRedirect(state), 303)
		return
	}
	state := r.FormValue("state")
	deleteState(w)
	if err := validateState(r, state); err != nil {
		logrus.Errorln("oauth: state缺少或已经失效")
		ctx = login.WithError(ctx, err)
		h.next.ServeHTTP(w, r.WithContext(ctx))
		return
	}
	source, err := h.conf.exchange(code, state)
	if err != nil {
		logrus.Errorf("oauth: 无法交换code: %s: %s", code, err)
		ctx = login.WithError(ctx, err)
		h.next.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	ctx = login.WithToken(ctx, &login.Token{

		Access:  source.AccessToken,
		Refresh: source.RefreshToken,
		Expires: time.Now().UTC().Add(
			time.Duration(source.Expires) * time.Second,
		),
	})

	h.next.ServeHTTP(w, r.WithContext(ctx))
}
