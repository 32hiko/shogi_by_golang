package shogi

import (
	. "logger"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

type IPlayer interface {
	Search(*TBan, int) (string, int)
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

func (player TSlidePlayer) Search(ban *TBan, ms int) (string, int) {
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

func (player TRandomPlayer) Search(ban *TBan, ms int) (string, int) {
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

func (player TKikiPlayer) Search(ban *TBan, ms int) (string, int) {
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

func (player TMainPlayer) Search(ban *TBan, ms int) (string, int) {
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

	// magic number
	width := 999
	depth := 3
	if ms < 300000 {
		depth = 2
	}
	if ms < 120000 {
		depth = 1
	}
	if ms < 60000 {
		// magic number
		depth = 0
	}

	// move, score := player.GetMainBestMove4(ban, &all_moves, (ms < 180000))
	move, score := player.GetMainBestMove3(ban, &all_moves, width, depth, true)

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

func GetScore(sfen string, teban TTeban, all_moves *map[int]*TMove, is_disp bool) <-chan string {
	score_channel := make(chan string)

	for key, move := range *all_moves {
		// ゴルーチンによりマルチコアを使い、時間短縮をはかる
		go scoreRoutine(sfen, teban, key, move, score_channel, is_disp)
	}

	return score_channel
}

func scoreRoutine(sfen string, teban TTeban, key int, move *TMove, score_channel chan string, is_disp bool) {
	logger := GetLogger()
	new_ban := FromSFEN(sfen)
	move_string := move.GetUSIMoveString()
	// 実際に動かしてみる
	new_ban.ApplyMove(move_string)
	result_sente, result_gote := new_ban.Analyze()
	score := Evaluate(result_sente, result_gote, teban)
	if move.IsForward(teban) {
		// 前に進む手を評価する
		score += 50
	}
	if new_ban.IsTadaMove(move) {
		// score -= 100000
		// logger.Trace("penalty move: " + move_string)
	}
	score /= 50
	if is_disp {
		Resp("info time 0 depth 1 nodes 1 score cp "+ToDisplayScore(score, teban)+" pv "+move_string, logger)
	}
	resp_string := s(key) + ":" + s(score)
	if IsOute(new_ban, !teban) {
		resp_string += ":Oute"
	}
	// "key_of_all_moves:score:Oute(if oute)"
	score_channel <- resp_string
}

type TNodeScore struct {
	Key       int
	Moves     string
	Score     int
	RespScore int
}

func NewNodeScore(key int, moves string, score int, resp_score int) *TNodeScore {
	node_score := TNodeScore{
		Key:       key,
		Moves:     moves,
		Score:     score,
		RespScore: resp_score,
	}
	return &node_score
}

func (player TMainPlayer) GetMainBestMove4(ban *TBan, all_moves *map[int]*TMove, hurry bool) (*TMove, int) {
	logger := GetLogger()
	teban := *(ban.Teban)
	current_sfen := ban.ToSFEN(false)

	// 現局面から、自分の全候補手を評価
	score_channel := GetScore(current_sfen, teban, all_moves, true)
	move_key_score_map := make(map[int]int)
	// ゴルーチンの結果待ち
	for i := 0; i < len(*all_moves); i++ {
		result := <-score_channel
		// "key_of_all_moves:score:Oute(if oute)"
		result_arr := strings.Split(result, ":")
		key, _ := strconv.Atoi(result_arr[0])
		score, _ := strconv.Atoi(result_arr[1])
		move_key_score_map[key] = score
	}

	node_score_map := make(map[int]*TNodeScore)
	for key, score := range move_key_score_map {
		next_ban := FromSFEN(current_sfen)
		move := (*all_moves)[key]
		move_string := move.GetUSIMoveString()
		next_ban.ApplyMove(move_string)
		next_moves := MakeAllMoves(next_ban)
		if len(next_moves) == 0 {
			// 相手の手がないなら詰み（一手詰み）
			return (*all_moves)[key], 99999
		}
		next_sfen := next_ban.ToSFEN(false)
		// 自分の候補手を指した後の局面から、相手の全候補手を評価
		next_score_channel := GetScore(next_sfen, !teban, &next_moves, false)
		next_ks_map := make(map[int]int)
		// ゴルーチンの結果待ち
		for i := 0; i < len(next_moves); i++ {
			result := <-next_score_channel
			// "key_of_all_moves:score:Oute(if oute)"
			result_arr := strings.Split(result, ":")
			key, _ := strconv.Atoi(result_arr[0])
			score, _ := strconv.Atoi(result_arr[1])
			next_ks_map[key] = score
		}
		// next_ks_mapから、最大の評価値の手を取得する
		resp_score := -99999
		resp_key := 0
		for k, s := range next_ks_map {
			if resp_score < s {
				resp_score = s
				resp_key = k
			}
		}
		// TNodeScoreに保存する
		if resp_key != 0 {
			next_move := next_moves[resp_key].GetUSIMoveString()
			moves_string := move_string + " " + next_move
			ns := NewNodeScore(key, moves_string, score, resp_score)
			node_score_map[key] = ns
		}
	}

	// TNodeScoreのうち、一部をさらに読む
	if !hurry {
		logger.Trace("before CloseUp, length: " + s(len(node_score_map)))
		new_map := CloseUp(current_sfen, teban, &node_score_map)
		logger.Trace("ok")
		node_score_map = new_map
		logger.Trace("after  CloseUp, length: " + s(len(node_score_map)))
	}

	// TNodeScoreのうち、RespScoreが最小の手を取得し、返す
	return_score := 99999
	return_key := 0
	for key, entry := range node_score_map {
		Resp("info time 0 depth 1 nodes 1 score cp "+ToDisplayScore(entry.RespScore, teban)+" pv "+entry.Moves, logger)
		if return_score > entry.RespScore {
			return_score = entry.RespScore
			return_key = key
		} else {
			if return_score == entry.RespScore {
				return_org_score := node_score_map[return_key].Score
				if return_org_score < entry.Score {
					return_score = entry.RespScore
					return_key = key
				}
			}
		}
	}
	return (*all_moves)[return_key], node_score_map[return_key].Score
}

func CloseUp(sfen string, teban TTeban, score_map *map[int]*TNodeScore) map[int]*TNodeScore {
	logger := GetLogger()
	width_max := 10

	// RespScoreが小さい上位width件に絞り込む
	var new_map map[int]*TNodeScore = nil
	if len(*score_map) <= width_max {
		// 件数がもともと少ないならこの工程は不要
		new_map = *score_map
	} else {
		// 件数が多い場合、新しいmapに詰め替える
		new_map = make(map[int]*TNodeScore)
		// ソートのためのスライス
		slice := make([]int, len(*score_map))
		for _, v := range *score_map {
			slice = append(slice, v.RespScore)
		}
		sort.Sort(sort.IntSlice(slice))
		var border_value int = slice[width_max-1]
		logger.Trace("slice: " + s(slice))
		logger.Trace("border_value: " + s(border_value))
		for k, v := range *score_map {
			if v.RespScore <= border_value {
				new_map[k] = v
				logger.Trace("add key, value: " + s(k) + " " + s(v.RespScore))
			}
		}
	}
	logger.Trace("CloseUp width: " + s(len(new_map)))
	for k, v := range new_map {
		ban := FromSFEN(sfen)
		moves_arr := strings.Split(v.Moves, " ")
		ban.ApplyMove(moves_arr[0])
		ban.ApplyMove(moves_arr[1])
		next_moves := MakeAllMoves(ban)
		if len(next_moves) == 0 {
			// 手がないならこちらが詰み
			new_v := v
			new_v.RespScore = 99999
			new_map[k] = new_v
		}
		next_sfen := ban.ToSFEN(false)
		// 相手の候補手を指した後の局面から、自分の全候補手を評価
		next_score_channel := GetScore(next_sfen, teban, &next_moves, false)
		next_ks_map := make(map[int]int)
		// ゴルーチンの結果待ち
		for i := 0; i < len(next_moves); i++ {
			result := <-next_score_channel
			// "key_of_all_moves:score:Oute(if oute)"
			result_arr := strings.Split(result, ":")
			key, _ := strconv.Atoi(result_arr[0])
			score, _ := strconv.Atoi(result_arr[1])
			next_ks_map[key] = score
		}
		// next_ks_mapから、最大の評価値の手を取得する
		resp_score := -99999
		resp_key := 0
		for k, s := range next_ks_map {
			if resp_score < s {
				resp_score = s
				resp_key = k
			}
		}
		// TNodeScoreに保存する
		if resp_key != 0 {
			next_move := next_moves[resp_key].GetUSIMoveString()
			moves_string := v.Moves + " " + next_move
			new_v := v
			// 低いの優先なので反転させる
			new_v.RespScore = -resp_score
			new_v.Moves = moves_string
			new_map[k] = new_v
		}
	}
	logger.Trace("CloseUp final width: " + s(len(new_map)))
	return new_map
}

func (player TMainPlayer) GetMainBestMove3(ban *TBan, all_moves *map[int]*TMove, width int, depth int, is_disp bool) (*TMove, int) {
	logger := GetLogger()
	teban := *(ban.Teban)
	current_sfen := ban.ToSFEN(false)
	score_channel := GetScore(current_sfen, teban, all_moves, is_disp)

	oute_map := make(map[int]int)
	better_moves_map := make(map[int]int)

	if depth == 0 {
		// 緊急避難ロジック
		return (*all_moves)[0], 28
	}

	if width == 999 {
		width = len(*all_moves) / 2
	}
	if width > 32 {
		width = 32
	}

	// logger.Trace("------start------")
	// ゴルーチンの結果待ち
	for i := 0; i < len(*all_moves); i++ {
		result := <-score_channel
		// "key_of_all_moves:score:Oute(if oute)"
		result_arr := strings.Split(result, ":")
		key, _ := strconv.Atoi(result_arr[0])
		score, _ := strconv.Atoi(result_arr[1])
		if len(result_arr) == 3 {
			// 王手フラグあり。候補手とする
			oute_map[key] = score
		} else {
			// 上位width件だけ候補手とする
			PutToMap(&better_moves_map, key, score, width)
		}
	}

	// 王手を候補手に追加
	for k, s := range oute_map {
		better_moves_map[k] = s
	}

	next_width := width / 4
	next_depth := depth - 1

	current_move_key := 0
	current_score := 0
	if depth >= 2 {
		// depthが2以上なら、絞り込んだ結果を元に、相手の手番でdepth-1手先まで読む。
		current_min := 99999
		// logger.Trace("depth: " + s(depth))
		// logger.Trace("moves: " + s(len(better_moves_map)))
		for key, score := range better_moves_map {
			new_ban := FromSFEN(current_sfen)
			move := (*all_moves)[key]
			move_string := move.GetUSIMoveString()
			new_ban.ApplyMove(move_string)
			// logger.Trace("before GetMainBestMove3 " + move_string)
			next_moves := MakeAllMoves(new_ban)
			if len(next_moves) == 0 {
				// 手がないので詰み。下のは読んだ先で詰みがある場合？
				current_move_key = key
				current_score = 99999
				logger.Trace("[BestMove3] tsumi: " + move_string)
				break
			}
			next_best_move, count := player.GetMainBestMove3(new_ban, &next_moves, next_width, next_depth, false)
			if next_best_move == nil {
				// 手がないのはつまり詰み。
				current_move_key = key
				current_score = -99999
				logger.Trace("[BestMove3] tsumi: " + move_string)
				break
			}
			next_best_move_string := next_best_move.GetUSIMoveString()
			if is_disp {
				Resp("info time 0 depth 1 nodes 1 score cp "+ToDisplayScore(count, teban)+" pv "+move_string+" "+next_best_move_string, logger)
			}
			if current_min > count {
				current_min = count
				current_move_key = key
				current_score = score
			} else {
				if current_min == count {
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
		for key, score := range better_moves_map {
			if score > max {
				max = score
				current_move_key = key
				current_score = score
			}
		}
	}
	// logger.Trace("[BestMove3] score: " + s(current_score))
	selected_move := (*all_moves)[current_move_key]
	// logger.Trace("------ end ------")
	return selected_move, current_score
}

func PutToMap(m *map[int]int, k int, s int, w int) {
	// 項目が上限に達している
	if len(*m) == w {
		min_k := 0
		min_s := 99999
		// 現在の最小の項目を取得する
		for mk, ms := range *m {
			if ms < min_s {
				min_k = mk
				min_s = ms
			}
		}
		// 現在の最小より新しい項目が大きい場合、入れ替えのため最小を削除
		if min_s < s {
			delete(*m, min_k)
			// 項目を追加する
			(*m)[k] = s
		}
	} else {
		// 項目を追加する
		(*m)[k] = s
	}
}

func IsOute(ban *TBan, aite_teban TTeban) bool {
	aite_gyoku := FindGyoku(ban, aite_teban)
	oute_kiki := ban.AllMasu[aite_gyoku.Position].GetAiteKiki(aite_teban)
	return len(*oute_kiki) > 0
}

func Evaluate(result_sente map[string]int, result_gote map[string]int, teban TTeban) int {
	// logger := GetLogger()
	sente_point := DoEvaluate(result_sente, (teban == Sente))
	gote_point := DoEvaluate(result_gote, (teban == Gote))
	point := 0
	if teban {
		point = sente_point - gote_point
	} else {
		point = gote_point - sente_point
	}
	// logger.Trace("  Evaluate: " + s(sente_point) + "," + s(gote_point))
	return point
}

func DoEvaluate(result map[string]int, is_aiteban bool) int {
	point := 0
	point += result["kiki"] * 5
	point += result["kikiMasu"] * 20
	point += result["koma"] * 5000
	point += result["himoKoma"] * 5
	point += result["ukiKoma"] * -5
	if is_aiteban {
		point += result["atariKoma"] * -10000
		point += result["tadaKoma"] * -10000
	}
	point += result["nariKoma"] * 10
	point += result["mochigomaCount"] * 1000
	return point
}

func MakeAllMoves(ban *TBan) map[int]*TMove {
	teban := *(ban.Teban)
	tegoma := ban.GetTebanKoma(teban)
	koma_moves := make(map[TKomaId]*TMoves)

	// 1.自殺手（玉）を探し、除外する
	jigyoku := FindGyoku(ban, teban)
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

func FindGyoku(ban *TBan, teban TTeban) *TKoma {
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
	teban := *(ban.Teban)
	// 両王手でなければ、王手かけてる駒を取る手か、合い駒する
	if len(*oute_kiki) == 1 {
		for target_id, _ := range *oute_kiki {
			// 1個しかないのにforを使う強引実装
			target_koma := ban.AllKoma[target_id]
			target_kiki := ban.AllMasu[target_koma.Position].GetAiteKiki(!teban)
			if len(*target_kiki) > 0 {
				for toru_id, _ := range *target_kiki {
					// toru_koma := ban.AllKoma[toru_id]
					for _, move := range (*koma_moves)[toru_id].Map {
						if move.ToPosition == target_koma.Position {
							AddMove(all_moves, move)
						}
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
