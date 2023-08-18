package read60s

import (
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	api = "https://api.2xb.cn/zaob"
)

func init() {
	zero.OnFullMatchGroup([]string{"今日新闻", "早报", "60s"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			data, err := web.GetData(api)
			if err != nil {
				log.Errorln("[read60s]", err)
				ctx.Send("获取早报图片失败了呢，可能是服务器网络出问题了罢。")
				return
			}
			apiMsg := gjson.Get(binary.BytesToString(data), "msg").String()
			if apiMsg == "Success" {
				imageUrl := gjson.Get(binary.BytesToString(data), "imageUrl").String()
				ctx.Send(message.Image(imageUrl))
			} else {
				ctx.Send("")
			}
		})
}
