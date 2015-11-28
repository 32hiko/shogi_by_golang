package main

import (
	"bufio"
	"fmt"
	"os"
	s "strings"
)

// alias
var p = fmt.Println

func main() {
	// TODO: logは機能化する（出力ON/OFF切り替え）
	log, _ := os.Create("log")
	defer log.Close()

	// temp logic
	i := 0

	// 将棋所とのやりとり
	// TODO:いつでも返答すべきコマンドは常時listenするイメージで。
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		log.WriteString(text + "\n")
		log.Sync()
		switch text {
		// エンジン登録時は、usiとquitのみ入力される。
		case "usi":
			// TODO: 切り出す
			p("id name shogi01 0.0.1")
			p("id author 32hiko")
			p("usiok")
		case "quit":
			// TODO 終了前処理
			os.Exit(0)
		case "setoption name USI_Ponder value true":
			// TODO 設定を保存する
		case "setoption name USI_Hash value 256":
			// TODO 設定を保存する
		case "isready":
			p("readyok")
		case "usinewgame":
			// TODO: モードを切り替えるべきか。
		case "gameover":
			// TODO: 対局待ち状態に戻る。
		default:
			if s.HasPrefix(text, "position") {
				// TODO: 盤面を更新する
				p("info string " + text)
			} else if s.HasPrefix(text, "go") {
				// TODO: ここで思考し、手を返す。以下は飛車を動かすだけの暫定ロジック。
				p("info string " + text)
				if i%2 == 0 {
					p("bestmove 8b7b")
				} else {
					p("bestmove 7b8b")
				}
				i++
			}
		}
	}
}
