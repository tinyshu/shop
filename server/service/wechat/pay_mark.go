package wechat

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"fresh-shop/server/global"
	"fresh-shop/server/model/shop"

	"gorm.io/gorm"
)

// ErrOrderCancelled 本地已取消，拒绝自动入账（人工处理）
var ErrOrderCancelled = errors.New("订单已取消，拒绝入账")

// MarkPaidInput 微信侧支付成功后的入账入参（回调与查单共用）
type MarkPaidInput struct {
	OrderSn       string
	TotalFeeFen   int64
	PayTime       time.Time
	OpenID        string
	TransactionID string
	Attach        string
}

// MarkOrderPaidFromWechat 金额校验 + 条件更新 status=0→1；已支付幂等成功。
func MarkOrderPaidFromWechat(in MarkPaidInput) error {
	logPrefix := fmt.Sprintf("订单支付入账: 订单号：%s, ", in.OrderSn)
	var order shop.Order
	if errors.Is(global.DB.Where("order_sn = ?", in.OrderSn).First(&order).Error, gorm.ErrRecordNotFound) {
		global.SugarLog.Errorf(logPrefix + "订单不存在 \n")
		return errors.New("订单不存在")
	}
	if order.Status != nil && *order.Status == 1 {
		global.SugarLog.Infof(logPrefix + "订单已支付(幂等) \n")
		return nil
	}
	if order.StatusCancel != nil && *order.StatusCancel != 0 {
		global.SugarLog.Errorf(logPrefix+"订单已取消 status_cancel=%d，拒绝自动入账 \n", *order.StatusCancel)
		return ErrOrderCancelled
	}
	if order.Status == nil || *order.Status != 0 {
		st := -1
		if order.Status != nil {
			st = *order.Status
		}
		global.SugarLog.Errorf(logPrefix+"订单状态不正确 Status：%d \n", st)
		return errors.New("订单状态不正确")
	}

	debug := global.Config.WechatPay.Debug
	if !NotifyAmountMatches(order.Total, in.TotalFeeFen, debug) {
		expectedFen, _ := ExpectedPayFen(order.Total, debug)
		global.SugarLog.Errorf(logPrefix+"金额不符 expectedFen=%d notifyFen=%d orderTotal=%.2f debug=%v \n",
			expectedFen, in.TotalFeeFen, order.Total, debug)
		return errors.New("支付金额与订单不符")
	}

	finishStr := fmt.Sprintf("%.2f", float64(in.TotalFeeFen)/100)
	finish, err := strconv.ParseFloat(finishStr, 64)
	if err != nil {
		global.SugarLog.Errorf(logPrefix+"finish 转换失败 totalFeeFen=%d \n", in.TotalFeeFen)
		return err
	}

	updates := map[string]interface{}{
		"finish":          finish,
		"status":          1,
		"pay_time":        in.PayTime,
		"payment_openid":  in.OpenID,
		"payment_info":    in.Attach,
		"transation_id":   in.TransactionID,
	}
	res := global.DB.Model(&shop.Order{}).
		Where("order_sn = ? AND status = ?", in.OrderSn, 0).
		Updates(updates)
	if res.Error != nil {
		global.SugarLog.Errorf(logPrefix+"条件更新失败, err:%s \n", res.Error.Error())
		return res.Error
	}
	if res.RowsAffected == 0 {
		// 并发下可能已被写成已支付
		var again shop.Order
		if err := global.DB.Where("order_sn = ?", in.OrderSn).First(&again).Error; err != nil {
			return err
		}
		if again.Status != nil && *again.Status == 1 {
			global.SugarLog.Infof(logPrefix + "条件更新 RowsAffected=0，已是已支付(幂等) \n")
			return nil
		}
		global.SugarLog.Errorf(logPrefix + "条件更新未影响行且非已支付 \n")
		return errors.New("订单状态不允许入账")
	}

	global.SugarLog.Infof(logPrefix + "支付成功")
	return nil
}
