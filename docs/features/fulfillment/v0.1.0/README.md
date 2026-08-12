# fulfillment v0.1.0（M2-1 / FUL-06 人工标记退款完成）

| 端 | 文档 | 状态 |
|----|------|------|
| server | [design.md](./server/design.md) / [tasks.md](./server/tasks.md) | **已实现** |
| web | [design.md](./web/design.md) / [tasks.md](./web/tasks.md) | **已实现** |

部署：执行 `sql/migrations/20260812_fulfillment_m2_1_mark_refund.sql`（注册 API + Casbin），重启后端后在订单列表使用「标记退款完成」。
