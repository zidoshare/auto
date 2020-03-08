package drone

import (
	"auto/gitlab"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func loadConfigFilePaths(configDir string) (map[string]string, error) {
	configs := make(map[string]string)
	if err := filepath.Walk(configDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".drone.yml") {
			configs[strings.TrimRight(info.Name(), ".drone.yml")] = p
		}
		return nil
	}); err != nil {
		return nil, errors.New(fmt.Sprintf("cannot parses default config files(*.drone.yml): %s", err))
	}
	return configs, nil
}

func getYml(project, configDir string) (string, error) {
	gitNodes, err := gitlab.Files(project)
	if err != nil {
		logrus.Error(err)
	}
	configs, err := loadConfigFilePaths(configDir)
	if err != nil {
		return "", err
	}
	var resultConf string
	for _, node := range gitNodes {
		if configs[node.Name] != "" {
			if resultConf != "" {
				return "", errors.New("仓库与配置标识文件匹配到多个")
			}
			resultConf = configs[node.Name]
		}
	}
	result, err := ioutil.ReadFile(resultConf)
	return string(result), err
}
