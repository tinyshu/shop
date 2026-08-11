# industry-config v0.1.0

| 文件 | 说明 |
|------|------|
| [server/design.md](./server/design.md) | M0-1 设计（已通过） |
| [server/tasks.md](./server/tasks.md) | 任务清单（已全部 ✅） |

**已落地**：
- `server/model/system/sys_config.go` — `name` size 64
- `server/service/common/feature_config.go` — Feature* 封装
- `sql/migrations/20260811_industry_config_m0_1.sql` — 已有库迁移
- `sql/fresh-shop.sql` — 新库种子已含三键

**已有库请执行**：`new_shop/sql/migrations/20260811_industry_config_m0_1.sql`

下一切片：说「按 roadmap 做 M0-2」或「只做 M1-1」。
