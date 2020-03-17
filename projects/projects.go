//Package projects 工程项目文件相关
package projects

import (
	"net/http"
)

//ProjectClient projects client抽象
type Client interface {
	Files(fullPath string) ([]FileNode, error)
}

//FileNode 项目文件抽象
type FileNode struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	WEBUrl string `json:"webUrl"`
}

type GitlabConfig struct {
	Server      string
	AccessToken string
}

type Config struct {
	Gitlab GitlabConfig
}

//GetClient 获取projects client以帮助处理项目相关api
func GetClient(cfg Config) Client {
	if cfg.Gitlab.Server != "" {
		return GitlabClient{
			Server:      cfg.Gitlab.Server,
			AccessToken: cfg.Gitlab.AccessToken,
			Client:      &http.Client{},
		}
	}
	return nil
}

//FileNames 从路径中获取所有文件名集合
func FileNames(c Client, fullPath string) ([]string, error) {
	nodes, err := c.Files(fullPath)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, node := range nodes {
		result = append(result, node.Name)
	}
	return result, nil
}
