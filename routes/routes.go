package routes

import (
	"auto/serv"

	"github.com/kataras/iris/v12"
)

func Routes(app *iris.Application) {
	//drone回调地址
	app.Get("/api/auto/drone/callback/config", serv.DroneYmlCallback)
	//客户端需要的template地址
	app.Get("/api/auto/drone/template", serv.Yml)
}
