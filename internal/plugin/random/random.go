package random

import (
	_ "embed"
	b64 "encoding/base64"
	"fmt"
	"math/rand"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

//go:embed assets/cxk.mp3
var cxkBytes []byte

func init() {
	engine := zero.New()
	engine.OnFullMatchGroup([]string{"掷硬币", "/coin"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			replyTextMessage := "掷出了正面。"
			if (rand.Int() % 2) == 0 {
				replyTextMessage = "掷出了反面。"
			}
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text(replyTextMessage))
		})
	engine.OnFullMatch("只因币").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Record("base64://" + b64.StdEncoding.EncodeToString(cxkBytes)))
		})
	engine.OnFullMatchGroup([]string{"掷骰子", "/dice"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			replyTextMessage := fmt.Sprintf("掷出了 %d 点。", (rand.Int()%2)+1)
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text(replyTextMessage))
		})
}
