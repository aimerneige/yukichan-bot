package imgsave

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aimerneige/yukichan-bot/internal/pkg/common"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const imgRootDir = "./data/imgsave"

func init() {
	initTargetDir()
	engine := zero.New()
	common.DefaultSingle.Apply(engine)
	engine.OnPrefix("加图").FirstPriority().SetBlock(true).Handle(handleImgSave)
	engine.OnPrefix("创建图库", zero.OwnerPermission).SecondPriority().SetBlock(true).Handle(func(ctx *zero.Ctx) {})
	engine.OnPrefix("来张").SecondPriority().SetBlock(true).Handle(func(ctx *zero.Ctx) {})
	engine.UseMidHandler(common.DefaultSpeedLimit)
}

func allowedImgList() [][]string {
	return [][]string{
		{"新田慧海", "新田恵海", "慧海", "恵海", "emi"},
		{"高坂穗乃果", "高坂穂乃果 ", "穗乃果", "穂乃果", "ほのか", "honoka"},
		{"表情包"},
	}
}

func handleImgSave(ctx *zero.Ctx) {
	msg := ctx.Event.Message
	if len(msg) < 2 {
		return
	}
	if msg[0].Type != "text" {
		return
	}
	userTarget := parseUserTarget(msg[0].Data["text"])
	dirTarget := toTargetDirName(userTarget)
	if dirTarget == "" {
		ctx.Send(fmt.Sprintf("图片库「%s」不存在！可联系机器人管理员创建。", userTarget))
		return
	}
	imgList := []UserImage{}
	for _, ele := range msg {
		if ele.Type == "image" {
			imgList = append(imgList, UserImage{
				Name:   ele.Data["file_unique"],
				Url:    ele.Data["url"],
				Sender: ctx.Event.Sender.ID,
			})
		}
	}
	if len(imgList) < 1 {
		ctx.Send("请发送图片！")
		// TODO
		return
	}
	successCount := 0
	failCount := 0
	targetDirPath := path.Join(imgRootDir, dirTarget)
	for _, img := range imgList {
		filePath := path.Join(targetDirPath, img.getRealFileName())
		if downloadImage(img.Url, filePath) {
			successCount++
		} else {
			failCount++
		}
	}
	ctx.Send(fmt.Sprintf("任务完成，已存入「%s」。成功保存 %d 张图片，失败 %d 张", dirTarget, successCount, failCount))
}

type UserImage struct {
	Name   string
	Url    string
	Sender int64
}

func (i UserImage) getRealFileName() string {
	timeStr := time.Now().Format("2006-01-02-15-04-05")
	return fmt.Sprintf("%s-%d-%s", timeStr, i.Sender, i.Name)
}

func downloadImage(imgUrl, filePath string) bool {
	response, err := http.Get(imgUrl)
	if err != nil {
		log.Errorln("[imgsave] fail to download image", err)
		return false
	}
	defer response.Body.Close()
	file, err := os.Create(filePath)
	if err != nil {
		log.Errorln("[imgsave] fail to create image", err)
		return false
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Errorln("[imgsave] fail to save image", err)
		return false
	}
	return true
}

func parseUserTarget(userSource string) string {
	userTarget := strings.TrimSpace(userSource)
	userTarget = strings.TrimPrefix(userTarget, "加图")
	userTarget = strings.TrimSpace(userTarget)
	userTarget = strings.ToLower(userTarget)
	return userTarget
}

func toTargetDirName(s string) string {
	for _, dir := range allowedImgList() {
		for _, name := range dir {
			if s == name {
				return dir[0]
			}
		}
	}
	return ""
}

func initTargetDir() {
	for _, dir := range allowedImgList() {
		dirPath := path.Join(imgRootDir, dir[0])
		_, err := os.Stat(dirPath)
		if os.IsNotExist(err) {
			os.Mkdir(dirPath, 0755)
		}
	}
}
