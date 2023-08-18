package morning

import (
	"time"

	"github.com/go-co-op/gocron"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var ()

func init() {
	engine := zero.New()

	go func() {
		s := gocron.NewScheduler(time.UTC)
		s.Every(1).Day().At("00:00").Do(func() {

		})
		s.StartAsync()
	}()

	engine.OnCommand("启用早安").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			// ...
		})

	engine.OnCommand("关闭早安").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			// ...
		})
}
