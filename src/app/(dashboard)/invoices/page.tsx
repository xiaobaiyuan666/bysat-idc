import { redirect } from "next/navigation";

import { getAdminConsoleUrl } from "@/lib/admin-console";

export default function InvoicesPage() {
  redirect(getAdminConsoleUrl("/billing/invoices"));
}
