package shogi

import (
	"fmt"
)

type TMove struct {
	FromId TKomaId
	ToX    byte
	ToY    byte
	ToId   TKomaId
	IsValid bool
}

func NewMove(from_id TKomaId, position complex64, to_id TKomaId) *TMove {
	var to_x byte = byte(real(position))
	var to_y byte = byte(imag(position))
	move := TMove{
		FromId: from_id,
		ToX:    to_x,
		ToY:    to_y,
		ToId:   to_id,
		IsValid: true,
	}
	return &move
}

func NewMove2(from_id TKomaId, to_x byte, to_y byte, to_id TKomaId) *TMove {
	move := TMove{
		FromId: from_id,
		ToX:    to_x,
		ToY:    to_y,
		ToId:   to_id,
		IsValid: true,
	}
	return &move
}

func (move TMove) getToAsComplex() complex64 {
	return complex(float32(move.ToX), float32(move.ToY))
}

func (move TMove) Display() string {
	return fmt.Sprint(move.ToX) + ", " + fmt.Sprint(move.ToY)
}
