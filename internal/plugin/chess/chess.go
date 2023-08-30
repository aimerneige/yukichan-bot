package chess

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aimerneige/yukichan-bot/internal/plugin/chess/database"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

//go:embed assets/cheese.jpeg
var cheeseData []byte

//go:embed assets/help.txt
var helpString string

const sqliteDBFile = "./data/chess/chess.db"

func init() {
	database.InitDatabase(sqliteDBFile)
	engine := zero.New()
	engine.OnFullMatchGroup([]string{"下棋", "chess"}, zero.OnlyGroup).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			userName := ctx.Event.Sender.NickName
			groupCode := ctx.Event.GroupID
			if replyMessage := Game(groupCode, userUin, userName); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"认输", "resign"}, zero.OnlyGroup).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			groupCode := ctx.Event.GroupID
			if replyMessage := Resign(groupCode, userUin); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"和棋", "draw"}, zero.OnlyGroup).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			groupCode := ctx.Event.GroupID
			if replyMessage := Draw(groupCode, userUin); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"中断", "abort"}, zero.OnlyGroup, zero.AdminPermission).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			groupCode := ctx.Event.GroupID
			if replyMessage := Abort(groupCode); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"盲棋", "blind"}, zero.OnlyGroup).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			userName := ctx.Event.Sender.NickName
			groupCode := ctx.Event.GroupID
			if replyMessage := Blindfold(groupCode, userUin, userName); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnRegex("^[!|！]([0-8]|[R|N|B|Q|K|O|a-h|x]|[-|=|+])+$", zero.OnlyGroup).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			groupCode := ctx.Event.GroupID
			userMsgStr := ctx.State["regex_matched"].([]string)[0]
			userMsgStr = strings.Replace(userMsgStr, "！", "!", 1)
			moveStr := userMsgStr[1:]
			if replyMessage := Play(userUin, groupCode, moveStr); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"排行榜", "ranking"}).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if replyMessage := Ranking(); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"等级分", "rate"}).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userUin := ctx.Event.UserID
			userName := ctx.Event.Sender.NickName
			if replyMessage := Rate(userUin, userName); len(replyMessage) >= 1 {
				ctx.Send(replyMessage)
			}
		})
	engine.OnFullMatchGroup([]string{"帮助", "help"}, zero.OnlyGroup).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.Send(helpString)
		})
	engine.OnPrefixGroup([]string{"清空等级分", ".clean.rate"}, zero.SuperUserPermission).
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			if playerUin, err := strconv.ParseInt(strings.TrimSpace(args), 10, 64); err == nil && playerUin > 0 {
				if replyMessage := CleanUserRate(playerUin); len(replyMessage) >= 1 {
					ctx.Send(replyMessage)
				}
			} else {
				ctx.Send(fmt.Sprintf("解析失败「%s」不是正确的 QQ 号。", args))
			}
		})
	engine.OnPrefix("/pgn2gif").
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			const PATTERN = "([0-9]|[R|N|B|Q|K|O|a-h|x]|[-|.|=|+|#|/| |\n])+"
			reg := regexp.MustCompile(PATTERN)
			if reg.FindString(args) == args {
				userUin := ctx.Event.UserID
				if replyMessage := GenerateGIF(userUin, args); len(replyMessage) >= 1 {
					ctx.Send(replyMessage)
				}
			}
		})
	engine.OnFullMatch("cheese").
		SetPriority(2).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(
				message.Text("Chess Cheese Cheese Chess"),
				message.Image("base64://"+base64.StdEncoding.EncodeToString(cheeseData)),
			)
		})
}
