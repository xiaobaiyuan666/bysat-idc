import { Card } from "@/components/ui/card";
import { StatusChip } from "@/components/ui/status-chip";
import { cycleLabel, formatCurrency, formatDate, orderStatusLabel } from "@/lib/format";
import { requirePortalUser } from "@/lib/portal-auth";
import { getPortalOrdersData } from "@/lib/portal-data";

const orderSourceLabelMap: Record<string, string> = {
  portal: "前台门户",
  sales: "销售录入",
  admin: "后台录入",
  "billing-engine": "计费引擎",
  "mofang-sync": "魔方云同步",
};

const orderTypeLabelMap: Record<string, string> = {
  new: "新购",
  renew: "续费",
  upgrade: "升级",
  manual: "手工",
  import: "导入",
};

function getSearchParam(
  params: Record<string, string | string[] | undefined>,
  key: string,
) {
  const value = params[key];
  return Array.isArray(value) ? value[0] : value;
}

export default async function PortalOrdersPage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const user = await requirePortalUser();
  const orders = await getPortalOrdersData(user.customerId);
  const params = await searchParams;
  const message = getSearchParam(params, "created") ? "订单已创建。" : null;

  return (
    <section className="space-y-6">
      <Card>
        <h2 className="text-2xl font-semibold text-[var(--ink)]">订单中心</h2>
        <p className="mt-3 text-sm leading-7 text-[var(--muted-ink)]">
          查看门户下单、后台录入、计费生成以及从魔方云同步导入的订单记录。
        </p>
        {message ? (
          <div className="mt-4 rounded-[20px] bg-[var(--accent-soft)] px-4 py-3 text-sm text-[var(--ink)]">
            {message}
          </div>
        ) : null}
      </Card>

      <div className="space-y-5">
        {orders.length === 0 ? (
          <Card>
            <div className="rounded-[22px] bg-white/70 px-5 py-10 text-sm text-[var(--muted-ink)]">
              暂无订单记录。
            </div>
          </Card>
        ) : null}

        {orders.map((order) => (
          <Card key={order.id}>
            <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
              <div>
                <div className="flex flex-wrap items-center gap-3">
                  <h3 className="text-xl font-semibold text-[var(--ink)]">{order.orderNo}</h3>
                  <StatusChip
                    label={orderStatusLabel(order.status)}
                    value={order.status}
                  />
                </div>
                <p className="mt-2 text-sm text-[var(--muted-ink)]">
                  下单时间 {formatDate(order.createdAt)} / 来源{" "}
                  {orderSourceLabelMap[order.source] ?? order.source} / 应付{" "}
                  {formatCurrency(order.totalAmount)}
                </p>
              </div>
              <div className="grid gap-3 text-sm text-[var(--muted-ink)] sm:grid-cols-3">
                <div>已付: {formatCurrency(order.paidAmount)}</div>
                <div>支付截止: {formatDate(order.dueDate)}</div>
                <div>订单类型: {orderTypeLabelMap[order.orderType] ?? order.orderType}</div>
              </div>
            </div>

            <div className="mt-6 grid gap-4 xl:grid-cols-2">
              {order.items.map((item) => (
                <div
                  key={item.id}
                  className="rounded-[24px] border border-[var(--border)] bg-white/70 p-5"
                >
                  <p className="text-lg font-semibold text-[var(--ink)]">{item.title}</p>
                  <p className="mt-1 text-sm text-[var(--muted-ink)]">
                    {item.product.name} / {item.plan?.region.name ?? "未分配地域"}
                    {item.plan?.zone ? ` / ${item.plan.zone.name}` : ""}
                  </p>
                  <div className="mt-4 grid gap-3 text-sm text-[var(--muted-ink)] sm:grid-cols-3">
                    <div>数量: {item.quantity}</div>
                    <div>周期: {cycleLabel(item.cycle)}</div>
                    <div>金额: {formatCurrency(item.totalAmount)}</div>
                  </div>
                </div>
              ))}
            </div>

            {order.invoices.length > 0 ? (
              <div className="mt-6 grid gap-4 xl:grid-cols-2">
                {order.invoices.map((invoice) => (
                  <div
                    key={invoice.id}
                    className="rounded-[24px] border border-[var(--border)] bg-white/70 p-5 text-sm text-[var(--muted-ink)]"
                  >
                    <p className="font-semibold text-[var(--ink)]">{invoice.invoiceNo}</p>
                    <p className="mt-1">
                      待付: {formatCurrency(invoice.totalAmount - invoice.paidAmount)}
                    </p>
                    <p className="mt-1">到期: {formatDate(invoice.dueDate)}</p>
                  </div>
                ))}
              </div>
            ) : null}
          </Card>
        ))}
      </div>
    </section>
  );
}
