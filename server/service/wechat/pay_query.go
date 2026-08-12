package wechat

import (
	"errors"
	"fmt"

	"fresh-shop/server/global"

	"github.com/silenceper/wechat/v2/pay/notify"
	orderPay "github.com/silenceper/wechat/v2/pay/order"
)

// QueryOrderByOutTradeNo 按商户订单号查微信订单。
// mock=true 时不调微信，返回 ResultCode=SUCCESS、TradeState=NOTPAY。
func QueryOrderByOutTradeNo(orderSn string) (notify.PaidResult, error) {
	if orderSn == "" {
		return notify.PaidResult{}, errors.New("订单号为空")
	}
	cfg := global.Config.WechatPay.Compensate
	if cfg.Mock {
		rc := "SUCCESS"
		tsc := "SUCCESS"
		ts := "NOTPAY"
		fee := 0
		return notify.PaidResult{
			ReturnCode: &rc,
			ResultCode: &tsc,
			TradeState: &ts,
			OutTradeNo: &orderSn,
			TotalFee:   &fee,
		}, nil
	}
	if global.WxPay == nil {
		return notify.PaidResult{}, errors.New("微信支付未初始化")
	}
	o := global.WxPay.GetOrder()
	result, err := o.QueryOrder(&orderPay.QueryParams{OutTradeNo: orderSn})
	if err != nil {
		global.SugarLog.Errorf("微信查单失败 orderSn=%s err=%s \n", orderSn, err.Error())
		return result, err
	}
	return result, nil
}

func tradeStateOf(r notify.PaidResult) string {
	if r.TradeState == nil {
		return ""
	}
	return *r.TradeState
}

func ptrStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func ptrInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func formatQueryBrief(r notify.PaidResult) string {
	return fmt.Sprintf("trade_state=%s total_fee=%d transaction_id=%s",
		tradeStateOf(r), ptrInt(r.TotalFee), ptrStr(r.TransactionID))
}
