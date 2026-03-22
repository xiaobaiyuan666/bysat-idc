# 公有云业务系统架构说明

## 目标定位

这套系统不是单一的 IDC 财务后台，而是按“运营后台 + 客户门户 + 云资源编排 + 财务结算”组织的完整业务平台：

- 运营后台负责产品、客户、订单、账单、支付、工单、云架构和资源运营。
- 客户门户负责下单、续费、支付、工单协同和云资源自助管理。
- 服务与资源动作层负责把业务动作映射到云平台或本地模拟交付。
- 计费引擎负责续费账单、宽限期、停机和自动续费逻辑。
- 通知与任务中心负责消息模板、发送队列和异步任务执行。

## 系统分层

### 1. 展示层

- 管理后台：Vue 3 + Element Plus，运行在 `admin-ui/`
- 客户门户：Next.js App Router，运行在 `src/app/portal/`

### 2. 业务层

- 服务动作层：统一执行实例开关机、重装、改密、VNC、救援、锁定、续费等动作
- 资源动作层：统一执行快照、备份、安全组规则等资源操作
- 财务结算层：统一处理账单支付、余额扣减、支付回调和续费账单创建
- 通知与队列层：统一处理通知模板、通知消息和异步投递任务

### 3. 数据层

- Prisma + SQLite
- 当前模型已覆盖：
  - 客户与门户账号
  - 产品与售卖方案
  - 订单、订单项、账单、支付、余额流水
  - 退款记录、通知模板、通知消息、异步任务
  - 服务实例
  - VPC、IP、磁盘、快照、备份、安全组
  - 工单、审计日志、平台回执、计费作业

## 核心模块

### 云架构中心

统一维护公有云售卖和交付的底层对象：

- 地域
- 可用区
- 规格
- 镜像
- 售卖方案
- 门户账号

售卖方案不只是价格，还包含云平台交付参数：

- 节点标识 `node`
- 系统标识 `os`
- CPU / 内存
- 系统盘 / 数据盘
- 网络类型
- 出入带宽
- 流量限制
- IP 数量
- 防御峰值

关键入口：

- [architecture route](/C:/Users/Administrator/Desktop/IDC/src/app/api/architecture/route.ts)
- [ArchitecturePage.vue](/C:/Users/Administrator/Desktop/IDC/admin-ui/src/pages/ArchitecturePage.vue)
- [ArchitecturePlansSection.vue](/C:/Users/Administrator/Desktop/IDC/admin-ui/src/components/architecture/ArchitecturePlansSection.vue)
- [cloud-plan-config.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/cloud-plan-config.ts)

### 订单与财务

当前订单链路已支持：

- 门户下单
- 自动生成订单、订单项、服务实例、账单
- 余额支付
- 手工收款
- 退款回冲
- 支付回调日志
- 余额流水
- 续费账单生成

关键入口：

- [portal.ts](/C:/Users/Administrator/Desktop/IDC/src/actions/portal.ts)
- [payment-service.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/payment-service.ts)
- [billing-engine.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/billing-engine.ts)

### 通知与任务中心

当前系统已补齐：

- 通知模板
- 手工通知创建
- 订单、支付、退款、工单、服务操作通知入队
- 通知任务队列执行
- 客户门户通知中心

关键入口：

- [notification-service.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/notification-service.ts)
- [notifications route](/C:/Users/Administrator/Desktop/IDC/src/app/api/notifications/route.ts)
- [notifications process route](/C:/Users/Administrator/Desktop/IDC/src/app/api/notifications/process/route.ts)
- [NotificationsPage.vue](/C:/Users/Administrator/Desktop/IDC/admin-ui/src/pages/NotificationsPage.vue)
- [portal notifications page](/C:/Users/Administrator/Desktop/IDC/src/app/portal/(console)/notifications/page.tsx)

### 服务管理

服务实例是业务中枢，关联：

- 客户
- 产品
- 售卖方案
- 订单
- 账单
- 工单
- 资源对象

后台与门户已共用统一的服务动作层：

- 同步状态
- 开机 / 关机 / 重启 / 硬重启 / 硬关机
- 重装系统
- 重置密码
- 获取 VNC
- 进入 / 退出救援模式
- 锁定 / 解锁
- 续费
- 暂停 / 终止

关键入口：

- [service-operations.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/service-operations.ts)
- [services action route](/C:/Users/Administrator/Desktop/IDC/src/app/api/services/[id]/action/route.ts)
- [Portal service detail](/C:/Users/Administrator/Desktop/IDC/src/app/portal/(console)/services/[id]/page.tsx)

### 资源中心

当前已建模并打通：

- VPC
- IP 地址
- 云硬盘
- 快照
- 备份
- 安全组

后台与门户已共用统一的资源动作层：

- 磁盘创建快照
- 磁盘创建备份
- 磁盘挂载 / 卸载
- 设为系统盘
- 快照恢复 / 删除
- 备份恢复 / 归档 / 删除
- 安全组新增规则 / 删除规则 / 删除安全组

关键入口：

- [resource-operations.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/resource-operations.ts)
- [resources action route](/C:/Users/Administrator/Desktop/IDC/src/app/api/resources/[resourceType]/[id]/action/route.ts)
- [ResourcesPage.vue](/C:/Users/Administrator/Desktop/IDC/admin-ui/src/pages/ResourcesPage.vue)

### 客户门户

当前门户已形成完整业务闭环：

- 登录
- 仪表盘
- 云产品目录
- 下单
- 订单中心
- 账单中心
- 钱包与余额流水
- 工单中心
- 实例池
- 单实例控制台

单实例控制台当前可见并可操作：

- 实例基础配置
- 规格、镜像、带宽、防护
- VPC 与 IP
- 磁盘、快照、备份
- 安全组及规则
- 关联账单
- 关联工单
- 平台回执
- 审计轨迹
- 续费账单生成

关键入口：

- [portal-shell.tsx](/C:/Users/Administrator/Desktop/IDC/src/components/portal-shell.tsx)
- [portal-data.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/portal-data.ts)
- [portal services list](/C:/Users/Administrator/Desktop/IDC/src/app/portal/(console)/services/page.tsx)
- [portal service detail](/C:/Users/Administrator/Desktop/IDC/src/app/portal/(console)/services/[id]/page.tsx)

## 关键业务流程

### 商品上架流程

1. 运营后台维护地域、可用区、规格和镜像。
2. 配置售卖方案与云平台交付参数。
3. 门户商品页读取公开方案并展示。

### 新购开通流程

1. 客户在门户下单。
2. 系统创建订单、订单项、服务实例和账单。
3. 账单支付完成后进入服务生命周期。
4. 服务动作层继续对接云平台开通或人工交付。

### 续费流程

1. 计费引擎或门户为实例生成续费账单。
2. 客户在门户支付续费账单。
3. 服务续期并重新纳入计费计划。

### 运维流程

1. 客户在实例控制台发起实例或资源操作。
2. 门户动作先校验资源归属关系。
3. 服务动作层 / 资源动作层统一执行业务动作。
4. 平台回执写入 `ProviderSyncLog`。
5. 业务审计写入 `AuditLog`。
6. 页面回显操作结果、任务号和控制台地址。

## 关键代码入口

- [schema.prisma](/C:/Users/Administrator/Desktop/IDC/prisma/schema.prisma)
- [portal.ts](/C:/Users/Administrator/Desktop/IDC/src/actions/portal.ts)
- [service-operations.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/service-operations.ts)
- [resource-operations.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/resource-operations.ts)
- [billing-engine.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/billing-engine.ts)
- [payment-service.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/payment-service.ts)
- [notification-service.ts](/C:/Users/Administrator/Desktop/IDC/src/lib/notification-service.ts)

## 当前边界

这套系统现在已经是“完整商用骨架 + 主流程可跑通”，但距离正式生产交付仍有明确边界：

- 魔方云磁盘、快照、备份、安全组的真实生产接口映射还需继续补齐
- 真实支付渠道签名、渠道原路退款、红冲、发票和税务流程还未完整落地
- 通知通道正式对接、失败重试和更复杂的异步编排还需要进一步工程化
- 部署、安全加固、监控告警和自动化备份还需要补生产方案

## 建议的下一阶段

- 补真实魔方云资源级 API 映射
- 补支付渠道抽象、退款和发票链路
- 增加消息通知中心与任务队列
- 增加报表导出、对账和财务稽核能力
- 完善门户多角色权限与操作授权
