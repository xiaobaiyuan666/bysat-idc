import "dotenv/config";

import bcrypt from "bcryptjs";

import {
  AdminRole,
  BillingCycle,
  CreditTransactionType,
  CustomerStatus,
  CustomerType,
  InvoiceStatus,
  InvoiceType,
  NotificationChannel,
  NotificationPriority,
  NotificationStatus,
  OrderStatus,
  PaymentMethod,
  PaymentSignType,
  PaymentStatus,
  PrismaClient,
  ProductCategory,
  ProductStatus,
  ProviderSyncStatus,
  ProviderType,
  RefundStatus,
  ServiceStatus,
  TaskRunStatus,
  TicketPriority,
  TicketStatus,
} from "@prisma/client";

process.env.DATABASE_URL ??= "file:./dev.db";

const db = new PrismaClient();

async function main() {
  await db.asyncTaskJob.deleteMany();
  await db.notificationMessage.deleteMany();
  await db.notificationTemplate.deleteMany();
  await db.refundRecord.deleteMany();
  await db.paymentGatewayConfig.deleteMany();
  await db.invoiceTaxProfile.deleteMany();
  await db.customerFollowUp.deleteMany();
  await db.customerCertification.deleteMany();
  await db.customerContact.deleteMany();
  await db.customerUser.deleteMany();
  await db.ticketReply.deleteMany();
  await db.ticket.deleteMany();
  await db.paymentCallbackLog.deleteMany();
  await db.payment.deleteMany();
  await db.billingJob.deleteMany();
  await db.invoice.deleteMany();
  await db.orderItem.deleteMany();
  await db.serviceSecurityGroupRule.deleteMany();
  await db.serviceSecurityGroup.deleteMany();
  await db.serviceSnapshot.deleteMany();
  await db.serviceBackup.deleteMany();
  await db.serviceDisk.deleteMany();
  await db.serviceIpAddress.deleteMany();
  await db.serviceInstance.deleteMany();
  await db.serviceVpcNetwork.deleteMany();
  await db.cloudPlan.deleteMany();
  await db.cloudImage.deleteMany();
  await db.cloudFlavor.deleteMany();
  await db.cloudZone.deleteMany();
  await db.cloudRegion.deleteMany();
  await db.order.deleteMany();
  await db.creditTransaction.deleteMany();
  await db.product.deleteMany();
  await db.customer.deleteMany();
  await db.providerSyncLog.deleteMany();
  await db.auditLog.deleteMany();
  await db.billingSetting.deleteMany();
  await db.adminUser.deleteMany();

  const passwordHash = await bcrypt.hash(
    process.env.SEED_ADMIN_PASSWORD ?? "Admin123!",
    10,
  );
  const portalPasswordHash = await bcrypt.hash(
    process.env.SEED_PORTAL_PASSWORD ?? "Portal123!",
    10,
  );

  const now = new Date();
  const daysAgo = (days: number) =>
    new Date(now.getTime() - 1000 * 60 * 60 * 24 * days);
  const daysLater = (days: number) =>
    new Date(now.getTime() + 1000 * 60 * 60 * 24 * days);

  await db.billingSetting.create({
    data: {
      id: "default",
      invoiceLeadDays: 7,
      graceDays: 3,
      autoSuspendDays: 7,
      autoTerminateDays: 30,
      autoRenewByBalance: true,
      allowNegativeBalance: false,
      invoicePrefix: "INV",
      invoiceIssuerName: "上海星云算力科技有限公司",
      invoiceTaxNo: "91310000IDC2026001X",
      financeEmail: "finance@idc.local",
      defaultTaxRate: 13,
    },
  });

  await db.invoiceTaxProfile.createMany({
    data: [
      {
        code: "VAT13",
        name: "增值税 13%",
        taxRate: 13,
        description: "适用于标准云服务器、存储和网络资源。",
        isDefault: true,
      },
      {
        code: "VAT6",
        name: "增值税 6%",
        taxRate: 6,
        description: "适用于部分技术服务和增值服务。",
        isDefault: false,
      },
      {
        code: "ZERO",
        name: "零税率",
        taxRate: 0,
        description: "适用于内部测试、赠送账单和免税场景。",
        isDefault: false,
      },
    ],
  });

  await db.paymentGatewayConfig.createMany({
    data: [
      {
        method: PaymentMethod.ALIPAY,
        name: "支付宝当面付",
        merchantId: "208810100000001",
        appId: "2026000000000001",
        apiBaseUrl: "https://openapi.alipay.com/gateway.do",
        signType: PaymentSignType.HMAC_SHA256,
        callbackSecret: "alipay-demo-hmac-secret",
        callbackHeader: "x-payment-signature",
        notifyUrl: "https://billing.idc.local/api/payments/callback",
        returnUrl: "https://portal.idc.local/portal/invoices",
        isEnabled: true,
        remark: "演示渠道，按 HMAC-SHA256 校验回调签名。",
      },
      {
        method: PaymentMethod.WECHAT,
        name: "微信支付 V3",
        merchantId: "1900000109",
        appId: "wx1234567890demo",
        apiBaseUrl: "https://api.mch.weixin.qq.com",
        signType: PaymentSignType.HMAC_SHA256,
        callbackSecret: "wechat-demo-hmac-secret",
        callbackHeader: "x-payment-signature",
        notifyUrl: "https://billing.idc.local/api/payments/callback",
        returnUrl: "https://portal.idc.local/portal/invoices",
        isEnabled: true,
        remark: "演示渠道，模拟统一签名回调。",
      },
      {
        method: PaymentMethod.BANK_TRANSFER,
        name: "银行转账对账通道",
        merchantId: "BANK-OFFLINE-001",
        signType: PaymentSignType.HEADER_SECRET,
        callbackSecret: "bank-transfer-secret",
        callbackHeader: "x-payment-secret",
        notifyUrl: "https://billing.idc.local/api/payments/callback",
        isEnabled: true,
        remark: "用于对账系统或内部打款回调。",
      },
    ],
  });

  const [admin, finance, support, operations] = await Promise.all([
    db.adminUser.create({
      data: {
        name: "System Admin",
        email: process.env.SEED_ADMIN_EMAIL ?? "admin@idc.local",
        role: AdminRole.SUPER_ADMIN,
        passwordHash,
        lastLoginAt: daysAgo(1),
      },
    }),
    db.adminUser.create({
      data: {
        name: "Finance Lead",
        email: "finance@idc.local",
        role: AdminRole.FINANCE,
        passwordHash,
        lastLoginAt: daysAgo(1),
      },
    }),
    db.adminUser.create({
      data: {
        name: "Support Desk",
        email: "support@idc.local",
        role: AdminRole.SUPPORT,
        passwordHash,
        lastLoginAt: daysAgo(2),
      },
    }),
    db.adminUser.create({
      data: {
        name: "Ops Center",
        email: "ops@idc.local",
        role: AdminRole.OPERATIONS,
        passwordHash,
        lastLoginAt: daysAgo(1),
      },
    }),
  ]);

  await db.notificationTemplate.createMany({
    data: [
      {
        code: "ORDER_CREATED",
        name: "新订单创建通知",
        channel: NotificationChannel.SYSTEM,
        subject: "订单 {{order_no}} 已创建",
        content:
          "尊敬的 {{customer_name}}，您的订单 {{order_no}} 已创建，对应服务 {{service_no}} 与账单 {{invoice_no}} 已生成。",
        variables: JSON.stringify(["customer_name", "order_no", "service_no", "invoice_no"]),
      },
      {
        code: "PAYMENT_SUCCESS",
        name: "收款成功通知",
        channel: NotificationChannel.SYSTEM,
        subject: "账单 {{invoice_no}} 收款成功",
        content:
          "您好，账单 {{invoice_no}} 已完成收款，到账金额 {{amount}} 元，支付方式 {{method}}。",
        variables: JSON.stringify(["invoice_no", "amount", "method"]),
      },
      {
        code: "PAYMENT_REFUND",
        name: "退款通知",
        channel: NotificationChannel.SYSTEM,
        subject: "退款 {{refund_no}} 已处理",
        content:
          "您好，收款 {{payment_no}} 已发起退款，退款金额 {{amount}} 元，退款原因：{{reason}}。",
        variables: JSON.stringify(["payment_no", "refund_no", "amount", "reason"]),
      },
      {
        code: "TICKET_CREATED",
        name: "工单创建通知",
        channel: NotificationChannel.SYSTEM,
        subject: "工单 {{ticket_no}} 已创建",
        content:
          "您好，工单 {{ticket_no}} 已创建，主题为“{{subject}}”，客服团队将尽快处理。",
        variables: JSON.stringify(["ticket_no", "subject"]),
      },
    ],
  });

  const [customerA, customerB, customerC, customerD] = await Promise.all([
    db.customer.create({
      data: {
        customerNo: "CUS-202603-001",
        name: "StarGalaxy Data",
        companyName: "Shanghai StarGalaxy Data Co., Ltd.",
        email: "ops@stargalaxy.cn",
        phone: "13800000001",
        contactWechat: "stargalaxy-ops",
        type: CustomerType.COMPANY,
        status: CustomerStatus.ACTIVE,
        level: "KA",
        creditBalance: 2_800_000,
        tags: "cloud,managed-service,ka",
        notes: "Key account with managed Kubernetes and production workload.",
      },
    }),
    db.customer.create({
      data: {
        customerNo: "CUS-202603-002",
        name: "CloudMatrix Tech",
        companyName: "Hangzhou CloudMatrix Technology Co., Ltd.",
        email: "finance@cloudmatrix.cn",
        phone: "13800000002",
        contactQQ: "2311999",
        type: CustomerType.COMPANY,
        status: CustomerStatus.ACTIVE,
        level: "VIP",
        creditBalance: 1_560_000,
        tags: "object-storage,cdn,vip",
        notes: "Uses storage, CDN and backup products across multiple regions.",
      },
    }),
    db.customer.create({
      data: {
        customerNo: "CUS-202603-003",
        name: "NorthShore Network",
        companyName: "Shenzhen NorthShore Network Ltd.",
        email: "partner@nshore.cn",
        phone: "13800000003",
        type: CustomerType.RESELLER,
        status: CustomerStatus.OVERDUE,
        level: "Reseller",
        creditBalance: 520_000,
        tags: "reseller,security,bgp",
        notes: "Regional reseller with overdue renewal invoice.",
      },
    }),
    db.customer.create({
      data: {
        customerNo: "CUS-202603-004",
        name: "Wu Feng",
        email: "wufeng@example.com",
        phone: "13800000004",
        type: CustomerType.PERSONAL,
        status: CustomerStatus.ACTIVE,
        level: "Standard",
        creditBalance: 120_000,
        tags: "lab,developer",
        notes: "Personal lab customer for database and backup testing.",
      },
    }),
  ]);

  await db.customerUser.createMany({
    data: [
      {
        customerId: customerA.id,
        name: "StarGalaxy 运维",
        email: customerA.email,
        phone: customerA.phone,
        passwordHash: portalPasswordHash,
        role: "OWNER",
        isOwner: true,
        lastLoginAt: daysAgo(1),
      },
      {
        customerId: customerB.id,
        name: "CloudMatrix 财务",
        email: customerB.email,
        phone: customerB.phone,
        passwordHash: portalPasswordHash,
        role: "BILLING",
        isOwner: true,
        lastLoginAt: daysAgo(2),
      },
      {
        customerId: customerC.id,
        name: "NorthShore 技术",
        email: customerC.email,
        phone: customerC.phone,
        passwordHash: portalPasswordHash,
        role: "TECH",
        isOwner: true,
        lastLoginAt: daysAgo(3),
      },
      {
        customerId: customerD.id,
        name: "Wu Feng",
        email: customerD.email,
        phone: customerD.phone,
        passwordHash: portalPasswordHash,
        role: "OWNER",
        isOwner: true,
        lastLoginAt: daysAgo(4),
      },
    ],
  });

  await db.customerContact.createMany({
    data: [
      {
        customerId: customerA.id,
        name: "Li Ming",
        department: "运维部",
        role: "技术负责人",
        email: "li.ming@stargalaxy.cn",
        phone: "13910000001",
        isPrimary: true,
        notes: "负责生产集群和紧急故障升级。",
      },
      {
        customerId: customerA.id,
        name: "Zhou Yun",
        department: "采购部",
        role: "采购经理",
        email: "procurement@stargalaxy.cn",
        phone: "13910000002",
        isPrimary: false,
        notes: "负责合同、对账与续费审批。",
      },
      {
        customerId: customerB.id,
        name: "Chen Yu",
        department: "财务部",
        role: "财务主管",
        email: "billing@cloudmatrix.cn",
        phone: "13910000003",
        isPrimary: true,
        notes: "主要处理发票与收款核销。",
      },
      {
        customerId: customerC.id,
        name: "Luo Peng",
        department: "渠道事业部",
        role: "渠道经理",
        email: "partner-ops@nshore.cn",
        phone: "13910000004",
        isPrimary: true,
        notes: "负责代理节点续费和渠道折扣。",
      },
      {
        customerId: customerD.id,
        name: "Wu Feng",
        department: "个人",
        role: "账户所有者",
        email: "wufeng@example.com",
        phone: "13800000004",
        isPrimary: true,
        notes: "个人测试账户主联系人。",
      },
    ],
  });

  await Promise.all([
    db.customerCertification.create({
      data: {
        customerId: customerA.id,
        subjectType: "COMPANY",
        subjectName: customerA.companyName ?? customerA.name,
        businessLicenseNo: "91310118MA1KAX001A",
        status: "VERIFIED",
        submittedAt: daysAgo(90),
        verifiedAt: daysAgo(88),
        reviewNote: "企业实名资料齐全，已通过审核。",
      },
    }),
    db.customerCertification.create({
      data: {
        customerId: customerB.id,
        subjectType: "COMPANY",
        subjectName: customerB.companyName ?? customerB.name,
        businessLicenseNo: "91330106MA2CFX882B",
        status: "PENDING",
        submittedAt: daysAgo(2),
        reviewNote: "等待补充法人身份证照片。",
      },
    }),
    db.customerCertification.create({
      data: {
        customerId: customerC.id,
        subjectType: "COMPANY",
        subjectName: customerC.companyName ?? customerC.name,
        businessLicenseNo: "91440300MA5F9N889C",
        status: "REJECTED",
        submittedAt: daysAgo(12),
        reviewNote: "营业执照扫描件过期，需要重新上传。",
      },
    }),
    db.customerCertification.create({
      data: {
        customerId: customerD.id,
        subjectType: "PERSONAL",
        subjectName: customerD.name,
        idNumber: "320311199201015678",
        status: "VERIFIED",
        submittedAt: daysAgo(30),
        verifiedAt: daysAgo(29),
        reviewNote: "个人实名已完成。",
      },
    }),
  ]);

  await db.customerFollowUp.createMany({
    data: [
      {
        customerId: customerA.id,
        type: "VISIT",
        title: "季度扩容需求确认",
        content: "客户计划在下个季度将香港节点扩容到三套业务集群，需要提前预留库存和价格锁定。",
        nextFollowAt: daysLater(10),
        operatorId: operations.id,
        operatorName: operations.name,
        operatorRole: operations.role,
      },
      {
        customerId: customerB.id,
        type: "EMAIL",
        title: "补充实名材料提醒",
        content: "已通过邮件提醒客户补充法人身份证与授权书，等待客户回传。",
        nextFollowAt: daysLater(2),
        operatorId: finance.id,
        operatorName: finance.name,
        operatorRole: finance.role,
      },
      {
        customerId: customerC.id,
        type: "CALL",
        title: "逾期账单电话催缴",
        content: "已电话联系渠道经理，客户承诺本周内处理安全节点续费账单。",
        nextFollowAt: daysLater(3),
        operatorId: support.id,
        operatorName: support.name,
        operatorRole: support.role,
      },
      {
        customerId: customerD.id,
        type: "NOTE",
        title: "测试客户保留",
        content: "保留为产品和账单联调测试客户，不纳入营销计划。",
        operatorId: admin.id,
        operatorName: admin.name,
        operatorRole: admin.role,
      },
    ],
  });

  const [productA, productB, productC, productD] = await Promise.all([
    db.product.create({
      data: {
        code: "CLOUD-GEN2-2C4G",
        name: "云服务器通用型 2C4G",
        category: ProductCategory.CLOUD_SERVER,
        status: ProductStatus.ACTIVE,
        billingCycle: BillingCycle.MONTHLY,
        price: 68_800,
        setupFee: 0,
        stock: 120,
        autoProvision: true,
        providerType: ProviderType.MOFANG_CLOUD,
        providerProductId: "cloud-gen2-2c4g",
        regionTemplate: "cn-hk-1",
        description: "适用于网站、应用与中间件场景的通用云服务器。",
      },
    }),
    db.product.create({
      data: {
        code: "CLOUD-SEC-4C8G",
        name: "高防云节点 4C8G",
        category: ProductCategory.SECURITY,
        status: ProductStatus.ACTIVE,
        billingCycle: BillingCycle.MONTHLY,
        price: 168_800,
        setupFee: 9_900,
        stock: 36,
        autoProvision: true,
        providerType: ProviderType.MOFANG_CLOUD,
        providerProductId: "shield-4c8g",
        regionTemplate: "cn-bgp-1",
        description: "适用于高防和 BGP 业务场景的安全型云实例。",
      },
    }),
    db.product.create({
      data: {
        code: "BM-PRO-8C32G",
        name: "裸金属性能型 8C32G",
        category: ProductCategory.BARE_METAL,
        status: ProductStatus.ACTIVE,
        billingCycle: BillingCycle.MONTHLY,
        price: 399_000,
        setupFee: 15_000,
        stock: 12,
        autoProvision: false,
        providerType: ProviderType.MANUAL,
        providerProductId: "bm-pro-8c32g",
        regionTemplate: "cn-sh-1",
        description: "适用于数据库和高 IO 业务的人工交付裸金属产品。",
      },
    }),
    db.product.create({
      data: {
        code: "OBJ-PLUS-5TB",
        name: "对象存储增强版 5TB",
        category: ProductCategory.STORAGE,
        status: ProductStatus.ACTIVE,
        billingCycle: BillingCycle.MONTHLY,
        price: 238_000,
        setupFee: 0,
        stock: 80,
        autoProvision: true,
        providerType: ProviderType.MOFANG_CLOUD,
        providerProductId: "obj-plus-5tb",
        regionTemplate: "cn-sz-1",
        description: "支持生命周期管理的高可靠对象存储方案。",
      },
    }),
    db.product.create({
      data: {
        code: "NET-EIP-50M",
        name: "弹性公网 IP 50M",
        category: ProductCategory.NETWORK,
        status: ProductStatus.ACTIVE,
        billingCycle: BillingCycle.MONTHLY,
        price: 19_900,
        setupFee: 0,
        stock: 300,
        autoProvision: true,
        providerType: ProviderType.MOFANG_CLOUD,
        providerProductId: "eip-50m",
        regionTemplate: "cn-hk-1",
        description: "附带 50 Mbps 带宽包的公网 IP 资源。",
      },
    }),
  ]);

  const [regionHK, regionBGP, regionSZ, flavorA, flavorB, flavorC, imageA, imageB, imageC] =
    await Promise.all([
      db.cloudRegion.create({
        data: {
          code: "cn-hk-1",
          name: "中国香港一区",
          location: "中国香港",
          providerType: ProviderType.MOFANG_CLOUD,
          providerRegionId: "mf-hk-1",
          description: "低时延国际业务地域。",
          sortOrder: 10,
        },
      }),
      db.cloudRegion.create({
        data: {
          code: "cn-bgp-1",
          name: "华南高防 BGP",
          location: "中国深圳",
          providerType: ProviderType.MOFANG_CLOUD,
          providerRegionId: "mf-bgp-1",
          description: "高防和 BGP 线路专用地域。",
          sortOrder: 20,
        },
      }),
      db.cloudRegion.create({
        data: {
          code: "cn-sz-1",
          name: "深圳存储一区",
          location: "中国深圳",
          providerType: ProviderType.MOFANG_CLOUD,
          providerRegionId: "mf-sz-1",
          description: "对象存储与备份业务主地域。",
          sortOrder: 30,
        },
      }),
      db.cloudFlavor.create({
        data: {
          code: "GEN2-2C4G",
          name: "通用型 2C4G",
          category: ProductCategory.CLOUD_SERVER,
          cpu: 2,
          memoryGb: 4,
          storageGb: 80,
          bandwidthMbps: 30,
          description: "适合网站和中间件场景。",
          sortOrder: 10,
        },
      }),
      db.cloudFlavor.create({
        data: {
          code: "SEC-4C8G",
          name: "高防型 4C8G",
          category: ProductCategory.SECURITY,
          cpu: 4,
          memoryGb: 8,
          storageGb: 120,
          bandwidthMbps: 50,
          description: "适合高防节点和 BGP 业务。",
          sortOrder: 20,
        },
      }),
      db.cloudFlavor.create({
        data: {
          code: "BM-8C32G",
          name: "裸金属 8C32G",
          category: ProductCategory.BARE_METAL,
          cpu: 8,
          memoryGb: 32,
          storageGb: 960,
          bandwidthMbps: 100,
          description: "数据库和高 IO 业务规格。",
          sortOrder: 30,
        },
      }),
      db.cloudImage.create({
        data: {
          code: "ubuntu-2204",
          name: "Ubuntu 22.04 LTS",
          osType: "Linux",
          version: "22.04",
          architecture: "x86_64",
          description: "默认 Linux 公共镜像。",
          isPublic: true,
        },
      }),
      db.cloudImage.create({
        data: {
          code: "debian-12",
          name: "Debian 12",
          osType: "Linux",
          version: "12",
          architecture: "x86_64",
          description: "适用于通用业务节点。",
          isPublic: true,
        },
      }),
      db.cloudImage.create({
        data: {
          code: "windows-2022",
          name: "Windows Server 2022",
          osType: "Windows",
          version: "2022",
          architecture: "x86_64",
          description: "适用于 Windows 业务实例。",
          isPublic: false,
        },
      }),
    ]);

  const [zoneHKA, zoneBGPA, zoneSZA] = await Promise.all([
    db.cloudZone.create({
      data: {
        regionId: regionHK.id,
        code: "hk-a",
        name: "香港可用区 A",
        providerZoneId: "mf-hk-a",
        sortOrder: 10,
      },
    }),
    db.cloudZone.create({
      data: {
        regionId: regionBGP.id,
        code: "bgp-a",
        name: "BGP 可用区 A",
        providerZoneId: "mf-bgp-a",
        sortOrder: 10,
      },
    }),
    db.cloudZone.create({
      data: {
        regionId: regionSZ.id,
        code: "sz-a",
        name: "深圳可用区 A",
        providerZoneId: "mf-sz-a",
        sortOrder: 10,
      },
    }),
  ]);

  const [planA, planB, planC, planD] = await Promise.all([
    db.cloudPlan.create({
      data: {
        productId: productA.id,
        regionId: regionHK.id,
        zoneId: zoneHKA.id,
        flavorId: flavorA.id,
        imageId: imageA.id,
        code: "HK-GEN2-2C4G",
        name: "香港通用型 2C4G",
        billingCycle: BillingCycle.MONTHLY,
        salePrice: 68_800,
        marketPrice: 88_000,
        stock: 120,
        isPublic: true,
        isActive: true,
        configOptions: JSON.stringify({
          node: "hk-node-01",
          os: "ubuntu-2204",
          cpu: 2,
          memory: 4,
          system_disk_size: 80,
          network_type: "public",
          bw: 30,
          flow_way: "month",
          ip_num: 1,
        }),
        description: "适合网站、应用和中间件场景。",
      },
    }),
    db.cloudPlan.create({
      data: {
        productId: productB.id,
        regionId: regionBGP.id,
        zoneId: zoneBGPA.id,
        flavorId: flavorB.id,
        imageId: imageB.id,
        code: "BGP-SEC-4C8G",
        name: "华南高防 4C8G",
        billingCycle: BillingCycle.MONTHLY,
        salePrice: 168_800,
        marketPrice: 198_800,
        setupFee: 9_900,
        stock: 36,
        isPublic: true,
        isActive: true,
        configOptions: JSON.stringify({
          node: "bgp-shield-01",
          os: "debian-12",
          cpu: 4,
          memory: 8,
          system_disk_size: 120,
          network_type: "bgp",
          bw: 50,
          flow_way: "month",
          ip_num: 1,
          peak_defence: 80,
        }),
        description: "适合高防业务、BGP 节点和安全网关。",
      },
    }),
    db.cloudPlan.create({
      data: {
        productId: productC.id,
        regionId: regionBGP.id,
        zoneId: zoneBGPA.id,
        flavorId: flavorC.id,
        imageId: imageC.id,
        code: "BGP-BM-8C32G",
        name: "高防裸金属 8C32G",
        billingCycle: BillingCycle.MONTHLY,
        salePrice: 399_000,
        marketPrice: 429_000,
        setupFee: 15_000,
        stock: 12,
        isPublic: false,
        isActive: true,
        configOptions: JSON.stringify({
          node: "bm-bgp-01",
          os: "windows-2022",
          cpu: 8,
          memory: 32,
          system_disk_size: 240,
          data_disk_size: 960,
          network_type: "bgp",
          bw: 100,
          flow_way: "month",
          ip_num: 1,
          peak_defence: 120,
        }),
        description: "适合数据库和高 IO 裸金属业务。",
      },
    }),
    db.cloudPlan.create({
      data: {
        productId: productD.id,
        regionId: regionSZ.id,
        zoneId: zoneSZA.id,
        imageId: imageA.id,
        code: "SZ-OBJ-5TB",
        name: "深圳对象存储 5TB",
        billingCycle: BillingCycle.MONTHLY,
        salePrice: 238_000,
        marketPrice: 258_000,
        stock: 80,
        isPublic: true,
        isActive: true,
        configOptions: JSON.stringify({
          node: "sz-storage-01",
          os: "ubuntu-2204",
          system_disk_size: 50,
          data_disk_size: 5120,
          network_type: "storage",
          bw: 20,
          flow_way: "month",
          ip_num: 1,
        }),
        description: "适合对象存储、归档与备份场景。",
      },
    }),
  ]);

  const [vpcA, vpcB] = await Promise.all([
    db.serviceVpcNetwork.create({
      data: {
        name: "SG Production VPC",
        region: "cn-hk-1",
        cidr: "10.21.0.0/16",
        gateway: "10.21.0.1",
        providerVpcId: "mf-vpc-001",
      },
    }),
    db.serviceVpcNetwork.create({
      data: {
        name: "BGP Security VPC",
        region: "cn-bgp-1",
        cidr: "10.66.0.0/16",
        gateway: "10.66.0.1",
        providerVpcId: "mf-vpc-002",
      },
    }),
  ]);

  const orderA = await db.order.create({
    data: {
      orderNo: "ORD-20260319-001",
      customerId: customerA.id,
      status: OrderStatus.ACTIVE,
      totalAmount: 137_600,
      paidAmount: 137_600,
      currency: "CNY",
      source: "portal",
      orderType: "new",
      dueDate: daysLater(7),
      paidAt: daysAgo(22),
      notes: "Production cluster initial order.",
    },
  });

  const orderB = await db.order.create({
    data: {
      orderNo: "ORD-20260319-002",
      customerId: customerB.id,
      status: OrderStatus.ACTIVE,
      totalAmount: 238_000,
      paidAmount: 238_000,
      currency: "CNY",
      source: "portal",
      orderType: "new",
      dueDate: daysLater(2),
      paidAt: daysAgo(30),
      notes: "Storage service launched last month.",
    },
  });

  const orderC = await db.order.create({
    data: {
      orderNo: "ORD-20260319-003",
      customerId: customerC.id,
      status: OrderStatus.ACTIVE,
      totalAmount: 178_700,
      paidAmount: 178_700,
      currency: "CNY",
      source: "sales",
      orderType: "new",
      dueDate: daysLater(4),
      paidAt: daysAgo(25),
      notes: "Reseller security node opening order.",
    },
  });

  const orderD = await db.order.create({
    data: {
      orderNo: "ORD-20260319-004",
      customerId: customerD.id,
      status: OrderStatus.ACTIVE,
      totalAmount: 414_000,
      paidAmount: 414_000,
      currency: "CNY",
      source: "admin",
      orderType: "new",
      dueDate: daysAgo(2),
      paidAt: daysAgo(45),
      notes: "Bare metal lab order created by operations.",
    },
  });

  const serviceA = await db.serviceInstance.create({
    data: {
      serviceNo: "SRV-202603-001",
      customerId: customerA.id,
      productId: productA.id,
      planId: planA.id,
      orderId: orderA.id,
      vpcNetworkId: vpcA.id,
      name: "StarGalaxy Web 集群 A",
      hostname: "web-a.stargalaxy.cn",
      providerType: ProviderType.MOFANG_CLOUD,
      providerResourceId: "mf-web-001",
      region: "cn-hk-1",
      billingCycle: BillingCycle.MONTHLY,
      status: ServiceStatus.ACTIVE,
      ipAddress: "203.0.113.21",
      cpuCores: 2,
      memoryGb: 4,
      storageGb: 80,
      monthlyCost: 68_800,
      nextDueDate: daysLater(7),
      activatedAt: daysAgo(22),
      configSnapshot: JSON.stringify({
        os: "Ubuntu 22.04",
        backupPolicy: "daily",
        line: "BGP",
      }),
      lastSyncAt: daysAgo(1),
    },
  });

  const serviceB = await db.serviceInstance.create({
    data: {
      serviceNo: "SRV-202603-002",
      customerId: customerB.id,
      productId: productD.id,
      planId: planD.id,
      orderId: orderB.id,
      name: "CloudMatrix 对象存储",
      hostname: "storage.cloudmatrix.cn",
      providerType: ProviderType.MOFANG_CLOUD,
      providerResourceId: "mf-obj-003",
      region: "cn-sz-1",
      billingCycle: BillingCycle.MONTHLY,
      status: ServiceStatus.OVERDUE,
      storageGb: 5_000,
      monthlyCost: 238_000,
      nextDueDate: daysAgo(10),
      activatedAt: daysAgo(30),
      configSnapshot: JSON.stringify({
        bucket: "matrix-assets",
        versioning: true,
        archiveTier: true,
      }),
      lastSyncAt: daysAgo(2),
    },
  });

  const serviceC = await db.serviceInstance.create({
    data: {
      serviceNo: "SRV-202603-003",
      customerId: customerC.id,
      productId: productB.id,
      planId: planB.id,
      orderId: orderC.id,
      vpcNetworkId: vpcB.id,
      name: "NorthShore 高防节点",
      hostname: "shield.nshore.cn",
      providerType: ProviderType.MOFANG_CLOUD,
      providerResourceId: "mf-ddos-002",
      region: "cn-bgp-1",
      billingCycle: BillingCycle.MONTHLY,
      status: ServiceStatus.ACTIVE,
      ipAddress: "198.51.100.23",
      cpuCores: 4,
      memoryGb: 8,
      storageGb: 120,
      monthlyCost: 168_800,
      nextDueDate: daysLater(4),
      activatedAt: daysAgo(25),
      configSnapshot: JSON.stringify({
        antiDdos: "50G",
        line: "BGP",
        monitor: "standard",
      }),
      lastSyncAt: daysAgo(1),
    },
  });

  const serviceD = await db.serviceInstance.create({
    data: {
      serviceNo: "SRV-202603-004",
      customerId: customerD.id,
      productId: productC.id,
      planId: planC.id,
      orderId: orderD.id,
      name: "WuFeng 数据库实验机",
      hostname: "db-lab.wufeng.dev",
      providerType: ProviderType.MANUAL,
      region: "cn-sh-1",
      billingCycle: BillingCycle.MONTHLY,
      status: ServiceStatus.SUSPENDED,
      cpuCores: 8,
      memoryGb: 32,
      storageGb: 960,
      monthlyCost: 399_000,
      nextDueDate: daysAgo(2),
      activatedAt: daysAgo(45),
      configSnapshot: JSON.stringify({
        rack: "R12",
        uplink: "10Gbps",
        raid: "RAID10",
      }),
      lastSyncAt: daysAgo(1),
    },
  });

  await Promise.all([
    db.orderItem.create({
      data: {
        orderId: orderA.id,
        productId: productA.id,
        planId: planA.id,
        serviceId: serviceA.id,
        title: productA.name,
        quantity: 2,
        unitPrice: 68_800,
        cycle: BillingCycle.MONTHLY,
        totalAmount: 137_600,
      },
    }),
    db.orderItem.create({
      data: {
        orderId: orderB.id,
        productId: productD.id,
        planId: planD.id,
        serviceId: serviceB.id,
        title: productD.name,
        quantity: 1,
        unitPrice: 238_000,
        cycle: BillingCycle.MONTHLY,
        totalAmount: 238_000,
      },
    }),
    db.orderItem.create({
      data: {
        orderId: orderC.id,
        productId: productB.id,
        planId: planB.id,
        serviceId: serviceC.id,
        title: productB.name,
        quantity: 1,
        unitPrice: 168_800,
        cycle: BillingCycle.MONTHLY,
        totalAmount: 178_700,
      },
    }),
    db.orderItem.create({
      data: {
        orderId: orderD.id,
        productId: productC.id,
        planId: planC.id,
        serviceId: serviceD.id,
        title: productC.name,
        quantity: 1,
        unitPrice: 399_000,
        cycle: BillingCycle.MONTHLY,
        totalAmount: 414_000,
      },
    }),
  ]);

  const invoiceA = await db.invoice.create({
    data: {
      invoiceNo: "INV-20260319-001",
      customerId: customerA.id,
      orderId: orderA.id,
      serviceId: serviceA.id,
      type: InvoiceType.ORDER,
      status: InvoiceStatus.PAID,
      subtotal: 137_600,
      taxRate: 6,
      taxAmount: 8_256,
      totalAmount: 145_856,
      paidAmount: 145_856,
      dueDate: daysAgo(20),
      issuedAt: daysAgo(23),
      paidAt: daysAgo(22),
      remark: "Initial order invoice paid by Alipay.",
    },
  });

  const invoiceB = await db.invoice.create({
    data: {
      invoiceNo: "INV-20260319-002",
      customerId: customerB.id,
      orderId: orderB.id,
      serviceId: serviceB.id,
      type: InvoiceType.RENEWAL,
      status: InvoiceStatus.OVERDUE,
      subtotal: 238_000,
      taxRate: 0,
      taxAmount: 0,
      totalAmount: 238_000,
      paidAmount: 0,
      dueDate: daysAgo(10),
      issuedAt: daysAgo(17),
      remark: "Storage renewal invoice overdue.",
    },
  });

  const invoiceC = await db.invoice.create({
    data: {
      invoiceNo: "INV-20260319-003",
      customerId: customerD.id,
      orderId: orderD.id,
      serviceId: serviceD.id,
      type: InvoiceType.RENEWAL,
      status: InvoiceStatus.PARTIAL,
      subtotal: 399_000,
      taxRate: 0,
      taxAmount: 0,
      totalAmount: 399_000,
      paidAmount: 100_000,
      dueDate: daysAgo(2),
      issuedAt: daysAgo(8),
      remark: "Lab server renewal invoice partially paid.",
    },
  });

  const invoiceD = await db.invoice.create({
    data: {
      invoiceNo: "INV-20260319-004",
      customerId: customerC.id,
      orderId: orderC.id,
      serviceId: serviceC.id,
      type: InvoiceType.MANUAL,
      status: InvoiceStatus.PAID,
      subtotal: 9_900,
      taxRate: 0,
      taxAmount: 0,
      totalAmount: 9_900,
      paidAmount: 9_900,
      dueDate: daysAgo(25),
      issuedAt: daysAgo(25),
      paidAt: daysAgo(25),
      remark: "Initial setup fee invoice.",
    },
  });

  await Promise.all([
    db.payment.create({
      data: {
        paymentNo: "PAY-20260319-001",
        customerId: customerA.id,
        invoiceId: invoiceA.id,
        orderId: orderA.id,
        method: PaymentMethod.ALIPAY,
        status: PaymentStatus.SUCCESS,
        amount: 145_856,
        transactionNo: "ALI202603190001",
        paidAt: daysAgo(22),
      },
    }),
    db.payment.create({
      data: {
        paymentNo: "PAY-20260319-002",
        customerId: customerD.id,
        invoiceId: invoiceC.id,
        orderId: orderD.id,
        method: PaymentMethod.BANK_TRANSFER,
        status: PaymentStatus.SUCCESS,
        amount: 100_000,
        transactionNo: "BANK202603190003",
        paidAt: daysAgo(1),
      },
    }),
    db.payment.create({
      data: {
        paymentNo: "PAY-20260319-003",
        customerId: customerC.id,
        invoiceId: invoiceD.id,
        orderId: orderC.id,
        method: PaymentMethod.OFFLINE,
        status: PaymentStatus.SUCCESS,
        amount: 9_900,
        transactionNo: "OFF202603190004",
        paidAt: daysAgo(25),
      },
    }),
  ]);

  const seededPayments = await db.payment.findMany({
    orderBy: {
      createdAt: "asc",
    },
  });

  await db.refundRecord.create({
    data: {
      refundNo: "REF-20260319-001",
      paymentId: seededPayments[0]!.id,
      invoiceId: invoiceA.id,
      customerId: customerA.id,
      processedById: finance.id,
      amount: 5_000,
      method: PaymentMethod.BALANCE,
      status: RefundStatus.PENDING,
      reason: "客户申请核减测试费用",
      payload: JSON.stringify({ originalPaymentNo: seededPayments[0]!.paymentNo }),
    },
  });

  await db.paymentCallbackLog.createMany({
    data: [
      {
        paymentId: seededPayments[0]?.id,
        method: PaymentMethod.ALIPAY,
        invoiceNo: invoiceA.invoiceNo,
        paymentNo: seededPayments[0]?.paymentNo,
        transactionNo: seededPayments[0]?.transactionNo,
        callbackStatus: "SUCCESS",
        payload: JSON.stringify({ tradeStatus: "TRADE_SUCCESS", amount: 145856 }),
        message: "官方支付回调已受理。",
        isHandled: true,
        handledAt: daysAgo(22),
      },
      {
        method: PaymentMethod.WECHAT,
        invoiceNo: invoiceB.invoiceNo,
        transactionNo: "WX202603190008",
        callbackStatus: "PENDING",
        payload: JSON.stringify({ tradeState: "USERPAYING", amount: 238000 }),
        message: "等待二次回调确认。",
        isHandled: false,
      },
    ],
  });

  const ipA = await db.serviceIpAddress.create({
    data: {
      serviceId: serviceA.id,
      address: "203.0.113.21",
      version: "IPv4",
      type: "primary",
      bandwidthMbps: 30,
      providerIpId: "mf-ip-001",
      isPrimary: true,
    },
  });

  await db.serviceIpAddress.create({
    data: {
      serviceId: serviceA.id,
      address: "2001:db8:21::21",
      version: "IPv6",
      type: "primary",
      providerIpId: "mf-ip6-001",
      isPrimary: false,
    },
  });

  await db.serviceIpAddress.create({
    data: {
      serviceId: serviceC.id,
      address: "198.51.100.23",
      version: "IPv4",
      type: "primary",
      bandwidthMbps: 50,
      providerIpId: "mf-ip-002",
      isPrimary: true,
    },
  });

  await db.serviceIpAddress.create({
    data: {
      serviceId: serviceD.id,
      address: "203.0.113.88",
      version: "IPv4",
      type: "primary",
      bandwidthMbps: 100,
      isPrimary: true,
      status: "SUSPENDED",
    },
  });

  await db.serviceDisk.create({
    data: {
      serviceId: serviceA.id,
      name: "system",
      type: "SSD",
      sizeGb: 80,
      mountPoint: "/",
      providerDiskId: "mf-disk-001",
      isSystem: true,
    },
  });

  const diskA2 = await db.serviceDisk.create({
    data: {
      serviceId: serviceA.id,
      name: "data",
      type: "SSD",
      sizeGb: 200,
      mountPoint: "/data",
      providerDiskId: "mf-disk-002",
      isSystem: false,
    },
  });

  const diskC1 = await db.serviceDisk.create({
    data: {
      serviceId: serviceC.id,
      name: "system",
      type: "SSD",
      sizeGb: 120,
      mountPoint: "/",
      providerDiskId: "mf-disk-003",
      isSystem: true,
    },
  });

  const diskD1 = await db.serviceDisk.create({
    data: {
      serviceId: serviceD.id,
      name: "raid-data",
      type: "SAS",
      sizeGb: 960,
      mountPoint: "/data",
      isSystem: false,
      status: "SUSPENDED",
    },
  });

  await Promise.all([
    db.serviceSnapshot.create({
      data: {
        serviceId: serviceA.id,
        sourceDiskId: diskA2.id,
        name: "web-data-20260317",
        providerSnapshotId: "mf-snap-001",
        sizeGb: 200,
      },
    }),
    db.serviceSnapshot.create({
      data: {
        serviceId: serviceC.id,
        sourceDiskId: diskC1.id,
        name: "shield-root-20260318",
        providerSnapshotId: "mf-snap-002",
        sizeGb: 120,
      },
    }),
    db.serviceSnapshot.create({
      data: {
        serviceId: serviceD.id,
        sourceDiskId: diskD1.id,
        name: "db-lab-prepatch",
        sizeGb: 960,
        status: "ARCHIVED",
      },
    }),
  ]);

  await Promise.all([
    db.serviceBackup.create({
      data: {
        serviceId: serviceA.id,
        name: "web-cluster-nightly",
        providerBackupId: "mf-bak-001",
        sizeGb: 45,
        expiresAt: daysLater(6),
      },
    }),
    db.serviceBackup.create({
      data: {
        serviceId: serviceB.id,
        name: "storage-bucket-weekly",
        providerBackupId: "mf-bak-002",
        sizeGb: 320,
        expiresAt: daysLater(20),
      },
    }),
    db.serviceBackup.create({
      data: {
        serviceId: serviceD.id,
        name: "db-lab-manual",
        sizeGb: 88,
        expiresAt: daysLater(3),
        status: "READY",
      },
    }),
  ]);

  const sgA = await db.serviceSecurityGroup.create({
    data: {
      serviceId: serviceA.id,
      vpcNetworkId: vpcA.id,
      name: "web-cluster-sg",
      providerSecurityGroupId: "mf-sg-001",
    },
  });

  const sgB = await db.serviceSecurityGroup.create({
    data: {
      serviceId: serviceC.id,
      vpcNetworkId: vpcB.id,
      name: "shield-node-sg",
      providerSecurityGroupId: "mf-sg-002",
    },
  });

  await db.serviceSecurityGroup.create({
    data: {
      vpcNetworkId: vpcB.id,
      name: "bgp-vpc-default",
      providerSecurityGroupId: "mf-sg-003",
    },
  });

  await db.serviceSecurityGroupRule.createMany({
    data: [
      {
        securityGroupId: sgA.id,
        direction: "ingress",
        protocol: "tcp",
        portRange: "80,443",
        sourceCidr: "0.0.0.0/0",
        description: "HTTP and HTTPS",
      },
      {
        securityGroupId: sgA.id,
        direction: "ingress",
        protocol: "tcp",
        portRange: "22",
        sourceCidr: "10.21.0.0/16",
        description: "Ops SSH within VPC",
      },
      {
        securityGroupId: sgB.id,
        direction: "ingress",
        protocol: "tcp",
        portRange: "80,443",
        sourceCidr: "0.0.0.0/0",
        description: "Public service traffic",
      },
      {
        securityGroupId: sgB.id,
        direction: "ingress",
        protocol: "udp",
        portRange: "1-65535",
        sourceCidr: "0.0.0.0/0",
        description: "Game shield traffic",
      },
    ],
  });

  await db.creditTransaction.createMany({
    data: [
      {
        customerId: customerA.id,
        operatorId: finance.id,
        type: CreditTransactionType.RECHARGE,
        amount: 3_000_000,
        balanceAfter: 3_000_000,
        description: "银行转账充值",
      },
      {
        customerId: customerA.id,
        operatorId: finance.id,
        type: CreditTransactionType.CONSUME,
        amount: -200_000,
        balanceAfter: 2_800_000,
        description: `服务费用结算 ${invoiceA.invoiceNo}`,
      },
      {
        customerId: customerB.id,
        operatorId: finance.id,
        type: CreditTransactionType.RECHARGE,
        amount: 1_800_000,
        balanceAfter: 1_800_000,
        description: "季度余额充值",
      },
      {
        customerId: customerB.id,
        operatorId: finance.id,
        type: CreditTransactionType.CONSUME,
        amount: -240_000,
        balanceAfter: 1_560_000,
        description: "上一期续费从余额中扣除",
      },
      {
        customerId: customerC.id,
        operatorId: admin.id,
        type: CreditTransactionType.RECHARGE,
        amount: 700_000,
        balanceAfter: 700_000,
        description: "代理商预存款",
      },
      {
        customerId: customerC.id,
        operatorId: admin.id,
        type: CreditTransactionType.ADJUSTMENT,
        amount: -180_000,
        balanceAfter: 520_000,
        description: "渠道返点结算",
      },
      {
        customerId: customerD.id,
        operatorId: finance.id,
        type: CreditTransactionType.RECHARGE,
        amount: 200_000,
        balanceAfter: 200_000,
        description: "个人账户充值",
      },
      {
        customerId: customerD.id,
        operatorId: finance.id,
        type: CreditTransactionType.CONSUME,
        amount: -80_000,
        balanceAfter: 120_000,
        description: "服务部分续费扣款",
      },
    ],
  });

  const ticketA = await db.ticket.create({
    data: {
      ticketNo: "TIC-20260319-001",
      customerId: customerC.id,
      serviceId: serviceC.id,
      assignedToId: support.id,
      subject: "代理节点续费后状态未同步",
      priority: TicketPriority.URGENT,
      status: TicketStatus.PROCESSING,
      summary: "客户已支付开通费用，当前希望确认下一周期续费窗口是否已同步。",
      lastReplyAt: daysAgo(1),
    },
  });

  const ticketB = await db.ticket.create({
    data: {
      ticketNo: "TIC-20260319-002",
      customerId: customerA.id,
      serviceId: serviceA.id,
      assignedToId: operations.id,
      subject: "Web 集群数据盘扩容到 300G",
      priority: TicketPriority.HIGH,
      status: TicketStatus.OPEN,
      summary: "客户申请扩容数据盘，并需要同步生成对应账单。",
      lastReplyAt: daysAgo(1),
    },
  });

  await db.ticket.create({
    data: {
      ticketNo: "TIC-20260319-003",
      customerId: customerB.id,
      serviceId: serviceB.id,
      assignedToId: finance.id,
      subject: "逾期账单催缴跟进",
      priority: TicketPriority.NORMAL,
      status: TicketStatus.WAITING_CUSTOMER,
      summary: "财务已对客户逾期的存储续费账单进行催缴提醒。",
      lastReplyAt: daysAgo(3),
    },
  });

  await db.ticketReply.createMany({
    data: [
      {
        ticketId: ticketA.id,
        authorType: "CUSTOMER",
        authorName: customerC.name,
        content: "请确认代理商安全节点是否已经进入下一计费周期。",
      },
      {
        ticketId: ticketA.id,
        authorType: "STAFF",
        authorName: support.name,
        content: "运维正在核查云平台回调与计费状态。",
      },
      {
        ticketId: ticketB.id,
        authorType: "CUSTOMER",
        authorName: customerA.name,
        content: "我们需要在周末发布前完成应用数据盘扩容。",
      },
      {
        ticketId: ticketB.id,
        authorType: "STAFF",
        authorName: operations.name,
        content: "内部备注：请准备扩容报价和维护窗口。",
        isInternal: true,
      },
    ],
  });

  await db.providerSyncLog.createMany({
    data: [
      {
        providerType: ProviderType.MOFANG_CLOUD,
        action: "sync",
        resourceType: "instance",
        resourceId: serviceA.providerResourceId,
        status: ProviderSyncStatus.SUCCESS,
        message: "实例状态与公网 IP 已完成同步。",
        requestBody: JSON.stringify({ localService: serviceA.serviceNo }),
        responseBody: JSON.stringify({ status: "ACTIVE", ip: ipA.address }),
      },
      {
        providerType: ProviderType.MOFANG_CLOUD,
        action: "suspend",
        resourceType: "instance",
        resourceId: serviceB.providerResourceId,
        status: ProviderSyncStatus.SUCCESS,
        message: "续费账单逾期，服务已暂停。",
        requestBody: JSON.stringify({ invoiceNo: invoiceB.invoiceNo }),
        responseBody: JSON.stringify({ status: "SUSPENDED", reason: "OVERDUE" }),
      },
      {
        providerType: ProviderType.MOFANG_CLOUD,
        action: "sync",
        resourceType: "instance",
        resourceId: serviceC.providerResourceId,
        status: ProviderSyncStatus.SUCCESS,
        message: "安全节点的 CPU、内存与状态已刷新。",
        requestBody: JSON.stringify({ localService: serviceC.serviceNo }),
        responseBody: JSON.stringify({ status: "ACTIVE", cpu: 4, memoryGb: 8 }),
      },
      {
        providerType: ProviderType.MOFANG_CLOUD,
        action: "backup",
        resourceType: "instance",
        resourceId: serviceB.providerResourceId,
        status: ProviderSyncStatus.PENDING,
        message: "备份任务已在云平台侧排队执行。",
        requestBody: JSON.stringify({ backup: "storage-bucket-weekly" }),
        responseBody: JSON.stringify({ jobStatus: "PENDING" }),
      },
    ],
  });

  await db.billingJob.createMany({
    data: [
      {
        jobType: "MARK_OVERDUE",
        status: TaskRunStatus.SUCCESS,
        customerId: customerB.id,
        serviceId: serviceB.id,
        invoiceId: invoiceB.id,
        message: "续费账单已标记逾期，服务状态已切换为逾期。",
        payload: JSON.stringify({ invoiceNo: invoiceB.invoiceNo }),
      },
      {
        jobType: "AUTO_SUSPEND",
        status: TaskRunStatus.SUCCESS,
        customerId: customerD.id,
        serviceId: serviceD.id,
        invoiceId: invoiceC.id,
        message: "裸金属服务在部分支付后仍逾期，已自动暂停。",
        payload: JSON.stringify({ invoiceNo: invoiceC.invoiceNo }),
      },
      {
        jobType: "RESOURCE_SYNC",
        status: TaskRunStatus.SUCCESS,
        customerId: customerA.id,
        serviceId: serviceA.id,
        message: "云平台实例信息已同步到本地资源台账。",
        payload: JSON.stringify({ providerResourceId: serviceA.providerResourceId }),
      },
    ],
  });

  await db.auditLog.createMany({
    data: [
      {
        adminUserId: admin.id,
        module: "auth",
        action: "login",
        targetType: "adminUser",
        targetId: admin.id,
        summary: "超级管理员从后台完成登录。",
      },
      {
        adminUserId: finance.id,
        module: "billing",
        action: "mark-overdue",
        targetType: "invoice",
        targetId: invoiceB.id,
        summary: "财务将存储续费账单标记为逾期。",
        detail: JSON.stringify({ invoiceNo: invoiceB.invoiceNo, customer: customerB.name }),
      },
      {
        adminUserId: operations.id,
        module: "service",
        action: "suspend",
        targetType: "service",
        targetId: serviceD.id,
        summary: "运维暂停了一台人工开通的裸金属服务。",
        detail: JSON.stringify({ serviceNo: serviceD.serviceNo, reason: "overdue renewal" }),
      },
      {
        adminUserId: support.id,
        module: "ticket",
        action: "reply",
        targetType: "ticket",
        targetId: ticketA.id,
        summary: "客服回复了代理商续费同步工单。",
        detail: JSON.stringify({ ticketNo: ticketA.ticketNo }),
      },
    ],
  });

  const seededTemplates = await db.notificationTemplate.findMany({
    orderBy: {
      createdAt: "asc",
    },
  });

  const messageA = await db.notificationMessage.create({
    data: {
      templateId: seededTemplates.find((item) => item.code === "PAYMENT_SUCCESS")?.id,
      customerId: customerA.id,
      createdById: finance.id,
      channel: NotificationChannel.SYSTEM,
      priority: NotificationPriority.NORMAL,
      status: NotificationStatus.SENT,
      recipient: customerA.email,
      recipientName: customerA.name,
      subject: "账单 INV-20260319-001 收款成功",
      content: "账单 INV-20260319-001 已完成收款，到账金额 1458.56 元。",
      module: "payment",
      relatedType: "invoice",
      relatedId: invoiceA.id,
      sentAt: daysAgo(22),
    },
  });

  const messageB = await db.notificationMessage.create({
    data: {
      templateId: seededTemplates.find((item) => item.code === "ORDER_CREATED")?.id,
      customerId: customerB.id,
      createdById: admin.id,
      channel: NotificationChannel.EMAIL,
      priority: NotificationPriority.HIGH,
      status: NotificationStatus.PENDING,
      recipient: customerB.email,
      recipientName: customerB.name,
      subject: "订单 ORD-20260319-002 已创建",
      content: "订单 ORD-20260319-002 已创建，待您完成续费账单支付。",
      module: "order",
      relatedType: "order",
      relatedId: orderB.id,
    },
  });

  await db.asyncTaskJob.createMany({
    data: [
      {
        queueName: "notifications",
        jobType: "SEND_SYSTEM",
        module: "payment",
        status: TaskRunStatus.SUCCESS,
        targetType: "notification",
        targetId: messageA.id,
        notificationId: messageA.id,
        result: "站内通知已推送给 ops@stargalaxy.cn",
        attempts: 1,
        availableAt: daysAgo(22),
        executedAt: daysAgo(22),
      },
      {
        queueName: "notifications",
        jobType: "SEND_EMAIL",
        module: "order",
        status: TaskRunStatus.PENDING,
        targetType: "notification",
        targetId: messageB.id,
        notificationId: messageB.id,
        payload: JSON.stringify({ recipient: customerB.email }),
        attempts: 0,
        availableAt: new Date(),
      },
    ],
  });

  console.log("Seed completed");
  console.log(`Admin account: ${process.env.SEED_ADMIN_EMAIL ?? "admin@idc.local"}`);
  console.log(`Admin password: ${process.env.SEED_ADMIN_PASSWORD ?? "Admin123!"}`);
  console.log("Finance account: finance@idc.local");
  console.log("Support account: support@idc.local");
  console.log("Operations account: ops@idc.local");
  console.log(`Portal password: ${process.env.SEED_PORTAL_PASSWORD ?? "Portal123!"}`);
  console.log("Portal account: ops@stargalaxy.cn");
  console.log("Portal account: finance@cloudmatrix.cn");
}

main()
  .catch((error) => {
    console.error(error);
    process.exit(1);
  })
  .finally(async () => {
    await db.$disconnect();
  });
