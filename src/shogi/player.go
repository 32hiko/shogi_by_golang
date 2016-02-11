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
	teban := *(ban.Teban)
	// logger.Trace("[RandomPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " +s(teban))
	tegoma := ban.GetTebanKoma(teban)

	// 自玉に王手がかかっているかどうかチェックする
	var oute_kiki *map[TKomaId]string
	var gyoku_id TKomaId
	gyoku_map := ban.FindKoma(teban, Gyoku)
	for _, gyoku := range *gyoku_map {
		// 1個しかないのにforを使う強引実装
		gyoku_id = gyoku.Id
		masu := ban.AllMasu[gyoku.Position]
		oute_kiki = masu.GetAiteKiki(teban)
	}

	all_moves := make(map[byte]*TMove)
	if len(*oute_kiki) > 0 {
		// 王手を回避しないと
		// 暫定的に、玉が逃げる手だけのランダムで
		for _, move := range ban.AllMoves[gyoku_id].Map {
			AddMove(&all_moves, move)
		}
		// 逃げる手がないと現状即投了してしまう。TODO: 王手かけてる駒を取るか、合い駒するか
		if len(*oute_kiki) == 1 {
			// 王手かけてる駒を取る手を探す
			for target_id, _ := range *oute_kiki {
				// 1個しかないのにforを使う強引実装
				target_koma := ban.AllKoma[target_id]
				for koma_id, _ := range *tegoma {
					for _, move := range ban.AllMoves[koma_id].Map {
						if move.ToPosition == target_koma.Position {
							AddMove(&all_moves, move)
						}
					}
				}
			}
		}
	} else {
		// 今までどおりランダム
		for koma_id, _ := range *tegoma {
			// logger.Trace("[RandomPlayer] koma_id: " + s(koma_id))
			for _, move := range ban.AllMoves[koma_id].Map {
				AddMove(&all_moves, move)
			}
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
