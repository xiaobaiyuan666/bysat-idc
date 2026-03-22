import {
  NotificationChannel,
  NotificationPriority,
  NotificationStatus,
  TaskRunStatus,
  type NotificationMessage,
  type NotificationTemplate,
} from "@prisma/client";

import { db } from "@/lib/db";

type TemplateVariables = Record<
  string,
  string | number | boolean | null | undefined
>;

function serializePayload(payload: unknown) {
  if (payload === undefined || payload === null) {
    return undefined;
  }

  return typeof payload === "string" ? payload : JSON.stringify(payload, null, 2);
}

function renderTemplateText(
  input: string | null | undefined,
  variables: TemplateVariables,
) {
  if (!input) {
    return undefined;
  }

  return input.replace(/\{\{\s*([\w.]+)\s*\}\}/g, (_, key: string) => {
    const value = variables[key];
    return value === undefined || value === null ? "" : String(value);
  });
}

async function resolveTemplate(
  templateCode?: string,
  variables: TemplateVariables = {},
) {
  if (!templateCode) {
    return {
      template: null,
      subject: undefined,
      content: undefined,
    };
  }

  const template = await db.notificationTemplate.findUnique({
    where: {
      code: templateCode,
    },
  });

  if (!template || !template.isActive) {
    return {
      template: null,
      subject: undefined,
      content: undefined,
    };
  }

  return {
    template,
    subject: renderTemplateText(template.subject, variables),
    content: renderTemplateText(template.content, variables),
  };
}

export async function queueNotification(input: {
  templateCode?: string;
  customerId?: string | null;
  createdById?: string | null;
  channel?: NotificationChannel;
  priority?: NotificationPriority;
  recipient: string;
  recipientName?: string;
  subject?: string;
  content?: string;
  module?: string;
  relatedType?: string;
  relatedId?: string;
  variables?: TemplateVariables;
  payload?: unknown;
}) {
  const variables = input.variables ?? {};
  const rendered = await resolveTemplate(input.templateCode, variables);
  const channel =
    input.channel ?? rendered.template?.channel ?? NotificationChannel.SYSTEM;
  const subject = rendered.subject ?? input.subject;
  const content = rendered.content ?? input.content;

  if (!content) {
    throw new Error("NOTIFICATION_CONTENT_REQUIRED");
  }

  return db.$transaction(async (tx) => {
    const message = await tx.notificationMessage.create({
      data: {
        templateId: rendered.template?.id,
        customerId: input.customerId ?? undefined,
        createdById: input.createdById ?? undefined,
        channel,
        priority: input.priority ?? NotificationPriority.NORMAL,
        status: NotificationStatus.PENDING,
        recipient: input.recipient,
        recipientName: input.recipientName,
        subject,
        content,
        module: input.module,
        relatedType: input.relatedType,
        relatedId: input.relatedId,
        payload: serializePayload({
          templateCode: input.templateCode,
          variables,
          extra: input.payload,
        }),
      },
    });

    const task = await tx.asyncTaskJob.create({
      data: {
        queueName: "notifications",
        jobType: `SEND_${channel}`,
        module: input.module ?? "notification",
        status: TaskRunStatus.PENDING,
        targetType: "notification",
        targetId: message.id,
        notificationId: message.id,
        payload: serializePayload({
          recipient: input.recipient,
          channel,
          subject,
        }),
      },
    });

    return {
      message,
      task,
      template: rendered.template,
    };
  });
}

function dispatchResultMessage(
  channel: NotificationChannel,
  notification: Pick<NotificationMessage, "recipient">,
) {
  if (channel === NotificationChannel.EMAIL) {
    return `邮件通知已投递到 ${notification.recipient}`;
  }

  if (channel === NotificationChannel.SMS) {
    return `短信通知已发送到 ${notification.recipient}`;
  }

  if (channel === NotificationChannel.WEBHOOK) {
    return `Webhook 通知已回调 ${notification.recipient}`;
  }

  return `站内通知已推送给 ${notification.recipient}`;
}

export async function processPendingNotifications(limit = 20) {
  const pendingTasks = await db.asyncTaskJob.findMany({
    where: {
      queueName: "notifications",
      status: TaskRunStatus.PENDING,
      availableAt: {
        lte: new Date(),
      },
      notificationId: {
        not: null,
      },
    },
    include: {
      notification: true,
    },
    orderBy: {
      createdAt: "asc",
    },
    take: limit,
  });

  let sent = 0;
  let failed = 0;

  for (const task of pendingTasks) {
    const notification = task.notification;

    if (!notification) {
      await db.asyncTaskJob.update({
        where: { id: task.id },
        data: {
          status: TaskRunStatus.FAILED,
          executedAt: new Date(),
          attempts: task.attempts + 1,
          errorMessage: "关联通知不存在",
        },
      });
      failed += 1;
      continue;
    }

    try {
      const resultMessage = dispatchResultMessage(notification.channel, notification);

      await db.$transaction(async (tx) => {
        await tx.notificationMessage.update({
          where: { id: notification.id },
          data: {
            status: NotificationStatus.SENT,
            sentAt: new Date(),
            errorMessage: null,
          },
        });

        await tx.asyncTaskJob.update({
          where: { id: task.id },
          data: {
            status: TaskRunStatus.SUCCESS,
            executedAt: new Date(),
            attempts: task.attempts + 1,
            result: resultMessage,
            errorMessage: null,
          },
        });
      });

      sent += 1;
    } catch (error) {
      const message = error instanceof Error ? error.message : "通知发送失败";

      await db.$transaction(async (tx) => {
        await tx.notificationMessage.update({
          where: { id: notification.id },
          data: {
            status: NotificationStatus.FAILED,
            errorMessage: message,
          },
        });

        await tx.asyncTaskJob.update({
          where: { id: task.id },
          data: {
            status: TaskRunStatus.FAILED,
            executedAt: new Date(),
            attempts: task.attempts + 1,
            errorMessage: message,
          },
        });
      });

      failed += 1;
    }
  }

  return {
    processed: pendingTasks.length,
    sent,
    failed,
  };
}

export type NotificationCenterData = {
  templates: NotificationTemplate[];
  messages: NotificationMessage[];
};
