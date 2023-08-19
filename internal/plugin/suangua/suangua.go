package suangua

import (
	"embed"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/fs"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

//go:embed assets/gua
var guaImages embed.FS

//go:embed assets/64.json
var guaResultJSONData []byte

func init() {
	zero.OnPrefix("算卦").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			things := strings.TrimSpace(args)
			resultIndex := uint32(0)
			if things != "" {
				resultIndex = calHash(things, ctx.Event.UserID)
			}
			explain, imgB64, err := getSuanguaMessage(resultIndex)
			if err != nil {
				ctx.Send("发生了玄学事故！算卦失败了！")
				return
			}
			msgText := message.Text(explain)
			msgImage := message.Image("base64://" + imgB64)
			ctx.SendChain(msgText, msgImage)
		})
}

func getSuanguaMessage(i uint32) (explain string, imgB64 string, err error) {
	var guaResultJSONObj []string
	if err = json.Unmarshal(guaResultJSONData, &guaResultJSONObj); err != nil {
		log.Errorln("[suangua]", "Assets JSON unmarshal error!", err)
		return
	}
	explain = guaResultJSONObj[i]
	imgPath := path.Join("assets/gua", fmt.Sprintf("%d.jpg", i))
	imgData, err := fs.ReadFile(guaImages, imgPath)
	if err != nil {
		log.Errorln("[suangua]", fmt.Sprintf("Fail to read img %s", imgPath), err)
		return
	}
	imgB64 = b64.StdEncoding.EncodeToString(imgData)
	return
}

// calHash 计算 Hash
func calHash(things string, uin int64) uint32 {
	unixTime := uint32(time.Now().Unix() / 10000)
	thingsHash := stringHash(things)
	uinHash := stringHash(fmt.Sprint(uin))
	return (unixTime+thingsHash+uinHash)%64 + 1
}

func stringHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
