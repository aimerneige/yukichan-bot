package donate

import (
	_ "embed"

	b64 "encoding/base64"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

//go:embed assets/alipay.jpg
var alipay []byte

//go:embed assets/wechat.jpg
var wechat []byte

func init() {
	zero.OnCommandGroup([]string{"捐赠", "donate"}).
		Handle(func(ctx *zero.Ctx) {
			alipayB64 := b64.StdEncoding.EncodeToString(alipay)
			wechatB64 := b64.StdEncoding.EncodeToString(wechat)
			ctx.SendChain(
				message.Text("感谢捐赠！"),
				message.Image("base64://"+alipayB64),
				message.Image("base64://"+wechatB64),
			)
		})
}
