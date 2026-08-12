# fulfillment — server 任务清单（v0.3.0 / M2-3）

基于 [design.md](./design.md)。全局约束：只改 `new_shop/`；C 端入参 `{ ID }`（订单主键）；收货时间 C 端用服务端 now；条件更新 `status=2 AND status_cancel=0`；禁止 `Save` 整单改 status；Casbin 对齐 `cancelOrder`（勿抄 `updateOrder`）；不调微信、不回库存。

---

## 执行顺序

1. ✅ 任务 1 — `order_receive.go` ConfirmOrder + completeReceive
2. ✅ 任务 2 — API + 路由
3. ✅ 任务 3 — `UpdateOrderDelivery` 条件更新
4. ✅ 任务 4 — Casbin / sys_apis 幂等 SQL
5. ✅ 任务 5 — 单测 + `go build` + 文档进度

---

## 任务 1：order_receive.go — C 端确认收货服务 `✅`

文件：`new_shop/server/service/shop/order_receive.go`（新建）

- `ConfirmOrder(userId, orderId)`：归属校验、已退款拒绝
- `completeReceive`：`WHERE status=2 AND status_cancel=0` → `status=3`；有发货单写 `receipt_time` / `deliver_count+1`；`gift_points>0` 才发积分

---

## 任务 2：API + 路由 `✅`

- `api/v1/shop/order.go`：`ConfirmOrder`，body `{ ID }`
- `router/shop/order.go`：`POST confirmOrder`（OperationRecord 组）

---

## 任务 3：order_delivery.go — 管理端收货 `✅`

有 `receiptTime` 时走 `completeReceive`，禁止 `Save` 整单改 status；无 `receiptTime` 只改发货单。

---

## 任务 4：迁移 SQL `✅`

`sql/migrations/20260812_fulfillment_m2_3_confirm_order.sql`：sys_apis + casbin（从 cancelOrder 角色复制）

---

## 任务 5：验证 `✅`

- `order_receive_test.go` 常量
- `go test ./service/shop/`、`go build`
- 文档进度 + [link_test](../../../link_test/fulfillment/v0.3.0/link_test.md)
