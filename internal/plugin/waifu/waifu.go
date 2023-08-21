// This file is modified from ZeroBot-Plugin.
// View this link for more detail:
// https://github.com/FloatTech/ZeroBot-Plugin/tree/master/plugin/aiwife

package waifu

import (
	"fmt"
	"math/rand"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	bed = "https://www.thiswaifudoesnotexist.net/example-%d.jpg"
)

func init() {
	zero.OnCommand("waifu").Handle(func(ctx *zero.Ctx) {
		miku := rand.Intn(100000) + 1
		ctx.SendChain(message.At(ctx.Event.UserID), message.Image(fmt.Sprintf(bed, miku)))
	})
}
