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
	// マスの位置（複素数）をキーに、マスへのポインタを持つマップ
	AllMasu map[complex64]*TMasu
	// 駒IDをキーに、駒へのポインタを持つマップ
	AllKoma   map[TKomaId]*TKoma
	SenteKoma map[TKomaId]*TKoma
	GoteKoma  map[TKomaId]*TKoma
}

func NewBan() *TBan {
	all_masu := make(map[complex64]*TMasu)
	// マスを初期化する
	var x, y byte = 1, 1
	for y <= 9 {
		x = 1
		for x <= 9 {
			pos := getComplex64(x, y)
			all_masu[pos] = NewMasu(pos, 0)
			x++
		}
		y++
	}
	// 持ち駒用
	all_masu[complex(0, 0)] = NewMasu(complex(0, 0), 0)

	ban := TBan{
		AllMasu:   all_masu,
		AllKoma:   make(map[TKomaId]*TKoma),
		SenteKoma: make(map[TKomaId]*TKoma),
		GoteKoma:  make(map[TKomaId]*TKoma),
	}
	return &ban
}

// 駒が持つデータ、マスが持つデータは今後も検討要
type TMasu struct {
	// マスの座標
	Position complex64
	// 駒があれば駒のId
	KomaId TKomaId
	// 駒があれば駒の合法手。駒同士の関係は、必ず盤（マス）を介する作りとする。
	Moves *map[byte]*TMove
	// このマスに利かせている駒のIdを入れる。ヒートマップを作るため
	SenteKiki *map[TKomaId]string // temp
	GoteKiki  *map[TKomaId]string // temp
}

func NewMasu(position complex64, koma_id TKomaId) *TMasu {
	moves := make(map[byte]*TMove)
	s_kiki := make(map[TKomaId]string)
	g_kiki := make(map[TKomaId]string)
	masu := TMasu{
		Position:  position,
		KomaId:    koma_id,
		Moves:     &moves,
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

func getComplex64(x byte, y byte) complex64 {
	return complex(float32(x), float32(y))
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

	// 駒の合法手
	if koma.CanFarMove() {
		// 香、角、飛、馬、龍の遠利き部分。駒の有無も考慮しつつ作成する。
		far_moves := ban.CreateFarMovesAndKiki(koma)
		ban.AllMasu[koma.Position].Moves = far_moves
	} else {
		// 駒から、その駒の機械的な利き先を取得する
		all_moves := koma.GetAllMoves()
		// 機械的な利き先のうち自陣営の駒がいるマスを除き、有効な指し手となるマスを保存する
		ban.CheckMovesAndKiki(koma, all_moves)
		valid_moves := deleteInvalidMoves(all_moves)
		ban.AllMasu[koma.Position].Moves = valid_moves
	}

	// 自マスに、他の駒からの利きとしてIdが入っている場合で、香、角、飛の場合は先の利きを止める
	ban.DeleteCloseMovesAndKiki(koma, Sente)
	ban.DeleteCloseMovesAndKiki(koma, Gote)
}

func (ban TBan) DeleteCloseMovesAndKiki(koma *TKoma, is_sente TTeban) {
	kiki_map := ban.AllMasu[koma.Position].GetKiki(is_sente)
	if len(*kiki_map) > 0 {
		for koma_id, _ := range *kiki_map {
			target_koma := ban.AllKoma[koma_id]
			target_moves := ban.AllMasu[target_koma.Position].Moves
			if koma.IsSente != is_sente {
				// komaが敵陣営なら、komaの位置への手でkomaが取れることを手に保存する。
				var saved bool = false
				for _, move := range *target_moves {
					if move.getToAsComplex() == koma.Position {
						move.ToId = koma.Id
						saved = true
						if target_koma.CanFarMove() {
							markInvalidMoves(target_koma.Position, koma.Position, target_moves)
						}
						break
					}
				}
				if !saved {
					// もともと自陣営の駒に利かせていたところ、その駒を取られた場合はmoveが存在していない。
					AddMove(target_moves, NewMove(target_koma.Id, koma.Position, koma.Id))
				}
			} else {
				// komaが自陣営なら、komaの位置への手は合法でなくなるので削除が必要。
				for _, move := range *target_moves {
					if move.getToAsComplex() == koma.Position {
						move.IsValid = false
						if target_koma.CanFarMove() {
							markInvalidMoves(target_koma.Position, koma.Position, target_moves)
						}
						break
					}
				}
			}
			if target_koma.CanFarMove() {
				// komaが先手後手問わず、香、角、飛からkomaへの延長線上への利きは削除が必要。
				ban.DeleteFarKiki(target_koma.Position, koma.Position, is_sente)
			}
			valid_moves := deleteInvalidMoves(target_moves)
			ban.AllMasu[target_koma.Position].Moves = valid_moves
		}
	}
}

// 香、角、飛の、1方向の利きを遮られた時の利きを削除する
func (ban TBan) DeleteFarKiki(from complex64, to complex64, is_sente TTeban) {
	// 削除する方向を求める
	diff := to - from
	x := real(diff)
	y := imag(diff)
	source_masu := ban.AllMasu[from]

	// 完全にコピペ。
	if (x != 0) && (y != 0) {
		// 角の場合、左上、右上、左下、右下のどれかを削除する
		if (x > 0) && (y > 0) {
			// 左下
			temp_to := to
			ban.DeleteAllFarKiki(source_masu, temp_to, complex(1, 1), is_sente)
		} else if (x > 0) && (y < 0) {
			// 左上
			temp_to := to
			ban.DeleteAllFarKiki(source_masu, temp_to, complex(1, -1), is_sente)
		} else if (x < 0) && (y > 0) {
			// 右下
			temp_to := to
			ban.DeleteAllFarKiki(source_masu, temp_to, complex(-1, 1), is_sente)
		} else if (x < 0) && (y < 0) {
			// 右上
			temp_to := to
			ban.DeleteAllFarKiki(source_masu, temp_to, complex(-1, -1), is_sente)
		}
	} else {
		// 香、飛の場合、上、下、左、右のどれかを削除する
		if (x == 0) && (y > 0) {
			// 下
			temp_to := to
			ban.DeleteAllFarKiki(source_masu, temp_to, complex(0, 1), is_sente)
		} else if (x == 0) && (y < 0) {
			// 上
			temp_to := to
			ban.DeleteAllFarKiki(source_masu, temp_to, complex(0, -1), is_sente)
		} else if (x > 0) && (y == 0) {
			// 左
			temp_to := to
			ban.DeleteAllFarKiki(source_masu, temp_to, complex(1, 0), is_sente)
		} else if (x < 0) && (y == 0) {
			// 右
			temp_to := to
			ban.DeleteAllFarKiki(source_masu, temp_to, complex(-1, 0), is_sente)
		}
	}
}

// ↑の正味の処理
func (ban TBan) DeleteAllFarKiki(masu *TMasu, temp_to complex64, delta complex64, is_sente TTeban) {
	for {
		temp_to += delta
		if isValidMove(temp_to) {
			temp_masu := ban.AllMasu[temp_to]
			temp_masu.DeleteKiki(masu.KomaId, is_sente)
		} else {
			break
		}
	}
}

// 香、角、飛の、全方向の手と利きを生成する。
func (ban TBan) CreateFarMovesAndKiki(koma *TKoma) *map[byte]*TMove {
	moves := make(map[byte]*TMove)
	var i byte = 0
	if koma.Promoted {
		switch koma.Kind {
		case Kaku:
			ban.CreateNMovesAndKiki(koma, move_ne, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_se, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_nw, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_sw, &i, &moves)
			ban.Create1MoveAndKiki(koma, move_n, &i, &moves)
			ban.Create1MoveAndKiki(koma, move_s, &i, &moves)
			ban.Create1MoveAndKiki(koma, move_e, &i, &moves)
			ban.Create1MoveAndKiki(koma, move_w, &i, &moves)
		case Hi:
			ban.CreateNMovesAndKiki(koma, move_n, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_s, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_e, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_w, &i, &moves)
			ban.Create1MoveAndKiki(koma, move_ne, &i, &moves)
			ban.Create1MoveAndKiki(koma, move_se, &i, &moves)
			ban.Create1MoveAndKiki(koma, move_nw, &i, &moves)
			ban.Create1MoveAndKiki(koma, move_sw, &i, &moves)
		}
	} else {
		switch koma.Kind {
		case Kyo:
			ban.CreateNMovesAndKiki(koma, move_n, &i, &moves)
		case Kaku:
			ban.CreateNMovesAndKiki(koma, move_ne, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_se, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_nw, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_sw, &i, &moves)
		case Hi:
			ban.CreateNMovesAndKiki(koma, move_n, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_s, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_e, &i, &moves)
			ban.CreateNMovesAndKiki(koma, move_w, &i, &moves)
		}
	}
	return &moves
}

// 馬、龍の、成ってできた1マス分だけの手と利きを生成する。処理的にはNと同じ。
func (ban TBan) Create1MoveAndKiki(koma *TKoma, delta complex64, map_key *byte, moves *map[byte]*TMove) {
	temp_move := koma.Position
	if koma.IsSente {
		temp_move += delta
	} else {
		temp_move -= delta
	}
	if isValidMove(temp_move) {
		// 利き先マスに、自駒のIdを保存する
		kiki_masu := ban.AllMasu[temp_move]
		kiki_masu.SaveKiki(koma.Id, koma.IsSente)
		target_id := kiki_masu.KomaId
		// 利き先マスに、駒があるかないか
		if target_id != 0 {
			// 駒がある
			target_koma := ban.AllKoma[target_id]
			if target_koma.IsSente == koma.IsSente {
				// 自陣営の駒のあるマスには指せない
				return
			} else {
				// 相手の駒は取れる。その先には動けない
				(*moves)[*map_key] = NewMove(koma.Id, temp_move, target_koma.Id)
				*map_key++
				return
			}
		} else {
			// 駒がないなら指せて、その先をまた確認する
			(*moves)[*map_key] = NewMove(koma.Id, temp_move, 0)
			*map_key++
		}
	} else {
		return
	}
}

// 香、角、飛の、ある1方向の手と利きを生成する。
func (ban TBan) CreateNMovesAndKiki(koma *TKoma, delta complex64, map_key *byte, moves *map[byte]*TMove) {
	temp_move := koma.Position
	for {
		if koma.IsSente {
			temp_move += delta
		} else {
			temp_move -= delta
		}
		if isValidMove(temp_move) {
			// 利き先マスに、自駒のIdを保存する
			kiki_masu := ban.AllMasu[temp_move]
			kiki_masu.SaveKiki(koma.Id, koma.IsSente)
			target_id := kiki_masu.KomaId
			// 利き先マスに、駒があるかないか
			if target_id != 0 {
				// 駒がある
				target_koma := ban.AllKoma[target_id]
				if target_koma.IsSente == koma.IsSente {
					// 自陣営の駒のあるマスには指せない
					return
				} else {
					// 相手の駒は取れる。その先には動けない
					(*moves)[*map_key] = NewMove(koma.Id, temp_move, target_koma.Id)
					*map_key++
					return
				}
			} else {
				// 駒がないなら指せて、その先をまた確認する
				(*moves)[*map_key] = NewMove(koma.Id, temp_move, 0)
				*map_key++
			}
		} else {
			return
		}
	}
}

// 駒の合法手と利き先マスをチェックする（香、角、飛を除く）
func (ban TBan) CheckMovesAndKiki(koma *TKoma, moves *map[byte]*TMove) {
	for _, move := range *moves {
		temp_pos := move.getToAsComplex()
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

// コピペの嵐、リファクタ要！！！if文の回数とかも切り詰めないと
func markInvalidMoves(from complex64, to complex64, moves *map[byte]*TMove) {
	// 削除する方向を求める
	diff := to - from
	x := real(diff)
	y := imag(diff)

	if (x != 0) && (y != 0) {
		// 角の場合、左上、右上、左下、右下のどれかを削除する
		if (x > 0) && (y > 0) {
			// 左下
			for _, move := range *moves {
				temp_to := move.getToAsComplex()
				temp_diff := temp_to - to
				temp_x := real(temp_diff)
				temp_y := imag(temp_diff)
				if (temp_x > 0) && (temp_y > 0) {
					move.IsValid = false
				}
			}
		} else if (x > 0) && (y < 0) {
			// 左上
			for _, move := range *moves {
				temp_to := move.getToAsComplex()
				temp_diff := temp_to - to
				temp_x := real(temp_diff)
				temp_y := imag(temp_diff)
				if (temp_x > 0) && (temp_y < 0) {
					move.IsValid = false
				}
			}
		} else if (x < 0) && (y > 0) {
			// 右下
			for _, move := range *moves {
				temp_to := move.getToAsComplex()
				temp_diff := temp_to - to
				temp_x := real(temp_diff)
				temp_y := imag(temp_diff)
				if (temp_x < 0) && (temp_y > 0) {
					move.IsValid = false
				}
			}
		} else if (x < 0) && (y < 0) {
			// 右上
			for _, move := range *moves {
				temp_to := move.getToAsComplex()
				temp_diff := temp_to - to
				temp_x := real(temp_diff)
				temp_y := imag(temp_diff)
				if (temp_x < 0) && (temp_y < 0) {
					move.IsValid = false
				}
			}
		}
	} else {
		// 香、飛の場合、上、下、左、右のどれかを削除する
		if (x == 0) && (y > 0) {
			// 下
			for _, move := range *moves {
				temp_to := move.getToAsComplex()
				temp_diff := temp_to - to
				temp_x := real(temp_diff)
				temp_y := imag(temp_diff)
				if (temp_x == 0) && (temp_y > 0) {
					move.IsValid = false
				}
			}
		} else if (x == 0) && (y < 0) {
			// 上
			for _, move := range *moves {
				temp_to := move.getToAsComplex()
				temp_diff := temp_to - to
				temp_x := real(temp_diff)
				temp_y := imag(temp_diff)
				if (temp_x == 0) && (temp_y < 0) {
					move.IsValid = false
				}
			}
		} else if (x > 0) && (y == 0) {
			// 左
			for _, move := range *moves {
				temp_to := move.getToAsComplex()
				temp_diff := temp_to - to
				temp_x := real(temp_diff)
				temp_y := imag(temp_diff)
				if (temp_x > 0) && (temp_y == 0) {
					move.IsValid = false
				}
			}
		} else if (x < 0) && (y == 0) {
			// 右
			for _, move := range *moves {
				temp_to := move.getToAsComplex()
				temp_diff := temp_to - to
				temp_x := real(temp_diff)
				temp_y := imag(temp_diff)
				if (temp_x < 0) && (temp_y == 0) {
					move.IsValid = false
				}
			}
		}
	}

}

func deleteInvalidMoves(org *map[byte]*TMove) *map[byte]*TMove {
	deleted := make(map[byte]*TMove)
	var i byte = 0
	for _, move := range *org {
		if move.IsValid {
			deleted[i] = move
			i++
		}
	}
	return &deleted
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
	from := str2Position(from_str)
	to := str2Position(to_str)
	logger := GetLogger()
	logger.Trace("from: " + s(from) + ", to: " + s(to))
	// こちらのmoveを実行する
	ban.DoMove(from, to, promote)
}

// 7g -> 7+7i
func str2Position(str string) complex64 {
	int_x, _ := strconv.Atoi(str[0:1])
	float_x := float32(int_x)
	char_y := str[1:2]
	float_y := float32(strings.Index("0abcdefghi", char_y))
	return complex(float_x, float_y)
}

func (ban TBan) DoMove(from complex64, to complex64, promote bool) {
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
	moves := from_masu.Moves
	var move *TMove = nil
	for _, value := range *moves {
		if value.getToAsComplex() == to {
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
			ban.CaptureKoma(move.ToId)
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
	target_masu.Moves = nil

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
	target_koma.Position = complex(0, 0)

	// 成りフラグをoff
	target_koma.Promoted = false
}

// 駒を移動させる際の移動元からの削除
func (ban TBan) RemoveKoma(koma_id TKomaId) {
	// 移動させる駒
	target_koma := ban.AllKoma[koma_id]

	// 駒の合法手を削除
	target_masu := ban.AllMasu[target_koma.Position]
	target_masu.Moves = nil

	// 駒の利きを削除
	ban.DeleteAllKiki(target_koma)

	// 駒のあった場所からIdを削除
	target_masu.KomaId = 0

	// 駒がどいたことにより、そのマスに利かせていた駒の利きを再評価
	ban.RefreshMovesAndKiki(target_masu, Sente)
	ban.RefreshMovesAndKiki(target_masu, Gote)
}

func (ban TBan) RefreshMovesAndKiki(masu *TMasu, is_sente TTeban) {
	kiki := masu.GetKiki(is_sente)
	for kiki_koma_id, _ := range *kiki {
		kiki_koma := ban.AllKoma[kiki_koma_id]
		moves := ban.AllMasu[kiki_koma.Position].Moves
		if kiki_koma.CanFarMove() {
			// 利きがFarMoveによるものかそうでないか判断しにくいので、手も利きもいったん全削除→全追加
			kiki_koma_masu := ban.AllMasu[kiki_koma.Position]
			kiki_koma_masu.Moves = nil
			ban.DeleteAllKiki(kiki_koma)
			far_moves := ban.CreateFarMovesAndKiki(kiki_koma)
			kiki_koma_masu.Moves = far_moves
		} else {
			// 利きを元に、どいたマスへの手を追加する
			AddMove(moves, NewMove(kiki_koma_id, masu.Position, 0))
		}
	}
}

func (ban TBan) DeleteAllKiki(koma *TKoma) {
	// 駒の利きを削除
	// 駒からgetAllMoveで全利き候補を取ってそこから削除するのが論理的か。
	// 少なくとも、前段の合法手には味方への利きが含まれていない。
	moves_4_delete_kiki := koma.GetAllMoves()
	for _, move := range *moves_4_delete_kiki {
		kiki_masu := ban.AllMasu[move.getToAsComplex()]
		kiki_masu.DeleteKiki(koma.Id, koma.IsSente)
	}
}

func AddMove(moves *map[byte]*TMove, move *TMove) {
	var i byte = 0
	for ; ; i++ {
		_, exists := (*moves)[i]
		if !exists {
			(*moves)[i] = move
			break
		}
	}
}

func (ban TBan) Display() string {
	var str string = ""
	// display ban
	var x, y byte = 9, 1
	for y <= 9 {
		x = 9
		for x >= 1 {
			koma_id := ban.AllMasu[getComplex64(x, y)].KomaId
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
		if isValidMove(koma.Position) {
			moves := ban.AllMasu[koma.Position].Moves
			if len(*moves) > 0 {
				var index byte = 0
				for index < byte(len(*moves)) {
					item := (*moves)[index]
					if item == nil {
						index++
						continue
					}
					temp_pos := item.getToAsComplex()
					str += s(temp_pos)
					str += ", "
					index++
				}
			}
		}
		str += "\n"
	}
	// display kiki
	var xx, yy byte = 9, 1
	for yy <= 9 {
		xx = 9
		for xx >= 1 {
			pos := getComplex64(xx, yy)
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
