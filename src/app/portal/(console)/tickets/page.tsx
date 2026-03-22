import { portalCreateTicketAction, portalReplyTicketAction } from "@/actions/portal";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { StatusChip } from "@/components/ui/status-chip";
import { SubmitButton } from "@/components/ui/submit-button";
import { Textarea } from "@/components/ui/textarea";
import {
  formatDateTime,
  ticketPriorityLabel,
  ticketStatusLabel,
} from "@/lib/format";
import { requirePortalUser } from "@/lib/portal-auth";
import { getPortalTicketsData } from "@/lib/portal-data";

function getSearchParam(
  params: Record<string, string | string[] | undefined>,
  key: string,
) {
  const value = params[key];
  return Array.isArray(value) ? value[0] : value;
}

export default async function PortalTicketsPage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const user = await requirePortalUser();
  const params = await searchParams;
  const { tickets, services } = await getPortalTicketsData(user.customerId);

  const message =
    getSearchParam(params, "created")
      ? "工单已提交。"
      : getSearchParam(params, "replied")
        ? "工单回复已提交。"
        : getSearchParam(params, "error")
          ? "工单操作失败，请检查填写内容。"
          : null;

  return (
    <section className="space-y-6">
      <Card>
        <h2 className="text-2xl font-semibold text-[var(--ink)]">工单支持</h2>
        <p className="mt-3 text-sm leading-7 text-[var(--muted-ink)]">
          客户可针对服务实例提交工单，支持多轮回复和持续跟踪。内部备注不会在门户显示。
        </p>
        {message ? (
          <div className="mt-4 rounded-[20px] bg-[var(--accent-soft)] px-4 py-3 text-sm text-[var(--ink)]">
            {message}
          </div>
        ) : null}
      </Card>

      <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
        <Card>
          <h3 className="text-xl font-semibold text-[var(--ink)]">提交新工单</h3>
          <form action={portalCreateTicketAction} className="mt-5 space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium text-[var(--ink)]">关联服务</label>
              <select
                name="serviceId"
                className="w-full rounded-2xl border border-[var(--border)] bg-white px-4 py-3 text-sm text-[var(--ink)] outline-none"
              >
                <option value="">通用咨询</option>
                {services.map((service) => (
                  <option key={service.id} value={service.id}>
                    {service.serviceNo} / {service.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-[var(--ink)]">优先级</label>
              <select
                name="priority"
                defaultValue="NORMAL"
                className="w-full rounded-2xl border border-[var(--border)] bg-white px-4 py-3 text-sm text-[var(--ink)] outline-none"
              >
                <option value="LOW">低</option>
                <option value="NORMAL">中</option>
                <option value="HIGH">高</option>
                <option value="URGENT">紧急</option>
              </select>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-[var(--ink)]">工单标题</label>
              <Input
                name="subject"
                placeholder="例如：华东节点实例重启后无法访问"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-[var(--ink)]">问题描述</label>
              <Textarea
                name="summary"
                placeholder="请描述现象、影响范围、期望处理方式以及问题开始时间"
                required
              />
            </div>

            <SubmitButton className="w-full">提交工单</SubmitButton>
          </form>
        </Card>

        <div className="space-y-5">
          {tickets.length === 0 ? (
            <Card>
              <div className="rounded-[22px] bg-white/70 px-5 py-10 text-sm text-[var(--muted-ink)]">
                暂无工单记录。
              </div>
            </Card>
          ) : null}

          {tickets.map((ticket) => (
            <Card key={ticket.id}>
              <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
                <div>
                  <div className="flex flex-wrap items-center gap-3">
                    <h3 className="text-xl font-semibold text-[var(--ink)]">{ticket.subject}</h3>
                    <StatusChip
                      label={ticketStatusLabel(ticket.status)}
                      value={ticket.status}
                    />
                  </div>
                  <p className="mt-2 text-sm text-[var(--muted-ink)]">
                    {ticket.ticketNo} / 优先级 {ticketPriorityLabel(ticket.priority)} / 服务{" "}
                    {ticket.service?.serviceNo ?? "通用咨询"}
                  </p>
                </div>
                <div className="text-sm text-[var(--muted-ink)]">
                  指派: {ticket.assignedTo?.name ?? "待受理"}
                </div>
              </div>

              <div className="mt-5 space-y-4">
                {ticket.replies.map((reply) => (
                  <div
                    key={reply.id}
                    className="rounded-[22px] border border-[var(--border)] bg-white/70 p-4"
                  >
                    <div className="flex flex-col gap-1 md:flex-row md:items-center md:justify-between">
                      <p className="font-semibold text-[var(--ink)]">{reply.authorName}</p>
                      <p className="text-sm text-[var(--muted-ink)]">
                        {formatDateTime(reply.createdAt)}
                      </p>
                    </div>
                    <p className="mt-3 text-sm leading-6 text-[var(--muted-ink)]">
                      {reply.content}
                    </p>
                  </div>
                ))}
              </div>

              {ticket.status !== "CLOSED" ? (
                <form action={portalReplyTicketAction} className="mt-5 space-y-3">
                  <input type="hidden" name="ticketId" value={ticket.id} />
                  <label className="text-sm font-medium text-[var(--ink)]">继续回复</label>
                  <Textarea name="content" className="min-h-24" required />
                  <SubmitButton>提交回复</SubmitButton>
                </form>
              ) : null}
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
