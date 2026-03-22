import Link from "next/link";

import { Card } from "@/components/ui/card";
import { StatusChip } from "@/components/ui/status-chip";
import {
  cycleLabel,
  formatCurrency,
  formatDate,
  invoiceStatusLabel,
  serviceStatusLabel,
  ticketStatusLabel,
} from "@/lib/format";
import { requirePortalUser } from "@/lib/portal-auth";
import { getPortalDashboardData } from "@/lib/portal-data";

export default async function PortalDashboardPage() {
  const user = await requirePortalUser();
  const data = await getPortalDashboardData(user.customerId);

  return (
    <>
      <section className="grid gap-4 xl:grid-cols-4">
        <div className="portal-metric-tile">
          <p className="portal-metric-label">账户余额</p>
          <p className="portal-metric-value">
            {formatCurrency(data.metrics.creditBalance)}
          </p>
          <p className="portal-metric-hint">可直接用于支付待结算账单和实例续费。</p>
        </div>
        <div className="portal-metric-tile">
          <p className="portal-metric-label">运行中服务</p>
          <p className="portal-metric-value">{data.metrics.activeServices}</p>
          <p className="portal-metric-hint">含已开通和正在交付中的云服务实例。</p>
        </div>
        <div className="portal-metric-tile">
          <p className="portal-metric-label">待支付金额</p>
          <p className="portal-metric-value">
            {formatCurrency(data.metrics.openInvoiceAmount)}
          </p>
          <p className="portal-metric-hint">请及时结算，避免服务进入逾期或暂停状态。</p>
        </div>
        <div className="portal-metric-tile">
          <p className="portal-metric-label">处理中工单</p>
          <p className="portal-metric-value">{data.metrics.openTickets}</p>
          <p className="portal-metric-hint">支持查看回复进度并继续补充业务说明。</p>
        </div>
      </section>

      <section className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
        <Card>
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-xl font-semibold text-[var(--ink)]">近期服务实例</h3>
              <p className="mt-2 text-sm leading-7 text-[var(--muted-ink)]">
                展示最近活跃的云服务、网络与关联资源状态。
              </p>
            </div>
            <Link
              href="/portal/services"
              className="text-sm font-semibold text-[var(--accent)]"
            >
              查看全部
            </Link>
          </div>

          <div className="mt-6 space-y-4">
            {data.services.map((service) => (
              <div
                key={service.id}
                className="rounded-[22px] border border-[var(--border)] bg-[var(--panel-soft)] p-5"
              >
                <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                  <div>
                    <p className="text-lg font-semibold text-[var(--ink)]">{service.name}</p>
                    <p className="mt-1 text-sm text-[var(--muted-ink)]">
                      {service.serviceNo} / {service.product.name} /{" "}
                      {service.plan?.region.name ?? service.region ?? "未分配地域"}
                    </p>
                  </div>
                  <StatusChip
                    label={serviceStatusLabel(service.status)}
                    value={service.status}
                  />
                </div>

                <div className="mt-4 grid gap-3 text-sm text-[var(--muted-ink)] md:grid-cols-4">
                  <div>计费周期: {cycleLabel(service.billingCycle)}</div>
                  <div>下次到期: {formatDate(service.nextDueDate)}</div>
                  <div>
                    规格:
                    {" "}
                    {service.plan?.flavor
                      ? `${service.plan.flavor.cpu}C / ${service.plan.flavor.memoryGb}G`
                      : `${service.cpuCores ?? "-"}C / ${service.memoryGb ?? "-"}G`}
                  </div>
                  <div>月费用: {formatCurrency(service.monthlyCost)}</div>
                </div>
              </div>
            ))}
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-xl font-semibold text-[var(--ink)]">推荐云产品</h3>
              <p className="mt-2 text-sm leading-7 text-[var(--muted-ink)]">
                可直接进入产品目录完成选购与下单。
              </p>
            </div>
            <Link
              href="/portal/store"
              className="text-sm font-semibold text-[var(--accent)]"
            >
              去选购
            </Link>
          </div>

          <div className="mt-6 space-y-4">
            {data.plans.map((plan) => (
              <div
                key={plan.id}
                className="rounded-[22px] border border-[var(--border)] bg-[var(--panel-soft)] p-5"
              >
                <p className="text-lg font-semibold text-[var(--ink)]">{plan.name}</p>
                <p className="mt-1 text-sm text-[var(--muted-ink)]">
                  {plan.product.name} / {plan.region.name}
                  {plan.zone ? ` / ${plan.zone.name}` : ""}
                </p>
                <p className="mt-3 text-2xl font-semibold text-[var(--ink)]">
                  {formatCurrency(plan.salePrice)}
                  <span className="ml-2 text-sm font-medium text-[var(--muted-ink)]">
                    / {cycleLabel(plan.billingCycle)}
                  </span>
                </p>
                <p className="mt-3 text-sm text-[var(--muted-ink)]">
                  {plan.flavor
                    ? `${plan.flavor.cpu}C ${plan.flavor.memoryGb}G ${plan.flavor.storageGb ?? 0}G`
                    : "自定义规格"}
                </p>
              </div>
            ))}
          </div>
        </Card>
      </section>

      <section className="grid gap-6 xl:grid-cols-2">
        <Card>
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-xl font-semibold text-[var(--ink)]">待处理账单</h3>
              <p className="mt-2 text-sm leading-7 text-[var(--muted-ink)]">
                优先处理逾期和临近到期账单，避免影响业务交付。
              </p>
            </div>
            <Link
              href="/portal/invoices"
              className="text-sm font-semibold text-[var(--accent)]"
            >
              账单中心
            </Link>
          </div>

          <div className="mt-5 space-y-4">
            {data.invoices.map((invoice) => (
              <div
                key={invoice.id}
                className="flex flex-col gap-3 rounded-[22px] border border-[var(--border)] bg-[var(--panel-soft)] p-5 md:flex-row md:items-center md:justify-between"
              >
                <div>
                  <p className="font-semibold text-[var(--ink)]">{invoice.invoiceNo}</p>
                  <p className="mt-1 text-sm text-[var(--muted-ink)]">
                    应付 {formatCurrency(invoice.totalAmount - invoice.paidAmount)} / 到期{" "}
                    {formatDate(invoice.dueDate)}
                  </p>
                </div>
                <StatusChip
                  label={invoiceStatusLabel(invoice.status)}
                  value={invoice.status}
                />
              </div>
            ))}
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-xl font-semibold text-[var(--ink)]">最近工单</h3>
              <p className="mt-2 text-sm leading-7 text-[var(--muted-ink)]">
                支持查看处理状态、最后回复和关联服务实例。
              </p>
            </div>
            <Link
              href="/portal/tickets"
              className="text-sm font-semibold text-[var(--accent)]"
            >
              工单中心
            </Link>
          </div>

          <div className="mt-5 space-y-4">
            {data.tickets.map((ticket) => (
              <div
                key={ticket.id}
                className="rounded-[22px] border border-[var(--border)] bg-[var(--panel-soft)] p-5"
              >
                <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                  <div>
                    <p className="font-semibold text-[var(--ink)]">{ticket.subject}</p>
                    <p className="mt-1 text-sm text-[var(--muted-ink)]">
                      {ticket.ticketNo} / {ticket.service?.serviceNo ?? "通用咨询"}
                    </p>
                  </div>
                  <StatusChip
                    label={ticketStatusLabel(ticket.status)}
                    value={ticket.status}
                  />
                </div>
                <p className="mt-4 text-sm leading-7 text-[var(--muted-ink)]">
                  {ticket.replies[0]?.content ?? ticket.summary ?? "暂无回复内容"}
                </p>
              </div>
            ))}
          </div>
        </Card>
      </section>
    </>
  );
}
