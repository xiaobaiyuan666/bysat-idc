import { redirect } from "next/navigation";

import { getAdminConsoleUrl } from "@/lib/admin-console";

export default function OrdersPage() {
  redirect(getAdminConsoleUrl("/orders/list"));
}
