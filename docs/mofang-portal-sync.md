# 魔方云同步到门户

这条同步链现在不再只把魔方云实例拉到后台服务表，而是会把门户侧需要的关键数据一起准备好。

## 已同步内容

- 客户资料：同步到 `Customer`
- 门户账号：自动补齐到 `CustomerUser`
- 云实例：同步到 `ServiceInstance`
- 资源明细：同步 IP、磁盘、快照、备份、安全组、VPC
- 导入订单：为同步实例补一张本地导入订单，门户订单页可直接看到
- 门户通知：实例新建或状态变化时，自动写入客户通知中心
- 续费预览：门户账单页会根据实例到期时间显示下一期续费预估

## 门户账号规则

- 每个同步客户会自动创建一个门户主账号
- 登录邮箱优先使用同步后的客户邮箱
- 如果魔方云没有邮箱，则回退为 `mf-<remote-id>@sync.local`
- 初始密码默认读取环境变量 `MOFANG_SYNC_PORTAL_DEFAULT_PASSWORD`
- 如果未配置，默认密码为 `Portal123!`

## 导入订单规则

- 每台同步实例会自动创建或更新一张本地导入订单
- 订单 `source` 固定为 `mofang-sync`
- 订单 `orderType` 固定为 `import`
- 订单金额默认取实例同步过来的 `monthlyCost`
- 已经存在的人工订单不会被同步逻辑覆盖
- 门户账单页里的“续费预览”不会伪造正式账单，只做续费提醒

## 同步结果新增字段

`/api/services/sync/provider` 现在会额外返回：

- `createdPortalUsers`
- `updatedPortalUsers`
- `queuedPortalNotifications`
- `createdImportOrders`
- `updatedImportOrders`

每条同步明细还会包含：

- `customerId`
- `customerName`
- `customerEmail`
- `portalEmail`
- `orderNo`
- `orderStatus`

## 当前边界

- 魔方云不会直接提供本系统里的支付、退款、发票与余额流水，所以这些财务数据仍由本地业务系统维护
- 当前同步重点是“客户、门户账号、服务、资源、通知、导入订单、续费预览”这条真实业务线
- 如果继续往正式运营版推进，下一步建议补：
  - 同步任务中心与失败重试页
  - 客户级差异比对与冲突处理
  - 门户重置密码与安全设置
  - 定时增量同步和变更审计
