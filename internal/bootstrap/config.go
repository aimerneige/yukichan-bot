package bootstrap

import (
	"fmt"
	"os"

	"github.com/aimerneige/yukichan-bot/internal/config"
)

func InitConfig(configPath string) {
	if err := config.InitConfig(configPath); err != nil {
		fmt.Printf("Fail to read config file in path \"%s\".\nFatal Error, Exiting...\n", configPath)
		os.Exit(1)
	}
}
