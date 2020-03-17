package api

import (
	"auto/authentication"
	"auto/cmd/server/config"

	"github.com/kataras/iris/v12"
)

func ServeLogin(app *iris.Application, cfg config.AutoConfig) {
	app.Get("/login", authentication.CreateHandlers(authentication.Config{
		Server: cfg.Server.Host,
		Gitlab: authentication.GitlabConfig{
			Host:         cfg.Gitlab.Host,
			ClientID:     cfg.Gitlab.ClientID,
			ClientSecret: cfg.Gitlab.ClientSecret,
			SkipVerify:   cfg.Gitlab.SkipVerify,
		},
	})...)
}
