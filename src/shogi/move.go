package shogi

import ()

type TMove struct {
	FromId       TKomaId
	FromPosition TPosition
	ToPosition   TPosition
	ToId         TKomaId
	IsValid      bool
}

func NewMove(from_id TKomaId, from_position TPosition, to_position TPosition, to_id TKomaId) *TMove {
	move := TMove{
		FromId:       from_id,
		FromPosition: from_position,
		ToPosition:   to_position,
		ToId:         to_id,
		IsValid:      true,
	}
	return &move
}

func (move TMove) Display() string {
	return "FromId: " + s(move.FromId) + ", FromPosition: " + s(move.FromPosition) + ", ToPosition: " + s(move.ToPosition)
}
