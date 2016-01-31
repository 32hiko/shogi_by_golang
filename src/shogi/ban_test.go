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
