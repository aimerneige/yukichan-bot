package openai

import (
	"fmt"
	"os"

	"github.com/aimerneige/yukichan-bot/internal/pkg/common"
	"github.com/aimerneige/yukichan-bot/internal/plugin/openai/query"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"gopkg.in/yaml.v3"
)

const configFilePath = "./config/openai.yaml"

type OpenaiConfig struct {
	Secret OpenaiSecret `yaml:"secret"`
}

type OpenaiSecret struct {
	Account string `yaml:"account"`
	AppId   string `yaml:"appid"`
	Token   string `yaml:"token"`
	AesKey  string `yaml:"aeskey"`
}

func init() {
	confData, err := os.ReadFile(configFilePath)
	if err != nil {
		logrus.Errorln("[openai]", "Fail to read config file", err)
		return
	}
	var config OpenaiConfig
	if err := yaml.Unmarshal(confData, &config); err != nil {
		logrus.Errorln("[openai]", "Fail to unmarshal config data", err)
		return
	}
	logrus.Debugln(config)

	engine := zero.New()
	engine.OnPrefix("//").
		SetPriority(2).
		SetBlock(true).Handle(func(ctx *zero.Ctx) {
		userQuery := ctx.State["args"].(string)
		secret := config.Secret
		accessToken := query.GetToken(secret.Account, secret.AppId, secret.Token)
		response := query.SendQueryRequest(query.ApiReq{
			Query:    userQuery,
			Env:      "online",
			UserName: ctx.Event.Sender.NickName,
			Avatar:   fmt.Sprintf("https://q2.qlogo.cn/headimg_dl?dst_uin=%d&spec=100", ctx.Event.Sender.ID),
			Userid:   fmt.Sprintf("qq_%d", ctx.Event.Sender.ID),
		}, accessToken, secret.Token, secret.AesKey)
		ctx.Send(response)
	})
	engine.UseMidHandler(common.DefaultSpeedLimit)
}
