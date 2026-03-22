import Link from "next/link";
import { notFound } from "next/navigation";

import {
  portalCreateRenewalInvoiceAction,
  portalManageResourceAction,
  portalManageServiceAction,
} from "@/actions/portal";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { StatusChip } from "@/components/ui/status-chip";
import { SubmitButton } from "@/components/ui/submit-button";
import { parseCloudPlanConfig } from "@/lib/cloud-plan-config";
import {
  cycleLabel,
  formatCurrency,
  formatDate,
  formatDateTime,
  invoiceStatusLabel,
  serviceStatusLabel,
  ticketStatusLabel,
} from "@/lib/format";
import { requirePortalUser } from "@/lib/portal-auth";
import { getPortalServiceDetail } from "@/lib/portal-data";
import {
  type ResourceType,
  type SupportedResourceAction as SupportedResourceManagementAction,
} from "@/lib/resource-operations";
import {
  isSupportedServiceAction,
  serviceActionLabelMap,
  type SupportedServiceAction,
} from "@/lib/service-operations";

const resourceActionLabels: Record<ResourceType, Record<string, string>> = {
  disks: {
    createSnapshot: "创建快照",
    createBackup: "创建备份",
    attach: "挂载磁盘",
    detach: "卸载磁盘",
    setBoot: "设为系统盘",
  },
  snapshots: {
    restore: "恢复快照",
    deleteSnapshot: "删除快照",
  },
  backups: {
    restore: "恢复备份",
    expireNow: "立即归档",
    deleteBackup: "删除备份",
  },
  "security-groups": {
    addRule: "新增安全组规则",
    deleteRule: "删除安全组规则",
    deleteGroup: "删除安全组",
  },
};

function getPortalServiceActionLabel(action: SupportedServiceAction) {
  return serviceActionLabelMap[action] ?? action;
}

function getPortalResourceActionLabel(resourceType: ResourceType, action: string) {
  return resourceActionLabels[resourceType]?.[action] ?? action;
}

function parseRuntimeConfig(input?: string | null) {
  if (!input) {
    return {};
  }

  try {
    const parsed = JSON.parse(input) as Record<string, unknown>;
    return parsed && typeof parsed === "object" ? parsed : {};
  } catch {
    return {};
  }
}

function readParam(
  params: Record<string, string | string[] | undefined>,
  key: string,
) {
  const value = params[key];
  return typeof value === "string" ? value : undefined;
}

function numericValue(value: unknown) {
  const normalized = Number(value ?? 0);
  return Number.isFinite(normalized) ? normalized : 0;
}

function ServiceActionButton({
  serviceId,
  action,
  tone = "ghost",
}: {
  serviceId: string;
  action: SupportedServiceAction;
  tone?: "primary" | "secondary" | "danger" | "ghost";
}) {
  return (
    <form action={portalManageServiceAction}>
      <input type="hidden" name="serviceId" value={serviceId} />
      <input type="hidden" name="action" value={action} />
      <SubmitButton tone={tone}>{getPortalServiceActionLabel(action)}</SubmitButton>
    </form>
  );
}

function ResourceActionButton({
  resourceType,
  resourceId,
  action,
  tone = "ghost",
}: {
  resourceType: ResourceType;
  resourceId: string;
  action: SupportedResourceManagementAction;
  tone?: "primary" | "secondary" | "danger" | "ghost";
}) {
  return (
    <form action={portalManageResourceAction}>
      <input type="hidden" name="resourceType" value={resourceType} />
      <input type="hidden" name="resourceId" value={resourceId} />
      <input type="hidden" name="action" value={action} />
      <SubmitButton tone={tone}>{getPortalResourceActionLabel(resourceType, action)}</SubmitButton>
    </form>
  );
}

export default async function PortalServiceDetailPage({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>;
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const portalUser = await requirePortalUser();
  const resolvedParams = await params;
  const resolvedSearchParams = await searchParams;
  const detail = await getPortalServiceDetail(portalUser.customerId, resolvedParams.id);

  if (!detail) {
    notFound();
  }

  const service = detail.service;
  const planConfig = parseCloudPlanConfig(service.plan?.configOptions);
  const runtimeConfig = parseRuntimeConfig(service.configSnapshot);
  const config = {
    ...planConfig,
    ...runtimeConfig,
  };

  const actionValue = readParam(resolvedSearchParams, "action");
  const actionLabelParam = readParam(resolvedSearchParams, "label");
  const actionLabel =
    actionLabelParam ||
    (actionValue && isSupportedServiceAction(actionValue)
      ? getPortalServiceActionLabel(actionValue)
      : "实例操作");

  const resultStatus = readParam(resolvedSearchParams, "status");
  const resultMessage = readParam(resolvedSearchParams, "message");
  const consoleUrl = readParam(resolvedSearchParams, "consoleUrl");
  const taskId = readParam(resolvedSearchParams, "taskId");
  const invoiceNo = readParam(resolvedSearchParams, "invoiceNo");

  const cpu = numericValue(config.cpu ?? service.plan?.flavor?.cpu ?? service.cpuCores);
  const memory = numericValue(
    config.memory ?? service.plan?.flavor?.memoryGb ?? service.memoryGb,
  );
  const systemDisk = numericValue(
    config.system_disk_size ?? service.plan?.flavor?.storageGb ?? service.storageGb,
  );
  const bandwidth = numericValue(config.bw ?? service.plan?.flavor?.bandwidthMbps);

  return (
    <section className="space-y-6">
      <Card>
        <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
          <div>
            <Link
              href="/portal/services"
              className="text-sm font-semibold text-[var(--accent-strong)]"
            >
              返回实例列表
            </Link>
            <div className="mt-3 flex flex-wrap items-center gap-3">
              <h2 className="text-2xl font-semibold text-[var(--ink)]">{service.name}</h2>
              <StatusChip
                label={serviceStatusLabel(service.status)}
                value={service.status}
              />
            </div>
            <p className="mt-2 text-sm text-[var(--muted-ink)]">
              {service.serviceNo} / {service.product.name} /{" "}
              {service.plan?.name ?? "基础商品"} / 主机名 {service.hostname ?? "未设置"}
            </p>
          </div>

          <div className="grid gap-3 text-sm text-[var(--muted-ink)] sm:grid-cols-2">
            <div>计费周期: {cycleLabel(service.billingCycle)}</div>
            <div>月费用: {formatCurrency(service.monthlyCost)}</div>
            <div>下次到期: {formatDate(service.nextDueDate)}</div>
            <div>最后同步: {formatDateTime(service.lastSyncAt)}</div>
          </div>
        </div>

        {resultStatus && resultMessage ? (
          <div
            className={`mt-5 rounded-[22px] px-4 py-4 text-sm ${
              resultStatus === "success"
                ? "bg-[var(--accent-soft)] text-[var(--ink)]"
                : "bg-[rgba(185,76,55,0.12)] text-[#9d3e2d]"
            }`}
          >
            <p className="font-semibold">{actionLabel}</p>
            <p className="mt-1 leading-6">{resultMessage}</p>
            {taskId ? <p className="mt-2">任务 ID: {taskId}</p> : null}
            {consoleUrl ? (
              <p className="mt-2">
                控制台地址:
                <a
                  href={consoleUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="ml-1 font-semibold text-[var(--accent-strong)]"
                >
                  打开 VNC / 控制台
                </a>
              </p>
            ) : null}
            {invoiceNo ? (
              <p className="mt-2">
                续费账单:
                <Link
                  href="/portal/invoices"
                  className="ml-1 font-semibold text-[var(--accent-strong)]"
                >
                  {invoiceNo}
                </Link>
              </p>
            ) : null}
          </div>
        ) : null}
      </Card>

      <div className="grid gap-6 xl:grid-cols-4">
        <Card className="xl:col-span-2">
          <h3 className="text-xl font-semibold text-[var(--ink)]">实例概览</h3>
          <div className="mt-5 grid gap-4 md:grid-cols-2">
            <div className="rounded-[22px] bg-white/70 p-4">
              <p className="text-xs uppercase tracking-[0.22em] text-[var(--muted-ink)]">
                地域与网络
              </p>
              <div className="mt-3 space-y-2 text-sm text-[var(--muted-ink)]">
                <div>地域: {service.plan?.region.name ?? service.region ?? "未设置"}</div>
                <div>可用区: {service.plan?.zone?.name ?? "默认"}</div>
                <div>VPC: {service.vpcNetwork?.name ?? "未绑定"}</div>
                <div>
                  主 IP: {service.ipAddress ?? service.ipAddresses[0]?.address ?? "未分配"}
                </div>
                <div>网络类型: {String(config.network_type ?? "默认网络")}</div>
              </div>
            </div>

            <div className="rounded-[22px] bg-white/70 p-4">
              <p className="text-xs uppercase tracking-[0.22em] text-[var(--muted-ink)]">
                规格与系统
              </p>
              <div className="mt-3 space-y-2 text-sm text-[var(--muted-ink)]">
                <div>计算规格: {cpu > 0 && memory > 0 ? `${cpu}C / ${memory}G` : "未定义"}</div>
                <div>系统盘: {systemDisk > 0 ? `${systemDisk}G` : "未定义"}</div>
                <div>默认镜像: {service.plan?.image?.name ?? "默认镜像"}</div>
                <div>系统标识: {String(config.os ?? service.plan?.image?.code ?? "-")}</div>
                <div>平台节点: {String(config.node ?? "-")}</div>
              </div>
            </div>

            <div className="rounded-[22px] bg-white/70 p-4">
              <p className="text-xs uppercase tracking-[0.22em] text-[var(--muted-ink)]">
                带宽与防护
              </p>
              <div className="mt-3 space-y-2 text-sm text-[var(--muted-ink)]">
                <div>出带宽: {bandwidth > 0 ? `${bandwidth} Mbps` : "未定义"}</div>
                <div>
                  入带宽:{" "}
                  {numericValue(config.in_bw) > 0 ? `${numericValue(config.in_bw)} Mbps` : "未定义"}
                </div>
                <div>IP 数量: {numericValue(config.ip_num) || service.ipAddresses.length || 1}</div>
                <div>
                  流量限制:{" "}
                  {numericValue(config.flow_limit) > 0
                    ? `${numericValue(config.flow_limit)} GB`
                    : "不限"}
                </div>
                <div>
                  防御峰值:{" "}
                  {numericValue(config.peak_defence) > 0
                    ? `${numericValue(config.peak_defence)} Gbps`
                    : "未配置"}
                </div>
              </div>
            </div>

            <div className="rounded-[22px] bg-white/70 p-4">
              <p className="text-xs uppercase tracking-[0.22em] text-[var(--muted-ink)]">
                资源规模
              </p>
              <div className="mt-3 space-y-2 text-sm text-[var(--muted-ink)]">
                <div>磁盘: {service.disks.length}</div>
                <div>快照: {service.snapshots.length}</div>
                <div>备份: {service.backups.length}</div>
                <div>安全组: {service.securityGroups.length}</div>
                <div>关联工单: {service.tickets.length}</div>
              </div>
            </div>
          </div>
        </Card>

        <Card className="xl:col-span-2">
          <h3 className="text-xl font-semibold text-[var(--ink)]">实例操作</h3>
          <p className="mt-2 text-sm leading-7 text-[var(--muted-ink)]">
            自助执行云主机控制、系统重装、密码重置和续费账单生成。
          </p>

          <div className="mt-5 flex flex-wrap gap-3">
            <ServiceActionButton serviceId={service.id} action="sync" />
            <ServiceActionButton serviceId={service.id} action="powerOn" tone="primary" />
            <ServiceActionButton serviceId={service.id} action="powerOff" tone="secondary" />
            <ServiceActionButton serviceId={service.id} action="reboot" tone="secondary" />
            <ServiceActionButton serviceId={service.id} action="hardReboot" tone="danger" />
            {service.status === "SUSPENDED" ? (
              <ServiceActionButton serviceId={service.id} action="unsuspend" tone="primary" />
            ) : (
              <ServiceActionButton serviceId={service.id} action="suspend" tone="secondary" />
            )}
            <ServiceActionButton serviceId={service.id} action="getVnc" tone="primary" />
            <ServiceActionButton serviceId={service.id} action="rescueStart" />
            <ServiceActionButton serviceId={service.id} action="rescueStop" />
            <ServiceActionButton serviceId={service.id} action="lock" />
            <ServiceActionButton serviceId={service.id} action="unlock" />
          </div>

          <div className="mt-6 grid gap-4 lg:grid-cols-2">
            <form
              action={portalManageServiceAction}
              className="rounded-[22px] border border-[var(--border)] bg-white/70 p-5"
            >
              <input type="hidden" name="serviceId" value={service.id} />
              <input type="hidden" name="action" value="resetPassword" />
              <p className="text-base font-semibold text-[var(--ink)]">重置实例密码</p>
              <p className="mt-2 text-sm leading-6 text-[var(--muted-ink)]">
                为云平台实例提交新的登录密码，适用于忘记密码或运维交接场景。
              </p>
              <div className="mt-4 space-y-2">
                <label className="text-sm font-medium text-[var(--ink)]">新密码</label>
                <Input
                  name="password"
                  type="password"
                  required
                  minLength={8}
                  placeholder="至少 8 位，建议包含大小写字母和数字"
                />
              </div>
              <div className="mt-4">
                <SubmitButton>提交重置</SubmitButton>
              </div>
            </form>

            <form
              action={portalManageServiceAction}
              className="rounded-[22px] border border-[var(--border)] bg-white/70 p-5"
            >
              <input type="hidden" name="serviceId" value={service.id} />
              <input type="hidden" name="action" value="reinstall" />
              <p className="text-base font-semibold text-[var(--ink)]">重装系统</p>
              <p className="mt-2 text-sm leading-6 text-[var(--muted-ink)]">
                可选填写镜像编码，留空则按当前售卖方案或平台默认镜像重装。
              </p>
              <div className="mt-4 space-y-2">
                <label className="text-sm font-medium text-[var(--ink)]">镜像编码</label>
                <Input
                  name="imageId"
                  defaultValue={String(config.os ?? service.plan?.image?.code ?? "")}
                  placeholder="例如 ubuntu-2204"
                />
              </div>
              <div className="mt-4">
                <SubmitButton tone="danger">提交重装</SubmitButton>
              </div>
            </form>
          </div>

          <div className="mt-6 rounded-[22px] border border-[var(--border)] bg-white/70 p-5">
            <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <p className="text-base font-semibold text-[var(--ink)]">续费与账务</p>
                <p className="mt-2 text-sm leading-6 text-[var(--muted-ink)]">
                  为当前实例生成续费账单，生成后可在账单中心直接使用余额结清。
                </p>
                {detail.renewalInvoice ? (
                  <p className="mt-2 text-sm text-[var(--muted-ink)]">
                    当前已有续费账单 {detail.renewalInvoice.invoiceNo} / 状态{" "}
                    {invoiceStatusLabel(detail.renewalInvoice.status)}
                  </p>
                ) : null}
              </div>
              <form action={portalCreateRenewalInvoiceAction}>
                <input type="hidden" name="serviceId" value={service.id} />
                <SubmitButton tone="primary">
                  {detail.renewalInvoice ? "查看续费状态" : "生成续费账单"}
                </SubmitButton>
              </form>
            </div>
          </div>
        </Card>
      </div>

      <div className="grid gap-6 xl:grid-cols-3">
        <Card className="xl:col-span-2">
          <h3 className="text-xl font-semibold text-[var(--ink)]">资源清单</h3>

          <div className="mt-5 grid gap-4 lg:grid-cols-2">
            <div className="rounded-[22px] border border-[var(--border)] bg-white/70 p-5">
              <p className="text-sm font-semibold text-[var(--ink)]">IP 地址</p>
              <div className="mt-4 space-y-2 text-sm text-[var(--muted-ink)]">
                {service.ipAddresses.length === 0 ? <div>暂无公网地址</div> : null}
                {service.ipAddresses.map((ip) => (
                  <div key={ip.id} className="rounded-2xl bg-[rgba(19,34,56,0.04)] px-3 py-3">
                    <div className="font-medium text-[var(--ink)]">{ip.address}</div>
                    <div className="mt-1">
                      {ip.version} / {ip.type} / {ip.bandwidthMbps ?? 0} Mbps
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="rounded-[22px] border border-[var(--border)] bg-white/70 p-5">
              <p className="text-sm font-semibold text-[var(--ink)]">磁盘</p>
              <div className="mt-4 space-y-2 text-sm text-[var(--muted-ink)]">
                {service.disks.length === 0 ? <div>暂无磁盘</div> : null}
                {service.disks.map((disk) => (
                  <div key={disk.id} className="rounded-2xl bg-[rgba(19,34,56,0.04)] px-3 py-3">
                    <div className="font-medium text-[var(--ink)]">
                      {disk.name} / {disk.sizeGb}G
                    </div>
                    <div className="mt-1">
                      {disk.type} / {disk.mountPoint ?? "未挂载"} /{" "}
                      {disk.isSystem ? "系统盘" : "数据盘"}
                    </div>
                    <div className="mt-3 flex flex-wrap gap-2">
                      <ResourceActionButton
                        resourceType="disks"
                        resourceId={disk.id}
                        action="createSnapshot"
                      />
                      <ResourceActionButton
                        resourceType="disks"
                        resourceId={disk.id}
                        action="createBackup"
                        tone="secondary"
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="rounded-[22px] border border-[var(--border)] bg-white/70 p-5">
              <p className="text-sm font-semibold text-[var(--ink)]">快照</p>
              <div className="mt-4 space-y-2 text-sm text-[var(--muted-ink)]">
                {service.snapshots.length === 0 ? <div>暂无快照</div> : null}
                {service.snapshots.map((snapshot) => (
                  <div
                    key={snapshot.id}
                    className="rounded-2xl bg-[rgba(19,34,56,0.04)] px-3 py-3"
                  >
                    <div className="font-medium text-[var(--ink)]">{snapshot.name}</div>
                    <div className="mt-1">
                      源磁盘 {snapshot.sourceDisk?.name ?? "-"} / {snapshot.sizeGb ?? 0}G /{" "}
                      {formatDateTime(snapshot.createdAt)}
                    </div>
                    <div className="mt-3 flex flex-wrap gap-2">
                      <ResourceActionButton
                        resourceType="snapshots"
                        resourceId={snapshot.id}
                        action="restore"
                        tone="secondary"
                      />
                      <ResourceActionButton
                        resourceType="snapshots"
                        resourceId={snapshot.id}
                        action="deleteSnapshot"
                        tone="danger"
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="rounded-[22px] border border-[var(--border)] bg-white/70 p-5">
              <p className="text-sm font-semibold text-[var(--ink)]">备份</p>
              <div className="mt-4 space-y-2 text-sm text-[var(--muted-ink)]">
                {service.backups.length === 0 ? <div>暂无备份</div> : null}
                {service.backups.map((backup) => (
                  <div key={backup.id} className="rounded-2xl bg-[rgba(19,34,56,0.04)] px-3 py-3">
                    <div className="font-medium text-[var(--ink)]">{backup.name}</div>
                    <div className="mt-1">
                      {backup.sizeGb ?? 0}G / 到期 {formatDate(backup.expiresAt)}
                    </div>
                    <div className="mt-3 flex flex-wrap gap-2">
                      <ResourceActionButton
                        resourceType="backups"
                        resourceId={backup.id}
                        action="restore"
                        tone="secondary"
                      />
                      <ResourceActionButton
                        resourceType="backups"
                        resourceId={backup.id}
                        action="expireNow"
                      />
                      <ResourceActionButton
                        resourceType="backups"
                        resourceId={backup.id}
                        action="deleteBackup"
                        tone="danger"
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </Card>

        <Card>
          <h3 className="text-xl font-semibold text-[var(--ink)]">安全组</h3>
          <div className="mt-5 space-y-4">
            {service.securityGroups.length === 0 ? (
              <div className="rounded-[22px] bg-white/70 px-4 py-5 text-sm text-[var(--muted-ink)]">
                暂无安全组配置。
              </div>
            ) : null}
            {service.securityGroups.map((group) => (
              <div
                key={group.id}
                className="rounded-[22px] border border-[var(--border)] bg-white/70 p-5"
              >
                <p className="font-semibold text-[var(--ink)]">{group.name}</p>
                <p className="mt-1 text-sm text-[var(--muted-ink)]">
                  规则数 {group.rules.length}
                </p>
                <div className="mt-4 space-y-2 text-sm text-[var(--muted-ink)]">
                  {group.rules.length === 0 ? <div>暂无规则</div> : null}
                  {group.rules.map((rule) => (
                    <div key={rule.id} className="rounded-2xl bg-[rgba(19,34,56,0.04)] px-3 py-3">
                      <div className="font-medium text-[var(--ink)]">
                        {rule.direction} / {rule.protocol} / {rule.portRange}
                      </div>
                      <div className="mt-1">{rule.sourceCidr}</div>
                      <form action={portalManageResourceAction} className="mt-3">
                        <input type="hidden" name="resourceType" value="security-groups" />
                        <input type="hidden" name="resourceId" value={group.id} />
                        <input type="hidden" name="action" value="deleteRule" />
                        <input type="hidden" name="ruleId" value={rule.id} />
                        <SubmitButton tone="danger">删除规则</SubmitButton>
                      </form>
                    </div>
                  ))}
                </div>
                <form action={portalManageResourceAction} className="mt-4 space-y-3">
                  <input type="hidden" name="resourceType" value="security-groups" />
                  <input type="hidden" name="resourceId" value={group.id} />
                  <input type="hidden" name="action" value="addRule" />
                  <div className="grid gap-3 md:grid-cols-2">
                    <Input name="direction" defaultValue="ingress" placeholder="方向 ingress / egress" />
                    <Input name="protocol" defaultValue="tcp" placeholder="协议 tcp / udp / icmp" />
                    <Input name="portRange" defaultValue="22" placeholder="端口，例如 22 或 80-443" />
                    <Input name="sourceCidr" defaultValue="0.0.0.0/0" placeholder="来源网段" />
                  </div>
                  <Input name="description" placeholder="规则说明，例如 SSH 放行" />
                  <SubmitButton tone="primary">新增安全组规则</SubmitButton>
                </form>
              </div>
            ))}
          </div>
        </Card>
      </div>

      <div className="grid gap-6 xl:grid-cols-2">
        <Card>
          <div className="flex items-center justify-between gap-3">
            <h3 className="text-xl font-semibold text-[var(--ink)]">账单与工单</h3>
            <div className="flex gap-3">
              <Link
                href="/portal/invoices"
                className="inline-flex items-center justify-center rounded-2xl border border-[var(--border)] px-4 py-2.5 text-sm font-semibold text-[var(--muted-ink)] transition hover:bg-[rgba(19,34,56,0.04)]"
              >
                账单中心
              </Link>
              <Link
                href="/portal/tickets"
                className="inline-flex items-center justify-center rounded-2xl border border-[var(--border)] px-4 py-2.5 text-sm font-semibold text-[var(--muted-ink)] transition hover:bg-[rgba(19,34,56,0.04)]"
              >
                工单中心
              </Link>
            </div>
          </div>

          <div className="mt-5 space-y-4">
            {service.invoices.map((invoice) => (
              <div
                key={invoice.id}
                className="rounded-[22px] border border-[var(--border)] bg-white/70 p-5"
              >
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <p className="font-semibold text-[var(--ink)]">{invoice.invoiceNo}</p>
                    <p className="mt-1 text-sm text-[var(--muted-ink)]">
                      {formatDate(invoice.dueDate)} / {formatCurrency(invoice.totalAmount)}
                    </p>
                  </div>
                  <StatusChip
                    label={invoiceStatusLabel(invoice.status)}
                    value={invoice.status}
                  />
                </div>
              </div>
            ))}

            {service.tickets.map((ticket) => (
              <div
                key={ticket.id}
                className="rounded-[22px] border border-[var(--border)] bg-white/70 p-5"
              >
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <p className="font-semibold text-[var(--ink)]">{ticket.subject}</p>
                    <p className="mt-1 text-sm text-[var(--muted-ink)]">
                      {ticket.ticketNo} / {formatDateTime(ticket.updatedAt)}
                    </p>
                  </div>
                  <StatusChip
                    label={ticketStatusLabel(ticket.status)}
                    value={ticket.status}
                  />
                </div>
              </div>
            ))}
          </div>
        </Card>

        <Card>
          <h3 className="text-xl font-semibold text-[var(--ink)]">操作轨迹</h3>

          <div className="mt-5 space-y-4">
            <div>
              <p className="text-sm font-semibold text-[var(--ink)]">平台回执</p>
              <div className="mt-3 space-y-3">
                {detail.providerLogs.length === 0 ? (
                  <div className="rounded-[22px] bg-white/70 px-4 py-5 text-sm text-[var(--muted-ink)]">
                    暂无平台回执记录。
                  </div>
                ) : null}
                {detail.providerLogs.map((log) => (
                  <div
                    key={log.id}
                    className="rounded-[22px] border border-[var(--border)] bg-white/70 p-4"
                  >
                    <div className="flex flex-wrap items-center justify-between gap-3">
                      <p className="font-semibold text-[var(--ink)]">{log.action}</p>
                      <StatusChip label={log.status} value={log.status} />
                    </div>
                    <p className="mt-2 text-sm text-[var(--muted-ink)]">
                      {log.message ?? "无详细消息"}
                    </p>
                    <p className="mt-1 text-xs text-[var(--muted-ink)]">
                      {formatDateTime(log.syncedAt)}
                    </p>
                  </div>
                ))}
              </div>
            </div>

            <div>
              <p className="text-sm font-semibold text-[var(--ink)]">业务审计</p>
              <div className="mt-3 space-y-3">
                {detail.auditLogs.length === 0 ? (
                  <div className="rounded-[22px] bg-white/70 px-4 py-5 text-sm text-[var(--muted-ink)]">
                    暂无业务审计记录。
                  </div>
                ) : null}
                {detail.auditLogs.map((log) => (
                  <div
                    key={log.id}
                    className="rounded-[22px] border border-[var(--border)] bg-white/70 p-4"
                  >
                    <p className="font-semibold text-[var(--ink)]">{log.summary}</p>
                    <p className="mt-2 text-sm text-[var(--muted-ink)]">
                      模块 {log.module} / 动作 {log.action}
                    </p>
                    <p className="mt-1 text-xs text-[var(--muted-ink)]">
                      {formatDateTime(log.createdAt)}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </Card>
      </div>
    </section>
  );
}
