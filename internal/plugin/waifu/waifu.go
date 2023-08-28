// This file is modified from ZeroBot-Plugin.
// View this link for more detail:
// https://github.com/FloatTech/ZeroBot-Plugin/tree/master/plugin/aiwife

package waifu

import (
	"encoding/base64"
	"fmt"
	"math/rand"

	"github.com/FloatTech/floatbox/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	bed = "https://www.thiswaifudoesnotexist.net/example-%d.jpg"
)

func init() {
	zero.OnCommand("waifu").
		SetPriority(8).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if data, err := web.GetData(fmt.Sprintf(bed, rand.Intn(100000)+1)); err == nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Image("base64://"+base64.StdEncoding.EncodeToString(data)))
			}
		})
}
