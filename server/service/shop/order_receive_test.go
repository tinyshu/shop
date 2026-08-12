package shop

import "testing"

func TestReceivePredicates(t *testing.T) {
	// 确认收货 WHERE：status=2 AND status_cancel=0 → status=3
	if OrderStatusShipped != 2 || OrderStatusReceived != 3 || StatusCancelNone != 0 {
		t.Fatalf("receive constants drifted: shipped=%d received=%d cancelNone=%d",
			OrderStatusShipped, OrderStatusReceived, StatusCancelNone)
	}
	if StatusRefundDone != 2 {
		t.Fatalf("StatusRefundDone=%d want 2", StatusRefundDone)
	}
}
