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
	case "Kiki":
		return NewKikiPlayer()
	case "Main":
		return NewMainPlayer()
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
	logger.Trace("[RandomPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " + s(teban))

	all_moves := MakeAllMoves(ban)

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

/*
 * 利きが通っているマスの数だけで評価してみる。
 * →相手の大駒の近くに駒を寄せるだけだった。
 */
type TKikiPlayer struct {
}

func NewKikiPlayer() *TKikiPlayer {
	player := TKikiPlayer{}
	return &player
}

func (player TKikiPlayer) Search(ban *TBan) string {
	logger := GetLogger()
	teban := *(ban.Teban)
	logger.Trace("[KikiPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " + s(teban))

	all_moves := MakeAllMoves(ban)

	moves_count := len(all_moves)
	logger.Trace("[KikiPlayer] moves: " + s(moves_count))
	if moves_count == 0 {
		return "resign"
	}

	move := GetMaxKikiMove(ban, &all_moves)
	return move.GetUSIMoveString()
}

func GetMaxKikiMove(ban *TBan, all_moves *map[byte]*TMove) *TMove {
	logger := GetLogger()
	current_sfen := ban.ToSFEN()
	current_max := -81
	var current_move_key byte = 0
	for key, move := range *all_moves {
		new_ban := FromSFEN(current_sfen)
		new_ban.ApplyMove(move.GetUSIMoveString())
		count := new_ban.CountKikiMasu(*(ban.Teban))
		count -= new_ban.CountKikiMasu(!*(ban.Teban))
		if current_max < count {
			logger.Trace("[KikiPlayer] count: " + s(count))
			current_max = count
			current_move_key = key
		}
	}
	return (*all_moves)[current_move_key]
}

type TMainPlayer struct {
}

func NewMainPlayer() *TMainPlayer {
	player := TMainPlayer{}
	return &player
}

func (player TMainPlayer) Search(ban *TBan) string {
	logger := GetLogger()
	teban := *(ban.Teban)
	logger.Trace("[MainPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " + s(teban))

	all_moves := MakeAllMoves(ban)

	moves_count := len(all_moves)
	logger.Trace("[MainPlayer] moves: " + s(moves_count))
	if moves_count == 0 {
		return "resign"
	}

	move := GetMainBestMove(ban, &all_moves)
	return move.GetUSIMoveString()
}

func GetMainBestMove(ban *TBan, all_moves *map[byte]*TMove) *TMove {
	logger := GetLogger()
	teban := *(ban.Teban)
	current_sfen := ban.ToSFEN()
	current_max := -81

	// 最終手に反応するための準備
	last_move_map := make(map[TPosition]string)
	if ban.LastMoveTo != nil {
		logger.Trace("[MainPlayer] LastMoveTo is: " + s(*(ban.LastMoveTo)))
		last_move_masu := ban.AllMasu[*(ban.LastMoveTo)]
		last_move_koma_moves := ban.AllMoves[last_move_masu.KomaId]
		for _, move := range last_move_koma_moves.Map {
			last_move_map[move.ToPosition] = ""
		}
	}

	var current_move_key byte = 0
	for key, move := range *all_moves {
		// 実際に動かしてみる
		new_ban := FromSFEN(current_sfen)
		new_ban.ApplyMove(move.GetUSIMoveString())

		// 利いているマスの数の評価
		masu_count := new_ban.CountKikiMasu(teban)
		masu_count -= new_ban.CountKikiMasu(!teban)
		masu_count *= 10 // 調整パラメーター

		// 駒得（1手しか読まないので駒の枚数だけ）
		teban_koma := new_ban.GetTebanKoma(teban)
		komadoku_point := len(*teban_koma)
		komadoku_point *= 50

		// タダ捨てを抑止したい
		// 単純に利きの数だけだと、自分の駒の利いてる範囲でうろつくだけになる。タダの地点だけマイナスするほうがよさそう
		move_masu := new_ban.AllMasu[move.ToPosition]
		teban_kiki := move_masu.GetKiki(teban)
		aite_kiki := move_masu.GetAiteKiki(teban)
		tada_point := len(*teban_kiki) - len(*aite_kiki)
		if tada_point < 0 {
			tada_point *= 200
		} else {
			tada_point *= 10
		}

		// 相手の手に反応するため、最後の手の利きを重く捉える
		// これ入れても、変な反応をして弱くなる
		_, ok := last_move_map[move.FromPosition]
		escape_point := 0
		if ok {
			escape_point = 100
		}

		forward_point := 0
		if move.IsForward(teban) {
			forward_point = 100
		}

		count := masu_count + komadoku_point + tada_point + escape_point + forward_point
		if current_max < count {
			logger.Trace("[MainPlayer] count: " + s(count))
			current_max = count
			current_move_key = key
		}
	}
	return (*all_moves)[current_move_key]
}

func MakeAllMoves(ban *TBan) map[byte]*TMove {
	teban := *(ban.Teban)
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
		RespondOute(ban, &koma_moves, jigyoku, oute_kiki, &all_moves)
	} else {
		// 今までどおり全部の手から
		for _, moves := range koma_moves {
			for _, move := range moves.Map {
				AddMove(&all_moves, move)
			}
		}
	}
	return all_moves
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

func RespondOute(ban *TBan, koma_moves *map[TKomaId]*TMoves, jigyoku *TKoma, oute_kiki *map[TKomaId]TKiki, all_moves *map[byte]*TMove) {
	logger := GetLogger()
	// 玉が逃げる手
	for _, move := range (*koma_moves)[jigyoku.Id].Map {
		AddMove(all_moves, move)
	}
	// 両王手でなければ、王手かけてる駒を取る手か、合い駒する
	if len(*oute_kiki) == 1 {
		for target_id, _ := range *oute_kiki {
			// 1個しかないのにforを使う強引実装
			target_koma := ban.AllKoma[target_id]
			for _, moves := range *koma_moves {
				for _, move := range moves.Map {
					if move.ToPosition == target_koma.Position {
						AddMove(all_moves, move)
						logger.Trace("[MainPlayer] RespondOute move: " + move.Display())
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
				for _, moves := range *koma_moves {
					for _, move := range moves.Map {
						_, ok := aida_map[move.ToPosition]
						if ok {
							AddMove(all_moves, move)
						}
					}
				}
			}
		}
	}
}
