package service

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/aimerneige/yukichan-bot/internal/config"
	"github.com/aimerneige/yukichan-bot/internal/plugin/chess/database"
	"github.com/notnil/chess"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
)

//go:embed scripts/board2svg.py
var pythonScriptBoard2SVG string

//go:embed scripts/pgn2gif.py
var pythonScriptPGN2GIF string

var instance *chessService

const (
	eloDefault   = 500
	inkscapePath = "./bin/inkscape"
)

type chessService struct {
	gameRooms map[int64]chessRoom
}

type chessRoom struct {
	chessGame    *chess.Game
	whitePlayer  int64
	whiteName    string
	blackPlayer  int64
	blackName    string
	drawPlayer   int64
	lastMoveTime int64
	isBlindfold  bool
	whiteErr     bool // 违例记录（盲棋用）
	blackErr     bool
}

func init() {
	instance = &chessService{
		gameRooms: make(map[int64]chessRoom, 1),
	}
}

// Game 下棋
func Game(groupCode, senderUin int64, senderName string) message.Message {
	return createGame(false, groupCode, senderUin, senderName)
}

// Blindfold 盲棋
func Blindfold(groupCode, senderUin int64, senderName string) message.Message {
	return createGame(true, groupCode, senderUin, senderName)
}

// Abort 中断对局
func Abort(groupCode int64) message.Message {
	if _, ok := instance.gameRooms[groupCode]; ok {
		return abortGame(groupCode, "对局已被管理员中断，游戏结束。")
	}
	return simpleText("对局不存在，发送「下棋」或「chess」可创建对局。")
}

// Draw 和棋
func Draw(groupCode, senderUin int64) message.Message {
	if room, ok := instance.gameRooms[groupCode]; ok {
		if senderUin == room.whitePlayer || senderUin == room.blackPlayer {
			room.lastMoveTime = time.Now().Unix()
			if room.drawPlayer == 0 {
				room.drawPlayer = senderUin
				instance.gameRooms[groupCode] = room
				return textWithAt(senderUin, "请求和棋，发送「和棋」或「draw」接受和棋。走棋视为拒绝和棋。")
			}
			if room.drawPlayer == senderUin {
				return textWithAt(senderUin, "已发起和棋请求，请勿重复发送。")
			}
			room.chessGame.Draw(chess.DrawOffer)
			chessString := getChessString(room)
			eloString := ""
			if len(room.chessGame.Moves()) > 4 {
				dbService := NewDBService(database.GetDB())
				if err := dbService.CreatePGN(chessString, room.whitePlayer, room.blackPlayer, room.whiteName, room.blackName); err != nil {
					log.Errorln("[chess]", "Fail to create PGN.", err)
				}
				whiteScore, blackScore := 0.5, 0.5
				elo, err := getELOString(room, whiteScore, blackScore)
				if err != nil {
					log.Errorln("[chess]", "Fail to get eloString.", eloString, err)
				}
				eloString = elo
			}
			replyMsg := textWithAt(senderUin, "接受和棋，游戏结束。\n"+eloString+chessString)
			// TODO 替换 GIF 生成库
			// gif, msg, err := generateGIF(groupCode, chessString)
			// if err != nil {
			// 	log.Errorln("[chess]", "Fail to generate GIF.", err)
			// 	replyMsg = append(replyMsg, message.Text("\n\n[GIF - ERROR]\n"+msg))
			// } else {
			// 	replyMsg = append(replyMsg, message.Image("base64://"+gif))
			// }
			delete(instance.gameRooms, groupCode)
			return replyMsg
		}
		return textWithAt(senderUin, "不是对局中的玩家，无法请求和棋。")
	}
	return simpleText("对局不存在，发送「下棋」或「chess」可创建对局。")
}

// Resign 认输
func Resign(groupCode, senderUin int64) message.Message {
	if room, ok := instance.gameRooms[groupCode]; ok {
		// 检查是否是当前游戏玩家
		if senderUin == room.whitePlayer || senderUin == room.blackPlayer {
			// 如果对局未建立，中断对局
			if room.whitePlayer == 0 || room.blackPlayer == 0 {
				delete(instance.gameRooms, groupCode)
				return simpleText("对局已释放。")
			}
			var resignColor chess.Color
			if senderUin == room.whitePlayer {
				resignColor = chess.White
			} else {
				resignColor = chess.Black
			}
			if isAprilFoolsDay() {
				if resignColor == chess.White {
					resignColor = chess.Black
				} else {
					resignColor = chess.White
				}
			}
			room.chessGame.Resign(resignColor)
			chessString := getChessString(room)
			eloString := ""
			if len(room.chessGame.Moves()) > 4 {
				dbService := NewDBService(database.GetDB())
				if err := dbService.CreatePGN(chessString, room.whitePlayer, room.blackPlayer, room.whiteName, room.blackName); err != nil {
					log.Errorln("[chess]", "Fail to create PGN.", err)
				}
				whiteScore, blackScore := 1.0, 1.0
				if resignColor == chess.White {
					whiteScore = 0.0
				} else {
					blackScore = 0.0
				}
				elo, err := getELOString(room, whiteScore, blackScore)
				if err != nil {
					log.Errorln("[chess]", "Fail to get eloString.", eloString, err)
				}
				eloString = elo
			}
			replyMsg := textWithAt(senderUin, "认输，游戏结束。\n"+eloString+chessString)
			if isAprilFoolsDay() {
				replyMsg = textWithAt(senderUin, "对手认输，游戏结束，你胜利了。\n"+eloString+chessString)
			}
			// TODO 替换 GIF 生成库
			// gif, msg, err := generateGIF(groupCode, chessString)
			// if err != nil {
			// 	log.Errorln("[chess]", "Fail to generate GIF.", err)
			// 	replyMsg = append(replyMsg, message.Text("\n\n[GIF - ERROR]\n"+msg))
			// } else {
			// 	replyMsg = append(replyMsg, message.Image("base64://"+gif))
			// }
			delete(instance.gameRooms, groupCode)
			return replyMsg
		}
		return textWithAt(senderUin, "不是对局中的玩家，无法认输。")
	}
	return simpleText("对局不存在，发送「下棋」或「chess」可创建对局。")
}

// Play 走棋
func Play(senderUin int64, groupCode int64, moveStr string) message.Message {
	if room, ok := instance.gameRooms[groupCode]; ok {
		// 不是对局中的玩家，忽略消息
		if (senderUin != room.whitePlayer) && (senderUin != room.blackPlayer) && !isAprilFoolsDay() {
			return nil
		}
		// 对局未建立
		if (room.whitePlayer == 0) || (room.blackPlayer == 0) {
			return textWithAt(senderUin, "请等候其他玩家加入游戏。")
		}
		// 需要对手走棋
		if ((senderUin == room.whitePlayer) && (room.chessGame.Position().Turn() != chess.White)) || ((senderUin == room.blackPlayer) && (room.chessGame.Position().Turn() != chess.Black)) {
			return textWithAt(senderUin, "请等待对手走棋。")
		}
		room.lastMoveTime = time.Now().Unix()
		// 走棋
		if err := room.chessGame.MoveStr(moveStr); err != nil {
			// 指令错误时检查
			if !room.isBlindfold {
				// 未开启盲棋，提示指令错误
				return simpleText(fmt.Sprintf("移动「%s」违规，请检查，格式请参考「代数记谱法」(Algebraic notation)。", moveStr))
			}
			// 开启盲棋，判断违例情况
			var currentPlayerColor chess.Color
			if senderUin == room.whitePlayer {
				currentPlayerColor = chess.White
			} else {
				currentPlayerColor = chess.Black
			}
			// 第一次违例，提示
			_flag := false
			if (currentPlayerColor == chess.White) && !room.whiteErr {
				room.whiteErr = true
				instance.gameRooms[groupCode] = room
				_flag = true
			}
			if (currentPlayerColor == chess.Black) && !room.blackErr {
				room.blackErr = true
				instance.gameRooms[groupCode] = room
				_flag = true
			}
			if _flag {
				return simpleText(fmt.Sprintf("移动「%s」违例，再次违例会立即判负。", moveStr))
			}
			// 出现多次违例，判负
			room.chessGame.Resign(currentPlayerColor)
			chessString := getChessString(room)
			replyMsg := textWithAt(senderUin, "违例两次，游戏结束。\n"+chessString)
			gif, msg, err := generateGIF(groupCode, chessString)
			if err != nil {
				log.Errorln("[chess]", "Fail to generate GIF.", err)
				replyMsg = append(replyMsg, message.Text("\n\n[GIF - ERROR]\n"+msg))
			} else {
				replyMsg = append(replyMsg, message.Image("base64://"+gif))
			}
			delete(instance.gameRooms, groupCode)
			return replyMsg
		}
		// 走子之后，视为拒绝和棋
		if room.drawPlayer != 0 {
			room.drawPlayer = 0
			instance.gameRooms[groupCode] = room
		}
		// 生成棋盘图片
		var boardImgEle message.MessageSegment
		if !room.isBlindfold {
			boardImgB64, ok, errMsg := getBoardElement(groupCode)
			boardImgEle = message.Image("base64://" + boardImgB64)
			if !ok {
				return errorText(errMsg)
			}
		}
		// 检查游戏是否结束
		if room.chessGame.Method() != chess.NoMethod {
			whiteScore, blackScore := 0.5, 0.5
			msg := "游戏结束，"
			switch room.chessGame.Method() {
			case chess.FivefoldRepetition:
				msg += "和棋，因为五次重复走子。\n"
			case chess.SeventyFiveMoveRule:
				msg += "和棋，因为七十五步规则。\n"
			case chess.InsufficientMaterial:
				msg += "和棋，因为不可能将死。\n"
			case chess.Stalemate:
				msg += "和棋，因为逼和（无子可动和棋）。\n"
			case chess.Checkmate:
				var winner string
				if room.chessGame.Position().Turn() == chess.White {
					whiteScore = 0.0
					blackScore = 1.0
					winner = "黑方"
				} else {
					whiteScore = 1.0
					blackScore = 0.0
					winner = "白方"
				}
				msg += winner
				msg += "胜利，因为将杀。\n"
			}
			chessString := getChessString(room)
			eloString := ""
			// 若走子次数超过 4 认为是有效对局，存入数据库
			if len(room.chessGame.Moves()) > 4 {
				dbService := NewDBService(database.GetDB())
				if err := dbService.CreatePGN(chessString, room.whitePlayer, room.blackPlayer, room.whiteName, room.blackName); err != nil {
					log.Errorln("[chess]", "Fail to create PGN.", err)
				}
				// 仅有效对局才会计算等级分
				elo, err := getELOString(room, whiteScore, blackScore)
				if err != nil {
					log.Errorln("[chess]", "Fail to get eloString.", eloString, err)
				}
				eloString = elo
			}
			replyMsg := simpleText(msg + eloString + chessString)
			if !room.isBlindfold {
				replyMsg = append(replyMsg, boardImgEle)
			}
			// TODO 替换 GIF 生成库
			// gif, msg, err := generateGIF(groupCode, chessString)
			// if err != nil {
			// 	log.Errorln("[chess]", "Fail to generate GIF.", err)
			// 	replyMsg = append(replyMsg, message.Text("\n\n[GIF - ERROR]\n"+msg))
			// } else {
			// 	replyMsg = append(replyMsg, message.Image("base64://"+gif))
			// }
			delete(instance.gameRooms, groupCode)
			return replyMsg
		}
		// 提示玩家继续游戏
		var currentPlayer int64
		if room.chessGame.Position().Turn() == chess.White {
			currentPlayer = room.whitePlayer
		} else {
			currentPlayer = room.blackPlayer
		}
		return append(textWithAt(currentPlayer, "对手已走子，游戏继续。"), boardImgEle)
	}
	return textWithAt(senderUin, "对局不存在，发送「下棋」或「chess」可创建对局。")
}

// Ranking 排行榜
func Ranking() message.Message {
	ranking, err := getRankingString()
	if err != nil {
		log.Errorln("[chess]", "Fail to get player ranking.", err)
		return simpleText("服务器错误，无法获取排行榜信息。请联系开发者修 bug。\n反馈地址 https://github.com/aimerneige/yukichan-bot/issues\n")
	}
	return simpleText(ranking)
}

// Rate 获取等级分
func Rate(senderUin int64, senderName string) message.Message {
	dbService := NewDBService(database.GetDB())
	rate, err := dbService.GetELORateByUin(senderUin)
	if err == gorm.ErrRecordNotFound {
		return simpleText("没有查找到等级分信息。请至少进行一局对局。")
	}
	if err != nil {
		logrus.Errorln("[chess]", "Fail to get player rank.", err)
		return simpleText("服务器错误，无法获取等级分信息。请联系开发者修 bug。\n反馈地址 https://github.com/aimerneige/yukichan-bot/issues\n")
	}
	return simpleText(fmt.Sprintf("玩家「%s」目前的等级分：%d", senderName, rate))
}

// createGame 创建游戏
func createGame(isBlindfold bool, groupCode int64, senderUin int64, senderName string) message.Message {
	if room, ok := instance.gameRooms[groupCode]; ok {
		if room.blackPlayer != 0 {
			// 检测对局是否已存在超过 6 小时
			if (time.Now().Unix() - room.lastMoveTime) > 21600 {
				autoAbortMsg := abortGame(groupCode, "对局已存在超过 6 小时，游戏结束。")
				autoAbortMsg = append(autoAbortMsg, message.Text("\n\n已有对局已被中断，如需创建新对局请重新发送指令。"))
				autoAbortMsg = append(autoAbortMsg, message.At(senderUin))
				return autoAbortMsg
			}
			// 对局在进行
			msg := textWithAt(senderUin, "对局已在进行中，无法创建或加入对局，当前对局玩家为：")
			if room.whitePlayer != 0 {
				msg = append(msg, message.At(room.whitePlayer))
			}
			if room.blackPlayer != 0 {
				msg = append(msg, message.At(room.blackPlayer))
			}
			msg = append(msg, message.Text("，群主或管理员发送「中断」或「abort」可中断对局（自动判和）。"))
			return msg
		}
		if senderUin == room.whitePlayer {
			return textWithAt(senderUin, "请等候其他玩家加入游戏。")
		}
		if room.isBlindfold && !isBlindfold {
			return simpleText("已创建盲棋对局，请加入或等待盲棋对局结束之后创建普通对局。")
		}
		if !room.isBlindfold && isBlindfold {
			return simpleText("已创建普通对局，请加入或等待普通对局结束之后创建盲棋对局。")
		}
		room.blackPlayer = senderUin
		room.blackName = senderName
		instance.gameRooms[groupCode] = room
		var boardImgEle message.MessageSegment
		if !room.isBlindfold {
			boardImgB64, ok, errMsg := getBoardElement(groupCode)
			if !ok {
				return errorText(errMsg)
			}
			boardImgEle = message.Image("base64://" + boardImgB64)
		}
		if isBlindfold {
			return append(simpleText("黑棋已加入对局，请白方下棋。"), message.At(room.whitePlayer))
		}
		return append(simpleText("黑棋已加入对局，请白方下棋。"), message.At(room.whitePlayer), boardImgEle)
	}
	instance.gameRooms[groupCode] = chessRoom{
		chessGame:    chess.NewGame(),
		whitePlayer:  senderUin,
		whiteName:    senderName,
		blackPlayer:  0,
		blackName:    "",
		drawPlayer:   0,
		lastMoveTime: time.Now().Unix(),
		isBlindfold:  isBlindfold,
		whiteErr:     false,
		blackErr:     false,
	}
	if isBlindfold {
		return simpleText("已创建新的盲棋对局，发送「盲棋」或「blind」可加入对局。")
	}
	return simpleText("已创建新的对局，发送「下棋」或「chess」可加入对局。")
}

// abortGame 中断游戏
func abortGame(groupCode int64, hint string) message.Message {
	room := instance.gameRooms[groupCode]
	room.chessGame.Draw(chess.DrawOffer)
	chessString := getChessString(room)
	if len(room.chessGame.Moves()) > 4 {
		dbService := NewDBService(database.GetDB())
		if err := dbService.CreatePGN(chessString, room.whitePlayer, room.blackPlayer, room.whiteName, room.blackName); err != nil {
			log.Errorln("[chess]", "Fail to create PGN.", err)
		}
	}
	delete(instance.gameRooms, groupCode)
	msg := simpleText(hint)
	if room.whitePlayer != 0 {
		msg = append(msg, message.At(room.whitePlayer))
	}
	if room.blackPlayer != 0 {
		msg = append(msg, message.At(room.blackPlayer))
	}
	msg = append(msg, message.Text("\n\n"+chessString))
	return msg
}

// getBoardElement 获取棋盘图片的消息内容
func getBoardElement(groupCode int64) (string, bool, string) {
	if room, ok := instance.gameRooms[groupCode]; ok {
		var uciStr string
		// 将最后一步走子转化为 uci 字符串
		moves := room.chessGame.Moves()
		if len(moves) != 0 {
			uciStr = moves[len(moves)-1].String()
		} else {
			uciStr = "None"
		}
		tempFileDir := config.Conf.TempDir["chess"]
		svgFilePath := path.Join(tempFileDir, fmt.Sprintf("%d.svg", groupCode))
		pngFilePath := path.Join(tempFileDir, fmt.Sprintf("%d.png", groupCode))
		// 调用 python 脚本生成 svg 文件
		if err := exec.Command("python", "-c", pythonScriptBoard2SVG, room.chessGame.FEN(), svgFilePath, uciStr).Run(); err != nil {
			log.Debugln("[chess]", "python", " ", "-c", " ", "python_script_board2svg", " ", room.chessGame.FEN(), " ", svgFilePath, " ", uciStr)
			log.Errorln("[chess]", "Unable to generate svg file.", err)
			return "", false, "无法生成 svg 图片，请检查 python-chess 库是否被正确安装。"
		}
		// 调用 inkscape 将 svg 图片转化为 png 图片
		if err := exec.Command(inkscapePath, "-w", "720", "-h", "720", svgFilePath, "-o", pngFilePath).Run(); err != nil {
			log.Errorln("[chess]", "Unable to convert to png.", err)
			return "", false, "无法生成 png 图片，请检查 inkscape 安装情况及其依赖 libfuse。"
		}
		// 尝试读取 png 图片
		imgData, err := os.ReadFile(pngFilePath)
		if err != nil {
			log.Errorln("[chess]", fmt.Sprintf("Unable to read image file in %s.", pngFilePath), err)
			return "", false, "无法读取 png 图片"
		}
		imgB64 := base64.StdEncoding.EncodeToString(imgData)
		return imgB64, true, ""
	}

	log.Debugln(fmt.Sprintf("No room for groupCode %d.", groupCode))
	return "", false, "对局不存在"
}

// getELOString 获得玩家等级分的文本内容
func getELOString(room chessRoom, whiteScore, blackScore float64) (string, error) {
	if room.whitePlayer == 0 || room.blackPlayer == 0 {
		return "", nil
	}
	eloString := "玩家等级分：\n"
	dbService := NewDBService(database.GetDB())
	if err := updateELORate(room.whitePlayer, room.blackPlayer, room.whiteName, room.blackName, whiteScore, blackScore, dbService); err != nil {
		eloString += "发生错误，无法更新等级分。"
		return eloString, err
	}
	whiteRate, blackRate, err := getELORate(room.whitePlayer, room.blackPlayer, dbService)
	if err != nil {
		eloString += "发生错误，无法获取等级分。"
		return eloString, err
	}
	eloString += fmt.Sprintf("%s：%d\n%s：%d\n\n", room.whiteName, whiteRate, room.blackName, blackRate)
	return eloString, nil
}

// getRankingString 获取等级分排行榜的文本内容
func getRankingString() (string, error) {
	dbService := NewDBService(database.GetDB())
	eloList, err := dbService.GetHighestRateList()
	if err != nil {
		return "", err
	}
	ret := "当前等级分排行榜：\n\n"
	for _, elo := range eloList {
		ret += fmt.Sprintf("%s: %d\n", elo.Name, elo.Rate)
	}
	return ret, nil
}

func simpleText(msg string) message.Message {
	return []message.MessageSegment{message.Text(msg)}
}

func textWithAt(target int64, msg string) message.Message {
	if target == 0 {
		return simpleText("@全体成员 " + msg)
	}
	return []message.MessageSegment{message.At(target), message.Text(msg)}
}

func errorText(errMsg string) message.Message {
	return simpleText("发生错误，请联系开发者修 bug。\n反馈地址 https://github.com/aimerneige/yukichan-bot/issues\n错误信息：" + errMsg)
}

// updateELORate 更新 elo 等级分
// 当数据库中没有玩家的等级分信息时，自动新建一条记录
func updateELORate(whiteUin, blackUin int64, whiteName, blackName string, whiteScore, blackScore float64, dbService *DBService) error {
	whiteRate, err := dbService.GetELORateByUin(whiteUin)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// create white elo
		if err := dbService.CreateELO(whiteUin, whiteName, eloDefault); err != nil {
			return err
		}
		whiteRate = eloDefault
	}
	blackRate, err := dbService.GetELORateByUin(blackUin)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// create black elo
		if err := dbService.CreateELO(blackUin, blackName, eloDefault); err != nil {
			return err
		}
		blackRate = eloDefault
	}
	whiteRate, blackRate = CalculateNewRate(whiteRate, blackRate, whiteScore, blackScore)
	// 更新白棋玩家的 ELO 等级分
	if err := dbService.UpdateELOByUin(whiteUin, whiteName, whiteRate); err != nil {
		return err
	}
	// 更新黑棋玩家的 ELO 等级分
	if err := dbService.UpdateELOByUin(blackUin, blackName, blackRate); err != nil {
		return err
	}

	return nil
}

// generateGIF 使用 PGN 字符串生成对局的 GIF
func generateGIF(groupCode int64, pgnStr string) (string, string, error) {
	tempFileDir := config.Conf.TempDir["chess"]
	if err := exec.Command("python", "-c", pythonScriptPGN2GIF, pgnStr, tempFileDir, fmt.Sprintf("%d", groupCode)).Run(); err != nil {
		log.Debugln("[chess]", "python", " ", "-c", " ", "python_sript_pgn2gif", " ", pgnStr, " ", tempFileDir, " ", fmt.Sprintf("%d", groupCode))
		return "", "生成 gif 时发生错误", err
	}
	gifFilePath := path.Join(tempFileDir, fmt.Sprintf("%d.gif", groupCode))
	gifData, err := os.ReadFile(gifFilePath)
	if err != nil {
		log.Errorln("[chess]", fmt.Sprintf("Unable to read gif file in %s.", gifFilePath), err)
		return "", "读取 gif 时发生错误", err
	}
	gifB64 := base64.StdEncoding.EncodeToString(gifData)
	return gifB64, "", err
}

// getChessString 获取 PGN 字符串
func getChessString(room chessRoom) string {
	game := room.chessGame
	siteString := "[Site \"github.com/aimerneige/MiraiChess\"]\n"
	dataString := fmt.Sprintf("[Date \"%s\"]\n", time.Now().Format("2006-01-02"))
	whiteString := fmt.Sprintf("[White \"%s\"]\n", room.whiteName)
	blackString := fmt.Sprintf("[Black \"%s\"]\n", room.blackName)
	chessString := game.String()

	return siteString + dataString + whiteString + blackString + chessString
}

// getELORate 获取玩家的 ELO 等级分
func getELORate(whiteUin, blackUin int64, dbService *DBService) (whiteRate int, blackRate int, err error) {
	whiteRate, err = dbService.GetELORateByUin(whiteUin)
	if err != nil {
		return
	}
	blackRate, err = dbService.GetELORateByUin(blackUin)
	if err != nil {
		return
	}
	return
}

// isAprilFoolsDay 判断当前时间是否为愚人节期间
func isAprilFoolsDay() bool {
	now := time.Now()
	return now.Month() == 4 && now.Day() == 1
}
