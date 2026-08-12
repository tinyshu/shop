package wechat

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"fresh-shop/server/global"
	"fresh-shop/server/model/shop"

	"github.com/silenceper/wechat/v2/pay/notify"
)

const (
	ActionMarkedPaid      = "marked_paid"
	ActionAlreadyPaid     = "already_paid"
	ActionStillUnpaid     = "still_unpaid"
	ActionAmountMismatch  = "amount_mismatch"
	ActionCancelledSkip   = "cancelled_skip"
	ActionQuerySkippedCap = "query_skipped_cap"
	ActionQueryFailed     = "query_failed"
)

// SyncPayResult 管理端 / 补偿结果
type SyncPayResult struct {
	LocalStatus      int    `json:"localStatus"`
	WechatTradeState string `json:"wechatTradeState"`
	Action           string `json:"action"`
	OrderSn          string `json:"orderSn"`
}

var queryCountMu sync.Mutex
var queryCountBySn = map[string]int{}

func compensateCfg() (enableAdmin bool, maxPerOrder int) {
	c := global.Config.WechatPay.Compensate
	return c.AdminSyncEnable, c.MaxQueryPerOrder
}

func allowQuery(orderSn string) bool {
	_, maxPer := compensateCfg()
	if maxPer <= 0 {
		return true
	}
	queryCountMu.Lock()
	defer queryCountMu.Unlock()
	return queryCountBySn[orderSn] < maxPer
}

func incrQuery(orderSn string) {
	queryCountMu.Lock()
	defer queryCountMu.Unlock()
	queryCountBySn[orderSn]++
}

// SyncWechatPayByOrder 单笔：查微信并尝试入账（管理端）
func (s *WechatService) SyncWechatPayByOrder(orderID uint, orderSn string) (SyncPayResult, error) {
	if !global.Config.WechatPay.Compensate.AdminSyncEnable {
		return SyncPayResult{}, errors.New("管理端支付同步已关闭(wechatPay.compensate.adminSyncEnable)")
	}
	var order shop.Order
	db := global.DB.Model(&shop.Order{})
	if orderID > 0 {
		if err := db.Where("id = ?", orderID).First(&order).Error; err != nil {
			return SyncPayResult{}, errors.New("订单不存在")
		}
	} else if orderSn != "" {
		if err := db.Where("order_sn = ?", orderSn).First(&order).Error; err != nil {
			return SyncPayResult{}, errors.New("订单不存在")
		}
	} else {
		return SyncPayResult{}, errors.New("请传入 orderId 或 orderSn")
	}
	return compensateOne(order)
}

func compensateOne(order shop.Order) (SyncPayResult, error) {
	res := SyncPayResult{OrderSn: order.OrderSn}
	if order.Status != nil {
		res.LocalStatus = *order.Status
	}
	if order.Status != nil && *order.Status == 1 {
		res.Action = ActionAlreadyPaid
		res.WechatTradeState = "SUCCESS"
		return res, nil
	}
	if order.StatusCancel != nil && *order.StatusCancel != 0 {
		res.Action = ActionCancelledSkip
		global.SugarLog.Errorf("掉单补偿跳过：订单已取消 orderSn=%s status_cancel=%d \n", order.OrderSn, *order.StatusCancel)
		return res, nil
	}
	if order.Status == nil || *order.Status != 0 {
		return res, errors.New("订单非待支付状态")
	}

	if !allowQuery(order.OrderSn) {
		res.Action = ActionQuerySkippedCap
		global.SugarLog.Warnf("掉单补偿跳过：超过 maxQueryPerOrder orderSn=%s \n", order.OrderSn)
		return res, nil
	}

	incrQuery(order.OrderSn)
	paid, err := QueryOrderByOutTradeNo(order.OrderSn)
	if err != nil {
		res.Action = ActionQueryFailed
		return res, fmt.Errorf("微信查单失败: %w", err)
	}
	res.WechatTradeState = tradeStateOf(paid)
	global.SugarLog.Infof("掉单查单 orderSn=%s %s \n", order.OrderSn, formatQueryBrief(paid))

	if res.WechatTradeState != "SUCCESS" {
		res.Action = ActionStillUnpaid
		return res, nil
	}

	in, err := markInputFromPaidResult(paid)
	if err != nil {
		res.Action = ActionQueryFailed
		return res, err
	}
	err = MarkOrderPaidFromWechat(in)
	if err != nil {
		if errors.Is(err, ErrOrderCancelled) {
			res.Action = ActionCancelledSkip
			return res, nil
		}
		if err.Error() == "支付金额与订单不符" {
			res.Action = ActionAmountMismatch
			return res, err
		}
		return res, err
	}
	res.LocalStatus = 1
	res.Action = ActionMarkedPaid
	return res, nil
}

func markInputFromPaidResult(paid notify.PaidResult) (MarkPaidInput, error) {
	if paid.OutTradeNo == nil || *paid.OutTradeNo == "" {
		return MarkPaidInput{}, errors.New("查单结果缺少 out_trade_no")
	}
	if paid.TotalFee == nil {
		return MarkPaidInput{}, errors.New("查单结果缺少 total_fee")
	}
	in := MarkPaidInput{
		OrderSn:       *paid.OutTradeNo,
		TotalFeeFen:   int64(*paid.TotalFee),
		OpenID:        ptrStr(paid.OpenID),
		TransactionID: ptrStr(paid.TransactionID),
		Attach:        ptrStr(paid.Attach),
		PayTime:       time.Now(),
	}
	if paid.TimeEnd != nil && *paid.TimeEnd != "" {
		if t, err := time.Parse("20060102150405", *paid.TimeEnd); err == nil {
			in.PayTime = t
		}
	}
	return in, nil
}

// RunCompensateScan 定时扫描待付单并查单补单
func RunCompensateScan() {
	cfg := global.Config.WechatPay.Compensate
	if !cfg.Enable {
		return
	}
	minAge := cfg.MinAgeMinutes
	if minAge <= 0 {
		minAge = 5
	}
	batch := cfg.BatchSize
	if batch <= 0 {
		batch = 20
	}
	cutoff := time.Now().Add(-time.Duration(minAge) * time.Minute)
	var list []shop.Order
	err := global.DB.Where("status = ? AND (status_cancel = 0 OR status_cancel IS NULL)", 0).
		Where("created_at < ?", cutoff).
		Order("created_at asc").
		Limit(batch).
		Find(&list).Error
	if err != nil {
		global.SugarLog.Errorf("掉单扫描查询失败: %s \n", err.Error())
		return
	}
	global.SugarLog.Infof("掉单扫描开始 count=%d minAge=%dm batch=%d \n", len(list), minAge, batch)
	for _, o := range list {
		r, e := compensateOne(o)
		if e != nil {
			global.SugarLog.Errorf("掉单补偿失败 orderSn=%s action=%s err=%s \n", o.OrderSn, r.Action, e.Error())
			continue
		}
		global.SugarLog.Infof("掉单补偿完成 orderSn=%s action=%s trade=%s \n", o.OrderSn, r.Action, r.WechatTradeState)
	}
}
