<div align="center">
  <img src="img/yukichan.svg" alt="ゆき酱" width = "200">
  <br>

  <h1>ゆき酱</h1>

  [![GitHub](https://img.shields.io/github/license/aimerneige/yukichan-bot)](https://raw.githubusercontent.com/aimerneige/yukichan-bot/main/LICENSE)

  ゆき酱是使用 [ZeroBot](https://github.com/wdvxdr1123/ZeroBot) 构建的「**国际象棋**」聊天机器人。

  <img src="https://counter.seku.su/cmoe?name=YukiChan-Bot&theme=r34" /><br>

</div>

> 本机器人主要用于自用，开发过程中没有考虑通用性，按照个人喜好加了很多彩蛋和私货，且部分功能参考了社区内的其他机器人，如果您想要一个通用且功能更加完善的机器人，推荐查阅 [FloatTech/ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin)，本仓库的原创插件也会尽量同步更新到这个仓库。

## Star History

[![Star Trend](https://api.star-history.com/svg?repos=aimerneige/yukichan-bot&type=Timeline)](https://seladb.github.io/StarTrack-js/#/preload?r=aimerneige,yukichan-bot)

## 如何编译

本项目使用 Makefile 管理编译流程，使用如下指令即可快速编译可执行文件：

```bash
make build
```

使用如下指令快速运行并测试程序：

```bash
make run
```

更多信息请查阅 `Makefile`

## 如何使用

本项目符合 [OneBot](https://github.com/howmanybots/onebot) 标准，可基于以下项目与机器人框架/平台进行交互
| 项目地址                                                                    | 平台                                          | 核心作者       |
| --------------------------------------------------------------------------- | --------------------------------------------- | -------------- |
| [go-cqhttp](https://github.com/Mrs4s/go-cqhttp)                             | [MiraiGo](https://github.com/Mrs4s/MiraiGo)   | Mrs4s          |
| [onebot-kotlin](https://github.com/yyuueexxiinngg/onebot-kotlin)            | [Mirai](https://github.com/mamoe/mirai)       | yyuueexxiinngg |
| [oicq/http-api](https://github.com/takayama-lily/oicq/tree/master/http-api) | [OICQ](https://github.com/takayama-lily/oicq) | takayama       |

## 如何部署

```bash
# 安装 GNU Make
sudo apt install -y make
# 安装 Go 1.21
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
# 安装 Python
sudo apt install -y python-is-python3 python3-pip	
# 安装 pip 包
pip install python-chess
git clone https://github.com/dn1z/pgn2gif.git && cd pgn2gif && sudo python setup.py install
# 安装 yukichan-bot
git clone https://github.com/aimerneige/yukichan-bot
# 编译 yukichan-bot
cd yukichan-bot && make build
# 安装 Inkscape 可选
wget https://inkscape.org/gallery/item/42330/Inkscape-0e150ed-x86_64.AppImage
chmod +x Inkscape-0e150ed-x86_64.AppImage
sudo mv Inkscape-0e150ed-x86_64.AppImage /usr/bin/inkscape
sudo add-apt-repository universe
sudo apt install -y libfuse2
# 启动项目
make run
```

## 插件及用法

<details>
<summary>点击展开查看插件及其用法</summary>

<details><summary>✅ alipay 支付宝到账语音生成</summary>

- 支付宝到账 114514

</details>
<details><summary>✅ bilibili 哔哩哔哩相关功能</summary>

> 解析群内 bilibili 链接

</details>
<details><summary>✅ blacklist 黑名单</summary>

> 拒绝为被加入黑名单的用户提供服务

</details>
<details><summary>✅ chess 国际象棋</summary>

> 群内发送「**帮助**」或「**help**」查看详细使用帮助

</details>
<details><summary>✅ donate 捐赠二维码</summary>

- /donate
- /捐赠

</details>
<details><summary>✅ fadian 每日发癫</summary>

- 每日发癫 小乌贼

</details>
<details><summary>✅ fortune 求签</summary>

- 求签 代码无 bug

> 注：机器人不会变卦

</details>
<details><summary>✅ github GitHub 仓库信息</summary>

> 群内接收到 GitHub 仓库链接时自动解析并发送仓库信息的图片

</details>
<details><summary>✅ manager 简易群管</summary>

> 群内发送「**群管帮助**」查看详细使用帮助

</details>
<details><summary>✅ music 点歌</summary>

- 点歌 My Dearest

</details>
<details><summary>✅ random 随机事件生成器</summary>

- /coin
- 掷硬币
- /dice
- 掷骰子

</details>
<details><summary>✅ read60s 每天 60 秒读懂世界</summary>

- 60s
- 早报
- 今日新闻

</details>
<details><summary>✅ setu 色图</summary>

- /setu

> 注：不公开的服务

</details>
<details><summary>✅ suangua 算卦</summary>

- 算卦 代码无 bug

> 注：机器人不会变卦

</details>
<details><summary>✅ tarot 塔罗牌</summary>

- 塔罗
- 今日运势
- 塔罗占卜
- 抽塔罗牌 3

</details>
<details><summary>✅ waifu 随机 AI 老婆</summary>

- /waifu

</details>
<details><summary>✅ wangyiyun 网易云热评</summary>

- 来份网易云热评
- /wyy

</details>
</details>

## 常见问题

### 是否会支持群内多盘对局同时进行

每个群内同时只能存在一盘对局，如果有多盘对局同时进行的需求可以 fork 之后自己改。\
本项目主要是希望提供一个在群内下棋的环境，重要的是大家一起围观、交流和讨论棋局，而不是单纯实现对局。太多的对局同时进行不仅会导致群消息过多炸群，而且也不利于交流。如果只是需要下棋，chess.com 的邀请链接完全可以创建无限的棋局。<sub>~绝对不是开发者懒得写~</sub>

## 交流群

点击链接或扫码加入 QQ 群:

[857066811](https://qm.qq.com/cgi-bin/qm/qr?k=rMtw1SlmoFOp08i5Zw5bM361ljIyzVA-&authKey=9OUzro5oH5CnnFaAbIMwa60987+8ZMwu5GvUAlFUzDIQKVL91z9zUhWp6m1Kayf8&noverify=0)

![qrcode 857066811](img/qr-code.png)

## LICENSE

<a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
<img src="https://www.gnu.org/graphics/agplv3-155x51.png">
</a>

本项目使用 `AGPLv3` 协议开源，您可以在 [GitHub](https://github.com/aimerneige/yukichan-bot) 获取本项目源代码。为了整个社区的良性发展，我们强烈建议您做到以下几点：

- **间接接触（包括但不限于使用 `Http API` 或 跨进程技术）到本项目的软件使用 `AGPLv3` 开源**
- **不鼓励，不支持一切商业使用**

## Open Source

- [wdvxdr1123/ZeroBot](https://github.com/wdvxdr1123/ZeroBot)
- [FloatTech/ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin)
