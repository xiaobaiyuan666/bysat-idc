<script setup lang="ts">
import { ArrowLeft, Edit, Money, Plus } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";

import { http } from "@/api/http";
import { useAuthStore } from "@/stores/auth";
import { formatCurrency, formatDate, formatDateTime } from "@/utils/format";

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const loading = ref(false);
const detail = ref<any | null>(null);

const customerDialogVisible = ref(false);
const balanceDialogVisible = ref(false);
const contactDialogVisible = ref(false);
const followUpDialogVisible = ref(false);
const certificationSubmitting = ref(false);

const customerForm = reactive({
  name: "",
  companyName: "",
  email: "",
  phone: "",
  contactQQ: "",
  contactWechat: "",
  type: "COMPANY",
  status: "ACTIVE",
  level: "standard",
  tags: "",
  notes: "",
});

const balanceForm = reactive({
  amount: 0,
  description: "",
});

const contactForm = reactive({
  id: "",
  name: "",
  department: "",
  role: "",
  email: "",
  phone: "",
  isPrimary: false,
  notes: "",
});

const certificationForm = reactive({
  subjectType: "PERSONAL",
  subjectName: "",
  idNumber: "",
  businessLicenseNo: "",
  status: "PENDING",
  submittedAt: "",
  reviewNote: "",
});

const followUpForm = reactive({
  id: "",
  type: "NOTE",
  title: "",
  content: "",
  nextFollowAt: "",
});

const canManageCustomers = computed(() =>
  Boolean(authStore.user?.permissions.includes("customers.manage")),
);

const primaryContact = computed(
  () => detail.value?.customer?.contacts?.find((item: any) => item.isPrimary) ?? null,
);

const textMap = {
  customerType: { PERSONAL: "个人", COMPANY: "企业", RESELLER: "代理商" },
  customerStatus: { ACTIVE: "正常", SUSPENDED: "暂停", OVERDUE: "逾期", ARCHIVED: "归档" },
  orderStatus: { PENDING: "待支付", PAID: "已付款", PROVISIONING: "开通中", ACTIVE: "已生效", CANCELLED: "已取消", REFUNDED: "已退款" },
  serviceStatus: { PENDING: "待开通", PROVISIONING: "开通中", ACTIVE: "运行中", SUSPENDED: "已暂停", OVERDUE: "已逾期", TERMINATED: "已终止", EXPIRED: "已到期", FAILED: "异常" },
  invoiceStatus: { DRAFT: "草稿", ISSUED: "待支付", PARTIAL: "部分支付", PAID: "已支付", OVERDUE: "已逾期", VOID: "已作废" },
  creditType: { RECHARGE: "充值", CONSUME: "消费", REFUND: "退款", ADJUSTMENT: "调整", AUTO_RENEW: "自动续费" },
  portalRole: { OWNER: "主账户", TECH: "技术联系人", BILLING: "财务联系人", READONLY: "只读成员" },
  followUpType: { NOTE: "备注", CALL: "电话", VISIT: "拜访", EMAIL: "邮件" },
  certificationStatus: { PENDING: "待审核", VERIFIED: "已通过", REJECTED: "已驳回" },
  certificationType: { PERSONAL: "个人实名", COMPANY: "企业实名" },
  ticketStatus: { OPEN: "待处理", PROCESSING: "处理中", WAITING_CUSTOMER: "待客户回复", RESOLVED: "已解决", CLOSED: "已关闭" },
} as const;

function labelOf(map: Record<string, string>, value?: string | null) {
  if (!value) return "-";
  return map[value] ?? value;
}

function tagTypeForStatus(status?: string | null) {
  switch (status) {
    case "ACTIVE":
    case "SUCCESS":
    case "PAID":
    case "VERIFIED":
    case "RESOLVED":
      return "success";
    case "PENDING":
    case "PROVISIONING":
    case "PARTIAL":
    case "PROCESSING":
      return "warning";
    case "SUSPENDED":
    case "OVERDUE":
    case "EXPIRED":
    case "FAILED":
    case "VOID":
    case "CANCELLED":
    case "REFUNDED":
    case "CLOSED":
    case "TERMINATED":
    case "REJECTED":
      return "danger";
    default:
      return "info";
  }
}

function fillCustomerForm() {
  const customer = detail.value?.customer;
  if (!customer) return;
  Object.assign(customerForm, {
    name: customer.name,
    companyName: customer.companyName || "",
    email: customer.email,
    phone: customer.phone || "",
    contactQQ: customer.contactQQ || "",
    contactWechat: customer.contactWechat || "",
    type: customer.type,
    status: customer.status,
    level: customer.level,
    tags: customer.tags || "",
    notes: customer.notes || "",
  });
}

function fillCertificationForm() {
  const certification = detail.value?.customer?.certification;
  Object.assign(certificationForm, {
    subjectType: certification?.subjectType || detail.value?.customer?.type || "PERSONAL",
    subjectName: certification?.subjectName || detail.value?.customer?.companyName || detail.value?.customer?.name || "",
    idNumber: certification?.idNumber || "",
    businessLicenseNo: certification?.businessLicenseNo || "",
    status: certification?.status || "PENDING",
    submittedAt: certification?.submittedAt ? String(certification.submittedAt).slice(0, 10) : "",
    reviewNote: certification?.reviewNote || "",
  });
}

function resetContactForm() {
  Object.assign(contactForm, { id: "", name: "", department: "", role: "", email: "", phone: "", isPrimary: false, notes: "" });
}

function resetFollowUpForm() {
  Object.assign(followUpForm, { id: "", type: "NOTE", title: "", content: "", nextFollowAt: "" });
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get(`/customers/${route.params.id}`);
    detail.value = data.data;
    fillCustomerForm();
    fillCertificationForm();
  } finally {
    loading.value = false;
  }
}

async function submitCustomer() {
  await http.put(`/customers/${route.params.id}`, customerForm);
  ElMessage.success("客户信息已更新");
  customerDialogVisible.value = false;
  await loadData();
}

async function submitBalance() {
  await http.post(`/customers/${route.params.id}/balance`, balanceForm);
  ElMessage.success("客户余额已调整");
  balanceDialogVisible.value = false;
  await loadData();
}

function openContactDialog(contact?: any) {
  if (!contact) {
    resetContactForm();
  } else {
    Object.assign(contactForm, {
      id: contact.id,
      name: contact.name,
      department: contact.department || "",
      role: contact.role || "",
      email: contact.email || "",
      phone: contact.phone || "",
      isPrimary: contact.isPrimary,
      notes: contact.notes || "",
    });
  }
  contactDialogVisible.value = true;
}

async function submitContact() {
  if (contactForm.id) {
    await http.put(`/customers/${route.params.id}/contacts/${contactForm.id}`, contactForm);
    ElMessage.success("联系人已更新");
  } else {
    await http.post(`/customers/${route.params.id}/contacts`, contactForm);
    ElMessage.success("联系人已新增");
  }
  contactDialogVisible.value = false;
  await loadData();
}

async function removeContact(contact: any) {
  await ElMessageBox.confirm(`确认删除联系人「${contact.name}」吗？`, "删除联系人", { type: "warning" });
  await http.delete(`/customers/${route.params.id}/contacts/${contact.id}`);
  ElMessage.success("联系人已删除");
  await loadData();
}

async function submitCertification() {
  certificationSubmitting.value = true;
  try {
    await http.put(`/customers/${route.params.id}/certification`, certificationForm);
    ElMessage.success("实名认证资料已保存");
    await loadData();
  } finally {
    certificationSubmitting.value = false;
  }
}

function openFollowUpDialog(followUp?: any) {
  if (!followUp) {
    resetFollowUpForm();
  } else {
    Object.assign(followUpForm, {
      id: followUp.id,
      type: followUp.type,
      title: followUp.title,
      content: followUp.content,
      nextFollowAt: followUp.nextFollowAt ? String(followUp.nextFollowAt).slice(0, 10) : "",
    });
  }
  followUpDialogVisible.value = true;
}

async function submitFollowUp() {
  if (followUpForm.id) {
    await http.put(`/customers/${route.params.id}/follow-ups/${followUpForm.id}`, followUpForm);
    ElMessage.success("跟进记录已更新");
  } else {
    await http.post(`/customers/${route.params.id}/follow-ups`, followUpForm);
    ElMessage.success("跟进记录已新增");
  }
  followUpDialogVisible.value = false;
  await loadData();
}

async function removeFollowUp(followUp: any) {
  await ElMessageBox.confirm(`确认删除跟进记录「${followUp.title}」吗？`, "删除跟进记录", { type: "warning" });
  await http.delete(`/customers/${route.params.id}/follow-ups/${followUp.id}`);
  ElMessage.success("跟进记录已删除");
  await loadData();
}

onMounted(async () => {
  await authStore.ensureUser();
  await loadData();
});

watch(() => route.params.id, () => {
  void loadData();
});
</script>

<template>
  <div class="page-body" v-loading="loading">
    <div class="toolbar-row">
      <div class="toolbar-actions">
        <el-button :icon="ArrowLeft" @click="router.push('/customers')">返回客户列表</el-button>
        <el-button v-if="canManageCustomers" type="primary" :icon="Edit" @click="customerDialogVisible = true; fillCustomerForm()">编辑客户</el-button>
        <el-button v-if="canManageCustomers" type="warning" :icon="Money" @click="balanceDialogVisible = true">调整余额</el-button>
        <el-button v-if="canManageCustomers" :icon="Plus" @click="openContactDialog()">新增联系人</el-button>
        <el-button v-if="canManageCustomers" :icon="Plus" @click="openFollowUpDialog()">新增跟进</el-button>
      </div>
    </div>

    <div v-if="detail" class="page-body">
      <div class="detail-grid">
        <div class="detail-hero">
          <div class="toolbar-row" style="margin-bottom: 12px">
            <div>
              <div class="detail-hero-title">{{ detail.customer.name }}</div>
              <div class="detail-hero-copy">
                {{ detail.customer.customerNo }} / {{ detail.customer.companyName || "个人客户" }} / {{ labelOf(textMap.customerType, detail.customer.type) }}
              </div>
            </div>
            <el-tag :type="tagTypeForStatus(detail.customer.status)" size="large">{{ labelOf(textMap.customerStatus, detail.customer.status) }}</el-tag>
          </div>

          <div class="summary-strip">
            <div class="summary-card"><div class="metric-label">累计开票</div><h3>{{ formatCurrency(detail.summary.totalBilled) }}</h3><p>历史账单总金额</p></div>
            <div class="summary-card"><div class="metric-label">累计实收</div><h3>{{ formatCurrency(detail.summary.totalPaid) }}</h3><p>已入账金额</p></div>
            <div class="summary-card"><div class="metric-label">账户余额</div><h3>{{ formatCurrency(detail.customer.creditBalance) }}</h3><p>可用于续费与抵扣</p></div>
            <div class="summary-card"><div class="metric-label">待收金额</div><h3>{{ formatCurrency(detail.summary.totalOutstanding) }}</h3><p>全部未结清账单合计</p></div>
          </div>
        </div>

        <div class="detail-side">
          <el-card class="page-card">
            <div class="section-heading"><h2>客户摘要</h2><el-tag effect="plain">{{ detail.customer.level }}</el-tag></div>
            <div class="detail-meta-list">
              <div class="detail-meta-item"><label>邮箱</label><strong>{{ detail.customer.email }}</strong></div>
              <div class="detail-meta-item"><label>电话</label><strong>{{ detail.customer.phone || "-" }}</strong></div>
              <div class="detail-meta-item"><label>主联系人</label><strong>{{ primaryContact?.name || "-" }}</strong></div>
              <div class="detail-meta-item"><label>实名认证</label><strong>{{ labelOf(textMap.certificationStatus, detail.customer.certification?.status || "PENDING") }}</strong></div>
              <div class="detail-meta-item"><label>下次跟进</label><strong>{{ formatDate(detail.summary.nextFollowUpAt) }}</strong></div>
            </div>
          </el-card>
        </div>
      </div>

      <el-card class="page-card">
        <el-tabs class="content-tabs">
          <el-tab-pane label="订单与服务">
            <el-row :gutter="16">
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>最近订单</h2></div>
                <el-table :data="detail.customer.orders" stripe>
                  <el-table-column label="订单号" min-width="180">
                    <template #default="{ row }">
                      <el-button text type="primary" @click="router.push(`/orders/${row.id}`)">{{ row.orderNo }}</el-button>
                      <div class="muted-line">{{ formatDateTime(row.createdAt) }}</div>
                    </template>
                  </el-table-column>
                  <el-table-column label="金额" width="120"><template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template></el-table-column>
                  <el-table-column label="状态" width="120">
                    <template #default="{ row }"><el-tag :type="tagTypeForStatus(row.status)">{{ labelOf(textMap.orderStatus, row.status) }}</el-tag></template>
                  </el-table-column>
                </el-table>
              </el-col>
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>最近服务</h2></div>
                <el-table :data="detail.customer.services" stripe>
                  <el-table-column label="服务" min-width="220">
                    <template #default="{ row }">
                      <el-button text type="primary" @click="router.push(`/services/${row.id}`)">{{ row.serviceNo }}</el-button>
                      <div class="muted-line">{{ row.name }} / {{ row.product.name }}</div>
                    </template>
                  </el-table-column>
                  <el-table-column label="资源" width="140"><template #default="{ row }">{{ row._count.ipAddresses }} IP / {{ row._count.disks }} 磁盘</template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="tagTypeForStatus(row.status)">{{ labelOf(textMap.serviceStatus, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane label="财务">
            <el-row :gutter="16">
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>最近账单</h2></div>
                <el-table :data="detail.customer.invoices" stripe>
                  <el-table-column label="账单号" min-width="180">
                    <template #default="{ row }">
                      <el-button text type="primary" @click="router.push(`/invoices/${row.id}`)">{{ row.invoiceNo }}</el-button>
                      <div class="muted-line">到期 {{ formatDate(row.dueDate) }}</div>
                    </template>
                  </el-table-column>
                  <el-table-column label="总额" width="120"><template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="tagTypeForStatus(row.status)">{{ labelOf(textMap.invoiceStatus, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>余额流水</h2></div>
                <el-table :data="detail.customer.creditTransactions" stripe>
                  <el-table-column label="时间" width="170"><template #default="{ row }">{{ formatDateTime(row.createdAt) }}</template></el-table-column>
                  <el-table-column label="类型" width="120"><template #default="{ row }">{{ labelOf(textMap.creditType, row.type) }}</template></el-table-column>
                  <el-table-column label="金额" width="120"><template #default="{ row }">{{ formatCurrency(row.amount) }}</template></el-table-column>
                  <el-table-column label="说明" min-width="180" prop="description" />
                </el-table>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane label="联系人与实名">
            <el-row :gutter="16">
              <el-col :xl="13" :span="24">
                <div class="section-heading"><h2>联系人</h2></div>
                <el-table :data="detail.customer.contacts" stripe>
                  <el-table-column label="联系人" min-width="170">
                    <template #default="{ row }">
                      <div class="primary-line">{{ row.name }}</div>
                      <div class="muted-line">{{ row.department || "-" }} / {{ row.role || "-" }}</div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="email" label="邮箱" min-width="180" />
                  <el-table-column prop="phone" label="电话" width="140" />
                  <el-table-column label="主联系人" width="100"><template #default="{ row }"><el-tag v-if="row.isPrimary" type="success">是</el-tag><span v-else>-</span></template></el-table-column>
                  <el-table-column v-if="canManageCustomers" label="操作" width="150">
                    <template #default="{ row }">
                      <el-button text type="primary" @click="openContactDialog(row)">编辑</el-button>
                      <el-button text type="danger" @click="removeContact(row)">删除</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </el-col>
              <el-col :xl="11" :span="24">
                <div class="section-heading"><h2>实名认证</h2><el-tag :type="tagTypeForStatus(certificationForm.status)">{{ labelOf(textMap.certificationStatus, certificationForm.status) }}</el-tag></div>
                <el-form label-position="top">
                  <el-row :gutter="12">
                    <el-col :span="12"><el-form-item label="认证类型"><el-select v-model="certificationForm.subjectType" style="width: 100%"><el-option label="个人实名" value="PERSONAL" /><el-option label="企业实名" value="COMPANY" /></el-select></el-form-item></el-col>
                    <el-col :span="12"><el-form-item label="审核状态"><el-select v-model="certificationForm.status" style="width: 100%"><el-option label="待审核" value="PENDING" /><el-option label="已通过" value="VERIFIED" /><el-option label="已驳回" value="REJECTED" /></el-select></el-form-item></el-col>
                  </el-row>
                  <el-form-item label="认证主体"><el-input v-model="certificationForm.subjectName" /></el-form-item>
                  <el-form-item label="证件号 / 执照号"><el-input v-model="certificationForm.businessLicenseNo" :placeholder="labelOf(textMap.certificationType, certificationForm.subjectType)" /></el-form-item>
                  <el-form-item label="审核备注"><el-input v-model="certificationForm.reviewNote" type="textarea" :rows="4" /></el-form-item>
                  <el-button v-if="canManageCustomers" type="primary" :loading="certificationSubmitting" @click="submitCertification">保存实名认证资料</el-button>
                </el-form>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane label="跟进与工单">
            <el-row :gutter="16">
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>跟进记录</h2></div>
                <el-table :data="detail.customer.followUps" stripe>
                  <el-table-column label="标题" min-width="200">
                    <template #default="{ row }">
                      <div class="primary-line">{{ row.title }}</div>
                      <div class="muted-line">{{ row.content }}</div>
                    </template>
                  </el-table-column>
                  <el-table-column label="类型" width="100"><template #default="{ row }">{{ labelOf(textMap.followUpType, row.type) }}</template></el-table-column>
                  <el-table-column label="跟进人" width="120" prop="operatorName" />
                  <el-table-column label="下次跟进" width="120"><template #default="{ row }">{{ formatDate(row.nextFollowAt) }}</template></el-table-column>
                  <el-table-column v-if="canManageCustomers" label="操作" width="150">
                    <template #default="{ row }">
                      <el-button text type="primary" @click="openFollowUpDialog(row)">编辑</el-button>
                      <el-button text type="danger" @click="removeFollowUp(row)">删除</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </el-col>
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>最近工单</h2></div>
                <el-table :data="detail.customer.tickets" stripe>
                  <el-table-column prop="ticketNo" label="工单号" min-width="160" />
                  <el-table-column prop="subject" label="主题" min-width="200" />
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="tagTypeForStatus(row.status)">{{ labelOf(textMap.ticketStatus, row.status) }}</el-tag></template></el-table-column>
                  <el-table-column label="更新时间" width="170"><template #default="{ row }">{{ formatDateTime(row.updatedAt) }}</template></el-table-column>
                </el-table>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane label="审计日志">
            <el-table :data="detail.auditLogs" stripe>
              <el-table-column label="时间" width="180"><template #default="{ row }">{{ formatDateTime(row.createdAt) }}</template></el-table-column>
              <el-table-column label="动作" width="160" prop="action" />
              <el-table-column label="摘要" min-width="220" prop="summary" />
              <el-table-column label="详情" min-width="280" prop="detail" show-overflow-tooltip />
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>

    <el-dialog v-model="customerDialogVisible" title="编辑客户" width="720px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12"><el-form-item label="客户名称"><el-input v-model="customerForm.name" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="企业名称"><el-input v-model="customerForm.companyName" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12"><el-form-item label="联系邮箱"><el-input v-model="customerForm.email" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="联系电话"><el-input v-model="customerForm.phone" /></el-form-item></el-col>
        </el-row>
        <el-form-item label="标签"><el-input v-model="customerForm.tags" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="customerForm.notes" type="textarea" :rows="4" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="customerDialogVisible = false">取消</el-button><el-button type="primary" @click="submitCustomer">保存</el-button></template>
    </el-dialog>

    <el-dialog v-model="balanceDialogVisible" title="调整客户余额" width="520px">
      <div v-if="detail" class="muted-line" style="margin-bottom: 12px">{{ detail.customer.name }} / 当前余额 {{ formatCurrency(detail.customer.creditBalance) }}</div>
      <el-form label-position="top">
        <el-form-item label="调整金额"><el-input-number v-model="balanceForm.amount" :precision="2" style="width: 100%" /></el-form-item>
        <el-form-item label="调整说明"><el-input v-model="balanceForm.description" type="textarea" :rows="4" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="balanceDialogVisible = false">取消</el-button><el-button type="primary" @click="submitBalance">提交</el-button></template>
    </el-dialog>

    <el-dialog v-model="contactDialogVisible" :title="contactForm.id ? '编辑联系人' : '新增联系人'" width="560px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12"><el-form-item label="姓名"><el-input v-model="contactForm.name" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="部门"><el-input v-model="contactForm.department" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12"><el-form-item label="岗位"><el-input v-model="contactForm.role" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="电话"><el-input v-model="contactForm.phone" /></el-form-item></el-col>
        </el-row>
        <el-form-item label="邮箱"><el-input v-model="contactForm.email" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="contactForm.notes" type="textarea" :rows="3" /></el-form-item>
        <el-form-item><el-switch v-model="contactForm.isPrimary" active-text="设为主联系人" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="contactDialogVisible = false">取消</el-button><el-button type="primary" @click="submitContact">保存</el-button></template>
    </el-dialog>

    <el-dialog v-model="followUpDialogVisible" :title="followUpForm.id ? '编辑跟进记录' : '新增跟进记录'" width="620px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12"><el-form-item label="类型"><el-select v-model="followUpForm.type" style="width: 100%"><el-option label="备注" value="NOTE" /><el-option label="电话" value="CALL" /><el-option label="拜访" value="VISIT" /><el-option label="邮件" value="EMAIL" /></el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="下次跟进"><el-input v-model="followUpForm.nextFollowAt" placeholder="YYYY-MM-DD" /></el-form-item></el-col>
        </el-row>
        <el-form-item label="标题"><el-input v-model="followUpForm.title" /></el-form-item>
        <el-form-item label="内容"><el-input v-model="followUpForm.content" type="textarea" :rows="5" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="followUpDialogVisible = false">取消</el-button><el-button type="primary" @click="submitFollowUp">保存</el-button></template>
    </el-dialog>
  </div>
</template>
