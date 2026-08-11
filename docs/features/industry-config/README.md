# industry-config — 行业开关 / 功能配置

> 用 `sys_config` 控制冻品/小 B 能力，默认走标准 B2C。对应阶段二 **M0**。  
> 总序见 [phase2-module-roadmap.md](../../phase2-module-roadmap.md)。

## 管什么

- 功能开关键约定（审核、月结、物流模式等）
- Server 侧统一读取与默认值（缺键不崩）
- 后续切片：审核可关（GEN-03）、管理端编辑、物流模式与履约联动

## 版本

| 版本 | 端 | 说明 | 状态 |
|------|-----|------|------|
| [v0.1.0](./v0.1.0/) | server | 配置键约定 + 读取封装 | **已实现** |
| 后续 | server / web / uniapp | M0-2 GEN-03、M0-3 管理端 | 未开始 |

## 路线图摘要

见同模块 [roadmap.md](./roadmap.md)。
