# payment v0.1.0（M1-1 / PAY-02）

| 文件 | 说明 |
|------|------|
| [server/design.md](./server/design.md) | 回调金额校验（已实现） |
| [server/tasks.md](./server/tasks.md) | 任务清单 ✅ |

**落地**：`service/wechat/pay_amount.go` + `NotifyLogic` 入账前校验。

下一刀：M1-2（支付流水 + 幂等）或 M0-3。
