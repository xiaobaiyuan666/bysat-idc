# 魔方财务深度拆解（二）

本文只做三件事：

1. 拆清 `下单 -> 账单 -> 支付 -> 开通 -> 上游 hostid 回写` 主链。
2. 拆清 `mf_cloud` v10 下单页的页面结构、字段矩阵、接口矩阵、状态回填。
3. 拆清用户中心 `servicedetail` 在 `zjmf_api / manual / dcimcloud / 其他模块` 下的动作差异。

说明：

- 本文以三层证据交叉校验：
  - 可读源码：`C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref`
  - 本地 3.7.4 模板与前端逻辑：`C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4`
  - 已验证的真实上下游接口行为
- `3.5.8 decoded` 负责解释逻辑。
- `3.7.4` 负责解释新版本页面和交互。
- 未直接在源码里证实、但由代码高概率推出的地方，会标注“推断”。

## 1. 主链总览

这套系统的主链不是“先开通再补账单”，而是严格的：

`购物参数归一化 -> 购物车会话 -> 结算校验 -> 发票/订单/服务实例落库 -> 余额/信用额/网关支付 -> 发票 Paid -> processPaidInvoice -> Host::create/unsuspend -> 模块/上游回写`

核心文件：

- `C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\Shop.php`
- `C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\home\controller\CartController.php`
- `C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\home\controller\PayController.php`
- `C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\Invoices.php`
- `C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\common\logic\Host.php`

## 2. 购物车入口链

### 2.1 `Shop::checkProductToArr`

位置：

- `Shop.php:80`

作用：

- 把页面输入归一化成标准购物项结构。
- 这一步还不会创建订单、账单、服务，只负责“把商品参数整理成合法购物车项”。

它做的校验和整理包括：

- 校验 `pid` 是否存在。
- 校验产品是否存在可售价格：
  - `ProductModel->checkProductPrice($pid, $billingcycle, $currencyid)`
- 如果页面显式传了 `serverid`，并且商品绑定了 `server_group`，校验该节点是否属于商品服务器组。
- 拉取当前商品绑定的全部配置项：
  - `product_config_options`
  - `product_config_links`
- 按配置项类型分别处理：
  - 数量型配置：按 `qty_minimum / qty_maximum` 修正
  - 单选/下拉型配置：若提交值不存在，则回落到第一个可用子项
  - 某些特殊选项类型允许空值
- 拉取商品自定义字段：
  - `customfields`
- 校验自定义字段：
  - 必填
  - 下拉值合法性
  - 正则表达式
- 记录以下基础业务参数：
  - `pid`
  - `billingcycle`
  - `configoptions`
  - `customfield`
  - `qty`
  - `allow_qty`
  - `pay_type`
  - `os`
  - `host`
  - `password`
  - `hostid`

返回结构：

- 成功：`status=success, data=<normalized cart item>`
- 失败：`status=error, msg=<具体错误>`

### 2.2 `CartController::addToShop`

位置：

- `CartController.php:2800`

对应路由：

- `POST /cart/add_to_shop`

这一步是购物页正式入购物车的入口。

它在调用 `Shop::addProduct` 之前，额外做了更强的业务校验：

- 购买前必须绑定手机号：
  - `buyProductMustBindPhone($uid)`
- 请求参数校验：
  - `pid`
  - `serverid`
  - `configoption`
  - `customfield`
  - `currencyid`
  - `qty`
  - `os`
  - `hostid`
  - `checkout`
- 若没传 `billingcycle`，自动用产品默认周期。
- 检查试用次数限制：
  - `judgeOntrialNum`
- 检查本地库存。
- 如果是 API 购买，还会检查：
  - 当前客户对该产品的允许持有数量
  - 对接上游试用数量限制
- 对 `api_type = zjmf_api / resource` 商品，还会远程检查上游库存：
  - `zjmfCurl(..., "cart/stock_control", ["pid" => upstream_pid], "GET")`
- 校验主机名规则。
- 校验密码规则。
- 若是开发者应用类商品，强制校验关联 `hostid`。

最后才调用：

- `Shop::addProduct(...)`

对应行为：

- `checkout = 0`：
  - 正常加购物车
  - 触发 `shopping_cart_add_product` hook
- `checkout = 1`：
  - 仍然先进购物车会话，但返回当前购物项位置 `i`
  - 供后续直接结算页使用

## 3. 结算落库主链

### 3.1 `CartController::settle`

位置：

- `CartController.php:1873`

对应路由：

- `POST /cart/settle`

这是最关键的总入口。

#### 3.1.1 输入来源

它支持三种输入来源：

- 当前用户的 `cart_session`
- `pos[]` 指定购物车内部分条目
- 直接传 `cart_data`，绕过既有购物车，临时结算单商品

#### 3.1.2 前置校验

这一步会做一大批“真正下单前”的校验：

- 购物车不能为空
- 客户持有数量限制
- 实名限制
- 手机绑定限制
- 支付网关是否可用
- 优惠码是否有效：
  - 开始时间
  - 过期时间
  - 最大使用次数
  - 仅新客 / 仅老客
  - 每客户仅一次
  - 依赖商品
- 商品是否停售
- 商品库存
- 商品价格周期是否有效
- 试用前置条件：
  - 实名
  - 邮箱
  - 手机
  - 微信

#### 3.1.3 对每个购物项的价格构造

结算时会把购物项先转成 `product_item`，再批量落库。

`product_item` 里包含：

- 商品基础信息
- 开通所需服务字段
- 发票项草稿
- 配置项快照
- 自定义字段快照
- 初装费、循环费、折扣、优惠码结果

处理细节包括：

- 按周期从 `pricing` 取：
  - `setup fee`
  - `recurring fee`
- 遍历配置项并叠加价格：
  - 数量型
  - 阶梯数量型
  - 普通单选型
- OS 类配置还会回写：
  - `os`
  - `os_url`
- 如果商品是 `zjmf_api` 且启用了上游版本和百分比加价：
  - 整个商品价和配置项价都会按 `upstream_price_value` 比例重算
- 如果商品是 `resource`：
  - 还会按资源用户等级折扣再折一次
- 然后叠加：
  - 用户等级折扣
  - 优惠码折扣

#### 3.1.4 发票、订单、服务对象的真实落库顺序

顺序非常明确：

1. 若本次需要付款，先创建 `invoices`
2. 再创建 `orders`
3. 再逐个商品、逐份数量创建 `host`
4. 再为每个 `host` 写：
   - `invoice_items`
   - `host_config_options`
   - `customfieldsvalues`
5. 最后更新：
   - 购物车
   - 优惠码使用次数
   - 导航分组显示状态
   - 发票跳转 URL

关键代码位置：

- `CartController.php:2511-2678`

细节：

- `invoices`：
  - `status = Unpaid`
  - `type = product`
- `orders`：
  - `status = Pending`
  - 绑定 `invoiceid`
- `host`：
  - `domainstatus = Pending`
  - 先存在本地，再等待支付后开通
- `invoice_items`：
  - 与 `hostid` 绑定
  - 一个服务实例对应一组发票项

#### 3.1.5 自动开通挂钩点

产品创建后不会立即一律开通，而是先看商品 `auto_setup`：

- `order`
  - 下单后立刻排队开通
- `payment`
  - 支付后再开通

代码位置：

- `CartController.php:2590-2594`
- `CartController.php:2740-2760`

#### 3.1.6 零元订单特殊路径

如果 `total == 0`：

- 有发票：
  - 直接把发票标记为 `Paid`
  - 然后调用 `Invoices::processPaidInvoice`
- 无发票：
  - 直接对 `create_after_pay` 里的服务实例发起开通

返回：

- `status = 1001`
- `data.hostid = [新建服务列表]`

## 4. 支付处理主链

### 4.1 余额支付 `PayController::applyCredit`

位置：

- `PayController.php:341`

对应路由：

- `POST /apply_credit`

它的职责不是“简单扣余额”，而是完整地：

- 校验发票是否属于当前用户
- 校验余额是否足够
- 计算当前发票已使用余额
- 计算本次可抵扣金额
- 更新 `invoices.credit / total / status / payment_status / paid_time`
- 更新客户余额
- 写 `credit_log`
- 若发票完全付清：
  - 调用 `Invoices::processPaidInvoice`
  - 返回 `status = 1001`
  - 返回当前发票里 `type=host` 的 `rel_id[]`

一个非常关键的细节：

- 如果这是上下游场景，并携带：
  - `downstream_url`
  - `downstream_token`
  - `downstream_id`
- 那它会把这些信息写入新 host 的 `stream_info`

这说明上下游开通链不是纯 API 调用，而是把“上下游回调关系”长期写在服务对象上。

### 4.2 信用额支付 `PayController::applyCreditLimit`

位置：

- `PayController.php:439`

对应路由：

- `POST /apply_credit_limit`

这条链和余额支付类似，但多了：

- 信用额功能是否开启
- 是否允许该客户使用信用额
- 当前已占用信用额
- 未结清信用额账单
- 信用额剩余额度

支付成功后同样会：

- 把发票标记 Paid
- 必要时扣部分余额
- 调用 `Invoices::processPaidInvoice`
- 返回 `status = 1001`

## 5. 发票 Paid 后处理

### 5.1 `Invoices::processPaidInvoice`

位置：

- `Invoices.php:314`

逻辑：

- 如果开了异步配置：
  - 推 `InvoicePaid` 队列
- 否则直接执行：
  - `processPaidInvoiceFinal`

### 5.2 `Invoices::processPaidInvoiceFinal`

位置：

- `Invoices.php:323`

这一层是“付款之后真正改变业务对象状态”的总分发器。

先做公共动作：

- 读取发票
- 强制要求发票已经是 `Paid`
- 生成 `invoice_num`
- 清 `is_delete`
- 拉取所有 `invoice_items`

然后按 `invoice_items.type` 分发：

#### `type = host`

做这些事：

- 更新 `host.regdate`
- 更新 `host.nextduedate / nextinvoicedate`
- 通知管理员
- 通知客户
- 发短信
- 如果商品开通条件是 `payment`：
  - 推 `AutoCreate`
- 如果服务当前是 `Suspended`：
  - 自动调用 `Host->unsuspend`
  - 并写 `RunMap`

#### `type = renew`

- 调 `Renew::renewHandle`

#### `type = recharge`

- 给客户 `credit` 增加余额
- 写余额日志

#### `type = upgrade`

- 调 `Upgrade::doUpgrade`

#### `type = credit_limit`

- 把预付款账单挂接到当前信用额账单
- 找出因信用额欠费暂停的服务并批量 `unsuspend`

#### `type = combine`

- 递归把被合并的原始账单也标记 Paid 并继续处理

#### 其他

- `voucher`
- `contract`
- `user_grade`

说明：

- 这就是魔方财务真正的业务中台核心。
- 支付成功本身不代表业务完成，真正的“开通/续费/升级/解暂停”都在这里分发。

## 6. 服务开通主链

### 6.1 `Host::create`

位置：

- `Host.php:21`

逻辑很薄：

- 如果开启自动开通队列：
  - 推 `AutoCreate`
- 否则直接走：
  - `createFinal`

### 6.2 `Host::createFinal`

位置：

- `Host.php:29`

这是最值得对标的服务开通总线。

它先做通用动作：

- 查询本地 `host + product + client`
- 用缓存锁避免重复开通：
  - `HOST_DEFAULT_ACTION_CREATE_{id}`
- 通过 `Provision->getParams` 生成模块参数
- 触发 `before_module_create` hook

然后按 `api_type` 分流。

#### 6.2.1 `api_type = zjmf_api`

这是最关键的上下游开通链。

精确顺序如下：

1. 远程调用 `/user_info`
   - 目的是拿上游用户和币种等信息，确认会话可用。

2. 生成本地下游 token
   - 写入本地 `host.stream_info.token`
   - 组合：
     - `downstream_url`
     - `downstream_token`
     - `downstream_id`

3. 先调用远程 `/cart/clear`
   - 不是简单清购物车。
   - 它更像“尝试确认上游是否已经为这个 downstream host 建好了资源”。

4. `/cart/clear` 返回两种常见情况：

- 情况 A：直接带 `hostid`
  - 说明上游已存在该实例
  - 直接回写本地：
    - `host.dcimid = hostid`

- 情况 B：带 `invoiceid`
  - 说明上游已有账单待支付
  - 直接尝试调用：
    - `/apply_credit`
  - 如果返回 `status = 1001`
    - 代表支付完成并返回 `hostid`
    - 回写本地 `dcimid`

- 情况 C：什么都没有
  - 进入完整重新下单链

5. 完整重新下单链：

- 准备远程购物参数：
  - `pid = upstream_pid`
  - `billingcycle`
  - `host`
  - `password`
  - `currencyid`
  - `qty = 1`
  - `configoption`
  - `customfield`

- 其中 `configoption` 来源于本地：
  - `host_config_options`
  - `product_config_options.upstream_id`
  - `product_config_options_sub.upstream_id`

- `customfield` 来源于本地：
  - `customfieldsvalues`
  - `customfields.upstream_id`

6. 远程调用 `/cart/add_to_shop`

7. 然后调用 `/cart/settle`

8. `/cart/settle` 返回后：

- 若已经直接返回 `hostid`
  - 回写本地 `dcimid`
- 若返回 `invoiceid`
  - 继续支付链

9. 支付链：

- 先 `/apply_credit`
- 若不足：
  - 再 `/apply_credit_limit`
- 若还不行：
  - 再回退 `/apply_credit use_credit=0`
  - 并返回失败

10. 一旦任一路径拿到上游 `hostid`

- 立即回写：
  - `host.dcimid`

这说明：

- 本地下游系统在开通上游商品时，并不是“调一个 create host 接口”。
- 它是把上游也当成一个完整财务系统来使用：
  - 上游购物车
  - 上游结算
  - 上游余额支付
  - 上游信用额支付
  - 上游 hostid 回写

#### 6.2.2 `api_type = resource`

- 调用 `resourceHostCreate`
- 成功则回写：
  - `dcimid`
  - `domainstatus`

#### 6.2.3 本地模块商品

- `type = dcim`
  - `Dcim->createAccount`
- `type = dcimcloud`
  - `DcimCloud->createAccount`
- 其他
  - `Provision->createAccount`

### 6.3 开通成功后的本地收尾

如果模块返回成功：

- 触发 `after_module_create`
- 本地非 `zjmf_api/resource` 商品会：
  - 把 `host.domainstatus` 更新为 `Active`
  - 发欢迎邮件
  - 发短信
- 记录活动日志
- 推送服务状态变化

如果失败：

- 触发 `after_module_create_failed`
- 记录失败日志
- 维持本地 Pending

## 7. `mf_cloud` v10 下单页接口矩阵

核心文件：

- `C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4\public\themes\cart\default\configureproduct-v10-cloud.tpl`
- `C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4\public\themes\cart\default\v10\api\mf_cloud.js`
- `C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4\public\themes\cart\default\v10\js\mf_cloud.js`

接口清单：

- `GET /product/common_cloud`
  - 商品列表/云产品入口
- `GET /v10/host/product/{id}/mf_cloud/order_page`
  - 下单页主配置数据
- `GET /v10/host/product/{id}/mf_cloud/image`
  - 镜像树
- `POST /v10/host/product/{id}/mf_cloud/duration`
  - 当前配置可选周期
- `POST /product/{id}/config_option`
  - 实时计算价格与预览
- `POST /product/settle`
  - 直接去结算页
- `POST /promo_code/apply`
  - 应用优惠码
- `GET /cart`
  - 购物车概览
- `GET /product/{id}/v10/host/mf_cloud/line/{line_id}`
  - 某线路的详细带宽/流量/IP/防御规则
- `GET /ssh_key`
  - SSH key 列表
- `GET /security_group`
  - 安全组列表
- `GET /v10/host/product/{id}/mf_cloud/vpc_network/search`
  - VPC 列表
- `GET /client_level/product/{id}/amount`
  - 用户等级折扣
- `POST /login/v10/auth`
  - 下单页轻登录状态
- `POST /cart/v10/add`
  - 加购物车
- `GET /cart/v10/edit`
  - 编辑购物车回填
- `POST /cart/v10/edit`
  - 更新购物车
- `GET /v10/host/product/{id}/self_defined_field/order_page`
  - 自定义字段

## 8. `mf_cloud` v10 页面结构矩阵

### 8.1 顶层结构

页面有两个顶层模式：

- `fast`
  - 推荐套餐
- `custom`
  - 自定义配置

状态变量：

- `activeName = fast | custom`
- `showFast`
- `isUpdate`
- `isCustom`

### 8.2 页面状态对象

主状态在 `mf_cloud.js:800+`。

最核心的是 `params`：

- `data_center_id`
- `cpu`
- `memory`
- `image_id`
- `system_disk.size`
- `system_disk.disk_type`
- `data_disk[]`
- `backup_num`
- `snap_num`
- `line_id`
- `bw`
- `flow`
- `peak_defence`
- `ip_num`
- `ipv6_num`
- `duration_id`
- `network_type`
- `name`
- `ssh_key_id`
- `security_group_id`
- `security_group_protocol[]`
- `password`
- `re_password`
- `vpc.id`
- `vpc.ips`
- `port`
- `notes`
- `auto_renew`
- `resource_package_id`
- `ip_mac_bind_enable`
- `nat_acl_limit_enable`
- `nat_web_limit_enable`
- `ipv6_num_enable`

其他关键状态：

- 地域层：
  - `country`
  - `countryName`
  - `city`
  - `area_name`
- 套餐层：
  - `recommendList`
  - `packageId`
  - `cloudIndex`
- 镜像层：
  - `imageList`
  - `curImage`
  - `curImageId`
  - `imageName`
  - `version`
- 网络层：
  - `lineList`
  - `lineDetail`
  - `lineType`
  - `netName`
  - `vpcList`
  - `vpc_ips`
- 登录层：
  - `login_way`
  - `sshList`
- 定价层：
  - `cycleList`
  - `totalPrice`
  - `preview`
  - `discount`
  - `duration`
- 回填层：
  - `position`
  - `backfill`

### 8.3 `fast` 推荐套餐模式

模板位置：

- `configureproduct-v10-cloud.tpl:21-324`

主要页面块：

- 国家/城市选择
- 推荐套餐卡片
- 镜像选择
- 网络类型切换
- VPC 选择或新建
- 登录方式
- 备注
- IP/MAC 绑定
- NAT 转发
- NAT 建站
- IPv6 开关
- 自定义字段

推荐套餐卡片展示的内容已经很完整：

- 名称
- CPU
- 内存
- 系统盘容量/类型
- 数据盘容量/类型
- 带宽或流量
- 免费盘
- 防御值
- IPv4 数量
- IPv6 数量
- GPU 数量与型号
- 是否不支持升级

`handlerFast()` 的真实逻辑：

- 从 `country -> city -> area -> recommend_config` 拉扁平推荐套餐
- 初始化：
  - `packageId`
  - `data_center_id`
  - `cpu`
  - `gpu_num`
  - `memory`
  - `line_id`
  - `bw / flow`
  - `ipv6_num`
  - `peak_defence`
  - `system_disk`
  - `data_disk`
- 再调用 `getCycleList()`

### 8.4 `custom` 自定义模式

模板位置：

- `configureproduct-v10-cloud.tpl:329-818`

页面块比 `fast` 更重：

- 资源包
- 国家/城市
- 可用区
- CPU
- 内存
- GPU
- 镜像
- 存储表格
  - 系统盘类型
  - 系统盘容量
  - 数据盘类型
  - 数据盘容量
  - 增删数据盘
- 网络类型
- VPC
- 线路
- 带宽或流量
- IPv4
- IPv6
- 防御
- 登录方式
- SSH key
- 密码规则校验
- SSH 端口
- 自动续费
- IP/MAC 绑定
- NAT
- 自定义字段

这说明 `mf_cloud` 商品页不是“简单配置项表单”，而是一套专门的云资源规划器。

### 8.5 网络与 VPC 规划

网络类型：

- `normal`
- `vpc`

当 `network_type = vpc` 时：

- 可选已有 VPC
- 也可新建 VPC

新建 VPC 有两种方式：

- 自动规划 `plan_way = 0`
- 手动规划 `plan_way = 1`

手动规划字段：

- 第一段网段前缀
- 第二段
- 第三段
- 第四段
- 掩码长度

### 8.6 登录方式与密码规则

登录方式变量：

- `login_way`

可见选项：

- 自动生成密码
- 手动设置密码
- 代码里保留了 SSH key 逻辑，但模板部分版本里被注释或按条件隐藏

校验逻辑：

- 长度
- 禁止非法字符
- 首字符限制
- 必须包含大小写和数字
- 确认密码一致
- SSH key 必填

### 8.7 状态回填与编辑模式

这页不是只支持新购，也支持：

- 从购物车回填编辑
- 从 `sessionStorage.product_information` 回填

它会额外保存一批“只为回显存在”的字段：

- `activeName`
- `country`
- `countryName`
- `city`
- `curImage`
- `curImageId`
- `imageName`
- `version`
- `cloudIndex`
- `login_way`
- `groupName`

也就是说，购物车里存的不只是交易必要字段，还存了页面 UI 状态。

### 8.8 提交链

#### 立即购买

- `submitOrder()`
- 先做：
  - 表单校验
  - `formatData()`
  - 自定义字段校验
- 然后把完整配置写到：
  - `sessionStorage.product_information`
- 最后跳：
  - `settlement.htm?id={product_id}`

#### 加购物车

- `addCart()`
- 把 `params + UI 回填字段 + 自定义字段` 打包
- `formatSwitch()` 把布尔开关转成 `0/1`
- POST `/cart/v10/add`
- 成功后跳转购物车

#### 编辑购物车

- `changeCart()`
- 与加购物车类似，但会带：
  - `position`
- POST `/cart/v10/edit`

#### 价格实时计算

- `changeConfig()`
- POST `/product/{id}/config_option`
- 返回：
  - `price`
  - `preview`
  - `duration`

## 9. 用户中心服务详情分支矩阵

核心文件：

- `C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\home\controller\HostController.php`
- `C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\public\themes\clientarea\default\assets\js\servicedetail.js`

### 9.1 分支条件

服务详情不是一套统一页，而是按两层分支：

- 第一层：`api_type`
  - `zjmf_api`
  - `manual`
  - `resource`
  - `whmcs`
  - 其他本地模块
- 第二层：`type`
  - `dcimcloud`
  - `dcim`
  - 其他模块

### 9.2 `api_type = zjmf_api`

核心行为：

- 服务详情会去上游拉：
  - `/host/header?host_id=dcimid`

上游返回的数据直接用于本地页面：

- `module_button`
- `module_chart`
- `module_power_status`
- `reinstall_random_port`
- `dcim.flow_packet_use_list`
- `dcimcloud.nat_acl`
- `dcimcloud.nat_web`

结论：

- 对于 `zjmf_api` 服务，本地详情页是一个“上游页面聚合器”。
- 本地不只存状态，也实时吃上游按钮和展示数据。

### 9.3 `api_type = manual`

核心行为：

- 通过 `UpperReaches` 取控制能力：
  - `modulePowerStatus`
  - `moduleClientButton`

还会返回：

- `manual.id`
- `manual.name`

结论：

- `manual` 是资源池/手动资源页面，不是云实例页面。

### 9.4 本地 `dcimcloud`

核心行为：

- 由 `DcimCloud` 直接提供：
  - `moduleClientButton`
  - `moduleClientArea`
  - `chart`
  - `getNatInfo`
  - `supportReinstallRandomPort`

返回扩展字段：

- `dcimcloud.nat_acl`
- `dcimcloud.nat_web`
- `module_power_status = true`
- `reinstall_random_port = true/false`

这说明 `dcimcloud` 本地商品在用户中心其实有专属 UI 区域，而不是单纯一个按钮栏。

### 9.5 其他模块

- `dcim`
  - 调 `Dcim->moduleClientButton / moduleClientArea`
- 普通模块
  - 调 `Provision->chart`
  - `Provision->checkDefineFunc($host_id, "Status")`

## 10. 用户中心动作矩阵

### 10.1 统一入口

前端入口：

- `service_module_button(_this, id, host_data_type)`

每个按钮通过 `data-*` 传：

- `func`
- `type`
- `desc`

### 10.2 `reinstall`

真实流程：

1. 若是 `dcim` / `dcimcloud`
   - 先请求 `/dcim/check_reinstall`
2. 展示本周免费重装次数、已用次数、剩余次数
3. 若免费次数不足
   - 可走购买重装次数
4. 打开重装弹窗
5. 校验：
   - 已备份确认
   - 密码规则
   - 端口规则
6. 如启用了二次验证
   - 先弹二次验证
7. 最终提交到：
   - `/provision/default`
   - 或 `dcim` 对应专用接口

### 10.3 `crack_pass`

真实流程：

1. 按密码规则自动生成新密码
2. 打开改密弹窗
3. 校验密码
4. 某些场景需要勾选强制关机
5. 如启用二次验证，先走二次验证
6. 最终提交

### 10.4 `rescue_system` for `dcimcloud`

这是最特殊的一条。

专属弹窗：

- `moduleDcimCloudRescue`

流程：

1. 自动生成临时密码
2. 根据当前状态决定是否显示“强制关机确认”
3. 校验密码规则
4. 可能触发二次验证
5. 提交救援动作
6. 刷新电源状态轮询

### 10.5 `vnc`

`dcim` 路线里明确存在：

- 请求 `/dcim/novnc`
- 返回 `password + url`
- 浏览器新开窗口进入 noVNC

这证明老系统前端并不是自己实现 VNC 协议，而是“后端生成控制台会话 -> 前端打开 noVNC 页”。

### 10.6 二次验证机制

多类动作共享这套机制：

- `reinstall`
- `crack_pass`
- `vnc`
- `rescue`

前端会先判断：

- `isNeedSecond(func)`

若需要：

- 先走 `getSecondModal(func, callback)`
- 拿到验证码后再真正提交动作

## 11. 当前已经可以确定的重建要求

基于本轮拆解，后续重建系统时必须满足：

- 商品模型必须同时承载：
  - 售卖属性
  - 开通属性
  - 上游映射属性
  - 云资源模板属性
- 订单链必须严格分层：
  - `cart -> order -> invoice -> invoice_items -> host`
- 支付成功后必须统一进入一个 `invoice paid dispatcher`
- `zjmf_api` 服务开通必须按“上游购物车/结算/余额/信用额”的真实顺序建模
- `mf_cloud` 商品页必须独立于普通产品页，不能退化成“配置项表格”
- 服务详情必须按：
  - `api_type`
  - `type`
  - `control_mode`
  组合渲染

## 12. 下一层继续拆的内容

下一轮继续拆四块：

1. 后台商品编辑页的“同步商品/同步配置项/同步价格/同步镜像模板”按钮链。
2. 后台魔方云商品映射页与 `dcimcloud` 商品构建逻辑。
3. 用户中心 `apimanage` 页面和 `api_user_product` 免费商品发放逻辑。
4. 后台订单详情、账单详情、服务详情三张工作台的对象关系。

## 13. 上下游 API 管理中心的页面结构

这一组页面很容易被误判成“只是一个 API 开关”，但实际它是老系统里独立的一条业务面：

- 上游财务接口源管理
- 下游客户 API 开通与锁定
- 免费产品发放
- API 日志
- 全局 API 条件设置

### 13.1 后台入口方式

后台不是一堆传统多页 PHP，而是：

- 路由壳：
  - `C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\data\route\admin.php`
- 入口视图控制器：
  - `C:\Users\Administrator\Desktop\IDC\zjmf-manger-decoded-ref\app\admin\controller\ViewResourceController.php`
- 后台 SPA 再分出子路由：
  - `zjmf-api`
  - `add-supplier`
  - `customer-view/api-overview`
  - `customer-view/log`
  - `api-setup`

这说明后台“上下游/API 管理”在架构上本来就是一个前端壳下的子系统。

### 13.2 后台页面树

#### 上游 API / 供应商列表页

能力：

- 上游接口源列表
- 新增上游接口源
- 编辑上游接口源
- 删除
- 刷新状态
- 查看汇总
- 查看授信/余额

主要控制器：

- `app/admin/controller/ZjmfFinanceApiController.php`

主要页面包：

- `public/admin/js/ZjmfApi~31ecd969.a1bf15c6.js`
- `public/admin/js/addSupplier~f71cff67.cbc03931.js`

#### 客户 API 概要页

这页不是客户详情里的普通 tab，而是专门的 API 业务工作台。

能力：

- API 开通
- 锁定 / 解锁
- 重置 API 密钥
- 查看 API 概要
- 管理免费产品

主要页面包：

- `public/admin/js/CustomerAdd~31ecd969.16507cba.js`

#### 客户 API 日志页

这页支持按来源切换：

- `WEB`
- `API`

切到 `API` 来源后，后台会调用专门的 API 日志接口，而不是普通站内日志。

主要页面包：

- `public/admin/js/CustomerLog~31ecd969.c8516017.js`

#### API 全局设置页

这一页不在 `ZjmfFinanceApiController`，而在系统配置控制器。

作用：

- API 总开关
- 实名前置条件
- 绑手机前置条件

主要控制器：

- `app/admin/controller/ConfigGeneralController.php`

### 13.3 用户中心页面树

#### `/apimanage`

这页是用户中心自己的 API 业务中心，不是简单展示页。

页面控制器：

- `app/home/controller/ViewClientsController.php`

模板：

- `public/themes/clientarea/default/apimanage.tpl`

功能：

- 未开通态：
  - 勾选协议
  - 立即开通
  - 提示实名/手机条件
- 已开通态：
  - 展示 API key
  - 展示开通时间
  - 展示产品数 / 代理数 / 今日请求数
  - 展示 7 日请求图
  - 展示免费产品列表
  - 重置密钥
  - 关闭 API

#### `/apilog`

页面控制器：

- `app/home/controller/ViewClientsController.php`

模板：

- `public/themes/clientarea/default/apilog.tpl`

功能：

- API 请求日志查看
- 不直接在模板里调接口，而是页面控制器内部再调用：
  - `app/home/controller/ZjmfFinanceApiController.php`

### 13.4 关键动作和控制器

#### 开通 API

- 后台：
  - `admin/ZjmfFinanceApiController`
- 用户中心：
  - `home/ZjmfFinanceApiController`

#### 重置 API 密钥

- 后台：
  - 后台客户 API 概要页触发
- 用户中心：
  - `apimanage` 页面触发

#### API 日志

- 后台：
  - 客户日志页切换到 `API` source
- 用户中心：
  - `apilog`

#### 免费产品

- 后台可编辑：
  - 新增
  - 删除
  - 修改
- 用户中心只展示

### 13.5 对重建的直接要求

如果要 1:1 重建，必须把“API 管理中心”当成独立业务域，而不是只在客户详情里塞两个按钮。

至少要重建这些页面：

- 后台：
  - 上游接口源列表
  - 新增/编辑接口源
  - 客户 API 概要工作台
  - 客户 API 日志工作台
  - API 全局设置页
- 用户中心：
  - API 管理首页
  - API 日志页
  - 免费产品列表
