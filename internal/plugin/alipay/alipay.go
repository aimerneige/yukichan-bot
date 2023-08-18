package alipay

import (
	"fmt"
	"strings"

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
			ctx.SendChain(message.Record(fmt.Sprintf(api, strings.TrimSpace(args))))
		})
}
