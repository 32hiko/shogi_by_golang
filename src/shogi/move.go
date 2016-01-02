package shogi

import ()

type TMove struct {
	FromId     TKomaId
	ToPosition TPosition
	ToId       TKomaId
	IsValid    bool
}

func NewMove(from_id TKomaId, position TPosition, to_id TKomaId) *TMove {
	move := TMove{
		FromId:     from_id,
		ToPosition: position,
		ToId:       to_id,
		IsValid:    true,
	}
	return &move
}

// 現状、使われてない。
func (move TMove) Display() string {
	to_x := real(move.ToPosition)
	to_y := imag(move.ToPosition)
	return s(to_x) + ", " + s(to_y)
}
