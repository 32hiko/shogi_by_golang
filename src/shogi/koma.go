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

func (kind TKind) toString() string {
	return disp_map[kind]
}

type TKoma struct {
	Id       TKomaId
	Kind     TKind
	Position [2]byte
	Side     bool
	Promoted bool
	MoveTo   *[][2]byte
}

func NewKoma(id TKomaId, kind TKind, position [2]byte, side bool) *TKoma {
	koma := TKoma{
		Id:       id,
		Kind:     kind,
		Position: position,
		Side:     side,
		Promoted: false,
		MoveTo:   nil,
	}
	return &koma
}

func (koma TKoma) Display() string {
	var side_str string
	if koma.Side {
		side_str = "▲"
	} else {
		side_str = "△"
	}
	return side_str + koma.Kind.toString()
}
