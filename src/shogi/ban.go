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
	AllMasu        map[TPosition]*TMasu
	AllKoma        map[TKomaId]*TKoma
	AllMoves       map[TKomaId]*TMoves
	SenteKoma      map[TKomaId]*TKoma
	GoteKoma       map[TKomaId]*TKoma
	SenteMochigoma *TMochigoma
	GoteMochigoma  *TMochigoma
	Teban          *TTeban
	Tesuu          *int
	EmptyMasu      []TPosition
	FuDropSente    []byte
	FuDropGote     []byte
	LastMoveTo     *TPosition
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
	tesuu := 0
	teban := Sente

	ban := TBan{
		AllMasu:        all_masu,
		AllKoma:        make(map[TKomaId]*TKoma),
		AllMoves:       make(map[TKomaId]*TMoves),
		SenteKoma:      make(map[TKomaId]*TKoma),
		GoteKoma:       make(map[TKomaId]*TKoma),
		SenteMochigoma: NewMochigoma(),
		GoteMochigoma:  NewMochigoma(),
		Teban:          &teban,
		Tesuu:          &tesuu,
		LastMoveTo:     nil,
	}
	return &ban
}
func FromSFEN(sfen string) *TBan {
	// 例：lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1
	// -は両者持ち駒がない場合。ある場合は、S2Pb3pのように表記。（先手銀1歩2、後手角1歩3）最後の数字は手数。
	split_str := strings.Split(sfen, " ")
	ban := NewBan()

	// 手番
	teban := TTeban(strings.Index("bw", split_str[1]) == 0)
	*(ban.Teban) = teban

	// 盤面
	ban.PutSFENKoma(split_str[0])

	// 持ち駒
	ban.SetSFENMochigoma(split_str[2])

	// 手数
	tesuu := 0
	if len(split_str) > 3 {
		tesuu, _ = strconv.Atoi(split_str[3])
	}
	*(ban.Tesuu) = tesuu

	p("teban: " + s(teban))
	p("tesuu: " + s(tesuu))
	return ban
}

func (ban TBan) ToSFEN(need_tesuu bool) string {
	// 例：lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1
	var str string = ""

	// 盤上
	var empties int = 0
	var x, y byte = 9, 1
	for y <= 9 {
		x = 9
		for x >= 1 {
			koma_id := ban.AllMasu[Bytes2TPosition(x, y)].KomaId
			if koma_id == 0 {
				empties++
			} else {
				koma := ban.AllKoma[koma_id]
				if empties > 0 {
					str += s(empties)
					empties = 0
				}
				k_tmp := koma.GetUSIDropString()
				if !koma.IsSente {
					k_tmp = strings.ToLower(k_tmp)
				}
				if koma.Promoted {
					k_tmp = "+" + k_tmp
				}
				str += k_tmp
			}
			x--
		}
		if empties > 0 {
			str += s(empties)
			empties = 0
		}
		str += "/"
		y++
	}
	str += " "

	// 手番
	if *(ban.Teban) == Sente {
		str += "b"
	} else {
		str += "w"
	}
	str += " "

	// 持ち駒
	var mochi_str string = ""
	for kind, count := range ban.SenteMochigoma.Map {
		if count != 0 {
			if count != 1 {
				mochi_str += s(count)
			}
			mochi_str += kind.GetUSIKind()
		}
	}
	for kind, count := range ban.GoteMochigoma.Map {
		if count != 0 {
			if count != 1 {
				mochi_str += s(count)
			}
			mochi_str += strings.ToLower(kind.GetUSIKind())
		}
	}
	if mochi_str == "" {
		str += "-"
	} else {
		str += mochi_str
	}

	// 手数
	if need_tesuu {
		str += " "
		str += s(*(ban.Tesuu))
	}
	return str
}

// 駒が持つデータ、マスが持つデータは今後も検討要
type TMasu struct {
	// マスの座標
	Position TPosition
	// 駒があれば駒のId
	KomaId TKomaId
	// このマスに利かせている駒のIdを入れる。ヒートマップを作るため
	SenteKiki *map[TKomaId]TKiki // temp
	GoteKiki  *map[TKomaId]TKiki // temp
}

func NewMasu(position TPosition, koma_id TKomaId) *TMasu {
	s_kiki := make(map[TKomaId]TKiki)
	g_kiki := make(map[TKomaId]TKiki)
	masu := TMasu{
		Position:  position,
		KomaId:    koma_id,
		SenteKiki: &s_kiki,
		GoteKiki:  &g_kiki,
	}
	return &masu
}

func (masu TMasu) SaveKiki(koma_id TKomaId, is_sente TTeban, tesuu int) {
	kiki := masu.GetKiki(is_sente)
	(*kiki)[koma_id] = TKiki(tesuu)
}

func (masu TMasu) DeleteKiki(koma_id TKomaId, is_sente TTeban) {
	kiki := masu.GetKiki(is_sente)
	delete(*kiki, koma_id)
}

func (masu TMasu) GetKiki(is_sente TTeban) *map[TKomaId]TKiki {
	if is_sente {
		return masu.SenteKiki
	} else {
		return masu.GoteKiki
	}
}

func (masu TMasu) GetAiteKiki(is_sente TTeban) *map[TKomaId]TKiki {
	if is_sente {
		return masu.GoteKiki
	} else {
		return masu.SenteKiki
	}
}

type TKiki int

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

func (ban TBan) GetFuDrop(teban TTeban) []byte {
	if teban {
		return ban.FuDropSente
	} else {
		return ban.FuDropGote
	}
}

func (ban TBan) GetMochigoma(teban TTeban) *TMochigoma {
	if teban {
		return ban.SenteMochigoma
	} else {
		return ban.GoteMochigoma
	}
}

func (ban TBan) SetSFENMochigoma(sfen_mochigoma string) {
	// 1文字ずつチェックする。
	var count int = 0
	for i := 0; i < len(sfen_mochigoma); i++ {
		char := sfen_mochigoma[i : i+1]
		// まず-かどうか
		if char == "-" {
			// 持ち駒なし、明示的に初期化が必要であればここですること
			return
		}
		num := strings.Index("0123456789", char)
		if num == -1 {
			// 数字ではないので、その駒を持っている。
			kind, teban := str2KindAndTeban(char)
			if count == 0 {
				count = 1
			}
			// 持ち駒表を更新
			target_mochigoma := *(ban.GetMochigoma(teban))
			target_mochigoma.Map[kind] = count
			// 持ち駒を生成
			for j := 0; j < count; j++ {
				koma_id := TKomaId(len(ban.AllKoma) + 1)
				new_mochigoma := NewKoma(koma_id, kind, 0, 0, teban)
				ban.AllKoma[koma_id] = new_mochigoma
				taban_koma := *(ban.GetTebanKoma(teban))
				taban_koma[koma_id] = new_mochigoma
			}
			count = 0
		} else {
			// 次の文字が駒であることが確定。枚数を取得して次の文字をチェックする
			if count != 0 {
				// まずないはずだが、歩を10枚以上持っている場合。
				count = count*10 + num
			} else {
				count = num
			}
		}
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

func (ban TBan) PutSFENKoma(sfen string) {
	arr := strings.Split(sfen, "/")
	var y byte = 1
	var x byte = 9
	var koma_id TKomaId = 1
	for _, line := range arr {
		x = 9
		promote := false
		// 1文字ずつチェックする。
		for i := 0; i < len(line); i++ {
			char := line[i : i+1]
			// まず数字かどうか
			num := strings.Index("0123456789", char)
			if num == -1 {
				// 数字ではないので駒が存在するマス。
				plus := strings.Index("+", char)
				if plus == 0 {
					// +は次の文字が成り駒であることを意味する。
					promote = true
					continue
				}
				kind, teban := str2KindAndTeban(char)
				koma := NewKoma(koma_id, kind, x, y, teban)
				if promote {
					koma.Promoted = true
					promote = false
				}
				ban.PutKoma(koma)
				koma_id++
				x--
			} else {
				// 空きマス分飛ばす
				x -= byte(num)
			}
		}
		y++
	}
}

// 駒を配置し、合法手、利きマスデータを更新する
func (ban TBan) PutKoma(koma *TKoma) {
	//logger := GetLogger()
	//logger.Trace("PutKoma id: " + s(koma.Id))
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
	ban.AllMoves[koma.Id] = ban.CreateFarMovesAndKiki(koma)

	// 自マスに、他の駒からの利きとしてIdが入っている場合で、香、角、飛の場合は先の利きを止める
	// こちらも、龍や馬の周囲に駒を打った場合、的外れな方向の手や利きを消そうとしてしまうので、作り直しとする
	ban.DeleteCloseMovesAndKiki(koma, Sente)
	ban.DeleteCloseMovesAndKiki(koma, Gote)
}

func (ban TBan) DeleteCloseMovesAndKiki(koma *TKoma, is_sente TTeban) {
	// logger := GetLogger()
	// logger.Trace("DeleteCloseMovesAndKiki id: " + s(koma.Id) + ", sente?: " + s(is_sente))
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
						}
					}
					if !saved {
						// もともと自陣営の駒に利かせていたところ、その駒を取られた場合はmoveが存在していない。
						ban.AddNewMoves(target_koma, koma.Position, koma.Id)
					}
				} else {
					// komaが自陣営なら、komaの位置への手は合法でなくなるので削除が必要。
					for _, move := range ban.AllMoves[koma_id].Map {
						if move.ToPosition == koma.Position {
							move.IsValid = false
							// 成らず、成る手両方を削除する必要があるのでbreakできない
						}
					}
				}
				ban.AllMoves[koma_id] = ban.AllMoves[koma_id].DeleteInvalidMoves()
			}
		}
	}
}

func (ban TBan) AddNewMoves(from_koma *TKoma, to_pos TPosition, to_id TKomaId) {
	move := NewMove(from_koma, to_pos, to_id)
	if !from_koma.Promoted {
		// 成っていない駒が成れる場合は成る手を追加
		can_promote, promote_move := move.CanPromote(from_koma.IsSente)
		if can_promote {
			ban.AllMoves[from_koma.Id].Add(promote_move)
		}
		// 行き場のない手でなければ追加
		if from_koma.CanMove(to_pos) {
			ban.AllMoves[from_koma.Id].Add(move)
		}
	} else {
		ban.AllMoves[from_koma.Id].Add(move)
	}
}

func AddNewMoves2Slice(slice *[]*TMove, from_koma *TKoma, to_pos TPosition, to_id TKomaId) {
	move := NewMove(from_koma, to_pos, to_id)
	if !from_koma.Promoted {
		// 成っていない駒が成れる場合は成る手を追加
		can_promote, promote_move := move.CanPromote(from_koma.IsSente)
		if can_promote {
			*slice = append(*slice, promote_move)
		}
		// 行き場のない手でなければ追加
		if from_koma.CanMove(to_pos) {
			*slice = append(*slice, move)
		}
	} else {
		*slice = append(*slice, move)
	}
}

// 手と利きを生成する。
func (ban TBan) CreateFarMovesAndKiki(koma *TKoma) *TMoves {
	// logger := GetLogger()
	// logger.Trace("CreateFarMovesAndKiki id: " + s(koma.Id))
	moves := NewMoves()
	if koma.Promoted {
		switch koma.Kind {
		case Kaku:
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_ne))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_se))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_nw))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_sw))
			m, _ := ban.Create1MoveAndKiki(koma, move_n, false)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_s, false)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_e, false)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_w, false)
			moves.AddAll(m)
		case Hi:
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_n))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_s))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_e))
			moves.AddAll(ban.CreateNMovesAndKiki(koma, move_w))
			m, _ := ban.Create1MoveAndKiki(koma, move_ne, false)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_se, false)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_nw, false)
			moves.AddAll(m)
			m, _ = ban.Create1MoveAndKiki(koma, move_sw, false)
			moves.AddAll(m)
		default:
			// と、杏、圭、全
			deltas := move_to_map[Kin]
			for _, delta := range deltas {
				m, _ := ban.Create1MoveAndKiki(koma, delta, false)
				moves.AddAll(m)
			}
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
		default:
			// 歩、桂、銀、金、玉
			deltas := move_to_map[koma.Kind]
			for _, delta := range deltas {
				m, _ := ban.Create1MoveAndKiki(koma, delta, false)
				moves.AddAll(m)
			}
		}
	}
	return moves
}

// 1マス分だけの手と利きを生成する。
func (ban TBan) Create1MoveAndKiki(koma *TKoma, delta TPosition, is_far bool) ([]*TMove, bool) {
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
		kiki_masu.SaveKiki(koma.Id, koma.IsSente, *(ban.Tesuu))
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
				AddNewMoves2Slice(&slice, koma, to_pos, target_koma.Id)
				// 王手の場合で遠利きの場合、利きを1マスだけ貫通させてみるテスト
				if target_koma.Kind == Gyoku && is_far {
					var saki TPosition
					if koma.IsSente {
						saki = to_pos + (delta.Vector())
					} else {
						saki = to_pos - (delta.Vector())
					}
					if saki.IsValidMove() {
						saki_masu := ban.AllMasu[saki]
						saki_masu.SaveKiki(koma.Id, koma.IsSente, *(ban.Tesuu))
					}
				}
				return slice, false
			}
		} else {
			// 駒がないなら指せる
			AddNewMoves2Slice(&slice, koma, to_pos, 0)
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
		moves, keep := ban.Create1MoveAndKiki(koma, delta_base, true)
		slice = append(slice, moves...)
		if keep {
			delta_base += delta
		} else {
			break
		}
	}
	return slice
}

// USI形式のmoveを反映させる。
func (ban *TBan) ApplyMove(usi_move string) {
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
	var to TPosition

	// これから反映する手数
	*(ban.Tesuu) += 1

	// logger := GetLogger()
	// 駒を打つかどうか
	is_drop := strings.Index(from_str, "*")
	if is_drop == -1 {
		// 打たない
		from := str2Position(from_str)
		to = str2Position(to_str)

		// logger.Trace("from: " + s(from) + ", to: " + s(to))
		// こちらのmoveを実行する
		ban.DoMove(from, to, promote)
	} else {
		// "*"を含む＝打つ。先手の銀打ちならS*,後手の銀打ちならs*で始め、打つマスの表記は同じ。
		// のはずだが、将棋所では先後問わず駒の種類が大文字になっている模様。
		kind, _ := str2KindAndTeban(from_str)
		// その手当て
		teban := *(ban.Teban)
		to = str2Position(to_str)

		// logger.Trace("駒打: " + teban_map[teban] + disp_map[kind] + ", to: " + s(to))
		ban.DoDrop(teban, kind, to)
	}

	// 最後の手を保存しておく
	ban.LastMoveTo = &to
	// 前に生成した、打つ手を全部リセットする。
	ban.ResetDropMoves()
	// 持ち駒がないなら、これを省略する。
	// 駒を打つ手を生成するために、空いているマスや二歩のチェックをする
	ban.CheckEmptyMasu()
	// 打つ手を生成する
	ban.CreateAllMochigomaMoves()
	// 指し手の反映が終わり、相手の手番に
	*(ban.Teban) = !*(ban.Teban)
}

func (ban TBan) ResetDropMoves() {
	for koma_id, koma := range ban.AllKoma {
		if koma.Position == Mochigoma {
			ban.AllMoves[koma_id] = NewMoves()
		}
	}
}

func (ban *TBan) CheckEmptyMasu() {
	logger := GetLogger()
	empty_masu := make([]TPosition, 81)
	fu_drop_sente := make([]byte, 0)
	fu_drop_gote := make([]byte, 0)
	var x, y byte = 1, 1
	for x <= 9 {
		fu_sente := false
		fu_gote := false
		y = 1
		for y <= 9 {
			pos := Bytes2TPosition(x, y)
			// logger.Trace("CheckEmptyMasu pos: " + s(pos))
			masu := ban.AllMasu[pos]
			if masu.KomaId == 0 {
				// 空いたマスを保存
				empty_masu = append(empty_masu, pos)
				// logger.Trace("CheckEmptyMasu append. pos: " + s(pos))
			} else {
				koma := ban.AllKoma[masu.KomaId]
				if koma.Position != pos {
					// ありえないが、バグとしてありえるので予め
					logger.Trace("CheckEmptyMasu Ghost Koma Id: " + s(koma.Id))
					masu.KomaId = 0
					empty_masu = append(empty_masu, pos)
				} else {
					// その列の歩の有無をチェック
					if koma.Kind == Fu {
						if koma.IsSente {
							fu_sente = true
						} else {
							fu_gote = true
						}
					}
				}
			}
			y++
		}
		if !fu_sente {
			fu_drop_sente = append(fu_drop_sente, x)
		}
		if !fu_gote {
			fu_drop_gote = append(fu_drop_gote, x)
		}
		x++
	}
	// 独自のstructに値を保存するには、レシーバをアドレス表記にする必要がある。
	ban.EmptyMasu = empty_masu
	ban.FuDropSente = fu_drop_sente
	ban.FuDropGote = fu_drop_gote
	// logger.Trace("CheckEmptyMasu ok: " + s(empty_masu))
}

func (ban TBan) CreateAllMochigomaMoves() {
	ban.DoCreateAllMochigomaMoves(Sente)
	ban.DoCreateAllMochigomaMoves(Gote)
}

func (ban TBan) DoCreateAllMochigomaMoves(teban TTeban) {
	all_mochigoma := *(ban.GetMochigoma(teban))
	for kind, num := range all_mochigoma.Map {
		if num > 0 {
			koma := ban.FindMochiKoma(teban, kind)
			ban.AllMoves[koma.Id] = ban.CreateMochigomaMoves(koma)
		}
	}
}

func (ban TBan) CreateMochigomaMoves(koma *TKoma) *TMoves {
	moves := NewMoves()
	// 空いているマスを探す
	for _, pos := range ban.EmptyMasu {
		if koma.CanMove(pos) {
			// 歩、香、桂の場合、行き場のないマスには打てない
			if koma.Kind == Fu {
				// 歩の場合、二歩は禁止
				fu_drop := ban.GetFuDrop(koma.IsSente)
				to_x := byte(real(pos))
				ok := false
				for _, x := range fu_drop {
					if to_x == x {
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
				// TODO: 歩の場合、打ち歩詰めは禁止
			}
			moves.Add(NewMove(koma, pos, 0))
		}
	}
	return moves
}

func (ban TBan) FilterSuicideMoves(gyoku *TKoma) *TMoves {
	new_moves := NewMoves()
	// 玉の動ける先に相手の利きがないか調べ、利きがあるならその手は自殺手として削除する
	for _, move := range ban.AllMoves[gyoku.Id].Map {
		kiki := ban.AllMasu[move.ToPosition].GetAiteKiki(gyoku.IsSente)
		// TODO 当たっている利きが遠利きかどうか確認する処理も必要。
		// 通常、利きは貫通しないが、玉の場合は貫通するようにしておけば、このロジックでもいいかな？
		// →現状、遠利きについてはそれでもいい。
		if len(*kiki) == 0 {
			new_moves.Add(move)
		}
	}
	return new_moves
}

func (ban TBan) FilterPinnedMoves(gyoku *TKoma, moves *(map[TKomaId]*TMoves)) {
	aite_koma := ban.GetTebanKoma(!gyoku.IsSente)
	for _, koma := range *aite_koma {
		// 相手の駒のうち、遠利きのある駒を探す
		if koma.CanFarMove() && koma.Position != Mochigoma {
			all_moves := koma.GetAllMoves()
			for _, move := range all_moves.Map {
				// 王が遠利きの筋に入っている場合、間に駒があるか調べる。
				// 龍や馬の近い利きも含んでしまっているが、for文が即終了するので問題ないはず。
				if gyoku.Position == move.ToPosition {
					aida := gyoku.Position - koma.Position
					var pinned_koma *TKoma = nil
					is_pinned := false
					aida_map := make(map[TPosition]string)
					// ピンされていても、ピンしている駒を取る手は可能とする。
					aida_map[koma.Position] = ""
					// 相手の駒から王までの間を、相手の駒のとなりから調べていく。
					for p := koma.Position + aida.Vector(); p != gyoku.Position; p += aida.Vector() {
						aida_map[p] = ""
						masu := ban.AllMasu[p]
						if masu.KomaId != 0 {
							// 縛っている駒がすでにあり、次の駒があった場合、駒の縛りはない
							if is_pinned {
								pinned_koma = nil
								break
							} else {
								// 縛っている駒がない時に王の陣営の駒があったら縛る
								k := ban.AllKoma[masu.KomaId]
								if k.IsSente == gyoku.IsSente {
									pinned_koma = k
									is_pinned = true
								} else {
									// 相手の陣営の駒があったら縛りはない
									break
								}
							}
						}
					}
					if pinned_koma != nil {
						logger := GetLogger()
						logger.Trace("DoDeleteSuicideMoves pinned: " + pinned_koma.Display() + "id: " + s(pinned_koma.Id))
						pinned_moves := NewMoves()
						for _, move := range ban.AllMoves[pinned_koma.Id].Map {
							// ピンされている駒は、ピンしている駒を取る手か、利き筋の中でのみ移動できる。
							_, ok := aida_map[move.ToPosition]
							if ok {
								pinned_moves.Add(move)
							}
						}
						(*moves)[pinned_koma.Id] = pinned_moves
					}
				}
			}
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
	index := strings.Index("PLNSGBRKplnsgbrk", char)
	teban := TTeban(index < 8)
	var kind TKind
	switch index {
	case 0, 8:
		kind = Fu
	case 1, 9:
		kind = Kyo
	case 2, 10:
		kind = Kei
	case 3, 11:
		kind = Gin
	case 4, 12:
		kind = Kin
	case 5, 13:
		kind = Kaku
	case 6, 14:
		kind = Hi
	case 7, 15:
		kind = Gyoku
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
		logger.Trace("ERROR!! Koma at: " + s(from) + ",no Move exists to: " + s(to))
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
	mm := *(ban.GetMochigoma((target_koma.IsSente)))
	if mm.Map[target_koma.Kind] == 0 {
		mm.Map[target_koma.Kind] = 1
	} else {
		mm.Map[target_koma.Kind] += 1
	}

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
				ban.AddNewMoves(kiki_koma, masu.Position, 0)
			} else {
				// どいた駒が敵陣営の場合、手は元々あるので、追加する必要はないが、その駒を取れなくなる。
				for _, move := range ban.AllMoves[kiki_koma_id].Map {
					if move.ToPosition == masu.Position {
						move.ToId = 0
					}
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
	koma := ban.FindMochiKoma(teban, kind)
	if koma == nil {
		return
	}
	// 打つ駒に座標を設定する
	koma.Position = to

	mm := *(ban.GetMochigoma(teban))
	if mm.Map[koma.Kind] == 0 {
		// ありえないが、ガードしておく。
		mm.Map[koma.Kind] = 0
	} else {
		mm.Map[koma.Kind] -= 1
	}
	// 駒を配置する
	ban.PutKoma(koma)
}

func (ban TBan) FindKoma(teban TTeban, kind TKind) *(map[TKomaId]*TKoma) {
	teban_koma_map := ban.GetTebanKoma(teban)
	result := make(map[TKomaId]*TKoma)
	for _, koma := range *teban_koma_map {
		if koma.Kind == kind {
			result[koma.Id] = koma
		}
	}
	return &result
}

func (ban TBan) FindMochiKoma(teban TTeban, kind TKind) *TKoma {
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

func (ban TBan) CountKikiMasu(teban TTeban) int {
	var count int = 0
	for _, masu := range ban.AllMasu {
		kiki_map := masu.GetKiki(teban)
		if len(*kiki_map) > 0 {
			count++
		}
	}
	return count
}

func (ban TBan) Analyze() map[string]int {
	var result = make(map[string]int)
	// 利きの数
	result["Sente:kiki"] = 0
	result["Gote:kiki"] = 0
	// 利きマスの数
	result["Sente:kikiMasu"] = 0
	result["Gote:kikiMasu"] = 0
	// 駒の数
	result["Sente:koma"] = 0
	result["Gote:koma"] = 0
	// ひも付き駒の数
	result["Sente:himoKoma"] = 0
	result["Gote:himoKoma"] = 0
	// 浮き駒の数
	result["Sente:ukiKoma"] = 0
	result["Gote:ukiKoma"] = 0
	// あたりされてる駒の数
	result["Sente:atariKoma"] = 0
	result["Gote:atariKoma"] = 0
	for _, masu := range ban.AllMasu {
		sente_kiki := masu.GetKiki(Sente)
		gote_kiki := masu.GetKiki(Gote)
		sente_kiki_len := len(*sente_kiki)
		gote_kiki_len := len(*gote_kiki)
		if sente_kiki_len > 0 {
			result["Sente:kiki"] += sente_kiki_len
			result["Sente:kikiMasu"]++
		}
		if gote_kiki_len > 0 {
			result["Gote:kiki"] += gote_kiki_len
			result["Gote:kikiMasu"]++
		}
		if masu.KomaId != 0 {
			koma := ban.AllKoma[masu.KomaId]
			if koma.IsSente {
				result["Sente:koma"]++
				if sente_kiki_len > 0 {
					result["Sente:himoKoma"]++
				} else {
					result["Sente:ukiKoma"]++
				}
				if gote_kiki_len > 0 {
					result["Sente:atariKoma"]++
				}
			} else {
				result["Gote:koma"]++
				if gote_kiki_len > 0 {
					result["Gote:himoKoma"]++
				} else {
					result["Gote:ukiKoma"]++
				}
				if sente_kiki_len > 0 {
					result["Gote:atariKoma"]++
				}
			}
		}
	}
	return result
}

func (ban TBan) Display() string {
	var str string = ""
	// display ban
	str += "後手の持ち駒："
	for kind, count := range ban.GoteMochigoma.Map {
		if count != 0 {
			str += disp_map[kind] + s(count) + ", "
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
	for kind, count := range ban.SenteMochigoma.Map {
		if count != 0 {
			str += disp_map[kind] + s(count) + ", "
		}
	}
	str += "\n"
	// display move
	var debug bool = false
	if debug {
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
				for k, v := range *(masu.SenteKiki) {
					str += ban.AllKoma[k].Display()
					str += "(" + s(v) + ")"
					str += ", "
				}
				for k, v := range *(masu.GoteKiki) {
					str += ban.AllKoma[k].Display()
					str += "(" + s(v) + ")"
					str += ", "
				}
				str += "\n"
				xx--
			}
			yy++
		}
	}
	return str
}
