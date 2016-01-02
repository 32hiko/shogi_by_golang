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

func (move TMove) Display() string {
	return "FromId: " + s(move.FromId) + ", ToPosition: " + s(move.ToPosition)
}
