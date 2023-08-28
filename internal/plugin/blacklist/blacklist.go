package blacklist

import (
	"os"

	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"gopkg.in/yaml.v3"
)

const configFilePath = "./config/blacklist.yaml"

type BlacklistConfig struct {
	List []int64 `yaml:"blacklist"`
}

func init() {
	confData, err := os.ReadFile(configFilePath)
	if err != nil {
		logrus.Errorln("[blacklist]", "Fail to read config file", err)
		return
	}
	var config BlacklistConfig
	if err := yaml.Unmarshal(confData, &config); err != nil {
		logrus.Errorln("[blacklist]", "Fail to unmarshal config data", err)
		return
	}
	zero.OnMessage(zero.CheckUser(config.List...)).SetPriority(1).SetBlock(true)
}
