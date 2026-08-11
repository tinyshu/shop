# industry-config — uniapp 任务清单（v0.2.0 / M0-2）

基于 [design.md](./design.md)。信任 server 生效态；补本地缓存刷新。

---

## 执行顺序

1. ✅ 任务 1 — `utils/audit.js` 刷新 isAudit
2. ✅ 任务 2 — 商品/分类等页改用刷新；`onShow` 补刷
3. ✅ 任务 3 — loginPop 保持依赖服务端（确认即可）

---

## 任务 1：utils/audit.js `✅`

新建 `refreshUserAudit(vm)`：本地已是 1 则 true；否则拉 `getUserAuditStatus` 并 `setUser`。

## 任务 2：接入页面 `✅`

`detail` / `goods` / `category` / `index` / `quickPay` / `pointGoods`：用 util；有 token 时 `onShow` 再刷一次。

## 任务 3：loginPop `✅`

逻辑不变：`auditStatus !== 1` 才跳转 memberInfo（免审时服务端已为 1）。
