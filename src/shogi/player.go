package shogi

import (
	. "logger"
	"math/rand"
	"time"
)

type IPlayer interface {
	Search(*TBan) string
}

func NewPlayer(name string) IPlayer {
	switch name {
	case "Slide":
		return NewSlidePlayer()
	case "Random":
		return NewRandomPlayer()
	default:
		return nil
	}
}

/*
 * ただ飛車を左右に動かすだけ。しかも後手専用
 */
type TSlidePlayer struct {
	i *int
}

func NewSlidePlayer() *TSlidePlayer {
	var count int = 0
	player := TSlidePlayer{}
	player.i = &count
	return &player
}

func (player TSlidePlayer) Search(ban *TBan) string {
	var te string
	if *(player.i)%2 == 0 {
		te = "8b7b"
	} else {
		te = "7b8b"
	}
	*(player.i)++
	return te
}

/*
 * ランダム指し
 */
type TRandomPlayer struct {
}

func NewRandomPlayer() *TRandomPlayer {
	player := TRandomPlayer{}
	return &player
}

func (player TRandomPlayer) Search(ban *TBan) string {
	logger := GetLogger()
	teban := TTeban(*(ban.Tesuu)%2 == 0)
	// logger.Trace("[RandomPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " +s(teban))
	tegoma := ban.GetTebanKoma(teban)
	all_moves := make(map[byte]*TMove)

	for koma_id, _ := range *tegoma {
		// logger.Trace("[RandomPlayer] koma_id: " + s(koma_id))

		for _, move := range ban.AllMoves[koma_id].Map {
			AddMove(&all_moves, move)
		}

	}
	moves_count := len(all_moves)
	logger.Trace("[RandomPlayer] moves: " + s(moves_count))
	if moves_count == 0 {
		return "resign"
	}
	rand.Seed(time.Now().UnixNano())
	random_index := rand.Intn(len(all_moves))
	random_move := all_moves[byte(random_index)]
	return random_move.GetUSIMoveString()
}
