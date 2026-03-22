import Link from "next/link";

import { Card } from "@/components/ui/card";
import { StatusChip } from "@/components/ui/status-chip";
import {
  cycleLabel,
  formatCurrency,
  formatDate,
  serviceStatusLabel,
} from "@/lib/format";
import { requirePortalUser } from "@/lib/portal-auth";
import { getPortalServicesData } from "@/lib/portal-data";

function isExpiringSoon(value?: Date | string | null) {
  if (!value) {
    return false;
  }

  const diff = new Date(value).getTime() - Date.now();
  return diff >= 0 && diff <= 1000 * 60 * 60 * 24 * 15;
}

export default async function PortalServicesPage() {
  const user = await requirePortalUser();
  const services = await getPortalServicesData(user.customerId);

  const metrics = {
    total: services.length,
    active: services.filter((service) => service.status === "ACTIVE").length,
    expiringSoon: services.filter((service) => isExpiringSoon(service.nextDueDate)).length,
    snapshots: services.reduce((sum, service) => sum + service.snapshots.length, 0),
  };

  return (
    <section className="space-y-6">
      <Card>
        <h2 className="portal-section-title">服务实例</h2>
        <p className="portal-section-subtitle">
          从服务池进入单实例控制台，查看规格、网络、磁盘、快照、备份、账单与操作轨迹。
        </p>

        <div className="mt-5 grid gap-4 md:grid-cols-4">
          <div className="portal-metric-tile">
            <p className="portal-metric-label">总实例数</p>
            <p className="portal-metric-value">{metrics.total}</p>
          </div>
          <div className="portal-metric-tile">
            <p className="portal-metric-label">运行中</p>
            <p className="portal-metric-value">{metrics.active}</p>
          </div>
          <div className="portal-metric-tile">
            <p className="portal-metric-label">15 天内到期</p>
            <p className="portal-metric-value">{metrics.expiringSoon}</p>
          </div>
          <div className="portal-metric-tile">
            <p className="portal-metric-label">快照总数</p>
            <p className="portal-metric-value">{metrics.snapshots}</p>
          </div>
        </div>
      </Card>

      <div className="space-y-5">
        {services.length === 0 ? (
          <Card>
            <div className="rounded-[22px] bg-[var(--panel-soft)] px-5 py-10 text-sm text-[var(--muted-ink)]">
              暂无服务实例。
            </div>
          </Card>
        ) : null}

        {services.map((service) => (
          <Card key={service.id}>
            <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
              <div>
                <div className="flex flex-wrap items-center gap-3">
                  <h3 className="text-xl font-semibold text-[var(--ink)]">{service.name}</h3>
                  <StatusChip
                    label={serviceStatusLabel(service.status)}
                    value={service.status}
                  />
                </div>
                <p className="mt-2 text-sm text-[var(--muted-ink)]">
                  {service.serviceNo} / {service.product.name} /{" "}
                  {service.plan?.name ?? "基础商品"}
                </p>
              </div>
              <div className="grid gap-3 text-sm text-[var(--muted-ink)] sm:grid-cols-3">
                <div>计费周期: {cycleLabel(service.billingCycle)}</div>
                <div>月费用: {formatCurrency(service.monthlyCost)}</div>
                <div>下次到期: {formatDate(service.nextDueDate)}</div>
              </div>
            </div>

            <div className="mt-6 grid gap-4 xl:grid-cols-4">
              <div className="rounded-[22px] border border-[var(--border)] bg-[var(--panel-soft)] p-5 xl:col-span-2">
                <p className="text-sm font-semibold text-[var(--ink)]">基础规格</p>
                <div className="mt-4 grid gap-3 text-sm text-[var(--muted-ink)] sm:grid-cols-2">
                  <div>地域: {service.plan?.region.name ?? service.region ?? "未设置"}</div>
                  <div>可用区: {service.plan?.zone?.name ?? "默认"}</div>
                  <div>
                    规格:
                    {" "}
                    {service.plan?.flavor
                      ? `${service.plan.flavor.cpu}C / ${service.plan.flavor.memoryGb}G`
                      : `${service.cpuCores ?? "-"}C / ${service.memoryGb ?? "-"}G`}
                  </div>
                  <div>镜像: {service.plan?.image?.name ?? "默认镜像"}</div>
                  <div>主机名: {service.hostname ?? "未设置"}</div>
                  <div>最近同步: {formatDate(service.lastSyncAt)}</div>
                </div>
              </div>

              <div className="rounded-[22px] border border-[var(--border)] bg-[var(--panel-soft)] p-5">
                <p className="text-sm font-semibold text-[var(--ink)]">网络资源</p>
                <div className="mt-4 space-y-2 text-sm text-[var(--muted-ink)]">
                  <div>VPC: {service.vpcNetwork?.name ?? "未绑定"}</div>
                  <div>主 IP: {service.ipAddress ?? service.ipAddresses[0]?.address ?? "未分配"}</div>
                  <div>弹性 IP: {service.ipAddresses.length}</div>
                  <div>安全组: {service.securityGroups.length}</div>
                </div>
              </div>

              <div className="rounded-[22px] border border-[var(--border)] bg-[var(--panel-soft)] p-5">
                <p className="text-sm font-semibold text-[var(--ink)]">存储资源</p>
                <div className="mt-4 space-y-2 text-sm text-[var(--muted-ink)]">
                  <div>磁盘: {service.disks.length}</div>
                  <div>快照: {service.snapshots.length}</div>
                  <div>备份: {service.backups.length}</div>
                  <div>工单: {service.tickets.length}</div>
                </div>
              </div>
            </div>

            <div className="mt-5 flex flex-wrap gap-3">
              <Link
                href={`/portal/services/${service.id}`}
                className="inline-flex items-center justify-center rounded-2xl bg-[var(--accent)] px-4 py-2.5 text-sm font-semibold text-white shadow-[0_14px_34px_rgba(36,104,242,0.24)] transition hover:bg-[var(--accent-strong)]"
              >
                进入实例控制台
              </Link>
              <Link
                href="/portal/invoices"
                className="inline-flex items-center justify-center rounded-2xl border border-[var(--border)] bg-[var(--panel-strong)] px-4 py-2.5 text-sm font-semibold text-[var(--muted-ink)] transition hover:bg-[var(--panel-soft)]"
              >
                查看账单
              </Link>
              <Link
                href="/portal/tickets"
                className="inline-flex items-center justify-center rounded-2xl border border-[var(--border)] bg-[var(--panel-strong)] px-4 py-2.5 text-sm font-semibold text-[var(--muted-ink)] transition hover:bg-[var(--panel-soft)]"
              >
                提交工单
              </Link>
            </div>
          </Card>
        ))}
      </div>
    </section>
  );
}
