# 魔方财务 3.x 活体拆解记录（2026-03-21）

本文不是泛分析，而是把三套证据对齐后的“活体结果”落档：

1. 本地跑通的 3.7.4 老财务系统
2. 可读源码仓库 `zjmf-manger-decoded-ref`（3.5.8）
3. 真实上下游财务接口与真实云主机回包

目标是继续逼近“1:1 重建”，尤其是这四条链：

- 后台商品架构
- 上下游财务 API 对接
- 用户中心服务详情
- 魔方云 `dcimcloud` 商品购买与实例动作

## 1. 本次实测样本

### 1.1 本地老系统

- 前台地址：`http://127.0.0.1:8091/`
- 后台地址：`http://127.0.0.1:8091/admin01/`
- 数据库：`zjmf_legacy374`

### 1.2 真实上下游财务接口

- 站点：`https://www.xiaobaiyun.cn`
- 已验证接口：
  - `POST /zjmf_api_login`
  - `GET /host/list`
  - `GET /host/header`
  - `GET /cart/all`
  - `GET /cart/get_product_config`

### 1.3 本地演示用户与演示服务

为了把用户中心完整跑起来，我在本地造了一套最小闭环数据：

- 演示用户：
  - 邮箱：`legacy-demo@example.com`
  - 密码：`Legacy123!`
- 本地服务：
  - `host.id = 1`
  - `host.type = dcimcloud`
  - `host.api_type = zjmf_api`
  - `host.productid = 3`
  - `host.dcimid = 8346`

这里最关键的点是：

- 本地 `host.dcimid = 8346`
- 对应的是“上游财务 host id”
- 不是底层云实例 id `7469`

这是老魔方最容易搞错的地方。  
对于 `api_type = zjmf_api` 的服务，后续 `host/header`、`host/details` 这类接口，调用上游财务站时传的是：

- `host_id = 本地 host.dcimid`

也就是：

- `8346`

不是底层 provider instance id：

- `7469`

## 2. 用户登录链路

### 2.1 登录页不是单表单

本地前台登录页默认激活的是：

- 手机登录 tab

邮箱登录表单是隐藏的，所以之前如果直接填：

- `#emailInp`
- `#emailPwdInp`

浏览器自动化会失败，因为表单未显示。

登录页 HTML 已保存：

- `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-login-page.html`

### 2.2 密码不是明文提交

登录页提交前会调用：

- `encryptPass(id)`

源码位置：

- `C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4\public\themes\clientarea\default\assets\js\public.js`

核心逻辑：

- AES-CBC
- key：`idcsmart.finance`
- iv：`9311019310287172`

也就是说：

- 前台邮箱/手机登录
- 不是直接把明文密码 POST 给 `/login`
- 而是前端先 AES 加密，再提交

### 2.3 本地真实登录已打通

我已经用真实登录链跑通了演示用户登录，并成功访问：

- `/clientarea`
- `/servicedetail?id=1`

生成的活体页面文件：

- `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-clientarea-live.png`
- `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-clientarea-live.html`
- `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-servicedetail-live.png`
- `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-servicedetail-live.html`
- `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-servicedetail-live.txt`

## 3. 用户中心首页结构

实测截图：

- `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-clientarea-live.png`

页面结构可以确认是：

### 3.1 左侧一级导航

- 用户中心
- 产品与服务
- 账户管理
- 财务管理
- 技术支持

这说明老魔方用户中心不是零散页面，而是稳定的“账户控制台”。

### 3.2 首页内容块

当前本地演示首页已经能看到这些区块：

- 用户信息卡
- 账户余额卡
- 已出产品数
- 我的产品表格
- 公告通知区

其中“我的产品”表格列已经实测存在：

- 产品名称
- 到期时间
- IP

这说明用户中心首页不是统计页，而是一个“业务入口 + 当前资源概览页”。

## 4. 本地服务列表接口：`/host/list`

本地真实回包已保存：

- `C:\Users\Administrator\Desktop\IDC\output\local-user\host-list.json`

关键字段已经确认：

- `id`
- `api_type`
- `zjmf_api_id`
- `domain`
- `domainstatus`
- `regdate`
- `dedicatedip`
- `assignedips`
- `nextduedate`
- `firstpaymentamount`
- `amount`
- `billingcycle`
- `os`
- `dcimid`
- `productname`
- `product_type`
- `pid`
- `invoice_id`

同时还会把高频配置项扁平化回列表：

- `node`
- `cpu`
- `memory`
- `network_type`
- `bw`
- `ip_num`
- `data_disk_size`
- `port`
- `in_bw`

这说明用户中心的“产品列表”不是只显示订单数据，而是：

- 服务基本信息
- 财务状态
- 到期信息
- 关键配置摘要

都混合在一条 `host/list` 回包里。

## 5. `host/header` 与 `host/details` 的真实区别

本地真实回包：

- `C:\Users\Administrator\Desktop\IDC\output\local-user\host-header-1.json`
- `C:\Users\Administrator\Desktop\IDC\output\local-user\host-details-1.json`

上游真实回包：

- `C:\Users\Administrator\Desktop\IDC\output\upstream\host-header-8346.json`

### 5.1 `host/header` 是当前服务详情页主数据源

本地 `servicedetail?id=1` 对应页面真正依赖的是：

- `host/header`

不是 `host/details`。

这一点已经被三层证据同时证实：

- 源码 `ViewClientsController::servicedetail()`
- `zjmfcloud.tpl`
- 本地活体页面

### 5.2 `host/header` 比 `host/details` 多出来的关键能力

`host/header` 会返回：

- `module_button.control`
- `module_button.console`
- `module_client_area`
- `module_chart`
- `module_power_status`
- `password_rule`
- `second`
- `system_button`
- `cancelist`
- `dcimcloud`
- `cloud_os`
- `cloud_os_group`

而本地 `host/details` 当前回包里：

- `module_button.control = []`
- `module_button.console = []`
- `module_client_area = []`

所以：

- `host/details` 更像轻量详情接口
- `host/header` 才是当前服务详情工作台的真实合同

### 5.3 本地 `host/header` 对上游数据做了二次加工

上游真实 `host/header(8346)` 返回的 `module_client_area` 包含：

- `snapshot`
- `security_groups`
- `setting`

但本地 `host/header(1)` 实际只保留：

- `snapshot`
- `setting`

`security_groups` 被过滤掉了。  
这和源码里对上游 `module_client_area` 的筛选逻辑一致。

另外一个明显差异：

- 上游 `host_data.password` 有值
- 本地 `host_data.password` 为空字符串

说明本地服务详情页不会原样暴露上游 root 密码，而是额外走：

- “重置密码”
- “登录信息”

这类控制入口。

## 6. `zjmfcloud` 服务详情页的真实页面结构

实测截图：

- `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-servicedetail-live.png`

### 6.1 顶部资源头卡

头卡信息已经实测包含：

- 产品名称
- 服务状态
- 机器名/域名
- 操作系统
- CPU
- 内存
- 主 IP

右侧动作区实测存在：

- `VNC`
- 电源动作下拉
- `重装系统`

### 6.2 左侧基础信息卡

当前 live 页面已经确认有这些卡：

- 登录信息
  - 用户名
  - 密码占位
  - 重置密码入口
  - 端口
- 价格与续费
  - 首购价格
  - 自动余额续费开关
  - 续费
  - 申请停用
  - 订购日期
  - 支付周期
  - 到期时间
- 配置摘要
  - 区域
  - 网络类型
  - 带宽
  - IP 数量
  - 查看更多信息
- 备注信息

### 6.3 右侧 tab

本地 live 页面已确认 tab 至少包括：

- 图表
- 产品升降级
- 快照/备份
- 设置
- 日志
- 财务

这说明服务详情不是单对象展示页，而是服务工作台。

## 7. 服务详情页的真实异步子页

### 7.1 NAT 子页

本地抓取：

- `C:\Users\Administrator\Desktop\IDC\output\local-user\service-nat-1.html`

当前 NAT 为空，只返回站点尾部标记，说明：

- 当前样本机器没有 NAT 规则
- 前台依然会尝试异步请求 NAT 区块

对应请求：

- `/servicedetail?action=nat&id=1`

### 7.2 快照/备份子页

本地抓取：

- `C:\Users\Administrator\Desktop\IDC\output\local-user\provision-custom-snapshot-1.html`

这个片段非常关键，因为它证明：

- `snapshot` tab 不是本地静态模板
- 而是通过 `/provision/custom/content?id=1&key=snapshot` 动态返回 HTML 片段

并且这个 HTML 片段本身还带完整 JS 行为。

已确认动作包括：

- `createSnap`
- `deleteSnap`
- `restoreSnap`
- `createBackup`
- `deleteBackup`
- `restoreBackup`

这些动作最终都 POST 到：

- `/provision/custom/1`

也就是说：

- 前台 tab 内容是 HTML fragment
- fragment 内部自己再绑定 provider 动作

### 7.3 设置子页

本地抓取：

- `C:\Users\Administrator\Desktop\IDC\output\local-user\provision-custom-setting-1.html`

已经确认这里不是“账号设置”，而是云主机设备设置，包含：

- 挂载 ISO
- 启动项（Boot Order）

对应动作：

- `mountIso`
- `unmountIso`
- `setBootOrder`

同样最终都 POST 到：

- `/provision/custom/1`

### 7.4 `remoteInfo` 探测

本地抓取：

- `C:\Users\Administrator\Desktop\IDC\output\local-user\provision-custom-remoteInfo-1.json`

当前返回：

```json
{"status":200,"msg":"获取成功","data":{"rescue":0}}
```

这个接口的作用不是渲染内容，而是：

- 探测当前实例是否处于救援模式
- 动态决定页面上哪些按钮显示/隐藏

## 8. 动作层真实结构

当前服务详情页上的实例动作，不是一个接口一个页面按钮直连，而是三层结构：

### 8.1 页面按钮层

按钮来源：

- `module_button.control`
- `module_button.console`

当前样本实际动作：

- `on`
- `off`
- `reboot`
- `hard_off`
- `hard_reboot`
- `reinstall`
- `crack_pass`
- `rescue_system`
- `exitRescue`
- `vnc`

### 8.2 通用动作提交层

大多数动作最终走：

- `/provision/default`

而模块自定义页签/子功能更多走：

- `/provision/custom/{host_id}`

### 8.3 页面内动作分流

前端不是所有按钮都直接提交。  
有些动作先打开模态框，再二次提交：

- `reinstall`
- `crack_pass`
- `rescue_system`

这类动作体现的是：

- Provider 能力统一
- 页面交互分层

而不是简单 REST 按钮。

## 9. 本地导入商品 3 的购买页结构

我把本地导入的上游云商品 `pid=3` 购买页跑出来了：

- `C:\Users\Administrator\Desktop\IDC\output\local-user\cart-configureproduct-pid3.html`

### 9.1 页面不是 v10 云规划器

这个页面是老版 `configureproduct.tpl` 结构，不是 v10 `mf_cloud` 大规划器。

但它仍然完整反映了：

- 同步后的本地配置项树
- 周期选择
- 服务端订单汇总

### 9.2 已确认的本地配置项

当前商品 3 已同步出的配置项有：

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
- 周期

这说明本地商品同步成功后，用户购买页的配置结构完全来自本地：

- `product_config_groups`
- `product_config_options`
- `product_config_options_sub`

不是前台实时直接去上游拿表单。

### 9.3 操作系统不是单纯下拉

操作系统这块实际是“两段式”：

1. OS 分组下拉
2. 分组下的具体镜像下拉

这和前台 JS 完全吻合：

- `configoption_os_group`
- `configoption_os`
- `osGroupChange()`
- `configoption_os()`

### 9.4 带宽是滑块 + 数值输入

当前商品 3 的带宽项不是普通单选，而是：

- range slider
- 数值输入框
- 阶梯区间价格

已确认：

- 最小 `20`
- 最大 `200`
- 按 `qty_stage=1`
- 对应区间价通过 `data-sub` 挂在输入上

这说明老魔方本地配置项体系本身已经能表达：

- 区间价格
- 阶梯价格
- 数量范围

### 9.5 右侧价格不是前端静态算

购买页右侧汇总来自服务端 HTML 回包，不是浏览器纯前端计算。

真实请求：

- `POST /cart?action=ordersummary`

本地真实回包：

- `C:\Users\Administrator\Desktop\IDC\output\local-user\cart-ordersummary-pid3-default.html`

当前默认配置实际价格结构已经确认：

- 商品基础价格：`￥20`
- CPU 4核：`￥5`
- IP 数量 1个：`￥70`
- 总价：`￥95`

这说明：

- 周期切换
- 配置变更
- 折扣/价格汇总

最终都回到服务端重新计算，再把 HTML 片段返回前台。

## 10. 后台商品同步链的真实结论

结合 decoded 源码和后台页面分析，可以确认：

### 10.1 商品创建和商品同步不是一回事

`ProductController::create()` 只是：

- 先建本地商品壳
- 初始化价格记录
- 建导航绑定

它不会在创建时就完成：

- `api_type = zjmf_api`
- `upstream_pid`
- `zjmf_api_id`
- `upper_reaches_id`

这些真正发生在：

- 后台商品编辑页保存
- “同步商品信息”链路

### 10.2 同步商品不是轻更新，而是破坏性重建

`ProductModel::syncProductInfo()` 的本质是：

- 用上游商品定义重建本地商品销售模型

它会重建这些对象：

- `products`
- `pricing`
- `customfields`
- `product_config_groups`
- `product_config_links`
- `product_config_options`
- `product_config_options_sub`
- `pricing(type=configoptions)`

同时它还会把历史服务的：

- `host_config_options`

按照 `upstream_id` 映射回新本地配置项 id。

这说明老魔方“同步商品”不是简单字段回填，而是一个完整商品编译器。

### 10.3 自动化接口和配置项模板是强耦合的

后台商品编辑页里：

- `api_type`
- `server_group`
- `upstream_pid`
- `zjmf_api_id`
- `upper_reaches_id`

不只是决定“用哪个 provider”，还会反向决定：

- 生成什么本地配置项模板
- 购物页渲染什么配置项
- 服务实例最终保存什么配置快照

这正是“商品架构”与“资源对接方式”强耦合的根源。

### 10.4 本地后台 `edit_product_page/3` 活体结果

我已经直接调用了本地后台真实接口：

- `GET /admin01/edit_product_page/3`

回包文件：

- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\edit-product-page-3.json`

它一次性返回了商品编辑页的大部分核心数据：

- `product`
- `pricing`
- `modules`
- `server_group`
- `config_groups`
- `config_links`
- `customfields`
- `all_product_data`
- `upgrade_product_ids`
- `upstream_product_pricings`

当前商品 3 的关键绑定已经在 live 回包里得到证实：

- `type = dcimcloud`
- `api_type = zjmf_api`
- `zjmf_api_id = 1`
- `upstream_pid = 96`
- `server_group = 1`
- `upstream_price_type = percent`
- `upstream_price_value = 100.00`
- `auto_setup = payment`
- `config_options_upgrade = 1`

同时还能看到：

- 本地商品价格矩阵 `pricing`
- 本地绑定的配置组 `config_links = [3]`
- 上游商品价格矩阵 `upstream_product_pricings`
- 服务器组选项里 `zjmf_api` 与 `normal` 是分开展示的

这说明后台商品编辑页在设计上就是一个“聚合工作台”，不是 tab 各自零散请求。

### 10.5 `get_upstream_products` 返回的是“导航分组 + 商品列表”

我已实测本地后台真实接口：

- `GET /admin01/get_upstream_products?id=1`

回包文件：

- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\get-upstream-products-1.json`

这个接口不是扁平商品数组，而是：

- 上游商品分组
- 分组下商品列表

也就是：

- `[{ id, name, products: [...] }]`

这和真实上游 `/cart/all` 的组织方式一致。  
所以后台“选择上游商品”弹窗的数据源，本质上复用的是上游商城导航结构，不是后台自己再造一套商品树。

### 10.6 `get_upstream_price` 返回的是完整配置树，不只是价格

我已实测本地后台真实接口：

- `GET /admin01/product/get_upstream_price?pid=3`

回包文件：

- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\get-upstream-price-pid3.json`

这个接口返回的不是单独价格，而是：

- `product_pricing`
- `config_groups`
- `flag`

其中 `config_groups` 里已经包含：

- 配置组
- 配置项
- 子项
- 每个子项的 `pricing(type=configoptions)`

也就是说，后台“查看上游价格/配置”不是只看周期价，而是看整棵上游配置树。

当前 live 回包已经确认的上游配置项包括：

- `node|区域`
- `os|操作系统`
- `cpu|CPU`
- `memory|内存`
- `network_type|网络类型`
- `bw|带宽`
- `ip_num|IP数量`
- `data_disk_size|数据盘`
- `port|远程端口`
- `in_bw|流入带宽`

这和本地前台购买页同步后的表单结构完全吻合。

### 10.7 “同步商品信息”按钮已实测成功

我已经对本地商品 3 实际发起了一次：

- `POST /admin01/product/sync_product_info`

实测回包：

- `status = 200`
- `msg = 同步数据成功`

回包文件：

- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\sync-product-info-pid3.json`

这进一步证明：

- 后台按钮链不是空壳
- 本地环境已经可以真实重放“同步上游商品”动作
- 后续可以继续围绕这条链观察同步前后数据库、配置项 id、服务快照回挂行为

### 10.8 数据库落库快照已经证实“商品定义”和“服务快照”是两层

我已经直接查询了本地数据库里的几张关键表，输出文件：

- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\product3-db-snapshot.tsv`
- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\product3-host-config.tsv`

当前已经证实：

- `shd_products`
  - 保存商品主定义
  - 包含 `type/api_type/zjmf_api_id/upstream_pid/server_group/upstream_price_type/upstream_price_value/auto_setup`
- `shd_pricing(type=product)`
  - 保存商品基础周期价格
- `shd_product_config_links`
  - 保存商品绑定了哪个配置组
- `shd_product_config_options`
  - 保存配置项定义，例如：
    - `node|区域`
    - `os|操作系统`
    - `cpu|CPU`
    - `memory|内存`
    - `bw|带宽`
    - `ip_num|IP数量`
- `shd_product_config_options_sub`
  - 保存配置项子项，例如：
    - `4|4核`
    - `2048|2G`
    - `vpc|VPC网络`
    - `1|1个`
    - `30G/SSD`
- `shd_pricing(type=configoptions)`
  - 保存每个子项的周期价格
- `shd_host_config_options`
  - 保存服务实例最终选中的配置快照

这里最关键的 live 证据是：

- `shd_host_config_options` 在这套老系统里不是 `hostid`
- 而是：
  - `relid`
  - `configid`
  - `optionid`
  - `qty`

也就是说，实例配置快照并不是 JSON 大字段，而是关系表快照。

当前演示服务 `relid = 1` 已实测保存了这组选择：

- 区域：`configid 17 -> optionid 95`
- 操作系统：`18 -> 96`
- CPU：`19 -> 116`
- 内存：`20 -> 121`
- 网络类型：`21 -> 126`
- 带宽：`22 -> 128`
- IP 数量：`23 -> 130`
- 数据盘：`24 -> 136`
- 远程端口：`25 -> 142`
- 流入带宽：`26 -> 143`

这和前台购买页默认提交值、服务详情页显示值是一一对应的。

## 11. 对重建系统的直接指导

根据这轮活体结果，后续重建必须遵守这些原则：

### 11.1 用户中心服务详情要做成“聚合工作台”

后端必须先聚合出一份类似 `host/header` 的大对象，再渲染页面。  
不要把它拆成一堆前端各自乱拉的小接口。

### 11.2 Provider 不能只做实例动作

至少要支持两层：

- 通用动作层：`/provision/default`
- 模块自定义层：`/provision/custom/{id}`

因为快照、备份、设置、救援、远程信息探测都走第二层。

### 11.3 商品同步必须重建本地销售模型

新系统里“同步上游商品”不能只存一个上游 product id。  
必须同步并重建：

- 商品基础信息
- 周期价格矩阵
- 配置项组
- 配置项子项
- 配置项价格
- 自定义字段

### 11.4 用户中心购买页价格重算必须服务端参与

老魔方并不是前端自己算出总价。  
前端只是提交当前配置快照，真正的价格汇总仍由服务端返回 HTML/结果对象。

### 11.5 `zjmf_api` 与 `dcimcloud` 是叠加关系

一个本地商品/服务可以同时满足：

- 上游交付方式是 `zjmf_api`
- 商品类型又是 `dcimcloud`

这意味着：

- “交付来源”
- “资源类型”

在模型上必须拆开。

## 12. 后台三张核心工作台：订单 / 账单 / 服务

这轮我已经把本地后台的真实接口跑出来了，不再只是读源码。

活体文件：

- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\order-1.json`
- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\invoice-summary-800000.json`
- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\clients-services-uid1-host1.json`
- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\invoice-addpay-page-800000.json`
- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\invoice-option-page-800000.json`
- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\invoice-refund-page-800000.json`
- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\clients-services-suspend-page-host1.json`
- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\order-invoice-host-chain.tsv`

### 12.1 后台订单详情：`GET /admin01/orders/:id`

当前演示样本：

- `GET /admin01/orders/1`

它返回的不是一张订单主表，而是一个聚合对象，至少包含：

- `order_notes`
- `subtotal`
- `i_status`
- `credit`
- `use_credit_limit`
- `i_payment`
- `username`
- `id`
- `uid`
- `status`
- `create_time`
- `amount`
- `promo_code`
- `payment`
- `invoiceid`
- `invoiceid_zh`

并且它还会直接带出订单下的服务对象数组：

- `server[]`

`server[]` 中当前演示样本已经确认有：

- `id`
- `billingcycle`
- `amount`
- `serverid`
- `username`
- `password`
- `type`
- `name`
- `domainstatus`
- `invoice_id`
- `runcreate`
- `sendwolcome`
- `server_group`

这说明后台订单详情页不是只看订单，而是直接把“订单交付结果”一起展示。

### 12.2 后台账单详情：`GET /admin01/invoice/summary/:id`

当前演示样本：

- `GET /admin01/invoice/summary/800000`

这个接口已经直接证明账单工作台的对象分层：

- `invoices`
- `invoice_items`
- `accounts`
- `exists_pay`

其中 `invoices` 里除了账单自身字段，还会补：

- `user_credit`
- `is_open_credit_limit`
- `credit_limit_balance`
- `pay_amount`
- `surplus`
- `credit_zh`
- `surplus_zh`
- `status_zh`

`invoice_items` 当前样本里是：

- `type = host`
- `rel_id = 1`
- `description = ECS - 成都 · 西信 | 弹性云 ...`
- `amount = 304.00`

这说明账单详情页里每条账单项本身就能反向定位到：

- 对应服务
- 对应购买周期
- 对应实例摘要

### 12.3 后台服务详情：`GET /admin01/clients_services?uid=...&hostselect=...`

当前演示样本：

- `GET /admin01/clients_services?uid=1&hostselect=1`

这是我目前抓到的老魔方后台最“值钱”的一个接口。  
它不是只返回服务对象，而是把服务工作台需要的整包数据一次性聚合好。

返回对象至少包含：

- `billing_cycle`
- `uid`
- `hostid`
- `host_data`
- `server_list`
- `host_option_config`
- `config_array`
- `product`
- `zjmf_api`
- `upstream_host`
- `reinstall_random_port`
- `module_button`
- `module_power_status`
- `module_admin_main_area`
- `module_admin_area`
- `module_upgrade`
- `custom_field`
- `custom_field_value`
- `domain_status_list`
- `dcim`
- `credit`

### 12.4 后台服务详情真正显示的是“三层数据”

从这个接口可以确认，后台服务详情页同时展示三层对象：

1. 本地服务实例：`host_data`
2. 本地商品定义：`product`
3. 上游远端摘要：`upstream_host`

当前样本里最关键的差异字段是：

- `host_data.dcimid = 8346`
- `host_data.upstream_configoption = {"upstream_host_id":8346,"upstream_provider_id":7469,...}`
- `upstream_host.domainstatus = Active`
- `dcim.url = https://www.xiaobaiyun.cn/#/cloud-server-info?id=8346`

这说明后台服务详情页本身就已经在做“本地服务对象”和“上游对象”的合并视图。

### 12.5 后台服务详情里的配置不是前端拼的

`host_option_config` 返回的是实例最终选择：

- `relid`
- `configid`
- `optionid`
- `qty`

`config_array` 返回的是商品当前可选配置树：

- 每个配置项的定义
- 子项
- 子项价格
- `show_pricing`
- `option_name_first`
- 以及 `linkage` 相关信息

所以后台服务详情是通过：

- `host_option_config`
- `config_array`

这两层一起做“当前配置 + 可变更配置”的展示，不是把实例快照简单塞成文本。

### 12.6 后台服务动作按钮是真正的 provider 能力镜像

当前服务详情接口返回的 `module_button` 已实测包含：

- `create`
- `suspend`
- `terminate`
- `on`
- `off`
- `reboot`
- `hard_off`
- `hard_reboot`
- `reinstall`
- `crack_pass`
- `rescue_system`
- `exitRescue`
- `vnc`
- `sync`

也就是说后台服务详情页和用户中心服务详情页一样，都是“按 provider 能力镜像按钮”，不是一套写死按钮。

## 13. 财务主链对象关系已实测打通

我已经直接查了本地 DB 样本链：

- `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\order-invoice-host-chain.tsv`

当前演示样本的真实关系是：

- `orders.id = 1`
- `orders.invoiceid = 800000`
- `invoices.id = 800000`
- `invoice_items.invoice_id = 800000`
- `invoice_items.type = host`
- `invoice_items.rel_id = 1`
- `host.id = 1`

这条链可以明确写成：

`orders -> invoices -> invoice_items -> host`

不是：

- 订单直接指向 host
- 或账单直接保存 host JSON

这是后续重建时必须保持的结构。

### 13.1 订单详情和账单详情各自的职责边界

订单详情页聚焦：

- 订单状态
- 客户
- 账单号
- 交付对象（服务数组）
- 手工审核/状态流转

账单详情页聚焦：

- 账单状态
- 账单项
- 收款记录
- 余额/信用额
- 退款
- 账单选项修改

服务详情页聚焦：

- 实例生命周期
- 配置项快照
- provider 动作
- 上游映射

这三张工作台对象边界在 live 数据里已经非常清楚。

### 13.2 账单动作的真实接口层

当前已实测的账单辅助接口有：

- `GET /admin01/invoice/addpay_page/:id`
- `GET /admin01/invoice/option_page/:id`
- `GET /admin01/invoice/refund_page?id=:id`

这些接口分别提供：

- 剩余可入账金额
- 账单选项编辑所需字段
- 可退款账户流水列表与客户余额

这说明后台账单工作台不是一个大接口包办所有动作，而是：

- 主详情接口
- 多个动作面板接口

联合构成。

### 13.3 服务动作里“暂停”有独立面板

当前已实测：

- `GET /admin01/clients_services/host_suspend?hostid=1`

返回：

- 当前暂停原因
- 暂停原因类型
- 可选原因字典

可选原因已确认包括：

- 到期
- 用量超额
- 未实名认证
- 其他

这说明后台服务动作不是直接点按钮就提交，而是很多动作先经过动作面板取参数。

## 14. 本轮最关键的证据文件

- 本地用户中心首页截图：  
  `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-clientarea-live.png`
- 本地服务详情截图：  
  `C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-servicedetail-live.png`
- 本地 `host/list`：  
  `C:\Users\Administrator\Desktop\IDC\output\local-user\host-list.json`
- 本地 `host/header`：  
  `C:\Users\Administrator\Desktop\IDC\output\local-user\host-header-1.json`
- 本地 `host/details`：  
  `C:\Users\Administrator\Desktop\IDC\output\local-user\host-details-1.json`
- 快照/备份子页：  
  `C:\Users\Administrator\Desktop\IDC\output\local-user\provision-custom-snapshot-1.html`
- 设置子页：  
  `C:\Users\Administrator\Desktop\IDC\output\local-user\provision-custom-setting-1.html`
- 远程信息探测：  
  `C:\Users\Administrator\Desktop\IDC\output\local-user\provision-custom-remoteInfo-1.json`
- 商品 3 购买页：  
  `C:\Users\Administrator\Desktop\IDC\output\local-user\cart-configureproduct-pid3.html`
- 商品 3 订单汇总：  
  `C:\Users\Administrator\Desktop\IDC\output\local-user\cart-ordersummary-pid3-default.html`
- 后台订单详情：  
  `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\order-1.json`
- 后台账单详情：  
  `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\invoice-summary-800000.json`
- 后台服务详情：  
  `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\clients-services-uid1-host1.json`
- 后台订单/账单/服务关系快照：  
  `C:\Users\Administrator\Desktop\IDC\output\legacy-admin\order-invoice-host-chain.tsv`

## 15. 下一步拆解顺序

下一轮最该继续深挖的是三块：

1. 后台订单详情、账单详情、服务详情三张工作台在前端页面上的真实 tab 和按钮布局
2. 服务详情里的动作请求体，尤其是 `reinstall / crack_pass / rescue_system / vnc / sync`
3. 商品编辑页 `自动开通 / 产品配置 / 升级选项 / 同步商品信息` 的字段级结构继续下钻
