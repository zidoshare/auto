package routes

import (
	"auto/config"
	"auto/gitlab"
	"auto/login"
	"crypto/tls"
	"net/http"

	"github.com/kataras/iris/v12/context"
)

func createLoginHandlers() []context.Handler {
	return []context.Handler{
		provideGitlabLogin(),
		login.HandleLogin,
	}
}

func provideGitlabLogin() context.Handler {
	cfg := config.Config()
	if cfg.Gitlab.ClientID == "" {
		return nil
	}
	return gitlab.CreateMiddleware(&gitlab.Config{
		ClientID:     cfg.Gitlab.ClientID,
		ClientSecret: cfg.Gitlab.ClientSecret,
		RedirectURL:  cfg.Server.Addr + "/login",
		Server:       cfg.Gitlab.Host,
		Client:       defaultClient(cfg.Gitlab.SkipVerify),
	})
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
