# industry-config — server 任务清单（v0.2.0 / M0-2）

基于 [design.md](./design.md)。全局约束：只改 `new_shop/`；不改 CreateOrder 审核校验；用 `common.FeatureEnabled(KeyUserAudit, false)`。

---

## 执行顺序

1. ✅ 任务 1 — feature_config 辅助函数
2. ✅ 任务 2 — LoginWx 注册默认态 + 登录回写
3. ✅ 任务 3 — GetAuditStatus / TokenNext 生效态
4. ✅ 任务 4 — 单测与编译

---

## 任务 1：feature_config.go `✅`

文件：`new_shop/server/service/common/feature_config.go`

- `UserAuditRequired() bool`
- `EffectiveAuditStatus(stored int8) int8`（免审恒为 1）

## 任务 2：sys_user.go LoginWx `✅`

文件：`new_shop/server/service/system/sys_user.go`

- 新用户：`AuditStatus = 1`（免审）或 `0`（需审）
- 老用户登录后：免审且 ≠1 则 UPDATE 并改内存字段

## 任务 3：sys_user API `✅`

文件：`new_shop/server/api/v1/system/sys_user.go`

- `GetAuditStatus`：免审返回 `auditStatus=1`（userId=0 仍为 0）
- `TokenNext`：claims/响应用 `EffectiveAuditStatus`

## 任务 4：验证 `✅`

- `go test ./service/common/`
- `go build`
