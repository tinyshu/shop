# fulfillment — uniapp 任务清单（v0.3.0 / M2-3）

基于 [design.md](./design.md)。`api/order.js` URL 已是 `POST /order/confirmOrder`，入参 `{ ID }`，无需改。不改为 `updateOrderDelivery`。

---

## 执行顺序

1. ✅ 任务 1 — `detail.vue` 确认收货按钮补未取消条件

---

## 任务 1：detail.vue `✅`

文件：`new_shop/uniapp/pages/order/detail.vue`

确认收货按钮：`v-if="order.status === 2 && order.statusCancel === 0"`。`toConfirmOrder` 仍传 `{ ID }`。列表注释按钮本版不解开。
