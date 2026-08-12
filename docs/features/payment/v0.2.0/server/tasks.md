# payment — server 任务清单（v0.2.0）

基于 [design.md](./design.md)。全局约束：只改 `new_shop/`；回调与查单共用 `MarkOrderPaidFromWechat`；本版不做流水表/自动退款。

另增加查单配置：`enable` / `spec` / `minAgeMinutes` / `batchSize` / `maxQueryPerOrder` / `mock` / `adminSyncEnable`。

---

## 执行顺序

1. ✅ 任务 1 — `WechatPay.Compensate` 配置
2. ✅ 任务 2 — `pay_mark.go` 条件更新入账
3. ✅ 任务 3 — `pay_query.go` 查单封装（含 mock）
4. ✅ 任务 4 — `pay_compensate.go` 单笔/扫描 + 次数帽
5. ✅ 任务 5 — `NotifyLogic` 改调 Mark
6. ✅ 任务 6 — `syncWechatPay` API + 路由
7. ✅ 任务 7 — 定时注册 compensate
8. ✅ 任务 8 — 单测 + `go build` + 文档进度

---

## 任务 1：config 查单配置 `✅`

文件：`new_shop/server/config/wechat.go`、`config.example.yaml`（本地 `config.yaml` 同步，gitignore）

### 1.1 结构体 `✅`

`WechatPayCompensate`：`Enable` / `Spec` / `MinAgeMinutes` / `BatchSize` / `MaxQueryPerOrder` / `Mock` / `AdminSyncEnable`

### 1.2 yaml 示例默认 `✅`

| 键 | 默认 | 说明 |
|----|------|------|
| enable | false | 定时总开关（测试建议关） |
| spec | @every 5m | cron |
| minAgeMinutes | 5 | 创建超过 N 分钟才扫 |
| batchSize | 20 | 每轮上限 |
| maxQueryPerOrder | 3 | 单笔进程内查微信次数；0=不限制 |
| mock | false | true=不调微信，视为未付 |
| adminSyncEnable | true | 管理端一键同步 |

---

## 任务 2：pay_mark.go `✅`

文件：`new_shop/server/service/wechat/pay_mark.go`（新建）

- `MarkPaidInput` + `MarkOrderPaidFromWechat`
- 金额校验复用 M1-1；`UPDATE ... WHERE order_sn=? AND status=0`
- 已取消 → `ErrOrderCancelled`；已支付幂等成功

---

## 任务 3：pay_query.go `✅`

文件：`new_shop/server/service/wechat/pay_query.go`（新建）

- `QueryOrderByOutTradeNo`：`mock` 短路；否则 `WxPay.GetOrder().QueryOrder`

---

## 任务 4：pay_compensate.go `✅`

文件：`new_shop/server/service/wechat/pay_compensate.go`（新建）

- `SyncWechatPayByOrder` / `compensateOne` / `RunCompensateScan`
- 进程内 `maxQueryPerOrder` 计数；已取消不入账

---

## 任务 5：NotifyLogic `✅`

文件：`new_shop/server/service/wechat/wechat.go`

- 组装 input → `MarkOrderPaidFromWechat`；`ErrOrderCancelled` 回 SUCCESS 停重试

---

## 任务 6：API + 路由 `✅`

文件：`api/v1/shop/order.go`、`router/shop/order.go`

- `POST /order/syncWechatPay` body：`orderId` 或 `orderSn`

---

## 任务 7：定时任务 `✅`

文件：`initialize/timer.go`

- `registerWechatPayCompensate`：仅 `compensate.enable=true` 时注册

---

## 任务 8：验证 `✅`

- `go test ./service/wechat/`
- `go build`
