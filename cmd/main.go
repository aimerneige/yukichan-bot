package main

import (
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"

	_ "github.com/aimerneige/yukichan-bot/internal/plugin/fadian"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/github"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/match"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/read60s"
	_ "github.com/aimerneige/yukichan-bot/internal/plugin/waifu"
)

// Path variables are overrides during build
// See the Makefile for more information.
var (
	ConfPath string = "./config/config.yaml"
)

var (
	WSServer      = "ws://127.0.0.1:6701"
	WSToken       = ""
	NickName      = []string{"ゆき酱"}
	CommandPrefix = "/"
	SuperUsers    = []int64{1227427929}
)

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	zero.RunAndBlock(&zero.Config{
		NickName:      NickName,
		CommandPrefix: CommandPrefix,
		SuperUsers:    SuperUsers,
		Driver: []zero.Driver{
			// 反向 WS
			driver.NewWebSocketServer(16, WSServer, WSToken),
		},
	}, nil)
}
