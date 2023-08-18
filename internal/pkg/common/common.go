package common

import (
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/extension/single"
)

var DefaultSingle = single.New(
	single.WithKeyFn(func(ctx *zero.Ctx) int64 {
		return ctx.Event.UserID
	}),
	single.WithPostFn[int64](func(ctx *zero.Ctx) {
		ctx.Send("您有操作正在执行，请稍后再试！")
	}),
)

var DefaultSpeedLimit = func(ctx *zero.Ctx) bool {
	if !rate.NewManager[int64](time.Minute*1, 1).Load(ctx.Event.UserID).Acquire() {
		ctx.Send("您的请求太快，请稍后再试！")
		return false
	}
	return true
}
