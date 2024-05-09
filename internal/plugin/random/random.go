package random

import (
	"embed"
	b64 "encoding/base64"
	"fmt"
	"io/fs"
	"math/rand"
	"path"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const CardImgPath = "assets/card"

//go:embed assets/cxk.mp3
var cxkBytes []byte

//go:embed assets/card
var cardEmbedFS embed.FS

func init() {
	engine := zero.New()
	engine.OnFullMatchGroup([]string{"掷硬币", "/coin"}).
		SetPriority(6).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			replyTextMessage := "掷出了正面。"
			if (rand.Int() % 2) == 0 {
				replyTextMessage = "掷出了反面。"
			}
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text(replyTextMessage))
		})
	engine.OnFullMatch("只因币").
		SetPriority(6).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Record("base64://" + b64.StdEncoding.EncodeToString(cxkBytes)))
		})
	engine.OnFullMatchGroup([]string{"掷骰子", "/dice"}).
		SetPriority(6).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			replyTextMessage := fmt.Sprintf("掷出了 %d 点。", (rand.Int()%6)+1)
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text(replyTextMessage))
		})
	engine.OnFullMatchGroup([]string{"抽扑克", "/card"}).
		SetPriority(6).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.Send(drawCard(1))
		})
}

func drawCard(number int) message.Message {
	cardImages, err := fs.ReadDir(cardEmbedFS, CardImgPath)
	log.Errorln(cardImages)
	if err != nil {
		log.Errorln("[random]", "Fail to read card images", err)
		return []message.MessageSegment{message.Text("发生错误，无法读取扑克图片")}
	}
	cardImagePaths := make([]string, 0, 78)
	for _, card := range cardImages {
		if !card.IsDir() {
			_imgPath := path.Join(CardImgPath, card.Name())
			cardImagePaths = append(cardImagePaths, _imgPath)
		}
	}
	cardResult := randomDraw(cardImagePaths, number)
	log.Debugln("[tarot]", fmt.Sprintf("cards: %v", cardResult))
	imgs := make([]message.MessageSegment, number)
	for i, imgPath := range cardResult {
		imgData, err := fs.ReadFile(cardEmbedFS, imgPath)
		// 读取图片
		if err != nil {
			log.Errorln("[tarot]", "Fail to read card image", err)
			imgs[i] = message.Text("[ERROR] 读取图片失败，请查阅后台日志。\n")
			continue
		}
		imgs[i] = message.Image("base64://" + b64.StdEncoding.EncodeToString(imgData))
	}
	return imgs
}

// randomDraw 随机不放回抽取
func randomDraw(s []string, k int) []string {
	n := len(s)
	if k > n {
		k = n
	}

	result := make([]string, k)
	for i := 0; i < k; i++ {
		j := rand.Intn(n-i) + i
		result[i] = s[j]
		s[i], s[j] = s[j], s[i]
	}

	return result
}
