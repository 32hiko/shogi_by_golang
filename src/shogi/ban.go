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
			all_masu[getComplex64(x, y)] = NewMasu(0)
			x++
		}
		y++
	}
	// 持ち駒用
	all_masu[complex(0, 0)] = NewMasu(0)

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
	// 駒があれば駒のId
	KomaId TKomaId
	// 駒があれば駒の合法手。駒同士の関係は、必ず盤（マス）を介する作りとする。
	Moves *map[byte]*TMove
	// このマスに利かせている駒のIdを入れる。ヒートマップを作るため
	SenteKiki map[TKomaId]string // temp
	GoteKiki  map[TKomaId]string // temp
}

func NewMasu(koma_id TKomaId) *TMasu {
	moves := make(map[byte]*TMove)
	masu := TMasu{
		KomaId:    koma_id,
		Moves:     &moves,
		SenteKiki: make(map[TKomaId]string),
		GoteKiki:  make(map[TKomaId]string),
	}
	return &masu
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

	if koma.canFarMove() {
		far_moves := ban.CreateFarMovesAndKiki(koma)
		ban.AllMasu[koma.Position].Moves = far_moves
	} else {
		// 駒から、その駒の機械的な利き先を取得する
		all_moves := koma.getAllMove()
		// 機械的な利き先のうち自陣営の駒がいるマスを除き、有効な指し手となるマスを保存する
		ban.CheckMovesAndKiki(koma, all_moves)
		valid_moves := deleteInvalidMoves(all_moves)
		ban.AllMasu[koma.Position].Moves = valid_moves
	}

	// 自マスに、他の駒からの利きとしてIdが入っている場合で、香、角、飛の場合は先の利きを止める
	s_map := ban.AllMasu[koma.Position].SenteKiki
	if len(s_map) > 0 {
		for koma_id, _ := range s_map {
			target_koma := ban.AllKoma[koma_id]
			target_moves := ban.AllMasu[target_koma.Position].Moves
			if koma.IsSente {
				// komaが先手なら、komaの位置への手は合法でなくなるので削除が必要。
				for _, move := range *target_moves {
					if move.getToAsComplex() == koma.Position {
						move.IsValid = false
						if target_koma.canFarMove() {
							markInvalidMoves(target_koma.Position, koma.Position, target_moves)
						}
						break
					}
				}
			} else {
				// komaが後手なら、komaの位置への手でkomaが取れることを手に保存する。
				for _, move := range *target_moves {
					if move.getToAsComplex() == koma.Position {
						move.ToId = koma.Id
						if target_koma.canFarMove() {
							markInvalidMoves(target_koma.Position, koma.Position, target_moves)
						}
						break
					}
				}
			}
			if target_koma.canFarMove() {
				// komaが先手後手問わず、香、角、飛からkomaへの延長線上への利きは削除が必要。
				deleteFarKiki(target_koma.Position, koma.Position, &ban, true)

			}
			valid_moves := deleteInvalidMoves(target_moves)
			ban.AllMasu[target_koma.Position].Moves = valid_moves
		}
	}
	g_map := ban.AllMasu[koma.Position].GoteKiki
	if len(g_map) > 0 {
		for koma_id, _ := range g_map {
			target_koma := ban.AllKoma[koma_id]
			target_moves := ban.AllMasu[target_koma.Position].Moves
			if koma.IsSente {
				// komaが先手なら、komaの位置への手でkomaが取れることを手に保存する。
				for _, move := range *target_moves {
					if move.getToAsComplex() == koma.Position {
						move.ToId = koma.Id
						if target_koma.canFarMove() {
							markInvalidMoves(target_koma.Position, koma.Position, target_moves)
						}
						break
					}
				}
			} else {
				// komaが後手なら、komaの位置への手は合法でなくなるので削除が必要。
				for _, move := range *target_moves {
					if move.getToAsComplex() == koma.Position {
						move.IsValid = false
						if target_koma.canFarMove() {
							markInvalidMoves(target_koma.Position, koma.Position, target_moves)
						}
						break
					}
				}
			}
			if target_koma.canFarMove() {
				// komaが先手後手問わず、香、角、飛からkomaへの延長線上への利きは削除が必要。
				deleteFarKiki(target_koma.Position, koma.Position, &ban, false)
			}
			valid_moves := deleteInvalidMoves(target_moves)
			ban.AllMasu[target_koma.Position].Moves = valid_moves
		}
	}

}

func saveKiki(masu *TMasu, koma_id TKomaId, is_sente bool) {
	if is_sente {
		masu.SenteKiki[koma_id] = ""
	} else {
		masu.GoteKiki[koma_id] = ""
	}
}

func deleteFarKiki(from complex64, to complex64, ban *TBan, is_sente bool) {
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
			for {
				temp_to += complex(1, 1)
				if isValidMove(temp_to) {
					temp_masu := ban.AllMasu[temp_to]
					if is_sente {
						delete(temp_masu.SenteKiki, source_masu.KomaId)
					} else {
						delete(temp_masu.GoteKiki, source_masu.KomaId)
					}
				} else {
					break
				}
			}
		} else if (x > 0) && (y < 0) {
			// 左上
			temp_to := to
			for {
				temp_to += complex(1, -1)
				if isValidMove(temp_to) {
					temp_masu := ban.AllMasu[temp_to]
					if is_sente {
						delete(temp_masu.SenteKiki, source_masu.KomaId)
					} else {
						delete(temp_masu.GoteKiki, source_masu.KomaId)
					}
				} else {
					break
				}
			}
		} else if (x < 0) && (y > 0) {
			// 右下
			temp_to := to
			for {
				temp_to += complex(-1, 1)
				if isValidMove(temp_to) {
					temp_masu := ban.AllMasu[temp_to]
					if is_sente {
						delete(temp_masu.SenteKiki, source_masu.KomaId)
					} else {
						delete(temp_masu.GoteKiki, source_masu.KomaId)
					}
				} else {
					break
				}
			}
		} else if (x < 0) && (y < 0) {
			// 右上
			temp_to := to
			for {
				temp_to += complex(-1, -1)
				if isValidMove(temp_to) {
					temp_masu := ban.AllMasu[temp_to]
					if is_sente {
						delete(temp_masu.SenteKiki, source_masu.KomaId)
					} else {
						delete(temp_masu.GoteKiki, source_masu.KomaId)
					}
				} else {
					break
				}
			}
		}
	} else {
		// 香、飛の場合、上、下、左、右のどれかを削除する
		if (x == 0) && (y > 0) {
			// 下
			temp_to := to
			for {
				temp_to += complex(0, 1)
				if isValidMove(temp_to) {
					temp_masu := ban.AllMasu[temp_to]
					if is_sente {
						delete(temp_masu.SenteKiki, source_masu.KomaId)
					} else {
						delete(temp_masu.GoteKiki, source_masu.KomaId)
					}
				} else {
					break
				}
			}
		} else if (x == 0) && (y < 0) {
			// 上
			temp_to := to
			for {
				temp_to += complex(0, -1)
				if isValidMove(temp_to) {
					temp_masu := ban.AllMasu[temp_to]
					if is_sente {
						delete(temp_masu.SenteKiki, source_masu.KomaId)
					} else {
						delete(temp_masu.GoteKiki, source_masu.KomaId)
					}
				} else {
					break
				}
			}
		} else if (x > 0) && (y == 0) {
			// 左
			temp_to := to
			for {
				temp_to += complex(1, 0)
				if isValidMove(temp_to) {
					temp_masu := ban.AllMasu[temp_to]
					if is_sente {
						delete(temp_masu.SenteKiki, source_masu.KomaId)
					} else {
						delete(temp_masu.GoteKiki, source_masu.KomaId)
					}
				} else {
					break
				}
			}
		} else if (x < 0) && (y == 0) {
			// 右
			temp_to := to
			for {
				temp_to += complex(-1, 0)
				if isValidMove(temp_to) {
					temp_masu := ban.AllMasu[temp_to]
					if is_sente {
						delete(temp_masu.SenteKiki, source_masu.KomaId)
					} else {
						delete(temp_masu.GoteKiki, source_masu.KomaId)
					}
				} else {
					break
				}
			}
		}
	}
}

func (ban TBan) CreateFarMovesAndKiki(koma *TKoma) *map[byte]*TMove {
	moves := make(map[byte]*TMove)
	var i byte = 0
	switch koma.Kind {
	case Kyo:
		createNMovesAndKiki(&ban, koma, move_n, &i, &moves)
	case Kaku:
		createNMovesAndKiki(&ban, koma, move_ne, &i, &moves)
		createNMovesAndKiki(&ban, koma, move_se, &i, &moves)
		createNMovesAndKiki(&ban, koma, move_nw, &i, &moves)
		createNMovesAndKiki(&ban, koma, move_sw, &i, &moves)
	case Hi:
		createNMovesAndKiki(&ban, koma, move_n, &i, &moves)
		createNMovesAndKiki(&ban, koma, move_s, &i, &moves)
		createNMovesAndKiki(&ban, koma, move_e, &i, &moves)
		createNMovesAndKiki(&ban, koma, move_w, &i, &moves)
	}
	return &moves
}

func createNMovesAndKiki(ban *TBan, koma *TKoma, delta complex64, map_key *byte, moves *map[byte]*TMove) {
	temp_move := koma.Position
	for {
		if koma.IsSente {
			temp_move += delta
		} else {
			temp_move -= delta
		}
		if isValidMove(temp_move) {
			// 利き先マスに、自駒のIdを保存する
			saveKiki(ban.AllMasu[temp_move], koma.Id, koma.IsSente)
			target_id := ban.AllMasu[temp_move].KomaId
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
		saveKiki(ban.AllMasu[temp_pos], koma.Id, koma.IsSente)
		target_id := ban.AllMasu[temp_pos].KomaId
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
	putAllKoma(ban)
	return ban
}

func putAllKoma(ban *TBan) {
	// 駒を1つずつ生成する
	var koma_id TKomaId = 1

	// 後手
	var side bool = false
	// 香
	ban.PutKoma(NewKoma(koma_id, Kyo, 1, 1, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kyo, 9, 1, side))
	koma_id++
	// 桂
	ban.PutKoma(NewKoma(koma_id, Kei, 2, 1, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kei, 8, 1, side))
	koma_id++
	// 銀
	ban.PutKoma(NewKoma(koma_id, Gin, 3, 1, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Gin, 7, 1, side))
	koma_id++
	// 金
	ban.PutKoma(NewKoma(koma_id, Kin, 4, 1, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kin, 6, 1, side))
	koma_id++
	// 王
	ban.PutKoma(NewKoma(koma_id, Gyoku, 5, 1, side))
	koma_id++
	// 角
	ban.PutKoma(NewKoma(koma_id, Kaku, 2, 2, side))
	koma_id++
	// 飛
	ban.PutKoma(NewKoma(koma_id, Hi, 8, 2, side))
	koma_id++
	// 歩
	var x byte = 1
	for x <= 9 {
		ban.PutKoma(NewKoma(koma_id, Fu, x, 3, side))
		koma_id++
		x++
	}

	// 先手
	side = true
	// 香
	ban.PutKoma(NewKoma(koma_id, Kyo, 1, 9, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kyo, 9, 9, side))
	koma_id++
	// 桂
	ban.PutKoma(NewKoma(koma_id, Kei, 2, 9, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kei, 8, 9, side))
	koma_id++
	// 銀
	ban.PutKoma(NewKoma(koma_id, Gin, 3, 9, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Gin, 7, 9, side))
	koma_id++
	// 金
	ban.PutKoma(NewKoma(koma_id, Kin, 4, 9, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kin, 6, 9, side))
	koma_id++
	// 王
	ban.PutKoma(NewKoma(koma_id, Gyoku, 5, 9, side))
	koma_id++
	// 角
	ban.PutKoma(NewKoma(koma_id, Kaku, 8, 8, side))
	koma_id++
	// 飛
	ban.PutKoma(NewKoma(koma_id, Hi, 2, 8, side))
	koma_id++
	// 歩
	x = 1
	for x <= 9 {
		ban.PutKoma(NewKoma(koma_id, Fu, x, 7, side))
		koma_id++
		x++
	}
}

// 7g -> 7+7i
func str2Position(str string) complex64 {
	int_x, _ := strconv.Atoi(str[0:1])
	float_x := float32(int_x)
	char_y := str[1:2]
	float_y := float32(strings.Index("0abcdefghi", char_y))
	return complex(float_x, float_y)
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

func (ban TBan) DoMove(from complex64, to complex64, promote bool) {
	// fromにある駒を取得
	from_masu := ban.AllMasu[from]
	from_koma := ban.AllKoma[from_masu.KomaId]
	if from_koma == nil {
		// 盤と手、どちらかがおかしい
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
	moves_4_delete_kiki := target_koma.getAllMove()
	for _, move := range *moves_4_delete_kiki {
		kiki_masu := ban.AllMasu[move.getToAsComplex()]
		if target_koma.IsSente {
			delete(kiki_masu.SenteKiki, koma_id)
		} else {
			delete(kiki_masu.GoteKiki, koma_id)
		}
	}

	// 駒の持ち主を入れ替える
	if target_koma.IsSente {
		target_koma.IsSente = false
		delete(ban.SenteKoma, koma_id)
		ban.GoteKoma[koma_id] = target_koma
	} else {
		target_koma.IsSente = true
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
	// 駒からgetAllMoveで全利き候補を取ってそこから削除するのが論理的か。
	// 少なくとも、前段の合法手には味方への利きが含まれていない。
	moves_4_delete_kiki := target_koma.getAllMove()
	for _, move := range *moves_4_delete_kiki {
		kiki_masu := ban.AllMasu[move.getToAsComplex()]
		if target_koma.IsSente {
			delete(kiki_masu.SenteKiki, koma_id)
		} else {
			delete(kiki_masu.GoteKiki, koma_id)
		}
	}
	// 駒のあった場所からIdを削除
	target_masu.KomaId = 0

	// 駒がどいたことにより、そのマスに利かせていた駒の利きを再評価
	for sente_kiki_id, _ := range target_masu.SenteKiki {
		kiki_koma := ban.AllKoma[sente_kiki_id]
		moves := ban.AllMasu[kiki_koma.Position].Moves
		if kiki_koma.canFarMove() {
			// 利かせてた駒→取り除いた駒へのルートの延長線上のmoveを生成する
			valid_moves := ban.CreateDeltaFarMovesAndKiki(sente_kiki_id, kiki_koma.Position, target_koma.Position)
			for _, move := range *valid_moves {
				AddMove(moves, move)
			}
		} else {
			AddMove(moves, NewMove(sente_kiki_id, target_koma.Position, 0))
		}
	}
	for gote_kiki_id, _ := range target_masu.GoteKiki {
		kiki_koma := ban.AllKoma[gote_kiki_id]
		moves := ban.AllMasu[kiki_koma.Position].Moves
		// 香、角、飛の場合はその先も追加する必要があるかも？
		if kiki_koma.canFarMove() {
			// 利かせてた駒→取り除いた駒へのルートの延長線上のmoveを生成する
			valid_moves := ban.CreateDeltaFarMovesAndKiki(gote_kiki_id, kiki_koma.Position, target_koma.Position)
			for _, move := range *valid_moves {
				AddMove(moves, move)
			}
		} else {
			AddMove(moves, NewMove(gote_kiki_id, target_koma.Position, 0))
		}
	}
}

// 香、角、飛の利きマスの、どいたマスからの手を再作成する
func (ban TBan) CreateDeltaFarMovesAndKiki(koma_id TKomaId, from complex64, to complex64) *map[byte]*TMove {
	far_moves := make(map[byte]*TMove)
	// 作成する方向を求める
	diff := to - from
	x := real(diff)
	y := imag(diff)
	var i byte = 0
	koma := ban.AllKoma[koma_id]

	// 完全にコピペ。
	if (x != 0) && (y != 0) {
		// 角の場合、左上、右上、左下、右下のどれかを作成する
		if (x > 0) && (y > 0) {
			// 左下
			temp_to := to
			for {
				if isValidMove(temp_to) {
					saveKiki(ban.AllMasu[temp_to], koma_id, koma.IsSente)
					target_id := ban.AllMasu[temp_to].KomaId
					if target_id != 0 {
						// 駒がある
						target_koma := ban.AllKoma[target_id]
						if target_koma.IsSente == koma.IsSente {
							// 自陣営の駒のあるマスには指せない
						} else {
							// 相手の駒は取れる。その先には動けない
							far_moves[i] = NewMove(koma.Id, temp_to, target_koma.Id)
							break
						}
					} else {
						// 駒がないなら指せて、その先をまた確認する
						far_moves[i] = NewMove(koma.Id, temp_to, 0)
						i++
					}
				} else {
					break
				}
				// どいたマスそのものを評価する
				temp_to += complex(1, 1)
			}
		} else if (x > 0) && (y < 0) {
			// 左上
			temp_to := to
			for {
				if isValidMove(temp_to) {
					saveKiki(ban.AllMasu[temp_to], koma_id, koma.IsSente)
					target_id := ban.AllMasu[temp_to].KomaId
					if target_id != 0 {
						// 駒がある
						target_koma := ban.AllKoma[target_id]
						if target_koma.IsSente == koma.IsSente {
							// 自陣営の駒のあるマスには指せない
						} else {
							// 相手の駒は取れる。その先には動けない
							far_moves[i] = NewMove(koma.Id, temp_to, target_koma.Id)
							break
						}
					} else {
						// 駒がないなら指せて、その先をまた確認する
						far_moves[i] = NewMove(koma.Id, temp_to, 0)
						i++
					}
				} else {
					break
				}
				// どいたマスそのものを評価する
				temp_to += complex(1, -1)
			}
		} else if (x < 0) && (y > 0) {
			// 右下
			temp_to := to
			for {
				if isValidMove(temp_to) {
					saveKiki(ban.AllMasu[temp_to], koma_id, koma.IsSente)
					target_id := ban.AllMasu[temp_to].KomaId
					if target_id != 0 {
						// 駒がある
						target_koma := ban.AllKoma[target_id]
						if target_koma.IsSente == koma.IsSente {
							// 自陣営の駒のあるマスには指せない
						} else {
							// 相手の駒は取れる。その先には動けない
							far_moves[i] = NewMove(koma.Id, temp_to, target_koma.Id)
							break
						}
					} else {
						// 駒がないなら指せて、その先をまた確認する
						far_moves[i] = NewMove(koma.Id, temp_to, 0)
						i++
					}
				} else {
					break
				}
				// どいたマスそのものを評価する
				temp_to += complex(-1, 1)
			}
		} else if (x < 0) && (y < 0) {
			// 右上
			temp_to := to
			for {
				if isValidMove(temp_to) {
					saveKiki(ban.AllMasu[temp_to], koma_id, koma.IsSente)
					target_id := ban.AllMasu[temp_to].KomaId
					if target_id != 0 {
						// 駒がある
						target_koma := ban.AllKoma[target_id]
						if target_koma.IsSente == koma.IsSente {
							// 自陣営の駒のあるマスには指せない
						} else {
							// 相手の駒は取れる。その先には動けない
							far_moves[i] = NewMove(koma.Id, temp_to, target_koma.Id)
							break
						}
					} else {
						// 駒がないなら指せて、その先をまた確認する
						far_moves[i] = NewMove(koma.Id, temp_to, 0)
						i++
					}
				} else {
					break
				}
				// どいたマスそのものを評価する
				temp_to += complex(-1, -1)
			}
		}
	} else {
		// 香、飛の場合、上、下、左、右のどれかを削除する
		if (x == 0) && (y > 0) {
			// 下
			temp_to := to
			for {
				if isValidMove(temp_to) {
					saveKiki(ban.AllMasu[temp_to], koma_id, koma.IsSente)
					target_id := ban.AllMasu[temp_to].KomaId
					if target_id != 0 {
						// 駒がある
						target_koma := ban.AllKoma[target_id]
						if target_koma.IsSente == koma.IsSente {
							// 自陣営の駒のあるマスには指せない
						} else {
							// 相手の駒は取れる。その先には動けない
							far_moves[i] = NewMove(koma.Id, temp_to, target_koma.Id)
							break
						}
					} else {
						// 駒がないなら指せて、その先をまた確認する
						far_moves[i] = NewMove(koma.Id, temp_to, 0)
						i++
					}
				} else {
					break
				}
				// どいたマスそのものを評価する
				temp_to += complex(0, 1)
			}
		} else if (x == 0) && (y < 0) {
			// 上
			temp_to := to
			for {
				if isValidMove(temp_to) {
					saveKiki(ban.AllMasu[temp_to], koma_id, koma.IsSente)
					target_id := ban.AllMasu[temp_to].KomaId
					if target_id != 0 {
						// 駒がある
						target_koma := ban.AllKoma[target_id]
						if target_koma.IsSente == koma.IsSente {
							// 自陣営の駒のあるマスには指せない
						} else {
							// 相手の駒は取れる。その先には動けない
							far_moves[i] = NewMove(koma.Id, temp_to, target_koma.Id)
							break
						}
					} else {
						// 駒がないなら指せて、その先をまた確認する
						far_moves[i] = NewMove(koma.Id, temp_to, 0)
						i++
					}
				} else {
					break
				}
				// どいたマスそのものを評価する
				temp_to += complex(0, -1)
			}
		} else if (x > 0) && (y == 0) {
			// 左
			temp_to := to
			for {
				if isValidMove(temp_to) {
					saveKiki(ban.AllMasu[temp_to], koma_id, koma.IsSente)
					target_id := ban.AllMasu[temp_to].KomaId
					if target_id != 0 {
						// 駒がある
						target_koma := ban.AllKoma[target_id]
						if target_koma.IsSente == koma.IsSente {
							// 自陣営の駒のあるマスには指せない
						} else {
							// 相手の駒は取れる。その先には動けない
							far_moves[i] = NewMove(koma.Id, temp_to, target_koma.Id)
							break
						}
					} else {
						// 駒がないなら指せて、その先をまた確認する
						far_moves[i] = NewMove(koma.Id, temp_to, 0)
						i++
					}
				} else {
					break
				}
				// どいたマスそのものを評価する
				temp_to += complex(1, 0)
			}
		} else if (x < 0) && (y == 0) {
			// 右
			temp_to := to
			for {
				if isValidMove(temp_to) {
					saveKiki(ban.AllMasu[temp_to], koma_id, koma.IsSente)
					target_id := ban.AllMasu[temp_to].KomaId
					if target_id != 0 {
						// 駒がある
						target_koma := ban.AllKoma[target_id]
						if target_koma.IsSente == koma.IsSente {
							// 自陣営の駒のあるマスには指せない
						} else {
							// 相手の駒は取れる。その先には動けない
							far_moves[i] = NewMove(koma.Id, temp_to, target_koma.Id)
							break
						}
					} else {
						// 駒がないなら指せて、その先をまた確認する
						far_moves[i] = NewMove(koma.Id, temp_to, 0)
						i++
					}
				} else {
					break
				}
				// どいたマスそのものを評価する
				temp_to += complex(-1, 0)
			}
		}
	}
	return &far_moves
}

func createFarMoves(koma_id TKomaId, from complex64, to complex64) *map[byte]*TMove {
	far_moves := make(map[byte]*TMove)
	// 作成する方向を求める
	diff := to - from
	x := real(diff)
	y := imag(diff)
	var i byte = 0

	// 完全にコピペ。
	if (x != 0) && (y != 0) {
		// 角の場合、左上、右上、左下、右下のどれかを作成する
		if (x > 0) && (y > 0) {
			// 左下
			temp_to := to
			for {
				temp_to += complex(1, 1)
				if isValidMove(temp_to) {
					far_moves[i] = NewMove(koma_id, temp_to, 0)
					i++
				} else {
					break
				}
			}
		} else if (x > 0) && (y < 0) {
			// 左上
			temp_to := to
			for {
				temp_to += complex(1, -1)
				if isValidMove(temp_to) {
					far_moves[i] = NewMove(koma_id, temp_to, 0)
					i++
				} else {
					break
				}
			}
		} else if (x < 0) && (y > 0) {
			// 右下
			temp_to := to
			for {
				temp_to += complex(-1, 1)
				if isValidMove(temp_to) {
					far_moves[i] = NewMove(koma_id, temp_to, 0)
					i++
				} else {
					break
				}
			}
		} else if (x < 0) && (y < 0) {
			// 右上
			temp_to := to
			for {
				temp_to += complex(-1, -1)
				if isValidMove(temp_to) {
					far_moves[i] = NewMove(koma_id, temp_to, 0)
					i++
				} else {
					break
				}
			}
		}
	} else {
		// 香、飛の場合、上、下、左、右のどれかを削除する
		if (x == 0) && (y > 0) {
			// 下
			temp_to := to
			for {
				temp_to += complex(0, 1)
				if isValidMove(temp_to) {
					far_moves[i] = NewMove(koma_id, temp_to, 0)
					i++
				} else {
					break
				}
			}
		} else if (x == 0) && (y < 0) {
			// 上
			temp_to := to
			for {
				temp_to += complex(0, -1)
				if isValidMove(temp_to) {
					far_moves[i] = NewMove(koma_id, temp_to, 0)
					i++
				} else {
					break
				}
			}
		} else if (x > 0) && (y == 0) {
			// 左
			temp_to := to
			for {
				temp_to += complex(1, 0)
				if isValidMove(temp_to) {
					far_moves[i] = NewMove(koma_id, temp_to, 0)
					i++
				} else {
					break
				}
			}
		} else if (x < 0) && (y == 0) {
			// 右
			temp_to := to
			for {
				temp_to += complex(-1, 0)
				if isValidMove(temp_to) {
					far_moves[i] = NewMove(koma_id, temp_to, 0)
					i++
				} else {
					break
				}
			}
		}
	}
	return &far_moves
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
				str += "[" + ban.AllKoma[koma_id].Display() + "]"
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
			for k, _ := range masu.SenteKiki {
				str += ban.AllKoma[k].Display()
				str += ", "
			}
			for k, _ := range masu.GoteKiki {
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
