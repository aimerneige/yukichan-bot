package github

import (
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	repoAPI  = "https://api.github.com/repos/"
	imageAPI = "https://opengraph.githubassets.com/1a/"
)

func init() {
	zero.OnPrefix("https://github.com/").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			repoInfo := ctx.Event.Message.ExtractPlainText()[19:]
			// 去除域名后内容过短
			if len(repoInfo) <= 2 {
				return
			}
			// 删除末尾的 /
			if repoInfo[len(repoInfo)-1] == '/' {
				repoInfo = repoInfo[:len(repoInfo)-1]
			}
			// 检查是否同时含有用户名和仓库名
			repoInfoSlice := strings.Split(repoInfo, "/")
			if len(repoInfoSlice) != 2 {
				return
			}
			// 检查仓库是否存在
			data, err := web.GetData(repoAPI + repoInfo)
			if err != nil {
				log.Errorln("[github]", "Fail to check repo status", err)
				return
			}
			repoStatusMessage := gjson.Get(binary.BytesToString(data), "message").String()
			if repoStatusMessage == "Not Found" {
				return
			}
			// 发送仓库图片
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("https://github.com/"+repoInfo), message.Image(imageAPI+repoInfo))
		})
}
