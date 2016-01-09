package shogi

import (
	. "logger"
	"time"
	"math/rand"
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
	logger.Trace("[RandomPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " +s(teban))
	tegoma := ban.GetTebanKoma(teban)
	all_moves := make(map[byte]*TMove)
	// いったんは打つ手を考えない
	for koma_id, koma := range *tegoma {
		logger.Trace("[RandomPlayer] koma_id: " + s(koma_id))
		if koma.Position != Mochigoma {
			masu := ban.AllMasu[koma.Position]
			for _, move := range *(masu.Moves) {
				AddMove(&all_moves, move)
			}
		}
	}
	rand.Seed(time.Now().UnixNano())
	random_index := rand.Intn(len(all_moves))
	random_move := all_moves[byte(random_index)]
	from := random_move.FromPosition
	to := random_move.ToPosition
	return_str := position2str(from) + position2str(to)
	if random_move.Promote {
		return_str += "+"
	}
	return return_str
}