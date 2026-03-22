export function getAdminConsoleUrl(path = "/") {
  const baseUrl = (process.env.IDC_FINANCE_ADMIN_BASE_URL ?? "http://localhost:5177").replace(
    /\/$/,
    "",
  );
  const targetPath = path.startsWith("/") ? path : `/${path}`;
  return `${baseUrl}${targetPath}`;
}
