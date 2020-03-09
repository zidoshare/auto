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

	"sync"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

type droneExPlugin struct {
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
	var resultConf []string
	for _, name := range fileNames {
		if v, ok := configs[name]; ok {
			logrus.Debugf("匹配:%s -> %s", name, v)
			resultConf = append(resultConf, v)
		}
	}
	if len(resultConf) > 1 {
		return "", errors.New(fmt.Sprintf("仓库与配置标识文件匹配到多个配置文件:%v", resultConf))
	} else if len(resultConf) < 1 {
		return "", errors.New(fmt.Sprintf("仓库与配置标识文件未匹配到任何配置文件"))
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
			fileNames, err := gitlab.FileNames(req.Repo.Slug)
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

var (
	handler http.Handler
	once    sync.Once
)

func GetDroneCallbackHandler() http.Handler {
	once.Do(func() {
		handler = config.Handler(&droneExPlugin{}, autoConfig.Config().Drone.Secret, logrus.StandardLogger())
	})
	return handler
}
