import { redirect } from "next/navigation";

import { getAdminConsoleUrl } from "@/lib/admin-console";

export default function CustomersPage() {
  redirect(getAdminConsoleUrl("/customer/list"));
}
