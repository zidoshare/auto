//Package authentication 认证相关包
package authentication

import (
	"auto/authentication/gitlab"
	"crypto/tls"
	"net/http"

	"github.com/kataras/iris/v12/context"
)

type Config struct {
	Server string
	Gitlab GitlabConfig
}

type GitlabConfig struct {
	Host         string
	ClientID     string
	ClientSecret string
	SkipVerify   bool
}

func CreateHandlers(cfg Config) context.Handlers {
	var handlers context.Handlers
	if cfg.Gitlab.ClientID == "" {
		c := &gitlab.GitlabClient{
			ClientID:     cfg.Gitlab.ClientID,
			ClientSecret: cfg.Gitlab.ClientSecret,
			RedirectURL:  cfg.Server + "/login",
			Server:       cfg.Gitlab.Host,
			Client:       defaultClient(cfg.Gitlab.SkipVerify),
		}
		handlers = append(handlers, gitlab.CreateMiddleware(c))
	}

	return handlers
}

func defaultClient(skipVerify bool) *http.Client {
	client := &http.Client{}
	client.Transport = defaultTransport(skipVerify)
	return client
}

func defaultTransport(skipverify bool) http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipverify,
		},
	}
}
