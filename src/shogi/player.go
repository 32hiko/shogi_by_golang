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
	logger.Trace("[RandomPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " +s(teban))
	tegoma := ban.GetTebanKoma(teban)
	koma_moves := make(map[TKomaId]*TMoves)

	// 1.自殺手（玉）を探し、除外する
	jigyoku := FindJiGyoku(ban, teban)
	jigyoku_moves := ban.FilterSuicideMoves(jigyoku)
	koma_moves[jigyoku.Id] = jigyoku_moves

	// 2.自殺手（pin）を探し、除外する
	ban.FilterPinnedMoves(jigyoku, &koma_moves)
	MergeMoves(&koma_moves, tegoma, ban)

	// 3.自玉に王手がかかっているかチェックする
	oute_kiki := ban.AllMasu[jigyoku.Position].GetAiteKiki(teban)

	all_moves := make(map[byte]*TMove)
	if len(*oute_kiki) > 0 {
		// 王手を回避する
		// 玉が逃げる手
		for _, move := range jigyoku_moves.Map {
			AddMove(&all_moves, move)
		}
		// 両王手でなければ、王手かけてる駒を取る手か、合い駒する
		if len(*oute_kiki) == 1 {
			for target_id, _ := range *oute_kiki {
				// 1個しかないのにforを使う強引実装
				target_koma := ban.AllKoma[target_id]
				for _, moves := range koma_moves {
					for _, move := range moves.Map {
						if move.ToPosition == target_koma.Position {
							AddMove(&all_moves, move)
						}
					}
				}
				// 王手かけてる駒が遠利きなら
				if target_koma.CanFarMove() {
					aida_map := make(map[TPosition]string)
					aida := jigyoku.Position - target_koma.Position
					for p := target_koma.Position + aida.Vector(); p != jigyoku.Position; p += aida.Vector() {
						aida_map[p] = ""
					}
					// 合い駒になる手を探す
					for _, moves := range koma_moves {
						for _, move := range moves.Map {
							_, ok := aida_map[move.ToPosition]
							if ok {
								AddMove(&all_moves, move)
							}
						}
					}
				}
			}
		}
	} else {
		// 今までどおり全部の手からランダム
		for _, moves := range koma_moves {
			for _, move := range moves.Map {
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

func FindJiGyoku(ban *TBan, teban TTeban) *TKoma {
	var gyoku *TKoma
	gyoku_map := ban.FindKoma(teban, Gyoku)
	for _, g := range *gyoku_map {
		// 1個しかないのにforを使う強引実装
		gyoku = g
	}
	return gyoku
}

func MergeMoves(moves *map[TKomaId]*TMoves, tegoma *(map[TKomaId]*TKoma), ban *TBan) {
	for koma_id, _ := range *tegoma {
		_, ok := (*moves)[koma_id]
		if !ok {
			(*moves)[koma_id] = ban.AllMoves[koma_id]
		}
	}
}

/**
 * 利きに反応するようにしてみたが、いちいち当たった駒が逃げようとするも逃げられないので弱くなった。サンプルとして保存しておく。
 **/
func (player TRandomPlayer) SearchSample(ban *TBan) string {
	logger := GetLogger()
	teban := *(ban.Teban)
	// logger.Trace("[RandomPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " +s(teban))
	tegoma := ban.GetTebanKoma(teban)

	// 自玉に王手がかかっているかどうかチェックする
	var oute_kiki *map[TKomaId]TKiki
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
		// 王手を回避する
		// 玉が逃げる手
		for _, move := range ban.AllMoves[gyoku_id].Map {
			AddMove(&all_moves, move)
		}
		// 王手かけてる駒を取る手
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
		// TODO:合い駒
	} else {
		aimed_koma := make(map[TKomaId]*TKoma)
		// 今までどおり全部の手からランダム
		for koma_id, _ := range *tegoma {
			// logger.Trace("[RandomPlayer] koma_id: " + s(koma_id))
			koma := ban.AllKoma[koma_id]
			masu := ban.AllMasu[koma.Position]
			aite_kiki := masu.GetAiteKiki(teban)
			jibun_kiki := masu.GetKiki(teban)
			if len(*aite_kiki) > len(*jibun_kiki) {
				logger.Trace(koma.Display() + " is aimed.")
				aimed_koma[koma_id] = koma
			}
		}
		if len(aimed_koma) > 0 {
			for koma_id, _ := range aimed_koma {
				for _, move := range ban.AllMoves[koma_id].Map {
					// 相手の利きが上回るマスには進まないでみる
					saki_masu := ban.AllMasu[move.ToPosition]
					aite_kiki := saki_masu.GetAiteKiki(teban)
					jibun_kiki := saki_masu.GetKiki(teban)
					if len(*aite_kiki) > len(*jibun_kiki) {
						logger.Trace("cut move [" + s(move.FromPosition) + " to " + s(move.ToPosition) + "]")
					} else {
						AddMove(&all_moves, move)
					}
				}
			}
		}
		if len(all_moves) == 0 {
			for koma_id, _ := range *tegoma {
				for _, move := range ban.AllMoves[koma_id].Map {
					// 相手の利きが上回るマスには進まないでみる
					saki_masu := ban.AllMasu[move.ToPosition]
					aite_kiki := saki_masu.GetAiteKiki(teban)
					jibun_kiki := saki_masu.GetKiki(teban)
					if len(*aite_kiki) > len(*jibun_kiki) {
						logger.Trace("cut move [" + s(move.FromPosition) + " to " + s(move.ToPosition) + "]")
					} else {
						AddMove(&all_moves, move)
					}
				}
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
