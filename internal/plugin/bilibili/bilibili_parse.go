// This file is modified from ZeroBot-Plugin.
// View this link for more detail:
// https://github.com/FloatTech/ZeroBot-Plugin/tree/master/plugin/bilibili

package bilibili

import (
	"regexp"

	bz "github.com/FloatTech/AnimeAPI/bilibili"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	searchVideo      = `bilibili.com\\?/video\\?/(?:av(\d+)|([bB][vV][0-9a-zA-Z]+))`
	searchDynamic    = `(t.bilibili.com|m.bilibili.com\\?/dynamic)\\?/(\d+)`
	searchArticle    = `bilibili.com\\?/read\\?/(?:cv|mobile\\?/)(\d+)`
	searchLiveRoom   = `live.bilibili.com\\?/(\d+)`
	searchVideoRe    = regexp.MustCompile(searchVideo)
	searchDynamicRe  = regexp.MustCompile(searchDynamic)
	searchArticleRe  = regexp.MustCompile(searchArticle)
	searchLiveRoomRe = regexp.MustCompile(searchLiveRoom)
)

func init() {
	engine := zero.New()
	engine.OnRegex(`((b23|acg).tv|bili2233.cn)/[0-9a-zA-Z]+`).
		SetPriority(9).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			url := ctx.State["regex_matched"].([]string)[0]
			realurl, err := bz.GetRealURL("https://" + url)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			switch {
			case searchVideoRe.MatchString(realurl):
				ctx.State["regex_matched"] = searchVideoRe.FindStringSubmatch(realurl)
				handleVideo(ctx)
			case searchDynamicRe.MatchString(realurl):
				ctx.State["regex_matched"] = searchDynamicRe.FindStringSubmatch(realurl)
				handleDynamic(ctx)
			case searchArticleRe.MatchString(realurl):
				ctx.State["regex_matched"] = searchArticleRe.FindStringSubmatch(realurl)
				handleArticle(ctx)
			case searchLiveRoomRe.MatchString(realurl):
				ctx.State["regex_matched"] = searchLiveRoomRe.FindStringSubmatch(realurl)
				handleLive(ctx)
			}
		})
	engine.OnRegex(searchVideo).SetPriority(9).SetBlock(true).Handle(handleVideo)
	engine.OnRegex(searchDynamic).SetPriority(9).SetBlock(true).Handle(handleDynamic)
	engine.OnRegex(searchArticle).SetPriority(9).SetBlock(true).Handle(handleArticle)
	engine.OnRegex(searchLiveRoom).SetPriority(9).SetBlock(true).Handle(handleLive)
}

func handleVideo(ctx *zero.Ctx) {
	id := ctx.State["regex_matched"].([]string)[1]
	if id == "" {
		id = ctx.State["regex_matched"].([]string)[2]
	}
	card, err := bz.GetVideoInfo(id)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	msg, err := videoCard2msg(card)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	ctx.SendChain(msg...)
}

func handleDynamic(ctx *zero.Ctx) {
	msg, err := dynamicDetail(ctx.State["regex_matched"].([]string)[2])
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	ctx.SendChain(msg...)
}

func handleArticle(ctx *zero.Ctx) {
	card, err := bz.GetArticleInfo(ctx.State["regex_matched"].([]string)[1])
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	ctx.SendChain(articleCard2msg(card, ctx.State["regex_matched"].([]string)[1])...)
}

func handleLive(ctx *zero.Ctx) {
	card, err := bz.GetLiveRoomInfo(ctx.State["regex_matched"].([]string)[1])
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
		return
	}
	ctx.SendChain(liveCard2msg(card)...)
}
