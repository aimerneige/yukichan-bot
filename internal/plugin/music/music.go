package music

import (
	"io"
	"net/http"
	"net/url"

	"github.com/aimerneige/yukichan-bot/internal/pkg/common"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := zero.New()
	common.DefaultSingle.Apply(engine)
	engine.OnPrefixGroup([]string{"点歌", "music"}).
		SetPriority(5).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				args = "My Dearest"
			}
			musicId := queryNeteaseMusic(args)
			if musicId == 0 {
				musicId = int64(825343)
			}
			ctx.SendChain(message.Music("163", musicId))
		})
	engine.UseMidHandler(common.DefaultSpeedLimit)
}

func queryNeteaseMusic(musicName string) int64 {
	client := http.Client{}
	req, err := http.NewRequest("GET", "http://music.163.com/api/search/get?type=1&s="+url.QueryEscape(musicName), nil)
	if err != nil {
		return 0
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36 Edg/87.0.664.66")
	res, err := client.Do(req)
	if err != nil {
		return 0
	}
	data, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return 0
	}
	return gjson.ParseBytes(data).Get("result.songs.0.id").Int()
}
