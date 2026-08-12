# industry-config — web 任务清单（v0.3.0 / M0-3）

基于 [design.md](./design.md)。全局约束：只改 `new_shop/`；不新增 sysConfig HTTP；菜单用幂等 SQL。

---

## 执行顺序

1. ✅ 任务 1 — 功能开关页 `featureConfig.vue`
2. ✅ 任务 2 — 菜单迁移 SQL
3. ✅ 任务 3 — 模块文档进度

---

## 任务 1：featureConfig.vue `✅`

文件：`new_shop/web/src/view/superAdmin/featureConfig/featureConfig.vue`

- 拉 `groupType=feature`；Switch / Radio；`updateSysConfig` 即时保存

## 任务 2：菜单 SQL `✅`

文件：`new_shop/sql/migrations/20260812_industry_config_m0_3_menu.sql`

## 任务 3：文档 `✅`

模块 README / roadmap / phase2 进度
