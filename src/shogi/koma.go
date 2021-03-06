package shogi

import ()

type TKomaId byte
type TKind byte
type TTeban bool
type TPosition complex64

const (
	Fu TKind = iota
	Kyo
	Kei
	Gin
	Kin
	Kaku
	Hi
	Gyoku
)

const (
	Sente TTeban = true
	Gote  TTeban = false
)

func tebanString(teban TTeban) string {
	if teban == Sente {
		return "Sente"
	} else {
		return "Gote"
	}
}

var teban_map = map[TTeban]string{
	Sente: "▲",
	Gote:  "△",
}

const (
	Mochigoma TPosition = complex(0, 0)
)

var disp_map = map[TKind]string{
	Fu:    "歩",
	Kyo:   "香",
	Kei:   "桂",
	Gin:   "銀",
	Kin:   "金",
	Kaku:  "角",
	Hi:    "飛",
	Gyoku: "玉",
}

var promoted_disp_map = map[TKind]string{
	Fu:   "と",
	Kyo:  "杏",
	Kei:  "圭",
	Gin:  "全",
	Kaku: "馬",
	Hi:   "龍",
}

var usi_drop_map = map[TKind]string{
	Fu:    "P",
	Kyo:   "L",
	Kei:   "N",
	Gin:   "S",
	Kin:   "G",
	Kaku:  "B",
	Hi:    "R",
	Gyoku: "K",
}

var kind_point_map = map[TKind]int{
	Fu:    1,
	Kyo:   3,
	Kei:   4,
	Gin:   5,
	Kin:   6,
	Kaku:  12,
	Hi:    14,
	Gyoku: 0,
}

// 将棋だけど東西南北で。直接画面には出ないし。
var move_n TPosition = complex(0, -1)
var move_s TPosition = complex(0, 1)
var move_e TPosition = complex(-1, 0)
var move_w TPosition = complex(1, 0)
var move_ne TPosition = move_n + move_e
var move_nw TPosition = move_n + move_w
var move_se TPosition = move_s + move_e
var move_sw TPosition = move_s + move_w
var move_kei_e TPosition = complex(-1, -2)
var move_kei_w TPosition = complex(1, -2)

// 何マス先でも進める系は、ロジックで。

var move_to_map = map[TKind][]TPosition{
	Fu:    []TPosition{move_n},
	Kei:   []TPosition{move_kei_e, move_kei_w},
	Gin:   []TPosition{move_n, move_ne, move_nw, move_se, move_sw},
	Kin:   []TPosition{move_n, move_ne, move_nw, move_e, move_w, move_s},
	Gyoku: []TPosition{move_n, move_ne, move_nw, move_e, move_w, move_s, move_se, move_sw},
}

func (kind TKind) toString(promoted bool) string {
	if promoted {
		return promoted_disp_map[kind]
	} else {
		return disp_map[kind]
	}
}

type TMochigoma struct {
	Map map[TKind]int
}

func NewMochigoma() *TMochigoma {
	mochigoma := TMochigoma{
		Map: make(map[TKind]int),
	}
	return &mochigoma
}

func (mochigoma TMochigoma) GetTotalCount() int {
	count := 0
	for _, v := range mochigoma.Map {
		count += v
	}
	return count
}

type TKoma struct {
	Id       TKomaId
	Kind     TKind
	Position TPosition
	IsSente  TTeban
	Promoted bool
}

// 駒の生成は対局開始前にやればいいので変換とかやってもいいでしょう
func NewKoma(id TKomaId, kind TKind, x byte, y byte, isSente TTeban) *TKoma {
	koma := TKoma{
		Id:       id,
		Kind:     kind,
		Position: Bytes2TPosition(x, y),
		IsSente:  isSente,
		Promoted: false,
	}
	return &koma
}

func (koma TKoma) Display() string {
	if koma.IsSente == Gote && koma.Kind == Gyoku {
		return "△王"
	}
	return teban_map[koma.IsSente] + koma.Kind.toString(koma.Promoted)
}

func (koma TKoma) CanFarMove() bool {
	if koma.Promoted {
		if koma.Kind == Kaku || koma.Kind == Hi {
			return true
		} else {
			return false
		}
	} else {
		return move_to_map[koma.Kind] == nil
	}
}

// 他の駒関係なく、盤上で移動できる先を洗い出す
func (koma TKoma) GetAllMoves() *TMoves {
	moves := NewMoves()

	if koma.Promoted {
		switch koma.Kind {
		case Kaku:
			moves.AddAll(koma.CreateFarMovesFromDelta(move_ne))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_se))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_nw))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_sw))
			moves.AddAll(koma.CreateMovesFromDelta(move_n))
			moves.AddAll(koma.CreateMovesFromDelta(move_s))
			moves.AddAll(koma.CreateMovesFromDelta(move_e))
			moves.AddAll(koma.CreateMovesFromDelta(move_w))
		case Hi:
			moves.AddAll(koma.CreateFarMovesFromDelta(move_n))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_s))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_e))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_w))
			moves.AddAll(koma.CreateMovesFromDelta(move_ne))
			moves.AddAll(koma.CreateMovesFromDelta(move_se))
			moves.AddAll(koma.CreateMovesFromDelta(move_nw))
			moves.AddAll(koma.CreateMovesFromDelta(move_sw))
		default:
			// と、杏、圭、全
			deltas := move_to_map[Kin]
			for _, delta := range deltas {
				moves.AddAll(koma.CreateMovesFromDelta(delta))
			}
		}
	} else {
		switch koma.Kind {
		case Kyo:
			moves.AddAll(koma.CreateFarMovesFromDelta(move_n))
		case Kaku:
			moves.AddAll(koma.CreateFarMovesFromDelta(move_ne))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_se))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_nw))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_sw))
		case Hi:
			moves.AddAll(koma.CreateFarMovesFromDelta(move_n))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_s))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_e))
			moves.AddAll(koma.CreateFarMovesFromDelta(move_w))
		default:
			// 歩、桂、銀、金、玉
			deltas := move_to_map[koma.Kind]
			for _, delta := range deltas {
				moves.AddAll(koma.CreateMovesFromDelta(delta))
			}
		}
	}
	return moves
}

func (koma TKoma) CreateFarMovesFromDelta(delta TPosition) []*TMove {
	slice := make([]*TMove, 0)
	delta_base := delta
	for {
		moves := koma.CreateMovesFromDelta(delta_base)
		if len(moves) > 0 {
			slice = append(slice, moves...)
			delta_base += delta
		} else {
			break
		}
	}
	return slice
}

func (koma TKoma) CreateMovesFromDelta(delta TPosition) []*TMove {
	slice := make([]*TMove, 0)
	var to_pos TPosition
	if koma.IsSente {
		to_pos = koma.Position + delta
	} else {
		to_pos = koma.Position - delta
	}
	if to_pos.IsValidMove() {
		AddNewMoves2Slice(&slice, &koma, to_pos, 0)
	}
	return slice
}

func (koma TKoma) CanMove(to_pos TPosition) bool {
	to_y := imag(to_pos)
	var can_move bool = true
	if koma.IsSente {
		if (koma.Kind == Fu) || (koma.Kind == Kyo) {
			if to_y == 1 {
				can_move = false
			}
		} else if koma.Kind == Kei {
			if to_y <= 2 {
				can_move = false
			}
		}
	} else {
		if (koma.Kind == Fu) || (koma.Kind == Kyo) {
			if to_y == 9 {
				can_move = false
			}
		} else if koma.Kind == Kei {
			if to_y >= 8 {
				can_move = false
			}
		}
	}
	return can_move
}

func (koma TKoma) GetUSIDropString() string {
	k := usi_drop_map[koma.Kind]
	// 打つときは、駒の種類はすべて大文字で。
	return k
}

func (kind TKind) GetUSIKind() string {
	return usi_drop_map[kind]
}

func (koma TKoma) GetKomaPoint() int {
	p := kind_point_map[koma.Kind]
	if koma.Promoted {
		p += 6
	}
	return p
}

func (position TPosition) IsValidMove() bool {
	x := real(position)
	y := imag(position)
	return (0 < x) && (x < 10) && (0 < y) && (y < 10)
}

func (position TPosition) Vector() TPosition {
	x := real(position)
	y := imag(position)
	if x > 0 {
		x = 1
	}
	if x < 0 {
		x = -1
	}
	if y > 0 {
		y = 1
	}
	if y < 0 {
		y = -1
	}
	return TPosition(complex(x, y))
}
