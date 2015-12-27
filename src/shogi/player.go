package shogi

import ()

type TPlayer struct {
	i *int
}

// 当然、Playerは何種類か用意する予定。interfaceとかちゃんと使って。
func NewPlayer() *TPlayer {
	var count int = 0
	player := TPlayer{}
	player.i = &count
	return &player
}

func (player TPlayer) Search(ban *TBan) string {
	var te string
	if *(player.i)%2 == 0 {
		te = "bestmove 8b7b"
	} else {
		te = "bestmove 7b8b"
	}
	*(player.i)++
	return te
}
