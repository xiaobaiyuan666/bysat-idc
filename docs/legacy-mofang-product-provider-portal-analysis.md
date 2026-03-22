# 魔方财务 3.7.4 深度拆解：商品架构、资源对接、用户中心

更新时间：2026-03-21

## 0. 本次使用的三层证据

这份文档不是只靠页面猜的，而是同时用了三层证据：

1. 本地已跑通的 3.7.4 老系统
2. 你给的 decoded 仓库 3.5.8 可读源码
3. 本地数据库结构与前台模板/后台 bundle

本地 decoded 仓库路径：

- [zjmf-manger-decoded-ref](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref)

这个仓库 README 已明确说明它对应的是 `3.5.8`，所以它的作用不是直接替代你本地跑起来的 `3.7.4`，而是补出那些在 3.7.4 里被 ionCube 掉的逻辑层。

本次重点参考的可读逻辑目录：

- [app/common/logic](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic)
- [app/admin/controller](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/admin/controller)
- [app/home/controller](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/home/controller)

关键可读逻辑文件：

- [Product.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/Product.php)
- [Pricing.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/Pricing.php)
- [ConfigOptions.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/ConfigOptions.php)
- [Shop.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/Shop.php)
- [Provision.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/Provision.php)
- [Host.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/Host.php)
- [Renew.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/Renew.php)
- [Upgrade.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/Upgrade.php)
- [Dcim.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/Dcim.php)
- [DcimCloud.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/DcimCloud.php)
- [UpperReaches.php](C:/Users/Administrator/Desktop/IDC/zjmf-manger-decoded-ref/app/common/logic/UpperReaches.php)

## 1. 结论先行

这套老版魔方财务不是“先有云实例，再补财务页”，而是反过来：

`客户 -> 商品 -> 配置项/价格矩阵 -> 购物车/订单 -> 账单/支付 -> 服务实例 -> Provider 动作 -> 工单/日志`

魔方云、DCIM、通用接口、上游资源都不是独立业务主线，而是挂在下面四层：

1. 自动化接口层
2. 商品资源映射层
3. 服务交付层
4. 异步动作/同步日志层

这也是后续重建必须严格照搬的核心结构。

## 2. 商品架构

### 2.1 商品创建第一层字段

后台 `商品管理 -> 新增商品` 的第一步表单已经确认，只有 4 个核心字段：

- 商品名称
- 商品类型
- 商品组
- 会员中心导航分类

这说明商品在系统里一开始就同时绑定了两件事：

- 销售归属：商品组
- 前台展示归属：会员中心导航分类

不是简单的“后台创建产品记录”。

来源：

- 页面快照：[page-2026-03-20T19-37-04-365Z.yml](C:/Users/Administrator/Desktop/IDC/.playwright-cli/page-2026-03-20T19-37-04-365Z.yml)
- 商品管理页：[ProductServer~31ecd969.8b33ff37.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/admin01/js/ProductServer~31ecd969.8b33ff37.js)

### 2.2 商品主表设计

核心商品表是 `shd_products`，它不是纯展示型产品表，而是“售卖 + 交付 + 上游绑定”混合模型。

关键字段：

- `id`
- `type`
- `gid`
- `name`
- `description`
- `pay_type`
- `pay_method`
- `auto_setup`
- `server_type`
- `server_group`
- `api_type`
- `zjmf_api_id`
- `upper_reaches_id`
- `upstream_pid`
- `config_option1..24`

关键判断：

- `type` 决定商品大类
- `gid` 决定商品分组
- `auto_setup` 决定付款后是否自动开通
- `server_type/server_group` 决定落到哪类节点/资源组
- `api_type/zjmf_api_id/upper_reaches_id/upstream_pid` 决定商品最终走哪种外部资源交付方式

也就是说，一个商品对象从设计上就已经知道自己是：

- 本地模块商品
- 魔方 DCIM 商品
- 魔方云商品
- 上游资源商品
- 财务/API 型商品

### 2.3 配置项和价格矩阵

商品不是直接定一口价，而是由“基础价格 + 配置项 + 周期 + 资源规则”拼出来。

核心表：

- `shd_pricing`
- `shd_product_config_groups`
- `shd_product_config_links`
- `shd_product_config_options`
- `shd_product_config_options_sub`
- `shd_host_config_options`

#### `shd_pricing`

这是统一价格矩阵表，覆盖：

- 商品
- 附加项
- 配置项
- 域名注册/转入/续费

确认字段：

- `type`
- `currency`
- `relid`
- `onetime`
- `hour`
- `day`
- `monthly`
- `quarterly`
- `semiannually`
- `annually`
- `biennially`
- `triennially`
- 更长年限价格字段
- 对应 `setupfee`

关键结论：

- 魔方的“价格体系”不是按产品表直接存，而是按 `type + relid + currency` 做统一矩阵
- 商品、配置项、附加项共用一套计费引擎

#### 配置项组链路

- `shd_product_config_groups`：配置项组
- `shd_product_config_links`：配置组和商品的绑定关系
- `shd_product_config_options`：组内配置项
- `shd_product_config_options_sub`：配置项可选值
- `shd_host_config_options`：服务实例实际选中的配置快照

确认字段特点：

- 配置组支持 `global`
- 配置项支持 `option_type`
- 配置项支持 `upgrade`
- 配置项支持 `upstream_id`
- 配置项支持 `auto`
- 配置项支持 `senior`

关键结论：

- 配置项不是只为前台下单存在，也直接参与升级、上游映射和服务实例快照
- 服务实例一旦开通，选中的配置项会固化到 `shd_host_config_options`

## 3. 资源与 Provider 对接方式

### 3.1 老系统不是直接把魔方云当实例表

后台菜单位置已经确认：

- 商品设置
  - 自动化接口
    - 通用接口
    - 魔方 DCIM
    - 魔方云

这说明魔方云在后台首先是“接口/节点配置中心”，不是“客户云主机列表”。

来源：

- 后台菜单返回：[legacy-admin-login-response.json](C:/Users/Administrator/Desktop/IDC/docs/legacy-admin-login-response.json)
- 路由 chunk：[app~5a11b65b.061d499b.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/admin01/js/app~5a11b65b.061d499b.js)

### 3.2 魔方云接口管理页真实字段

`#/zjmfcloud` 页面已经实测确认是“接口管理页”。

新增接口弹窗字段：

- 名称
- 高级
- 地址
- 账号
- 密码
- 是否 https
- 是否启用

页面提示明确写了：

- 地址要填写“域名或 IP + 端口 + 后台管理员路径”
- 如果要使用 VNC，必须配置成域名并保证证书正常

这说明老系统对魔方云接口的建模就是：

- 一个被后台管理的远端控制面节点
- 节点本身承载实例操作能力
- VNC 是节点级接入要求，不是服务级拼接逻辑

来源：

- 页面快照：[page-2026-03-20T19-31-21-950Z.yml](C:/Users/Administrator/Desktop/IDC/.playwright-cli/page-2026-03-20T19-31-21-950Z.yml)

### 3.3 节点表与扩展表

资源接入不是一张表，而是基础节点 + 模块扩展。

核心表：

- `shd_servers`
- `shd_server_groups`
- `shd_server_groups_rel`
- `shd_dcim_servers`

#### `shd_servers`

这是通用自动化接口节点表，关键字段：

- `gid`
- `name`
- `ip_address`
- `hostname`
- `username`
- `password`
- `secure`
- `port`
- `active`
- `server_type`
- `type`
- `link_status`

#### `shd_dcim_servers`

这是 DCIM/魔方云扩展配置表，关键字段：

- `serverid`
- `reinstall_times`
- `buy_times`
- `reinstall_price`
- `auth`
- `api_status`
- `area`
- `os`
- `bill_type`

关键结论：

- 老系统把“通用接口节点”和“DCIM/魔方云能力扩展”分开了
- 魔方云不是单独另一套表，而是基于节点表扩一层能力字段

### 3.4 商品如何绑定资源

商品不是直接绑定“实例”，而是绑定“节点/资源来源类型”。

主要通过这些字段和表：

- `shd_products.server_type`
- `shd_products.server_group`
- `shd_products.api_type`
- `shd_products.zjmf_api_id`
- `shd_products.upper_reaches_id`
- `shd_products.upstream_pid`

以及：

- `shd_zjmf_finance_api`
- `shd_upper_reaches`
- `shd_upper_reaches_res`
- `shd_upper_manual_info`

也就是说，商品交付来源有至少四种：

1. 通用模块节点
2. DCIM/魔方云节点
3. 上游财务/API 对接
4. 手动/上游资源池

### 3.5 资源动作如何执行

真正执行服务动作的不是商品页，而是服务与模块动作层。

已确认的后台动作入口：

- `admin/provision/default`
- `admin/provision/custom`
- `admin/dcimCloud/server/*`
- `admin/upper/dcim_client/*`
- `admin/upper/ipmi/*`

已确认的前台动作入口：

- `home/provision/execute`
- `home/provision/customFunc`
- `openapi/Host/module/*`

关键结论：

- 后台、前台、OpenAPI 三套入口，最终都汇聚到“Host 模块动作模型”
- 魔方云能力在对外接口上不会单独叫 `mf_cloud` 控制器，而是统一落在 `Host module` 动作体系里

### 3.6 异步交付与运行日志

服务不是支付成功后立即同步完成一切，而是进入交付和动作记录层。

核心表：

- `shd_module_queue`
- `shd_run_maping`
- `shd_zjmf_pushhost`

职责：

- `shd_module_queue`：模块动作异步队列
- `shd_run_maping`：服务动作运行记录/状态映射
- `shd_zjmf_pushhost`：服务推送/回调重试日志

这说明老系统在架构上明确区分：

- 财务事件
- 服务对象
- 资源交付动作
- 动作结果与回调

## 4. 用户中心架构

### 4.1 用户中心不是一个单页，而是一整套主题系统

前台主题目录已经确认：

- [public/themes/clientarea/default](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/themes/clientarea/default)
- [public/themes/cart/default](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/themes/cart/default)

这套前台实际上拆成三层：

1. 客户中心主题 `clientarea`
2. 购物车/下单主题 `cart`
3. v10 商品配置与服务详情增强层

### 4.2 用户中心页面树

根据模板和路由，前台主要模块是：

- 登录与注册
  - `login.tpl`
  - `register.tpl`
  - `pwreset.tpl`
  - `bind.tpl`
  - `security.tpl`
  - `loginlog.tpl`
  - `apilog.tpl`
  - `apimanage.tpl`
- 控制台与账户
  - `clientarea.tpl`
  - `details.tpl`
  - `systemlog.tpl`
  - `message.tpl`
- 服务中心
  - `service.tpl`
  - `servicedetail.tpl`
  - `servicedetail-v10-common.tpl`
  - `servicedetail-v10-cloud.tpl`
  - `servicedetail-v10-dcim.tpl`
  - `service_cloud.tpl`
  - `service_cloud_server.tpl`
  - `service_server.tpl`
  - `service_product.tpl`
  - `service_hosting.tpl`
  - `service_domain.tpl`
  - `service_ssl.tpl`
  - `service_sms.tpl`
  - `service_soft.tpl`
  - `service_cdn.tpl`
- 财务中心
  - `billing.tpl`
  - `viewbilling.tpl`
  - `combinebilling.tpl`
  - `transaction.tpl`
  - `addfunds.tpl`
  - `credit.tpl`
  - `creditdetail.tpl`
  - `invoicelist.tpl`
  - `invoiceapply.tpl`
  - `invoiceaddress.tpl`
  - `invoicecompany.tpl`
  - `contract.tpl`
  - `contracthost.tpl`
- 工单中心
  - `supporttickets.tpl`
  - `submitticket.tpl`
  - `viewticket.tpl`
- 分销中心
  - `affiliates.tpl`
  - `affiliates/*`
- 内容中心
  - `news*.tpl`
  - `knowledgebase*.tpl`
  - `downloads*.tpl`

关键结论：

- 用户中心不是简单的“订单/账单/服务”三页
- 它还包含 API、登录安全、实名、合同、发票、分销、消息、知识库

### 4.3 服务详情是前台真正的工作台

前台服务详情动作脚本已经确认：

- 产品升降级
- 配置升降级
- 流量包订购
- 续费
- 模块通用按钮
- 重装系统
- 重置密码
- 救援模式
- VNC / KVM / BMC 等控制台打开

关键实现特征：

- 客户端动作不直接散落在页面模板里，而是统一 POST 到 `/provision/default` 或 `/provision/custom/:id`
- `service_module_button` 根据 `func`、`type`、`desc` 组装动作
- 对 `reinstall`、`crack_pass`、`rescue_system` 这类高风险动作会先弹确认或验证模态框
- 如果返回 `res.data.url`，前台直接 `window.open(res.data.url)` 打开外部控制台

来源：

- [servicedetail.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/themes/clientarea/default/assets/js/servicedetail.js)

### 4.4 v10 购物车明确拆了魔方云和 DCIM

购物车 v10 目录里已经确认有：

- [mf_cloud.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/themes/cart/default/v10/api/mf_cloud.js)
- [mf_dcim.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/themes/cart/default/v10/api/mf_dcim.js)
- `configureproduct-v10-cloud.tpl`
- `configureproduct-v10-dcim.tpl`
- `configureproduct-v10-common.tpl`

这说明商品配置前台也不是统一一套表单，而是按资源类型拆了不同配置器。

`mf_cloud.js` 已确认的前台能力：

- 获取云产品列表
- 获取订购页面配置
- 获取镜像列表
- 获取周期价格
- 动态计算价格
- 结算商品
- 使用优惠码
- 获取购物车
- 获取线路详情
- 获取 SSH Key
- 获取安全组
- 获取 VPC
- 获取客户等级折扣
- 加入购物车
- 编辑购物车
- 读取商品自定义字段

关键接口：

- `/product/common_cloud`
- `/v10/host/product/{id}/mf_cloud/order_page`
- `/v10/host/product/{id}/mf_cloud/image`
- `/v10/host/product/{id}/mf_cloud/duration`
- `/product/{id}/config_option`
- `/product/settle`
- `/promo_code/apply`
- `/product/{id}/v10/host/mf_cloud/line/{line_id}`
- `/v10/host/product/{id}/mf_cloud/vpc_network/search`
- `/cart/v10/add`
- `/cart/v10/edit`

关键结论：

- 魔方云商品前台下单器本身就是“云资源配置器”，不是普通购物车表单
- 前台配置器直接暴露了镜像、线路、SSH Key、安全组、VPC 这些云概念

## 5. 订单支付后如何进入交付

综合数据库、路由、前台脚本和后台结构，目前可确定的交付顺序是：

1. 商品在后台先完成商品定义、价格矩阵、配置项绑定、资源来源绑定
2. 前台购物车按资源类型调用不同配置器
3. 结算时生成订单、账单、账单项
4. 支付动作只先完成财务闭环
5. 账单支付后再驱动服务实例 `shd_host`
6. 服务实例再通过 Provider 动作层调用通用接口、DCIM、魔方云、上游资源
7. 动作结果进入队列、运行日志、回调日志

关键对象顺序：

`shd_products -> shd_pricing -> shd_product_config_* -> shd_orders -> shd_invoices -> shd_invoice_items -> shd_host -> shd_module_queue / shd_run_maping`

## 6. 对新系统重建的直接要求

如果要做接近 1:1 的新系统，商品与资源这块至少必须照这个方式重建：

### 6.1 商品域

- 商品组
- 商品
- 会员中心导航分类
- 统一价格矩阵
- 配置项组
- 配置项
- 配置项可选值
- 商品升级映射
- 商品订购设置

### 6.2 资源接入域

- 通用接口节点
- DCIM 节点
- 魔方云节点
- 上游供应商
- 上游资源池
- 手动资源池
- Provider 元数据
- Provider 动作日志
- 异步任务队列

### 6.3 用户中心

- 客户控制台首页
- 服务列表
- 服务详情工作台
- 账单/合并账单/充值/余额/信用额
- 工单
- 发票与合同
- API 管理与 API 日志
- 登录日志与安全设置
- 分销中心

## 7. 当前还缺的最后两块

这轮已经把骨架拆出来了，但还有两块还要继续补：

1. 前台真实登录后的导航和页面跳转实测
2. 商品编辑页第二层、第三层的详细字段取证

完成这两块后，就可以把这套老系统的：

- 商品建模
- 资源交付模型
- 用户中心 IA
- 前后台动作链

整理成真正可重建的 1:1 设计稿。
