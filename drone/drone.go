package drone

import (
	autoConfig "auto/config"
	"auto/gitlab"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kataras/iris/v12"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

type droneExPlugin struct {
	gitlabConfig gitlab.ApiConfig
}

//getYml 根据项目和配置目录匹配相关的配置文件
func GetYmlByFileNames(fileNames []string, configDir string) (string, error) {
	configs := make(map[string]string)
	if err := filepath.Walk(configDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".drone.yml") {
			configs[strings.TrimSuffix(info.Name(), ".drone.yml")] = p
		}
		return nil
	}); err != nil {
		return "", errors.New(fmt.Sprintf("解析配置文件目录(*.drone.yml)发生错误:%s", err))
	}
	if autoConfig.Config().Server.Debug {
		var logs = "匹配配置标识文件:{"
		for k, v := range configs {
			logs += fmt.Sprintf("%s->%s,", k, v)
		}
		if len(configs) > 0 {
			logs = logs[:len(logs)-1]
		}
		logs += "}"
		logrus.Debug(logs)
	}
	resultConf := make([]string, 0)
	for _, name := range fileNames {
		if v, ok := configs[name]; ok {
			logrus.Debugf("匹配:%s -> %s", name, v)
			resultConf = append(resultConf, v)
		}
	}
	if l := len(resultConf); l != 1 {
		if l > 1 {
			return "", errors.New(fmt.Sprintf("仓库与配置标识文件匹配到多个配置文件:%v", resultConf))
		} else if l == 0 {
			return "", errors.New(fmt.Sprintf("仓库与配置标识文件未匹配到任何配置文件"))
		}
	}
	result, err := ioutil.ReadFile(resultConf[0])
	return string(result), err
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
			fileNames, err := p.gitlabConfig.FileNames(req.Repo.Slug)
			if err != nil {
				return nil, err
			}
			data, err := GetYmlByFileNames(fileNames, autoConfig.Config().Drone.YmlDir)
			return &drone.Config{
				Data: data,
			}, err
		}
	}
	return nil, nil
}

type Config struct {
	ymlDir          string
	callbackHandler http.Handler
}

func New(secret, ymlDir string, gitlabConfig gitlab.ApiConfig) *Config {
	plugin := &droneExPlugin{
		gitlabConfig,
	}
	return &Config{
		ymlDir:          ymlDir,
		callbackHandler: config.Handler(plugin, secret, logrus.StandardLogger()),
	}
}

func (c *Config) Yml(ctx iris.Context) {
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
	data, err := GetYmlByFileNames(fileNames, c.ymlDir)
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

func (c *Config) Callback(ctx iris.Context) {
	c.callbackHandler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
