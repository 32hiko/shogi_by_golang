package shogi

import (
	// . "logger"
	"testing"
)

func TestFromSFEN(t *testing.T) {
	// logger := GetLogger()
	// logger.Trace("")
	sfen1 := "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	p(sfen1 + " start")
	ban1 := FromSFEN(sfen1)
	p(ban1.Display())

	sfen2 := "4k4/1r5b1/9/9/9/9/PPPPPPPPP/9/LNSGKGSNL w BR9p2l2n2s2g 123"
	p(sfen2 + " start")
	ban2 := FromSFEN(sfen2)
	p(ban2.Display())
	p("TestFromSFEN")
}

func TestCountKikiMasu(t *testing.T) {
	ban1 := CreateInitialState()
	count := ban1.CountKikiMasu(Sente)
	if count != 30 {
		t.Errorf("count: actual:[%v] expected:[%v]", count, 30)
	}
	p("TestCountKiki")
}

func TestAnalyze(t *testing.T) {
	ban1 := CreateInitialState()
	result_sente, result_gote := ban1.Analyze()
	for k, v := range result_sente {
		p(k + ": " + s(v))
	}
	for k, v := range result_gote {
		p(k + ": " + s(v))
	}
	p("TestAnalyze")
}

func TestPlaceKoma(t *testing.T) {
	ban1 := NewBan()

	koma1 := NewKoma(1, Fu, 2, 7, Sente)
	ban1.PlaceKoma(koma1)

	koma2 := NewKoma(2, Kyo, 1, 1, Gote)
	ban1.PlaceKoma(koma2)

	koma3 := NewKoma(3, Kei, 2, 9, Sente)
	koma3.Promoted = true
	ban1.PlaceKoma(koma3)

	for a := 0; a < 2; a++ {
		for b := 0; b < 14; b++ {
			for c := 0; c < 18; c++ {
				if ban1.Koma[a][b][c] != 0 {
					p("ban1.Koma[" + s(a) + "][" + s(b) + "][" + s(c) + "]=[" + s(ban1.Koma[a][b][c]) + "]")
				}
			}
		}
	}
	if ban1.Koma[0][0][0] != 27 {
		t.Errorf("actual:[%v] expected:[%v]", ban1.Koma[0][0][0], 27)
	}
	if ban1.Koma[1][1][0] != 11 {
		t.Errorf("actual:[%v] expected:[%v]", ban1.Koma[1][1][0], 11)
	}
	if ban1.Koma[0][10][0] != 29 {
		t.Errorf("actual:[%v] expected:[%v]", ban1.Koma[0][10][0], 29)
	}
	p("TestPlaceKoma")
}
