package main

import (
	"testing"
)

func Test_handle_go_text_sample_time_up(t *testing.T) {
	result_time_up := handle_go_text_sample("go btime 1000 wtime 1500 binc 10000 winc 10000")
	// assert
	if result_time_up != "bestmove resign" {
		t.Errorf("actual:[%v] expected:[%v]", result_time_up, "bestmove resign")
	}
}

func Test_handle_go_text_sample_in_time(t *testing.T) {
	result_in_time := handle_go_text_sample("go btime 2000 wtime 1500 binc 10000 winc 10000")
	// assert
	if result_in_time != "bestmove 7c7d" {
		t.Errorf("actual:[%v] expected:[%v]", result_in_time, "bestmove 7c7d")
	}
}
