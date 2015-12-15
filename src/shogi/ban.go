package shogi

import (
	"fmt"
	// . "logger"
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
	masu := TMasu{
		KomaId:    koma_id,
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

	// ここでの目的
	// ・新しい駒を配置する
	// ・合法手の更新（各駒に、全部の合法手が洗い出されていること）
	// 　　新しい駒の合法手を洗い出すこと★遮られた時の処理
	// 　　新しい駒により、既存の駒の合法手を更新すること
	// 　　　★新しい駒、既存の駒問わず、合法手をリフレッシュする機能が必要！
	// ・利きマスの更新（各マスに、両陣営からの利き元をマーキングできていること）
	// 　　新しい駒から、利いているマスにマーキングできていること
	// 　　新しい駒が、遮った時でもマーキングできていること
	// 　　　★新しい駒、既存の駒問わず、マーキングをリフレッシュする機能が必要！

	// 以下デバッグ表示
	/**
	logger := GetLogger()
	var str string
	str += koma.Display()
	str += " id:"
	str += s(koma.Id)
	str += ", position:"
	str += s(koma.Position)
	**/

	// 駒から、その駒の機械的な利き先を取得する
	all_moves := *(koma.getAllMove())

	// 機械的な利き先のうち自陣営の駒がいるマスを除き、有効な指し手となるマスを保存する
	// 利き先と指し手の違いを意識すること。自陣営の駒がいるマスに利いていても、指せない。
	if len(all_moves) > 0 {
		for _, move := range all_moves {
			temp_pos := move.getToAsComplex()
			target_id := ban.AllMasu[temp_pos].KomaId
			if target_id != 0 {
				if ban.AllKoma[target_id].IsSente == koma.IsSente {
					if move.IsValid {
						// 利き先マスに、自駒のIdを保存する
						saveKiki(ban.AllMasu[temp_pos], koma.Id, koma.IsSente)
					}
					// 自陣営の駒がいるマスには指せない
					move.IsValid = false
					if koma.canFarMove() {
						// 香、角、飛の場合、その先のマスへの利き先もdelete要
						// 現在地と、削除したマスの座標から方向を特定し、削除していく。
						markInvalidMoves(koma.Position, temp_pos, &all_moves)
					}
				} else {
					// 指し手の取る駒
					move.ToId = target_id
					if move.IsValid {
						// 利き先マスに、自駒のIdを保存する
						saveKiki(ban.AllMasu[temp_pos], koma.Id, koma.IsSente)
					}
					if koma.canFarMove() {
						// 香、角、飛の場合、その先のマスへの利き先もdelete要
						// 現在地と、取れる駒の座標から方向を特定し、削除していく。
						markInvalidMoves(koma.Position, temp_pos, &all_moves)
					}
				}
			} else {
				// 利き先マスに、自駒のIdを保存する
				saveKiki(ban.AllMasu[temp_pos], koma.Id, koma.IsSente)
			}
		}
	}
	valid_moves := deleteInvalidMoves(&all_moves)
	ban.AllMasu[koma.Position].Moves = valid_moves

	// 以下デバッグ表示
	/**
	str += ", move:"
	if len(*valid_moves) > 0 {
		var index byte = 0
		for index < byte(len(*valid_moves)) {
			item := (*valid_moves)[index]
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
	**/

	// logger.Trace(str)
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
				deleteKiki(target_koma.Position, koma.Position, &ban, true)

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
				deleteKiki(target_koma.Position, koma.Position, &ban, false)
			}
			valid_moves := deleteInvalidMoves(target_moves)
			ban.AllMasu[target_koma.Position].Moves = valid_moves
		}
	}

}

func deleteKiki(from complex64, to complex64, ban *TBan, is_sente bool) {
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

// コピペの嵐、リファクタ要！！！
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

func saveKiki(masu *TMasu, koma_id TKomaId, is_sente bool) {
	if is_sente {
		masu.SenteKiki[koma_id] = ""
	} else {
		masu.GoteKiki[koma_id] = ""
	}
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

func (ban TBan) Display() string {
	var str string = ""
	// display ban
	var x, y byte = 9, 1
	for y <= 9 {
		x = 9
		for x >= 1 {
			koma_id := ban.AllMasu[getComplex64(x, y)].KomaId
			if koma_id == 0 {
				str += "[   ]"
			} else {
				str += "[" + ban.AllKoma[koma_id].Display() + "]"
			}
			x--
		}
		str += "\n"
		y++
	}
	// display move
	for _, koma := range ban.AllKoma {
		str += koma.Display()
		str += " id:"
		str += s(koma.Id)
		str += ", position:"
		str += s(koma.Position)
		str += ", move:"
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
