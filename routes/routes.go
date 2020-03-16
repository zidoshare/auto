package routes

import (
	"auto/config"
	"auto/drone"
	"auto/gitlab"
	"net/http"

	"github.com/kataras/iris/v12"
)

func Routes(app *iris.Application) {
	gitlabConfig := &gitlab.ApiConfig{
		Server: config.Config().Gitlab.Host,
		Client: &http.Client{},
	}
	droneConf := drone.New(config.Config().Drone.Secret, config.Config().Drone.YmlDir, *gitlabConfig)
	//drone回调地址
	app.Get("/drone/callback/config", droneConf.Callback)
	//客户端需要的template地址
	app.Get("/drone/template", droneConf.Yml)
	app.Get("/login", createLoginHandlers()...)
}
