# payment — server 任务清单（v0.1.0 / M1-1）

基于 [design.md](./design.md)。只改 `new_shop/server`；不做流水表/查单。

---

## 执行顺序

1. ✅ 任务 1 — `pay_amount.go` 期望分与比较
2. ✅ 任务 2 — `pay_amount_test.go`
3. ✅ 任务 3 — `NotifyLogic` 接入校验
4. ✅ 任务 4 — 编译验证 + 文档进度

---

## 任务 1：pay_amount.go `✅`

文件：`new_shop/server/service/wechat/pay_amount.go`（新建）

- `ExpectedPayFen(orderTotal float64, debug bool) (int64, error)`：debug→1；否则 `round(total*100)`；负/NaN/Inf → error
- `NotifyAmountMatches(orderTotal float64, totalFeeFen int64, debug bool) bool`

## 任务 2：单测 `✅`

文件：`pay_amount_test.go`

## 任务 3：NotifyLogic `✅`

文件：`wechat.go`

- TotalFee nil → error
- 金额不符 → Error 日志（order_sn、期望分、实付分）→ error，不 Save
- 已支付短路不变

## 任务 4：验证与文档 `✅`
