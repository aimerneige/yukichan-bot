package chess

import (
	_ "embed"
	"encoding/base64"
	"strings"

	"github.com/aimerneige/yukichan-bot/internal/config"
	"github.com/aimerneige/yukichan-bot/internal/plugin/chess/database"
	"github.com/aimerneige/yukichan-bot/internal/plugin/chess/service"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

//go:embed assets/cheese.jpeg
var cheeseData []byte

//go:embed assets/help.txt
var helpString string

func init() {
	sqliteDatabase := database.SqliteDatabase{FilePath: config.Conf.DataFile["chess"]}
	database.InitDatabase(sqliteDatabase)
	engine := zero.New()
	engine.OnFullMatchGroup([]string{"下棋", "chess"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			userName := ctx.Event.Sender.NickName
			groupCode := ctx.Event.GroupID
			if replyMessage := service.Game(groupCode, userUin, userName); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"认输", "resign"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			groupCode := ctx.Event.GroupID
			if replyMessage := service.Resign(groupCode, userUin); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"和棋", "draw"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			groupCode := ctx.Event.GroupID
			if replyMessage := service.Draw(groupCode, userUin); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"中断", "abort"}, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			groupCode := ctx.Event.GroupID
			if replyMessage := service.Abort(groupCode); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"盲棋", "blind"}, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			userName := ctx.Event.Sender.NickName
			groupCode := ctx.Event.GroupID
			if replyMessage := service.Blindfold(groupCode, userUin, userName); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnRegex("[!|！]([0-9]|[A-Z]|[a-z]|=|-)+", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			groupCode := ctx.Event.GroupID
			userMsgStr := ctx.Event.Message.ExtractPlainText()
			userMsgStr = strings.Replace(userMsgStr, "！", "!", 1)
			moveStr := userMsgStr[1:]
			if replyMessage := service.Play(userUin, groupCode, moveStr); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"排行榜", "ranking"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if replyMessage := service.Ranking(); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"等级分", "rate"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			userName := ctx.Event.Sender.NickName
			if replyMessage := service.Rate(userUin, userName); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"帮助", "help"}, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.Send(helpString)
	})
	engine.OnFullMatch("cheese").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(
			message.Text("Chess Cheese Cheese Chess"),
			message.Image("base64://"+base64.StdEncoding.EncodeToString(cheeseData)),
		)
	})
}
