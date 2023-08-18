package match

import (
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	zero.OnFullMatch("老婆").Handle(func(ctx *zero.Ctx) {
		ctx.Send("肥宅不要乱叫老婆啊！")
	})
}
