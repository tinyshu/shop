# 二次开发待补清单（Todo Gaps）

> **不是学习必做项**，是对照 `old_shop` 现状：功能缺失 / 半成品，接单或阶段二改造时按优先级补。  
> 支付细节展开：[payment-reliability.md](./payment-reliability.md)  
> 履约/售后上下文：[03b-order-fulfillment.md](./03b-order-fulfillment.md)

约定：

- **P0**：上线收钱前强烈建议  
- **P1**：零售 B2C 交付建议有；可报价为加项  
- **P2**：体验/行业定制，可后做  
- **首期可不做**：合同写清人工流程即可暂缓  

---

## 支付与资金

| ID | 项 | 现状 | 建议优先级 | 备注 |
|----|----|------|------------|------|
| PAY-01 | 掉单补偿（主动查单补 `status=1`） | **v0.2.0 已做** | P0 | 见 payment-reliability / features/payment/v0.2.0 |
| PAY-02 | 回调金额与订单金额校验 | 无 | P0 | |
| PAY-03 | 支付流水表 / 幂等条件更新 | 条件更新 **v0.2.0 已做**；流水表刻意后置 | P0 | |
| PAY-04 | 超时关单 + 与支付成功冲突处理 | 字段有、逻辑无 | P1 | 含回库存 |
| PAY-05 | **微信自动退款闭环** | 无 API；仅能改状态/人工商户平台退 | **P1**（零售必谈） | 见下方专节 |
| PAY-06 | 已支付取消：退款 + 回库存 | 注释有、实现无 | P1 | 与 PAY-05 同源 |

---

## 售后 / 履约

| ID | 项 | 现状 | 建议优先级 | 备注 |
|----|----|------|------------|------|
| FUL-01 | 售后：申请 → `status_refund` → 微信退款 → 回库存 | CRUD 脚手架 | **P1** | 与 PAY-05 同一能力 |
| FUL-06 | **人工退款后「标记退款完成」专用接口** | 仅靠通用 `updateOrder` / `updateOrderReturn` 改字段 | **P1**（人工退方案标配） | 见专节；可比自动退款先做 |
| FUL-02 | 小程序确认收货 API 对齐 | 前端打 `/order/confirmOrder`，后端无此路由 | P1 | 应对齐到 `PUT /orderDelivery/updateOrderDelivery` |
| FUL-03 | `OrderService.OrderDeliver` 空实现 | 发货实际走 delivery | P2 | 删或转发，避免误用 |
| FUL-04 | 快递单号等通用物流字段 | 现为配送员城配模型 | P2 | 按客户行业开关 |
| FUL-05 | **发货 / 取消竞态**（TOCTOU） | 查条件与更新不在同一原子操作 | **P1** | 见下方专节；可用乐观锁 `version` |
| FUL-06 | **「标记退款完成」专用接口**（人工打款后回写状态） | 仅靠通用 `updateOrder` / `updateOrderReturn` 改字段 | **P1** | 见下方专节；可先于自动微信退款落地 |

---

## 专节：微信自动退款（PAY-05 / FUL-01）

**缺失含义**：表和后台能记「售后/退款状态」，**不会**调微信把钱退回用户；现在只能人工在商户平台退款再改库。

**目标闭环（二次开发时）：**

```text
申请售后 → 订单 status_refund=退款中
  → 审核同意 → 调微信退款 API（原单号/金额）
  → 成功 → status_refund=已退 + 售后单完成 + 视规则回库存
  → 失败 → 标记失败 + 告警/人工
```

**何时必须做**：面向普通消费者、微信实时收款的 B2C 交付。  
**何时可暂缓**：客户书面接受「商户平台人工退款」；或以月结对账冲减为主。  
**报价**：标准版可写人工退；「一键微信退款」单列加项或增强包。

### 过渡方案：人工打款 + 标记退款完成（FUL-06）

即使不调微信退款 API，钱在商户平台退完后，**系统内仍要改状态**，否则列表/统计仍当正常单。

现状：靠 `PUT /order/updateOrder`、`PUT /orderReturn/updateOrderReturn` 通用编辑改 `statusRefund` 等——能凑合，但易误改其它字段、无校验、通常不回库存、无操作审计。

更稳妥的二次开发（**可先于 PAY-05 落地**）：

```text
运营在微信商户平台完成退款
  → 管理端点「标记退款完成」（专用接口，非整单 Save）
  → 仅允许：订单 status_refund → 已退款；售后单同步完成
  → 可选：按规则回库存、记操作人/时间
  → 拒绝非法状态跃迁（如未支付单、已标记过等）
```

与 PAY-05 关系：FUL-06 = 人工退钱 + 系统状态闭环；PAY-05 = 系统直接调微信打钱。标准交付可只做 FUL-06，自动退款作增强包。

---

## 专节：发货 / 取消竞态（FUL-05）

**问题**：`CreateOrderDelivery` 先在事务外 `First(status=1 AND status_cancel=0)`，再在事务里 `Save` 整单；`CancelOrder` 同样是读后写。中间若用户取消成功，发货仍可能把 `status` 写成 2，并用内存里的 `status_cancel=0` **覆盖**已取消标记。

```text
发货 First 通过 → 用户 Cancel 写 status_cancel=1 → 发货 Save(status=2, cancel=0)
→ 已发货且取消被冲掉
```

**二次开发方向（择一或组合）**：

1. **乐观锁（偏好）**：`shop_order` 增加 `version`（或复用已有版本字段）；更新时 `WHERE id=? AND version=?`，成功则 `version+1`，冲突则失败重试/提示。发货、取消、支付回调改状态都走同一套。  
2. **条件更新**：发货 `UPDATE ... SET status=2 WHERE id=? AND status=1 AND status_cancel=0`，`RowsAffected==0` 则拒绝；取消对称 `AND status<2 AND status_cancel=0`。  
3. 可选：事务内 `SELECT ... FOR UPDATE`（悲观锁），并发低时也够用。

**涉及代码**：`service/shop/order_delivery.go` → `CreateOrderDelivery`；`service/shop/order.go` → `CancelOrder`（支付回调改状态同类问题见 PAY-03）。

**优先级**：P1（状态机严谨）；表结构变更需迁移，放二次开发完善，学习阶段不必改 `old_shop`。

---

## 其它（学习后边写边加）

| ID | 项 | 现状 | 建议优先级 |
|----|----|------|------------|
| GEN-01 | 多规格购物车/下单 | `spec_item_id` 常写 0 | 按客户 |
| GEN-02 | 去冻品品牌 / 配置隔离 | 阶段二脚手架 | **脚手架已做**（见 debrand / config-isolation） |
| GEN-03 | 用户审核 `audit_status` 可开关 | 小 B 特色 | 按客户 |
| AUTH-01 | **`session_key` 不返回客户端** | `code2Session` 把 key 回给 UniApp，`loginWx` 再回传 | **P1**（登录安全） | 见下方专节 |
| AUTH-02 | **LoginWx 查用户：phone 索引 / 或 openid 方案** | 按 `phone` 查且无索引、无唯一；并发可重复用户 | **P1** | 见下方专节 |

模块 4～6 读完后，把新发现的缺口追加到本表。

### 专节：session_key 仅留服务端（AUTH-01）

**问题**：理想模型是 `session_key` 只存在 Gin（解密手机号用）；现状 `GET /wechat/code2Session` 把 `openid + session_key` 返回前端，`POST /base/loginWx` 再让前端把 `sessionKey` 传回。key 在客户端多暴露一层，可被截获后用于解密同会话敏感数据。

**二次开发方向**：

```text
code2Session：服务端用 code 换 session_key，写入服务端缓存（Redis，按 openid 或自建 loginTicket，短 TTL）
  → 响应给前端：openid（或仅 loginTicket），不返回 session_key
loginWx：前端只传 encryptedData、iv、openid/ticket（及必要字段）
  → 服务端取出缓存的 session_key 解密 → 查/注册用户 → 发 JWT → 删除/失效该 key
```

**涉及**：`api/v1/wechat`（Code2Session 响应）、`LoginWx` 入参、小程序登录页；见 [04-user-auth.md](./04-user-auth.md)。

**优先级**：P1（安全加固）；可与阶段二脚手架一起做，学习阶段不必改 `old_shop`。

### 专节：LoginWx 用户查找与唯一性（AUTH-02）

**现状**：`LoginWx` / `LoginByPhone` 用 `WHERE phone = ?` 判断「有没有这个用户」（不是判断是否已登录）。模型里 `phone` **无索引、无唯一**；`open_id` 有普通索引但未作查用户主路径。并发双请求都 First 不到时可能插入两条同号用户。

**说明**：openid 维度是「微信用户 × 某一个小程序 AppId」；不同客户小程序对同一微信用户 openid **不同**。单商户独立部署下用 openid 足够；跨多个 AppId 打通才需要 unionid。

---

#### 方案 A：继续用手机号（可保留，但必须补索引）

适用：业务就是「一个手机号一个会员」，客服/后台常按手机找人。

| 项 | 建议 |
|----|------|
| 普通索引 | `phone` 必加，否则用户多时 loginWx 变慢 |
| 唯一约束 | 非空手机唯一；空手机用 `NULL`（勿用空串占唯一槽），后台未绑手机用户可多条 NULL |
| 并发 | 仅 First→Insert 不够；唯一冲突时改为按 phone 再查登录；或事务内锁 |
| Register | 插入撞唯一 → 转登录，避免报错中断 |

最低交付：**索引 + 唯一（处理好空值）**；逻辑可仍按 phone。

---

#### 方案 B：openid 为主查找（完整优化建议）

适用：登录性能/微信身份更清晰；与「本小程序用户」一致；`open_id` 已有索引基础。

**库表**

1. `open_id`：**唯一索引**（C 端微信用户必填；管理端账号无 openid 用 NULL）  
2. `phone`：仍建议 **唯一（非空）+ 普通索引**——支付/客服/合并账号仍常用手机  
3. 可选预留 `union_id` 字段（暂不用，避免以后改表）

**登录主路径**

```text
code2Session 得 openid（+ session_key，配合 AUTH-01 只留服务端）
  → 用户授权手机号 → 解密得 phone
  → 查用户优先：WHERE open_id = ?
       · 命中：更新 phone（若变更）、登录时间 → 发 JWT
       · 未命中：再 WHERE phone = ?（合并历史按手机注册的老数据）
            · 命中：补写 open_id（若空）→ 发 JWT
            · 未命中：Register（写入 open_id + phone + 角色 1000 等）→ 发 JWT
```

**合并与冲突规则（必须写清）**

| 情况 | 处理 |
|------|------|
| 同 openid，手机变了 | 以微信最新手机为准更新；注意唯一约束冲突要提示/人工 |
| 同 phone，不同 openid | 极少见（换号段/脏数据）；拒绝自动合并或走人工，防账号被绑走 |
| 仅有 phone 老用户、无 openid | 首次微信登录按 phone 命中后回填 open_id |
| 管理端用户 | 不走 LoginWx；勿与 C 端抢 open_id 唯一 |

**并发**：对 `open_id` / `phone` 唯一约束兜底；Insert 冲突 → 再查一次登录。

**与方案 A 关系**：A 是小改（索引+唯一）；B 是查用户主路径升级。可先做 A 止血，再演进 B。做 B 时 A 的 phone 唯一仍建议保留。

**涉及**：`model/system/sys_user.go`、迁移 SQL、`LoginWx`/`Register`/`LoginByPhone`；见 [04-user-auth.md](./04-user-auth.md)。

**优先级**：P1；索引/唯一建议阶段二或首单前做；openid 主路径可列为登录增强。

---

## 与阶段的关系

| 阶段 | 对本清单 |
|------|----------|
| 阶段一学习 | 知道缺口即可，不必实现 |
| 阶段二通用脚手架 | 至少文档与开关位；PAY-01/02 优先评估 |
| 首单交付 | 按合同勾选；未做项写进交付说明 |
| 零售增强包 | PAY-05 + FUL-01 + FUL-02；状态机严谨含 FUL-05；人工退至少 FUL-06；登录 AUTH-01/02 |
