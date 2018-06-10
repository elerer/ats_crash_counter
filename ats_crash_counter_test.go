package main

import (
	"strings"
	"testing"
)

func TestDateToManagerLogFormt(t *testing.T) {
	date := dateToManagerLogFormt("Jun", "6")
	if strings.Compare(date, "Jun  6") != 0 {
		t.Error(
			"For", date,
			"expected", "Jun  6",
			"got", date,
		)
	}

	date = dateToManagerLogFormt("Jun", "16")
	if strings.Compare(date, "Jun 16") != 0 {
		t.Error(
			"For", date,
			"expected", "Jun 16",
			"got", date,
		)
	}
}
