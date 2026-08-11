# 配置隔离说明（每客户独立部署）

> 本产品为**单商户独立部署**：一客户一套库、一套小程序凭证、一套支付商户号。不是多租户 SaaS。

## 后端 `server/`

| 文件 | 用途 |
|------|------|
| [`config.example.yaml`](../server/config.example.yaml) | 模板（无真实密码）；复制为 `config.yaml` |
| `config.yaml` | **本地/服务器机密**，已在 [`.gitignore`](../.gitignore) 中忽略 |

必改项：

- `mysql` / `db-list`：库名、账号密码（建议库名如 `new_shop` 或客户名）
- `jwt.signing-key`：每客户随机 UUID，勿共用
- `wechat.appid` / `secret`：客户小程序
- `wechatPay.*`：客户商户号、密钥、证书、`notifyUrl`
- `system.env`：生产勿用 `develop`（否则 Casbin 放行，见 [06-ops-system.md](./06-ops-system.md)）

启动前：

```text
cd new_shop/server
copy config.example.yaml config.yaml   # Windows
# 编辑 config.yaml 后：
go run .
```

## 管理端 `web/`

| 文件 | 用途 |
|------|------|
| [`.env.example`](../web/.env.example) | 模板 |
| `.env.development` | 本地开发（已指向 127.0.0.1:48888） |
| `.env.production` | 生产 `VITE_BASE_PATH` 改为客户域名 |

## 小程序 `uniapp/`

| 文件 | 用途 |
|------|------|
| [`config/config.example.js`](../uniapp/config/config.example.js) | 模板说明 |
| `config/config.js` | 默认 `http://localhost:48888`；上线改客户 API 域名 |
| `manifest.json` | 填入客户微信小程序 AppId |

## 交付检查

- [ ] 独立 MySQL 库已导入 `sql/`
- [ ] `config.yaml` 未提交到公开仓库
- [ ] JWT / 微信 / 支付均为该客户凭证
- [ ] 管理端与小程序 API 地址指向该客户后端
