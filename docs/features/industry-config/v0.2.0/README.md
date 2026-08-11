# industry-config v0.2.0（M0-2 / GEN-03）

| 文件 | 说明 |
|------|------|
| [server/design.md](./server/design.md) | 审核开关接入（已实现） |
| [server/tasks.md](./server/tasks.md) | Server 任务 ✅ |
| [uniapp/design.md](./uniapp/design.md) | 小程序配合（已实现） |
| [uniapp/tasks.md](./uniapp/tasks.md) | UniApp 任务 ✅ |

**行为**：`feature.user_audit=0`（默认）时注册/登录生效态为已通过，可加购；改为 `1` 恢复小 B 审核拦截。

下一刀：M0-3 管理端开关，或 M1-1（PAY-02）。
