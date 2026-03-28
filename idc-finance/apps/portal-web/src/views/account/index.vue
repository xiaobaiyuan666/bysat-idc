<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { useRoute, useRouter } from "vue-router";
import {
  changePortalPassword,
  createContact,
  deleteContact,
  loadAccount,
  loadZhimaKYCRecords,
  submitIdentity,
  submitZhimaKYC,
  updateContact,
  updateProfile,
  type PortalAccount,
  type PortalContact,
  type PortalKYCRecord
} from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalCustomerStatus,
  formatPortalIdentityStatus,
  formatPortalMoney
} from "@/utils/business";

const localeStore = useLocaleStore();
const route = useRoute();
const router = useRouter();
const loading = ref(false);
const saving = ref(false);
const error = ref("");
const activeTab = ref("profile");
const account = ref<PortalAccount | null>(null);
const kycLoading = ref(false);
const kycSubmitting = ref(false);
const kycRecords = ref<PortalKYCRecord[]>([]);

const profileForm = reactive({
  name: "",
  email: "",
  mobile: "",
  remarks: ""
});

const contactForm = reactive({
  id: 0,
  name: "",
  email: "",
  mobile: "",
  roleName: "",
  isPrimary: false
});
const contactDialogVisible = ref(false);

const identityForm = reactive({
  identityType: "COMPANY",
  subjectName: "",
  certNo: "",
  countryCode: "CN"
});

const pluginKycForm = reactive({
  realName: "",
  idCardNo: "",
  mobile: ""
});

const securityForm = reactive({
  currentPassword: "",
  newPassword: "",
  confirmPassword: ""
});

const primaryContactCount = computed(() => account.value?.customer.contacts.filter(item => item.isPrimary).length ?? 0);
const accountActions = computed(() => [
  {
    title: copy.value.profileTab,
    description: pickLabel(localeStore.locale, "维护基础资料和联系信息。", "Maintain base profile and contact info."),
    action: () => (activeTab.value = "profile")
  },
  {
    title: copy.value.contactsTab,
    description: pickLabel(localeStore.locale, "新增或编辑联系人。", "Add or edit contacts."),
    action: () => (activeTab.value = "contacts")
  },
  {
    title: copy.value.identityTab,
    description: pickLabel(localeStore.locale, "提交或更新实名认证资料。", "Submit or update identity documents."),
    action: () => (activeTab.value = "identity")
  },
  {
    title: copy.value.securityTab,
    description: pickLabel(localeStore.locale, "更新登录密码和安全设置。", "Update password and security settings."),
    action: () => (activeTab.value = "security")
  }
]);
const accountChecklist = computed(() => [
  {
    label: copy.value.identity,
    value: account.value?.customer.identity?.verifyStatus
      ? formatPortalIdentityStatus(localeStore.locale, account.value.customer.identity?.verifyStatus)
      : pickLabel(localeStore.locale, "未提交", "Not Submitted")
  },
  {
    label: copy.value.contacts,
    value: String(account.value?.customer.contacts.length ?? 0)
  },
  {
    label: copy.value.primaryContact,
    value: String(primaryContactCount.value)
  },
  {
    label: copy.value.wallet,
    value: account.value ? formatPortalMoney(localeStore.locale, account.value.wallet.balance) : formatPortalMoney(localeStore.locale, 0)
  }
]);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "账户中心", "Account"),
  title: pickLabel(localeStore.locale, "账户中心", "Account"),
  subtitle: pickLabel(
    localeStore.locale,
    "统一维护基础资料、联系人、实名认证与安全设置，并与后台客户档案保持同步。",
    "Manage profile, contacts, identity, and security settings in one place."
  ),
  balance: pickLabel(localeStore.locale, "账户余额", "Balance"),
  status: pickLabel(localeStore.locale, "账户状态", "Status"),
  identity: pickLabel(localeStore.locale, "实名认证", "Identity"),
  contacts: pickLabel(localeStore.locale, "联系人", "Contacts"),
  primaryContact: pickLabel(localeStore.locale, "主联系人", "Primary Contact"),
  actionDesk: pickLabel(localeStore.locale, "账户工作台", "Account Desk"),
  actionDeskDesc: pickLabel(localeStore.locale, "把资料维护、联系人、实名和安全设置集中管理。", "Manage profile, contacts, identity, and security from one desk."),
  checklist: pickLabel(localeStore.locale, "账户检查项", "Account Checklist"),
  checklistDesc: pickLabel(localeStore.locale, "快速查看实名、联系人和钱包状态。", "Quickly review identity, contact, and wallet status."),
  profileTab: pickLabel(localeStore.locale, "基础资料", "Profile"),
  contactsTab: pickLabel(localeStore.locale, "联系人管理", "Contacts"),
  identityTab: pickLabel(localeStore.locale, "实名认证", "Identity"),
  securityTab: pickLabel(localeStore.locale, "安全设置", "Security"),
  save: pickLabel(localeStore.locale, "保存", "Save"),
  addContact: pickLabel(localeStore.locale, "新增联系人", "Add Contact"),
  edit: pickLabel(localeStore.locale, "编辑", "Edit"),
  remove: pickLabel(localeStore.locale, "删除", "Delete"),
  wallet: pickLabel(localeStore.locale, "钱包中心", "Wallet"),
  invoices: pickLabel(localeStore.locale, "财务中心", "Finance"),
  profileSaved: pickLabel(localeStore.locale, "资料已更新", "Profile updated"),
  identitySaved: pickLabel(localeStore.locale, "实名认证资料已提交", "Identity submitted"),
  passwordSaved: pickLabel(localeStore.locale, "密码已更新", "Password updated"),
  contactSaved: pickLabel(localeStore.locale, "联系人已保存", "Contact saved"),
  contactRemoved: pickLabel(localeStore.locale, "联系人已删除", "Contact removed"),
  loadError: pickLabel(localeStore.locale, "账户资料加载失败", "Failed to load account"),
  currentPassword: pickLabel(localeStore.locale, "当前密码", "Current password"),
  newPassword: pickLabel(localeStore.locale, "新密码", "New password"),
  confirmPassword: pickLabel(localeStore.locale, "确认新密码", "Confirm password"),
  name: pickLabel(localeStore.locale, "名称", "Name"),
  email: pickLabel(localeStore.locale, "邮箱", "Email"),
  mobile: pickLabel(localeStore.locale, "手机", "Mobile"),
  remarks: pickLabel(localeStore.locale, "备注", "Remarks"),
  roleName: pickLabel(localeStore.locale, "角色", "Role"),
  isPrimary: pickLabel(localeStore.locale, "设为主联系人", "Set as primary"),
  identityType: pickLabel(localeStore.locale, "实名类型", "Identity Type"),
  subjectName: pickLabel(localeStore.locale, "实名主体", "Subject"),
  certNo: pickLabel(localeStore.locale, "证件号码", "Certificate No."),
  countryCode: pickLabel(localeStore.locale, "国家 / 地区", "Country"),
  deleteConfirm: pickLabel(localeStore.locale, "确定删除这位联系人吗？", "Delete this contact?"),
  company: pickLabel(localeStore.locale, "企业", "Company"),
  personal: pickLabel(localeStore.locale, "个人", "Personal"),
  yes: pickLabel(localeStore.locale, "是", "Yes"),
  no: pickLabel(localeStore.locale, "否", "No"),
  cancel: pickLabel(localeStore.locale, "取消", "Cancel"),
  passwordMismatch: pickLabel(localeStore.locale, "两次输入的新密码不一致", "Passwords do not match")
}));

function fillProfile() {
  if (!account.value) return;
  profileForm.name = account.value.customer.name;
  profileForm.email = account.value.customer.email;
  profileForm.mobile = account.value.customer.mobile;
  profileForm.remarks = account.value.customer.remarks || "";
  identityForm.identityType = account.value.customer.identity?.identityType || "COMPANY";
  identityForm.subjectName = account.value.customer.identity?.subjectName || account.value.customer.name || "";
  identityForm.certNo = account.value.customer.identity?.certNo || "";
  identityForm.countryCode = account.value.customer.identity?.countryCode || "CN";
}

async function fetchAccount() {
  loading.value = true;
  error.value = "";
  try {
    account.value = await loadAccount();
    fillProfile();
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

async function saveProfile() {
  saving.value = true;
  try {
    account.value = await updateProfile(profileForm);
    fillProfile();
    ElMessage.success(copy.value.profileSaved);
  } finally {
    saving.value = false;
  }
}

function openContactDialog(contact?: PortalContact) {
  contactForm.id = contact?.id || 0;
  contactForm.name = contact?.name || "";
  contactForm.email = contact?.email || "";
  contactForm.mobile = contact?.mobile || "";
  contactForm.roleName = contact?.roleName || "";
  contactForm.isPrimary = contact?.isPrimary || false;
  contactDialogVisible.value = true;
}

async function saveContact() {
  saving.value = true;
  try {
    if (contactForm.id) {
      await updateContact(contactForm.id, contactForm);
    } else {
      await createContact(contactForm);
    }
    contactDialogVisible.value = false;
    await fetchAccount();
    ElMessage.success(copy.value.contactSaved);
  } finally {
    saving.value = false;
  }
}

async function removeContact(contact: PortalContact) {
  try {
    await ElMessageBox.confirm(copy.value.deleteConfirm, copy.value.remove, { type: "warning" });
  } catch {
    return;
  }
  await deleteContact(contact.id);
  await fetchAccount();
  ElMessage.success(copy.value.contactRemoved);
}

async function saveIdentity() {
  saving.value = true;
  try {
    await submitIdentity(identityForm);
    await fetchAccount();
    ElMessage.success(copy.value.identitySaved);
  } finally {
    saving.value = false;
  }
}

async function fetchPluginKYCRecords() {
  kycLoading.value = true;
  try {
    kycRecords.value = await loadZhimaKYCRecords(20);
  } catch {
    kycRecords.value = [];
  } finally {
    kycLoading.value = false;
  }
}

function pluginKycStatusLabel(status: string) {
  switch (status) {
    case "SUCCESS":
    case "APPROVED":
      return "已通过";
    case "FAILED":
    case "REJECTED":
      return "失败";
    case "PENDING":
      return "待处理";
    default:
      return status || "-";
  }
}

function pluginKycStatusType(status: string): "success" | "warning" | "danger" | "info" {
  switch (status) {
    case "SUCCESS":
    case "APPROVED":
      return "success";
    case "FAILED":
    case "REJECTED":
      return "danger";
    case "PENDING":
      return "warning";
    default:
      return "info";
  }
}

async function submitPluginKyc() {
  if (!pluginKycForm.realName.trim() || !pluginKycForm.idCardNo.trim()) {
    ElMessage.warning("请填写真实姓名和证件号码");
    return;
  }

  kycSubmitting.value = true;
  try {
    await submitZhimaKYC({
      realName: pluginKycForm.realName.trim(),
      idCardNo: pluginKycForm.idCardNo.trim(),
      mobile: pluginKycForm.mobile.trim()
    });
    ElMessage.success("插件实名核验已提交");
    pluginKycForm.realName = "";
    pluginKycForm.idCardNo = "";
    pluginKycForm.mobile = "";
    await fetchPluginKYCRecords();
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : "插件实名核验提交失败");
  } finally {
    kycSubmitting.value = false;
  }
}

async function savePassword() {
  if (securityForm.newPassword !== securityForm.confirmPassword) {
    ElMessage.error(copy.value.passwordMismatch);
    return;
  }
  saving.value = true;
  try {
    await changePortalPassword({
      currentPassword: securityForm.currentPassword,
      newPassword: securityForm.newPassword
    });
    securityForm.currentPassword = "";
    securityForm.newPassword = "";
    securityForm.confirmPassword = "";
    ElMessage.success(copy.value.passwordSaved);
  } finally {
    saving.value = false;
  }
}

function syncTabFromRoute() {
  const mapping: Record<string, string> = {
    "/account": "profile",
    "/account/profile": "profile",
    "/account/contacts": "contacts",
    "/account/identity": "identity",
    "/account/security": "security"
  };
  activeTab.value = mapping[route.path] ?? "profile";
}

onMounted(() => {
  syncTabFromRoute();
  void fetchAccount();
  void fetchPluginKYCRecords();
});

watch(
  () => route.path,
  () => {
    syncTabFromRoute();
  }
);

watch(activeTab, tab => {
  const mapping: Record<string, string> = {
    profile: "/account/profile",
    contacts: "/account/contacts",
    identity: "/account/identity",
    security: "/account/security"
  };
  const target = mapping[tab] || "/account/profile";
  if (route.path !== target) {
    void router.replace(target);
  }
});
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />

    <section v-if="account" class="portal-card">
      <div class="portal-card-head">
        <div>
          <div class="portal-badge">{{ copy.badge }}</div>
          <h1 class="portal-title" style="margin-top: 12px">{{ copy.title }}</h1>
          <p class="portal-subtitle">{{ copy.subtitle }}</p>
        </div>
      </div>

      <div class="portal-toolbar" style="margin-top: 20px">
        <el-button type="primary" plain @click="router.push('/wallet')">{{ copy.wallet }}</el-button>
        <el-button plain @click="router.push('/invoices')">{{ copy.invoices }}</el-button>
      </div>

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <article class="portal-stat">
          <h3>{{ copy.name }}</h3>
          <strong>{{ account.customer.name }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.status }}</h3>
          <strong>{{ formatPortalCustomerStatus(localeStore.locale, account.customer.status) }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.balance }}</h3>
          <strong>{{ formatPortalMoney(localeStore.locale, account.wallet.balance) }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.identity }}</h3>
          <strong>{{ formatPortalIdentityStatus(localeStore.locale, account.customer.identity?.verifyStatus) }}</strong>
        </article>
      </div>
    </section>

    <section v-if="account" class="portal-card">
      <div class="portal-grid portal-grid--four" style="margin-bottom: 20px">
        <article class="portal-stat">
          <h3>{{ copy.contacts }}</h3>
          <strong>{{ account.customer.contacts.length }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.primaryContact }}</h3>
          <strong>{{ primaryContactCount }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ pickLabel(localeStore.locale, '客户分组', 'Group') }}</h3>
          <strong>{{ account.customer.groupName || "-" }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ pickLabel(localeStore.locale, '客户等级', 'Level') }}</h3>
          <strong>{{ account.customer.levelName || "-" }}</strong>
        </article>
      </div>

      <section class="portal-grid portal-grid--two" style="margin-bottom: 20px">
        <article class="portal-card" style="padding: 0">
          <div style="padding: 20px 20px 0">
            <div class="portal-card-head">
              <div>
                <h2 class="portal-panel__title">{{ copy.actionDesk }}</h2>
                <div class="portal-panel__meta">{{ copy.actionDeskDesc }}</div>
              </div>
            </div>
          </div>
          <div class="portal-actions-grid" style="padding: 18px 20px 20px">
            <button
              v-for="item in accountActions"
              :key="item.title"
              type="button"
              class="portal-action-card"
              @click="item.action()"
            >
              <strong>{{ item.title }}</strong>
              <span>{{ item.description }}</span>
            </button>
          </div>
        </article>

        <article class="portal-card" style="padding: 0">
          <div style="padding: 20px 20px 0">
            <div class="portal-card-head">
              <div>
                <h2 class="portal-panel__title">{{ copy.checklist }}</h2>
                <div class="portal-panel__meta">{{ copy.checklistDesc }}</div>
              </div>
            </div>
          </div>
          <div class="portal-summary" style="margin: 18px 20px 20px">
            <div v-for="item in accountChecklist" :key="item.label" class="portal-summary-row">
              <span>{{ item.label }}</span>
              <strong>{{ item.value }}</strong>
            </div>
          </div>
        </article>
      </section>

      <el-tabs v-model="activeTab">
        <el-tab-pane :label="copy.profileTab" name="profile">
          <el-form label-position="top" class="portal-grid portal-grid--two" style="margin-top: 12px">
            <el-form-item :label="copy.name"><el-input v-model="profileForm.name" /></el-form-item>
            <el-form-item :label="copy.email"><el-input v-model="profileForm.email" /></el-form-item>
            <el-form-item :label="copy.mobile"><el-input v-model="profileForm.mobile" /></el-form-item>
            <el-form-item :label="copy.remarks"><el-input v-model="profileForm.remarks" type="textarea" :rows="3" /></el-form-item>
          </el-form>
          <div class="portal-toolbar" style="margin-top: 16px">
            <el-button type="primary" :loading="saving" @click="saveProfile">{{ copy.save }}</el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="copy.contactsTab" name="contacts">
          <div class="portal-card-head" style="margin-top: 12px">
            <div>
              <h2 class="portal-panel__title">{{ copy.contacts }}</h2>
            </div>
            <el-button type="primary" plain @click="openContactDialog()">{{ copy.addContact }}</el-button>
          </div>
          <el-table :data="account.customer.contacts" border style="margin-top: 18px">
            <el-table-column prop="name" :label="copy.name" min-width="140" />
            <el-table-column prop="roleName" :label="copy.roleName" min-width="140" />
            <el-table-column prop="email" :label="copy.email" min-width="220" />
            <el-table-column prop="mobile" :label="copy.mobile" min-width="160" />
            <el-table-column :label="copy.primaryContact" min-width="100">
              <template #default="{ row }">
                <el-tag :type="row.isPrimary ? 'success' : 'info'">{{ row.isPrimary ? copy.yes : copy.no }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column width="160" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openContactDialog(row)">{{ copy.edit }}</el-button>
                <el-button link type="danger" @click="removeContact(row)">{{ copy.remove }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane :label="copy.identityTab" name="identity">
          <el-form label-position="top" class="portal-grid portal-grid--two" style="margin-top: 12px">
            <el-form-item :label="copy.identityType">
              <el-select v-model="identityForm.identityType">
                <el-option :label="copy.company" value="COMPANY" />
                <el-option :label="copy.personal" value="PERSONAL" />
              </el-select>
            </el-form-item>
            <el-form-item :label="copy.countryCode"><el-input v-model="identityForm.countryCode" /></el-form-item>
            <el-form-item :label="copy.subjectName"><el-input v-model="identityForm.subjectName" /></el-form-item>
            <el-form-item :label="copy.certNo"><el-input v-model="identityForm.certNo" /></el-form-item>
          </el-form>
          <div class="portal-toolbar" style="margin-top: 16px">
            <el-button type="primary" :loading="saving" @click="saveIdentity">{{ copy.save }}</el-button>
          </div>

          <div class="portal-card" style="margin-top: 20px" v-loading="kycLoading">
            <div class="portal-card-head">
              <div>
                <h2 class="portal-panel__title">插件实名核验</h2>
                <div class="portal-panel__meta">补充调用芝麻实名认证插件做实名校验，并保留提交记录。</div>
              </div>
            </div>

            <el-form label-position="top" class="portal-grid portal-grid--three" style="margin-top: 16px">
              <el-form-item label="真实姓名">
                <el-input v-model="pluginKycForm.realName" />
              </el-form-item>
              <el-form-item label="证件号码">
                <el-input v-model="pluginKycForm.idCardNo" />
              </el-form-item>
              <el-form-item label="手机号（可选）">
                <el-input v-model="pluginKycForm.mobile" />
              </el-form-item>
            </el-form>

            <div class="portal-toolbar" style="margin-top: 8px">
              <el-button type="success" :loading="kycSubmitting" @click="submitPluginKyc">
                提交插件核验
              </el-button>
            </div>

            <el-table :data="kycRecords" border style="margin-top: 18px">
              <el-table-column prop="realName" label="真实姓名" min-width="140" />
              <el-table-column prop="idCardNo" label="证件号码" min-width="220" />
              <el-table-column prop="mobile" label="手机号" min-width="140" />
              <el-table-column label="状态" min-width="120">
                <template #default="{ row }">
                  <el-tag :type="pluginKycStatusType(row.status)">{{ pluginKycStatusLabel(row.status) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="message" label="返回信息" min-width="220" show-overflow-tooltip />
              <el-table-column prop="createdAt" label="提交时间" min-width="180" />
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="copy.securityTab" name="security">
          <el-form label-position="top" class="portal-grid portal-grid--two" style="margin-top: 12px">
            <el-form-item :label="copy.currentPassword">
              <el-input v-model="securityForm.currentPassword" type="password" show-password />
            </el-form-item>
            <el-form-item :label="copy.newPassword">
              <el-input v-model="securityForm.newPassword" type="password" show-password />
            </el-form-item>
            <el-form-item :label="copy.confirmPassword">
              <el-input v-model="securityForm.confirmPassword" type="password" show-password />
            </el-form-item>
          </el-form>
          <div class="portal-toolbar" style="margin-top: 16px">
            <el-button type="primary" :loading="saving" @click="savePassword">{{ copy.save }}</el-button>
          </div>
        </el-tab-pane>
      </el-tabs>
    </section>

    <el-dialog v-model="contactDialogVisible" :title="contactForm.id ? copy.edit : copy.addContact" width="520px">
      <el-form label-position="top">
        <el-form-item :label="copy.name"><el-input v-model="contactForm.name" /></el-form-item>
        <el-form-item :label="copy.roleName"><el-input v-model="contactForm.roleName" /></el-form-item>
        <el-form-item :label="copy.email"><el-input v-model="contactForm.email" /></el-form-item>
        <el-form-item :label="copy.mobile"><el-input v-model="contactForm.mobile" /></el-form-item>
        <el-form-item :label="copy.isPrimary">
          <el-switch v-model="contactForm.isPrimary" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="contactDialogVisible = false">{{ copy.cancel }}</el-button>
        <el-button type="primary" :loading="saving" @click="saveContact">{{ copy.save }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>
