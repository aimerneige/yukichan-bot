package tarot

import (
	"bytes"
	"embed"
	b64 "encoding/base64"
	"image"
	"image/png"
	"io/fs"
	"math/rand"
	"path"
	"strconv"

	"github.com/aimerneige/yukichan-bot/internal/pkg/common"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

//go:embed assets/deck
var deckEmbedFS embed.FS

//go:embed assets/tarot.jpg
var tarotCommandImg []byte

func init() {
	engine := zero.New()
	common.DefaultSingle.Apply(engine)

	engine.OnFullMatchGroup([]string{"塔罗", "塔罗牌", "tarot"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			helpText := message.Text("支持指令\n「运势预测」（单张牌预测运势）\n「塔罗占卜」（三张牌进行占卜）\n「抽塔罗牌 5」（抽取指定张数的塔罗牌）")
			helpImage := message.Image("base64://" + b64.StdEncoding.EncodeToString(tarotCommandImg))
			ctx.SendChain(helpText, helpImage)
		})
	engine.OnFullMatch("运势预测").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(drawCard(1)...)
		})
	engine.OnFullMatch("塔罗占卜").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(drawCard(3)...)
		})
	engine.OnPrefix("抽塔罗牌").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				return
			}
			count, err := strconv.ParseInt(args, 10, 64)
			if err != nil {
				ctx.Send("牌数解析失败")
				return
			}
			if count <= 0 || count > 8 {
				count = 1
			}
			ctx.SendChain(drawCard(int(count))...)
		})

	engine.UseMidHandler(common.DefaultSpeedLimit)
}

func drawCard(number int) []message.MessageSegment {
	theme := "classic"
	if (rand.Int() % 3) == 0 {
		theme = "bilibili"
	}
	deckPath := path.Join("./assets/deck/", theme)
	cardImages, err := fs.ReadDir(deckEmbedFS, deckPath)
	if err != nil {
		log.Errorln("[tarot]", "Fail to read cardImages.", err)
		return []message.MessageSegment{message.Text("发生错误，无法读取塔罗图片。")}
	}
	cardImagePaths := make([]string, 0, 78)
	for _, card := range cardImages {
		if !card.IsDir() {
			_imgPath := path.Join(deckPath, card.Name())
			cardImagePaths = append(cardImagePaths, _imgPath)
		}
	}
	cardResult := randomDraw(cardImagePaths, number)
	replyMsg := []message.MessageSegment{}
	replyMsg = append(replyMsg, message.Text("你抽到的卡牌如下："))
	for _, imgPath := range cardResult {
		// 读取图片
		imgData, err := fs.ReadFile(deckEmbedFS, imgPath)
		if err != nil {
			log.Errorln("[tarot]", "Fail to read img", err)
			replyMsg = append(replyMsg, message.Text("[ERROR] 读取图片失败\n"))
		}
		// 翻转图片，实现正逆位
		if (rand.Int() % 2) == 0 {
			flippedImageData, err := flipImage(imgData)
			if err != nil {
				log.Errorln("[tarot]", "Fail to flip image", err)
				continue
			} else {
				imgData = flippedImageData
			}
		}
		replyMsg = append(replyMsg, message.Image("base64://"+b64.StdEncoding.EncodeToString(imgData)))
	}
	return replyMsg
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

func flipImage(imageData []byte) ([]byte, error) {
	// Decode the []byte into an image.Image.
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Errorln("[tarot]", "Fail to decode image", err)
		return nil, err
	}

	// Flip the image vertically.
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			flipped.Set(bounds.Max.X-x-1, bounds.Max.Y-y-1, img.At(x, y))
		}
	}

	// Encode the flipped image as a []byte.
	var buf bytes.Buffer
	if err := png.Encode(&buf, flipped); err != nil {
		log.Errorln("[tarot]", "Fail to encode flipped image", err)
		return nil, err
	}

	// Return the flipped []byte image.
	flippedData := buf.Bytes()
	return flippedData, nil
}