# fulfillment — web 任务清单（v0.1.0 / M2-1）

基于 [design.md](./design.md)。

---

## 执行顺序

1. ✅ `api/order.js` → `markRefundDone`
2. ✅ `order.vue` / `orderMonth.vue` 按钮 + 二次确认

---

## 任务 1：API `✅`

`web/src/api/order.js`：`POST /order/markRefundDone`

## 任务 2：订单页 `✅`

- 显示：`status > 0 && statusRefund !== 2`
- 确认文案：已在微信商户平台退款？
