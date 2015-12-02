package shogi

import (
	"fmt"
)

// alias
var p = fmt.Println
var s = fmt.Sprint

type TBan struct {
	// 直感的な数字でわかるように10*10とし、0の要素は使わない。
	Masu [10][10]*TMasu
	// 駒IDをキーに、駒へのポインタを持つマップ
	AllKoma   map[TKomaId]*TKoma
	SenteKoma map[TKomaId]*TKoma
	GoteKoma  map[TKomaId]*TKoma
}

func NewBan() *TBan {
	masu := [10][10]*TMasu{}
	// マスを初期化する
	var x, y byte = 1, 1
	for y <= 9 {
		x = 1
		for x <= 9 {
			masu[x][y] = NewMasu(0)
			x++
		}
		y++
	}
	ban := TBan{
		Masu:      masu,
		AllKoma:   make(map[TKomaId]*TKoma),
		SenteKoma: make(map[TKomaId]*TKoma),
		GoteKoma:  make(map[TKomaId]*TKoma),
	}
	return &ban
}

func (ban TBan) PutKoma(koma *TKoma) {
	ban.AllKoma[koma.Id] = koma
	if koma.Side {
		ban.SenteKoma[koma.Id] = koma
	} else {
		ban.GoteKoma[koma.Id] = koma
	}
	ban.Masu[koma.Position[0]][koma.Position[1]].KomaId = koma.Id
}

type TMasu struct {
	KomaId    TKomaId
	Relations *Relations
}

func NewMasu(koma_id TKomaId) *TMasu {
	masu := TMasu{
		KomaId:    koma_id,
		Relations: nil,
	}
	return &masu
}

type Relations struct {
	SavingIds    *[]byte
	SavedIds     *[]byte
	TargetingIds *[]byte
	TargetedIds  *[]byte
}

func CreateInitialState() *TBan {
	ban := NewBan()
	// 駒を1つずつ生成する
	var koma_id TKomaId = 1

	// 後手
	var side bool = false
	// 香
	ban.PutKoma(NewKoma(koma_id, Kyo, [2]byte{1, 1}, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kyo, [2]byte{9, 1}, side))
	koma_id++
	// 桂
	ban.PutKoma(NewKoma(koma_id, Kei, [2]byte{2, 1}, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kei, [2]byte{8, 1}, side))
	koma_id++
	// 銀
	ban.PutKoma(NewKoma(koma_id, Gin, [2]byte{3, 1}, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Gin, [2]byte{7, 1}, side))
	koma_id++
	// 金
	ban.PutKoma(NewKoma(koma_id, Kin, [2]byte{4, 1}, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kin, [2]byte{6, 1}, side))
	koma_id++
	// 王
	ban.PutKoma(NewKoma(koma_id, Gyoku, [2]byte{5, 1}, side))
	koma_id++
	// 角
	ban.PutKoma(NewKoma(koma_id, Kaku, [2]byte{2, 2}, side))
	koma_id++
	// 飛
	ban.PutKoma(NewKoma(koma_id, Hi, [2]byte{8, 2}, side))
	koma_id++
	// 歩
	var x byte = 1
	for x <= 9 {
		ban.PutKoma(NewKoma(koma_id, Fu, [2]byte{x, 3}, side))
		koma_id++
		x++
	}

	// 先手
	side = true
	// 香
	ban.PutKoma(NewKoma(koma_id, Kyo, [2]byte{1, 9}, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kyo, [2]byte{9, 9}, side))
	koma_id++
	// 桂
	ban.PutKoma(NewKoma(koma_id, Kei, [2]byte{2, 9}, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kei, [2]byte{8, 9}, side))
	koma_id++
	// 銀
	ban.PutKoma(NewKoma(koma_id, Gin, [2]byte{3, 9}, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Gin, [2]byte{7, 9}, side))
	koma_id++
	// 金
	ban.PutKoma(NewKoma(koma_id, Kin, [2]byte{4, 9}, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kin, [2]byte{6, 9}, side))
	koma_id++
	// 王
	ban.PutKoma(NewKoma(koma_id, Gyoku, [2]byte{5, 9}, side))
	koma_id++
	// 角
	ban.PutKoma(NewKoma(koma_id, Kaku, [2]byte{8, 8}, side))
	koma_id++
	// 飛
	ban.PutKoma(NewKoma(koma_id, Hi, [2]byte{2, 8}, side))
	koma_id++
	// 歩
	x = 1
	for x <= 9 {
		ban.PutKoma(NewKoma(koma_id, Fu, [2]byte{x, 7}, side))
		koma_id++
		x++
	}

	// 各マスの関係を更新する
	return ban
}

func (ban TBan) Display() string {
	var str string = ""
	var x, y byte = 1, 1
	for y <= 9 {
		x = 1
		for x <= 9 {
			if ban.Masu[x][y].KomaId == 0 {
				str = str + "[   ]"
			} else {
				str = str + "[" + ban.AllKoma[ban.Masu[x][y].KomaId].Display() + "]"
			}
			x++
		}
		str = str + "\n"
		y++
	}
	return str
}
