package routes

import (
	"github.com/kataras/iris"
	"gitlab.scustartup.com/hnqc/auto/api/drone"
)

func Routes(app *iris.Application) {
	app.Get("/api/auto/drone/config", drone.Configuration)
}
