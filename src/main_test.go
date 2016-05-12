package main

import (
	"testing"
)

func Test_handle_text_time_up(t *testing.T) {
	result_time_up := handle_text("go btime 1000 wtime 1500 binc 10000 winc 10000")
	// assert
	if result_time_up != "bestmove time up" {
		t.Errorf("actual:[%v] expected:[%v]", result_time_up, "bestmove time up")
	}
}

func Test_handle_text_in_time(t *testing.T) {
	result_in_time := handle_text("go btime 2000 wtime 1500 binc 10000 winc 10000")
	// assert
	if result_in_time != "bestmove 7c7d" {
		t.Errorf("actual:[%v] expected:[%v]", result_in_time, "bestmove 7c7d")
	}
}
