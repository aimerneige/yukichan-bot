package fadian

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"strings"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

//go:embed assets/post.json
var postJSONData []byte

//go:embed assets/text.json
var textJSONData []byte

// PostJSON post json object
type PostJSON struct {
	Post []string `json:"post"`
}

// TextJSON text json object
type TextJSON struct {
	Text []string `json:"text"`
}

func init() {
	engine := zero.New()
	engine.OnPrefix("每日发癫").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := "小乌贼"
			nameString := strings.TrimSpace(ctx.Event.Message.ExtractPlainText()[12:])
			if nameString != "" {
				name = nameString
			}
			ctx.Send(getFadianPost(name))
		})
	engine.OnFullMatchGroup([]string{"小作文", "发大病"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.Send(getFadianText())
	})
}

func getFadianPost(name string) string {
	var postJSON PostJSON
	err := json.Unmarshal(postJSONData, &postJSON)
	if err != nil {
		log.Errorln("[fadian]", err)
		return "解析 JSON 失败，请查阅后台日志。"
	}
	postString := postJSON.Post[rand.Intn(len(postJSON.Post))]
	postString = strings.ReplaceAll(postString, "阿咪", name)
	return postString
}

func getFadianText() string {
	var textJSON TextJSON
	err := json.Unmarshal(postJSONData, &textJSON)
	if err != nil {
		log.Errorln("[fadian]", err)
		return "解析 JSON 失败，请查阅后台日志。"
	}
	postString := textJSON.Text[rand.Intn(len(textJSON.Text))]
	return postString
}
