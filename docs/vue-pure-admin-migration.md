# vue-pure-admin 接入研究与迁移方案

> 说明：本文档保留为历史迁移记录。当前默认后台入口已切换到 `idc-finance/apps/admin-web`，`admin-ui/` 不再参与默认开发流程。

## 本地基线

- 官方仓库已拉取到 [vue-pure-admin-ref](C:/Users/Administrator/Desktop/IDC/vue-pure-admin-ref)
- 当前参考版本：`6.3.0`
- 当前机器环境：
  - `node -v` 为 `v24.14.0`
  - 未安装 `pnpm`
  - `corepack pnpm` 访问 npm registry 超时，因此当前不走整仓直接安装路线

## 我确认下来的框架核心

### 1. 路由必须模块化
- 官方把业务路由拆在 `src/router/modules`
- 主路由入口只负责组装模块、登录守卫、标签页和权限挂载
- 这套做法适合 IDC 后台，不适合继续维持单文件大路由

参考：
- [router/index.ts](C:/Users/Administrator/Desktop/IDC/vue-pure-admin-ref/src/router/index.ts)
- [board.ts](C:/Users/Administrator/Desktop/IDC/vue-pure-admin-ref/src/router/modules/board.ts)

### 2. 布局不是单个侧栏，而是一整套后台壳层
- 左侧导航
- 顶部导航
- 标签页
- 页面容器
- 权限驱动菜单

参考：
- [layout/index.vue](C:/Users/Administrator/Desktop/IDC/vue-pure-admin-ref/src/layout/index.vue)

### 3. 权限体系分为三层
- 用户 store：登录态、用户信息、权限
- permission store：菜单树、可见路由
- 路由守卫：登录和权限校验

参考：
- [user.ts](C:/Users/Administrator/Desktop/IDC/vue-pure-admin-ref/src/store/modules/user.ts)
- [permission.ts](C:/Users/Administrator/Desktop/IDC/vue-pure-admin-ref/src/store/modules/permission.ts)

## 当前迁移结论

- 后台不再做“接近 pure-admin”，而是直接按它的结构重建
- 用户中心也要进入同一套 Vue 基座，但布局不能照搬后台，需要做门户布局
- 后端 API 仍保留现有 Next.js 服务层，不动后端业务域

## 已完成的后台替换

当前 `admin-ui` 已经切换到 `vue-pure-admin` 风格的核心骨架：
- [router/index.ts](C:/Users/Administrator/Desktop/IDC/admin-ui/src/router/index.ts)
- [workbench.ts](C:/Users/Administrator/Desktop/IDC/admin-ui/src/router/modules/workbench.ts)
- [business.ts](C:/Users/Administrator/Desktop/IDC/admin-ui/src/router/modules/business.ts)
- [finance.ts](C:/Users/Administrator/Desktop/IDC/admin-ui/src/router/modules/finance.ts)
- [cloud.ts](C:/Users/Administrator/Desktop/IDC/admin-ui/src/router/modules/cloud.ts)
- [platform.ts](C:/Users/Administrator/Desktop/IDC/admin-ui/src/router/modules/platform.ts)
- [layout/index.vue](C:/Users/Administrator/Desktop/IDC/admin-ui/src/layout/index.vue)
- [NavVertical.vue](C:/Users/Administrator/Desktop/IDC/admin-ui/src/layout/components/lay-sidebar/NavVertical.vue)
- [lay-navbar/index.vue](C:/Users/Administrator/Desktop/IDC/admin-ui/src/layout/components/lay-navbar/index.vue)
- [lay-tag/index.vue](C:/Users/Administrator/Desktop/IDC/admin-ui/src/layout/components/lay-tag/index.vue)
- [permission.ts](C:/Users/Administrator/Desktop/IDC/admin-ui/src/stores/modules/permission.ts)
- [multiTags.ts](C:/Users/Administrator/Desktop/IDC/admin-ui/src/stores/modules/multiTags.ts)
- [user.ts](C:/Users/Administrator/Desktop/IDC/admin-ui/src/stores/modules/user.ts)

这一步的结果是：
- 根布局已经由新的 pure-admin 壳层接管
- 菜单按领域分组
- 顶栏和标签页开始接管页面切换
- 登录页和后台首页已经按新基座实跑

## 下一阶段

### Phase 1
- 继续把剩余后台页面头部、筛选区、表格容器统一成 pure-admin 页面模板
- 把客户、订单、服务、账单四个核心详情页做成稳定工作台

### Phase 2
- 在同一套 Vue 基座中创建 `portal` 路由域
- 增加门户专用布局，而不是继续依赖当前 Next 页面壳层
- 逐步迁移用户中心的首页、订单、服务、账单、工单、钱包

### Phase 3
- 建立后台与门户共用的权限、标签页、通知、缓存策略
- 完成双端同基座运行

## 现阶段判断

- 后台已经进入“真迁移”阶段，不再是样式参考
- 用户中心还没迁过去，这是下一波主任务
- 这条路比整仓替换稳，因为我们保住了现有业务接口和数据模型
