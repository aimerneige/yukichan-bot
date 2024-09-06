package imgsave

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	b64 "encoding/base64"

	"github.com/aimerneige/yukichan-bot/internal/pkg/common"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gopkg.in/yaml.v3"
)

//go:embed assets/help.txt
var helpInfo string

const (
	configFilePath = "./config/imgsave.yaml"
)

type ImgsaveConfig struct {
	Rootdit   string     `yaml:"rootdir"`
	Repos     [][]string `yaml:"repos"`
	Blacklist []int64    `yaml:"blacklist"`
}

var config ImgsaveConfig

func init() {
	initConfig()
	initTargetDir()
	engine := zero.New()
	common.DefaultSingle.Apply(engine)
	engine.OnPrefix("加图").FirstPriority().SetBlock(true).Handle(handleImgSave)
	engine.OnPrefix("来张").SecondPriority().SetBlock(true).Handle(handleImgGet)
	engine.OnFullMatch("图库帮助").ThirdPriority().SetBlock(true).Handle(handleImgsaveHelp)
	engine.OnFullMatch("图库列表").ThirdPriority().SetBlock(true).Handle(handleRepoList)
	engine.OnPrefix("创建图库", zero.SuperUserPermission).SecondPriority().SetBlock(true).Handle(handleRepoCreate)
	engine.OnPrefix("图库加黑名单", zero.SuperUserPermission).ThirdPriority().SetBlock(true).Handle(handleBlacklistAdd)
	engine.OnPrefix("图库去黑名单", zero.SuperUserPermission).ThirdPriority().SetBlock(true).Handle(handleBlacklistDel)
	engine.UseMidHandler(common.DefaultSpeedLimit)
}

func handleImgSave(ctx *zero.Ctx) {
	msg := ctx.Event.Message
	if len(msg) < 2 {
		return
	}
	if msg[0].Type != "text" {
		return
	}
	if inBlacklist(ctx.Event.Sender.ID) {
		ctx.Send("你已被拉黑！如有疑问请联系 bot 管理。")
		return
	}
	userArgs := ctx.State["args"].(string)
	userTarget := strings.TrimSpace(userArgs)
	userTarget = strings.ToLower(userTarget)
	dirTarget := toTargetDirName(userTarget)
	if dirTarget == "" {
		ctx.Send(fmt.Sprintf("图库「%s」不存在！可联系机器人管理员创建。", userTarget))
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
		// TODO 多轮对话
		return
	}
	successCount := 0
	failCount := 0
	targetDirPath := path.Join(config.Rootdit, dirTarget)
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

func handleImgGet(ctx *zero.Ctx) {
	userArgs := ctx.State["args"].(string)
	userTarget := strings.TrimSpace(userArgs)
	userTarget = strings.ToLower(userTarget)
	dirTarget := toTargetDirName(userTarget)
	if dirTarget == "" {
		ctx.Send(fmt.Sprintf("图库「%s」不存在！可联系机器人管理员创建。", userTarget))
		return
	}
	dirPath := path.Join(config.Rootdit, dirTarget)
	imgB64, err := getLocalImage(dirPath)
	if err != nil {
		ctx.Send(fmt.Sprintf("读取本地图片失败了，错误信息：%v", err))
		return
	}
	m := message.Image("base64://" + imgB64)
	if id := ctx.Send(m).ID(); id == 0 {
		ctx.Send("ERROR: 可能被风控或读取图片用时过长，请耐心等待")
	}
}

func handleImgsaveHelp(ctx *zero.Ctx) {
	ctx.Send(helpInfo)
}

func handleRepoList(ctx *zero.Ctx) {
	var buffer bytes.Buffer
	for _, list := range config.Repos {
		buffer.WriteString(strings.Join(list, " "))
		buffer.WriteString("\n")
	}
	ctx.Send(buffer.String())
}

func handleRepoCreate(ctx *zero.Ctx) {
	data := ctx.State["args"].(string)
	if data == "" {
		ctx.Send("请提供图库名，可设置别称，需要用空格隔开。\n例：创建图库 高坂穗乃果 高坂穂乃果 穗乃果 穂乃果 ほのか honoka")
		return
	}
	repoNameList := strings.Split(data, " ")
	for _, name := range repoNameList {
		if target := toTargetDirName(name); target != "" {
			ctx.Send(fmt.Sprintf("已存在图库「%s」使用了相同的名称「%s」", target, name))
			return
		}
	}
	config.Repos = append(config.Repos, repoNameList)
	updateConfigToDisk()
	successMsg := fmt.Sprintf("成功创建图库「%s」", data)
	ctx.Send(successMsg)
}

func handleBlacklistAdd(ctx *zero.Ctx) {
	data := ctx.State["args"].(string)
	uid, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		ctx.Send(fmt.Sprintf("解析失败，「%s」不是准确的 QQ 号。", data))
		return
	}
	if inBlacklist(uid) {
		ctx.Send(fmt.Sprintf("添加失败，「%d」已经在黑名单里了。", uid))
		return
	}
	config.Blacklist = append(config.Blacklist, uid)
	updateConfigToDisk()
	ctx.Send(fmt.Sprintf("成功添加「%d」到图库黑名单。", uid))
}

func handleBlacklistDel(ctx *zero.Ctx) {
	data := ctx.State["args"].(string)
	ctx.Send(data)
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

func toTargetDirName(s string) string {
	for _, dir := range config.Repos {
		for _, name := range dir {
			if s == name {
				return dir[0]
			}
		}
	}
	return ""
}

func getLocalImage(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Errorln("[setu]", fmt.Sprintf("Fail to read dir %s", dir))
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no image found in dir %s", dir)
	}
	imgFile := files[rand.Intn(len(files))]
	// 检测是否读到文件夹，如果是则重试三次，否则报错
	for i := 0; i < 3 && imgFile.IsDir(); i++ {
		imgFile = files[rand.Intn(len(files))]
	}
	if imgFile.IsDir() {
		return "", fmt.Errorf("fail to get a file in dir %s", dir)
	}
	imgBytes, err := os.ReadFile(path.Join(dir, imgFile.Name()))
	if err != nil {
		log.Errorln("[setu]", fmt.Sprintf("Fail to read img file %s", imgFile.Name()), err)
		return "", err
	}
	return b64.StdEncoding.EncodeToString(imgBytes), nil
}

func inBlacklist(uid int64) bool {
	for _, id := range config.Blacklist {
		if uid == id {
			return true
		}
	}
	return false
}

func initConfig() {
	confData, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Panicln("[imgsace]", "fail to read config file", err)
		return
	}
	if err := yaml.Unmarshal(confData, &config); err != nil {
		log.Panicln("[imgsave]", "Fail to unmarshal config data", err)
		return
	}
	log.Infoln("[imgsave]", config)
}

func updateConfigToDisk() {
	file, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		log.Panicln("[imgsave]", "fail to open config file", err)
		return
	}
	defer file.Close()
	enc := yaml.NewEncoder(file)
	err = enc.Encode(config)
	if err != nil {
		log.Panicln("[imgsave]", "fail to encoding config", err)
		return
	}
}

func initTargetDir() {
	for _, dir := range config.Repos {
		dirPath := path.Join(config.Rootdit, dir[0])
		_, err := os.Stat(dirPath)
		if os.IsNotExist(err) {
			os.MkdirAll(dirPath, 0755)
		}
	}
}
