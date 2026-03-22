"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

import { requireUser } from "@/lib/auth";
import { db } from "@/lib/db";
import { customerSchema } from "@/lib/validation";
import { makeCode } from "@/lib/utils";

export async function createCustomerAction(formData: FormData) {
  await requireUser();

  const parsed = customerSchema.safeParse({
    name: formData.get("name"),
    companyName: formData.get("companyName") || undefined,
    email: formData.get("email"),
    phone: formData.get("phone") || undefined,
    type: formData.get("type"),
    notes: formData.get("notes") || undefined,
  });

  if (!parsed.success) {
    redirect("/customers");
  }

  await db.customer.create({
    data: {
      customerNo: makeCode("CUS"),
      name: parsed.data.name,
      companyName: parsed.data.companyName,
      email: parsed.data.email,
      phone: parsed.data.phone,
      type: parsed.data.type,
      notes: parsed.data.notes,
      status: "ACTIVE",
      level: parsed.data.type === "RESELLER" ? "代理" : "标准",
    },
  });

  revalidatePath("/customers");
  revalidatePath("/dashboard");
  redirect("/customers");
}
