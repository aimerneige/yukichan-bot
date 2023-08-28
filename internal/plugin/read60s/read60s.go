package read60s

import (
	b64 "encoding/base64"

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
	zero.OnFullMatchGroup([]string{"今日新闻", "早报", "60s"}).
		SetPriority(7).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			data, err := web.GetData(api)
			if err != nil {
				log.Errorln("[read60s]", err)
				ctx.Send("网络错误，获取早报信息失败。")
				return
			}
			apiMsg := gjson.Get(binary.BytesToString(data), "msg").String()
			if apiMsg == "Success" {
				imageUrl := gjson.Get(binary.BytesToString(data), "imageUrl").String()
				imageData, err := web.GetData(imageUrl)
				if err != nil {
					ctx.Send("网络错误，早报图片获取失败。")
				}
				ctx.Send(message.Image("base64://" + b64.StdEncoding.EncodeToString(imageData)))
			} else {
				ctx.Send("API 错误，无法获取早报图片。")
			}
		})
}
