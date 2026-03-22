import { redirect } from "next/navigation";

import { getAdminConsoleUrl } from "@/lib/admin-console";

export default function ProductsPage() {
  redirect(getAdminConsoleUrl("/catalog/products"));
}
