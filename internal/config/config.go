package config

import (
	"github.com/spf13/viper"
)

var GlobalConfig *Config

type Config struct {
	*viper.Viper
}

type CommonConfig struct {
	WSServer      string   `yaml:"ws_server"`
	WSToken       string   `yaml:"ws_token"`
	NickName      []string `yaml:"nickname"`
	CommandPrefix string   `yaml:"command_prefix"`
	SuperUsers    []int64  `yaml:"super_users"`
}

func InitConfig(configFilePath string) error {
	GlobalConfig = &Config{
		viper.New(),
	}
	GlobalConfig.SetConfigFile(configFilePath)
	return GlobalConfig.ReadInConfig()
}

func ReadCommonConfig() *CommonConfig {
	c := GlobalConfig.AllSettings()["common"]
	if common, ok := c.(CommonConfig); !ok {
		return GetDefaultCommonConfig()
	} else {
		return &common
	}
}

func GetDefaultCommonConfig() *CommonConfig {
	return &CommonConfig{
		WSServer:      "ws://127.0.0.1:6701",
		WSToken:       "",
		NickName:      []string{"ゆき酱"},
		CommandPrefix: "/",
		SuperUsers:    []int64{1227427929},
	}
}
