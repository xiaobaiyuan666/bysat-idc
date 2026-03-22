import { redirect } from "next/navigation";

import { getAdminConsoleUrl } from "@/lib/admin-console";

export default function ServicesPage() {
  redirect(getAdminConsoleUrl("/services/list"));
}
