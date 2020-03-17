package api

import (
	config2 "auto/cmd/server/config"
	"auto/drone"
	"auto/projects"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

//ServeDrone drone
func ServeDrone(app *iris.Application, projectClient projects.Client, cfg config2.AutoConfig) {
	client := drone.NewClient(drone.Config{
		PermittedNameSpaces: cfg.Gitlab.Namespace,
		ConfigDir:           cfg.Drone.YmlDir,
		ProjectClient:       projectClient,
		GitlabSecret:        cfg.Gitlab.ClientSecret,
		Debug:               cfg.Server.Debug,
	})
	//drone回调地址
	app.Get("/drone/callback/config", func(ctx context.Context) {
		client.Callback.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	})
	//客户端需要的template地址
	app.Get("/drone/templates", func(ctx iris.Context) {
		var fileNames []string
		err := ctx.ReadJSON(&fileNames)
		if err != nil {
			ctx.StatusCode(400)
		}
		if len(fileNames) == 0 {
			_, _ = ctx.JSON(map[string]interface{}{
				"code": 0,
				"data": "",
			})
			return
		}
		data, err := client.GetYmlByFileNames(fileNames)
		if err != nil {
			ctx.StatusCode(500)
			_, _ = ctx.JSON(map[string]interface{}{
				"code": 100,
				"msg":  err.Error(),
			})
			return
		}
		ctx.StatusCode(200)
		_, _ = ctx.Text(data)
	})
}
