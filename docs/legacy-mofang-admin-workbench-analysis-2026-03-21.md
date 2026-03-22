# 魔方财务后台工作台活体拆解（2026-03-21）

本文只记录本地老系统后台三张最核心工作台的活体结果：

- 产品订单列表
- 订单详情
- 账单详情
- 客户视角下的服务详情

证据来源同时包含三层：

1. 本地 3.7.4 老系统后台活体页面
2. 可读源码仓库 `zjmf-manger-decoded-ref`
3. 本地接口回包与数据库对象链

目标不是泛泛分析，而是为后续 1:1 重建后台 IA、页面工作台和对象主链提供直接依据。

## 1. 本次活体样本

- 本地后台地址：`http://127.0.0.1:8091/admin01/`
- 管理员：`admin`
- 样本订单：`order.id = 1`
- 样本账单：`invoice.id = 800000`
- 样本服务：`host.id = 1`
- 样本客户：`client.id = 1 / legacydemo`

对象链已通过数据库确认：

- `shd_orders.id = 1`
- `shd_orders.invoiceid = 800000`
- `shd_invoice_items.invoice_id = 800000`
- `shd_invoice_items.type = host`
- `shd_invoice_items.rel_id = 1`
- `shd_host.id = 1`

也就是：

`订单 -> 账单 -> 账单项 -> 服务`

证据文件：

- [order-invoice-host-chain.tsv](C:\Users\Administrator\Desktop\IDC\output\legacy-admin\order-invoice-host-chain.tsv)
- [order-1.json](C:\Users\Administrator\Desktop\IDC\output\legacy-admin\order-1.json)
- [invoice-summary-800000.json](C:\Users\Administrator\Desktop\IDC\output\legacy-admin\invoice-summary-800000.json)
- [clients-services-uid1-host1.json](C:\Users\Administrator\Desktop\IDC\output\legacy-admin\clients-services-uid1-host1.json)

## 2. 后台工作台总结构

这次确认到的后台不是“每个对象一张普通详情页”，而是三类工作台模型：

1. 列表工作台
2. 财务对象详情工作台
3. 客户对象内嵌服务工作台

对应关系：

- 订单：`#/order-list` + `#/order-detail?id={id}`
- 账单：`#/bill-management` + `#/bill-detail?id={id}&uid={uid}`
- 服务：不走独立一级 `service-detail`，而是嵌在客户工作台里  
  `#/customer-view/product-innerpage?id={uid}&hid={hostid}&fa=productList`

这点很关键：

- 订单详情是“订单工作台”
- 账单详情是“账单工作台”
- 服务详情在老系统里其实是“客户详情工作台中的产品/服务内页”

所以后续重建时，服务详情如果做成单独一级页面也可以，但 IA 层一定要保留“它本质从属于客户工作台”的关系。

## 3. 产品订单列表工作台

### 3.1 活体截图

- [legacy-admin-order-list-live.png](C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-admin-order-list-live.png)

### 3.2 路由

- 页面路由：`#/order-list`
- 活体最终 URL：`#/order-list?searchObj&page=1&limit=50&order=id&sort=desc`

### 3.3 布局结构

订单列表页是典型的运营后台列表工作台，结构固定：

1. 左侧业务菜单
2. 顶部状态切换 tabs
3. 页面说明 / 帮助文档入口
4. 主操作按钮区
5. 表格
6. 批量动作区
7. 汇总区
8. 分页区

### 3.4 状态条

活体确认存在 4 个 tab：

- 全部
- 待核验
- 已激活
- 已取消

这说明订单列表不是纯搜索页，而是“状态驱动列表页”。

### 3.5 主按钮

活体确认按钮：

- 添加新订单
- 高级搜索

### 3.6 表格字段

活体确认到的表格列：

- 勾选框
- ID
- 客户名
- 产品
- IP
- 下单时间
- 金额
- 付款状态/付款方式
- 状态
- 客户备注
- 提成/销售

样本订单行里可点击的跳转包括：

- 订单 ID -> `#/order-detail?id=1`
- 客户名 -> `#/customer-view/abstract?id=1`
- 产品名 -> `#/customer-view/product-innerpage?id=1&hid=1&fa=productList`

### 3.7 批量动作

活体确认存在：

- 核验通过
- 取消订单
- 删除订单

### 3.8 列表工作台建模结论

订单列表页是老魔方后台的标准模板之一，后续重建可抽象成统一组件：

- 状态 tabs
- 筛选条
- 表格
- 批量动作
- 汇总
- 分页

## 4. 订单详情工作台

### 4.1 活体截图

- [legacy-admin-order-detail-live.png](C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-admin-order-detail-live.png)

### 4.2 页面路由

- 页面：`#/order-detail?id=1`
- 页面标题：`订单详情 - 财务系统`

### 4.3 主 XHR

活体抓到页面只打了一个核心详情接口：

- `GET /admin01/orders/1`

证据：

- [.playwright-cli/network-2026-03-21T05-34-27-421Z.log](C:\Users\Administrator\Desktop\IDC\.playwright-cli\network-2026-03-21T05-34-27-421Z.log)

### 4.4 页面结构

订单详情是两段式：

1. 订单头信息区
2. 订单项目表格 + 动作区

#### 订单头信息区

左右两栏结构，活体确认字段：

- 客户
- 订单号
- 时间
- 优惠码
- IP 地址
- 账单信息
- 付款方式
- 金额
- 状态
- 客户备注

其中“账单信息”字段可直接跳账单详情：

- `#/bill-detail?id=800000&uid=1`

#### 订单项目表格

活体确认列：

- ID
- 条目
- 描述
- 付款周期
- 金额
- 状态
- 付款状态
- 操作

样本页里：

- 条目列跳服务详情
- 操作列包含“发送产品开通邮件”开关

### 4.5 页面动作

活体按钮：

- 核验通过
- 取消订单
- 删除订单

源码侧还确认存在：

- 保存订单备注
- 修改订单状态

### 4.6 对象关系

订单详情本身不负责支付和退款，但它是三跳关系的入口：

- 订单 -> 账单
- 订单 -> 服务
- 订单 -> 客户

因此它的职责不是“资金处理”，而是“订单交付视角聚合页”。

## 5. 账单详情工作台

### 5.1 活体截图

- [legacy-admin-bill-detail-live.png](C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-admin-bill-detail-live.png)

### 5.2 页面路由

- 页面：`#/bill-detail?id=800000&uid=1`
- 页面标题：`账单 - 财务系统`

### 5.3 主 XHR

活体抓到的页面请求比订单详情多，说明账单详情是典型“主详情 + 多动作面板”页面：

- `GET /admin01/invoice/summary/800000`
- `GET /admin01/common/get_email_tem?type=invoice`
- `GET /admin01/common/get_getways`
- `GET /admin01/invoice/option_page/800000`
- `GET /admin01/invoice/log_list?invoice_id=800000&uid=1`

证据：

- [.playwright-cli/network-2026-03-21T05-34-51-930Z.log](C:\Users\Administrator\Desktop\IDC\.playwright-cli\network-2026-03-21T05-34-51-930Z.log)

### 5.4 页面结构

账单详情工作台由 5 个区块组成：

1. 账单标题区
2. 账单信息卡
3. 交易明细表
4. 账单项目表
5. 操作日志表

#### 标题区

活体确认：

- 标题：`账单#800000`
- 右侧按钮：发送邮件

#### 账单信息卡

活体确认字段：

- 客户姓名
- 账单编号
- 账单生成时间
- 账单逾期时间
- 账单金额
- 余额支付
- 接口支付
- 支付时间
- 备注
- 状态文本
- 付款方式

右侧按钮：

- 编辑

#### 交易明细表

活体确认：

- 标题右侧按钮：手动入账
- 表格列：
  - 时间
  - 付款方式
  - 付款流水号
  - 金额
  - 操作

样本账单当前无交易数据，但动作区已存在。

#### 账单项目表

活体确认列：

- ID
- 描述
- 金额
- 操作

页面允许直接编辑：

- 描述
- 金额

同时存在：

- 删除账单项
- 新增空白账单项
- 保存更改
- 取消更改
- 返回

#### 操作日志表

列结构：

- ID
- 时间
- 描述
- 用户名
- IP 地址

### 5.5 页面职责

账单详情工作台不是单纯展示账单，而是后台财务操作中心，至少承担：

- 查看账单
- 编辑账单头
- 手工入账
- 编辑账单项
- 查看账单日志
- 发送账单邮件

源码还确认存在但样本页未直接操作的：

- 退款
- 余额/信用额支付
- 删除交易流水
- 复制账单
- 删除账单

### 5.6 与服务的关系

账单项 `type=host` 时，`rel_id` 直指服务对象。  
这意味着账单详情天然能回溯到服务生命周期。

## 6. 客户工作台中的服务详情

### 6.1 活体截图

- [legacy-admin-service-detail-live.png](C:\Users\Administrator\Desktop\IDC\output\playwright\legacy-admin-service-detail-live.png)

### 6.2 页面路由

不是一级 `/service-detail`，而是客户详情内页：

- `#/customer-view/product-innerpage?id=1&hid=1&fa=productList`

页面标题：

- `客户配置 - 财务系统`

### 6.3 主 XHR

活体抓到页面加载请求：

- `GET /admin01/clients_services?uid=1&hostselect=1`
- `GET /admin01/get_user?uid=1`
- `POST /admin01/provision/default`
- `GET /admin01/common/get_product_list`
- `GET /admin01/common/get_promo_code`
- `GET /admin01/common/get_getways`
- `GET /admin01/common/host_list?uid=1`

证据：

- [.playwright-cli/network-2026-03-21T05-35-15-418Z.log](C:\Users\Administrator\Desktop\IDC\.playwright-cli\network-2026-03-21T05-35-15-418Z.log)

这里最值得注意的是：

- 服务详情页本身就会主动请求 `provision/default`
- 说明它不是纯本地表单页，而是带 provider 扩展面板的聚合工作台

### 6.4 顶部客户工作台 tabs

服务详情不是孤立页，顶部仍保留客户工作台 tabs：

- 客户摘要
- 个人资料
- 产品/服务
- 账单
- 交易记录
- 工单
- 日志
- API 概览
- 通知日志
- 附件
- 跟进状态

这再次证明：  
服务详情在 IA 上属于“客户详情工作台的一个深层业务页”。

### 6.5 页面总体结构

当前服务详情页分成 4 块：

1. 服务选择与顶部快捷动作
2. 产品/接口基础信息区
3. 财务区
4. 产品配置区

#### 1. 服务选择与快捷动作

活体确认：

- 顶部可切换当前客户下的服务实例
- 快捷按钮：
  - 工单
  - 发送信息
  - 转移产品/服务

#### 2. 产品 / 接口基础信息区

活体确认字段：

- 订单号
- 商品/服务
- 接口
- IP 地址
- 端口
- 用户名
- 密码
- 其他 IP
- 主机名
- 客户备注
- 管理员备注
- 上游产品 ID

这里最关键的两个字段是：

- 接口
- 上游产品 ID

在当前样本里，它们直接暴露了这台服务并不是本地原生实例，而是：

- `api_type = zjmf_api`
- 上游商品 / 上游 host 已绑定

#### 3. 接口命令按钮组

活体确认存在这些按钮：

- 开通
- 暂停
- 删除
- 开机
- 关机
- 重启
- 硬关机
- 硬重启
- 重装系统
- 重置密码
- 救援系统
- 退出救援系统
- VNC
- 拉取信息

这是后台真正的 provider 动作主按钮区。

#### 4. 财务区

活体确认字段：

- 订购时间
- 首付金额
- 付款方式
- 到期时间
- 续费金额
- 余额自动续费
- 付款周期
- 状态
- 优惠码
- 成本

财务区右上动作：

- 交易流水
- 账单

这说明服务详情并不是只做运维动作，它同时承接服务对象上的财务维护。

#### 5. 产品配置区

活体确认当前 `dcimcloud` 商品的实例配置字段：

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

并存在按钮：

- 应用至接口

这说明后台服务详情不仅能看快照，还能把配置改动推送回 provider。

#### 6. 页面底部动作

活体确认：

- 保存更改
- 取消更改
- 删除记录

### 6.6 “拉取信息”按钮活体验证

这次做了真实点击，确认它不是直接执行，而是：

1. 先弹确认框  
   文案：`确定执行拉取信息操作?`
2. 点击确定后调用 provider 动作
3. 成功后出现提示  
   文案：`同步成功`
4. 页面随后重新请求 `clients_services`

活体网络证据：

- 再次 `POST /admin01/provision/default`
- 再次 `GET /admin01/clients_services`

说明这类按钮属于：

- 先确认
- 再走 provider
- 再刷新本地聚合页

而不是只改本地数据库。

## 7. 三张工作台的职责分工

### 7.1 订单详情

核心职责：

- 看订单头
- 看交付项
- 跳客户
- 跳服务
- 跳账单
- 审核 / 取消 / 删除

它更像“交付链入口”。

### 7.2 账单详情

核心职责：

- 看账单头
- 看交易明细
- 看账单项
- 直接做财务操作
- 看账单日志

它是“财务操作中心”。

### 7.3 服务详情

核心职责：

- 看服务基础信息
- 做 provider 动作
- 看并改服务财务字段
- 看并改实例配置
- 跳转客户其他业务维度

它是“客户工作台内的服务总工作台”。

## 8. 重建时必须保留的结构原则

### 8.1 后台不要做成一堆普通 CRUD 页

老魔方真正成熟的地方不是“页多”，而是对象分工清晰：

- 订单负责交付链
- 账单负责财务链
- 服务负责生命周期与资源链

### 8.2 服务详情必须保留三层混合视图

服务详情最少要同时承载：

1. 本地服务对象
2. 财务字段
3. Provider / 上游实例动作

否则就会退化成“一个实例信息页”，这也是现在很多重建版做得很丐的原因。

### 8.3 服务详情要放在客户工作台语境里

即使新系统可以额外提供独立服务详情路由，也要保留：

- 客户 tabs
- 客户上下文
- 服务只是客户资产的一部分

### 8.4 列表页必须是运营工作台

订单列表已经证明：

- 状态 tabs
- 批量动作
- 汇总
- 行内跳转

是老系统的基本风格，不是附加功能。

## 9. 当前已经钉死的后台对象主链

后台工作台可以按下面这条链完整重建：

```text
客户工作台
  -> 产品/服务内页
    -> 服务基础信息
    -> 财务字段
    -> Provider 动作
    -> 配置项快照 / 配置变更

订单列表
  -> 订单详情
    -> 账单详情
    -> 服务详情

账单详情
  -> 交易明细
  -> 账单项
  -> 服务 rel_id
```

## 10. 本轮最重要的结论

1. 后台服务详情并不是一级“服务模块详情页”，而是客户工作台里的深层业务内页。
2. 订单详情只打一条主详情接口，账单详情会同时加载编辑、网关、日志等辅助接口。
3. 服务详情页加载时就会触发 provider 扩展逻辑，不是纯本地表单。
4. “拉取信息”这类按钮已经活体验证为：确认框 -> provider 动作 -> 本地详情刷新。
5. 老魔方后台真正成熟的地方，在于三张工作台的职责边界非常清晰，而不是单纯功能多。

## 11. 服务详情按钮行为补充

这轮又补了两个高价值活体行为，专门用来区分“后台服务详情按钮到底是哪一类动作”。

### 11.1 重装系统

活体结果：

- 点击 `重装系统` 后，不会直接执行动作
- 页面会打开一个标题为 `重装系统` 的弹窗
- 弹窗内至少包含：
  - `操作系统` 选择框
  - `取消`
  - `确定`

因此它属于：

- 参数弹窗型动作

也就是：

`点击按钮 -> 选择参数 -> 确认提交 -> 调 provider`

### 11.2 VNC

活体结果：

- 点击 `VNC` 不会在当前页弹出普通对话框
- 浏览器会新开一个 tab
- 新 tab 地址形态是：

`/admin01/dcim/novnc?url=...&password=...&host_token=...`

这次实际打开的控制台地址里已经包含：

- 编码后的 `wss` 控制台地址
- `password`
- `host_token`

也就是说，老系统后台的 VNC 链路不是前端页面自己拼出来的，而是：

1. 服务详情页触发 provider 控制台动作
2. 后台生成 noVNC 专用地址
3. 新开 `dcim/novnc` 页面承接控制台

这对后续重建非常关键，因为这说明：

- VNC 不应该只是返回一段 `ws`
- 应该存在一个独立的控制台承接页

### 11.3 拉取信息

这轮已再次确认：

- 先弹确认框
- 再走 provider
- 成功后提示 `同步成功`
- 页面重新请求 `clients_services`

所以它属于：

- 确认后执行型动作

### 11.4 当前可确认的三类按钮

到目前为止，后台服务详情页按钮已经可以明确分成三类：

1. 确认后执行型  
   例：`拉取信息`
2. 参数弹窗型  
   例：`重装系统`
3. 新页跳转型  
   例：`VNC`
