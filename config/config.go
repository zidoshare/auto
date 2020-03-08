package config

import (
	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type AutoConfig struct {
	Server ServerConfig `toml:"server"`
	Drone  DroneConfig  `toml:"drone"`
	Gitlab GitlabConfig `toml:"gitlab"`
}

type ServerConfig struct {
	Port  string `toml:"port" short:"p" long:"server-port" description:"服务器绑定端口" value-name:"port"`
	Debug bool   `toml:"debug" short:"d" long:"server-debug" description:"是否开启debug模式" value-name:"debug"`
}

type DroneConfig struct {
	Secret string `toml:"secret" long:"drone-secret" description:"drone服务器api secret" value-name:"secret"`
	YmlDir string `toml:"yml_dir" long:"drone-yml-dir" description:"drone服务器所需要的默认配置文件所在位置" value-name:"dir"`
}

type GitlabConfig struct {
	Host        string   `toml:"host" long:"gitlab-host" description:"gitlab基本路径" value-name:"host"`
	AccessToken string   `toml:"access_token" long:"gitlab-access-token" description:"gitlab服务器accessToken" value-name:"<token>"`
	Namespace   []string `toml:"namespace" long:"gitlab-namespace" description:"gitlab服务器对应的namespace空间" value-name:"namespace"`
}

var (
	cfg  *AutoConfig
	once sync.Once
)

func getConfigFileFromExecutable(fileName string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return path.Join(dir, fileName)
}

func Config() *AutoConfig {
	once.Do(func() {
		var cmdConfig struct {
			ServerConfig
			DroneConfig
			GitlabConfig
			ConfigPath string `short:"c" long:"config" value-name:"[FILE]" description:"配置文件位置"`
		}
		parser := flags.NewParser(&cmdConfig, flags.Default|flags.IgnoreUnknown)
		if _, err := parser.Parse(); err != nil {
			if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}
		path := cmdConfig.ConfigPath
		if path == "" {
			path = os.Getenv("AUTO_CONF")
			if path == "" {
				path = getConfigFileFromExecutable("config.toml")
			}
		}
		info, err := os.Stat(path)
		cfg = &AutoConfig{
			Server: ServerConfig{
				Port: ":8002",
			},
		}
		if !os.IsNotExist(err) {
			if !info.IsDir() {
				if _, err := toml.DecodeFile(path, cfg); err != nil {
					logrus.Panic(err)
				}
			} else {
				logrus.Panic("config file is a directory")
			}
		}
		if err := mergo.Merge(&cfg.Server, cmdConfig.ServerConfig, mergo.WithOverride); err != nil {
			logrus.Panicf("cannot merge config server from command arguments to default config:%s", err)
		}

		if err := mergo.Merge(&cfg.Drone, cmdConfig.DroneConfig, mergo.WithOverride); err != nil {
			logrus.Panicf("cannot merge config drone from command arguments to default config:%s", err)
		}
		if err := mergo.Merge(&cfg.Gitlab, cmdConfig.GitlabConfig, mergo.WithOverride); err != nil {
			logrus.Panicf("cannot merge config gitlab from command arguments to default config:%s", err)
		}
	})
	return cfg
}
