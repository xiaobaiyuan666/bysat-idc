"use server";

import { redirect } from "next/navigation";

import { loginPortalWithPassword, logoutPortal } from "@/lib/portal-auth";
import { loginSchema } from "@/lib/validation";

export async function portalSignInAction(formData: FormData) {
  const parsed = loginSchema.safeParse({
    email: formData.get("email"),
    password: formData.get("password"),
  });

  if (!parsed.success) {
    redirect("/portal/login?error=1");
  }

  const success = await loginPortalWithPassword(parsed.data.email, parsed.data.password);

  if (!success) {
    redirect("/portal/login?error=1");
  }

  redirect("/portal");
}

export async function portalLogoutAction() {
  await logoutPortal();
  redirect("/portal/login");
}
