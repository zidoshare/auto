package serv

import (
	"auto/config"
	"auto/drone"

	"github.com/kataras/iris/v12"
)

func Yml(ctx iris.Context) {
	var fileNames []string
	err := ctx.ReadJSON(&fileNames)
	if err != nil {
		ctx.StatusCode(400)
	}
	if len(fileNames) == 0 {
		ctx.JSON(map[string]interface{}{
			"code": 0,
			"data": "",
		})
		return
	}
	data, err := drone.GetYmlByFileNames(fileNames, config.Config().Drone.YmlDir)
	if err != nil {
		ctx.StatusCode(500)
		ctx.JSON(map[string]interface{}{
			"code": 100,
			"msg":  err.Error(),
		})
		return
	}
	ctx.StatusCode(200)
	ctx.Text(data)
}

//DroneYmlCallback extension
func DroneYmlCallback(ctx iris.Context) {
	drone.GetDroneCallbackHandler().ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
