// This file is modified from ZeroBot-Plugin.
// View this link for more detail:
// https://github.com/FloatTech/ZeroBot-Plugin/tree/master/plugin/manager

package manager

import (
	"github.com/FloatTech/floatbox/math"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	hint = "====群管====\n" +
		"- 禁言@QQ 1分钟\n" +
		"- 解除禁言 @QQ\n" +
		"- 我要自闭 1分钟\n" +
		"- 开启全员禁言\n" +
		"- 解除全员禁言\n" +
		"- 修改名片@QQ XXX\n" +
		"- 修改头衔@QQ XXX\n" +
		"- 申请头衔 XXX\n" +
		"- 踢出群聊@QQ\n"
)

func init() {
	engine := zero.New()
	engine.OnFullMatch("群管帮助", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.Send(hint)
		})
	// 踢出群聊
	engine.OnRegex(`^踢出群聊.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetThisGroupKick(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被踢出群聊的人的qq
				false,
			)
			nickname := ctx.GetThisGroupMemberInfo( // 被踢出群聊的人的昵称
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被踢出群聊的人的qq
				false,
			).Get("nickname").Str
			ctx.SendChain(message.Text("残念~ " + nickname + " 被放逐"))
		})
	// 开启全体禁言
	engine.OnRegex(`^开启全员禁言$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetThisGroupWholeBan(true)
			ctx.SendChain(message.Text("全员自闭开始~"))
		})
	// 解除全员禁言
	engine.OnRegex(`^解除全员禁言$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetThisGroupWholeBan(false)
			ctx.SendChain(message.Text("全员自闭结束~"))
		})
	// 禁言
	engine.OnRegex(`^禁言.*?(\d+).*?\s(\d+)(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			duration := math.Str2Int64(ctx.State["regex_matched"].([]string)[2])
			switch ctx.State["regex_matched"].([]string)[3] {
			case "分钟":
				//
			case "小时":
				duration *= 60
			case "天":
				duration *= 60 * 24
			default:
				//
			}
			if duration >= 43200 {
				duration = 43199 // qq禁言最大时长为一个月
			}
			ctx.SetThisGroupBan(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 要禁言的人的qq
				duration*60, // 要禁言的时间（分钟）
			)
			ctx.SendChain(message.Text("小黑屋收留成功~"))
		})
	// 解除禁言
	engine.OnRegex(`^解除禁言.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetThisGroupBan(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 要解除禁言的人的qq
				0,
			)
			ctx.SendChain(message.Text("小黑屋释放成功~"))
		})
	// 自闭禁言
	engine.OnRegex(`^(我要自闭|禅定).*?(\d+)(.*)`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			duration := math.Str2Int64(ctx.State["regex_matched"].([]string)[2])
			switch ctx.State["regex_matched"].([]string)[3] {
			case "分钟", "min", "mins", "m":
				break
			case "小时", "hour", "hours", "h":
				duration *= 60
			case "天", "day", "days", "d":
				duration *= 60 * 24
			default:
				break
			}
			if duration >= 43200 {
				duration = 43199 // qq禁言最大时长为一个月
			}
			ctx.SetThisGroupBan(
				ctx.Event.UserID,
				duration*60, // 要自闭的时间（分钟）
			)
			ctx.SendChain(message.Text("那我就不手下留情了~"))
		})
	// 修改名片
	engine.OnRegex(`^修改名片.*?(\d+).+?\s*(.*)$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if len(ctx.State["regex_matched"].([]string)[2]) > 60 {
				ctx.SendChain(message.Text("名字太长啦！"))
				return
			}
			ctx.SetThisGroupCard(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被修改群名片的人
				ctx.State["regex_matched"].([]string)[2],                 // 修改成的群名片
			)
			ctx.SendChain(message.Text("嗯！已经修改了"))
		})
	// 修改头衔
	engine.OnRegex(`^修改头衔.*?(\d+).+?\s*(.*)$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			sptitle := ctx.State["regex_matched"].([]string)[2]
			if sptitle == "" {
				ctx.SendChain(message.Text("头衔不能为空！"))
				return
			} else if len(sptitle) > 18 {
				ctx.SendChain(message.Text("头衔太长啦！"))
				return
			}
			ctx.SetThisGroupSpecialTitle(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被修改群头衔的人
				sptitle, // 修改成的群头衔
			)
			ctx.SendChain(message.Text("嗯！已经修改了"))
		})
	// 申请头衔
	engine.OnRegex(`^申请头衔\s*(.*)$`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			sptitle := ctx.State["regex_matched"].([]string)[1]
			if sptitle == "" {
				ctx.SendChain(message.Text("头衔不能为空！"))
				return
			} else if len(sptitle) > 18 {
				ctx.SendChain(message.Text("头衔太长啦！"))
				return
			}
			ctx.SetThisGroupSpecialTitle(
				ctx.Event.UserID, // 被修改群头衔的人
				sptitle,          // 修改成的群头衔
			)
			ctx.SendChain(message.Text("嗯！不错的头衔呢~"))
		})
}
