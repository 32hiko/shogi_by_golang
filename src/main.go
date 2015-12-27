package main

import (
	"bufio"
	"fmt"
	. "logger"
	"os"
	. "shogi"
	s "strings"
)

// const
const PROGRAM_NAME = "shogi01"
const PROGRAM_VERSION = "0.0.1"
const AUTHOR = "32hiko"

// alias
var p = fmt.Println

func resp(str string, logger *Logger) {
	p(str)
	logger.Res(str)
}

func respUSI(logger *Logger) {
	resp("id name "+PROGRAM_NAME+" "+PROGRAM_VERSION, logger)
	resp("id author "+AUTHOR, logger)
	resp("usiok", logger)
}

func main() {
	// 独自のLoggerを使用
	InitLogger()
	logger := GetLogger()
	defer logger.Close()

	// master ban
	var master *TBan
	var tesuu int = 0
	player := NewPlayer()

	// 将棋所とのやりとり
	// TODO:いつでも返答すべきコマンドは常時listenするイメージで。GoRoutineとChannelを使えばよさげ
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		logger.Req(text)
		switch text {
		// エンジン登録時は、usiとquitのみ入力される。
		case "usi":
			respUSI(logger)
		case "quit":
			// TODO 終了前処理
			os.Exit(0)
		case "setoption name USI_Ponder value true":
			// TODO 設定を保存する
		case "setoption name USI_Hash value 256":
			// TODO 設定を保存する
		case "isready":
			master = CreateInitialState()
			resp("readyok", logger)
		case "usinewgame":
			// TODO: モードを切り替えるべきか。
		case "gameover":
			// TODO: 対局待ち状態に戻る。
		default:
			if s.HasPrefix(text, "position") {
				// TODO: 盤面を更新する
				split_text := s.Split(text, " ")
				// position startpos moves 7g7f 8b7b 2g2f
				// とりあえずは初期配置は通常で。
				for index, value := range split_text {
					if index < 3 {
						continue
					}
					// 何度も一手ずつ反映する必要はないので、スキップできるようにする。
					if index-3 < tesuu {
						continue
					}
					master.ApplyMove(value)
					logger.Trace(master.Display())
					tesuu++
				}
				resp("info string "+text, logger)
			} else if s.HasPrefix(text, "go") {
				bestmove_str := player.Search(master)
				resp(bestmove_str, logger)
			}
		}
	}
}
