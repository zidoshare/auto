package gitlab

import (
	"auto/oauth2"
	"net/http"
	"strings"

	"github.com/kataras/iris/v12/context"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Server       string
	Scope        []string
	Client       *http.Client
}

func CreateMiddleware(c *Config) context.Handler {
	server := normalizeAddress(c.Server)
	handler := oauth2.Handler{
		Conf: &oauth2.Config{
			BasicAuthOff:     true,
			Client:           c.Client,
			ClientID:         c.ClientID,
			ClientSecret:     c.ClientSecret,
			RedirectURL:      c.RedirectURL,
			AccessTokenURL:   server + "/oauth/token",
			AuthorizationURL: server + "/oauth/authorize",
			Scope:            c.Scope,
		}}
	return handler.Handle
}

func normalizeAddress(address string) string {
	if address == "" {
		return "https://gitlab.com"
	}
	return strings.TrimSuffix(address, "/")
}
