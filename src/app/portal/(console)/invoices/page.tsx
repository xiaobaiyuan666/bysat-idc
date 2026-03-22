import { portalPayInvoiceAction } from "@/actions/portal";
import { Card } from "@/components/ui/card";
import { StatusChip } from "@/components/ui/status-chip";
import { SubmitButton } from "@/components/ui/submit-button";
import {
  cycleLabel,
  formatCurrency,
  formatDate,
  invoiceStatusLabel,
  paymentMethodLabel,
  paymentStatusLabel,
  serviceStatusLabel,
} from "@/lib/format";
import { requirePortalUser } from "@/lib/portal-auth";
import { getPortalInvoicesData } from "@/lib/portal-data";

function getSearchParam(
  params: Record<string, string | string[] | undefined>,
  key: string,
) {
  const value = params[key];
  return Array.isArray(value) ? value[0] : value;
}

export default async function PortalInvoicesPage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const user = await requirePortalUser();
  const params = await searchParams;
  const data = await getPortalInvoicesData(user.customerId);
  const method = getSearchParam(params, "method") || "BALANCE";
  const onlineGateways = data.paymentGateways.filter((item) => item.method !== "BALANCE");

  const message =
    getSearchParam(params, "paid")
      ? `账单已通过 ${paymentMethodLabel(method)} 完成支付。`
      : getSearchParam(params, "error") === "balance"
        ? "账户余额不足，无法完成余额支付。"
        : getSearchParam(params, "error") === "gateway"
          ? "当前支付渠道未启用，请联系财务或切换其他支付方式。"
          : getSearchParam(params, "error")
            ? "支付失败，请稍后重试。"
            : null;

  return (
    <section className="space-y-6">
      <Card>
        <h2 className="text-2xl font-semibold text-[var(--ink)]">账单中心</h2>
        <p className="mt-3 text-sm leading-7 text-[var(--muted-ink)]">
          查看应收账单、续费账单、支付记录和从魔方云同步过来的续费预估。
        </p>
        {message ? (
          <div className="mt-4 rounded-[20px] bg-[var(--accent-soft)] px-4 py-3 text-sm text-[var(--ink)]">
            {message}
          </div>
        ) : null}
      </Card>

      <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
        <Card>
          <h3 className="text-xl font-semibold text-[var(--ink)]">应收账单</h3>
          <div className="mt-5 space-y-4">
            {data.invoices.length === 0 ? (
              <div className="rounded-[24px] border border-dashed border-[var(--border)] bg-white/70 p-5 text-sm text-[var(--muted-ink)]">
                当前没有待处理账单。
              </div>
            ) : null}

            {data.invoices.map((invoice) => {
              const outstanding = Math.max(invoice.totalAmount - invoice.paidAmount, 0);

              return (
                <div
                  key={invoice.id}
                  className="rounded-[24px] border border-[var(--border)] bg-white/70 p-5"
                >
                  <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                    <div>
                      <p className="text-lg font-semibold text-[var(--ink)]">{invoice.invoiceNo}</p>
                      <p className="mt-1 text-sm text-[var(--muted-ink)]">
                        {invoice.service?.name ?? invoice.order?.orderNo ?? "手工账单"} / 到期{" "}
                        {formatDate(invoice.dueDate)}
                      </p>
                    </div>
                    <StatusChip
                      label={invoiceStatusLabel(invoice.status)}
                      value={invoice.status}
                    />
                  </div>

                  <div className="mt-4 grid gap-3 text-sm text-[var(--muted-ink)] md:grid-cols-4">
                    <div>账单金额: {formatCurrency(invoice.totalAmount)}</div>
                    <div>已收金额: {formatCurrency(invoice.paidAmount)}</div>
                    <div>待收金额: {formatCurrency(outstanding)}</div>
                    <div>签发日期: {formatDate(invoice.issuedAt)}</div>
                  </div>

                  {outstanding > 0 ? (
                    <div className="mt-5 grid gap-3">
                      <div className="flex flex-wrap gap-3">
                        <form action={portalPayInvoiceAction}>
                          <input type="hidden" name="invoiceId" value={invoice.id} />
                          <input type="hidden" name="method" value="BALANCE" />
                          <SubmitButton>余额支付</SubmitButton>
                        </form>

                        {onlineGateways.map((gateway) => (
                          <form key={gateway.id} action={portalPayInvoiceAction}>
                            <input type="hidden" name="invoiceId" value={invoice.id} />
                            <input type="hidden" name="method" value={gateway.method} />
                            <SubmitButton tone="secondary">
                              {paymentMethodLabel(gateway.method)}
                            </SubmitButton>
                          </form>
                        ))}
                      </div>
                      <p className="text-xs text-[var(--muted-ink)]">
                        在线渠道会走后台配置好的签名与回调策略，余额支付会直接扣减账户余额。
                      </p>
                    </div>
                  ) : null}
                </div>
              );
            })}
          </div>
        </Card>

        <div className="space-y-6">
          <Card>
            <h3 className="text-xl font-semibold text-[var(--ink)]">支付方式</h3>
            <div className="mt-5 grid gap-4">
              <div className="rounded-[24px] border border-[var(--border)] bg-white/70 p-5">
                <p className="text-sm text-[var(--muted-ink)]">账户余额</p>
                <p className="mt-3 text-3xl font-semibold text-[var(--ink)]">
                  {formatCurrency(data.customer.creditBalance)}
                </p>
                <p className="mt-2 text-sm text-[var(--muted-ink)]">
                  可直接结清账单，也可用于自动续费。
                </p>
              </div>

              {onlineGateways.map((gateway) => (
                <div
                  key={gateway.id}
                  className="rounded-[24px] border border-[var(--border)] bg-white/70 p-5"
                >
                  <p className="font-semibold text-[var(--ink)]">{gateway.name}</p>
                  <p className="mt-1 text-sm text-[var(--muted-ink)]">
                    {paymentMethodLabel(gateway.method)} / {gateway.signType}
                  </p>
                  <p className="mt-3 text-xs leading-6 text-[var(--muted-ink)]">
                    商户号 {gateway.merchantId || "-"}，回调头 {gateway.callbackHeader || "-"}。
                  </p>
                </div>
              ))}
            </div>
          </Card>

          <Card>
            <h3 className="text-xl font-semibold text-[var(--ink)]">收款记录</h3>
            <div className="mt-5 space-y-4">
              {data.payments.length === 0 ? (
                <div className="rounded-[24px] border border-dashed border-[var(--border)] bg-white/70 p-5 text-sm text-[var(--muted-ink)]">
                  当前没有支付记录。
                </div>
              ) : null}

              {data.payments.map((payment) => (
                <div
                  key={payment.id}
                  className="rounded-[24px] border border-[var(--border)] bg-white/70 p-5"
                >
                  <div className="flex items-center justify-between gap-3">
                    <div>
                      <p className="font-semibold text-[var(--ink)]">{payment.paymentNo}</p>
                      <p className="mt-1 text-sm text-[var(--muted-ink)]">
                        {payment.invoice?.invoiceNo ?? "无关联账单"} / {formatDate(payment.paidAt)}
                      </p>
                    </div>
                    <StatusChip
                      label={paymentStatusLabel(payment.status)}
                      value={payment.status}
                    />
                  </div>
                  <div className="mt-4 grid gap-3 text-sm text-[var(--muted-ink)] md:grid-cols-2">
                    <div>支付方式: {paymentMethodLabel(payment.method)}</div>
                    <div>金额: {formatCurrency(payment.amount)}</div>
                  </div>
                </div>
              ))}
            </div>
          </Card>
        </div>
      </div>

      <Card>
        <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
          <div>
            <h3 className="text-xl font-semibold text-[var(--ink)]">续费预览</h3>
            <p className="mt-1 text-sm text-[var(--muted-ink)]">
              这些实例来自魔方云同步，当前没有未结清账单时，会在这里显示下一期续费预估。
            </p>
          </div>
          <div className="text-sm text-[var(--muted-ink)]">
            共 {data.renewalPreviews.length} 条
          </div>
        </div>

        <div className="mt-5 space-y-4">
          {data.renewalPreviews.length === 0 ? (
            <div className="rounded-[24px] border border-dashed border-[var(--border)] bg-white/70 p-5 text-sm text-[var(--muted-ink)]">
              当前没有续费预览，说明相关服务已经生成正式账单，或暂时无需续费。
            </div>
          ) : null}

          {data.renewalPreviews.map((preview) => (
            <div
              key={preview.id}
              className="rounded-[24px] border border-[var(--border)] bg-white/70 p-5"
            >
              <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                <div>
                  <p className="text-lg font-semibold text-[var(--ink)]">
                    {preview.serviceName}
                  </p>
                  <p className="mt-1 text-sm text-[var(--muted-ink)]">
                    {preview.serviceNo} / {preview.productName}
                    {preview.orderNo ? ` / 来源订单 ${preview.orderNo}` : ""}
                  </p>
                </div>
                <StatusChip
                  label={serviceStatusLabel(preview.status)}
                  value={preview.status}
                />
              </div>

              <div className="mt-4 grid gap-3 text-sm text-[var(--muted-ink)] md:grid-cols-4">
                <div>续费周期: {cycleLabel(preview.billingCycle)}</div>
                <div>预计金额: {formatCurrency(preview.amount)}</div>
                <div>预计续费日: {formatDate(preview.dueDate)}</div>
                <div>状态: {serviceStatusLabel(preview.status)}</div>
              </div>
            </div>
          ))}
        </div>
      </Card>
    </section>
  );
}
