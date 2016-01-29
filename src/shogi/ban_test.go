package shogi

import (
	"testing"
)

func TestFromSFEN(t *testing.T) {
	sfen1 := "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	p(sfen1 + " start")
	FromSFEN(sfen1)

	sfen2 := "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 123"
	p(sfen2 + " start")
	FromSFEN(sfen2)
	p("TestFromSFEN")
}
