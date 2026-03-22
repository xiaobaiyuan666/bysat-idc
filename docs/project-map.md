# 无穷云 IDC 系统项目结构图

## 项目定位

当前仓库只保留一条主线工程：

- 系统名称：无穷云 IDC 系统
- active workspace：[`idc-finance`](/Users/a1/Documents/codex/bysat-idc/idc-finance)
- 技术栈：`Go + Vue 3 + TypeScript + Element Plus`

## 根目录职责

- [`README.md`](/Users/a1/Documents/codex/bysat-idc/README.md)：项目入口和快速恢复说明
- [`AGENTS.md`](/Users/a1/Documents/codex/bysat-idc/AGENTS.md)：AI 协作规则和 continuity 约束
- [`docs`](/Users/a1/Documents/codex/bysat-idc/docs)：根目录级别的对接与续开发文档
- [`idc-finance`](/Users/a1/Documents/codex/bysat-idc/idc-finance)：唯一业务代码工作区
- [`package.json`](/Users/a1/Documents/codex/bysat-idc/package.json)：根目录统一启动脚本

## `idc-finance` 结构职责

### 运行入口

- `apps/api`：Go API 服务入口
- `apps/admin-web`：后台管理端
- `apps/portal-web`：用户端 / 门户

### 后端主体

- `internal/modules/customer`：客户、联系人、客户详情相关业务
- `internal/modules/catalog`：商品与价格矩阵相关业务
- `internal/modules/order`：订单、账单、支付、服务主链
- `internal/modules/provider`：Provider、魔方云、同步、远程实例动作
- `internal/modules/portal`：门户侧聚合数据与相关流程
- `internal/modules/report`：报表与工作台聚合逻辑
- `internal/platform`：认证、审计、配置、服务装配、中间件等基础设施

### 数据与脚本

- `migrations`：数据库迁移
- `seed`：演示或初始化数据
- `scripts`：开发和数据库相关脚本

## 前端主要目录

### 后台 `apps/admin-web/src`

- `views/customer`：客户工作台
- `views/catalog`：商品工作台
- `views/order`：订单工作台
- `views/billing`：账单工作台
- `views/service`：服务工作台
- `views/provider`：Provider 与魔方云相关页面
- `views/vnc`：实例控制台页面
- `api`：后台 API 请求封装
- `router/modules`：模块化路由
- `store/modules`：登录态、菜单、语言等前端状态

### 用户端 `apps/portal-web/src`

- `views/console`：用户控制台首页
- `views/store`：商品购买入口
- `views/services`：用户服务列表与服务页面
- `views/orders`：用户订单
- `views/invoices`：用户账单
- `views/tickets`：用户工单
- `views/wallet`：钱包
- `views/account`：账户资料

## 当前保留的关键文档

- [`docs/project-continuity.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-continuity.md)：如何恢复开发
- [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)：当前阶段状态
- [`docs/active-handoff.md`](/Users/a1/Documents/codex/bysat-idc/docs/active-handoff.md)：当前交接面板与下一步建议
- [`docs/cross-device-collaboration.md`](/Users/a1/Documents/codex/bysat-idc/docs/cross-device-collaboration.md)：跨端协作制度
- [`docs/handoff-template.md`](/Users/a1/Documents/codex/bysat-idc/docs/handoff-template.md)：阶段性交接模板
- [`docs/environment-setup.md`](/Users/a1/Documents/codex/bysat-idc/docs/environment-setup.md)：新设备环境恢复与启动说明
- [`docs/development-rules.md`](/Users/a1/Documents/codex/bysat-idc/docs/development-rules.md)：注释和同步规则
- [`docs/mofang-integration.md`](/Users/a1/Documents/codex/bysat-idc/docs/mofang-integration.md)：魔方云接口摘录
- [`docs/mofang-portal-sync.md`](/Users/a1/Documents/codex/bysat-idc/docs/mofang-portal-sync.md)：同步到门户的业务说明

## 当前明确废弃的内容

以下内容已经从仓库主线清理，不再作为继续开发依据：

- 旧根目录 Next 应用
- `admin-ui`
- 旧参考仓和 legacy 源码镜像
- 历史分析文档中仅用于拆解旧系统的材料

如果未来需要新增参考材料，不要直接把大体量参考仓重新塞回主仓库。
应该先确认它是否真的参与当前开发，再决定是否入库。
