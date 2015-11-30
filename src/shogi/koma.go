package shogi

import ()

const (
	Fu = iota
	Kyo
	Kei
	Gin
	Kin
	Kaku
	Hi
	Gyoku
)

var disp_map = map[byte]string{Fu:"歩", Kyo:"香", Kei:"桂", Gin:"銀", Kin:"金", Kaku:"角", Hi:"飛", Gyoku:"玉"}

type TKoma struct {
	Id       byte
	Kind     byte
	Position [2]byte
	Side     bool
	Promoted bool
	MoveTo   *[][2]byte
}

func NewKoma(id byte, kind byte, position [2]byte, side bool) *TKoma {
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
	return side_str + disp_map[koma.Kind]
}