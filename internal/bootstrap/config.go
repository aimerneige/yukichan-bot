package bootstrap

import (
	"fmt"
	"os"

	"github.com/aimerneige/yukichan-bot/internal/config"
	"gopkg.in/yaml.v3"
)

func InitConfig(confPath string) {
	confData, err := os.ReadFile(confPath)
	if err != nil {
		fmt.Printf("Fail to read config file in path \"%s\".\nFatal Error, Exiting...\n", confPath)
		os.Exit(1)
	}
	if err := yaml.Unmarshal(confData, &config.Conf); err != nil {
		fmt.Printf("Fail to unmarshal your config file. Please check your file in \"%s\".\nFatal Error, Exiting...\n", confPath)
		os.Exit(1)
	}
}
