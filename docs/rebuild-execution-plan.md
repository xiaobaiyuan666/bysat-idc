# IDC 云业务系统重建执行计划

## 1. 目标

基于已完成的逆向拆解结果，重建一套接近原系统业务结构的 IDC 财务与云业务管理系统。

本计划只解决三件事：

1. 明确重建范围与边界
2. 明确开发顺序与并行拆分
3. 明确每阶段交付物与验收标准

## 2. 重建原则

- 架构形态采用模块化单体，不先拆微服务
- 先还原业务结构，再优化技术实现
- 优先还原订单、账单、支付、服务生命周期主链
- 后台采用列表工作台 + 详情工作台
- 用户中心采用控制台首页 + 业务中心 + 单服务工作台
- 所有关键写操作必须可审计
- 所有异步动作必须有任务记录、失败重试和最终状态
- Provider 统一走适配器，不让业务层直接耦合魔方云或其他上游

## 3. 技术路线

- 后端：Go + Gin + GORM + Casbin + Redis + MySQL
- 异步与计划任务：Asynq + cron
- 后台前端：Vue 3 + TypeScript + Element Plus + vue-pure-admin
- 用户中心：Vue 3 + TypeScript + Element Plus + vue-pure-admin Portal Layout
- 接口协议：RESTful 为主，Webhook 作为异步回调补充

## 4. 总体阶段

### Phase 0：基础底座

- 初始化 API 服务
- 初始化后台前端骨架
- 初始化用户中心骨架
- 建立统一配置、日志、错误码、审计、上传、任务机制

交付物：

- 可运行的 API 服务
- 后台登录页与空白工作台
- 用户中心登录页与空白控制台
- 数据库初始化迁移

验收：

- 本地一键启动
- 后台与用户中心能正常登录
- 审计日志、菜单、角色、权限基础表可用

### Phase 1：权限与客户域

- 后台管理员认证
- RBAC
- 菜单树
- 数据权限骨架
- 客户、联系人、实名认证、客户分组、客户等级
- 客户详情工作台

交付物：

- 角色权限系统
- 客户模块全量页面
- 客户详情工作台

验收：

- 不同角色进入后台看到的菜单不同
- 客户 CRUD、联系人 CRUD、实名信息录入可用
- 客户详情页可聚合资料、服务、账单、工单入口

### Phase 2：商品与定价域

- 产品组、产品、价格、周期
- 配置项组、配置项、配置项值
- 产品自定义字段
- 升级映射
- 优惠码

交付物：

- 商品管理工作台
- 配置项管理工作台
- 价格计算服务

验收：

- 产品可配置周期与价格
- 配置项可影响价格
- 优惠码可通过规则校验

### Phase 3：购物车、订单与账单主链

- 购物车
- 下单结算
- 订单生成
- 账单生成
- 账单详情
- 后台订单审核

交付物：

- 前台购物车与结算页
- 后台订单列表与订单详情工作台
- 后台账单列表与账单详情工作台

验收：

- 用户可从商品下单到生成订单与账单
- 后台可审核订单
- 订单、账单、账单项关系正确

### Phase 4：支付、流水、余额、信用额

- 支付网关抽象
- 支付回调
- 余额支付
- 信用额支付
- 线下补录
- 流水与退款
- 合并账单

交付物：

- 支付中心
- 流水中心
- 退款中心
- 对账中心

验收：

- 账单可完成支付
- 支付后会写流水和审计
- 余额与信用额支付可用
- 可发起退款并正确回冲

### Phase 5：服务与资源域

- 服务实例模型
- 服务配置快照
- 服务详情工作台
- 续费、升降级、暂停、解除暂停、终止
- 服务器、服务器组、资源池

交付物：

- 服务列表页
- 服务详情工作台
- 节点与资源配置页

验收：

- 支付成功后能生成服务实例
- 服务详情能展示账期、状态、配置、日志
- 服务动作状态流转正确

### Phase 6：Provider 与自动化开通

- 定义 ProviderAdapter
- 实现 Mock Provider
- 实现魔方云 Provider
- 实现上游 API Provider
- 接入自动开通、自动续费、到期暂停、到期删除
- 接入 RunMap 与同步任务

交付物：

- Provider 抽象层
- 自动化任务中心
- 魔方云同步与动作链

验收：

- 下单支付后能自动开通
- 定时任务能生成续费账单、暂停过期服务
- Provider 动作有日志、有重试、有最终状态

### Phase 7：工单与售后

- 工单、回复、备注、部门、状态、预设回复
- 用户提单、回复、关闭、评价
- 后台接单、派单、转派、内部备注

交付物：

- 工单列表
- 工单详情工作台
- 部门与状态配置页

验收：

- 用户可提交并追踪工单
- 后台可处理、转派、关闭
- 工单与服务、客户关联正确

### Phase 8：分销、代理、合同、发票

- 分销链接与关系绑定
- 佣金、提现、审核
- 合同模板、合同签署、邮寄
- 发票抬头、地址、发票申请

交付物：

- 分销中心
- 合同中心
- 发票中心

验收：

- 推广用户可绑定关系并计算佣金
- 提现申请与审核可用
- 合同与发票流程可跑通

### Phase 9：报表、OpenAPI、上线准备

- 收入报表
- 客户报表
- 服务分布与风险报表
- OpenAPI
- Webhook
- 部署与监控

交付物：

- 报表中心
- OpenAPI 文档
- 部署文档

验收：

- 主要报表可查
- OpenAPI 可鉴权调用
- 本地、测试、生产部署方案清晰

## 5. 并行线程拆分

多线程对这个项目有帮助，但前提是每条线程都有清晰边界，不交叉写同一组文件。

### 线程 A：平台与权限

负责：

- 管理员认证
- 角色权限
- 菜单
- 审计
- 系统配置

写入范围：

- `internal/platform/*`
- `apps/admin-web/src/router`
- `apps/admin-web/src/store/auth*`

### 线程 B：客户与工单

负责：

- 客户
- 联系人
- 实名认证
- 客户详情工作台
- 工单系统

写入范围：

- `internal/modules/customer/*`
- `internal/modules/ticket/*`
- `apps/admin-web/src/views/customer/*`
- `apps/admin-web/src/views/ticket/*`
- `apps/portal-web/src/views/customer/*`
- `apps/portal-web/src/views/ticket/*`

### 线程 C：商品、订单、财务

负责：

- 商品
- 定价
- 购物车
- 订单
- 账单
- 支付
- 退款

写入范围：

- `internal/modules/catalog/*`
- `internal/modules/cart/*`
- `internal/modules/order/*`
- `internal/modules/billing/*`
- `internal/modules/payment/*`
- `apps/admin-web/src/views/order/*`
- `apps/admin-web/src/views/billing/*`
- `apps/portal-web/src/views/cart/*`
- `apps/portal-web/src/views/order/*`

### 线程 D：服务、资源、Provider

负责：

- 服务实例
- 续费
- 升降级
- 节点资源
- Provider
- 自动开通
- 同步任务

写入范围：

- `internal/modules/service/*`
- `internal/modules/provider/*`
- `internal/modules/resource/*`
- `apps/admin-web/src/views/service/*`
- `apps/admin-web/src/views/resource/*`
- `apps/portal-web/src/views/service/*`

### 线程 E：运营与开放能力

负责：

- 分销
- 合同
- 发票
- 报表
- OpenAPI

写入范围：

- `internal/modules/affiliate/*`
- `internal/modules/contract/*`
- `internal/modules/report/*`
- `internal/modules/openapi/*`
- `apps/admin-web/src/views/report/*`
- `apps/portal-web/src/views/affiliate/*`

## 6. 推荐执行方式

### 第一波

- 线程 A
- 线程 B
- 线程 C

原因：

- 这三条线能最快把“可登录、可管理客户、可下单、可出账单”打通。

### 第二波

- 线程 D

原因：

- 服务与 Provider 必须建立在订单和账单主链稳定之后。

### 第三波

- 线程 E

原因：

- 分销、合同、报表、开放接口属于上层扩展，不应阻塞核心交付。

## 7. 每波验收标准

### Wave 1 验收

- 后台管理员和前台用户均可登录
- 客户、商品、订单、账单四个主模块可完整 CRUD
- 能从购物车下单到生成账单
- 能在后台查看订单详情和账单详情

### Wave 2 验收

- 支付成功后自动生成服务
- 服务支持暂停、解除暂停、续费、升级
- Provider 动作有任务与日志
- 到期自动暂停与续费账单自动生成生效

### Wave 3 验收

- 工单闭环
- 分销提现闭环
- 报表可查
- OpenAPI 可调用

## 8. 风险与控制

- 风险：先做页面不做对象主链，最后会返工。
  - 控制：所有页面都必须先绑定真实对象和状态机。

- 风险：Provider 层和业务层耦合。
  - 控制：统一 `ProviderAdapter`，业务层禁止直连魔方云接口。

- 风险：支付、退款、余额、信用额之间数据不一致。
  - 控制：财务写操作必须事务化，且全部落流水与审计。

- 风险：多线程并行互相覆盖。
  - 控制：按模块目录切写入范围，主线程统一合并接口约定。

## 9. 下一步动作

下一步不再继续扩分析，直接进入代码生成准备：

1. 初始化新项目目录
2. 建基础迁移与通用包
3. 先做 Phase 0 和 Phase 1
4. 然后按 Wave 1 推进

