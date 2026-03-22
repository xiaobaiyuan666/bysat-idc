import { z } from "zod";

const billingCycleEnum = z.enum([
  "MONTHLY",
  "QUARTERLY",
  "SEMI_ANNUALLY",
  "ANNUALLY",
  "BIENNIALLY",
  "TRIENNIALLY",
  "ONETIME",
]);
const invoiceTypeEnum = z.enum(["ORDER", "RENEWAL", "MANUAL", "CREDIT"]);
const invoiceStatusEnum = z.enum(["DRAFT", "ISSUED", "PARTIAL", "PAID", "OVERDUE", "VOID"]);

const customerTypeEnum = z.enum(["PERSONAL", "COMPANY", "RESELLER"]);
const customerStatusEnum = z.enum(["ACTIVE", "SUSPENDED", "OVERDUE", "ARCHIVED"]);
const customerCertificationStatusEnum = z.enum(["PENDING", "VERIFIED", "REJECTED"]);
const customerCertificationTypeEnum = z.enum(["PERSONAL", "COMPANY"]);
const customerFollowUpTypeEnum = z.enum(["NOTE", "CALL", "VISIT", "EMAIL"]);
const providerTypeEnum = z.enum(["MOFANG_CLOUD", "MANUAL"]);
const productCategoryEnum = z.enum([
  "CLOUD_SERVER",
  "BARE_METAL",
  "STORAGE",
  "CDN",
  "SECURITY",
  "NETWORK",
  "LICENSE",
  "DOMAIN",
  "OTHER",
]);
const productStatusEnum = z.enum(["ACTIVE", "DRAFT", "DISABLED"]);
const paymentMethodEnum = z.enum([
  "BALANCE",
  "ALIPAY",
  "WECHAT",
  "BANK_TRANSFER",
  "OFFLINE",
  "PAYPAL",
  "STRIPE",
  "OTHER",
]);
const paymentSignTypeEnum = z.enum(["HEADER_SECRET", "HMAC_SHA256"]);
const paymentStatusEnum = z.enum(["PENDING", "SUCCESS", "FAILED", "REFUNDED"]);
const ticketPriorityEnum = z.enum(["LOW", "NORMAL", "HIGH", "URGENT"]);
const ticketStatusEnum = z.enum([
  "OPEN",
  "PROCESSING",
  "WAITING_CUSTOMER",
  "RESOLVED",
  "CLOSED",
]);
const notificationChannelEnum = z.enum(["SYSTEM", "EMAIL", "SMS", "WEBHOOK"]);
const notificationPriorityEnum = z.enum(["LOW", "NORMAL", "HIGH", "URGENT"]);
const serviceStatusEnum = z.enum([
  "PENDING",
  "PROVISIONING",
  "ACTIVE",
  "SUSPENDED",
  "OVERDUE",
  "TERMINATED",
  "EXPIRED",
  "FAILED",
]);

export const loginSchema = z.object({
  email: z.email("请输入正确的邮箱地址"),
  password: z.string().min(6, "密码长度不能少于 6 位"),
});

export const customerSchema = z.object({
  name: z.string().min(2, "客户名称至少需要 2 个字符"),
  companyName: z.string().max(120, "企业名称不能超过 120 个字符").optional().or(z.literal("")),
  email: z.email("请输入正确的邮箱地址"),
  phone: z.string().max(40, "联系电话不能超过 40 个字符").optional().or(z.literal("")),
  contactQQ: z.string().max(40, "QQ 不能超过 40 个字符").optional().or(z.literal("")),
  contactWechat: z
    .string()
    .max(40, "微信号不能超过 40 个字符")
    .optional()
    .or(z.literal("")),
  type: customerTypeEnum,
  status: customerStatusEnum.default("ACTIVE"),
  level: z.string().min(1, "客户等级不能为空").max(40, "客户等级不能超过 40 个字符"),
  tags: z.string().max(200, "标签不能超过 200 个字符").optional().or(z.literal("")),
  notes: z.string().max(1000, "备注不能超过 1000 个字符").optional().or(z.literal("")),
});

export const customerContactSchema = z.object({
  name: z.string().min(2, "联系人姓名至少需要 2 个字符").max(60, "联系人姓名不能超过 60 个字符"),
  department: z.string().max(60, "部门不能超过 60 个字符").optional().or(z.literal("")),
  role: z.string().max(60, "角色不能超过 60 个字符").optional().or(z.literal("")),
  email: z.email("请输入正确的邮箱地址").optional().or(z.literal("")),
  phone: z.string().max(40, "联系电话不能超过 40 个字符").optional().or(z.literal("")),
  isPrimary: z.boolean().default(false),
  notes: z.string().max(500, "备注不能超过 500 个字符").optional().or(z.literal("")),
});

export const customerCertificationSchema = z.object({
  subjectType: customerCertificationTypeEnum,
  subjectName: z.string().min(2, "认证主体至少需要 2 个字符").max(120, "认证主体不能超过 120 个字符"),
  idNumber: z.string().max(60, "身份证号不能超过 60 个字符").optional().or(z.literal("")),
  businessLicenseNo: z.string().max(80, "营业执照号不能超过 80 个字符").optional().or(z.literal("")),
  status: customerCertificationStatusEnum.default("PENDING"),
  submittedAt: z.string().optional().or(z.literal("")),
  reviewNote: z.string().max(500, "审核备注不能超过 500 个字符").optional().or(z.literal("")),
});

export const customerFollowUpSchema = z.object({
  type: customerFollowUpTypeEnum.default("NOTE"),
  title: z.string().min(2, "跟进标题至少需要 2 个字符").max(120, "跟进标题不能超过 120 个字符"),
  content: z.string().min(4, "跟进内容至少需要 4 个字符").max(5000, "跟进内容不能超过 5000 个字符"),
  nextFollowAt: z.string().optional().or(z.literal("")),
});

export const productSchema = z.object({
  code: z.string().min(2, "产品编码至少需要 2 个字符").max(60, "产品编码不能超过 60 个字符"),
  name: z.string().min(2, "产品名称至少需要 2 个字符").max(120, "产品名称不能超过 120 个字符"),
  category: productCategoryEnum,
  status: productStatusEnum.default("ACTIVE"),
  billingCycle: billingCycleEnum,
  price: z.number().int().nonnegative("产品价格不能小于 0"),
  setupFee: z.number().int().nonnegative("开通费不能小于 0"),
  stock: z.number().int().nonnegative("库存不能小于 0"),
  autoProvision: z.boolean().default(true),
  providerType: providerTypeEnum.default("MOFANG_CLOUD"),
  providerProductId: z
    .string()
    .max(80, "云平台商品 ID 不能超过 80 个字符")
    .optional()
    .or(z.literal("")),
  regionTemplate: z
    .string()
    .max(60, "默认地域不能超过 60 个字符")
    .optional()
    .or(z.literal("")),
  description: z.string().max(1000, "产品说明不能超过 1000 个字符").optional().or(z.literal("")),
});

export const orderSchema = z.object({
  customerId: z.string().min(1, "请选择客户"),
  productId: z.string().min(1, "请选择产品"),
  planId: z.string().optional(),
  serviceName: z.string().min(2, "服务名称至少需要 2 个字符").max(120, "服务名称不能超过 120 个字符"),
  cycle: billingCycleEnum,
  quantity: z.number().int().min(1, "购买数量至少为 1").max(100, "购买数量不能超过 100"),
  notes: z.string().max(1000, "订单备注不能超过 1000 个字符").optional().or(z.literal("")),
});

export const orderActionSchema = z.object({
  action: z.enum(["approve", "cancel"]),
  reason: z.string().max(200, "原因不能超过 200 个字符").optional().or(z.literal("")),
});

export const paymentSchema = z.object({
  invoiceId: z.string().min(1, "请选择账单"),
  method: paymentMethodEnum,
  amount: z.number().int().positive("收款金额必须大于 0"),
  transactionNo: z
    .string()
    .max(120, "交易流水号不能超过 120 个字符")
    .optional()
    .or(z.literal("")),
});

export const invoiceCreateSchema = z.object({
  customerId: z.string().min(1, "请选择客户"),
  serviceId: z.string().optional().or(z.literal("")),
  taxProfileId: z.string().optional().or(z.literal("")),
  type: invoiceTypeEnum.default("MANUAL"),
  subtotal: z.number().int().positive("账单金额必须大于 0"),
  taxRate: z.number().int().min(0, "税率不能小于 0").max(100, "税率不能超过 100"),
  dueDate: z.string().min(1, "请选择到期日"),
  remark: z.string().max(500, "备注不能超过 500 个字符").optional().or(z.literal("")),
  issueNow: z.boolean().default(true),
});

export const invoiceActionSchema = z.object({
  action: z.enum(["issue", "void"]),
  reason: z.string().max(200, "原因不能超过 200 个字符").optional().or(z.literal("")),
  dueDate: z.string().optional().or(z.literal("")),
  status: invoiceStatusEnum.optional(),
});

export const refundSchema = z.object({
  amount: z.number().int().positive("退款金额必须大于 0"),
  reason: z
    .string()
    .min(2, "退款原因至少需要 2 个字符")
    .max(200, "退款原因不能超过 200 个字符"),
  method: paymentMethodEnum.optional(),
  transactionNo: z
    .string()
    .max(120, "退款流水号不能超过 120 个字符")
    .optional()
    .or(z.literal("")),
});

export const ticketSchema = z.object({
  customerId: z.string().min(1, "请选择客户"),
  serviceId: z.string().optional(),
  subject: z.string().min(4, "工单标题至少需要 4 个字符").max(120, "工单标题不能超过 120 个字符"),
  priority: ticketPriorityEnum,
  summary: z.string().min(6, "问题描述至少需要 6 个字符").max(5000, "问题描述不能超过 5000 个字符"),
});

export const billingSettingsSchema = z.object({
  invoiceLeadDays: z.coerce.number().int().min(0, "提前出账天数不能小于 0").max(90, "提前出账天数不能超过 90"),
  graceDays: z.coerce.number().int().min(0, "宽限天数不能小于 0").max(30, "宽限天数不能超过 30"),
  autoSuspendDays: z.coerce.number().int().min(0, "自动暂停天数不能小于 0").max(90, "自动暂停天数不能超过 90"),
  autoTerminateDays: z.coerce.number().int().min(1, "自动释放天数至少为 1").max(180, "自动释放天数不能超过 180"),
  autoRenewByBalance: z.boolean(),
  allowNegativeBalance: z.boolean(),
  invoicePrefix: z.string().min(2, "发票前缀至少需要 2 个字符").max(12, "发票前缀不能超过 12 个字符"),
  invoiceIssuerName: z.string().min(2, "开票主体至少需要 2 个字符").max(80, "开票主体不能超过 80 个字符"),
  invoiceTaxNo: z.string().max(40, "税号不能超过 40 个字符").optional().or(z.literal("")),
  financeEmail: z.email("请输入正确的财务邮箱地址").optional().or(z.literal("")),
  defaultTaxRate: z.coerce.number().int().min(0, "默认税率不能小于 0").max(100, "默认税率不能超过 100"),
});

export const taxProfileSchema = z.object({
  code: z.string().min(2, "税率编码至少需要 2 个字符").max(30, "税率编码不能超过 30 个字符"),
  name: z.string().min(2, "税率名称至少需要 2 个字符").max(60, "税率名称不能超过 60 个字符"),
  taxRate: z.number().int().min(0, "税率不能小于 0").max(100, "税率不能超过 100"),
  description: z.string().max(500, "税务说明不能超过 500 个字符").optional().or(z.literal("")),
  isDefault: z.boolean().default(false),
  isActive: z.boolean().default(true),
});

export const paymentGatewaySchema = z.object({
  method: paymentMethodEnum,
  name: z.string().min(2, "渠道名称至少需要 2 个字符").max(60, "渠道名称不能超过 60 个字符"),
  merchantId: z.string().max(60, "商户号不能超过 60 个字符").optional().or(z.literal("")),
  appId: z.string().max(80, "应用标识不能超过 80 个字符").optional().or(z.literal("")),
  apiBaseUrl: z.string().max(200, "接口地址不能超过 200 个字符").optional().or(z.literal("")),
  signType: paymentSignTypeEnum.default("HEADER_SECRET"),
  callbackSecret: z.string().min(6, "回调密钥至少需要 6 个字符").max(120, "回调密钥不能超过 120 个字符"),
  callbackHeader: z.string().max(40, "签名头不能超过 40 个字符").optional().or(z.literal("")),
  notifyUrl: z.string().max(200, "回调地址不能超过 200 个字符").optional().or(z.literal("")),
  returnUrl: z.string().max(200, "返回地址不能超过 200 个字符").optional().or(z.literal("")),
  isEnabled: z.boolean().default(true),
  remark: z.string().max(500, "备注不能超过 500 个字符").optional().or(z.literal("")),
});

export const balanceAdjustmentSchema = z.object({
  amount: z.coerce.number(),
  description: z
    .string()
    .min(4, "调整说明至少需要 4 个字符")
    .max(200, "调整说明不能超过 200 个字符"),
});

export const paymentCallbackSchema = z
  .object({
    method: paymentMethodEnum,
    invoiceNo: z.string().optional(),
    paymentNo: z.string().optional(),
    transactionNo: z.string().optional(),
    amount: z.coerce.number().nonnegative().optional(),
    status: paymentStatusEnum.default("SUCCESS"),
    signature: z.string().optional(),
    payload: z.unknown().optional(),
    message: z.string().optional(),
  })
  .refine((value) => Boolean(value.invoiceNo || value.paymentNo || value.transactionNo), {
    message: "至少需要提供一个支付标识",
    path: ["invoiceNo"],
  });

export const manualNotificationSchema = z.object({
  customerId: z.string().optional().or(z.literal("")),
  channel: notificationChannelEnum.default("SYSTEM"),
  priority: notificationPriorityEnum.default("NORMAL"),
  recipient: z.string().min(2, "接收对象不能为空").max(120, "接收对象不能超过 120 个字符"),
  recipientName: z
    .string()
    .max(80, "接收方名称不能超过 80 个字符")
    .optional()
    .or(z.literal("")),
  subject: z.string().max(160, "通知标题不能超过 160 个字符").optional().or(z.literal("")),
  content: z.string().min(4, "通知内容至少需要 4 个字符").max(5000, "通知内容不能超过 5000 个字符"),
  module: z.string().max(40, "模块标识不能超过 40 个字符").optional().or(z.literal("")),
  relatedType: z.string().max(40, "关联类型不能超过 40 个字符").optional().or(z.literal("")),
  relatedId: z.string().max(80, "关联 ID 不能超过 80 个字符").optional().or(z.literal("")),
});

export const ticketReplySchema = z.object({
  content: z
    .string()
    .min(2, "回复内容至少需要 2 个字符")
    .max(5000, "回复内容不能超过 5000 个字符"),
  status: ticketStatusEnum.optional(),
  assignedToId: z.string().optional(),
  isInternal: z.boolean().default(false),
});

export const serviceUpdateSchema = z.object({
  name: z.string().min(2, "服务名称至少需要 2 个字符").max(120, "服务名称不能超过 120 个字符"),
  hostname: z.string().max(120, "主机名不能超过 120 个字符").optional().or(z.literal("")),
  region: z.string().max(60, "地域不能超过 60 个字符").optional().or(z.literal("")),
  planId: z.string().optional().or(z.literal("")),
  vpcNetworkId: z.string().optional().or(z.literal("")),
  status: serviceStatusEnum,
  monthlyCost: z.number().int().nonnegative("月费用不能小于 0"),
  nextDueDate: z.string().optional().or(z.literal("")),
});

export const customerUserSchema = z.object({
  customerId: z.string().min(1, "请选择所属客户"),
  name: z.string().min(2, "联系人名称至少需要 2 个字符").max(60, "联系人名称不能超过 60 个字符"),
  email: z.email("请输入正确的邮箱地址"),
  phone: z.string().max(40, "手机号不能超过 40 个字符").optional().or(z.literal("")),
  password: z
    .string()
    .min(6, "门户密码长度不能少于 6 位")
    .max(60, "门户密码长度不能超过 60 位")
    .optional(),
  role: z.string().min(2, "联系人角色不能为空").max(30, "联系人角色不能超过 30 个字符"),
  isOwner: z.boolean().default(false),
  isActive: z.boolean().default(true),
});

export const cloudRegionSchema = z.object({
  code: z.string().min(2, "地域编码至少需要 2 个字符").max(40, "地域编码不能超过 40 个字符"),
  name: z.string().min(2, "地域名称至少需要 2 个字符").max(60, "地域名称不能超过 60 个字符"),
  location: z.string().max(80, "地域位置不能超过 80 个字符").optional().or(z.literal("")),
  providerType: providerTypeEnum,
  providerRegionId: z
    .string()
    .max(80, "平台地域 ID 不能超过 80 个字符")
    .optional()
    .or(z.literal("")),
  description: z.string().max(500, "地域说明不能超过 500 个字符").optional().or(z.literal("")),
  isActive: z.boolean().default(true),
  sortOrder: z.number().int().min(0, "排序值不能小于 0").default(0),
});

export const cloudZoneSchema = z.object({
  regionId: z.string().min(1, "请选择所属地域"),
  code: z.string().min(2, "可用区编码至少需要 2 个字符").max(40, "可用区编码不能超过 40 个字符"),
  name: z.string().min(2, "可用区名称至少需要 2 个字符").max(60, "可用区名称不能超过 60 个字符"),
  providerZoneId: z
    .string()
    .max(80, "平台可用区 ID 不能超过 80 个字符")
    .optional()
    .or(z.literal("")),
  isActive: z.boolean().default(true),
  sortOrder: z.number().int().min(0, "排序值不能小于 0").default(0),
});

export const cloudFlavorSchema = z.object({
  code: z.string().min(2, "规格编码至少需要 2 个字符").max(40, "规格编码不能超过 40 个字符"),
  name: z.string().min(2, "规格名称至少需要 2 个字符").max(80, "规格名称不能超过 80 个字符"),
  category: productCategoryEnum,
  cpu: z.number().int().positive("CPU 核心数必须大于 0"),
  memoryGb: z.number().int().positive("内存容量必须大于 0"),
  storageGb: z.number().int().nonnegative("系统盘容量不能小于 0").optional(),
  bandwidthMbps: z.number().int().nonnegative("带宽不能小于 0").optional(),
  gpuCount: z.number().int().nonnegative("GPU 数量不能小于 0").default(0),
  description: z.string().max(500, "规格说明不能超过 500 个字符").optional().or(z.literal("")),
  isActive: z.boolean().default(true),
  sortOrder: z.number().int().min(0, "排序值不能小于 0").default(0),
});

export const cloudImageSchema = z.object({
  code: z.string().min(2, "镜像编码至少需要 2 个字符").max(40, "镜像编码不能超过 40 个字符"),
  name: z.string().min(2, "镜像名称至少需要 2 个字符").max(80, "镜像名称不能超过 80 个字符"),
  osType: z.string().min(2, "系统类型不能为空").max(40, "系统类型不能超过 40 个字符"),
  version: z.string().max(40, "版本号不能超过 40 个字符").optional().or(z.literal("")),
  architecture: z
    .string()
    .min(2, "架构不能为空")
    .max(20, "架构不能超过 20 个字符")
    .default("x86_64"),
  regionId: z.string().optional().or(z.literal("")),
  description: z.string().max(500, "镜像说明不能超过 500 个字符").optional().or(z.literal("")),
  isPublic: z.boolean().default(true),
  isActive: z.boolean().default(true),
});

export const cloudPlanSchema = z.object({
  productId: z.string().min(1, "请选择关联产品"),
  regionId: z.string().min(1, "请选择可售地域"),
  zoneId: z.string().optional().or(z.literal("")),
  flavorId: z.string().optional().or(z.literal("")),
  imageId: z.string().optional().or(z.literal("")),
  code: z.string().min(2, "方案编码至少需要 2 个字符").max(50, "方案编码不能超过 50 个字符"),
  name: z.string().min(2, "方案名称至少需要 2 个字符").max(120, "方案名称不能超过 120 个字符"),
  billingCycle: billingCycleEnum,
  salePrice: z.number().int().nonnegative("销售价不能小于 0"),
  marketPrice: z.number().int().nonnegative("划线价不能小于 0").optional(),
  setupFee: z.number().int().nonnegative("开通费不能小于 0").default(0),
  stock: z.number().int().nonnegative("库存不能小于 0").default(0),
  isPublic: z.boolean().default(true),
  isActive: z.boolean().default(true),
  providerNodeId: z.string().max(80, "平台节点标识不能超过 80 个字符").optional().or(z.literal("")),
  providerOsCode: z.string().max(80, "平台系统标识不能超过 80 个字符").optional().or(z.literal("")),
  configCpu: z.number().int().nonnegative("CPU 不能小于 0").optional(),
  configMemoryGb: z.number().int().nonnegative("内存不能小于 0").optional(),
  systemDiskSize: z.number().int().nonnegative("系统盘不能小于 0").optional(),
  dataDiskSize: z.number().int().nonnegative("数据盘不能小于 0").optional(),
  networkType: z.string().max(30, "网络类型不能超过 30 个字符").optional().or(z.literal("")),
  bandwidthMbps: z.number().int().nonnegative("出带宽不能小于 0").optional(),
  inboundBandwidthMbps: z.number().int().nonnegative("入带宽不能小于 0").optional(),
  flowLimitGb: z.number().int().nonnegative("流量限制不能小于 0").optional(),
  flowBillingMode: z.string().max(30, "流量计费方式不能超过 30 个字符").optional().or(z.literal("")),
  ipCount: z.number().int().nonnegative("IP 数量不能小于 0").optional(),
  peakDefenceGbps: z.number().int().nonnegative("防御峰值不能小于 0").optional(),
  description: z.string().max(1000, "方案说明不能超过 1000 个字符").optional().or(z.literal("")),
});

export const portalOrderSchema = z.object({
  planId: z.string().min(1, "请选择可购买方案"),
  quantity: z.coerce.number().int().min(1, "购买数量至少为 1").max(20, "购买数量不能超过 20"),
  serviceName: z.string().min(2, "实例名称至少需要 2 个字符").max(120, "实例名称不能超过 120 个字符"),
  notes: z.string().max(500, "订单备注不能超过 500 个字符").optional().or(z.literal("")),
});

export const portalTicketCreateSchema = z.object({
  serviceId: z.string().optional().or(z.literal("")),
  subject: z.string().min(4, "工单标题至少需要 4 个字符").max(120, "工单标题不能超过 120 个字符"),
  priority: ticketPriorityEnum.default("NORMAL"),
  summary: z.string().min(6, "问题描述至少需要 6 个字符").max(5000, "问题描述不能超过 5000 个字符"),
});

export const portalTicketReplySchema = z.object({
  content: z.string().min(2, "回复内容至少需要 2 个字符").max(5000, "回复内容不能超过 5000 个字符"),
});
