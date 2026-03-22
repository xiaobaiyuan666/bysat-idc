# 魔方财务老系统学习笔记

## 学习来源

- 官方使用文档总入口: https://www.idcsmart.com/wiki_list/479.html
- 官方开发文档总入口: https://www.idcsmart.com/wiki_list/478.html
- 官方接口文档: https://www.idcsmart.com/wiki_list/558.html
- 官方服务器模块文档: https://www.idcsmart.com/wiki_list/556.html
- 官方对接魔方云文档: https://www.idcsmart.com/wiki_list/669.html
- 官方财务系统使用导读: https://www.idcsmart.com/wiki_list/695.html
- 本地源码包: `C:\Users\Administrator\Desktop\IDC\idccw.zip`

## 结论先行

这套老的魔方财务，本质上不是“几个管理页”，而是一套完整的云业务运营平台。它的核心不是 UI，而是以下 8 条业务主链：

1. 客户与实名
2. 商品与配置项
3. 购物车与订单
4. 账单、支付、交易流水、余额
5. 服务实例生命周期
6. 工单协同
7. 插件、服务器模块、HOOK 扩展
8. 后台权限、员工、计划任务、通知模板

如果我们要做“和它一模一样，但更现代”的系统，应该复制的是这套领域结构、流程编排和页面工作台思路，而不是复制它的旧 PHP 实现。

## 从源码包反推出的系统结构

`idccw.zip` 根目录为 `3.7.4/`。核心 PHP 业务代码大量使用 ionCube 加密，不能直接阅读逻辑细节，但目录结构、模板、安装 SQL、权限表和控制器命名已经足够反推出系统架构。

### 代码层分区

- `app/admin/controller`
  后台控制器，覆盖客户、订单、账单、交易、工单、商品、服务器、优惠码、合同、插件、权限、员工、统计、自动任务等。
- `app/home/controller`
  前台客户中心控制器，覆盖购物车、订单、支付、账单、服务、工单、实名、充值、API 管理、发票、下载、帮助中心等。
- `app/api/controller`
  开放 API 与安装升级相关接口。
- `app/common/logic`
  公共业务逻辑层，包含 `Cart`、`Invoices`、`PaymentGateways`、`Pricing`、`Product`、`Provision`、`Renew`、`Host`、`DcimCloud`、`Dcim`、`Contract`、`Sms`、`SystemMessage` 等服务类。
- `public/admin/themes/default`
  后台主题模板。
- `public/themes/clientarea/default`
  客户中心模板。
- `public/install/thinkcmf.sql`
  安装 SQL，能直接反推出完整主表。
- `public/upgrade/*.sql`
  历次升级 SQL，能看到系统是长期迭代的而不是一次性产品。

### 后台控制器命名能反推出的领域

后台控制器里最关键的一批如下：

- 客户域: `UserManageController`、`ClientGroupController`、`ClientsContactsController`、`ClientCareController`
- 业务域: `OrderController`、`HostController`、`ClientsServicesController`、`CancelRequestController`
- 财务域: `InvoiceController`、`InvoiceItemsController`、`AccountController`、`CreditController`、`CreditLimitController`、`ContractController`
- 商品域: `ProductController`、`ConfigOptionsController`、`PromoCodeController`、`VoucherController`
- 资源域: `ConfigServersController`、`DcimCloudController`、`DcimController`
- 工单域: `TicketController`、`TicketDepartmentController`、`TicketStatusController`、`TicketPrereplyController`、`TicketDeliverController`
- 平台域: `PluginController`、`RbacController`、`CronController`、`ApiController`、`EmailTemplateController`、`SendMessageController`

这说明它的模块划分不是“页面驱动”，而是标准的业务域划分。

## 后台信息架构

结合官方使用文档、权限 SQL 和之前对线上后台的只读观察，可以把后台 IA 定义得非常清楚。

### 一级模块

- 客户
- 业务
- 财务
- 工单
- 功能
- 设置
- 资源与商店
- 反馈

### 客户模块

主要负责客户对象本身和销售跟进：

- 客户列表
- 实名认证
- 客户资源池
- 客户分组
- 客户等级
- 客户联系人
- 跟进记录
- 我的业绩
- 营销推送
- 推介计划

### 业务模块

核心是区分“订单”和“已开通业务”：

- 产品订单
- 续费订单
- 流量包订单
- 业务列表
- 产品暂停请求

这个拆分很关键。订单是交易对象，业务是服务对象，两者不能混。

### 财务模块

财务被拆成多层，不是一个流水表：

- 交易流水
- 账单管理
- 余额管理
- 信用额管理
- 提现审核
- 发票列表
- 合同列表
- 支付接口
- 优惠码
- 代金券

### 工单模块

- 工单列表
- 工单统计
- 工单部门
- 工单状态
- 预设回复
- 派单相关

### 设置模块

- 基础设置
- 系统设置
- 商品设置
- 站务设置
- 邮件模板
- 消息模板
- 服务器配置
- 可配置选项组
- 自动任务
- 权限与员工

## 客户端信息架构

从 `public/themes/clientarea/default` 可以反推出客户中心不是单页，而是一套完整门户。

### 主功能页

- `clientarea.tpl`
  客户首页工作台
- `service.tpl`
  服务列表
- `servicedetail.tpl`
  服务详情
- `billing.tpl`
  账单列表
- `viewbilling.tpl`
  账单详情
- `transaction.tpl`
  交易记录
- `credit.tpl`
  余额与信用额度
- `addfunds.tpl`
  充值
- `supporttickets.tpl`
  工单列表
- `viewticket.tpl`
  工单详情
- `apimanage.tpl`
  API 管理
- `invoiceapply.tpl`
  发票申请
- `contract.tpl`
  合同
- `downloads.tpl`
  下载中心
- `knowledgebase.tpl`
  帮助中心

### 客户中心首页工作台

从模板内容可以确认首页至少包含：

- 个人信息卡
- 账户安全状态
- 待处理工单数
- 未支付订单数
- 产品总数
- 当前余额
- 本月消费
- 未支付账单
- 充值入口
- 已开通产品分组统计
- 资源列表
- 公告通知

这不是简单 dashboard，而是客户自助门户的主控台。

### 服务列表设计

`service.tpl` 说明客户侧服务列表具备完整运营能力：

- 状态多选筛选
- 关键字搜索
- 表头排序
- IP 展示
- 到期日展示
- 周期价格展示
- 自动续费状态展示
- 当前系统展示
- 备注编辑
- 进入服务详情操作

### 账单列表设计

`billing.tpl` 说明客户侧账单页具备：

- 状态筛选
- 排序
- 查看账单
- 支付未付账单
- 删除未付账单
- 合并付款
- 分页和每页数量控制

这意味着客户端必须是“自助业务中心”，不是只读看板。

## 数据模型主干

`public/install/thinkcmf.sql` 里可直接读到主表。最关键的表分组如下。

### 客户与权限

- `shd_clients`
- `shd_contacts`
- `shd_client_groups`
- `shd_certification`
- `shd_certifi_person`
- `shd_certifi_company`
- `shd_user`
- `shd_role`
- `shd_auth_rule`
- `shd_role_user`

### 商品与定价

- `shd_products`
- `shd_product_groups`
- `shd_pricing`
- `shd_product_config_groups`
- `shd_product_config_options`
- `shd_product_config_options_sub`
- `shd_product_config_links`
- `shd_customfields`
- `shd_customfieldsvalues`

### 订单与账单财务

- `shd_orders`
- `shd_invoices`
- `shd_invoice_items`
- `shd_accounts`
- `shd_credit`
- `shd_credit_limit`
- `shd_pay_log`
- `shd_payment_gateways`
- `shd_promo_code`
- `shd_voucher`
- `shd_voucher_invoices`
- `shd_contract`
- `shd_contract_pdf`

### 服务与资源

- `shd_host`
- `shd_host_category`
- `shd_host_config_options`
- `shd_servers`
- `shd_server_groups`
- `shd_server_groups_rel`
- `shd_dcim_servers`
- `shd_dcim_flow_packet`
- `shd_dcim_buy_record`
- `shd_cancel_requests`
- `shd_upgrades`
- `shd_renew_cycle`

### 工单与通知

- `shd_ticket`
- `shd_ticket_reply`
- `shd_ticket_note`
- `shd_ticket_status`
- `shd_ticket_department`
- `shd_ticket_prereply`
- `shd_ticket_deliver`
- `shd_email_templates`
- `shd_email_log`
- `shd_message_template`
- `shd_message_log`
- `shd_system_message`

### 平台与扩展

- `shd_plugin`
- `shd_hook`
- `shd_hook_plugin`
- `shd_api`
- `shd_jobs`
- `shd_module_queue`
- `shd_cron_log`
- `shd_configuration`

### 对我们最重要的结论

老系统的表设计很重，但业务边界是对的。我们新系统要保留它的领域拆分，不要保留它“一个大系统很多耦合表”的实现方式。

## 后台权限系统

`downloads/database/auth_rule.sql` 很有价值，因为它不是文档，而是真实后台菜单和权限定义。

从这个文件可以确认：

- 菜单和权限共用一套 `auth_rule`
- 菜单通过 `pid` 形成树结构
- 每个菜单既有标题，也有跳转地址，也有权限标识
- 菜单可配置显示与隐藏
- 多语言直接挂在权限定义上

这意味着老系统的后台导航不是前端静态写死，而是权限驱动的动态菜单。

### 直接确认到的关键菜单

- 客户列表: `/customer-list`
- 产品订单: `/order-list`
- 续费订单: `/renewal-order`
- 流量包订单: `/dcim-traffic-log`
- 业务列表: `/customer-product`
- 交易流水: `/business-statement`
- 账单管理: `/bill-management`
- 工单列表: `/support-ticket`
- 商品管理: `/product-server`
- 服务器设置: `/server-settings`
- 可配置选项: `/configurable-option`
- 支付接口: `/payment-interface`
- 优惠码: `/promo-code`
- 员工管理: `/admin-management`

## 核心业务流程

### 1. 客户创建与实名

- 客户注册或后台创建客户
- 维护基础资料与联系人
- 完成个人或企业实名认证
- 分配销售和客户分组
- 进入可下单、可支付、可开票状态

### 2. 商品下单

- 客户进入购物车
- 选择商品、周期、配置项、自定义字段
- 计算价格、折扣、优惠码、代金券
- 生成订单
- 生成账单
- 进入支付或人工审核流程

### 3. 财务闭环

- 订单生成账单
- 账单支付后生成交易流水
- 流水影响客户余额
- 财务可补录、删除、查看明细
- 支持信用额度、充值、提现、发票、合同

老系统很强调：

- 账单是应收对象
- 交易流水是资金对象
- 余额是账户对象

这三层必须分开。

### 4. 服务开通与生命周期

- 订单审核通过
- 调服务器模块或魔方云接口开通
- 生成 `host` 服务记录
- 服务进入 `Active / Suspended / Pending / Cancelled` 等状态
- 后续支持续费、升降配、暂停、解除暂停、删除、取消请求

### 5. 工单协同

- 客户提交工单
- 工单按部门、状态、派单规则流转
- 支持预设回复、附件、备注、合并工单、状态变更
- 工单详情本质上是客服工作台

## 开发与扩展架构

官方开发文档把它的扩展思路说明得很清楚。

### 1. API

官方提供 REST 风格 API 文档，页面可直接看到控制器接口与参数结构。说明这套系统不是纯页面系统，而是开放平台。

### 2. 服务器模块

服务器模块文档明确了服务交付是通过“模块适配器”完成的，不是写死在系统里。

文档可以确认的典型能力包括：

- 产品参数与模块参数配置
- 创建实例
- 暂停实例
- 解除暂停
- 删除实例
- 重启
- 开关机
- 修改密码
- 使用量同步
- 状态同步

这就是老系统能对接多种云平台、VPS、裸金属、DCIM 的根基。

### 3. HOOK 系统

开发文档里有完整 HOOK 分类，覆盖至少：

- 产品
- 购物车
- 账单
- 客户
- 工单

这意味着大量业务增强不是通过改核心代码，而是通过事件插桩实现。

### 4. 插件系统

源码里存在完整插件目录，文档里也有插件开发说明。说明支付、实名、OAuth、通知、样式、第三方服务等都可以插件化。

## 魔方云对接学习结果

官方对接魔方云文档说明了几个关键点：

- 财务侧商品类型要选“魔方云”
- 需要先在财务后台配置第三方平台
- 需要录入接口地址、管理员账号、管理员密码
- 商品里要配置模块参数
- 下单后可以直接创建实例

实际文档中能看到典型模块参数：

- `node`
- `os`
- `cpu`
- `memory`
- `system_disk`
- `data_disk`
- `bw`
- `peak_defence`
- `gpu`

这说明财务系统本身并不管理底层资源，它负责：

- 商品售卖
- 参数映射
- 订单驱动开通
- 账单与生命周期触发

真正的资源创建由魔方云平台完成。

## 老系统真正值得复刻的地方

### 1. 工作台化详情页

客户详情、订单详情、工单详情、服务详情都不是单纯展示页，而是单对象工作台。这个思路必须保留。

### 2. 订单、账单、流水、服务四层拆分

这是整套系统最重要的正确分层。我们当前系统最需要向它靠拢的，也是这四层对象的严格分离。

### 3. 权限驱动菜单

后台导航、页面权限、按钮权限应统一挂在 RBAC 上，不应该前端硬编码。

### 4. 商品和配置项体系

商品、产品分组、周期价格、可配置选项、自定义字段、服务器模块参数，是云业务平台的底层骨架。

### 5. 客户端自助中心

客户端必须具备：

- 服务自助管理
- 账单查看与支付
- 余额充值
- 工单提交与回复
- API 管理
- 发票/合同申请

## 老系统不值得照搬的地方

### 1. 技术栈太旧

- PHP 旧式 MVC
- 大量 Smarty 模板
- jQuery + Bootstrap 交互
- 业务逻辑偏重控制器与大表

### 2. 前后端边界不清

页面里混了大量模板逻辑、脚本和业务判断，这种方式不适合我们当前 Vue + Next 的前后端分离路线。

### 3. 表结构偏重历史包袱

升级 SQL 很多，历史字段和兼容字段会比较重。我们应学习对象边界，不要机械复刻表结构。

## 对我们新系统的重构建议

### 后台按领域重构

建议后台固定拆成这些主域：

- 工作台
- 客户
- 订单
- 业务
- 财务
- 工单
- 商品
- 资源
- 设置
- 平台

### 每个主域统一页面范式

列表页固定包含：

- 状态 tabs
- 快捷筛选
- 高级搜索
- 表格
- 分页
- 批量动作

详情页固定包含：

- 摘要区
- 状态区
- 关键动作区
- 关联 tabs
- 时间线或日志

### 客户端按自助中心重构

客户端固定拆成：

- 控制台首页
- 云产品
- 账单
- 订单
- 钱包
- 工单
- API
- 发票与合同
- 安全与实名

### 与魔方云的边界

我们系统负责：

- 商品售卖
- 财务规则
- 订单和账单
- 生命周期策略
- 本地客户、日志、通知、权限

魔方云负责：

- 实例资源创建
- 实例运行状态
- 开关机与重启
- VNC
- 底层云资源

## 当前学习结论对应到我们的下一步

后续如果继续按这套蓝图优化，我们不应该再“零碎补页面”，而应该按下面顺序整体重构：

1. 先重构后台 IA 和菜单树
2. 再重构四个核心对象: 客户、订单、账单、服务
3. 再补工单工作台和财务工作台
4. 再把客户端按自助门户重做
5. 最后把魔方云同步、动作、日志彻底纳入统一生命周期

## 备注

本地 `idccw.zip` 的核心 PHP 逻辑被 ionCube 加密，因此以下判断主要来自：

- 目录结构
- 控制器和逻辑类命名
- 未加密模板
- 安装 SQL 与权限 SQL
- 官方使用文档
- 官方开发文档

这类判断属于基于源码骨架和官方资料的架构反推，不是逐行代码审计结论。
