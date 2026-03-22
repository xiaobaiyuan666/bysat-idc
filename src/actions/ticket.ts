"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

import { requireUser } from "@/lib/auth";
import { db } from "@/lib/db";
import { ticketSchema } from "@/lib/validation";
import { makeCode } from "@/lib/utils";

export async function createTicketAction(formData: FormData) {
  const user = await requireUser();

  const parsed = ticketSchema.safeParse({
    customerId: formData.get("customerId"),
    serviceId: formData.get("serviceId") || undefined,
    subject: formData.get("subject"),
    priority: formData.get("priority"),
    summary: formData.get("summary"),
  });

  if (!parsed.success) {
    redirect("/tickets");
  }

  const ticket = await db.ticket.create({
    data: {
      ticketNo: makeCode("TIC"),
      customerId: parsed.data.customerId,
      serviceId: parsed.data.serviceId,
      subject: parsed.data.subject,
      priority: parsed.data.priority,
      status: "OPEN",
      summary: parsed.data.summary,
      assignedToId: user.id,
      lastReplyAt: new Date(),
    },
  });

  await db.ticketReply.create({
    data: {
      ticketId: ticket.id,
      authorType: "STAFF",
      authorName: user.name,
      content: parsed.data.summary,
    },
  });

  ["/dashboard", "/tickets"].forEach((path) => revalidatePath(path));
  redirect("/tickets");
}
