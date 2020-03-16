package config

import (
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type AutoConfig struct {
	Server ServerConfig `toml:"server"`
	Drone  DroneConfig  `toml:"drone"`
	Gitlab GitlabConfig `toml:"gitlab"`
}

type ServerConfig struct {
	Host  string `toml:"host" default:"localhost:8002"`
	Proto string `toml:"proto" default:"http"`
	Port  string `toml:"port" default:":8002"`
	Debug bool   `toml:"debug"`
	Addr  string `toml:"-"`
}

type DroneConfig struct {
	Secret string `toml:"secret"`
	YmlDir string `toml:"yml_dir"`
}

type GitlabConfig struct {
	Host         string   `toml:"host"`
	ClientID     string   `toml:"client_id"`
	ClientSecret string   `toml:"client_secret"`
	SkipVerify   bool     `toml:"skip_verify"`
	AccessToken  string   `toml:"access_token"`
	Namespace    []string `toml:"namespace"`
}

var (
	cfg  *AutoConfig
	once sync.Once
)

func parseConfig() {
	path := os.Getenv("AUTO_CONF")
	if path == "" {
		path = "config.toml"
	}

	info, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if !info.IsDir() {
			if _, err := toml.DecodeFile(path, cfg); err != nil {
				logrus.Panic(err)
			}
		} else {
			logrus.Panicf("config file is a directory:%s", path)
		}
	} else {
		logrus.Panicf("config file is not exists:%s", path)
	}
	defaultAddr(cfg)
}

func Config() *AutoConfig {
	once.Do(parseConfig)
	return cfg
}

func defaultAddr(c *AutoConfig) {
	c.Server.Addr = c.Server.Proto + "://" + c.Server.Host
}
