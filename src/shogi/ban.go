package shogi

import (
	"fmt"
)

type TBan struct {
	// 直感的な数字でわかるように10*10とし、0の要素は使わない。
	Masu [10][10]*TMasu
	// 駒IDをキーに、駒へのポインタを持つマップ
	SenteKoma map[byte]*TKoma
	GoteKoma  map[byte]*TKoma
}

func NewBan() *TBan {
	masu := [10][10]*TMasu{}
	// マスを初期化する
	var x, y byte = 1, 1
	for x <= 9 {
		y = 1
		for y <= 9 {
			masu[x][y] = NewMasu(0)
			y++
		}
		x++
	}
	ban := TBan{
		Masu:      masu,
		SenteKoma: make(map[byte]*TKoma),
		GoteKoma:  make(map[byte]*TKoma),
	}
	return &ban
}

type TMasu struct {
	KomaId    byte
	Relations *Relations
}

func NewMasu(koma_id byte) *TMasu {
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
	fmt.Println("CreateInitialState()")
	ban := NewBan()
	// 駒を1つずつ生成する
	var koma_id byte = 1
	// まずは後手の歩9枚だけ
	var side bool = false
	var i byte = 1
	for i <= 9 {
		pos := [2]byte{i, 3}
		fu := NewFu(koma_id, pos, side)
		ban.GoteKoma[koma_id] = fu
		masu := ban.Masu[i][3]
		masu.KomaId = koma_id
		koma_id++
		i++
	}

	// 各マスの関係を更新する
	return ban
}

func (ban TBan) Display() string {
	var str string = ""
	var x, y byte = 1, 1
	for x <= 9 {
		y = 1
		for y <= 9 {
			if ban.Masu[x][y].KomaId != 0 {
				str = str + "(" + fmt.Sprint(x) + "," + fmt.Sprint(y) + ")=" + fmt.Sprint(ban.Masu[x][y].KomaId)
			}
			y++
		}
		x++
	}
	return str
}
