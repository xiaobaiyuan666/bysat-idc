import { redirect } from "next/navigation";

import { getAdminUiUrl } from "@/lib/admin-ui";

export default function ProductsPage() {
  redirect(getAdminUiUrl("/products"));
}
