package main

import (
	l "./logger"
	"bufio"
	"fmt"
	"os"
	s "strings"
)

// alias
var p = fmt.Println

func resp(str string, logger *l.Logger) {
	p(str)
	logger.Res(str)
}

func main() {
	// 独自のLoggerを使用
	var logger *l.Logger = l.GetLogger()
	defer logger.Close()

	// temp logic
	i := 0

	// 将棋所とのやりとり
	// TODO:いつでも返答すべきコマンドは常時listenするイメージで。
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		logger.Req(text)
		switch text {
		// エンジン登録時は、usiとquitのみ入力される。
		case "usi":
			// TODO: 切り出す
			resp("id name shogi01 0.0.1", logger)
			resp("id author 32hiko", logger)
			resp("usiok", logger)
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
