package log

import "testing"

func TestAll(t *testing.T) {
	for i := 0; i < 1; i++ {
		Info("hello info")
		Error("hello error")
		Warning("hello waring")
	}
}
