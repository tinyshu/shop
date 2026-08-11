# NOTICE

本项目（`new_shop`）是基于开源电商项目二次开发的**通用单商户微信小程序商城**脚手架。

## 上游来源

- 上游仓库：[fevrax/fresh-shop-group](https://github.com/fevrax/fresh-shop-group.git)
- 本地只读学习副本：仓库内 `old_shop/`（**请勿修改**；运行与交付以 `new_shop/` 为准）
- 上游声明协议以该仓库为准（README 标注 MIT）；本地 `old_shop` 副本中曾出现 Apache-2.0 文件，以官方仓库声明为准。

## 本仓库二次开发

- 产品目录：`new_shop/`（`server` / `web` / `uniapp` / `sql` / `docs`）
- 在保留同栈（Go Gin + Vue 管理端 + UniApp）的基础上，去行业品牌、配置隔离，并按 `docs/todo-gaps.md` 逐步补全生产能力。
- 对本目录中**新增与修改**的代码，版权归本项目贡献者；使用时请保留本 NOTICE 与根目录 `LICENSE`。

## 第三方组件

依赖各自许可证见各子项目 `go.mod` / `package.json` 及 `node_modules` 内声明（安装依赖后）。
