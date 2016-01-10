package shogi

import (
	"fmt"
	. "logger"
	"strconv"
	"strings"
)

// alias
var p = fmt.Println
var s = fmt.Sprint

type TBan struct {
	AllMasu   map[TPosition]*TMasu
	AllKoma   map[TKomaId]*TKoma
	AllMoves  map[TKomaId]*TMoves
	SenteKoma map[TKomaId]*TKoma
	GoteKoma  map[TKomaId]*TKoma
	Tesuu     *int
}

func NewBan() *TBan {
	all_masu := make(map[TPosition]*TMasu)
	// マスを初期化する
	var x, y byte = 1, 1
	for y <= 9 {
		x = 1
		for x <= 9 {
			pos := Bytes2TPosition(x, y)
			all_masu[pos] = NewMasu(pos, 0)
			x++
		}
		y++
	}
	// 持ち駒用
	all_masu[Mochigoma] = NewMasu(Mochigoma, 0)
	var tesuu int = 0

	ban := TBan{
		AllMasu:   all_masu,
		AllKoma:   make(map[TKomaId]*TKoma),
		AllMoves:  make(map[TKomaId]*TMoves),
		SenteKoma: make(map[TKomaId]*TKoma),
		GoteKoma:  make(map[TKomaId]*TKoma),
		Tesuu:     &tesuu,
	}
	return &ban
}

// 駒が持つデータ、マスが持つデータは今後も検討要
type TMasu struct {
	// マスの座標
	Position TPosition
	// 駒があれば駒のId
	KomaId TKomaId
	// このマスに利かせている駒のIdを入れる。ヒートマップを作るため
	SenteKiki *map[TKomaId]string // temp
	GoteKiki  *map[TKomaId]string // temp
}

func NewMasu(position TPosition, koma_id TKomaId) *TMasu {
	s_kiki := make(map[TKomaId]string)
	g_kiki := make(map[TKomaId]string)
	masu := TMasu{
		Position:  position,
		KomaId:    koma_id,
		SenteKiki: &s_kiki,
		GoteKiki:  &g_kiki,
	}
	return &masu
}

func (masu TMasu) SaveKiki(koma_id TKomaId, is_sente TTeban) {
	kiki := masu.GetKiki(is_sente)
	(*kiki)[koma_id] = ""
}

func (masu TMasu) DeleteKiki(koma_id TKomaId, is_sente TTeban) {
	kiki := masu.GetKiki(is_sente)
	delete(*kiki, koma_id)
}

func (masu TMasu) GetKiki(is_sente TTeban) *map[TKomaId]string {
	if is_sente {
		return masu.SenteKiki
	} else {
		return masu.GoteKiki
	}
}

func (masu TMasu) GetAiteKiki(is_sente TTeban) *map[TKomaId]string {
	if is_sente {
		return masu.GoteKiki
	} else {
		return masu.SenteKiki
	}
}

func Bytes2TPosition(x byte, y byte) TPosition {
	return TPosition(complex(float32(x), float32(y)))
}

func (ban TBan) GetTebanKoma(teban TTeban) *(map[TKomaId]*TKoma) {
	if teban {
		return &(ban.SenteKoma)
	} else {
		return &(ban.GoteKoma)
	}
}

func CreateInitialState() *TBan {
	ban := NewBan()
	ban.PutAllKoma()
	return ban
}

func (ban TBan) PutAllKoma() {
	// 駒を1つずつ生成する
	var koma_id TKomaId = 1

	// 後手
	var teban TTeban = Gote
	// 香
	ban.PutKoma(NewKoma(koma_id, Kyo, 1, 1, teban))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kyo, 9, 1, teban))
	koma_id++
	// 桂
	ban.PutKoma(NewKoma(koma_id, Kei, 2, 1, teban))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kei, 8, 1, teban))
	koma_id++
	// 銀
	ban.PutKoma(NewKoma(koma_id, Gin, 3, 1, teban))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Gin, 7, 1, teban))
	koma_id++
	// 金
	ban.PutKoma(NewKoma(koma_id, Kin, 4, 1, teban))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kin, 6, 1, teban))
	koma_id++
	// 王
	ban.PutKoma(NewKoma(koma_id, Gyoku, 5, 1, teban))
	koma_id++
	// 角
	ban.PutKoma(NewKoma(koma_id, Kaku, 2, 2, teban))
	koma_id++
	// 飛
	ban.PutKoma(NewKoma(koma_id, Hi, 8, 2, teban))
	koma_id++
	// 歩
	var x byte = 1
	for x <= 9 {
		ban.PutKoma(NewKoma(koma_id, Fu, x, 3, teban))
		koma_id++
		x++
	}

	// 先手
	teban = Sente
	// 香
	ban.PutKoma(NewKoma(koma_id, Kyo, 1, 9, teban))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kyo, 9, 9, teban))
	koma_id++
	// 桂
	ban.PutKoma(NewKoma(koma_id, Kei, 2, 9, teban))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kei, 8, 9, teban))
	koma_id++
	// 銀
	ban.PutKoma(NewKoma(koma_id, Gin, 3, 9, teban))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Gin, 7, 9, teban))
	koma_id++
	// 金
	ban.PutKoma(NewKoma(koma_id, Kin, 4, 9, teban))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kin, 6, 9, teban))
	koma_id++
	// 王
	ban.PutKoma(NewKoma(koma_id, Gyoku, 5, 9, teban))
	koma_id++
	// 角
	ban.PutKoma(NewKoma(koma_id, Kaku, 8, 8, teban))
	koma_id++
	// 飛
	ban.PutKoma(NewKoma(koma_id, Hi, 2, 8, teban))
	koma_id++
	// 歩
	x = 1
	for x <= 9 {
		ban.PutKoma(NewKoma(koma_id, Fu, x, 7, teban))
		koma_id++
		x++
	}
}

// 駒を配置し、合法手、利きマスデータを更新する
func (ban TBan) PutKoma(koma *TKoma) {
	// 駒が持っている位置を更新
	ban.AllKoma[koma.Id] = koma

	// ここは本来、駒の所有権が決まった時点でやる処理。初期化も、まず全部持ち駒にしてそれを打っていくのが正しい。
	if koma.IsSente {
		ban.SenteKoma[koma.Id] = koma
	} else {
		ban.GoteKoma[koma.Id] = koma
	}

	ban.AllMasu[koma.Position].KomaId = koma.Id

	// 配置した駒の合法手、利きを作成
	if koma.CanFarMove() {
		// 香、角、飛、馬、龍の遠利き部分。駒の有無も考慮しつつ作成する。
		ban.AllMoves[koma.Id] = ban.CreateFarMovesAndKiki(koma)
	} else {
		// TODO 同じロジックを使えるようになったはず。
		// 駒から、その駒の機械的な利き先を取得する
		all_moves := koma.GetAllMoves()
		// 機械的な利き先のうち自陣営の駒がいるマスを除き、有効な指し手となるマスを保存する
		ban.CheckMovesAndKiki(koma, all_moves)
		ban.AllMoves[koma.Id] = all_moves.DeleteInvalidMoves()
	}

	// 自マスに、他の駒からの利きとしてIdが入っている場合で、香、角、飛の場合は先の利きを止める
	// こちらも、龍や馬の周囲に駒を打った場合、的外れな方向の手や利きを消そうとしてしまうので、作り直しとする
	ban.DeleteCloseMovesAndKiki(koma, Sente)
	ban.DeleteCloseMovesAndKiki(koma, Gote)
}

func (ban TBan) DeleteCloseMovesAndKiki(koma *TKoma, is_sente TTeban) {
	kiki_map := ban.AllMasu[koma.Position].GetKiki(is_sente)
	if len(*kiki_map) > 0 {
		for koma_id, _ := range *kiki_map {
			target_koma := ban.AllKoma[koma_id]
			if target_koma.CanFarMove() {
				// 手も利きもいったん削除し、作りなおす
				// 利きがFarMoveによるものかそうでないか判断しにくいので、手も利きもいったん全削除→全追加
				ban.DeleteAllKiki(target_koma)
				ban.AllMoves[koma_id] = ban.CreateFarMovesAndKiki(target_koma)
			} else {
				if koma.IsSente != is_sente {
					// komaが敵陣営なら、komaの位置への手でkomaが取れることを手に保存する。
					var saved bool = false
					for _, move := range ban.AllMoves[koma_id].Map {
						if move.ToPosition == koma.Position {
							move.ToId = koma.Id
							saved = true
							break
						}
					}
					if !saved {
						// もともと自陣営の駒に利かせていたところ、その駒を取られた場合はmoveが存在していない。
						m := NewMove(target_koma.Id, target_koma.Position, koma.Position, koma.Id)
						ban.AllMoves[koma_id].Add(m)
						if !target_koma.Promoted {
							can_promote, promote_move := m.CanPromote(target_koma.IsSente)
							if can_promote {
								ban.AllMoves[koma_id].Add(promote_move)
							}
						}
					}
				} else {
					// komaが自陣営なら、komaの位置への手は合法でなくなるので削除が必要。
					for _, move := range ban.AllMoves[koma_id].Map {
						if move.ToPosition == koma.Position {
							move.IsValid = false
							break
						}
					}
				}
				ban.AllMoves[koma_id] = ban.AllMoves[koma_id].DeleteInvalidMoves()
			}
		}
	}
}

// 香、角、飛の、全方向の手と利きを生成する。
func (ban TBan) CreateFarMovesAndKiki(koma *TKoma) *TMoves {
	moves := NewMoves()
	if koma.Promoted {
		switch koma.Kind {
		case Kaku:
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_ne))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_se))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_nw))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_sw))
			m, _ := ban.Create1MoveAndKiki(koma, move_n)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_s)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_e)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_w)
			moves.AddAll(m)
		case Hi:
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_n))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_s))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_e))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_w))
			m, _ := ban.Create1MoveAndKiki(koma, move_ne)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_se)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_nw)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_sw)
			moves.AddAll(m)
		}
	} else {
		switch koma.Kind {
		case Kyo:
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_n))
		case Kaku:
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_ne))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_se))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_nw))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_sw))
		case Hi:
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_n))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_s))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_e))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_w))
		}
	}
	return moves
}

// 馬、龍の、成ってできた1マス分だけの手と利きを生成する。処理的にはNと同じ。
func (ban TBan) Create1MoveAndKiki(koma *TKoma, delta TPosition) ([]*TMove, bool) {
	slice := make([]*TMove, 0)
	var to_pos TPosition
	if koma.IsSente {
		to_pos = koma.Position + delta
	} else {
		to_pos = koma.Position - delta
	}
	if to_pos.IsValidMove() {
		// 利き先マスに、自駒のIdを保存する
		kiki_masu := ban.AllMasu[to_pos]
		kiki_masu.SaveKiki(koma.Id, koma.IsSente)
		target_id := kiki_masu.KomaId
		// 利き先マスに、駒があるかないか
		if target_id != 0 {
			// 駒がある
			target_koma := ban.AllKoma[target_id]
			if target_koma.IsSente == koma.IsSente {
				// 自陣営の駒のあるマスには指せない
				return slice, false
			} else {
				// 相手の駒は取れる。その先には動けない
				m := NewMove(koma.Id, koma.Position, to_pos, target_koma.Id)
				slice = append(slice, m)
				if !koma.Promoted {
					can_promote, promote_move := m.CanPromote(koma.IsSente)
					if can_promote {
						slice = append(slice, promote_move)
					}
				}
				return slice, false
			}
		} else {
			// 駒がないなら指せる
			m := NewMove(koma.Id, koma.Position, to_pos, 0)
			slice = append(slice, m)
			if !koma.Promoted {
				can_promote, promote_move := m.CanPromote(koma.IsSente)
				if can_promote {
					slice = append(slice, promote_move)
				}
			}
			return slice, true
		}
	}
	return slice, false
}

// 香、角、飛の、ある1方向の手と利きを生成する。
func (ban TBan) CreateNMovesAndKiki(koma *TKoma, delta TPosition) []*TMove {
	slice := make([]*TMove, 0)
	delta_base := delta
	for {
		moves, keep := ban.Create1MoveAndKiki(koma, delta_base)
		slice = append(slice, moves...)
		if keep {
			delta_base += delta
		} else {
			break
		}
	}
	return slice
}

// 駒の合法手と利き先マスをチェックする（香、角、飛を除く）
func (ban TBan) CheckMovesAndKiki(koma *TKoma, moves *TMoves) {
	for _, move := range moves.Map {
		temp_pos := move.ToPosition
		// 利き先マスに、自駒のIdを保存する
		kiki_masu := ban.AllMasu[temp_pos]
		kiki_masu.SaveKiki(koma.Id, koma.IsSente)
		target_id := kiki_masu.KomaId
		if target_id != 0 {
			if ban.AllKoma[target_id].IsSente == koma.IsSente {
				// 自陣営の駒がいるマスには指せない
				move.IsValid = false
			} else {
				// 指し手の取る駒
				move.ToId = target_id
			}
		}
	}
}

// USI形式のmoveを反映させる。
func (ban TBan) ApplyMove(usi_move string) {
	// usi_moveをこちらのmoveに変換する
	var from_str string
	var to_str string
	var promote bool
	if len(usi_move) == 5 {
		// 成り
		promote = true
	}
	if promote {
	}
	from_str = usi_move[0:2]
	to_str = usi_move[2:4]

	logger := GetLogger()
	// 駒を打つかどうか
	is_drop := strings.Index(from_str, "*")
	if is_drop == -1 {
		// 打たない
		from := str2Position(from_str)
		to := str2Position(to_str)

		logger.Trace("from: " + s(from) + ", to: " + s(to))
		// こちらのmoveを実行する
		ban.DoMove(from, to, promote)
	} else {
		// "*"を含む＝打つ。先手の銀打ちならS*,後手の銀打ちならs*で始め、打つマスの表記は同じ。
		kind, teban := str2KindAndTeban(from_str)
		to := str2Position(to_str)

		logger.Trace("駒打: " + teban_map[teban] + disp_map[kind] + ", to: " + s(to))
		ban.DoDrop(teban, kind, to)
	}

	ban.UpdateMoves()
	ban.DeleteSuicideMoves()
	*(ban.Tesuu) += 1
}

func (ban TBan) UpdateMoves() {
	// 両陣営の玉について、手を生成し直す（次のロジックによる削除が、自動では復元されないので）
	ban.DoUpdateMoves(Sente)
	ban.DoUpdateMoves(Gote)
}

func (ban TBan) DoUpdateMoves(teban TTeban) {
	teban_koma := ban.GetTebanKoma(teban)
	for _, koma := range *teban_koma {
		if koma.Kind == Gyoku {
			// TODO 同じロジックを使えるようになったはず。
			// 駒から、その駒の機械的な利き先を取得する
			all_moves := koma.GetAllMoves()
			// 機械的な利き先のうち自陣営の駒がいるマスを除き、有効な指し手となるマスを保存する
			ban.CheckMovesAndKiki(koma, all_moves)
			ban.AllMoves[koma.Id] = all_moves.DeleteInvalidMoves()
			logger := GetLogger()
			logger.Trace("DoUpdateMoves is sente: " + s(teban))
			break
		}
	}
}

func (ban TBan) DeleteSuicideMoves() {
	// 両陣営の玉について、自殺手を削除する
	ban.DoDeleteSuicideMoves(Sente)
	ban.DoDeleteSuicideMoves(Gote)
}

func (ban TBan) DoDeleteSuicideMoves(teban TTeban) {
	teban_koma := ban.GetTebanKoma(teban)
	for _, koma := range *teban_koma {
		if koma.Kind == Gyoku {
			for _, move := range ban.AllMoves[koma.Id].Map {
				kiki := ban.AllMasu[move.ToPosition].GetAiteKiki(teban)
				if len(*kiki) > 0 {
					// TODO 当たっている利きが遠利きかどうか確認する処理も必要。
					move.IsValid = false
					logger := GetLogger()
					logger.Trace("DoDeleteSuicideMoves is sente: " + s(teban))
				}
			}
			ban.AllMoves[koma.Id] = ban.AllMoves[koma.Id].DeleteInvalidMoves()
			break
		}
	}
}

// 7g -> 7+7i
func str2Position(str string) TPosition {
	int_x, _ := strconv.Atoi(str[0:1])
	byte_x := byte(int_x)
	char_y := str[1:2]
	byte_y := byte(strings.Index("0abcdefghi", char_y))
	return Bytes2TPosition(byte_x, byte_y)
}

// S* -> 銀、先手
func str2KindAndTeban(str string) (TKind, TTeban) {
	char := str[0:1]
	index := strings.Index("PLNSGBRplnsgbr", char)
	teban := TTeban(index < 7)
	var kind TKind
	switch index {
	case 0, 7:
		kind = Fu
	case 1, 8:
		kind = Kyo
	case 2, 9:
		kind = Kei
	case 3, 10:
		kind = Gin
	case 4, 11:
		kind = Kin
	case 5, 12:
		kind = Kaku
	case 6, 13:
		kind = Hi
	}
	return kind, teban
}

func position2str(pos TPosition) string {
	int_x := int(real(pos))
	int_y := int(imag(pos))
	str := "0abcdefghi"
	return s(int_x) + str[int_y:int_y+1]
}

func (ban TBan) DoMove(from TPosition, to TPosition, promote bool) {
	logger := GetLogger()
	// fromにある駒を取得
	from_masu := ban.AllMasu[from]
	from_koma := ban.AllKoma[from_masu.KomaId]
	if from_koma == nil {
		// 盤と手、どちらかがおかしい
		logger.Trace("ERROR!! no Koma exists at: " + s(from))
		return
	}

	// fromにある手を取得
	moves := ban.AllMoves[from_koma.Id]
	var move *TMove = nil
	for _, value := range moves.Map {
		if value.ToPosition == to {
			move = value
			break
		}
	}
	if move == nil {
		// 盤と手、どちらかがおかしい
		logger.Trace("ERROR!! no Move exists to: " + s(to))
		return
	} else {
		if move.ToId == 0 {
			// 相手の駒を取る手ではない
		} else {
			// 相手の駒を取る
			capture_koma := ban.AllKoma[move.ToId]
			if capture_koma.Position != to {
				logger.Trace("ERROR!! capture_koma id " + s(move.ToId) + " is at: " + s(capture_koma.Position))
			} else {
				ban.CaptureKoma(move.ToId)
			}
		}
	}

	// fromにある駒をいったん盤から取り除く
	ban.RemoveKoma(from_masu.KomaId)

	// 取り除いた駒をtoに置く
	from_koma.Position = to
	// 駒が成る場合
	if promote {
		from_koma.Promoted = true
	}
	ban.PutKoma(from_koma)
}

// 駒を取る
func (ban TBan) CaptureKoma(koma_id TKomaId) {
	// 取られる駒
	target_koma := ban.AllKoma[koma_id]

	// 駒の合法手を削除
	target_masu := ban.AllMasu[target_koma.Position]
	ban.AllMoves[koma_id] = NewMoves()

	// 駒の利きを削除
	// 駒からgetAllMoveで全利き候補を取ってそこから削除するのが論理的か。
	// 少なくとも、前段の合法手には味方への利きが含まれていない。
	ban.DeleteAllKiki(target_koma)

	// 駒の持ち主を入れ替える
	if target_koma.IsSente {
		target_koma.IsSente = Gote
		delete(ban.SenteKoma, koma_id)
		ban.GoteKoma[koma_id] = target_koma
	} else {
		target_koma.IsSente = Sente
		delete(ban.GoteKoma, koma_id)
		ban.SenteKoma[koma_id] = target_koma
	}
	// 駒のあった場所からIdを削除
	target_masu.KomaId = 0
	// 駒の場所を持ち駒とする
	target_koma.Position = Mochigoma

	// 成りフラグをoff
	target_koma.Promoted = false
}

// 駒を移動させる際の移動元からの削除
func (ban TBan) RemoveKoma(koma_id TKomaId) {
	// 移動させる駒
	target_koma := ban.AllKoma[koma_id]

	// 駒の合法手を削除
	target_masu := ban.AllMasu[target_koma.Position]
	ban.AllMoves[koma_id] = NewMoves()

	// 駒の利きを削除
	ban.DeleteAllKiki(target_koma)

	// 駒のあった場所からIdを削除
	target_masu.KomaId = 0

	// 駒がどいたことにより、そのマスに利かせていた駒の利きを再評価
	ban.RefreshMovesAndKiki(target_masu, Sente, target_koma.IsSente)
	ban.RefreshMovesAndKiki(target_masu, Gote, target_koma.IsSente)
}

func (ban TBan) RefreshMovesAndKiki(masu *TMasu, kiki_teban TTeban, removed_koma_teban TTeban) {
	kiki := masu.GetKiki(kiki_teban)
	for kiki_koma_id, _ := range *kiki {
		kiki_koma := ban.AllKoma[kiki_koma_id]
		if kiki_koma.CanFarMove() {
			// 利きがFarMoveによるものかそうでないか判断しにくいので、手も利きもいったん全削除→全追加
			ban.DeleteAllKiki(kiki_koma)
			ban.AllMoves[kiki_koma_id] = ban.CreateFarMovesAndKiki(kiki_koma)
		} else {
			if kiki_teban == removed_koma_teban {
				// 利きを元に、どいたマスへの手を追加する
				m := NewMove(kiki_koma_id, kiki_koma.Position, masu.Position, 0)
				ban.AllMoves[kiki_koma_id].Add(m)
				if !kiki_koma.Promoted {
					can_promote, promote_move := m.CanPromote(kiki_koma.IsSente)
					if can_promote {
						ban.AllMoves[kiki_koma_id].Add(promote_move)
					}
				}
			} else {
				// どいた駒が敵陣営の場合、手は元々あるので、追加する必要はないが、その駒を取れなくなる。
				for _,move := range ban.AllMoves[kiki_koma_id].Map {
					if move.ToPosition == masu.Position {
						move.ToId = 0
					}
					break
				}
			}
		}
	}
}

func (ban TBan) DeleteAllKiki(koma *TKoma) {
	// 駒の利きを削除
	// 駒からgetAllMoveで全利き候補を取ってそこから削除するのが論理的か。
	// 少なくとも、前段の合法手には味方への利きが含まれていない。
	moves_4_delete_kiki := koma.GetAllMoves()
	for _, move := range moves_4_delete_kiki.Map {
		kiki_masu := ban.AllMasu[move.ToPosition]
		kiki_masu.DeleteKiki(koma.Id, koma.IsSente)
	}
}

func (ban TBan) DoDrop(teban TTeban, kind TKind, to TPosition) {
	// 打つ駒を特定する
	koma := ban.FindKoma(teban, kind)
	if koma == nil {
		return
	}
	// 打つ駒に座標を設定する
	koma.Position = to
	// 駒を配置する
	ban.PutKoma(koma)
}

func (ban TBan) FindKoma(teban TTeban, kind TKind) *TKoma {
	teban_koma_map := ban.GetTebanKoma(teban)
	var found *TKoma = nil
	for _, koma := range *teban_koma_map {
		if koma.Kind == kind && koma.Position == Mochigoma {
			found = koma
			break
		}
	}
	if found == nil {
		logger := GetLogger()
		logger.Trace("ERROR!! koma to drop is not found.")
	}
	return found
}

func AddMove(moves *map[byte]*TMove, move *TMove) {
	// logger := GetLogger()
	var i byte = 0
	for ; ; i++ {
		_, exists := (*moves)[i]
		if !exists {
			(*moves)[i] = move
			// logger.Trace("AddMove[" + move.Display() + "]")
			break
		}
	}
}

func (ban TBan) Display() string {
	var str string = ""
	// display ban
	str += "後手の持ち駒："
	for _, koma := range ban.GoteKoma {
		if koma.Position == Mochigoma {
			str += disp_map[koma.Kind] + ", "
		}
	}
	str += "\n"
	var x, y byte = 9, 1
	for y <= 9 {
		x = 9
		for x >= 1 {
			koma_id := ban.AllMasu[Bytes2TPosition(x, y)].KomaId
			if koma_id == 0 {
				str += "[＿＿]"
			} else {
				koma := ban.AllKoma[koma_id]
				str += "[" + koma.Display() + "]"
			}
			x--
		}
		str += "\n"
		y++
	}
	str += "先手の持ち駒："
	for _, koma := range ban.SenteKoma {
		if koma.Position == Mochigoma {
			str += disp_map[koma.Kind] + ", "
		}
	}
	str += "\n"
	// display move
	var k TKomaId = 1
	for ; k <= 40; k++ {
		koma := ban.AllKoma[k]
		str += koma.Display()
		str += " id:"
		str += s(koma.Id)
		str += ", position:"
		str += s(koma.Position)
		str += ", move:"
		if koma.Position.IsValidMove() {
			moves := ban.AllMoves[koma.Id]
			for _, move := range moves.Map {
				str += s(move.ToPosition)
				if move.Promote {
					str += "+"
				}
				str += ", "
			}
		}
		str += "\n"
	}
	// display kiki
	var xx, yy byte = 9, 1
	for yy <= 9 {
		xx = 9
		for xx >= 1 {
			pos := Bytes2TPosition(xx, yy)
			str += "masu: "
			str += s(pos)
			str += " kiki: "
			masu := ban.AllMasu[pos]
			for k, _ := range *(masu.SenteKiki) {
				str += ban.AllKoma[k].Display()
				str += ", "
			}
			for k, _ := range *(masu.GoteKiki) {
				str += ban.AllKoma[k].Display()
				str += ", "
			}
			str += "\n"
			xx--
		}
		yy++
	}
	return str
}
