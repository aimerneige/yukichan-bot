package setu

import (
	"fmt"
	"math/rand"
	"os"
	"path"

	b64 "encoding/base64"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/aimerneige/yukichan-bot/internal/pkg/common"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	api = "https://api.lolicon.app/setu/v2"
)

var localSetu []string = []string{
	"/opt/yukichan/seturepo",
}

func init() {
	engine := zero.New()
	common.DefaultSingle.Apply(engine)
	engine.OnFullMatch("/setu", zero.SuperUserPermission).
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
	engine.OnFullMatch("setu", zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			dir := localSetu[rand.Int()%len(localSetu)]
			imgB64, err := getLocalSetu(dir)
			if err != nil {
				ctx.Send("读取本地图片失败了，请查阅后台日志。")
				return
			}
			m := message.Image("base64://" + imgB64)
			if id := ctx.Send(m).ID(); id == 0 {
				ctx.Send("ERROR: 可能被风控或读取图片用时过长，请耐心等待")
			}
		})
	engine.UseMidHandler(common.DefaultSpeedLimit)
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

func getLocalSetu(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Errorln("[setu]", fmt.Sprintf("Fail to read dir %s", dir))
		return "", err
	}
	imgFile := files[rand.Intn(len(files))]
	// 检测是否读到文件夹，如果是则重试三次，否则报错
	for i := 0; i < 3 && imgFile.IsDir(); i++ {
		imgFile = files[rand.Intn(len(files))]
	}
	if imgFile.IsDir() {
		return "", fmt.Errorf("Fail to get a file in dir %s", dir)
	}
	imgBytes, err := os.ReadFile(path.Join(dir, imgFile.Name()))
	if err != nil {
		log.Errorln("[setu]", fmt.Sprintf("Fail to read img file %s", imgFile.Name()), err)
		return "", err
	}
	return b64.StdEncoding.EncodeToString(imgBytes), nil
}
