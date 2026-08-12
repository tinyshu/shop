package shop

import "testing"

func TestResolveCancelType(t *testing.T) {
	if got := resolveCancelType(nil); got != StatusCancelUser {
		t.Fatalf("nil => %d want %d (用户取消)", got, StatusCancelUser)
	}
	user := StatusCancelUser
	if got := resolveCancelType(&user); got != StatusCancelUser {
		t.Fatalf("1 => %d want 1", got)
	}
	admin := StatusCancelAdmin
	if got := resolveCancelType(&admin); got != StatusCancelAdmin {
		t.Fatalf("2 => %d want 2 (后台取消)", got)
	}
	timeout := StatusCancelTimeout
	if got := resolveCancelType(&timeout); got != StatusCancelTimeout {
		t.Fatalf("3 => %d want 3 (超时取消)", got)
	}
}

func TestShipCancelPredicates(t *testing.T) {
	// 发货 WHERE：status=1 AND status_cancel=0；取消 WHERE：status<2 AND status_cancel=0
	if OrderStatusPaid != 1 || OrderStatusShipped != 2 || StatusCancelNone != 0 {
		t.Fatalf("status constants drifted: paid=%d shipped=%d cancelNone=%d",
			OrderStatusPaid, OrderStatusShipped, StatusCancelNone)
	}
}
