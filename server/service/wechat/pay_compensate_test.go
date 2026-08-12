package wechat

import (
	"testing"

	"fresh-shop/server/global"

	"github.com/silenceper/wechat/v2/pay/notify"
)

func TestMarkInputFromPaidResult(t *testing.T) {
	sn := "SN001"
	fee := 1234
	tid := "wx_tx_1"
	oid := "openid_1"
	att := "attach"
	te := "20260812120000"
	in, err := markInputFromPaidResult(notify.PaidResult{
		OutTradeNo:    &sn,
		TotalFee:      &fee,
		TransactionID: &tid,
		OpenID:        &oid,
		Attach:        &att,
		TimeEnd:       &te,
	})
	if err != nil {
		t.Fatal(err)
	}
	if in.OrderSn != sn || in.TotalFeeFen != 1234 || in.TransactionID != tid {
		t.Fatalf("unexpected input: %+v", in)
	}
	if in.PayTime.IsZero() {
		t.Fatal("PayTime should parse from TimeEnd")
	}
}

func TestMarkInputFromPaidResultMissingFee(t *testing.T) {
	sn := "SN002"
	_, err := markInputFromPaidResult(notify.PaidResult{OutTradeNo: &sn})
	if err == nil {
		t.Fatal("expected error for missing total_fee")
	}
}

func TestQueryCountCap(t *testing.T) {
	queryCountMu.Lock()
	queryCountBySn = map[string]int{}
	queryCountMu.Unlock()

	old := global.Config.WechatPay.Compensate.MaxQueryPerOrder
	global.Config.WechatPay.Compensate.MaxQueryPerOrder = 2
	defer func() { global.Config.WechatPay.Compensate.MaxQueryPerOrder = old }()

	sn := "cap-order"
	if !allowQuery(sn) {
		t.Fatal("first should allow")
	}
	incrQuery(sn)
	if !allowQuery(sn) {
		t.Fatal("second should allow")
	}
	incrQuery(sn)
	if allowQuery(sn) {
		t.Fatal("third should be capped")
	}
}
