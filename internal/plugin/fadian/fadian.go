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

// PostJSON post json object
type PostJSON struct {
	Post []string `json:"post"`
}

func init() {
	zero.OnPrefix("每日发癫").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			name := "小乌贼"
			nameString := strings.TrimSpace(ctx.Event.Message.ExtractPlainText()[12:])
			if nameString != "" {
				name = nameString
			}
			ctx.Send(getFadianText(name))
		})
}

func getFadianText(name string) string {
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
