package shop

import (
	"errors"
	"time"

	"fresh-shop/server/global"
	"fresh-shop/server/model/business"
	"fresh-shop/server/model/shop"
	sysModel "fresh-shop/server/model/system"
	"fresh-shop/server/service/common"

	"gorm.io/gorm"
)

// ConfirmOrder C 端确认收货：仅下单人，已发货且未取消。
func (orderService *OrderService) ConfirmOrder(userId uint, orderId uint) error {
	if orderId == 0 {
		return errors.New("订单ID无效")
	}
	var order shop.Order
	if err := global.DB.Where("id = ?", orderId).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	if order.UserId == nil || uint(*order.UserId) != userId {
		return errors.New("无权操作")
	}
	if order.StatusRefund != nil && *order.StatusRefund == StatusRefundDone {
		return errors.New("订单已退款，不可确认收货")
	}
	now := time.Now()
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return completeReceive(tx, order, now)
	})
}

// completeReceive 条件更新订单为已收货，并同步发货单/积分。
func completeReceive(tx *gorm.DB, order shop.Order, receiveTime time.Time) error {
	res := tx.Model(&shop.Order{}).
		Where("id = ? AND status = ? AND status_cancel = ?", order.ID, OrderStatusShipped, StatusCancelNone).
		Updates(map[string]interface{}{
			"status":       OrderStatusReceived,
			"receive_time": receiveTime,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("订单状态不允许确认收货")
	}

	var delivery shop.OrderDelivery
	dErr := tx.Where("order_id = ?", order.ID).First(&delivery).Error
	if dErr == nil {
		if err := tx.Model(&shop.OrderDelivery{}).Where("id = ?", delivery.ID).
			Update("receipt_time", receiveTime).Error; err != nil {
			return err
		}
		if delivery.DeliveryId != nil && *delivery.DeliveryId > 0 {
			if err := tx.Model(&business.UserDelivery{}).Where("id = ?", *delivery.DeliveryId).
				UpdateColumn("deliver_count", gorm.Expr("IFNULL(deliver_count,0) + 1")).Error; err != nil {
				return err
			}
		}
	} else if !errors.Is(dErr, gorm.ErrRecordNotFound) {
		return dErr
	}

	if order.GoodsArea != nil && *order.GoodsArea == 0 && order.GiftPoints > 0 {
		if order.UserId == nil {
			return errors.New("用户不存在")
		}
		var user sysModel.SysUser
		if err := tx.Where("id = ?", *order.UserId).First(&user).Error; err != nil {
			return errors.New("用户不存在")
		}
		f := common.NewFinance(0, 6, user.ID, user.Username, order.GiftPoints, order.OrderSn, user.ID, user.Username, "确认收货发放积分")
		if err := common.AccountUnifyDeduction(common.POINT, f); err != nil {
			global.SugarLog.Errorf("发放积分失败 UserFinance:%v, error: %v", f, err)
			return err
		}
	}
	return nil
}
