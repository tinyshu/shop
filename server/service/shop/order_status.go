package shop

// 订单履约状态 shop_order.status
const (
	OrderStatusUnpaid    = 0
	OrderStatusPaid      = 1 // 已付款待发货
	OrderStatusShipped   = 2
	OrderStatusReceived  = 3
)

// 取消状态 shop_order.status_cancel
const (
	StatusCancelNone    = 0
	StatusCancelUser    = 1
	StatusCancelAdmin   = 2
	StatusCancelTimeout = 3
)
