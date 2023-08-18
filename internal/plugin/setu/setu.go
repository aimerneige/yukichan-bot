package setu

import (
	"fmt"
	"time"

	b64 "encoding/base64"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/extension/single"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	api = "https://api.lolicon.app/setu/v2"
)

var (
	limit = rate.NewManager[int64](time.Second*30, 1)
)

func init() {
	engine := zero.New()

	single.New(
		single.WithKeyFn(func(ctx *zero.Ctx) int64 {
			return ctx.Event.UserID
		}),
		single.WithPostFn[int64](func(ctx *zero.Ctx) {
			ctx.Send("您有操作正在执行，请稍后再试!")
		}),
	).Apply(engine)

	engine.OnCommandGroup([]string{"来点色图", "setu"}, zero.SuperUserPermission).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			imgB64, err := getSetu()
			if err != nil {
				log.Errorln("[setu]", err)
				ctx.Send("色图娘好像挂掉了呢。")
				return
			}
			if len(imgB64) == 0 {
				log.Errorln("[setu]", err)
				ctx.Send("色图下载失败了捏。")
				return
			}
			m := message.Image("base64://" + imgB64)
			if id := ctx.Send(m).ID(); id == 0 {
				ctx.Send("ERROR: 可能被风控或下载图片用时过长，请耐心等待")
			}
		})

	engine.UseMidHandler(func(ctx *zero.Ctx) bool { // 限速器
		if !limit.Load(ctx.Event.UserID).Acquire() {
			ctx.Send("您的请求太快，请稍后重试0x0...")
			return false
		}
		return true
	})
}

func getSetu() (string, error) {
	data, err := web.GetData(api)
	if err != nil {
		return "", err
	}
	apiResponse := binary.BytesToString(data)
	apiError := gjson.Get(apiResponse, "error").String()
	if apiError != "" {
		return "", fmt.Errorf(apiError)
	}
	imgUrl := gjson.Get(apiResponse, "data.0.urls.original").String()
	log.Debugln("[setu]", "setu url:", imgUrl)
	imgData, err := web.GetData(imgUrl)
	if err != nil {
		return "", err
	}
	imgB64 := b64.StdEncoding.EncodeToString(imgData)
	log.Debugln("[setu]", "imgB64 length:", len(imgB64))
	return imgB64, nil
}
