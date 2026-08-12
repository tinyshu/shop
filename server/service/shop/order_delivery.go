package shop

import (
	"errors"
	"time"

	"fresh-shop/server/global"
	"fresh-shop/server/model/common/request"
	"fresh-shop/server/model/shop"
	shopReq "fresh-shop/server/model/shop/request"

	"gorm.io/gorm"
)

type OrderDeliveryService struct {
}

// CreateOrderDelivery 创建OrderDelivery记录（条件更新：仅待发货且未取消）
// Author [dalefeng](https://github.com/dalefeng)
func (orderDeliveryService *OrderDeliveryService) CreateOrderDelivery(orderDelivery shop.OrderDelivery) (err error) {
	if orderDelivery.OrderId == nil || *orderDelivery.OrderId == 0 {
		return errors.New("订单id参数错误")
	}
	orderId := *orderDelivery.OrderId
	now := time.Now()
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&shop.Order{}).
			Where("id = ? AND status = ? AND status_cancel = ?", orderId, OrderStatusPaid, StatusCancelNone).
			Updates(map[string]interface{}{
				"status":        OrderStatusShipped,
				"shipment_time": now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			var again shop.Order
			if qErr := tx.Where("id = ?", orderId).First(&again).Error; errors.Is(qErr, gorm.ErrRecordNotFound) {
				return errors.New("订单不存在")
			}
			global.SugarLog.Errorf("发货条件更新未影响行 orderId:%d", orderId)
			return errors.New("订单状态不允许发货")
		}
		return tx.Create(&orderDelivery).Error
	})
	return
}

// DeleteOrderDelivery 删除OrderDelivery记录
// Author [dalefeng](https://github.com/dalefeng)
func (orderDeliveryService *OrderDeliveryService) DeleteOrderDelivery(orderDelivery shop.OrderDelivery) (err error) {
	err = global.DB.Delete(&orderDelivery).Error
	return err
}

// DeleteOrderDeliveryByIds 批量删除OrderDelivery记录
// Author [dalefeng](https://github.com/dalefeng)
func (orderDeliveryService *OrderDeliveryService) DeleteOrderDeliveryByIds(ids request.IdsReq) (err error) {
	err = global.DB.Delete(&[]shop.OrderDelivery{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateOrderDelivery 订单收货（管理端）。有 receiptTime 时条件更新订单 status 2→3。
func (orderDeliveryService *OrderDeliveryService) UpdateOrderDelivery(orderDelivery shop.OrderDelivery) (err error) {
	if orderDelivery.OrderId == nil || *orderDelivery.OrderId == 0 {
		return errors.New("订单id参数错误")
	}
	if orderDelivery.ReceiptTime == nil {
		return global.DB.Save(&orderDelivery).Error
	}
	var order shop.Order
	if err = global.DB.Where("id = ?", *orderDelivery.OrderId).First(&order).Error; err != nil {
		global.SugarLog.Errorf("获取订单信息失败 orderId:%d, error: %v", *orderDelivery.OrderId, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	if order.StatusRefund != nil && *order.StatusRefund == StatusRefundDone {
		return errors.New("订单已退款，不可确认收货")
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if txErr := completeReceive(tx, order, *orderDelivery.ReceiptTime); txErr != nil {
			return txErr
		}
		if orderDelivery.ID != 0 {
			return tx.Save(&orderDelivery).Error
		}
		return nil
	})
}

// GetOrderDelivery 根据id获取OrderDelivery记录
// Author [dalefeng](https://github.com/dalefeng)
func (orderDeliveryService *OrderDeliveryService) GetOrderDelivery(id uint, orderId int) (orderDelivery shop.OrderDelivery, err error) {
	db := global.DB
	if id != 0 {
		db = db.Where("id = ?", id)
	} else if orderId != 0 {
		db = db.Where("order_id = ?", orderId)
	} else {
		return orderDelivery, errors.New("参数异常")
	}
	err = db.First(&orderDelivery).Error
	return
}

// GetOrderDeliveryInfoList 分页获取OrderDelivery记录
// Author [dalefeng](https://github.com/dalefeng)
func (orderDeliveryService *OrderDeliveryService) GetOrderDeliveryInfoList(info shopReq.OrderDeliverySearch) (list []shop.OrderDelivery, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&shop.OrderDelivery{})
	var orderDelivers []shop.OrderDelivery
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.StartScheduledTime != nil && info.EndScheduledTime != nil {
		db = db.Where("scheduled_time BETWEEN ? AND ? ", info.StartScheduledTime, info.EndScheduledTime)
	}
	if info.DeliverName != "" {
		db = db.Where("deliver_name LIKE ?", "%"+info.DeliverName+"%")
	}
	if info.DeliverMobile != "" {
		db = db.Where("deliver_mobile LIKE ?", "%"+info.DeliverMobile+"%")
	}
	if info.StartReceiptTime != nil && info.EndReceiptTime != nil {
		db = db.Where("receipt_time BETWEEN ? AND ? ", info.StartReceiptTime, info.EndReceiptTime)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&orderDelivers).Error
	return orderDelivers, total, err
}
