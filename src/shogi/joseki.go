package shogi

type TJoseki struct {
	FixOpening map[int]*TMove
}

func NewJoseki() *TJoseki {
	fix_opening := CreateFixOpening()
	joseki := TJoseki{
		FixOpening: fix_opening,
	}
	return &joseki
}

func CreateFixOpening() map[int]*TMove {
	m := make(map[int]*TMove)
	// 先手の場合、初手78飛
	// moveを作るために、ダミーの駒が必要。（直したい）
	m[1] = NewMove(NewKoma(1, Hi, 2, 8, Sente), TPosition(complex(7, 8)), 0)
	// 後手の場合、初手32飛
	m[2] = NewMove(NewKoma(2, Hi, 8, 2, Sente), TPosition(complex(3, 2)), 0)
	return m
}
