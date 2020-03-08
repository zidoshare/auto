package routes

import (
	"auto/drone"
	"github.com/kataras/iris"
)

func Routes(app *iris.Application) {
	app.Get("/api/auto/drone/config", drone.Configuration)
}
