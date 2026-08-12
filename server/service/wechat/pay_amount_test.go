package wechat

import "testing"

func TestExpectedPayFen(t *testing.T) {
	fen, err := ExpectedPayFen(12.34, false)
	if err != nil || fen != 1234 {
		t.Fatalf("12.34 => (%d,%v) want 1234", fen, err)
	}
	fen, err = ExpectedPayFen(0.01, false)
	if err != nil || fen != 1 {
		t.Fatalf("0.01 => (%d,%v) want 1", fen, err)
	}
	fen, err = ExpectedPayFen(99.99, true)
	if err != nil || fen != 1 {
		t.Fatalf("debug => (%d,%v) want 1", fen, err)
	}
	if _, err := ExpectedPayFen(-1, false); err == nil {
		t.Fatal("negative should error")
	}
}

func TestNotifyAmountMatches(t *testing.T) {
	if !NotifyAmountMatches(10.00, 1000, false) {
		t.Fatal("10.00 vs 1000 fen should match")
	}
	if NotifyAmountMatches(10.00, 999, false) {
		t.Fatal("10.00 vs 999 fen should not match")
	}
	if !NotifyAmountMatches(88.88, 1, true) {
		t.Fatal("debug should expect 1 fen")
	}
	if NotifyAmountMatches(-1, 0, false) {
		t.Fatal("invalid total should not match")
	}
}
