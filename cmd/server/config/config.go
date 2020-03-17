package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

//AutoConfig 配置聚合
type AutoConfig struct {
	Server ServerConfig `toml:"server"`
	Drone  DroneConfig  `toml:"drone"`
	Gitlab GitlabConfig `toml:"gitlab"`
}

//ServerConfig 服务器相关配置
type ServerConfig struct {
	Host  string `toml:"host" default:"localhost:8002"`
	Proto string `toml:"proto" default:"http"`
	Port  string `toml:"port" default:":8002"`
	Debug bool   `toml:"debug"`
	Addr  string `toml:"-"`
}

//DroneConfig drone相关配置
type DroneConfig struct {
	Secret string `toml:"secret"`
	YmlDir string `toml:"yml_dir"`
}

//GitlabConfig gitlab相关配置
type GitlabConfig struct {
	Host         string   `toml:"host"`
	ClientID     string   `toml:"client_id"`
	ClientSecret string   `toml:"client_secret"`
	SkipVerify   bool     `toml:"skip_verify"`
	AccessToken  string   `toml:"access_token"`
	Namespace    []string `toml:"namespace"`
}

//Get AutoConfig struct
func Get() (AutoConfig, error) {
	cfg := AutoConfig{}
	path := os.Getenv("AUTO_CONF")
	if path == "" {
		path = "config.toml"
	}

	info, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if !info.IsDir() {
			if _, err := toml.DecodeFile(path, &cfg); err != nil {
				return cfg, err
			}
		} else {
			logrus.Panicf("config file is a directory:%s", path)
		}
	} else {
		logrus.Panicf("config file is not exists:%s", path)
	}
	defaultAddr(cfg)
	return cfg, nil
}

//defaultAddr
func defaultAddr(c AutoConfig) {
	c.Server.Addr = c.Server.Proto + "://" + c.Server.Host
}
