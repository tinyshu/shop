# 去品牌清单（阶段二脚手架）

> 仅改 `new_shop/`，**未改** `old_shop/`。产品名暂用「通用商城」。

## 已改

| 位置 | 变更 |
|------|------|
| `web/src/core/config.js` | appName / 欢迎语 |
| `web/index.html` | keywords |
| `web/.../goods.vue` | 导出 Excel 文件名 |
| `server/core/server.go` | 启动欢迎语 |
| `server/service/system/sys_user.go` | 微信首登昵称前缀「用户」+尾号 |
| `server/model/system/sys_user.go` 等 | 去掉 `fs.ssooai.com` 默认头像 |
| `server/source/system/user.go` | 种子头像清空 |
| `uniapp/manifest.json` | name / description |
| `uniapp/pages.json` | 导航标题 |
| `uniapp/pages/index|my|memberInfo|category` | 分享标题、页标题、默认 logo |
| `uniapp/README.md` | 说明文案 |

## 未改 / 后续

| 项 | 说明 |
|----|------|
| `sql/*.sql` 种子类目名 | 仍可能含冻品类目文案，接单时换类目即可 |
| `docs/` 学习文档中的「启运/冻品」叙述 | 历史学习上下文，保留 |
| Go module 名 `fresh-shop/server` | 故意未全局 rename，避免第一次大范围替换；见根 README |
| 小程序 `manifest.appid` | 仍为原 `__UNI__*` 占位，客户部署时换成真实微信 AppId |
| 上传目录历史图片 | 若库里已有旧 CDN URL，需数据迁移，非脚手架范围 |

检索命令（应无业务命中）：在 `new_shop` 下对 `*.{go,vue,js,json,html}` 搜 `启运|ssooai`。
