# fulfillment — server 任务清单（v0.1.0 / M2-1）

基于 [design.md](./design.md)。全局约束：不调微信、不回库存；条件更新 + 事务。

---

## 执行顺序

1. ✅ 常量与 MarkRefundDone 服务
2. ✅ API + 路由
3. ✅ Casbin/API 注册 SQL（幂等）
4. ✅ 单测骨架 / go build

---

## 任务 1：order_refund.go `✅`

文件：`server/service/shop/order_refund.go`

- `StatusRefund*` / `ReturnStatus*` / `ReturnRefund*` 常量
- 售后完成态：`status=1`（对齐模型注释与字典 desc「1已退款」）
- `MarkRefundDone(orderId) → (result, error)`

## 任务 2：API + 路由 `✅`

- `api/v1/shop/order.go`：`MarkRefundDone`
- `router/shop/order.go`：`POST markRefundDone`（OperationRecord 组）

## 任务 3：权限 SQL `✅`

`sql/migrations/20260812_fulfillment_m2_1_mark_refund.sql`：sys_apis + casbin（对已有 updateOrder 的角色同步赋权）

## 任务 4：验证 `✅`

- `go test` / `go build`
