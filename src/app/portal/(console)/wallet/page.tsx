import { Card } from "@/components/ui/card";
import { customerStatusLabel, formatCurrency, formatDateTime } from "@/lib/format";
import { requirePortalUser } from "@/lib/portal-auth";
import { getPortalWalletData } from "@/lib/portal-data";

const transactionTypeMap: Record<string, string> = {
  RECHARGE: "充值",
  CONSUME: "消费",
  REFUND: "退款",
  ADJUSTMENT: "调账",
  AUTO_RENEW: "自动续费",
};

export default async function PortalWalletPage() {
  const user = await requirePortalUser();
  const customer = await getPortalWalletData(user.customerId);

  return (
    <section className="space-y-6">
      <Card>
        <h2 className="portal-section-title">余额与流水</h2>
        <p className="portal-section-subtitle">
          当前余额可直接用于支付账单。所有充值、消费、退款和自动续费都会记录到余额流水中。
        </p>
      </Card>

      <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
        <Card>
          <p className="text-sm text-[var(--muted-ink)]">当前可用余额</p>
          <p className="mt-4 text-4xl font-semibold text-[var(--ink)]">
            {formatCurrency(customer.creditBalance)}
          </p>
          <div className="mt-6 grid gap-4 md:grid-cols-2">
            <div className="rounded-[20px] bg-[var(--panel-soft)] p-4">
              <p className="text-xs uppercase tracking-[0.22em] text-[var(--muted-ink)]">
                客户编号
              </p>
              <p className="mt-3 font-semibold text-[var(--ink)]">{customer.customerNo}</p>
            </div>
            <div className="rounded-[20px] bg-[var(--panel-soft)] p-4">
              <p className="text-xs uppercase tracking-[0.22em] text-[var(--muted-ink)]">
                账户状态
              </p>
              <p className="mt-3 font-semibold text-[var(--ink)]">
                {customerStatusLabel(customer.status)}
              </p>
            </div>
          </div>
        </Card>

        <Card>
          <h3 className="text-xl font-semibold text-[var(--ink)]">最近 20 条流水</h3>
          <div className="mt-5 space-y-4">
            {customer.creditTransactions.length === 0 ? (
              <div className="rounded-[22px] bg-[var(--panel-soft)] px-5 py-8 text-sm text-[var(--muted-ink)]">
                暂无余额流水。
              </div>
            ) : null}

            {customer.creditTransactions.map((transaction) => (
              <div
                key={transaction.id}
                className="rounded-[22px] border border-[var(--border)] bg-[var(--panel-soft)] p-5"
              >
                <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                  <div>
                    <p className="font-semibold text-[var(--ink)]">
                      {transactionTypeMap[transaction.type] ?? transaction.type}
                    </p>
                    <p className="mt-1 text-sm text-[var(--muted-ink)]">
                      {transaction.description}
                    </p>
                  </div>
                  <div className="text-right">
                    <p
                      className={
                        transaction.amount >= 0
                          ? "text-lg font-semibold text-[var(--success)]"
                          : "text-lg font-semibold text-[var(--danger)]"
                      }
                    >
                      {transaction.amount >= 0 ? "+" : ""}
                      {formatCurrency(transaction.amount)}
                    </p>
                    <p className="mt-1 text-sm text-[var(--muted-ink)]">
                      余额 {formatCurrency(transaction.balanceAfter)}
                    </p>
                  </div>
                </div>
                <p className="mt-3 text-sm text-[var(--muted-ink)]">
                  {formatDateTime(transaction.createdAt)}
                </p>
              </div>
            ))}
          </div>
        </Card>
      </div>
    </section>
  );
}
