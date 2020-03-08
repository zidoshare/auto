package drone

import (
	autoConfig "auto/config"
	"context"
	"errors"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"net/http"
)

type droneExPlugin struct {
}

//根据项目空间、项目文件结构获取相应的默认配置文件
func (p *droneExPlugin) Find(ctx context.Context, req *config.Request) (*drone.Config, error) {
	namespaces := autoConfig.Config().Gitlab.Namespace
	if namespaces == nil {
		logrus.Debug("未配置gitlab namespace")
		return nil, nil
	}
	for _, namespace := range namespaces {
		if req.Repo.Namespace == namespace {
			data, err := getYml(req.Repo.Slug, autoConfig.Config().Drone.YmlDir)
			return &drone.Config{
				Data: data,
			}, err
		}
	}
	return nil, nil
}

var handler http.Handler

func init() {
	handler = config.Handler(&droneExPlugin{}, autoConfig.Config().Drone.Secret, logrus.StandardLogger())
}

//Configuration extension
func Configuration(ctx iris.Context) {
	handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
