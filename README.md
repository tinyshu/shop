# 通用单商户商城（new_shop）

基于开源 [fresh-shop-group](https://github.com/fevrax/fresh-shop-group.git) 同栈拷贝改造的**通用微信小程序商城**产品目录。  
学习文档与缺口清单在 [`docs/`](./docs/)。只读对照源码在仓库 `old_shop/`（勿改）。

## 目录

| 路径 | 说明 |
|------|------|
| `server/` | Go / Gin 后端（模块路径暂仍为 `fresh-shop/server`，后续可整体改名） |
| `web/` | Vue3 管理端 |
| `uniapp/` | UniApp 微信小程序 |
| `sql/` | 数据库脚本 |
| `docs/` | 阶段一学习文档 + 阶段二说明 |
| `LICENSE` / `NOTICE.md` | MIT 与上游致谢 |

## 快速启动

1. **后端**：Go 1.23+、MySQL 8 → 见 [docs/run-backend-now.md](./docs/run-backend-now.md)（工作目录改为 `new_shop/server`，配置用 `config.example.yaml` 复制为 `config.yaml`）
2. **管理端**：见文档中 web 启动步骤；工作目录 `new_shop/web`
3. **小程序**：见 [docs/run-uniapp-win10.md](./docs/run-uniapp-win10.md)；用 HBuilderX 打开 `new_shop/uniapp`

每客户部署：独立数据库、独立小程序 AppId/Secret、独立支付商户号、独立 `jwt.signing-key`。详见 [docs/config-isolation.md](./docs/config-isolation.md)。

## 阶段进度

- **阶段一**：模块学习文档已完成（`docs/00`～`06`）
- **阶段二（进行中）**：本目录脚手架 — 拷贝、去品牌、配置隔离；支付/退款等增强见 [docs/todo-gaps.md](./docs/todo-gaps.md)
