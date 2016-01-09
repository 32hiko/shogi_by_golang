package shogi

import (
	. "logger"
)

type TMoves struct {
	Map map[string]*TMove
}

func NewMoves() *TMoves {
	m := make(map[string]*TMove)
	moves := TMoves{
		Map: m,
	}
	return &moves
}

func (moves TMoves) Add(move *TMove) {
	logger := GetLogger()
	key := move.GetUSIMoveString()
	moves.Map[key] = move
	logger.Trace("Add key: " + key + ", move: [" + move.Display() + "]")
}

type TMove struct {
	FromId       TKomaId
	FromPosition TPosition
	ToPosition   TPosition
	ToId         TKomaId
	IsValid      bool
	Promote      bool
}

func NewMove(from_id TKomaId, from_position TPosition, to_position TPosition, to_id TKomaId) *TMove {
	move := TMove{
		FromId:       from_id,
		FromPosition: from_position,
		ToPosition:   to_position,
		ToId:         to_id,
		IsValid:      true,
		Promote:      false,
	}
	return &move
}

func (move TMove) CanPromote(teban TTeban) (bool, *TMove) {
	from_y := imag(move.FromPosition)
	to_y := imag(move.ToPosition)
	var can_promote bool = false
	var promote_move *TMove = nil
	if teban {
		if (from_y <= 3) || (to_y <= 3) {
			can_promote = true
		}
	} else {
		if (from_y >= 7) || (to_y >= 7) {
			can_promote = true
		}
	}
	if can_promote {
		promote_move = NewMove(move.FromId, move.FromPosition, move.ToPosition, move.ToId)
		promote_move.Promote = true
	}
	return can_promote, promote_move
}

func (move TMove) GetUSIMoveString() string {
	// TODO 打つ手（fromが0,0の場合）に対応する。
	from := move.FromPosition
	to := move.ToPosition
	return_str := position2str(from) + position2str(to)
	if move.Promote {
		return_str += "+"
	}
	return return_str
}

func (move TMove) Display() string {
	var str string = ""
	str += "FromId: " + s(move.FromId)
	str += ", FromPosition: " + s(move.FromPosition)
	str += ", ToPosition: " + s(move.ToPosition)
	str += ", Promote: " + s(move.Promote)
	return str
}
