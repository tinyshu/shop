package shop

import (
	"errors"
	"time"

	"fresh-shop/server/global"
	"fresh-shop/server/model/shop"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 订单退款状态 shop_order.status_refund
const (
	StatusRefundNone    = 0
	StatusRefunding     = 1
	StatusRefundDone    = 2
	StatusRefundFailed  = 3
)

// 售后单状态 shop_order_return.status（对齐模型注释 / 字典 desc：1=已退款）
const (
	ReturnStatusRejected = -1
	ReturnStatusPending  = 0
	ReturnStatusDone     = 1
)

// 售后退款状态 shop_order_return.refund_status
const (
	ReturnRefundNone    = 0
	ReturnRefunding     = 1
	ReturnRefundSuccess = 2
)

const (
	MarkRefundActionDone    = "marked_done"
	MarkRefundActionAlready = "already_done"
)

// MarkRefundDoneResult 标记退款完成结果
type MarkRefundDoneResult struct {
	OrderId       uint   `json:"orderId"`
	StatusRefund  int    `json:"statusRefund"`
	Action        string `json:"action"`
	ReturnSynced  bool   `json:"returnSynced"`
}

// MarkRefundDone 人工在微信商户平台退款后，将系统订单标记为已退款（不调微信、不回库存）。
func (orderService *OrderService) MarkRefundDone(orderId uint) (MarkRefundDoneResult, error) {
	res := MarkRefundDoneResult{OrderId: orderId, StatusRefund: StatusRefundDone}
	if orderId == 0 {
		return res, errors.New("订单ID无效")
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		var order shop.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", orderId).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("订单不存在")
			}
			return err
		}
		if order.Status == nil || *order.Status == 0 {
			return errors.New("订单未支付，不可标记退款完成")
		}
		curRefund := StatusRefundNone
		if order.StatusRefund != nil {
			curRefund = *order.StatusRefund
		}
		if curRefund == StatusRefundDone {
			res.Action = MarkRefundActionAlready
			res.ReturnSynced = syncOrderReturnIfNeeded(tx, orderId, false)
			return nil
		}
		if curRefund != StatusRefundNone && curRefund != StatusRefunding && curRefund != StatusRefundFailed {
			return errors.New("订单退款状态不允许标记完成")
		}

		upd := tx.Model(&shop.Order{}).
			Where("id = ? AND status_refund IN ?", orderId, []int{StatusRefundNone, StatusRefunding, StatusRefundFailed}).
			Update("status_refund", StatusRefundDone)
		if upd.Error != nil {
			return upd.Error
		}
		if upd.RowsAffected == 0 {
			var again shop.Order
			if err := tx.Where("id = ?", orderId).First(&again).Error; err != nil {
				return err
			}
			if again.StatusRefund != nil && *again.StatusRefund == StatusRefundDone {
				res.Action = MarkRefundActionAlready
				res.ReturnSynced = syncOrderReturnIfNeeded(tx, orderId, false)
				return nil
			}
			return errors.New("订单退款状态不允许标记完成")
		}

		res.Action = MarkRefundActionDone
		res.ReturnSynced = syncOrderReturnIfNeeded(tx, orderId, true)
		return nil
	})
	return res, err
}

// syncOrderReturnIfNeeded 同步售后单；forceUpdate=true 时写入完成态。
// 返回是否存在售后行（用于响应 returnSynced）。
func syncOrderReturnIfNeeded(tx *gorm.DB, orderId uint, forceUpdate bool) bool {
	var ret shop.OrderReturn
	err := tx.Where("order_id = ?", orderId).First(&ret).Error
	if err != nil {
		return false
	}
	if !forceUpdate {
		return true
	}
	now := time.Now()
	updates := map[string]interface{}{
		"refund_status": ReturnRefundSuccess,
		"process_time":  now,
	}
	// 非拒绝单：置售后完成态
	if ret.Status == nil || *ret.Status != ReturnStatusRejected {
		updates["status"] = ReturnStatusDone
	}
	_ = tx.Model(&shop.OrderReturn{}).Where("id = ?", ret.ID).Updates(updates).Error
	return true
}
