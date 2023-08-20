package alipay

import (
	b64 "encoding/base64"
	"fmt"
	"strings"

	"github.com/FloatTech/floatbox/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	api = "https://mm.cqu.cc/share/zhifubaodaozhang/mp3/%v.mp3"
)

func init() {
	zero.OnPrefixGroup([]string{"支付宝到账", "alipay"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			mp3Url := fmt.Sprintf(api, strings.TrimSpace(args))
			data, err := web.GetData(mp3Url)
			if err != nil {
				return
			}
			ctx.SendChain(message.Record("base64://" + b64.StdEncoding.EncodeToString(data)))
		})
}
