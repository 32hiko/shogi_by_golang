package shogi

import (
	. "logger"
	"math/rand"
	"time"
)

type IPlayer interface {
	Search(*TBan) (string, int)
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

func (player TSlidePlayer) Search(ban *TBan) (string, int) {
	var te string
	if *(player.i)%2 == 0 {
		te = "8b7b"
	} else {
		te = "7b8b"
	}
	*(player.i)++
	return te, 0
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

func (player TRandomPlayer) Search(ban *TBan) (string, int) {
	logger := GetLogger()
	teban := *(ban.Teban)
	logger.Trace("[RandomPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " + s(teban))

	all_moves := MakeAllMoves(ban)

	moves_count := len(all_moves)
	logger.Trace("[RandomPlayer] moves: " + s(moves_count))
	if moves_count == 0 {
		return "resign", 0
	}
	rand.Seed(time.Now().UnixNano())
	random_index := rand.Intn(len(all_moves))
	random_move := all_moves[random_index]
	return random_move.GetUSIMoveString(), 0
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

func (player TKikiPlayer) Search(ban *TBan) (string, int) {
	logger := GetLogger()
	teban := *(ban.Teban)
	logger.Trace("[KikiPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " + s(teban))

	all_moves := MakeAllMoves(ban)

	moves_count := len(all_moves)
	logger.Trace("[KikiPlayer] moves: " + s(moves_count))
	if moves_count == 0 {
		return "resign", 0
	}

	move := GetMaxKikiMove(ban, &all_moves)
	return move.GetUSIMoveString(), 0
}

func GetMaxKikiMove(ban *TBan, all_moves *map[int]*TMove) *TMove {
	logger := GetLogger()
	current_sfen := ban.ToSFEN(false)
	current_max := -99999
	current_move_key := 0
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
	Joseki *TJoseki
}

func NewMainPlayer() *TMainPlayer {
	joseki := NewJoseki()
	player := TMainPlayer{
		Joseki: joseki,
	}
	return &player
}

func (player TMainPlayer) Search(ban *TBan) (string, int) {
	logger := GetLogger()
	teban := *(ban.Teban)
	logger.Trace("[MainPlayer] ban.Tesuu: " + s(*(ban.Tesuu)) + ", teban: " + s(teban))

	all_moves := MakeAllMoves(ban)

	moves_count := len(all_moves)
	logger.Trace("[MainPlayer] moves: " + s(moves_count))
	if moves_count == 0 {
		return "resign", 0
	}

	joseki_move := player.GetJosekiMove(ban, &all_moves)
	if joseki_move != nil {
		return joseki_move.GetUSIMoveString(), 0
	}
	// move, score := player.GetMainBestMove2(ban, &all_moves)
	move, score := player.GetMainBestMove3(ban, &all_moves, 16, 1, true)
	return move.GetUSIMoveString(), score
}

func (player TMainPlayer) GetJosekiMove(ban *TBan, all_moves *map[int]*TMove) *TMove {
	logger := GetLogger()
	current_sfen := ban.ToSFEN(false)

	// 現在の局面に定跡が登録されているか確認する
	fix_move, fix_move_exists := player.Joseki.FixOpening[*(ban.Tesuu)+1]
	var fix_move_string string = ""
	if fix_move_exists {
		fix_move_string = fix_move.GetUSIMoveString()
		logger.Trace("[MainPlayer] fix_move_string is: " + fix_move_string)
	} else {
		sfen_joseki_move, exists := player.Joseki.SFENMap[current_sfen]
		if exists {
			fix_move_string = sfen_joseki_move.GetUSIMoveString()
			logger.Trace("[MainPlayer] sfen_joseki_move_string is: " + fix_move_string)
		}
	}

	// 定跡手が存在する場合、（念のため）手があることを確認して返す
	if fix_move_string != "" {
		for key, move := range *all_moves {
			move_string := move.GetUSIMoveString()
			if fix_move_string == move_string {
				return (*all_moves)[key]
			} else {
				continue
			}
		}
	}
	return nil
}

func (player TMainPlayer) GetMainBestMove3(ban *TBan, all_moves *map[int]*TMove, width int, depth int, is_disp bool) (*TMove, int) {
	logger := GetLogger()
	teban := *(ban.Teban)
	logger.Trace("[BestMove3] depth: " + s(depth))
	current_sfen := ban.ToSFEN(false)

	// 1手指して有力そうな数手は、相手の応手も考慮する
	better_moves_map := make(map[int]int)

	// all_movesに対してEvaluateし、上位width件に絞り込む。改善の余地あり
	for key, move := range *all_moves {
		new_ban := FromSFEN(current_sfen)
		move_string := move.GetUSIMoveString()
		// 実際に動かしてみる
		new_ban.ApplyMove(move_string)
		result_sente, result_gote := new_ban.Analyze()
		score := Evaluate(result_sente, result_gote, teban)
		// logger.Trace("    [MainPlayer] move: " + move_string + " score: " + s(score))
		if is_disp {
			Resp("info time 0 depth 1 nodes 1 score cp "+ToDisplayScore(score, teban)+" pv "+move_string, logger)
		}
		// 上位width件に絞り込む
		min := -99999
		if len(better_moves_map) >= width {
			temp_min := 99999
			for c, _ := range better_moves_map {
				if c < temp_min {
					temp_min = c
				}
			}
			min = temp_min
		}
		// 今保持している上位width件の評価値より高い手は保持する
		if min < score {
			// 保持している中で最低値を削除
			delete(better_moves_map, min)
			// 保存する
			_, ok := better_moves_map[score]
			if ok {
				score++
			}
			better_moves_map[score] = key
			// logger.Trace("    [MainPlayer] min: " + s(min) + " score: " + s(score))
		}
	}

	current_move_key := 0
	current_score := 0
	if depth >= 2 {
		// depthが2以上なら、絞り込んだ結果を元に、相手の手番でdepth-1手先まで読む。
		current_max := -99999
		for score, key := range better_moves_map {
			new_ban := FromSFEN(current_sfen)
			move := (*all_moves)[key]
			move_string := move.GetUSIMoveString()
			new_ban.ApplyMove(move_string)
			next_moves := MakeAllMoves(new_ban)
			next_best_move, _ := player.GetMainBestMove3(new_ban, &next_moves, width /2, depth -1, false)
			if next_best_move == nil {
				// 手がないのはつまり詰み。
				current_move_key = key
				logger.Trace("[BestMove3] tsumi: " + move_string)
				break
			}
			next_best_move_string := next_best_move.GetUSIMoveString()
			new_ban.ApplyMove(next_best_move_string)
			result_sente, result_gote := new_ban.Analyze()
			count := Evaluate(result_sente, result_gote, !teban)
			if is_disp {
				Resp("info time 0 depth 1 nodes 1 score cp "+ToDisplayScore(count, teban)+" pv "+move_string+" "+next_best_move_string, logger)
			}
			if current_max < count {
				current_max = count
				current_move_key = key
				current_score = score
			} else {
				if current_max == count {
					if current_score < score {
						current_move_key = key
						current_score = score
					}
				}
			}
		}
	} else {
		// depthが1なら、上位width件の中で最高の評価値の手を返す
		max := -99999
		for score, key := range better_moves_map {
			if score > max {
				max = score
				current_move_key = key
				current_score = score
			}
		}
	}
	logger.Trace("[BestMove3] score: " + s(current_score))
	selected_move := (*all_moves)[current_move_key]
	return selected_move, current_score
}

func Evaluate(result_sente map[string]int, result_gote map[string]int, teban TTeban) int {
	sente_point := DoEvaluate(result_sente)
	gote_point := DoEvaluate(result_gote)
	point := 0
	if teban {
		point = sente_point - gote_point
	} else {
		point = gote_point - sente_point
	}
	return point
}

func DoEvaluate(result map[string]int) int {
	point := 0
	point += result["kiki"] * 5
	point += result["kikiMasu"] * 5
	point += result["koma"] * 100
	point += result["himoKoma"] * 10
	point += result["ukiKoma"] * -10
	point += result["atariKoma"] * -100
	point += result["tadaKoma"] * -300
	point += result["nariKoma"] * 10
	point += result["mochigomaCount"] * 200
	return point
}

func MakeAllMoves(ban *TBan) map[int]*TMove {
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

	all_moves := make(map[int]*TMove)
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

func RespondOute(ban *TBan, koma_moves *map[TKomaId]*TMoves, jigyoku *TKoma, oute_kiki *map[TKomaId]TKiki, all_moves *map[int]*TMove) {
	// logger := GetLogger()
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
						// logger.Trace("[MainPlayer] RespondOute move: " + move.Display())
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

func ToDisplayScore(score int, teban TTeban) string {
	i := score
	if !teban {
		i *= -1
	}
	return s(i)
}
