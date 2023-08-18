package fortune

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	zero.OnPrefix("求签").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			things := strings.TrimSpace(args)
			if things == "" {
				return
			}
			fortuneResult := drawAFortuneStick(things, ctx.Event.UserID)
			ctx.Send(fmt.Sprintf("所求事项\"%s\"的求签结果为: %s", things, fortuneResult))
		})
}

func drawAFortuneStick(things string, uin int64) string {
	unixTime := uint32(time.Now().Unix() / 10000)
	thingsHash := stringHash(things)
	uinHash := stringHash(fmt.Sprint(uin))
	return getFortuneResult(unixTime + thingsHash + uinHash)
}

func getFortuneResult(hash uint32) string {
	var result string
	switch key := hash % 100; {
	case key < 2:
		result = "上吉" // 2
	case key < 10:
		result = "大吉" // 8
	case key < 38:
		result = "上上" // 28
	case key < 42:
		result = "上中" // 4
	case key < 45:
		result = "上平" // 3
	case key < 46:
		result = "上" // 1
	case key < 49:
		result = "中吉" // 3
	case key < 51:
		result = "中上" // 2
	case key < 57:
		result = "中中" // 6
	case key < 66:
		result = "中平" // 9
	case key < 71:
		result = "中" // 5
	case key < 72:
		result = "平中" // 1
	case key < 73:
		result = "平平" // 1
	case key < 74:
		result = "平" // 1
	case key < 99:
		result = "下" // 25
	case key < 100:
		result = "下下" // 1
	default:
		result = "大凶"
	}
	return result
}

func stringHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
