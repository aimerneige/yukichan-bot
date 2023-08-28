// This file is modified from ZeroBot-Plugin.
// View this link for more detail:
// https://github.com/FloatTech/ZeroBot-Plugin/tree/master/plugin/wangyiyun

package wangyiyun

import (
	"github.com/FloatTech/floatbox/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

const (
	wangyiyunURL     = "https://api.gmit.vip/Api/HotComments?format=text"
	wangyiyunReferer = "https://api.gmit.vip/"
	ua               = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
)

func init() {
	zero.OnFullMatchGroup([]string{"来份网易云热评", "/wyy"}).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			data, err := web.RequestDataWith(web.NewDefaultClient(), wangyiyunURL, "GET", wangyiyunReferer, ua, nil)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Text(helper.BytesToString(data)))
		})
}
