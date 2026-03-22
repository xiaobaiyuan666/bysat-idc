import { redirect } from "next/navigation";

import { getAdminConsoleUrl } from "@/lib/admin-console";

export default function DashboardPage() {
  redirect(getAdminConsoleUrl("/workbench"));
}
