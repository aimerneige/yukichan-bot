// This file is modified from ZeroBot-Plugin.
// View this link for more detail:
// https://github.com/FloatTech/ZeroBot-Plugin/tree/master/plugin/alipayvoice

package alipay

import (
	b64 "encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	api = "https://mm.cqu.cc/share/zhifubaodaozhang/mp3/%v.mp3"
)

func init() {
	zero.OnPrefixGroup([]string{"支付宝到账", "alipay"}).
		SetPriority(8).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			if moneyCount, err := strconv.ParseFloat(strings.TrimSpace(args), 64); err == nil && moneyCount > 0 {
				if data, err := web.GetData(fmt.Sprintf(api, moneyCount)); err == nil {
					ctx.SendChain(message.Record("base64://" + b64.StdEncoding.EncodeToString(data)))
				}
			}
		})
}
