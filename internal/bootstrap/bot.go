package bootstrap

import (
	"github.com/aimerneige/yukichan-bot/internal/config"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"

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

func StartBot() {
	common := config.ReadCommonConfig()
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
