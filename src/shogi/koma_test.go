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
}
