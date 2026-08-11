# 老版本（old_shop）本地运行服务端

本文只覆盖 **Go 后端 API** 在本机（含 Win10）跑通。管理端 / 小程序另文说明。  
源码只读：`old_shop/`，**不要改**业务代码；配置按你本机改 `config.yaml` 即可。

---

## 1. 环境要求

| 软件 | 最低要求 | 推荐 | 说明 |
|------|----------|------|------|
| Go | **1.23+** | 1.23 / 1.24 | `go.mod` 写明 `go 1.23`，更低版本无法按官方依赖编译 |
| MySQL | **5.7+** | **8.0** | SQL 使用 `datetime(3)`、`utf8mb4`、InnoDB 等，**5.5 不能导入** |
| Redis | — | 可选 | `use-redis: false` 时可先不装 |
| Git | — | 任意 | 可选 |

在 PowerShell 检查：

```powershell
go version
mysql --version
```

### 1.1 版本不达标时（务必先升级）

若你当前类似：

```text
go version go1.18.1 windows/amd64
mysql  Ver 14.14 Distrib 5.5.51   ← 太旧
```

| 现状 | 结论 | 建议 |
|------|------|------|
| Go 1.18 | **不满足** | 安装 [Go 1.23+](https://go.dev/dl/)，装完重开终端，再执行 `go version` 确认 |
| MySQL 5.5 | **不满足** | 安装 **MySQL 8.0**（或至少 5.7）。不要指望改 SQL 适配 5.5 |

**Go 升级注意（Win10）：**

1. 从官网下载 Windows 安装包（amd64）。
2. 安装后**新开**一个 PowerShell / CMD，再查版本（旧窗口仍可能是 1.18）。
3. 若仍显示 1.18：检查环境变量 `PATH` 是否指向新 Go 的 `bin`（常见 `C:\Program Files\Go\bin`），去掉旧路径。
4. 可同时装多个版本，但 `PATH` 里生效的必须是 1.23+。

**MySQL 升级注意：**

1. 5.5 → 8.0 **没有**可靠「原地小改就能用」的捷径，建议新装 MySQL 8.0。
2. 若本机还要保留 5.5：可用 **Docker** 只跑一个 MySQL 8 容器（端口 3306 或改成 3307，并在 `config.yaml` 写对应端口）。
3. 安装 8.0 后用新客户端执行 `mysql --version`，确认不再是 5.5，再建库导入 `fresh-shop.sql`。

满足上表后再往下做第 3 节。

---

## 2. 目录位置

```text
old_shop/
├── fresh-shop/
│   ├── server/          ← 后端（本文）
│   │   ├── main.go
│   │   ├── config.yaml  ← 必改配置
│   │   └── go.mod
│   └── sql/
│       └── fresh-shop.sql   ← 建库脚本
└── fresh-shop-uniapp/   ← 小程序（本文不涉及）
```

工作目录（按你机器路径）：

```text
d:\tiny_zimeiti\cursor_dev_pro\shop\old_shop\fresh-shop\server
```

---

## 3. 准备数据库

### 3.1 创建空库

用 MySQL 客户端或命令行（库名可改，需与 `config.yaml` 一致）：

```sql
CREATE DATABASE IF NOT EXISTS `fresh-shop`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;
```

### 3.2 导入脚本

在 PowerShell 中（按实际账号、路径改）：

```powershell
cd d:\tiny_zimeiti\cursor_dev_pro\shop\old_shop\fresh-shop
mysql -u root -p --default-character-set=utf8mb4 fresh-shop < .\sql\fresh-shop.sql
```

或用 Navicat / DBeaver / MySQL Workbench 打开 `sql\fresh-shop.sql` 导入到 `fresh-shop` 库。

导入成功后应能看到 `shop_goods`、`sys_users`、`shop_order` 等表。

> SQL 体积较大且含示例商品数据，导入可能需一两分钟，属正常。

---

## 4. 修改配置

编辑：

`old_shop/fresh-shop/server/config.yaml`

### 4.1 必改：MySQL（两处保持一致）

文件里有 **`mysql:`** 和 **`db-list:`** 两段，建议都改成你本机账号：

```yaml
system:
  addr: 48888          # 服务端口，默认即可
  db-type: mysql
  use-redis: false     # 本地可先关 Redis
  router-prefix: ""    # 本地保持空

mysql:
  path: localhost
  port: 3306
  config: charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai
  db-name: fresh-shop
  username: root
  password: 你的密码     # ← 改这里

db-list:
  - disable: false
    type: "mysql"
    alias-name: "freshShopMysql"
    path: localhost
    port: 3306
    db-name: "fresh-shop"
    username: "root"
    password: "你的密码"   # ← 同步改这里
```

当前仓库里示例密码可能是 `12345678`，以你本机 MySQL 为准。

### 4.2 微信（本地可先空着）

```yaml
wechat:
  appid: ''
  secret: ''

wechatPay:
  mchId: ''
  # ...
```

- **只测管理端登录、商品列表等**：可不配微信。  
- **测小程序登录 / 微信支付**：再填 AppID、Secret、商户号等；支付回调 `notifyUrl` 需公网可访问地址，纯本机调试通常要用内网穿透。

### 4.3 Redis

`system.use-redis: false` 时可不装 Redis。若改为 `true`，需本机起 Redis，并配置：

```yaml
redis:
  addr: 127.0.0.1:6379
  password: ""
```

---

## 5. 启动服务端

```powershell
cd d:\tiny_zimeiti\cursor_dev_pro\shop\old_shop\fresh-shop\server
go mod tidy
go run main.go
```

首次 `go mod tidy` / 下载依赖可能较慢，可配置代理（可选）：

```powershell
go env -w GOPROXY=https://goproxy.cn,direct
```

启动成功后，控制台一般会有路由注册、监听端口日志。默认监听：

| 项 | 地址 |
|----|------|
| API | http://localhost:48888 |
| 健康检查 | http://localhost:48888/health |
| Swagger | http://localhost:48888/swagger/index.html |

浏览器打开 `/health`，正常应返回 `"ok"`（或等价成功响应）。

编译成可执行文件再跑（可选）：

```powershell
go build -o server.exe main.go
.\server.exe
```

---

## 6. 验证是否跑通

1. 打开 http://localhost:48888/health  
2. 打开 http://localhost:48888/swagger/index.html 能看到接口文档  
3. （可选）公开商品列表：浏览器或 curl 访问类似  
   `http://localhost:48888/goods/getGoodsList?...`  
   （具体查询参数以 Swagger / 后续商品模块文档为准）

管理端默认账号（等你启动 web 时用，**不是**起 server 的必需步骤）：

- 用户名：`admin`  
- 密码：`123456`  
（以 SQL 种子数据为准；若改过库则按实际。）

---

## 7. 常见问题（Win10）

| 现象 | 处理 |
|------|------|
| `connect: connection refused` / 连不上库 | 确认 MySQL 服务已启动；`path/port/username/password/db-name` 是否正确 |
| `Unknown database 'fresh-shop'` | 先建库再导入 SQL |
| 端口被占用 | 改 `system.addr`，或结束占用 48888 的进程 |
| `go: go.mod requires go >= 1.23` | 升级本机 Go 到 1.23+（Go 1.18 不行） |
| 导入 SQL 报 `datetime(3)` / 语法错误 | 当前多半是 MySQL 5.5，需换 **5.7+/8.0** |
| 依赖下载失败 | 设置 `GOPROXY=https://goproxy.cn,direct` 后重试 `go mod tidy` |
| 中文乱码 / 时间不对 | 保持 `charset=utf8mb4` 与 `loc=Asia%2FShanghai` |
| 跨域（浏览器直接调 API） | 后端 `router.go` 里 CORS 中间件默认注释；管理端开发用 Vite 代理 `/goapi`，一般不直连测跨域 |

---

## 8. 和前端的关系（备忘）

| 端 | 如何指到本机 API |
|----|------------------|
| 管理端 `web` | `.env.development` 中 `VITE_SERVER_PORT=48888`，`VITE_BASE_API=/goapi`（Vite 代理到后端） |
| 小程序 | `fresh-shop-uniapp/config/config.js` 的 `baseUrl` 设为 `http://localhost:48888` 或本机局域网 IP |

服务端单独跑通后，再启 web / 小程序即可联调。

---

## 9. 检查清单

- [ ] Go 1.23+、MySQL 已安装  
- [ ] 已建库并导入 `fresh-shop.sql`  
- [ ] `config.yaml` 中 `mysql` 与 `db-list` 密码正确，`use-redis: false`（或不配 Redis）  
- [ ] `go run main.go` 无报错  
- [ ] `/health` 与 Swagger 可访问  

完成后即可继续模块学习，或再启管理端做界面验证。
