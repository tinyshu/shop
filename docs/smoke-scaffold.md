# 阶段二脚手架冒烟记录

日期：执行脚手架当日

| 检查项 | 结果 |
|--------|------|
| `go build`（`new_shop/server`） | 通过 |
| `GET http://127.0.0.1:48888/health` | 返回 `ok`（使用本地 `config.yaml` + 已有 MySQL） |
| 管理端产品名 | `web/src/core/config.js` → `通用商城`（需 `npm i && npm run serve` 看页面；脚手架未强制装依赖） |
| 小程序工程 | `uniapp/manifest.json` 名为「通用商城」；用 HBuilderX 打开 `new_shop/uniapp` 即可（微信登录待配真实 AppId） |

说明：完整管理端 UI 冒烟需在 `web/` 执行 `npm install`；未纳入本次强制步骤以免长时间安装阻塞。后端编译与 health 已验证 `new_shop/server` 可独立运行。
