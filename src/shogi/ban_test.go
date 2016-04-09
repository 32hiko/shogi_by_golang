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
	result := ban1.Analyze()
	p("Sente:kikiMasu: " + s(result["Sente:kikiMasu"]))
	p("Gote:kikiMasu: " + s(result["Gote:kikiMasu"]))
	p("TestAnalyze")
}
