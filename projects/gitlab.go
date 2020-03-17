package projects

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

//GitlabClient gitlab项目客户端实现
type GitlabClient struct {
	Server      string
	AccessToken string
	Client      *http.Client
}

//Files 仅查询根目录下所有文件（不包括目录）
func (c GitlabClient) Files(fullPath string) ([]FileNode, error) {
	var filesQL = `{"query":"query {
  project(fullPath:` + fullPath + `) {
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
	url := c.Server + "/api/graphql"
	req, err := http.NewRequest("POST", url, strings.NewReader(filesQL))
	if err != nil {
		logrus.Errorf("cannot get files from %s:%s", url, err)
	}
	req.Header["Authorization"] = []string{"Bearer " + c.AccessToken}
	req.Header["Content-Type"] = []string{"application/json"}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot get files from %s:%s", url, err)
	}
	nodes := make([]FileNode, 0)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot get files from %s:%s", url, err)
	}
	if err := json.Unmarshal(body, &nodes); err != nil {
		return nil, fmt.Errorf("cannot get files from %s:%s", url, err)
	}
	return nodes, nil
}
