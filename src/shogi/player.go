package shogi

import ()

type IPlayer interface {
	Search(*TBan) string
}

func NewPlayer(name string) IPlayer {
	switch name {
		case "Slide":
			return NewSlidePlayer()
		case "Random":
			return NewRandomPlayer()
		default:
			return nil
	}
}

/*
 * ただ飛車を左右に動かすだけ
 */
type TSlidePlayer struct {
	i *int
}

func NewSlidePlayer() *TSlidePlayer {
	var count int = 0
	player := TSlidePlayer{}
	player.i = &count
	return &player
}

func (player TSlidePlayer) Search(ban *TBan) string {
	var te string
	if *(player.i)%2 == 0 {
		te = "bestmove 8b7b"
	} else {
		te = "bestmove 7b8b"
	}
	*(player.i)++
	return te
}

/*
 * ランダム指し
 */
type TRandomPlayer struct {
}

func NewRandomPlayer() *TRandomPlayer {
	player := TRandomPlayer{}
	return &player
}

func (player TRandomPlayer) Search(ban *TBan) string {
	var te string = "bestmove "
	// ここを実装する
	return te
}