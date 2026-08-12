# 通用商城文档（new_shop）

- **阶段一**：学习 `old_shop/`（只读）→ 文档在本目录 `00`～`06`
- **阶段二**：产品代码在 [`../`](../)（`server` / `web` / `uniapp`），说明见 [根 README](../README.md)、[config-isolation.md](./config-isolation.md)、[debrand-checklist.md](./debrand-checklist.md)

源码只读参考：`old_shop/`（勿改）。运行与改造请用 **`new_shop/`**。

## 商业与定位

| 文档 | 说明 |
|------|------|
| [business-positioning.md](./business-positioning.md) | **行业认知、目标客户、服务模式、获客与阶段目标** |

## 阶段二推进

| 文档 | 说明 |
|------|------|
| [phase2-module-roadmap.md](./phase2-module-roadmap.md) | **二次开发细切片顺序、过关标准、与 todo-gaps 映射** ← 阶段二从这里推进 |
| [features/industry-config/](./features/industry-config/) | M0 行业开关（v0.1.0～**v0.3.0 管理端** 已实现） |
| [features/payment/](./features/payment/) | M1 支付：v0.1.0 / **v0.2.0 已实现**（条件更新+掉单查单） |
| [features/fulfillment/](./features/fulfillment/) | M2 履约：v0.1.0 / **v0.2.0 已实现**（标记退款 + 发货/取消条件更新） |

## 学习路线

| 文档 | 说明 |
|------|------|
| [learning-roadmap.md](./learning-roadmap.md) | **框架/代码学习顺序（按模块）** ← 学代码从这里看 |
| [payment-reliability.md](./payment-reliability.md) | **支付异常现状与二次开发补全清单**（掉单/对账/退款等） |
| [todo-gaps.md](./todo-gaps.md) | **功能缺失 / 待补 Todo**（含微信自动退款等，接单前对照） |
| [debrand-checklist.md](./debrand-checklist.md) | **阶段二去品牌清单** |
| [config-isolation.md](./config-isolation.md) | **每客户配置隔离（yaml / env / 小程序）** |
| [smoke-scaffold.md](./smoke-scaffold.md) | **阶段二脚手架冒烟记录** |

## 环境文档

| 文档 | 说明 |
|------|------|
| [run-backend-now.md](./run-backend-now.md) | **当前推荐：Go 1.23 + MySQL 8 后端启动步骤**（按序操作） |
| [run-uniapp-win10.md](./run-uniapp-win10.md) | **Win10 运行微信小程序客户端**（HBuilderX + 开发者工具） |
| [run-old-server-local.md](./run-old-server-local.md) | 环境说明与排错补充（含旧版升级说明） |

## 你怎么用

每轮只读 **当前模块** → 有问题就问 → 说「下一模块」再继续。  
你要掌握：表字段、架构、接口与流程；实现细节由文档代劳。

## 进度

| 次序 | 模块 | 文档 | 状态 |
|------|------|------|------|
| 0 | 整体架构总览 | [00-overview.md](./00-overview.md) | 已有 |
| 1 | 商品 | [01-goods.md](./01-goods.md) | 已写出，待你确认 |
| 2 | 购物车 | [02-cart.md](./02-cart.md) | 已写出，待你确认 |
| 3a | 订单·下单支付 | [03a-order-pay.md](./03a-order-pay.md) | 已写出，待你确认 |
| 3b | 订单·履约退货 | [03b-order-fulfillment.md](./03b-order-fulfillment.md) | 已写出，待你确认 |
| 4 | 用户/登录/地址/收藏 | [04-user-auth.md](./04-user-auth.md) | 已写出，待你确认 |
| 5 | 账户/余额/充值 | [05-account.md](./05-account.md) | 已写出，待你确认 |
| 6 | 运营与系统权限 | [06-ops-system.md](./06-ops-system.md) | 已写出，待你确认 |

详细说明与过关标准见 [learning-roadmap.md](./learning-roadmap.md)。

## 读完模块0可自测

1. 三端分别是什么、谁调谁？
2. 后端请求从路由到数据库大致经过哪几层？
3. Public 与 Private 有何区别？请求头带什么？
4. 本地默认后端端口是多少？
