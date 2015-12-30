package shogi

import ()

type TKomaId byte
type TKind byte

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

// 将棋だけど東西南北で。直接画面には出ないし。
var move_n complex64 = complex(0, -1)
var move_s complex64 = complex(0, 1)
var move_e complex64 = complex(-1, 0)
var move_w complex64 = complex(1, 0)
var move_ne complex64 = move_n + move_e
var move_nw complex64 = move_n + move_w
var move_se complex64 = move_s + move_e
var move_sw complex64 = move_s + move_w
var move_kei_e complex64 = complex(-1, -2)
var move_kei_w complex64 = complex(1, -2)

// 何マス先でも進める系は、ロジックで。

var move_to_map = map[TKind][]complex64{
	Fu:    []complex64{move_n},
	Kei:   []complex64{move_kei_e, move_kei_w},
	Gin:   []complex64{move_n, move_ne, move_nw, move_se, move_sw},
	Kin:   []complex64{move_n, move_ne, move_nw, move_e, move_w, move_s},
	Gyoku: []complex64{move_n, move_ne, move_nw, move_e, move_w, move_s, move_se, move_sw},
}

func (kind TKind) toString(promoted bool) string {
	if promoted {
		return promoted_disp_map[kind]
	} else {
		return disp_map[kind]
	}
}

type TKoma struct {
	Id       TKomaId
	Kind     TKind
	Position complex64
	IsSente  bool
	Promoted bool
}

// 駒の生成は対局開始前にやればいいので変換とかやってもいいでしょう
func NewKoma(id TKomaId, kind TKind, x byte, y byte, isSente bool) *TKoma {
	koma := TKoma{
		Id:       id,
		Kind:     kind,
		Position: complex(float32(x), float32(y)),
		IsSente:  isSente,
		Promoted: false,
	}
	return &koma
}

func (koma TKoma) Display() string {
	var side_str string
	if koma.IsSente {
		side_str = "▲"
	} else {
		side_str = "△"
		if koma.Kind == Gyoku {
			return "△王"
		}
	}
	return side_str + koma.Kind.toString(koma.Promoted)
}

// 他の駒関係なく、盤上で移動できる先を洗い出す
func (koma TKoma) GetAllMove() *map[byte]*TMove {
	all_move := make(map[byte]*TMove)
	var i byte = 0
	if koma.Promoted {
		switch koma.Kind {
		case Kaku:
			koma.CreateNMoves(move_ne, &i, &all_move)
			koma.CreateNMoves(move_se, &i, &all_move)
			koma.CreateNMoves(move_nw, &i, &all_move)
			koma.CreateNMoves(move_sw, &i, &all_move)
			koma.Create1Move(move_n, &i, &all_move)
			koma.Create1Move(move_s, &i, &all_move)
			koma.Create1Move(move_e, &i, &all_move)
			koma.Create1Move(move_w, &i, &all_move)
		case Hi:
			koma.CreateNMoves(move_n, &i, &all_move)
			koma.CreateNMoves(move_s, &i, &all_move)
			koma.CreateNMoves(move_e, &i, &all_move)
			koma.CreateNMoves(move_w, &i, &all_move)
			koma.Create1Move(move_ne, &i, &all_move)
			koma.Create1Move(move_se, &i, &all_move)
			koma.Create1Move(move_nw, &i, &all_move)
			koma.Create1Move(move_sw, &i, &all_move)
		default:
			// と、杏、圭、全
			moves := move_to_map[Kin]
			for _, pos := range moves {
				koma.Create1Move(pos, &i, &all_move)
			}
		}
	} else {
		switch koma.Kind {
		case Kyo:
			koma.CreateNMoves(move_n, &i, &all_move)
		case Kaku:
			koma.CreateNMoves(move_ne, &i, &all_move)
			koma.CreateNMoves(move_se, &i, &all_move)
			koma.CreateNMoves(move_nw, &i, &all_move)
			koma.CreateNMoves(move_sw, &i, &all_move)
		case Hi:
			koma.CreateNMoves(move_n, &i, &all_move)
			koma.CreateNMoves(move_s, &i, &all_move)
			koma.CreateNMoves(move_e, &i, &all_move)
			koma.CreateNMoves(move_w, &i, &all_move)
		default:
			// 歩、桂、銀、金、玉
			moves := move_to_map[koma.Kind]
			for _, pos := range moves {
				koma.Create1Move(pos, &i, &all_move)
			}
		}
	}
	return &all_move
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

func (koma TKoma) CreateNMoves(move complex64, i *byte, moves *map[byte]*TMove) {
	temp_move := koma.Position
	for {
		if koma.IsSente {
			temp_move += move
		} else {
			temp_move -= move
		}
		if isValidMove(temp_move) {
			(*moves)[*i] = NewMove(koma.Id, temp_move, 0)
			*i++
		} else {
			return
		}
	}
}

func (koma TKoma) Create1Move(move complex64, i *byte, moves *map[byte]*TMove) {
	temp_move := koma.Position
	if koma.IsSente {
		temp_move += move
	} else {
		temp_move -= move
	}
	if isValidMove(temp_move) {
		(*moves)[*i] = NewMove(koma.Id, temp_move, 0)
		*i++
	}
}

func isValidMove(pos complex64) bool {
	var x byte = byte(real(pos))
	var y byte = byte(imag(pos))
	return (0 < x) && (x < 10) && (0 < y) && (y < 10)
}