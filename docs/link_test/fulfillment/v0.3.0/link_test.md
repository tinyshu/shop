---
module: fulfillment
version: v0.3.0
date: 2026-08-12
base_url: http://127.0.0.1:48888
tags: [link_test, api, confirmOrder, M2-3]
---

# fulfillment — 确认收货接口链路测试文档

## 0. 前置条件

- `new_shop/server` 已启动（默认 `http://127.0.0.1:48888`，`router-prefix` 为空）
- 已执行 `sql/migrations/20260812_fulfillment_m2_3_confirm_order.sql`（sys_apis + casbin）
- 环境变量：
  - `API_BASE`（可选，覆盖 base_url）
  - `TOKEN`：C 端 JWT（`x-token`）；或 `SEED_MANIFEST` 中 `primary_user.token`
  - `ORDER_ID`（可选）：已发货且未取消、属于该用户的 `shop_order.id`，用于成功路径

本项目业务错误多为 **HTTP 200 + JSON `code`**：成功 `0`，失败 `7`，未登录 `401`。

## 1. 接口总览

| # | 方法 | 路径 | 鉴权 | 说明 |
|---|------|------|------|------|
| 1 | POST | `/order/confirmOrder` | JWT + Casbin（`x-token`） | C 端确认收货 |

管理端 `PUT /orderDelivery/updateOrderDelivery` 本版仅改服务层条件更新，路径不变，不在本文展开。

## 2. 接口明细

### 2.1 POST /order/confirmOrder — 确认收货

**鉴权**：`x-token: <jwt>`；角色须有 casbin `/order/confirmOrder` POST（与 cancelOrder 同角色）。

**请求 Body**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ID | uint | 是 | `shop_order.id`（不是发货单 id、不是 order_sn） |

```json
{ "ID": 123 }
```

**响应 200（业务成功）**

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | `0` |
| msg | string | `确认收货成功` |
| data | object | 可为空对象 |

**curl**

```bash
curl -sS "http://127.0.0.1:48888/order/confirmOrder" \
  -H "Content-Type: application/json" \
  -H "x-token: YOUR_TOKEN_HERE" \
  -d "{\"ID\": 123}"
```

**错误响应（仍为 HTTP 200）**

| JSON code | msg 示例 | 场景 |
|-----------|----------|------|
| 401 | 未登录或非法访问 | 无 token |
| 7 | 权限不足 | 角色无 casbin |
| 7 | 订单ID不能为空 | `ID` 为 0 |
| 7 | 订单不存在 | 无此单 |
| 7 | 无权操作 | 非下单人 |
| 7 | 订单已退款，不可确认收货 | `status_refund=2` |
| 7 | 订单状态不允许确认收货 | 未发货 / 已取消 / 已收货 / 重复确认 |

**断言要点**

- HTTP 200
- 无 token：`code === 401`
- `ID=0`：`code === 7` 且 msg 含「订单ID」
- 成功：`code === 0`，`msg === 确认收货成功`；再打一次同一 ID：`code === 7`（幂等失败，积分不双发）
- 成功后 `GET /order/findOrder?ID=` 该单 `status === 3` 且有 `receiveTime`

## 3. Python 脚本

| 脚本 | 说明 |
|------|------|
| `scripts/link_test/fulfillment/test_confirm_order_v030.py` | 覆盖本文 HTTP 接口 |

**运行**

```bash
pip install -r scripts/link_test/requirements.txt
set API_BASE=http://127.0.0.1:48888
set TOKEN=your_jwt
python scripts/link_test/fulfillment/test_confirm_order_v030.py
```

可选：`set ORDER_ID=123` 跑成功路径（会把该单写成已收货，勿对生产单使用）。

## 4. 与 seed 环境配合

当前仓库无 `scripts/seed/seed_manifest.json`。优先 `TOKEN` 环境变量；若日后有 manifest，脚本会读 `primary_user.token`。
