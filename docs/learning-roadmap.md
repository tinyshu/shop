# 框架与代码学习顺序（按模块）

> 目标：搞懂 **库表字段、整体架构、接口与流程**，为 `new_shop` 通用商城二次开发打底。  
> 源码只读：`old_shop/`。学习产出在 `new_shop/docs/`。  
> 你有 Go 经验：后端细读 service；前端以「会定位文件、会联调」为准。

---

## 怎么学（每模块固定三步）

每个模块按同一节奏，**一个模块过关再进下一个**：

1. **表**：该模块涉及的表 + 每个字段含义 / 枚举  
2. **接口**：路径、鉴权、谁调用（管理端 / 小程序）、关联表  
3. **流程**：关键时序或状态机（对照本地已跑通的三端验证）

配合方式：

- 说 **「开始模块 N」** → 我写该模块文档（如 `01-goods.md`）  
- 你阅读 + 对照管理端 / Swagger / 小程序 → 有问题就问  
- 说 **「下一模块」** → 再继续  

**不要**一次啃完所有代码细节。

---

## 学习总顺序（推荐）

| 次序 | 模块 | 文档 | 你要掌握什么 | 建议对照 |
|------|------|------|--------------|----------|
| 0 | 整体架构 | [00-overview.md](./00-overview.md) | 三端、分层、Public/Private、配置入口 | 已完成可复习 |
| 1 | 商品 | [01-goods.md](./01-goods.md) | 分类/品牌/标签/规格/图片表与接口 | 管理端商品；小程序列表详情 |
| 2 | 购物车 | `02-cart.md` | 加购、勾选、数量；与用户/商品关系 | 小程序购物车（需登录可后补） |
| 3a | 订单·下单支付 | [03a-order-pay.md](./03a-order-pay.md) | 建单、支付方式、回调、状态 0→1 | Swagger + 订单 service |
| 3b | 订单·履约退货 | `03b-order-fulfillment.md` | 发货、收货、退货、月结（可当开关学） | 管理端订单/发货 |
| 4 | 用户与登录 | `04-user-auth.md` | sys_users、微信登录、地址、收藏、审核字段 | 管理端用户；小程序登录逻辑 |
| 5 | 账户资金 | `05-account.md` | 余额/积分账户、充值、流水 | 管理端账户相关菜单 |
| 6 | 运营与系统 | `06-ops-system.md` | Banner、配送员、Casbin/菜单/字典/配置 | 管理端权限与设置 |
| — | 商业定位 | [business-positioning.md](./business-positioning.md) | 行业/客户/获客（非代码，随时可看） | — |

订单拆成 **3a / 3b** 两轮，避免一次过重；若你希望合并成一个 `03-order.md` 也可以说一声。

---

## 各模块学习要点

### 模块 0：整体架构（已有文档）

- 目录：`server` / `web` / `fresh-shop-uniapp`  
- 分层：`router → api → service → model`  
- 鉴权：JWT（`x-token`）+ Casbin  
- 入口：`main.go`、`config.yaml`、`initialize/router.go`  

**过关**：能画出三端如何连后端，并说出公开接口与需登录接口的差别。

---

### 模块 1：商品（下一个建议开始）

**表（优先）**

- `shop_category`、`shop_brand`、`shop_brand_category`  
- `shop_tags`、`shop_tags_goods`  
- `shop_goods`、`shop_goods_description`、`shop_goods_image`  
- `shop_goods_spec` / `_item` / `_value`  

**接口**

- 公开：`getGoodsList`、`findGoods`、分类/品牌/标签列表  
- 管理：商品 CRUD、导入导出等  

**流程**

- 单规格 vs 多规格如何落到表  
- 上架/首页/热销等标记如何影响列表  

**源码锚点**

- `server/router/shop/goods.go`、`category.go`、`brand.go`…  
- `server/model/shop/goods*.go`  
- `web/src/view/shop/goods/`、`web/src/api/goods.js`  
- `fresh-shop-uniapp/pages/goods/`、`api/goods.js`  

---

### 模块 2：购物车

**表**：`shop_cart`  

**接口**：`createCart`、`updateCart`、`getCartList`、全选/单选等  

**流程**：登录用户加购 → 勾选 → 进入下单（衔接到模块 3）  

**源码锚点**：`server/router/shop/cart.go`、`service/shop` 购物车；小程序 `pages/cart/`  

---

### 模块 3a：下单与支付

**表**：`shop_order`、`shop_order_details`（先盯这两张）  

**接口**：`createOrder`、`orderPay`、微信 `createPayData` / `pay/notify`  

**流程**：建单 → 调起支付 → 回调改状态；弄清 `status`、`payment`、`settlement_type`、`goods_area`  

**源码锚点**：`server/router/shop/order.go`、`router/wechat/wechat.go`；重点读 `service` 里下单/支付  

---

### 模块 3b：履约与退货

**表**：`shop_order_delivery`、`shop_order_return`、`shop_order_return_details`  

**接口**：发货/确认收货、退货、后台 `batchSettlement`（月结）  

**流程**：已支付 → 发货 → 收货；退款状态机；哪些是冻品/小 B 特色（二次开发可开关）  

---

### 模块 4：用户与登录

**表**：`sys_users`（含 `open_id`、`audit_status`）、`user_address`、`shop_favorites`  

**接口**：`/base/login`、`/base/loginWx`、`/wechat/code2Session`、地址与收藏 CRUD  

**流程**：管理端账号登录 vs 小程序微信手机号登录；审核字段对通用 B2C 的影响  

---

### 模块 5：账户资金

**表**：`user_account`、`user_account_group`、`sys_recharge`、`user_finance_*`  

**接口**：账户查询、充值、流水  

**流程**：余额支付与微信支付如何并存（和订单 `payment` 对照）  

---

### 模块 6：运营与系统（可粗看）

**表**：`shop_banner`、`user_delivery`；`sys_*` 菜单角色、`casbin_rule`、字典、`sys_config`  

**接口**：Banner、配送员；权限菜单 API  

**流程**：管理端动态菜单如何出来；Casbin 与 Private 路由的关系  

有 Go/后台经验可加快；接单部署时配置、权限会用到。

---

## 建议时间分配（灵活）

| 模块 | 相对权重 | 说明 |
|------|----------|------|
| 0 架构 | 已完成 | 复习即可 |
| 1 商品 | ★★★★ | 主数据，必扎实 |
| 2 购物车 | ★★ | 短，衔接订单 |
| 3a+3b 订单 | ★★★★★ | 接单核心 |
| 4 用户登录 | ★★★ | 微信配置与审核开关 |
| 5 账户 | ★★ | 余额场景用得到再深挖 |
| 6 运营系统 | ★★ | 粗通即可 |

---

## 环境与辅助文档

| 文档 | 用途 |
|------|------|
| [run-backend-now.md](./run-backend-now.md) | 后端启动 |
| [run-uniapp-win10.md](./run-uniapp-win10.md) | 小程序运行 |
| [business-positioning.md](./business-positioning.md) | 行业与获客（并行阅读） |
| [00-overview.md](./00-overview.md) | 架构总览 |

本地保持：**后端 + 管理端** 常开；小程序用于看 C 端列表即可，登录可放到模块 4。

---

## 模块过关自测（总览）

- [ ] 商品从分类到 SKU 涉及哪些表？  
- [ ] 购物车一行数据关键字段是什么？  
- [ ] 订单从创建到支付成功，状态与接口顺序？  
- [ ] 发货/退货改哪些表、哪些状态字段？  
- [ ] 小程序登录打哪些接口？`audit_status` 是什么？  
- [ ] 余额支付和微信支付在订单上如何区分？  
- [ ] 管理端菜单权限大致靠哪些表？  

全部能答后 → 进入阶段二：在 `new_shop` 落地通用商城二次开发。

---

## 下一步

**阶段二进行中**：产品代码已在 `new_shop/server|web|uniapp`。见 [根 README](../README.md)、[debrand-checklist.md](./debrand-checklist.md)、[config-isolation.md](./config-isolation.md)。

**按模块细切片推进**：[phase2-module-roadmap.md](./phase2-module-roadmap.md)。**M0（含 M0-3）与 payment v0.2.0 已完成**；PAY-04/05 暂缓。
