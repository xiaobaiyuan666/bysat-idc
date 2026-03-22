# 魔方财务 3.x 深度拆解：上下游、资源池、魔方云三条对接链

更新时间：2026-03-21

## 1. 结论先行

老魔方财务里最容易混淆的不是页面，而是三条完全不同的对接链：

1. `上下游财务 API`  
   对应 `zjmf_finance_api`，核心是“对接另一套魔方财务站点”，同步商品、下单、支付、服务头信息、升级、续费、工单等。

2. `资源池 / 代理资源`  
   对应 `zjmf_finance_api.is_resource = 1`、`AgentController`、`upper_reaches_res` 等，核心是“资源代理与资源巡检/售后”，不是普通客户商品同步。

3. `魔方云 / DCIMCloud`  
   对应 `dcimcloud` 商品类型、`DcimCloud.php`、后台 `zjmfcloud`/`zjmfcloud-product` 页面，核心是“节点能力 + 云实例动作 + 商品映射”。

这三条链共享一部分底层表和门面方法，但业务目的不同：

- `上下游财务 API` 解决“卖别人的商品、调用对方财务系统接口、把对方实例映射到我方服务”。
- `资源池 / 代理资源` 解决“拿上游资源库存做资源分销、资源巡检、售后、工单和资源检查”。
- `魔方云 / DCIMCloud` 解决“对接某个云资源节点本身，执行实例操作和商品参数配置”。

## 2. 本次证据来源

### 2.1 可读源码

- [app/zjmf.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\zjmf.php)
- [app/admin/controller/ZjmfFinanceApiController.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\admin\controller\ZjmfFinanceApiController.php)
- [app/home/controller/ZjmfFinanceApiController.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\home\controller\ZjmfFinanceApiController.php)
- [app/common/logic/Host.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\Host.php)
- [app/common/logic/Provision.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\Provision.php)
- [app/common/logic/DcimCloud.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\DcimCloud.php)
- [app/common/logic/UpperReaches.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\UpperReaches.php)
- [app/common/model/ProductModel.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\model\ProductModel.php)
- [public/themes/clientarea/default/apimanage.tpl](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\public\themes\clientarea\default\apimanage.tpl)
- [public/themes/clientarea/default/servicedetail.tpl](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\public\themes\clientarea\default\servicedetail.tpl)
- [public/themes/clientarea/default/servicedetail/zjmfcloud.tpl](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\public\themes\clientarea\default\servicedetail\zjmfcloud.tpl)

### 2.2 本地数据库表结构

已确认相关表：

- `shd_zjmf_finance_api`
- `shd_api_user_product`
- `shd_upper_reaches`
- `shd_upper_reaches_res`
- `shd_products`
- `shd_pricing`
- `shd_host`
- `shd_product_config_groups`
- `shd_product_config_options`

### 2.3 真实远端联调

对你提供的站点做了真实调用验证，但这里不记录密钥本身。

验证结果：

- `POST /zjmf_api_login`：返回 `200`，能拿到 JWT  
- `POST /resource_login`：返回 `404`
- `GET /cart/all`：返回 `200`，商品总数 `103`，币种 `CNY`
- `GET /cart/get_product_config?pid=27`：返回 `200`
- `GET /user_info`：返回 `200`

这说明你给的这套接入信息，按老源码的判定规则，属于：

- `type = zjmf_api`
- 不是 `is_resource = 1` 的资源代理登录链

## 3. 上下游财务 API 的真实鉴权方式

### 3.1 鉴权门面

核心函数在 [app/zjmf.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\zjmf.php)：

- `zjmfCurl($api_id, $path, $data, $timeout, $request)`
- `zjmfCurlHasFile(...)`
- `zjmfApiLogin($id, $url, $data, $force = false)`

### 3.2 登录分流规则

老系统会先查 `zjmf_finance_api` 表，再按 `is_resource` 判断登录地址：

- `is_resource = 1`  
  登录到 `/resource_login`
- `is_resource = 0`  
  登录到 `/zjmf_api_login`

`zjmfApiLogin()` 会把 JWT 缓存在：

- `zjmf_finance_jwt_{id}`

缓存时长是：

- `5400` 秒

如果业务调用返回 `405`，会强制重新登录一次，再重试接口。

### 3.3 这条链和 WHMCS 兼容模式的区别

`admin/ZjmfFinanceApiController::refreshStatus()` 里还能看到第三种模式：

- `type = whmcs`
- 检查接口为 `/modules/addons/idcsmart_api/api.php?action=/v1/check`
- 参数是 `apiname/apikey`

所以 `zjmf_finance_api` 实际支持 3 种外部源：

1. `zjmf_api`
2. `manual`
3. `whmcs`

## 4. 上下游财务 API 管的到底是什么

### 4.1 后台页面职责

核心控制器是 [app/admin/controller/ZjmfFinanceApiController.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\admin\controller\ZjmfFinanceApiController.php)。

它不是简单“保存一组接口地址”，而是完整的上下游商品与客户 API 管理中心。

已确认功能包括：

- 新增接口源
- 修改接口源
- 删除接口源
- 接口连通性检查
- 商品数量刷新
- 接口概览
- 下游客户 API 开通/关闭
- 下游客户 API 密钥重置
- 免费产品分配
- 上游商品与本地商品映射查看

### 4.2 创建接口时的真实检查逻辑

当后台新增 `type = zjmf_api` 的接口源时，会立即：

1. 保存 `hostname/username/password`
2. 调 `zjmfCurl($id, "/cart/all", [], 15, "GET")`
3. 如果成功，把 `product_num` 更新成远端商品数

也就是说：

- “上下游接口是否有效”的第一判断，不是 ping 域名
- 而是“能不能真实登录并拉到对方售卖商品”

### 4.3 用户中心 API 管理

用户中心控制器在 [app/home/controller/ZjmfFinanceApiController.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\home\controller\ZjmfFinanceApiController.php)，页面模板在 [apimanage.tpl](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\public\themes\clientarea\default\apimanage.tpl)。

用户中心这条链的功能点是：

- 开通 API 功能
- 关闭 API 功能
- 重置 API 密钥
- 查看 API 概览
- 查看 7 天 API 调用趋势
- 查看免费 API 商品
- 查看 API 调用日志

对应的核心数据表和字段：

- `clients.api_open`
- `clients.api_password`
- `clients.api_create_time`
- `clients.lock_reason`
- `clients.api_lock_time`
- `api_resource_log`
- `api_user_product`

### 4.4 这条链真正服务的是“下游客户 API 能力”

从 `summary()` 和 `apimanage.tpl` 能确认：

- 它不是给管理员自己看的监控页
- 它是给客户自己开的“下游 API 管理中心”

客户开通之后，可以：

- 拿到自己的 API 密钥
- 调用我方开放的 API
- 管理被赠送/授权的免费产品
- 查看自己的 API 调用日志和资源调用情况

这就是老系统里“上游管理下游”的真正含义。

## 5. 真实站点联调得到的商品同步结构

### 5.1 `/cart/all` 的返回结构

实测远端站点 `/cart/all` 返回的是：

- `count`
- `currency`
- `products`

其中 `products` 不是单个商品数组，而是：

- 导航分类对象数组

每个分类对象形态是：

- `id`
- `name`
- `products[]`

而分类下的每个商品至少有：

- `id`
- `type`
- `name`
- `description`

这证明老系统同步商品时，最外层先同步的是：

- 会员中心/商城导航分类
- 分类下商品列表

不是先扁平拉所有商品。

### 5.2 远端站点里真实商品类型

从实际返回看，首批商品大量是：

- `type = dcimcloud`

说明这个上游站点本身在卖的主要是：

- 魔方云/DCIMCloud 型商品

也就是说：

- `上下游财务 API` 和 `魔方云商品` 在真实业务里经常是叠加关系
- 上游站点对外提供的是“财务 API 商品同步”
- 但同步进来的商品本身可能是 `dcimcloud`

### 5.3 `/cart/get_product_config` 的返回结构

对实际商品 `pid=27` 的配置接口验证结果：

- `status = 200`
- `type = dcimcloud`
- `stock_control = 1`
- `qty = 10`
- `allow_qty = 1`

返回对象包含：

- `products`
- `product_pricings`
- `config_groups`
- `config_links`
- `customfields`
- `advanced`
- `flag`

这和本地数据库结构一一对应：

- `products` -> `shd_products`
- `product_pricings` -> `shd_pricing`
- `config_groups` -> `shd_product_config_groups`
- `config_links` -> `shd_product_config_links`
- `customfields` -> `shd_customfields`

### 5.4 真实配置项组长什么样

实际拿到的配置项组已经能证明“同步不是只同步价格”，而是同步完整商品配置器。

已确认的配置项包括：

- `node|节点id`
- `os|操作系统`
- CPU/内存类配置
- `IP数量`
- `data_disk_size|数据盘`
- `port|远程端口`

并且每个配置项下的子项都带：

- `option_name`
- `pricings[]`
- 不同周期价格

例如：

- 节点子项有单独月付/季付/半年/年付加价
- IP 数量有单独月付/季付/半年/年付加价
- 数据盘有数量区间和单价

这说明老魔方的“商品同步”其实是在复制：

1. 商品基础信息
2. 商品价格矩阵
3. 配置项组定义
4. 配置项子项定义
5. 各子项价格矩阵
6. 自定义字段
7. 库存与试用规则

## 6. 产品同步到本地的真实写库逻辑

关键函数在 [ProductModel::syncProductInfoToResource](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\model\ProductModel.php)。

它做的不是“存一个 upstream_pid”，而是整套商品同步：

- 校验 `api_type`
- 校验 `zjmf_finance_api_id`
- 校验 `upstream_pid`
- 拉远端商品详情
- 根据汇率和加价规则换算价格
- 回写本地 `products`
- 清空并重建本地 `pricing`
- 清空并重建本地 `customfields`
- 同步配置项组和选项
- 更新本地版本号和上游版本号

关键字段包括：

- `products.zjmf_api_id`
- `products.server_group`
- `products.upstream_pid`
- `products.upstream_version`
- `products.location_version`
- `products.upstream_product_shopping_url`
- `products.upstream_qty`
- `products.rate`

这意味着：

- 本地商品是“复制品 + 映射关系”，不是只保留一个远端引用
- 后台产品编辑页的“同步商品”按钮，本质上是在重建本地售卖模型

## 7. 上下游商品开通的真实业务链

关键逻辑在 [Host.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\Host.php)。

当 `host.api_type == "zjmf_api"` 时，开通链是：

1. `GET /user_info`
2. 生成本地下游 `token`
3. 调 `/cart/clear`，尝试清理并初始化远端购物态
4. 如果有远端 `invoiceid`，先尝试 `/apply_credit`
5. 否则：
   - `/cart/add_to_shop`
   - `/cart/settle`
   - `/apply_credit`
   - 不足时再 `/apply_credit_limit`
6. 成功后把远端 `hostid` 回写到本地 `host.dcimid`

并且会把以下信息带给上游：

- `downstream_url`
- `downstream_token`
- `downstream_id`

说明这不是单向采购，而是：

- “本地下游站点”与“上游财务站点”之间有一条回调/识别链

这也是为什么老系统里 `stream_info`、`run_maping`、`dcimid` 这些字段都很关键。

## 8. 上下游链与魔方云链如何叠加

### 8.1 财务 API 链处理“卖什么、怎么买、怎么买通”

`zjmf_api` 负责：

- 登录对方财务站
- 拉商品
- 拉配置项
- 结算
- 扣余额/信用额
- 创建远端服务
- 续费/升级/工单等二次业务动作

### 8.2 魔方云链处理“创建出来以后怎么控”

`dcimcloud` 负责：

- 镜像
- 区域/线路
- 实例开关机
- VNC/控制台
- 重装
- 重置密码
- 救援模式
- 自定义按钮

也就是说，对一个 `dcimcloud` 商品来说：

- 购物和结算可能走 `zjmf_api`
- 实例动作最终走 `DcimCloud` / `Provision` / 远端 `provision/custom/{id}`

这就是老系统里“财务接口”和“云资源接口”并行存在的原因。

## 9. 资源池 / UpperReaches 是第三条链

和上面两条不同，`UpperReaches` 这条链解决的是：

- 手动资源池
- 代理资源
- 上游服务器库存
- IPMI / DCIM Client 控制

已确认相关表：

- `shd_upper_reaches`
- `shd_upper_reaches_res`

以及服务详情里附带：

- `upper_reaches_id`
- `upper_reaches_res`
- `upper_reaches_control_mode`

从代码引用看，`UpperReaches` 支持的动作至少包括：

- 电源状态查询
- 客户端动作按钮输出
- 重装
- 重置密码
- 取消重装
- 获取可重装系统

所以它不是“魔方云的别名”，而是第三条独立资源控制链。

## 10. 对重建新系统的直接要求

如果要做接近 1:1 的新系统，资源与商品架构至少要拆成下面 3 层：

### 10.1 财务与上下游层

- `finance_api_sources`
- `finance_api_products`
- `finance_api_product_sync_jobs`
- `client_api_keys`
- `client_api_logs`
- `free_product_grants`

### 10.2 商品售卖层

- `product_groups`
- `products`
- `pricings`
- `config_groups`
- `config_links`
- `config_options`
- `config_option_values`
- `custom_fields`

### 10.3 资源与动作层

- `providers`
- `provider_nodes`
- `services`
- `service_provider_bindings`
- `service_actions`
- `run_maps`
- `resource_pools`
- `resource_pool_items`

并且后端动作要明确分三类：

1. `finance upstream action`
2. `cloud/provider action`
3. `resource pool action`

不能再像很多丐版系统那样把它们混成一个“服务操作”按钮组。

## 11. 目前已经可以确定的页面与实现优先级

优先级最高的后台页：

- 上下游接口管理
- 商品同步与映射页
- 商品编辑页
- 服务详情页
- 运行日志/同步日志

优先级最高的用户中心页：

- API 管理页
- 服务详情页
- 购物车配置页
- 账单与充值页

优先级最高的后端模块：

- 上下游财务 API 适配器
- 商品同步器
- 下单/结算/远端开通链
- 服务动作门面
- 资源池动作门面

## 12. 当前最重要的结论

对你现在的新系统来说，真正该抄的不是“某个实例页长什么样”，而是这套对象关系：

`接口源 -> 商品同步 -> 本地商品 -> 订单/账单 -> 本地服务 -> 远端 hostid(dcimid) -> 服务动作`

并且服务动作背后至少要支持三类驱动：

- `zjmf_api`
- `dcimcloud`
- `upper_reaches`

只有这样，系统完整性才会接近老魔方，而不是只长得像。

## 13. 后台商品编辑页的真实字段结构

### 13.1 商品编辑页不是简单表单

从 [ProductController.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\admin\controller\ProductController.php) 的 `edit_product_page` 注释和查询字段可以确认，商品编辑页至少包含下面几层数据：

- 基础信息
- 定价与周期
- 自动开通与终止策略
- 模块/接口绑定
- 配置项与升级
- 自定义字段
- 上游同步字段
- 前台展示字段

后台读取商品时会一次性拉这些核心字段：

- `type`
- `gid`
- `groupid`
- `name`
- `description`
- `hidden`
- `welcome_email`
- `stock_control`
- `qty`
- `tax`
- `host`
- `password`
- `is_domain`
- `is_bind_phone`
- `pay_type`
- `pay_method`
- `allow_qty`
- `auto_setup`
- `api_type`
- `upstream_price_type`
- `upstream_price_value`
- `zjmf_api_id`
- `upstream_pid`
- `server_group`
- `auto_terminate_email`
- `config_options_upgrade`
- `upgrade_email`
- `affiliateonetime`
- `affiliate_pay_type`
- `affiliate_pay_amount`
- `retired`
- `is_featured`
- `config_option1..24`
- `is_truename`
- `clientscount`
- `product_shopping_url`
- `product_group_url`
- `upper_reaches_id`
- `cancel_control`
- `upstream_auto_setup`
- `upstream_ontrial_status`

### 13.2 和上游同步直接相关的字段

对接上下游后，商品编辑页最关键的映射字段是：

- `api_type`
- `zjmf_finance_api_id`
- `upstream_pid`
- `upstream_price_type`
- `upstream_price_value`
- `server_group`
- `upper_reaches_id`

其中逻辑分层是：

- `zjmf_finance_api_id`  
  对应哪一个上游财务接口源
- `upstream_pid`  
  对应上游商品 ID
- `upstream_price_type`  
  决定成本价同步后如何转成本地售价
- `upstream_price_value`  
  具体比例或方案值
- `server_group`
  在很多同步商品场景下直接等于 `zjmf_finance_api_id`
- `upper_reaches_id`
  手动/资源池型上游来源

### 13.3 后台“同步商品”接口做什么

后台商品页里有一个明确的同步入口：

- `/admin/product/sync_product_info`

要求参数：

- `pid`
- `upstream_pid`
- `zjmf_finance_api_id`
- `api_type`
- `rate`
- `upstream_price_type`
- `upstream_price_value`

这个接口不是只更新一个映射字段，而是触发完整商品同步。

### 13.4 `dcimcloud` 商品可自动生成配置项模板

这是一个特别关键的设计点。

在 `ProductController.php` 里，`dcimcloud` 商品支持按节点能力自动生成一套标准配置项模板。源码里能直接看到的模板项包括：

- `area|区域`
- `os|操作系统`
- `cpu|CPU`
- `memory|内存`
- `system_disk_size|系统盘`
- `network_type|网络类型`
- `bw|带宽`
- `in_bw|流入带宽`
- `ip_num|IP数量`
- `data_disk_size|数据盘`
- `snap_num|快照数量`
- `backup_num|备份数量`
- `nat_acl_limit|NAT转发`
- `nat_web_limit|共享建站`
- `system_disk_io_limit|系统盘性能`
- `data_disk_io_limit|数据盘性能`
- `traffic_bill_type|流量计费方式`

这说明：

- 老系统的“云商品”不是后台手填几十个字段
- 而是有一套标准商品模板生成器
- 生成器会再去节点取 `区域` 和 `操作系统` 等动态选项

### 13.5 `dcim` 裸金属也有另一套自动配置模板

同一个控制器里还内置了 `dcim` 商品模板生成逻辑，示例项包括：

- `area|区域`
- `os|操作系统`
- `server_group|硬件配置`
- `ip_num|IP数量`
- `bw|带宽`

这进一步说明：

- 商品编辑页在老系统里其实承担了一部分“商品建模器”角色
- 不同产品类型背后接的是不同模板生成器

## 14. 用户中心服务详情页的真实分流逻辑

### 14.1 服务详情不是统一接口渲染

关键逻辑在 [HostController.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\home\controller\HostController.php)。

服务详情页返回的数据，会按服务来源分三大类处理：

1. `api_type = zjmf_api`
2. `api_type = manual`
3. 其他本地/模块类服务，再按 `type` 判断 `dcimcloud / dcim / normal`

### 14.2 `zjmf_api` 服务详情怎么取数据

当服务是上下游财务 API 开出来的，详情页会：

- 调 `/host/header`
- 传 `host_id = dcimid`
- 从远端拿回：
  - `module_button`
  - `module_chart`
  - `module_power_status`
  - `show_traffic_usage`
  - `reinstall_random_port`
  - 流量包信息
  - 头部摘要信息

这意味着：

- 对于上游同步来的服务，详情页很多运行信息根本不靠本地库
- 而是实时以 `dcimid` 向上游拉取

### 14.3 `manual` / `UpperReaches` 服务详情怎么取数据

当服务是手动资源池链，详情页会读取：

- `products.upper_reaches_id`
- `upper_reaches_res.hid = host_id`

然后回填：

- `manual = { id, name }`
- `host_data.upper_reaches_res`
- `host_data.upper_reaches_control_mode`
- `module_power_status`
- `module_button`

这条链不依赖远端 `host/header`。

### 14.4 本地 `dcimcloud` 服务详情怎么取数据

如果不是 `zjmf_api`，而是本地 `dcimcloud` 商品，详情页直接走 `DcimCloud` 逻辑对象：

- `moduleClientButton(dcimid)`
- `moduleClientArea(host_id)`
- `chart(dcimid, host_id)`
- `getNatInfo(host_id)`
- `supportReinstallRandomPort(host_id)`

并额外回传：

- `dcimcloud.nat_acl`
- `dcimcloud.nat_web`

这说明 `dcimcloud` 详情页本身就比普通模块多两块：

- NAT 转发
- 共享建站

## 15. 服务详情页前端动作矩阵

### 15.1 前端动作总入口

高频动作集中在 [servicedetail.js](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\public\themes\clientarea\default\assets\js\servicedetail.js)。

核心通用入口：

- `service_module_button(_this, id, host_data_type)`

这套动作先取按钮上的：

- `func`
- `type`
- `desc`

然后决定：

- 弹窗
- 二次验证
- 调哪个接口

### 15.2 用户侧可见的高频动作

已确认的动作包括：

- `reinstall`
- `crack_pass`
- `rescue_system`
- 电源类动作
- 购买重装次数
- 获取图表
- 获取电源状态

### 15.3 重装的真实前端流程

`func = reinstall` 时，不是直接 POST。

流程是：

1. 先检查是否是 `dcimcloud/dcim`
2. 调 `/dcim/check_reinstall`
3. 根据返回：
   - 显示免费重装次数
   - 或提示是否购买重装次数
4. 用户确认后再提交 `/provision/default`
5. 若开启二次验证，还会弹二次验证框

这说明重装动作被拆成了两层：

- 前置资格检查
- 真正执行动作

### 15.4 改密的真实前端流程

`func = crack_pass` 时：

- 前端先按密码规则自动生成随机密码
- 打开改密弹窗
- 必要时触发二次验证
- 最终走 `/provision/default`

### 15.5 魔方云救援模式的特殊流程

`host_data_type == 'dcimcloud' && func == 'rescue_system'` 时：

- 会打开单独的 `moduleDcimCloudRescue` 弹窗
- 自动生成临时密码
- 某些场景下要求勾选强制关机确认
- 提交时也支持二次验证

这说明 `dcimcloud` 在前端不是普通按钮集合，而是单独有一套救援模式 UI。

## 16. 三条链的动作能力边界

### 16.1 `UpperReaches` 能力边界

从 [UpperReaches.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\UpperReaches.php) 可直接确认它支持：

- `getOs`
- `on`
- `off`
- `reboot`
- `vnc`
- `reinstall`
- `getStatus`
- `CrackPassword`
- `moduleClientButton`
- `moduleAdminButton`
- `modulePowerStatus`

并且它内部又再分两类控制模式：

1. `ipmi`
2. `dcimClient`

其中 `dcimClient` 还额外支持：

- `CancelReinstall`
- `ReinstallStatus`
- `GetOs`

### 16.2 `DcimCloud` 能力边界

从 [DcimCloud.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\DcimCloud.php) 可确认它的能力明显更重。

除了基础动作：

- `on`
- `off`
- `reboot`
- `hardOff`
- `hardReboot`
- `vnc`
- `reinstall`
- `sync`
- `suspend`
- `unsuspend`
- `terminate`
- `CrackPassword`
- `managePanel`
- `createAccount`

还包含云资源域能力：

- `listSnapBackup`
- `createSnap`
- `deleteSnap`
- `restoreSnap`
- `createBackup`
- `deleteBackup`
- `restoreBackup`
- `mountIso`
- `unmountIso`
- `setBootOrder`
- `createSecurityGroup`
- `showSecurityRules`
- `createSecurityRule`
- `delSecurityGroup`
- `delSecurityRule`
- `linkSecurityGroup`
- `addNatAcl`
- `delNatAcl`
- `addNatWeb`
- `delNatWeb`
- `exitRescue`
- `remoteInfo`
- `getTrafficUsage`
- `rescue`

这就是为什么魔方云链不能和 `UpperReaches` 混写成一组通用按钮。

### 16.3 `Provision` 的角色

[Provision.php](C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\Provision.php) 本质上是统一门面。

它负责：

- 模块发现
- 客户区输出
- 管理区输出
- 默认按钮
- 客户按钮
- 管理按钮
- 图表
- 自定义函数分发
- 资源下载
- 工单创建/回复

也就是说：

- `Provision` 是总线
- `DcimCloud` 是强能力云模块
- `UpperReaches` 是资源池控制模块
- `zjmf_api` 是财务/开通/上游同步链

## 17. 对重建新系统的进一步要求

如果要按老魔方 1:1 结构去重建，至少要在设计上分出这些后台页面：

- 财务接口源管理页
- 上游商品同步页
- 商品编辑工作台
- 资源池/代理资源页
- 魔方云节点与商品映射页
- 服务详情工作台
- 用户中心 API 管理页
- 用户中心服务详情页

并且页面和后端关系要做到：

- 商品编辑页不仅改本地字段，还能触发上游同步
- 服务详情页不是一套固定按钮，而是按 `api_type + type + control_mode` 组合渲染
- 用户中心 API 管理页必须是独立业务中心，不是附属开关
