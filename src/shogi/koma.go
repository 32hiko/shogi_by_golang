package shogi

import ()

const (
	Fu = iota
	Kyo
	Kei
	Gin
	Kin
	Kaku
	Hisha
	Gyoku
)

type TKoma struct {
	Id       byte
	Kind     byte
	Position [2]byte
	Side     bool
	Promoted bool
	MoveTo   *[][2]byte
}

func NewFu(id byte, position [2]byte, side bool) *TKoma {
	fu := TKoma{
		Id:       id,
		Kind:     Fu,
		Position: position,
		Side:     side,
		Promoted: false,
		MoveTo:   nil,
	}
	return &fu
}
