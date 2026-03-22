import { AppShell } from "@/components/app-shell";
import { requireUser } from "@/lib/auth";

export default async function DashboardLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const user = await requireUser();

  return (
    <AppShell
      user={{
        name: user.name,
        email: user.email,
        role: user.role,
      }}
    >
      {children}
    </AppShell>
  );
}
