import { redirect } from "next/navigation";

import { getAdminUiUrl } from "@/lib/admin-ui";

export default function LoginPage() {
  redirect(getAdminUiUrl("/login"));
}
