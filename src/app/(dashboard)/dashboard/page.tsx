import { redirect } from "next/navigation";

import { getAdminUiUrl } from "@/lib/admin-ui";

export default function DashboardPage() {
  redirect(getAdminUiUrl("/dashboard"));
}
