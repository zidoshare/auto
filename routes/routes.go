package routes

import (
	"auto/drone"

	"github.com/kataras/iris/v12"
)

func Routes(app *iris.Application) {
	app.Get("/api/auto/drone/config", drone.Configuration)
}
