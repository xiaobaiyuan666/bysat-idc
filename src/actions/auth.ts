"use server";

import { redirect } from "next/navigation";

import { loginWithPassword, logout } from "@/lib/auth";
import { loginSchema } from "@/lib/validation";

export async function signInAction(formData: FormData) {
  const parsed = loginSchema.safeParse({
    email: formData.get("email"),
    password: formData.get("password"),
  });

  if (!parsed.success) {
    redirect("/login?error=1");
  }

  const success = await loginWithPassword(parsed.data.email, parsed.data.password);

  if (!success) {
    redirect("/login?error=1");
  }

  redirect("/dashboard");
}

export async function logoutAction() {
  await logout();
  redirect("/login");
}
