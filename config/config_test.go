package config

import (
	"os"
	"reflect"
	"testing"
)

func Test_config(t *testing.T) {
	tests := []struct {
		name   string
		before func()
		want   *AutoConfig
	}{
		{name: "test all config",
			before: func() {
				if err := os.Setenv("AUTO_CONF", "./not_exists.toml"); err != nil {
					t.Error(err)
				}
			}, want: &AutoConfig{
				Server: ServerConfig{
					Listen: ":8002",
					Debug:  false,
				},
				Drone: DroneConfig{
					Secret: "",
					YmlDir: "",
				},
				Gitlab: GitlabConfig{
					Host:        "",
					AccessToken: "",
					Namespace:   nil,
				},
			}},
		{
			name: "test with config file",
			before: func() {
				//use default config toml
				if err := os.Setenv("AUTO_CONF", ""); err != nil {
					t.Error(err)
				}

			},
			want: &AutoConfig{
				Server: ServerConfig{
					Listen: ":8003",
					Debug:  true,
				},
				Drone: DroneConfig{
					Secret: "xxx",
					YmlDir: "dir",
				},
				Gitlab: GitlabConfig{
					Host:        "http://gitlab.example.com",
					AccessToken: "token",
					Namespace:   nil,
				},
			},
		},
		{
			name: "test with arguments",
			before: func() {
				//use args
				os.Args = []string{"auto-server", "--server-listen", ":8005", "--drone-secret", "xx2"}
			},
			want: &AutoConfig{
				Server: ServerConfig{
					Listen: ":8005",
					Debug:  true,
				},
				Drone: DroneConfig{
					Secret: "xx2",
					YmlDir: "dir",
				},
				Gitlab: GitlabConfig{
					Host:        "http://gitlab.example.com",
					AccessToken: "token",
					Namespace:   nil,
				},
			},
		},
		{
			name: "use -c",
			before: func() {
				//use args
				os.Args = []string{"auto-server", "--server-listen", ":8005", "--drone-secret", "xx2", "-c", "config-2.toml"}
			},
			want: &AutoConfig{
				Server: ServerConfig{
					Listen: ":8005",
					Debug:  true,
				},
				Drone: DroneConfig{
					Secret: "xx2",
					YmlDir: "dir",
				},
				Gitlab: GitlabConfig{
					Host:        "http://gitlab.example.com",
					AccessToken: "token",
					Namespace:   nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			parseConfig()
			if got := Config(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
