package ytdlp

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/aimerneige/yukichan-bot/internal/pkg/common"
	"github.com/aimerneige/yukichan-bot/internal/pkg/utils"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const tempFileDir = "./temp/ytdlp/"

func init() {
	engine := zero.New()
	common.DefaultSingle.Apply(engine)
	engine.OnPrefix("/yt-dlp", zero.OnlyGroup, zero.SuperUserPermission).
		SetPriority(6).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			wd, err := os.Getwd()
			if err != nil {
				ctx.Send("系统错误，无法获取当前程序目录。")
			}
			if !utils.CommandExists("yt-dlp") {
				ctx.Send("目标服务器未安装 yt-dlp。")
				return
			}
			url := ctx.State["args"].(string)
			cmdVideoFileSize := exec.Command("yt-dlp", "--print", "%(filesize,filesize_approx)s", url)
			videoFileSizeInByte, err := cmdVideoFileSize.Output()
			if err != nil {
				ctx.Send("获取视频文件大小失败，可能是服务器无法访问 YouTube。")
				return
			}
			videoFileSize := string(videoFileSizeInByte)
			cmdVideoTitle := exec.Command("yt-dlp", "--print", "%(title)s", url)
			videoTitleInByte, err := cmdVideoTitle.Output()
			if err != nil {
				ctx.Send("获取视频标题失败，可能是服务器无法访问 YouTube。")
				return
			}
			videoTitle := string(videoTitleInByte)
			ctx.Send(fmt.Sprintf("视频标题：%s视频大小：%s即将开始下载视频，请稍候。", videoTitle, videoFileSize))
			fileName := fmt.Sprintf("%d_%d.mp4", ctx.Event.Sender.ID, time.Now().Unix())
			videoFilePath := path.Join(tempFileDir, fileName)
			cmdDownload := exec.Command("yt-dlp", url, "-o", videoFilePath)
			downloadLogInByte, err := cmdDownload.Output()
			downloadLog := string(downloadLogInByte)
			if err != nil {
				ctx.Send(fmt.Sprintf("视频文件下载失败，可能存在网络波动。\n%s", downloadLog))
				return
			}
			ctx.Send("文件下载成功，正在上传，请稍候。")
			fullVideoFilePath := path.Join(wd, videoFilePath)
			ctx.UploadThisGroupFile(fullVideoFilePath, fileName, "")
			os.Remove(videoFilePath)
		})
}
