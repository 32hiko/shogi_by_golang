package shogi

import (
	"testing"
)

func TestNewKoma(t *testing.T) {
	// NewKoma(id TKomaId, kind TKind, x byte, y byte, isSente TTeban)
	// test
	koma := NewKoma(1, Fu, 2, 3, Sente)

	// assert
	if koma.Id != 1 {
		t.Errorf("Id: actual:[%v] expected:[%v]", koma.Id, 1)
	}
	if koma.Kind != Fu {
		t.Errorf("Kind: actual:[%v] expected:[%v]", koma.Kind, Fu)
	}
	if koma.Position != complex(2, 3) {
		t.Errorf("Id: actual:[%v] expected:[%v]", koma.Position, complex(2, 3))
	}
	if koma.IsSente != Sente {
		t.Errorf("Id: actual:[%v] expected:[%v]", koma.IsSente, Sente)
	}
	if koma.Promoted != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", koma.Promoted, false)
	}
	p("TestNewKoma")
}

func TestDisplay(t *testing.T) {
	// prepare
	sente_fu := NewKoma(1, Fu, 1, 1, Sente)
	gote_kyo := NewKoma(2, Kyo, 1, 1, Gote)
	sente_kei := NewKoma(3, Kei, 1, 1, Sente)
	gote_gin := NewKoma(4, Gin, 1, 1, Gote)
	sente_kin := NewKoma(5, Kin, 1, 1, Sente)
	gote_kaku := NewKoma(6, Kaku, 1, 1, Gote)
	sente_hi := NewKoma(7, Hi, 1, 1, Sente)
	gote_gyoku := NewKoma(8, Gyoku, 1, 1, Gote)
	sente_gyoku := NewKoma(9, Gyoku, 1, 1, Sente)

	// test & assert
	if sente_fu.Display() != "▲歩" {
		t.Errorf("Id: actual:[%s] expected:[%s]", sente_fu.Display(), "▲歩")
	}
	if gote_kyo.Display() != "△香" {
		t.Errorf("Id: actual:[%s] expected:[%s]", gote_kyo.Display(), "△香")
	}
	if sente_kei.Display() != "▲桂" {
		t.Errorf("Id: actual:[%s] expected:[%s]", sente_kei.Display(), "▲桂")
	}
	if gote_gin.Display() != "△銀" {
		t.Errorf("Id: actual:[%s] expected:[%s]", gote_gin.Display(), "△銀")
	}
	if sente_kin.Display() != "▲金" {
		t.Errorf("Id: actual:[%s] expected:[%s]", sente_kin.Display(), "▲金")
	}
	if gote_kaku.Display() != "△角" {
		t.Errorf("Id: actual:[%s] expected:[%s]", gote_kaku.Display(), "△角")
	}
	if sente_hi.Display() != "▲飛" {
		t.Errorf("Id: actual:[%s] expected:[%s]", sente_hi.Display(), "▲飛")
	}
	if gote_gyoku.Display() != "△王" {
		t.Errorf("Id: actual:[%s] expected:[%s]", gote_gyoku.Display(), "△王")
	}
	if sente_gyoku.Display() != "▲玉" {
		t.Errorf("Id: actual:[%s] expected:[%s]", sente_gyoku.Display(), "▲玉")
	}
	p("TestDisplay")
}

func TestCanFarMove(t *testing.T) {
	// prepare
	sente_fu := NewKoma(1, Fu, 1, 1, Sente)
	gote_kyo := NewKoma(2, Kyo, 1, 1, Gote)
	sente_kei := NewKoma(3, Kei, 1, 1, Sente)
	gote_gin := NewKoma(4, Gin, 1, 1, Gote)
	sente_kin := NewKoma(5, Kin, 1, 1, Sente)
	gote_kaku := NewKoma(6, Kaku, 1, 1, Gote)
	sente_hi := NewKoma(7, Hi, 1, 1, Sente)
	gote_gyoku := NewKoma(8, Gyoku, 1, 1, Gote)
	gote_fu := NewKoma(9, Fu, 1, 1, Gote)
	gote_fu.Promoted = true
	sente_kyo := NewKoma(10, Kyo, 1, 1, Sente)
	sente_kyo.Promoted = true
	gote_kei := NewKoma(11, Kei, 1, 1, Gote)
	gote_kei.Promoted = true
	sente_gin := NewKoma(12, Gin, 1, 1, Sente)
	sente_gin.Promoted = true
	gote_hi := NewKoma(13, Hi, 1, 1, Gote)
	gote_hi.Promoted = true
	sente_kaku := NewKoma(14, Kaku, 1, 1, Sente)
	sente_kaku.Promoted = true

	if sente_fu.CanFarMove() != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", sente_fu.CanFarMove(), false)
	}
	if gote_kyo.CanFarMove() != true {
		t.Errorf("Id: actual:[%v] expected:[%v]", gote_kyo.CanFarMove(), true)
	}
	if sente_kei.CanFarMove() != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", sente_kei.CanFarMove(), false)
	}
	if gote_gin.CanFarMove() != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", gote_gin.CanFarMove(), false)
	}
	if sente_kin.CanFarMove() != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", sente_kin.CanFarMove(), false)
	}
	if gote_kaku.CanFarMove() != true {
		t.Errorf("Id: actual:[%v] expected:[%v]", gote_kaku.CanFarMove(), true)
	}
	if sente_hi.CanFarMove() != true {
		t.Errorf("Id: actual:[%v] expected:[%v]", sente_hi.CanFarMove(), true)
	}
	if gote_gyoku.CanFarMove() != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", gote_gyoku.CanFarMove(), false)
	}
	if gote_fu.CanFarMove() != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", gote_fu.CanFarMove(), false)
	}
	if sente_kyo.CanFarMove() != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", sente_kyo.CanFarMove(), false)
	}
	if gote_kei.CanFarMove() != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", gote_kei.CanFarMove(), false)
	}
	if sente_gin.CanFarMove() != false {
		t.Errorf("Id: actual:[%v] expected:[%v]", sente_gin.CanFarMove(), false)
	}
	if gote_hi.CanFarMove() != true {
		t.Errorf("Id: actual:[%v] expected:[%v]", gote_hi.CanFarMove(), true)
	}
	if sente_kaku.CanFarMove() != true {
		t.Errorf("Id: actual:[%v] expected:[%v]", sente_kaku.CanFarMove(), true)
	}
	p("TestCanFarMove")
}

func Test_getAiteBan(t *testing.T) {
	sente := Sente_i
	gote := Gote_i
	sente_no_aite_ban := sente.getAiteBan()
	if sente_no_aite_ban != Gote_i {
		t.Errorf("actual:[%v] expected:[%v]", sente_no_aite_ban, Gote_i)
	}
	gote_no_aite_ban := gote.getAiteBan()
	if gote_no_aite_ban != Sente_i {
		t.Errorf("actual:[%v] expected:[%v]", gote_no_aite_ban, Sente_i)
	}
	p("Test_getAiteBan")
}
