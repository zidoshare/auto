package drone

import (
	"context"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"net/http"
)

const defaultPipeline = `
`

type droneExPlugin struct {
}

func (p *droneExPlugin) Find(ctx context.Context, req *config.Request) (*drone.Config, error) {
	//根据项目空间、项目文件结构获取相应的默认配置文件
	if req.Repo.Namespace == "hnqc" {
		return &drone.Config{
			Data: defaultPipeline,
		}, nil
	}
	return nil, nil
}

type Spec struct {
	Secret string `envconfig:"DRONE_SECRET"`
}

var spec = new(Spec)
var handler http.Handler

func init() {
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	handler = config.Handler(&droneExPlugin{}, spec.Secret, logrus.StandardLogger())
}

//Configuration extension
func Configuration(ctx iris.Context) {
	handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
