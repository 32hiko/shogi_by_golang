package main

import (
	. "logger"
	. "shogi"
	"bufio"
	"fmt"
	"os"
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

func develop(logger *Logger) bool {
	state := CreateInitialState()
	state_str := state.Display()
	p(state_str)
	logger.Trace(state_str)
	resp("ok!", logger)
	return true
}

func main() {
	// 独自のLoggerを使用
	InitLogger()
	logger := GetLogger()
	defer logger.Close()

	// develp mode
	if develop(logger) {
		os.Exit(0)
	}

	// temp logic
	i := 0

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
			resp("readyok", logger)
		case "usinewgame":
			// TODO: モードを切り替えるべきか。
		case "gameover":
			// TODO: 対局待ち状態に戻る。
		default:
			if s.HasPrefix(text, "position") {
				// TODO: 盤面を更新する
				resp("info string "+text, logger)
			} else if s.HasPrefix(text, "go") {
				// TODO: ここで思考し、手を返す。以下は飛車を動かすだけの暫定ロジック。
				resp("info string "+text, logger)
				if i%2 == 0 {
					resp("bestmove 8b7b", logger)
				} else {
					resp("bestmove 7b8b", logger)
				}
				i++
			}
		}
	}
}
