# industry-config — server 任务清单

基于 [design.md](./design.md)（M0-1）。  
全局约束：只改 `new_shop/`；不改 LoginWx/下单；不改 `GetSysConfig` 缺键报错语义；后端 `service/common` 分层。

---

## 执行顺序

1. ✅ 任务 1 — model 扩 `name` 列宽（无依赖）
   - ✅ 1.1 修改 gorm size
2. ✅ 任务 2 — 迁移 + 种子 SQL（依赖任务 1）
   - ✅ 2.1 ALTER + 幂等 INSERT
3. ✅ 任务 3 — `feature_config.go`（依赖任务 1）
   - ✅ 3.1 常量 + FeatureEnabled / FeatureString
4. ✅ 任务 4 — 单测解析与默认值（依赖任务 3）
5. ✅ 任务 5 — 编译验证 + 文档进度（依赖 1～4）

---

## 任务 1：sys_config.go — name size 20→64 `✅ 已完成`

文件：`new_shop/server/model/system/sys_config.go`  
改动类型：修改

### 1.1 修改 Name 字段 gorm size `✅`

将 `size:20` 改为 `size:64`。

---

## 任务 2：SQL 迁移与种子 `✅ 已完成`

文件：`new_shop/sql/migrations/20260811_industry_config_m0_1.sql`（新建）  
改动类型：新建

### 2.1 ALTER + 三键 INSERT（幂等）`✅`

```sql
ALTER TABLE sys_config MODIFY COLUMN name varchar(64) ...;
INSERT ... WHERE NOT EXISTS (SELECT 1 FROM sys_config WHERE name = 'feature.user_audit');
-- settle_month / courier_mode 同理
```

默认：`feature.user_audit=0`，`feature.settle_month=0`，`feature.courier_mode=courier`，`group_type=feature`，`status=1`。

---

## 任务 3：feature_config.go — 读取封装 `✅ 已完成`

文件：`new_shop/server/service/common/feature_config.go`（新建）  
改动类型：新建

### 3.1 键常量与 Feature* `✅`

- 常量：`KeyUserAudit`、`KeySettleMonth`、`KeyCourierMode`
- `FeatureEnabled(key, defaultEnabled bool) bool`
- `FeatureString(key, defaultVal string) string`
- 缺键 / status≠1 / status nil / 非法布尔 → 默认值
- 真值：`1` / `true` / `on`（大小写不敏感）

---

## 任务 4：feature_config_test.go `✅ 已完成`

文件：`new_shop/server/service/common/feature_config_test.go`（新建）  
改动类型：新建

### 4.1 测 parseFeatureBool；无 DB 时 Feature* 返回默认 `✅`

---

## 任务 5：编译与文档 `✅ 已完成`

### 5.1 `go build` / `go test` 相关包 `✅`

### 5.2 更新 design 评审状态、roadmap、phase2 进度、v0.1.0 README `✅`
