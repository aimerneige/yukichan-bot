// This plugin call service from aiursoft
// View https://glm.aiursoft.cn/ for more detail.

package glm

import (
	"fmt"

	"github.com/FloatTech/floatbox/binary"
	"github.com/aimerneige/yukichan-bot/internal/pkg/common"
	"github.com/aimerneige/yukichan-bot/internal/pkg/web"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	api = "https://glm.aiursoft.cn/"
	// TODO: add history
	template = "{\"history\":[],\"prompt\":\"%s\"}"
)

func init() {
	engine := zero.New()
	common.DefaultSingle.Apply(engine)
	engine.OnPrefix("//").
		SetPriority(8).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			data, err := web.PostRequest(api, fmt.Sprintf(template, args))
			if err != nil {
				ctx.Send(fmt.Sprintln("[ERROR]:", err))
				return
			}
			apiResponse := binary.BytesToString(data)
			response := gjson.Get(apiResponse, "response").String()
			ctx.Send(response)
		})
	engine.UseMidHandler(common.DefaultSpeedLimit)
}
