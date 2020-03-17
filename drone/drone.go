package drone

import (
	"auto/projects"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

type Config struct {
	PermittedNameSpaces []string
	ConfigDir           string
	ProjectClient       projects.Client
	GitlabSecret        string
	Debug               bool
}

type Client struct {
	PermittedNameSpaces []string
	ConfigDir           string
	Debug               bool
	Callback            http.Handler
}

func NewClient(cfg Config) Client {
	client := Client{
		PermittedNameSpaces: cfg.PermittedNameSpaces,
		ConfigDir:           cfg.ConfigDir,
		Debug:               cfg.Debug,
	}
	plugin := &droneExPlugin{
		droneClient:         client,
		permittedNameSpaces: cfg.PermittedNameSpaces,
		client:              cfg.ProjectClient,
	}
	client.Callback = config.Handler(plugin, cfg.GitlabSecret, logrus.StandardLogger())
	return client
}

//GetYmlByFileNames 根据项目和配置目录匹配相关的配置文件
func (c Client) GetYmlByFileNames(fileNames []string) (string, error) {
	configs := make(map[string]string)
	if err := filepath.Walk(c.ConfigDir, func(p string, info os.FileInfo, err error) error {
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
	if c.Debug {
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

type droneExPlugin struct {
	droneClient         Client
	permittedNameSpaces []string
	client              projects.Client
}

//Find 根据项目空间、项目文件结构获取相应的默认配置文件
func (p *droneExPlugin) Find(_ context.Context, req *config.Request) (*drone.Config, error) {
	if p.permittedNameSpaces == nil {
		logrus.Debug("未配置gitlab namespace")
		return nil, nil
	}
	for _, namespace := range p.permittedNameSpaces {
		if req.Repo.Namespace == namespace {
			fileNames, err := projects.FileNames(p.client, req.Repo.Slug)
			if err != nil {
				return nil, err
			}
			data, err := p.droneClient.GetYmlByFileNames(fileNames)
			return &drone.Config{
				Data: data,
			}, err
		}
	}
	return nil, nil
}
