# IDC 财务与云业务系统重建蓝图（2026-03-21）

本文是基于以下证据收敛出的重建方案：

- 可读源码仓库 `zjmf-manger-decoded-ref`
- 本地运行中的老财务系统 3.7.4
- 真实上游财务接口与真实魔方云接口
- 已经抓下来的后台、用户中心、商品购买页、服务详情页、上下游管理页、数据库表结构、接口回包

目标不是“做一个像 IDC 系统的后台”，而是：

- 用 `Go + Vue 3 + vue-pure-admin`
- 重建一套接近老魔方财务 1:1 结构级别的 IDC 财务与云业务系统
- 同时保留现代实现质量、模块清晰度、可维护性

---

## 1. 结论：现在可以开始改造

可以开始。

不是因为我们已经把所有边角都 100% 看完，而是因为：

1. 主对象链已经完全明确  
   `客户 -> 订单 -> 账单 -> 账单项 -> 服务 -> Provider 动作`

2. 后台三大核心工作台已经看清  
   - 订单工作台
   - 账单工作台
   - 客户工作台里的服务工作台

3. 商品与购买链已经看清  
   - 普通商品
   - `zjmf_api` 上下游商品
   - `dcimcloud` 云商品
   - v10 `mf_cloud` 配置购买页

4. 用户中心服务详情链已经看清  
   - `host/header`
   - `host/details`
   - `snapshot`
   - `setting`
   - `remoteInfo`
   - `VNC`

5. 对接类型边界已经看清  
   - `zjmf_api`
   - `dcimcloud`
   - `manual`
   - `resource`
   - `upper_reaches`

剩余未完全抠透的部分，已经属于：

- 边角弹窗
- 个别管理项
- 一些低频操作
- 某些老版本 bundle 细节

这些不再影响主架构设计。

结论就是：

- 架构定型可以开始
- 改造可以开始
- 细节继续并行补证据

---

## 2. 我们要做的不是“复刻页面”，而是复刻系统骨架

老魔方真正成熟的地方不是皮肤，也不是某几个接口，而是这 6 个层次：

1. 领域对象非常清晰
2. 工作台分工非常清晰
3. 财务链驱动业务链
4. 服务对象同时承载财务、资源、生命周期
5. 上下游和魔方云是标准 provider 体系，不是散乱对接
6. 用户中心不是附属页，而是完整业务控制台

因此新系统的重建原则必须是：

- 先还原对象与工作台
- 再还原动作链
- 最后再还原视觉与交互密度

不能反过来。

---

## 3. 新系统总架构

重建主项目继续以：

- [idc-finance](/C:/Users/Administrator/Desktop/IDC/idc-finance)

为主线。

技术路线保持：

- 后端：Go
- 后台前端：Vue 3 + vue-pure-admin
- 用户中心：Vue 3 + vue-pure-admin portal layout
- 数据库：MySQL
- 异步任务：后续接 Redis + 队列

### 3.1 架构原则

- 架构形态：模块化单体
- 前后端分离
- 领域模块边界固定
- Provider 统一抽象
- 所有写动作可审计
- 所有异步动作可重试

### 3.2 核心模块划分

- 平台层
  - 认证
  - RBAC
  - 菜单
  - 审计
  - 配置中心
  - 文件上传
  - 消息模板
  - 队列与计划任务

- 客户层
  - 客户
  - 联系人 / 子账户
  - 实名认证
  - 客户分组
  - 客户等级
  - 客户跟进
  - 客户附件
  - 客户通知日志

- 商品层
  - 商品组
  - 商品
  - 周期价格矩阵
  - 可配置项组
  - 可配置项
  - 可配置项子项价格
  - 自定义字段
  - 升降级映射
  - 云商品模板

- 订单层
  - 购物车
  - 结算
  - 订单
  - 订单项
  - 订单审核
  - 取消请求

- 财务层
  - 账单
  - 账单项
  - 支付流水
  - 余额流水
  - 信用额
  - 退款
  - 发票
  - 合同
  - 对账

- 服务层
  - 服务实例
  - 服务配置快照
  - 服务动作日志
  - 续费
  - 升降级
  - 暂停 / 解除暂停
  - 删除 / 终止

- Provider 层
  - 通用 provider
  - `zjmf_api`
  - `dcimcloud`
  - `manual`
  - `resource`
  - `upper_reaches`

- 工单层
  - 工单
  - 回复
  - 备注
  - 部门
  - 状态
  - 转派

- 运营层
  - 分销 / 代理
  - 供应商管理
  - 商品同步
  - 资源池
  - 任务队列
  - 报表

- 用户中心层
  - 控制台首页
  - 产品 / 服务
  - 订单
  - 账单
  - 交易记录
  - 钱包 / 信用额
  - 工单
  - API 概览
  - 通知日志
  - 附件 / 跟进 / 合同 / 发票

---

## 4. 必须复刻的对象主链

这是整个系统的主脊梁：

```text
客户
  -> 订单
    -> 账单
      -> 账单项
        -> 服务
          -> Provider 资源
```

这条链的具体含义是：

### 4.1 客户

客户不是只存一张 `customers` 表，而是一个工作台对象：

- 基础资料
- 联系人
- 实名
- 服务
- 账单
- 交易
- 工单
- API
- 通知
- 附件

### 4.2 订单

订单不是支付对象，它是：

- 下单行为对象
- 交付前置对象
- 审核对象
- 与账单和服务之间的中继对象

### 4.3 账单

账单才是财务动作的核心对象：

- 支付
- 手工入账
- 退款
- 编辑账单项
- 发送账单邮件
- 查看账单日志

### 4.4 服务

服务不是“资源实例”这么简单，它同时承载：

- 本地业务身份
- 财务字段
- 订单来源
- 配置快照
- provider 映射
- 生命周期动作

### 4.5 Provider 资源

资源与服务不是同一个层次：

- 服务是本地业务对象
- Provider 资源是外部交付对象

这个层次在老魔方里非常重要。

例如：

- `host.id` 是本地服务 ID
- `host.dcimid` 在 `zjmf_api` 情况下对应的是上游财务 host id
- 真正底层云实例 id 可能是另外一层 provider 资源 id

新系统必须保留这个分层。

---

## 5. 后台 IA 重建方案

### 5.1 一级导航

一级导航不应该再按我们现在的“开发视角模块”来排，而要按老魔方的运营工作流排：

- 工作台
- 客户
- 业务
- 财务
- 工单
- 商品设置
- 资源与商店
- 功能 / 统计
- 设置

### 5.2 二级导航

#### 工作台

- 首页

#### 客户

- 客户列表
- 实名认证
- 客户资源池
- 我的业绩
- 推介计划
- 营销推送

#### 业务

- 产品订单
- 续费订单
- 流量包订单
- 业务列表
- 产品暂停请求

#### 财务

- 交易流水
- 账单管理
- 信用额管理
- 提现审核
- 发票列表
- 合同列表

#### 工单

- 工单列表
- 工单统计

#### 商品设置

- 商品管理
- 流量包管理
- 通用接口
- 魔方 DCIM
- 魔方云
- 全局可配置项
- 商品订购设置

#### 资源与商店

- 供应商管理
- 商品管理
- 产品管理
- 任务队列
- 服务器列表
- 订单列表
- 续费订单
- API 设置
- 应用商店

#### 功能 / 统计

- 年度收入统计
- 新客户
- 产品收入
- 收入排名
- 定时任务状态
- 任务队列
- 数据库状态

#### 设置

- 客户设置
- 财务设置
- 工单设置
- 系统设置
- 站务设置

### 5.3 后台页面模型

后台只保留两种主页面模型：

1. 列表工作台
2. 详情工作台

#### 列表工作台统一结构

- 顶部状态 tabs
- 基础筛选
- 高级筛选
- 表格
- 行内跳转
- 批量动作
- 汇总区
- 分页区

#### 详情工作台统一结构

- 顶部标题 / 状态 / 主动作
- 基础信息卡
- 子对象表格
- 动作按钮区
- 日志区

---

## 6. 后台三大核心工作台重建方案

### 6.1 订单工作台

#### 订单列表

字段参考老系统：

- ID
- 客户名
- 产品
- IP
- 下单时间
- 金额
- 付款状态 / 付款方式
- 状态
- 客户备注
- 提成 / 销售

批量动作：

- 核验通过
- 取消订单
- 删除订单

#### 订单详情

分两块：

1. 订单头
2. 订单项目

订单头字段：

- 客户
- 订单号
- 时间
- 优惠码
- IP
- 账单信息
- 付款方式
- 金额
- 状态
- 客户备注

订单项目字段：

- 条目
- 周期
- 金额
- 服务状态
- 支付状态
- 是否发送开通邮件

动作：

- 核验通过
- 取消订单
- 删除订单
- 保存备注
- 修改状态

### 6.2 账单工作台

#### 账单列表

状态分组：

- 未支付
- 已支付
- 已取消
- 已退款

字段：

- 账单号
- 客户
- 金额
- 已付
- 应付
- 付款方式
- 类型
- 创建时间
- 到期时间
- 状态

#### 账单详情

必须拆成 4 个区：

1. 账单信息卡
2. 交易明细
3. 账单项目
4. 操作日志

动作：

- 编辑账单头
- 发送邮件
- 手工入账
- 余额支付
- 信用额支付
- 退款
- 删除交易
- 编辑账单项
- 删除账单项
- 保存更改

### 6.3 服务工作台

这是最复杂的一页。

#### 入口关系

服务详情不能只做独立 `/services/:id`。

必须保留两种访问方式：

1. 客户详情内进入
2. 业务列表独立进入

#### 页面结构

必须拆成 4 大区：

1. 服务与客户上下文区
2. Provider 动作区
3. 财务字段区
4. 配置快照 / 配置变更区

#### 基础字段

- 订单号
- 商品 / 服务
- 接口
- IP
- 端口
- 用户名
- 密码
- 主机名
- 备注
- 上游产品 ID / 资源 ID

#### 动作区

动作必须按类型拆分：

- DB-only 动作
  - 保存更改
  - 转移服务
  - 删除记录

- 需要确认后执行的 provider 动作
  - 拉取信息
  - 删除
  - 暂停

- 需要参数弹窗的动作
  - 重装系统
  - 重置密码
  - 救援系统

- 需要新页承接的动作
  - VNC

#### 财务区

字段：

- 订购时间
- 首付金额
- 付款方式
- 到期时间
- 续费金额
- 自动续费
- 付款周期
- 状态
- 优惠码
- 成本

动作：

- 交易流水
- 账单
- 续费

#### 配置区

这里不能做成一个 JSON 文本框。

必须按商品配置项逐项渲染：

- 区域
- 操作系统
- CPU
- 内存
- 网络类型
- 带宽
- 流入带宽
- IP 数量
- 数据盘
- 远程端口

并保留：

- 应用至接口

这个动作。

---

## 7. 商品架构重建方案

### 7.1 商品不是一张表

老系统里商品本质是：

`商品定义 + 财务定义 + 交付定义 + 配置定义`

新系统必须拆成这些对象：

- product_groups
- products
- product_pricings
- product_config_groups
- product_config_options
- product_config_option_values
- product_custom_fields
- product_upgrade_maps
- product_provider_templates

### 7.2 商品字段必须分四层

#### 基础层

- 名称
- 标识
- 描述
- 分组
- 状态
- 排序

#### 财务层

- 周期价格
- 开户费
- 续费价
- 折扣
- 优惠码可用性
- 税率

#### 交付层

- `type`
- `api_type`
- provider id
- server group
- upstream product id
- upstream price type
- upstream price value
- upper_reaches_id
- 自动开通策略

#### 配置层

- 可配置项组
- 可配置项
- 子项价格
- 自定义字段
- 云模板默认值

### 7.3 三类商品必须分开建模

#### A. 本地普通商品

本地模块直接开通。

#### B. `zjmf_api` 上下游商品

本质是上游财务商品映射，关键能力：

- 同步商品
- 同步价格
- 同步配置项
- 同步自定义字段
- 远端结算
- 回写上游 host id

#### C. `dcimcloud / mf_cloud` 云商品

本质是云资源规划器商品，购买页不能按普通配置项页做。

必须独立支持：

- 推荐套餐
- 自定义配置
- 区域
- 镜像
- 网络
- VPC
- 带宽 / 流量
- IPv4 / IPv6
- 安全组
- NAT
- 磁盘规划

---

## 8. 上下游与 Provider 架构重建方案

### 8.1 Provider 统一抽象

必须有统一接口，例如：

- GetSummary
- Sync
- Create
- Suspend
- Unsuspend
- Terminate
- PowerOn
- PowerOff
- Reboot
- HardOff
- HardReboot
- Reinstall
- ResetPassword
- Rescue
- ExitRescue
- VNC
- ChangeConfig
- Renew

### 8.2 Provider 类型

#### `zjmf_api`

它不是底层实例 provider，而是：

- 上游财务 / 资源系统代理层

链路是：

- 同步商品
- 获取配置项
- 下单
- 上游结算
- 上游开通
- 回写上游 host id
- 再通过上游 host/header 驱动本地详情页

#### `dcimcloud`

它是：

- 本地原生云 provider

链路包括：

- 实例
- 快照
- 备份
- 安全组
- VNC
- 重装
- 救援

#### `manual`

它不是完整 provider，更像：

- 手工资源池映射层

#### `resource`

本质是：

- 资源型上游 provider

#### `upper_reaches`

它不是云实例动作层，而是：

- 供应商 / 上游资源池管理层

### 8.3 我们系统里必须分三层存储

不能只存一个 services 表糊住。

至少要有：

- 本地服务对象
- provider 绑定对象
- provider 资源快照 / 同步结果

例如：

- `services`
- `service_provider_bindings`
- `service_provider_snapshots`
- `service_actions`
- `service_sync_logs`

---

## 9. 用户中心重建方案

### 9.1 用户中心不是简化版后台

它是独立的业务控制台。

一级导航应至少包括：

- 用户中心首页
- 产品与服务
- 账户管理
- 财务管理
- 技术支持

### 9.2 用户中心首页

不是统计页，而是：

- 用户信息卡
- 账户余额
- 已出产品数量
- 我的产品
- 公告通知

### 9.3 服务列表

服务列表不能只显示订单。

必须显示：

- 服务状态
- 到期时间
- IP
- 核心配置摘要
- 产品名
- 周期

### 9.4 服务详情

必须按老系统保留：

- 头卡摘要
- 电源动作
- VNC
- 重装系统
- 登录信息
- 快照
- 配置
- 图表 / 流量
- 远程信息

这里必须按 `host/header` 驱动，而不是只用本地静态 DTO。

### 9.5 购物车与云商品购买页

普通商品和云商品必须分开。

#### 普通商品页

- 基础配置
- 周期选择
- 价格汇总

#### 云商品页

- 推荐套餐
- 自定义模式
- 区域
- 镜像
- CPU
- 内存
- 网络
- 带宽
- 磁盘
- 安全组
- NAT
- 登录方式
- 服务端 `ordersummary`

### 9.6 财务页

用户中心财务模块至少要有：

- 账单列表
- 账单详情
- 交易记录
- 余额
- 信用额
- 充值

### 9.7 工单页

必须是完整业务工单，不是留言板：

- 工单列表
- 提交工单
- 工单详情
- 回复
- 关闭
- 评价

---

## 10. 数据库重建方案

### 10.1 核心表

#### 客户域

- customers
- customer_contacts
- customer_identities
- customer_groups
- customer_levels
- customer_notes
- customer_attachments

#### 商品域

- product_groups
- products
- product_pricings
- product_config_groups
- product_config_options
- product_config_option_values
- product_custom_fields
- product_provider_templates
- product_upgrade_maps

#### 订单财务域

- orders
- order_items
- invoices
- invoice_items
- payments
- payment_logs
- wallet_ledgers
- credit_limits
- refunds
- invoice_logs

#### 服务域

- services
- service_config_snapshots
- service_custom_field_values
- service_provider_bindings
- service_provider_snapshots
- service_actions
- service_sync_logs
- service_renewals
- service_upgrades

#### Provider / 上下游域

- providers
- provider_accounts
- provider_products
- provider_product_mappings
- provider_sync_tasks
- provider_sync_logs
- upstream_orders
- upstream_order_items

#### 工单域

- tickets
- ticket_replies
- ticket_notes
- ticket_departments
- ticket_statuses
- ticket_transfers

#### 平台域

- admins
- roles
- permissions
- menus
- audit_logs
- settings
- jobs

### 10.2 关键关系

- orders.invoice_id -> invoices.id
- invoice_items.invoice_id -> invoices.id
- invoice_items.service_id -> services.id
- services.customer_id -> customers.id
- services.order_id -> orders.id
- services.product_id -> products.id
- services.provider_binding_id -> service_provider_bindings.id
- product_config_option_values -> service_config_snapshots

### 10.3 关键原则

- 服务配置快照必须结构化存储
- 不要只存一个大 JSON
- 账单项必须能反查服务
- provider 映射必须独立，不要直接混到 services 所有字段里

---

## 11. 接口重建思路

### 11.1 后台接口

#### 客户

- `GET /api/v1/admin/customers`
- `GET /api/v1/admin/customers/{id}`
- `GET /api/v1/admin/customers/{id}/services`
- `GET /api/v1/admin/customers/{id}/invoices`
- `GET /api/v1/admin/customers/{id}/transactions`

#### 订单

- `GET /api/v1/admin/orders`
- `GET /api/v1/admin/orders/{id}`
- `POST /api/v1/admin/orders/{id}/approve`
- `POST /api/v1/admin/orders/{id}/cancel`
- `DELETE /api/v1/admin/orders/{id}`

#### 账单

- `GET /api/v1/admin/invoices`
- `GET /api/v1/admin/invoices/{id}`
- `POST /api/v1/admin/invoices/{id}/receive-payment`
- `POST /api/v1/admin/invoices/{id}/refund`
- `PATCH /api/v1/admin/invoices/{id}`
- `POST /api/v1/admin/invoices/{id}/send-email`

#### 服务

- `GET /api/v1/admin/services`
- `GET /api/v1/admin/services/{id}`
- `PATCH /api/v1/admin/services/{id}`
- `POST /api/v1/admin/services/{id}/actions/sync`
- `POST /api/v1/admin/services/{id}/actions/vnc`
- `POST /api/v1/admin/services/{id}/actions/reinstall`
- `POST /api/v1/admin/services/{id}/actions/reset-password`

#### 商品

- `GET /api/v1/admin/products`
- `GET /api/v1/admin/products/{id}`
- `PATCH /api/v1/admin/products/{id}`
- `POST /api/v1/admin/products/{id}/sync-upstream-info`
- `POST /api/v1/admin/products/{id}/sync-upstream-price`
- `POST /api/v1/admin/products/{id}/sync-upstream-config`

#### Provider

- `GET /api/v1/admin/providers`
- `GET /api/v1/admin/providers/{id}/products`
- `POST /api/v1/admin/providers/{id}/sync-products`
- `POST /api/v1/admin/providers/{id}/sync-services`

### 11.2 用户中心接口

- `GET /api/v1/portal/dashboard`
- `GET /api/v1/portal/services`
- `GET /api/v1/portal/services/{id}/header`
- `GET /api/v1/portal/services/{id}/details`
- `GET /api/v1/portal/services/{id}/custom/{key}`
- `POST /api/v1/portal/services/{id}/actions/{action}`
- `GET /api/v1/portal/orders`
- `GET /api/v1/portal/invoices`
- `POST /api/v1/portal/orders/checkout`
- `POST /api/v1/portal/invoices/{id}/pay`

---

## 12. 对当前 `idc-finance` 的具体改造策略

### 12.1 不推倒重来，但要换骨架

当前 `idc-finance` 已经有一些基础：

- Go 后端
- admin-web
- portal-web
- 基本 auth
- 基础客户 / 商品 / 订单 / 账单 / 服务骨架

这些可以保留。

但要改的不是“补几个功能”，而是换成老魔方的工作台骨架。

### 12.2 需要重做的部分

#### A. 后台菜单 IA

当前后台菜单还是偏开发视角，要改成老魔方运营视角。

#### B. 页面模型

把现有零散页面收敛成：

- 列表工作台
- 详情工作台

#### C. 服务详情页

这是当前最该重做的一页。

必须按老魔方拆成：

- 客户上下文
- Provider 动作
- 财务区
- 配置区

#### D. 商品编辑页

必须支持：

- 价格矩阵
- 可配置项组
- 云模板
- 上游商品同步
- 上游价格同步
- 上游配置同步

#### E. 用户中心服务详情

必须按 `host/header` 结构重做，不要继续做普通信息页。

### 12.3 保留但下沉的部分

- 审计
- RBAC
- MySQL 持久化
- 工作台首页聚合

这些仍然保留，但不要继续抢主优先级。

---

## 13. 改造顺序

### Phase A：定骨架

1. 改后台菜单 IA
2. 改后台列表工作台骨架
3. 改后台详情工作台骨架
4. 改用户中心菜单 IA

### Phase B：重做三条主链

1. 订单工作台
2. 账单工作台
3. 服务工作台

### Phase C：重做商品与购买链

1. 商品编辑工作台
2. 普通商品下单链
3. `zjmf_api` 商品同步链
4. `mf_cloud` 云商品购买页

### Phase D：重做用户中心

1. 用户中心首页
2. 服务列表
3. 服务详情
4. 账单 / 交易 / 工单

### Phase E：重做 Provider 与上下游

1. `zjmf_api`
2. `dcimcloud`
3. `manual`
4. `resource`
5. `upper_reaches`

### Phase F：补全边角与运营能力

1. 发票
2. 合同
3. API 概览
4. 通知日志
5. 分销 / 代理
6. 统计报表

---

## 14. 下一轮应该怎么开工

如果直接进入改造，我建议第一轮不要碰太多页，而是先把骨架换正。

### 第一轮改造目标

只做 4 件事：

1. 重排后台 IA
2. 抽出统一列表工作台组件
3. 抽出统一详情工作台组件
4. 重做服务详情页骨架

原因很简单：

- 服务详情页是老魔方最复杂也最核心的一页
- 一旦这页骨架对了，订单详情、账单详情、用户中心服务详情都会跟着顺

### 第一轮改造后应达到的效果

- 后台菜单结构像老魔方
- 订单 / 账单 / 服务三条链的入口关系正确
- 服务详情页不再“像个普通实例页”
- 用户中心后面有了正确参照

---

## 15. 我的判断

现在已经不是“要不要继续理解”的阶段了。

正确判断是：

- 主骨架已经足够清楚
- 可以开始系统改造
- 同时保留一条并行研究线，继续补低频细节

也就是说，接下来最正确的做法是：

- 主线程：改 `idc-finance`
- 支线程：继续对照老魔方补边角工作台和低频动作

结论：

- 现在可以正式开始改造我们的系统
- 而且应该从后台工作台骨架和服务详情页开始
