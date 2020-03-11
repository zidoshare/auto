package gitlab

import (
	"auto/config"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

var client = &http.Client{}

type FileNode struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	WebUrl string `json:"webUrl"`
}

func FileNames(fullpath string) ([]string, error) {
	nodes, err := Files(fullpath)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, node := range nodes {
		result = append(result, node.Name)
	}
	return result, nil
}

//仅查询根目录下所有文件（不包括目录）
func Files(fullpath string) ([]FileNode, error) {
	var filesQL = `{"query":"query {
  project(fullPath:` + fullpath + `) {
    id
    name
    repository{
      tree{
        blobs{
          nodes{
            name
						path
						webUrl
          }
        }
      }
    }
  }
}"}`
	url := config.Config().Gitlab.Host + "/api/graphql"
	req, err := http.NewRequest("POST", url, strings.NewReader(filesQL))
	if err != nil {
		logrus.Errorf("cannot get files from %s:%s", url, err)
	}
	req.Header["Authorization"] = []string{"Bearer " + config.Config().Gitlab.AccessToken}
	req.Header["Content-Type"] = []string{"application/json"}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot get files from %s:%s", url, err))
	}
	nodes := make([]FileNode, 0)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot get files from %s:%s", url, err))
	}
	if err := json.Unmarshal(body, &nodes); err != nil {
		return nil, errors.New(fmt.Sprintf("cannot get files from %s:%s", url, err))
	}
	return nodes, nil
}
