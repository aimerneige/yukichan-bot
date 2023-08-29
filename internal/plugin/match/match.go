package match

import (
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	engine := zero.New()
	engine.OnFullMatch("老婆").
		SetPriority(3).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.Send("肥宅不要乱叫老婆啊！")
		})
	engine.OnFullMatchGroup([]string{"关于", "about"}).
		SetPriority(3).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.Send("ゆき酱国际象棋机器人\n本项目是使用 AGPL-3.0 开源协议授权的开源项目。\n开源地址及使用帮助：https://github.com/aimerneige/yukichan-bot\n捐赠支持开发：https://aimer.aiursoft.cn/zh/donate/")
		})
}
