"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

import { requireUser } from "@/lib/auth";
import { db } from "@/lib/db";
import { parseMoneyToCent } from "@/lib/format";
import { productSchema } from "@/lib/validation";

export async function createProductAction(formData: FormData) {
  await requireUser();

  const parsed = productSchema.safeParse({
    code: formData.get("code"),
    name: formData.get("name"),
    category: formData.get("category"),
    billingCycle: formData.get("billingCycle"),
    price: parseMoneyToCent(formData.get("price")),
    setupFee: parseMoneyToCent(formData.get("setupFee")),
    stock: Number(formData.get("stock") || 0),
    providerProductId: formData.get("providerProductId") || undefined,
    regionTemplate: formData.get("regionTemplate") || undefined,
    description: formData.get("description") || undefined,
  });

  if (!parsed.success) {
    redirect("/products");
  }

  await db.product.create({
    data: {
      code: parsed.data.code.trim().toUpperCase(),
      name: parsed.data.name,
      category: parsed.data.category,
      billingCycle: parsed.data.billingCycle,
      price: parsed.data.price,
      setupFee: parsed.data.setupFee,
      stock: parsed.data.stock,
      providerProductId: parsed.data.providerProductId,
      regionTemplate: parsed.data.regionTemplate,
      description: parsed.data.description,
      providerType: parsed.data.category === "BARE_METAL" ? "MANUAL" : "MOFANG_CLOUD",
      status: "ACTIVE",
    },
  });

  revalidatePath("/products");
  revalidatePath("/dashboard");
  redirect("/products");
}
