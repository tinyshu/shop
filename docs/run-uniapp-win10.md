# Win10 运行客户端（fresh-shop-uniapp）

本项目是 **HBuilderX 工程**（目录下无根级 `package.json`），在 Win10 上用 **HBuilderX + 微信开发者工具** 跑微信小程序。  
前提：后端已在 `http://localhost:48888` 正常运行。

---

## 1. 安装工具

| 工具 | 用途 | 下载 |
|------|------|------|
| **HBuilderX** | 打开/编译 UniApp | https://www.dcloud.io/hbuilderx.html（选 App 开发版更省事） |
| **微信开发者工具** | 小程序模拟器 | https://developers.weixin.qq.com/miniprogram/dev/devtools/download.html |

安装后：

1. 打开微信开发者工具，用微信扫码登录  
2. 设置里打开 **服务端口**（部分版本在「设置 → 安全设置 → 服务端口」），供 HBuilderX 唤起  

---

## 2. 导入项目

1. 打开 HBuilderX  
2. **文件 → 导入 → 从本地目录导入**  
3. 选择目录：

```text
d:\tiny_zimeiti\cursor_dev_pro\shop\old_shop\fresh-shop-uniapp
```

不要只导入子文件夹；根目录应能看到 `pages.json`、`manifest.json`、`config`、`pages`。

---

## 3. 配置 API 地址

编辑：`fresh-shop-uniapp/config/config.js`

本机联调保持：

```js
baseUrl: 'http://localhost:48888',
```

若用**真机预览**，`localhost` 指向手机自己，需改成电脑局域网 IP，例如：

```js
baseUrl: 'http://192.168.x.x:48888',
```

（手机与电脑同一 WiFi；Windows 防火墙放行 48888。）

---

## 4. 微信小程序 AppID（本地调试）

`manifest.json` → `mp-weixin` → `appid` 当前有示例值。本地可选：

| 方式 | 说明 |
|------|------|
| **测试号 / 自己的 AppID** | 在微信公众平台申请，改到 `manifest.json` 的 `mp-weixin.appid` |
| **开发者工具里选「测试号」** | 部分流程可不绑正式 AppID，以工具提示为准 |

后端微信登录另需在 `server/config.yaml` 配置：

```yaml
wechat:
  appid: '你的小程序AppID'
  secret: '你的小程序Secret'
```

**未配 AppID/Secret 时**：页面浏览、调公开接口可能可以；**微信登录 / 支付** 通常不可用。可先看商品列表等公开能力。

---

## 5. 运行到微信开发者工具

1. 确认后端窗口仍在跑（`48888`）  
2. HBuilderX 中选中项目根  
3. 菜单：**运行 → 运行到小程序模拟器 → 微信开发者工具**  
4. 首次会编译，并尝试自动打开微信开发者工具  

若未自动打开：

- 确认微信开发者工具已安装且开了服务端口  
- HBuilderX：**工具 → 设置 → 运行配置**，填微信开发者工具安装路径  
- 或手动打开微信开发者工具，导入编译产物目录（常见在项目下 `unpackage/dist/dev/mp-weixin`）

在微信开发者工具中建议勾选：

- **不校验合法域名、web-view（业务域名）、TLS 版本以及 HTTPS 证书**（本地调试）

`manifest.json` 里已有 `"urlCheck": false`，本地一般可直接请求 `http://localhost:48888`。

---

## 6. 验收是否联通

1. 小程序能打开首页  
2. 开发者工具 **Network** 里请求打到 `localhost:48888`（或你的 IP）  
3. 商品列表等公开接口有数据（库已导入 SQL）  

登录失败时优先查：后端 `wechat.appid/secret`、小程序 AppID 是否一致、后端日志报错。

---

## 7. 「构建」指什么

| 动作 | 怎么做 | 用途 |
|------|--------|------|
| **运行（开发）** | 上面第 5 步 | 日常联调 |
| **发行（正式构建）** | HBuilderX：**发行 → 小程序-微信** | 生成上传包，再在微信开发者工具上传 |

本地学习先会「运行」即可，不必先做发行。

---

## 8. 可选：跑 H5（浏览器）

项目 `manifest.json` 带了 `h5` 配置。HBuilderX：**运行 → 运行到浏览器 → Chrome**。  

注意：微信登录、`uni.requestPayment` 等在浏览器里不完整，**主路径仍以微信开发者工具为准**。

---

## 9. 常见问题

| 现象 | 处理 |
|------|------|
| HBuilderX 唤不起微信工具 | 开服务端口；配置工具安装路径；手动打开并导入 `unpackage/dist/dev/mp-weixin` |
| 请求失败 / 不在合法域名 | 开发者工具勾选「不校验合法域名」；确认后端已启动 |
| 真机看不到数据 | `baseUrl` 改成电脑局域网 IP，不要用 `localhost` |
| 登录失败 | 配置后端 `wechat.appid/secret`，与小程序 AppID 一致 |
| 端口被占用 / 编不过 | 看 HBuilderX 控制台报错；清理后重新运行 |

---

## 一页清单

- [ ] 已装 HBuilderX、微信开发者工具  
- [ ] 已导入 `old_shop/fresh-shop-uniapp`  
- [ ] `config/config.js` 的 `baseUrl` 指向本机后端  
- [ ] 后端 `48888` 在跑  
- [ ] 运行到微信开发者工具，首页能出、接口有响应  

管理端 + 后端 + 小程序都通后，可继续模块学习（商品 / 购物车 / 订单）。
