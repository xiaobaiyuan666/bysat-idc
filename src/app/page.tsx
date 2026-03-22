import { redirect } from "next/navigation";

import { getCurrentPortalUser } from "@/lib/portal-auth";

export default async function HomePage() {
  const portalUser = await getCurrentPortalUser();
  redirect(portalUser ? "/portal" : "/portal/login");
}
