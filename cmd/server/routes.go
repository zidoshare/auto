package main

import (
	"auto/cmd/server/api"
	"auto/cmd/server/config"
	"auto/projects"

	"github.com/kataras/iris/v12"
)

func routes(app *iris.Application, cfg config.AutoConfig) {
	client := projects.GetClient(projects.Config{
		Gitlab: projects.GitlabConfig{
			Server:      cfg.Gitlab.Host,
			AccessToken: cfg.Gitlab.AccessToken,
		}},
	)
	api.ServeDrone(app, client, cfg)
	api.ServeLogin(app, cfg)
}
