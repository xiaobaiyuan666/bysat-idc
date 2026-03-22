import { portalCreateOrderAction } from "@/actions/portal";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { SubmitButton } from "@/components/ui/submit-button";
import { Textarea } from "@/components/ui/textarea";
import { parseCloudPlanConfig } from "@/lib/cloud-plan-config";
import { cycleLabel, formatCurrency } from "@/lib/format";
import { getPortalCatalogData } from "@/lib/portal-data";

type PortalCatalogPlan = Awaited<ReturnType<typeof getPortalCatalogData>>[number];

function numericValue(value: unknown) {
  const normalized = Number(value ?? 0);
  return Number.isFinite(normalized) ? normalized : 0;
}

function getSearchParam(
  params: Record<string, string | string[] | undefined>,
  key: string,
) {
  const value = params[key];
  return Array.isArray(value) ? value[0] : value;
}

function getPlanSpec(plan: PortalCatalogPlan) {
  const config = parseCloudPlanConfig(plan.configOptions);

  return {
    cpu: numericValue(config.cpu ?? plan.flavor?.cpu),
    memory: numericValue(config.memory ?? plan.flavor?.memoryGb),
    systemDisk: numericValue(config.system_disk_size ?? plan.flavor?.storageGb),
    dataDisk: numericValue(config.data_disk_size),
    bandwidth: numericValue(config.bw ?? plan.flavor?.bandwidthMbps),
    ipCount: numericValue(config.ip_num) || 1,
    flowLimit: numericValue(config.flow_limit),
    networkType: String(config.network_type ?? "默认网络"),
    os: String(config.os ?? plan.image?.code ?? ""),
    node: String(config.node ?? ""),
  };
}

export default async function PortalStorePage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const plans = await getPortalCatalogData();
  const params = await searchParams;
  const message =
    getSearchParam(params, "error") === "order"
      ? "订单提交失败，请检查方案和数量。"
      : getSearchParam(params, "error") === "plan"
        ? "当前方案不可下单，请刷新后重试。"
        : null;

  return (
    <section className="space-y-6">
      <Card>
        <h2 className="text-2xl font-semibold text-[var(--ink)]">云产品目录</h2>
        <p className="mt-3 max-w-3xl text-sm leading-7 text-[var(--muted-ink)]">
          这里展示的是已发布的公有云销售方案。每个方案都绑定了地域、可用区、规格和镜像配置，
          前台下单后会生成正式订单、服务实例和待支付账单。
        </p>
        {message ? (
          <div className="mt-4 rounded-[20px] bg-[var(--accent-soft)] px-4 py-3 text-sm text-[var(--ink)]">
            {message}
          </div>
        ) : null}
      </Card>

      <div className="grid gap-6 2xl:grid-cols-2">
        {plans.map((plan) => {
          const spec = getPlanSpec(plan);

          return (
            <Card key={plan.id} className="overflow-hidden p-0">
              <div className="border-b border-[var(--border)] bg-[rgba(255,255,255,0.68)] px-6 py-5">
                <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                  <div>
                    <p className="text-xl font-semibold text-[var(--ink)]">{plan.name}</p>
                    <p className="mt-1 text-sm text-[var(--muted-ink)]">
                      {plan.product.name} / {plan.region.name}
                      {plan.zone ? ` / ${plan.zone.name}` : ""}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-2xl font-semibold text-[var(--ink)]">
                      {formatCurrency(plan.salePrice)}
                    </p>
                    <p className="text-sm text-[var(--muted-ink)]">
                      / {cycleLabel(plan.billingCycle)}
                    </p>
                  </div>
                </div>
              </div>

              <div className="grid gap-6 px-6 py-6 xl:grid-cols-[1fr_0.95fr]">
                <div className="space-y-4 text-sm text-[var(--muted-ink)]">
                  <div className="grid gap-3 sm:grid-cols-2">
                    <div className="rounded-[20px] bg-white/70 p-4">
                      <p className="text-xs uppercase tracking-[0.22em] text-[var(--muted-ink)]">
                        规格
                      </p>
                      <p className="mt-3 text-base font-semibold text-[var(--ink)]">
                        {spec.cpu > 0 && spec.memory > 0
                          ? `${spec.cpu}C / ${spec.memory}G`
                          : "自定义规格"}
                      </p>
                      <p className="mt-1">
                        系统盘 {spec.systemDisk}G / 带宽 {spec.bandwidth}Mbps
                      </p>
                      {spec.dataDisk > 0 ? <p className="mt-1">数据盘 {spec.dataDisk}G</p> : null}
                    </div>
                    <div className="rounded-[20px] bg-white/70 p-4">
                      <p className="text-xs uppercase tracking-[0.22em] text-[var(--muted-ink)]">
                        镜像
                      </p>
                      <p className="mt-3 text-base font-semibold text-[var(--ink)]">
                        {plan.image?.name ?? "系统镜像待选"}
                      </p>
                      <p className="mt-1">
                        {plan.image?.osType ?? "默认"} {plan.image?.version ?? ""}
                      </p>
                      <p className="mt-1">系统标识 {spec.os || "默认"}</p>
                    </div>
                  </div>

                  <div className="rounded-[20px] bg-white/70 p-4">
                    <p className="text-xs uppercase tracking-[0.22em] text-[var(--muted-ink)]">
                      售卖信息
                    </p>
                    <div className="mt-3 grid gap-3 sm:grid-cols-3">
                      <div>库存: {plan.stock}</div>
                      <div>开通费: {formatCurrency(plan.setupFee)}</div>
                      <div>公开销售: {plan.isPublic ? "是" : "否"}</div>
                    </div>
                    <div className="mt-3 grid gap-3 sm:grid-cols-3">
                      <div>网络: {spec.networkType}</div>
                      <div>IP 数量: {spec.ipCount}</div>
                      <div>
                        流量限制: {spec.flowLimit > 0 ? `${spec.flowLimit} GB` : "不限"}
                      </div>
                    </div>
                    {spec.node ? <p className="mt-3 leading-6">平台节点: {spec.node}</p> : null}
                    <p className="mt-3 leading-6">{plan.description || "暂无额外说明。"}</p>
                  </div>
                </div>

                <form
                  action={portalCreateOrderAction}
                  className="space-y-4 rounded-[24px] border border-[var(--border)] bg-white/80 p-5"
                >
                  <input type="hidden" name="planId" value={plan.id} />
                  <div>
                    <p className="text-lg font-semibold text-[var(--ink)]">提交订单</p>
                    <p className="mt-1 text-sm text-[var(--muted-ink)]">
                      下单后将自动生成服务实例和待支付账单。
                    </p>
                  </div>

                  <div className="space-y-2">
                    <label className="text-sm font-medium text-[var(--ink)]">实例名称</label>
                    <Input
                      name="serviceName"
                      defaultValue={`${plan.product.name}-${plan.region.code}`}
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <label className="text-sm font-medium text-[var(--ink)]">购买数量</label>
                    <Input
                      name="quantity"
                      type="number"
                      min="1"
                      max="20"
                      defaultValue="1"
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <label className="text-sm font-medium text-[var(--ink)]">订单备注</label>
                    <Textarea
                      name="notes"
                      className="min-h-24"
                      placeholder="可填写业务用途、预期开通时间或自定义说明"
                    />
                  </div>

                  <SubmitButton className="w-full">立即下单</SubmitButton>
                </form>
              </div>
            </Card>
          );
        })}
      </div>
    </section>
  );
}
