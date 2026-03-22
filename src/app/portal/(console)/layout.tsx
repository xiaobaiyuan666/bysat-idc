import { PortalShell } from "@/components/portal-shell";
import { requirePortalUser } from "@/lib/portal-auth";

export default async function PortalConsoleLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const user = await requirePortalUser();

  return (
    <PortalShell
      user={{
        name: user.name,
        email: user.email,
        role: user.role,
      }}
      customer={{
        name: user.customer.name,
        customerNo: user.customer.customerNo,
        creditBalance: user.customer.creditBalance,
      }}
    >
      {children}
    </PortalShell>
  );
}
