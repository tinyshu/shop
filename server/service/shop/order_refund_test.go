package shop

import "testing"

func TestRefundStatusConstants(t *testing.T) {
	if StatusRefundDone != 2 {
		t.Fatalf("StatusRefundDone=%d want 2", StatusRefundDone)
	}
	if ReturnStatusDone != 1 {
		t.Fatalf("ReturnStatusDone=%d want 1 (model/dict)", ReturnStatusDone)
	}
	if ReturnRefundSuccess != 2 {
		t.Fatalf("ReturnRefundSuccess=%d want 2", ReturnRefundSuccess)
	}
}
