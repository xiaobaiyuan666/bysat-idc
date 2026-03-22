<script setup lang="ts">
import { Plus, Refresh } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { onMounted, reactive, ref } from "vue";

import { http } from "@/api/http";
import { formatDate } from "@/utils/format";
import { getLabel, getStatusTagType, resourceStatusMap } from "@/utils/maps";

const loading = ref(false);
const activeTab = ref("vpcs");
const ruleDialogVisible = ref(false);
const currentSecurityGroupId = ref("");
const resources = ref<any>({
  vpcs: [],
  ips: [],
  disks: [],
  snapshots: [],
  backups: [],
  securityGroups: [],
});

const ruleForm = reactive({
  direction: "ingress",
  protocol: "tcp",
  portRange: "22",
  sourceCidr: "0.0.0.0/0",
  ruleAction: "ALLOW",
  description: "",
});

function resetRuleForm() {
  Object.assign(ruleForm, {
    direction: "ingress",
    protocol: "tcp",
    portRange: "22",
    sourceCidr: "0.0.0.0/0",
    ruleAction: "ALLOW",
    description: "",
  });
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/resources");
    resources.value = data.data;
  } finally {
    loading.value = false;
  }
}

async function runResourceAction(
  resourceType: string,
  id: string,
  action: string,
  payload?: Record<string, unknown>,
) {
  const { data } = await http.post(`/resources/${resourceType}/${id}/action`, {
    action,
    payload,
  });

  ElMessage.success(data.message || "操作已提交");
  await loadData();
}

async function createSnapshot(row: any) {
  try {
    const { value } = await ElMessageBox.prompt(
      "可选填写快照名称，留空则自动生成。",
      "创建快照",
      {
        inputValue: `${row.name}-快照`,
        confirmButtonText: "创建",
        cancelButtonText: "取消",
      },
    );

    await runResourceAction("disks", row.id, "createSnapshot", {
      name: value,
    });
  } catch {
    return;
  }
}

async function createBackup(row: any) {
  try {
    const { value } = await ElMessageBox.prompt(
      "可选填写备份名称，默认保留 30 天。",
      "创建备份",
      {
        inputValue: `${row.name}-备份`,
        confirmButtonText: "创建",
        cancelButtonText: "取消",
      },
    );

    await runResourceAction("disks", row.id, "createBackup", {
      name: value,
    });
  } catch {
    return;
  }
}

async function toggleDiskMount(row: any) {
  if (row.status === "DETACHED") {
    try {
      const { value } = await ElMessageBox.prompt(
        "请输入挂载点。",
        "挂载磁盘",
        {
          inputValue: row.mountPoint || "/data",
          confirmButtonText: "挂载",
          cancelButtonText: "取消",
        },
      );

      await runResourceAction("disks", row.id, "attach", {
        mountPoint: value,
      });
    } catch {
      return;
    }

    return;
  }

  try {
    await ElMessageBox.confirm(`确认卸载磁盘 ${row.name} 吗？`, "卸载磁盘", {
      type: "warning",
    });
  } catch {
    return;
  }

  await runResourceAction("disks", row.id, "detach");
}

async function setBootDisk(row: any) {
  try {
    await ElMessageBox.confirm(
      `确认将磁盘 ${row.name} 设为系统盘吗？当前实例的其他系统盘标记会被清除。`,
      "设为系统盘",
      {
        type: "warning",
      },
    );
  } catch {
    return;
  }

  await runResourceAction("disks", row.id, "setBoot");
}

async function restoreSnapshot(row: any) {
  try {
    await ElMessageBox.confirm(`确认使用快照 ${row.name} 执行回滚恢复吗？`, "快照恢复", {
      type: "warning",
    });
  } catch {
    return;
  }

  await runResourceAction("snapshots", row.id, "restore");
}

async function deleteSnapshot(row: any) {
  try {
    await ElMessageBox.confirm(`确认删除快照 ${row.name} 吗？`, "删除快照", {
      type: "warning",
    });
  } catch {
    return;
  }

  await runResourceAction("snapshots", row.id, "deleteSnapshot");
}

async function restoreBackup(row: any) {
  try {
    await ElMessageBox.confirm(`确认使用备份 ${row.name} 发起恢复吗？`, "备份恢复", {
      type: "warning",
    });
  } catch {
    return;
  }

  await runResourceAction("backups", row.id, "restore");
}

async function expireBackup(row: any) {
  try {
    await ElMessageBox.confirm(`确认将备份 ${row.name} 立即归档吗？`, "归档备份", {
      type: "warning",
    });
  } catch {
    return;
  }

  await runResourceAction("backups", row.id, "expireNow");
}

async function deleteBackup(row: any) {
  try {
    await ElMessageBox.confirm(`确认删除备份 ${row.name} 吗？`, "删除备份", {
      type: "warning",
    });
  } catch {
    return;
  }

  await runResourceAction("backups", row.id, "deleteBackup");
}

function openAddRule(row: any) {
  currentSecurityGroupId.value = row.id;
  resetRuleForm();
  ruleDialogVisible.value = true;
}

async function submitRule() {
  await runResourceAction("security-groups", currentSecurityGroupId.value, "addRule", {
    ...ruleForm,
  });

  ruleDialogVisible.value = false;
}

async function deleteRule(securityGroupId: string, ruleId: string) {
  try {
    await ElMessageBox.confirm("确认删除这条安全组规则吗？", "删除规则", {
      type: "warning",
    });
  } catch {
    return;
  }

  await runResourceAction("security-groups", securityGroupId, "deleteRule", {
    ruleId,
  });
}

async function deleteGroup(row: any) {
  try {
    await ElMessageBox.confirm(`确认删除安全组 ${row.name} 吗？`, "删除安全组", {
      type: "warning",
    });
  } catch {
    return;
  }

  await runResourceAction("security-groups", row.id, "deleteGroup");
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">资源中心</h1>
        <div class="page-subtitle">
          统一维护 VPC、IP、云硬盘、快照、备份和安全组，并直接执行挂载、回滚、归档和规则变更操作。
        </div>
      </div>
      <el-button :icon="Refresh" @click="loadData">刷新数据</el-button>
    </div>

    <el-card class="page-card">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="VPC 网络" name="vpcs">
          <el-table v-loading="loading" :data="resources.vpcs" stripe>
            <el-table-column prop="name" label="名称" min-width="180" />
            <el-table-column prop="region" label="地域" width="120" />
            <el-table-column prop="cidr" label="CIDR" width="160" />
            <el-table-column prop="gateway" label="网关" width="140" />
            <el-table-column label="服务 / 安全组" width="150">
              <template #default="{ row }">
                {{ row.services.length }} / {{ row.securityGroups.length }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">
                  {{ getLabel(resourceStatusMap, row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="IP 地址" name="ips">
          <el-table v-loading="loading" :data="resources.ips" stripe>
            <el-table-column prop="address" label="地址" min-width="180" />
            <el-table-column prop="version" label="版本" width="100" />
            <el-table-column prop="type" label="类型" width="120" />
            <el-table-column label="所属服务" min-width="180">
              <template #default="{ row }">
                {{ row.service?.serviceNo }} / {{ row.service?.customer?.name }}
              </template>
            </el-table-column>
            <el-table-column prop="bandwidthMbps" label="带宽(Mbps)" width="120" />
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">
                  {{ getLabel(resourceStatusMap, row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="云硬盘" name="disks">
          <el-table v-loading="loading" :data="resources.disks" stripe>
            <el-table-column prop="name" label="名称" min-width="180" />
            <el-table-column label="所属服务" min-width="180">
              <template #default="{ row }">
                {{ row.service?.serviceNo }} / {{ row.service?.customer?.name }}
              </template>
            </el-table-column>
            <el-table-column prop="type" label="介质" width="120" />
            <el-table-column prop="sizeGb" label="容量(GB)" width="120" />
            <el-table-column label="挂载信息" min-width="160">
              <template #default="{ row }">
                <div>{{ row.mountPoint || "未挂载" }}</div>
                <div class="muted-line">{{ row.isSystem ? "系统盘" : "数据盘" }}</div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">
                  {{ getLabel(resourceStatusMap, row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" min-width="260" fixed="right">
              <template #default="{ row }">
                <div class="action-grid">
                  <el-button size="small" @click="createSnapshot(row)">快照</el-button>
                  <el-button size="small" @click="createBackup(row)">备份</el-button>
                  <el-button size="small" type="warning" plain @click="toggleDiskMount(row)">
                    {{ row.status === "DETACHED" ? "挂载" : "卸载" }}
                  </el-button>
                  <el-button size="small" type="primary" plain @click="setBootDisk(row)">
                    设为系统盘
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="快照" name="snapshots">
          <el-table v-loading="loading" :data="resources.snapshots" stripe>
            <el-table-column prop="name" label="名称" min-width="180" />
            <el-table-column label="所属服务" min-width="180">
              <template #default="{ row }">
                {{ row.service?.serviceNo }}
              </template>
            </el-table-column>
            <el-table-column label="源磁盘" width="160">
              <template #default="{ row }">
                {{ row.sourceDisk?.name || "-" }}
              </template>
            </el-table-column>
            <el-table-column prop="sizeGb" label="容量(GB)" width="120" />
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">
                  {{ getLabel(resourceStatusMap, row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="创建时间" width="140">
              <template #default="{ row }">
                {{ formatDate(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <div class="action-grid">
                  <el-button size="small" type="primary" plain @click="restoreSnapshot(row)">
                    恢复
                  </el-button>
                  <el-button size="small" type="danger" plain @click="deleteSnapshot(row)">
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="备份" name="backups">
          <el-table v-loading="loading" :data="resources.backups" stripe>
            <el-table-column prop="name" label="名称" min-width="180" />
            <el-table-column label="所属服务" min-width="180">
              <template #default="{ row }">
                {{ row.service?.serviceNo }}
              </template>
            </el-table-column>
            <el-table-column prop="sizeGb" label="容量(GB)" width="120" />
            <el-table-column label="过期时间" width="140">
              <template #default="{ row }">
                {{ formatDate(row.expiresAt) }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">
                  {{ getLabel(resourceStatusMap, row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" min-width="240" fixed="right">
              <template #default="{ row }">
                <div class="action-grid">
                  <el-button size="small" type="primary" plain @click="restoreBackup(row)">
                    恢复
                  </el-button>
                  <el-button size="small" type="warning" plain @click="expireBackup(row)">
                    归档
                  </el-button>
                  <el-button size="small" type="danger" plain @click="deleteBackup(row)">
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="安全组" name="securityGroups">
          <el-table v-loading="loading" :data="resources.securityGroups" stripe>
            <el-table-column type="expand">
              <template #default="{ row }">
                <div class="expand-panel">
                  <div class="expand-title">规则明细</div>
                  <el-empty v-if="row.rules.length === 0" description="暂无规则" :image-size="56" />
                  <el-table v-else :data="row.rules" size="small" border>
                    <el-table-column prop="direction" label="方向" width="100" />
                    <el-table-column prop="protocol" label="协议" width="100" />
                    <el-table-column prop="portRange" label="端口" width="120" />
                    <el-table-column prop="sourceCidr" label="来源" min-width="180" />
                    <el-table-column prop="action" label="策略" width="100" />
                    <el-table-column prop="description" label="说明" min-width="160" />
                    <el-table-column label="操作" width="100" fixed="right">
                      <template #default="{ row: rule }">
                        <el-button
                          text
                          type="danger"
                          @click="deleteRule(row.id, rule.id)"
                        >
                          删除
                        </el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="名称" min-width="180" />
            <el-table-column label="所属服务 / VPC" min-width="220">
              <template #default="{ row }">
                {{ row.service?.serviceNo || "共享" }} / {{ row.vpcNetwork?.name || "-" }}
              </template>
            </el-table-column>
            <el-table-column label="规则数" width="100">
              <template #default="{ row }">
                {{ row.rules.length }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">
                  {{ getLabel(resourceStatusMap, row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="创建时间" width="140">
              <template #default="{ row }">
                {{ formatDate(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" min-width="180" fixed="right">
              <template #default="{ row }">
                <div class="action-grid">
                  <el-button size="small" type="primary" plain :icon="Plus" @click="openAddRule(row)">
                    加规则
                  </el-button>
                  <el-button size="small" type="danger" plain @click="deleteGroup(row)">
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="ruleDialogVisible" title="新增安全组规则" width="680px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="方向">
              <el-select v-model="ruleForm.direction" style="width: 100%">
                <el-option label="入站" value="ingress" />
                <el-option label="出站" value="egress" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="协议">
              <el-select v-model="ruleForm.protocol" style="width: 100%">
                <el-option label="TCP" value="tcp" />
                <el-option label="UDP" value="udp" />
                <el-option label="ICMP" value="icmp" />
                <el-option label="ALL" value="all" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="端口范围">
              <el-input v-model="ruleForm.portRange" placeholder="例如 22 或 80-443" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="来源网段">
              <el-input v-model="ruleForm.sourceCidr" placeholder="例如 0.0.0.0/0" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="策略">
              <el-select v-model="ruleForm.ruleAction" style="width: 100%">
                <el-option label="允许" value="ALLOW" />
                <el-option label="拒绝" value="DENY" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="说明">
              <el-input v-model="ruleForm.description" placeholder="例如 SSH 放行" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="ruleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRule">提交规则</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.muted-line {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 12px;
}

.action-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.expand-panel {
  padding: 8px 12px;
  background: var(--card-soft);
  border-radius: 12px;
}

.expand-title {
  margin-bottom: 12px;
  font-size: 13px;
  font-weight: 700;
  color: var(--text-primary);
}
</style>
