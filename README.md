# 无穷云 IDC 系统

当前仓库已经整理为单一主线：[`idc-finance`](/Users/a1/Documents/codex/bysat-idc/idc-finance)。

这套系统现在统一命名为“无穷云 IDC 系统”，当前代码基线为 AI 生成并在此仓库内持续演进。

## 当前入口

- 后台：`http://localhost:5177`
- 用户端：`http://localhost:5178`
- API：`http://localhost:18080/api/v1/health`

## 默认账号

- 后台：`admin / Admin123!`
- 用户端：`portal / Portal123!`

## 启动

在仓库根目录可以直接执行：

```bash
npm run dev:api
npm run dev:admin
npm run dev:portal
```

也可以直接进入 [`idc-finance`](/Users/a1/Documents/codex/bysat-idc/idc-finance) 目录启动。

## 说明

- 根目录只保留当前开发还会继续用到的入口和说明文件
- 当前保留的补充文档只剩魔方云接口与门户同步两份说明，位于 [`docs`](/Users/a1/Documents/codex/bysat-idc/docs)
- 服务详情里的 `VNC` 已迁移到新后台，入口在新版后台的服务详情页 `远端实例 -> 获取控制台`

## 续开发入口

如果你换了编辑器、机器，或者中断一段时间后要继续开发，先按这个顺序阅读：

1. [`docs/project-continuity.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-continuity.md)
2. [`docs/project-map.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-map.md)
3. [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)
4. [`docs/development-rules.md`](/Users/a1/Documents/codex/bysat-idc/docs/development-rules.md)
5. [`idc-finance/README.md`](/Users/a1/Documents/codex/bysat-idc/idc-finance/README.md)
