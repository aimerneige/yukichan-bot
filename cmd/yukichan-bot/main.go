package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"gopkg.in/yaml.v3"

	"github.com/aimerneige/yukichan-bot/internal/config"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/alipay"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/chess"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/donate"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/fadian"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/fortune"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/github"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/manager"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/match"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/music"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/random"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/read60s"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/setu"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/suangua"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/tarot"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/waifu"
)

// Variables are overrides during build
// See the Makefile for more information.
var (
	ConfPath   string = "./config/application.yaml"
	DebugLevel string = "info"
)

func init() {
	initLog()
	initConfig()
}

func main() {
	common := config.Conf.Common
	zero.RunAndBlock(&zero.Config{
		NickName:      common.NickName,
		CommandPrefix: common.CommandPrefix,
		SuperUsers:    common.SuperUsers,
		Driver: []zero.Driver{
			// 反向 WS
			driver.NewWebSocketServer(16, common.WSServer, common.WSToken),
		},
	}, nil)
}

func initLog() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
	})
	if debugLevel, err := log.ParseLevel(DebugLevel); err != nil {
		fmt.Printf("Fail to parse debug level \"%s\".\nFatal Error, Exiting...\n", DebugLevel)
		os.Exit(1)
	} else {
		log.SetLevel(debugLevel)
	}
}

func initConfig() {
	confData, err := os.ReadFile(ConfPath)
	if err != nil {
		fmt.Printf("Fail to read config file in path \"%s\".\nFatal Error, Exiting...\n", ConfPath)
		os.Exit(1)
	}
	if err := yaml.Unmarshal(confData, &config.Conf); err != nil {
		fmt.Printf("Fail to unmarshal your config file. Please check your file in \"%s\".\nFatal Error, Exiting...\n", ConfPath)
		os.Exit(1)
	}
}
