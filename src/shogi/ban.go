package shogi

import (
	"fmt"
	. "logger"
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

// 駒を配置する
func (ban TBan) PutKoma(koma *TKoma) {
	// 駒が持っている位置を更新
	ban.AllKoma[koma.Id] = koma
	// ここは本来、駒の所有権が決まった時点でやる処理。初期化も、まず全部持ち駒にしてそれを打っていくのが正しい。
	if koma.Side {
		ban.SenteKoma[koma.Id] = koma
	} else {
		ban.GoteKoma[koma.Id] = koma
	}
	ban.Masu[int(real(koma.Position))][int(imag(koma.Position))].KomaId = koma.Id

	// 以下デバッグ表示
	logger := GetLogger()
	var str string
	str += koma.Display()
	str += " id:"
	str += s(koma.Id)

	// 駒から、その駒の機械的な利き先を取得する
	all_move := *(koma.getAllMove())

	// 以下デバッグ表示
	str += ", position:"
	str += s(koma.Position)
	str += ", move:"
	if len(all_move) > 0 {
		var index byte = 0
		for index < byte(len(all_move)) {
			item := all_move[index]
			temp_pos := complex(float32(item.ToX), float32(item.ToY))
			str += s(temp_pos)
			str += ", "
			index++
		}
	}
	logger.Trace(str)
	// 自マスに、有効な利き先マスを保存する
	// 利き先マスに、自駒のIdを保存する
	// 自マスに、他の駒からの利きとしてIdが入っている場合で、香、角、飛、馬、龍の場合は先の利きを止める
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
	putAllKoma(ban)
	// 1マスずつ、合法手や利き等を初期化する
	// initializeRelations(ban)
	return ban
}

func putAllKoma(ban *TBan) {
	// 駒を1つずつ生成する
	var koma_id TKomaId = 1

	// 後手
	var side bool = false
	// 香
	ban.PutKoma(NewKoma(koma_id, Kyo, 1, 1, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kyo, 9, 1, side))
	koma_id++
	// 桂
	ban.PutKoma(NewKoma(koma_id, Kei, 2, 1, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kei, 8, 1, side))
	koma_id++
	// 銀
	ban.PutKoma(NewKoma(koma_id, Gin, 3, 1, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Gin, 7, 1, side))
	koma_id++
	// 金
	ban.PutKoma(NewKoma(koma_id, Kin, 4, 1, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kin, 6, 1, side))
	koma_id++
	// 王
	ban.PutKoma(NewKoma(koma_id, Gyoku, 5, 1, side))
	koma_id++
	// 角
	ban.PutKoma(NewKoma(koma_id, Kaku, 2, 2, side))
	koma_id++
	// 飛
	ban.PutKoma(NewKoma(koma_id, Hi, 8, 2, side))
	koma_id++
	// 歩
	var x byte = 1
	for x <= 9 {
		ban.PutKoma(NewKoma(koma_id, Fu, x, 3, side))
		koma_id++
		x++
	}

	// 先手
	side = true
	// 香
	ban.PutKoma(NewKoma(koma_id, Kyo, 1, 9, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kyo, 9, 9, side))
	koma_id++
	// 桂
	ban.PutKoma(NewKoma(koma_id, Kei, 2, 9, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kei, 8, 9, side))
	koma_id++
	// 銀
	ban.PutKoma(NewKoma(koma_id, Gin, 3, 9, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Gin, 7, 9, side))
	koma_id++
	// 金
	ban.PutKoma(NewKoma(koma_id, Kin, 4, 9, side))
	koma_id++
	ban.PutKoma(NewKoma(koma_id, Kin, 6, 9, side))
	koma_id++
	// 王
	ban.PutKoma(NewKoma(koma_id, Gyoku, 5, 9, side))
	koma_id++
	// 角
	ban.PutKoma(NewKoma(koma_id, Kaku, 8, 8, side))
	koma_id++
	// 飛
	ban.PutKoma(NewKoma(koma_id, Hi, 2, 8, side))
	koma_id++
	// 歩
	x = 1
	for x <= 9 {
		ban.PutKoma(NewKoma(koma_id, Fu, x, 7, side))
		koma_id++
		x++
	}
}

/*
func initializeRelations(ban *TBan) {
	// 全部の駒について、利きの範囲にいる駒（味方・敵）と合法手をチェックしていく
	// 利かされているかどうかは別途
	for id, koma := range ban.SenteKoma {

	}

}
*/

func (ban TBan) Display() string {
	var str string = ""
	var x, y byte = 1, 1
	for y <= 9 {
		x = 1
		for x <= 9 {
			if ban.Masu[x][y].KomaId == 0 {
				str += "[   ]"
			} else {
				str += "[" + ban.AllKoma[ban.Masu[x][y].KomaId].Display() + "]"
			}
			x++
		}
		str += "\n"
		y++
	}
	return str
}
