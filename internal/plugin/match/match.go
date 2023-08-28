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
			ctx.Send("新一代机器人\n🚧施工中🚧\n开源地址：https://github.com/aimerneige/yukichan-bot\n捐赠支持开发：https://aimer.aiursoft.cn/zh/donate/")
		})
}
