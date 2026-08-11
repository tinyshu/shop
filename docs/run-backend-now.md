# 老项目后端启动步骤（Go 1.23 + MySQL 8.0）

按顺序做完即可在本机跑起 `old_shop` 的 Gin 服务。  
源码只读，**不要改业务代码**；只需改本机数据库密码等配置。

---

## 启动前检查（新开一个 CMD）

安装/改 PATH 后请**关闭旧窗口，新开 CMD**，再执行：

```bat
go version
mysql --version
```

期望类似：

```text
go version go1.23.x windows/amd64
mysql  Ver 8.0.40 ...
```

若 `go` 找不到：把 Go 的 `bin` 加到系统/用户 PATH（常见 `C:\Program Files\Go\bin`），再重开 CMD。  
若 `mysql` 找不到：确认已加  
`C:\Program Files\MySQL\MySQL Server 8.0\bin`。

确认 MySQL 服务在跑（服务名一般是 `MySQL80`）：

```bat
sc query MySQL80
```

状态应为 `RUNNING`。

---

## 步骤 1：创建数据库

CMD 登录（密码为安装 MySQL 时设置的 root 密码）：

```bat
mysql -u root -p
```

在 MySQL 里执行：

```sql
CREATE DATABASE IF NOT EXISTS `fresh-shop`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;
EXIT;
```

---

## 步骤 2：导入 SQL

在 CMD 中（路径按你机器，密码交互输入）：

```bat
cd /d d:\tiny_zimeiti\cursor_dev_pro\shop\old_shop\fresh-shop
mysql -u root -p --default-character-set=utf8mb4 fresh-shop < sql\fresh-shop.sql
```

- 文件较大，可能要一两分钟。  
- 也可用 Navicat：选中库 `fresh-shop` → 运行 SQL 文件 → 选 `sql\fresh-shop.sql`。

导入后可用下面命令粗查：

```bat
mysql -u root -p -e "USE `fresh-shop`; SHOW TABLES LIKE 'shop_goods'; SELECT COUNT(*) FROM shop_goods;"
```

能看到 `shop_goods` 且有数据即成功。

> Navicat 若报 **1251 认证协议**：用 8.0 的 mysql 客户端执行一次（密码换成你的）：  
> `ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '你的密码'; FLUSH PRIVILEGES;`
>
> Navicat 若报 **1055 / only_full_group_by / PROFILING**：多半是客户端内部查询触发，**建库/导库往往其实已成功**。可先刷新左侧看有没有库 `fresh-shop`。本机开发可执行：  
> `SET GLOBAL sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));`  
> 或改用 CMD 导入 SQL，避免被 Navicat 干扰。

---

## 步骤 3：改后端配置

用编辑器打开：

`d:\tiny_zimeiti\cursor_dev_pro\shop\old_shop\fresh-shop\server\config.yaml`

### 3.1 确认系统段（一般不用改）

```yaml
system:
  addr: 48888
  db-type: mysql
  use-redis: false
  router-prefix: ""
```

### 3.2 改 MySQL 密码（两处都要改）

把 **`mysql:`** 和 **`db-list:`** 里的 `password` 改成你本机 root 密码（当前文件里示例是 `12345678`，若不一致必须改）：

```yaml
mysql:
  path: localhost
  port: 3306
  db-name: fresh-shop
  username: root
  password: 你的MySQL密码

db-list:
  - disable: false
    type: "mysql"
    alias-name: "freshShopMysql"
    path: localhost
    port: 3306
    db-name: "fresh-shop"
    username: "root"
    password: "你的MySQL密码"
```

微信相关可先留空，不影响先起服务、测接口。

保存文件。

---

## 步骤 4：安装依赖并启动

```bat
cd /d d:\tiny_zimeiti\cursor_dev_pro\shop\old_shop\fresh-shop\server

go env -w GOPROXY=https://goproxy.cn,direct

go mod tidy

go run main.go
```

- 首次 `go mod tidy` / 下载依赖可能较慢，属正常。  
- 控制台无 panic、出现路由注册成功、监听端口等信息即启动成功。  
- **不要关这个窗口**，关了服务就停了。

可选：编译后再跑

```bat
go build -o server.exe main.go
server.exe
```

---

## 步骤 5：验收

浏览器打开：

| 地址 | 期望 |
|------|------|
| http://localhost:48888/health | 返回 `ok` 一类成功响应 |
| http://localhost:48888/swagger/index.html | 能打开 Swagger 文档页 |

管理端账号（等以后起 web 时用，**不是**起后端必需）：`admin` / `123456`（以 SQL 种子为准）。

---

## 常见报错

| 现象 | 处理 |
|------|------|
| `Access denied for user 'root'` | `config.yaml` 两处密码与 MySQL root 不一致 |
| `Unknown database 'fresh-shop'` | 未建库或库名写错，回到步骤 1 |
| `connect: connection refused` | MySQL80 服务未启动 |
| 端口占用 / bind 失败 | 改 `system.addr`，或结束占用 48888 的进程 |
| `go` 不是内部或外部命令 | PATH 未生效，重开 CMD；确认 Go bin 在 PATH |
| 依赖下载失败 | 确认执行了 `go env -w GOPROXY=https://goproxy.cn,direct` |

---

## 一页清单

- [ ] `go version` ≥ 1.23，`mysql --version` 为 8.0.x  
- [ ] MySQL80 服务 Running  
- [ ] 已建库 `fresh-shop` 并导入 `sql\fresh-shop.sql`  
- [ ] `config.yaml` 中 `mysql` + `db-list` 密码正确  
- [ ] `go mod tidy` + `go run main.go` 无报错  
- [ ] `/health`、Swagger 可访问  

完成后后端即就绪；需要时再启管理端 `web` 或小程序联调。
