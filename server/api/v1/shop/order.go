package shop

import (
	"fresh-shop/server/global"
	"fresh-shop/server/model/common/request"
	"fresh-shop/server/model/common/response"
	"fresh-shop/server/model/shop"
	shopReq "fresh-shop/server/model/shop/request"
	"fresh-shop/server/service"
	"fresh-shop/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type OrderApi struct {
}

var orderService = service.ServiceGroupApp.ShopServiceGroup.OrderService
var wechatService = service.ServiceGroupApp.WechatServiceGroup.WechatService

// CreateOrder 创建待支付订单 Order
// @Tags Order
// @Summary 创建Order
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "创建Order"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /order/createOrder [post]
func (orderApi *OrderApi) CreateOrder(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if order.AddressId == 0 && *order.ShipmentType == 0 {
		response.FailWithMessage("请选择收货地址", c)
		return
	}
	userId := utils.GetUserID(c)
	order.UserId = utils.Pointer(int(userId))
	if orderResp, err := orderService.CreateOrder(order, c.ClientIP()); err != nil {
		global.Log.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(orderResp, c)
	}
}

// OrderPay 支付 Order, 返回微信支付所需要的参数
// @Tags Order
// @Summary 支付 Order, 返回微信支付所需要的参数
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "支付 Order, 返回微信支付所需要的参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /order/orderPay [post]
func (orderApi *OrderApi) OrderPay(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if order.ID == 0 {
		response.FailWithMessage("订单ID不能为空", c)
		return
	}
	userId := utils.GetUserID(c)
	order.UserId = utils.Pointer(int(userId))
	userClaims := utils.GetUserInfo(c)
	if orderResp, err := orderService.OrderPay(order, userClaims, c.ClientIP()); err != nil {
		global.Log.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(orderResp, c)
	}
}

// DeleteOrder 删除Order
// @Tags Order
// @Summary 删除Order
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "删除Order"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /order/deleteOrder [delete]
func (orderApi *OrderApi) DeleteOrder(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := orderService.DeleteOrder(order); err != nil {
		global.Log.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// CancelOrder 取消订单
// @Tags Order
// @Summary 取消订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "取消订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"取消成功"}"
// @Router /order/cancelOrder [delete]
func (orderApi *OrderApi) CancelOrder(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := orderService.CancelOrder(order); err != nil {
		global.Log.Error("取消失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("取消成功", c)
	}
}

// ConfirmOrder 确认收货（C 端，入参订单主键 ID）
// @Tags Order
// @Summary 确认收货
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "确认收货"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"确认收货成功"}"
// @Router /order/confirmOrder [post]
func (orderApi *OrderApi) ConfirmOrder(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if order.ID == 0 {
		response.FailWithMessage("订单ID不能为空", c)
		return
	}
	userId := utils.GetUserID(c)
	if err := orderService.ConfirmOrder(userId, order.ID); err != nil {
		global.Log.Error("确认收货失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("确认收货成功", c)
	}
}

// DeleteOrderByIds 批量删除Order
// @Tags Order
// @Summary 批量删除Order
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Order"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /order/deleteOrderByIds [delete]
func (orderApi *OrderApi) DeleteOrderByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := orderService.DeleteOrderByIds(IDS); err != nil {
		global.Log.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateOrder 更新Order
// @Tags Order
// @Summary 更新Order
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "更新Order"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /order/updateOrder [put]
func (orderApi *OrderApi) UpdateOrder(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := orderService.UpdateOrder(order); err != nil {
		global.Log.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// BatchSettlement 批量结算用户
// @Tags Order
// @Summary 批量结算用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "批量结算用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /order/batchSettlement [post]
func (orderApi *OrderApi) BatchSettlement(c *gin.Context) {
	var settlement shopReq.BatchSettlement
	err := c.ShouldBindJSON(&settlement)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := orderService.BatchSettlement(settlement); err != nil {
		global.Log.Error("更新失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// SyncWechatPay 管理端按订单查微信并补单（PAY-01）
// @Tags Order
// @Summary 同步微信支付状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body object true "orderId 或 orderSn"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /order/syncWechatPay [post]
func (orderApi *OrderApi) SyncWechatPay(c *gin.Context) {
	var req struct {
		OrderId uint   `json:"orderId"`
		OrderSn string `json:"orderSn"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := wechatService.SyncWechatPayByOrder(req.OrderId, req.OrderSn)
	if err != nil {
		global.Log.Error("同步微信支付失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// MarkRefundDone 管理端标记退款完成（商户平台人工退款后的状态闭环）
// @Tags Order
// @Summary 标记退款完成
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body object true "orderId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /order/markRefundDone [post]
func (orderApi *OrderApi) MarkRefundDone(c *gin.Context) {
	var req struct {
		OrderId uint `json:"orderId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := orderService.MarkRefundDone(req.OrderId)
	if err != nil {
		global.Log.Error("标记退款完成失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// FindOrder 用id查询Order
// @Tags Order
// @Summary 用id查询Order
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Order true "用id查询Order"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /order/findOrder [get]
func (orderApi *OrderApi) FindOrder(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindQuery(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if reorder, err := orderService.GetOrder(order.ID); err != nil {
		global.Log.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reorder": reorder}, c)
	}
}

// FindUserOrderStatus 获取用户订单中数量
// @Tags Order
// @Summary 获取用户订单中数量
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Order true "获取用户订单中数量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /order/findUserOrderStatus [get]
func (orderApi *OrderApi) FindUserOrderStatus(c *gin.Context) {
	settlementMonthStr := c.Query("settlementMonth")
	var settlementMonth time.Time
	var err error
	if settlementMonthStr != "" {
		// 格式化为时间
		settlementMonth, err = time.Parse("2006-01-02 15:04:05", settlementMonthStr)
		if err != nil {
			response.FailWithMessage("时间格式错误", c)
		}
	}
	userId := utils.GetUserID(c)
	if reorder, err := orderService.FindUserOrderStatus(userId, settlementMonth); err != nil {
		global.Log.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(reorder, c)
	}
}

// OrderStatus 获取订单状态
// @Tags Order
// @Summary 用id查询Order
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Order true "获取订单状态"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /order/orderStatus [get]
func (orderApi *OrderApi) OrderStatus(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindQuery(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if order.ID == 0 {
		response.FailWithMessage("订单ID不能为空", c)
		return
	}
	if reorder, err := orderService.OrderStatus(order.ID); err != nil {
		global.Log.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(reorder, c)
	}
}

// GetOrderList 分页获取Order列表
// @Tags Order
// @Summary 分页获取Order列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.OrderSearch true "分页获取Order列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /order/getOrderList [get]
func (orderApi *OrderApi) GetOrderList(c *gin.Context) {
	var pageInfo shopReq.OrderSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := orderService.GetOrderInfoList(pageInfo); err != nil {
		global.Log.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetOrderMonthStatistics 获取月结统计
// @Tags Order
// @Summary 获取月结统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.OrderSearch true "获取月结统计"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /order/GetOrderMonthStatistics [get]
func (orderApi *OrderApi) GetOrderMonthStatistics(c *gin.Context) {
	var pageInfo shopReq.OrderSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if resp, err := orderService.GetOrderMonthStatistics(pageInfo); err != nil {
		global.Log.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(resp, c)
	}
}

// GetUserOrderList 分页获取登录用户获取Order列表
// @Tags Order
// @Summary 页获取登录用户获取Order列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.OrderSearch true "页获取登录用户获取Order列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /order/getUserOrderList [get]
func (orderApi *OrderApi) GetUserOrderList(c *gin.Context) {
	var pageInfo shopReq.OrderSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	pageInfo.UserId = utils.Pointer(int(userId))
	if list, total, err := orderService.GetOrderInfoList(pageInfo); err != nil {
		global.Log.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
