package main

import (
	"bufio"
	. "logger"
	"os"
	. "shogi"
	"strconv"
	s "strings"
	"time"
)

// const
const PROGRAM_NAME = "HoneyWaffle"
const PROGRAM_VERSION = "2.0.0"
const AUTHOR = "Mitsuhiko Watanabe"

const NW_LAG_MS = 1000
const SAFETY_MS = 3000

var master *TBan = nil
var tesuu int = 0
var logger *Logger = nil
var player IPlayer

func respUSI(logger *Logger) {
	Resp("id name "+PROGRAM_NAME+" "+PROGRAM_VERSION, logger)
	Resp("id author "+AUTHOR, logger)
	Resp("usiok", logger)
}

func main() {
	InitLogger()
	logger = GetLogger()
	defer logger.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		logger.Req(text)

		switch text {
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
			// TODO ここも要見直し
			master = CreateInitialState()
			player = NewPlayer("Main")
			Resp("readyok", logger)
		case "usinewgame":
			// TODO: モードを切り替えるべきか。
		case "gameover":
			// TODO: 対局待ち状態に戻る。
		default:
			if s.HasPrefix(text, "position") {
				handle_pos_text(text)
			} else if s.HasPrefix(text, "go") {
				bestmove_str := handle_go_text(text)
				Resp(bestmove_str, logger)
			}
		}
	}
}

func handle_pos_text(text string) {
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
	tesuu = (*master.Tesuu)
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
			master.ApplyMove(value, true, true, true)
			logger.Trace(master.Display())
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
			master.ApplyMove(value, true, true, true)
			logger.Trace(master.Display())
			logger.Trace(master.ToSFEN(true))
		}
	}
}

func handle_go_text(text string) string {
	split_text := s.Split(text, " ")
	var btime string = "600000"
	var wtime string = "600000"
	if len(split_text) >= 5 {
		btime = split_text[2]
		wtime = split_text[4]
	}
	teban := *(master.Teban)
	var ms int = 0
	if teban {
		ms, _ = strconv.Atoi(btime)
	} else {
		ms, _ = strconv.Atoi(wtime)
	}

	ms -= NW_LAG_MS
	time_manager := time.NewTimer(time.Millisecond * time.Duration(ms))
	player_timer := time.NewTimer(time.Millisecond * time.Duration(ms-SAFETY_MS))
	ch := make(chan string)

	go func() {
		<-time_manager.C
		// TODO select better move before time up
		stop := player_timer.Stop()
		if stop {
			ch <- "resign"
		}
	}()
	go func() {
		bestmove, _ := player.Search(master, ms)
		if len(bestmove) < 6 {
			master.ApplyMove(bestmove, true, true, true)
			logger.Trace(master.Display())
			logger.Trace(master.ToSFEN(true))
		}
		ch <- bestmove
		<-player_timer.C
		time_manager.Stop()
	}()

	res := <-ch
	player_timer.Stop()
	time_manager.Stop()
	return "bestmove " + res
}

// タイマーの部分だけサンプル的に作ってみた。
func handle_go_text_sample(text string) string {
	response := ""
	// go btime 600000 wtime 600000 binc 10000 winc 10000
	split_text := s.Split(text, " ")
	btime := split_text[2]
	wtime := split_text[4]

	// temp logic
	limit_ms, _ := strconv.Atoi(btime)
	think_ms, _ := strconv.Atoi(wtime)
	timer1 := time.NewTimer(time.Millisecond * time.Duration(limit_ms))
	timer2 := time.NewTimer(time.Millisecond * time.Duration(think_ms))
	ch := make(chan string)

	// time manager
	go func() {
		<-timer1.C
		// TODO select better move before time up
		stop2 := timer2.Stop()
		if stop2 {
			ch <- "resign"
		}
	}()
	// search thread
	go func() {
		<-timer2.C
		stop1 := timer1.Stop()
		if stop1 {
			ch <- "7c7d"
		}
	}()
	res := <-ch
	response = "bestmove " + res

	return response
}
