package shogi

import (
	"testing"
)

func TestFixOpening(t *testing.T) {
	joseki := NewJoseki()
	m1 := joseki.FixOpening[1]
	m2 := joseki.FixOpening[2]

	p("1: " + m1.GetUSIMoveString())
	p("2: " + m2.GetUSIMoveString())
	p("TestFixOpening")
}
