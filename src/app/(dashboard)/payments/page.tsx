import { redirect } from "next/navigation";

import { getAdminConsoleUrl } from "@/lib/admin-console";

export default function PaymentsPage() {
  redirect(getAdminConsoleUrl("/reports/overview"));
}
