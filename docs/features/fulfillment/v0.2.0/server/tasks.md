# fulfillment — server 任务清单（v0.2.0 / M2-2）

基于 [design.md](./design.md)。全局约束：只改 `new_shop/`；无 schema；不改确认收货 / 回库存 / 退款；写库必须以 `WHERE` 条件为准，禁止 `Save` 整单。

---

## 执行顺序

1. ✅ 任务 1 — `CreateOrderDelivery` 条件更新
2. ✅ 任务 2 — `CancelOrder` 条件更新
3. ✅ 任务 3 — 发货失败文案透出
4. ✅ 任务 4 — 单测 + `go build` + 文档进度

---

## 任务 1：order_delivery.go `✅`

文件：`new_shop/server/service/shop/order_delivery.go`

- 事务内 `UPDATE ... WHERE status=1 AND status_cancel=0`
- `RowsAffected==0` → 订单不存在 / 状态不允许发货
- 成功再 `Create` delivery；**未改**确认收货

## 任务 2：order.go `✅`

文件：`new_shop/server/service/shop/order.go`

- `resolveCancelType`；`UPDATE ... WHERE status<2 AND status_cancel=0`
- `RowsAffected==0` → `订单不允许取消`

## 任务 3：order_delivery API `✅`

发货失败回传 `err.Error()`

## 任务 4：验证 `✅`

- `go test ./service/shop/`
- `go build`
