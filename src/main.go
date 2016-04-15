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
const PROGRAM_NAME = "HoneyWaffle"
const PROGRAM_VERSION = "0.0.3"
const AUTHOR = "Mitsuhiko Watanabe"

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
	//player := NewPlayer("Slide")
	//player := NewPlayer("Random")
	player := NewPlayer("Main")

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
				logger.Trace(text)
				split_text := s.Split(text, " ")
				// 通常の対局
				// position startpos moves 7g7f 8b7b 2g2f
				is_sfen := false
				if split_text[1] == "sfen" {
					is_sfen = true
					// 局面編集からの検討だとこのように
					// position sfen lnsgkgsnl/1r5b1/1pppppppp/p8/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1 moves 2g2f
					sfen_index := s.Index(text, "sfen")
					moves_index := s.Index(text, "moves")
					var sfen_str string
					if moves_index == -1 {
						sfen_str = text[sfen_index:]
					} else {
						sfen_str = text[sfen_index : moves_index-1]
					}
					master = FromSFEN(sfen_str)
				}
				if is_sfen {
					// こちらのルートはどうすればいいのか不明。デッドコピーとしておく。
					for index, value := range split_text {
						if index < 7 {
							continue
						}
						// 何度も一手ずつ反映する必要はないので、スキップしている。
						if index-7 < tesuu {
							continue
						}
						logger.Trace("to apply: " + value)
						master.ApplyMove(value)
						logger.Trace(master.Display())
						tesuu++
					}
					// resp("info string "+text, logger)
				} else {
					for index, value := range split_text {
						if index < 3 {
							continue
						}
						// 何度も一手ずつ反映する必要はないので、スキップしている。
						if index-3 < tesuu {
							continue
						}
						logger.Trace("to apply: " + value)
						master.ApplyMove(value)
						logger.Trace(master.Display())
						logger.Trace(master.ToSFEN())
						tesuu++
					}
					// resp("info string "+text, logger)
				}
			} else if s.HasPrefix(text, "go") {
				bestmove := player.Search(master)
				if len(bestmove) < 6 {
					master.ApplyMove(bestmove)
					if s.Index(bestmove, "*") == 1 {
						// 打つ手は、先後問わず駒の種類を大文字で返す仕様。
						from := bestmove[0:2]
						to := bestmove[2:]
						bestmove = s.ToUpper(from) + to
					}
					logger.Trace(master.Display())
					logger.Trace(master.ToSFEN())
					tesuu++
				}
				resp("info time 0 depth 1 nodes 1 score cp 28 pv "+bestmove, logger)
				bestmove_str := "bestmove " + bestmove
				resp(bestmove_str, logger)
			}
		}
	}
}
