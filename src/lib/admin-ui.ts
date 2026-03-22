export function getAdminUiUrl(path = "/") {
  const baseUrl = (process.env.ADMIN_UI_BASE_URL ?? "http://localhost:5173").replace(/\/$/, "");
  const targetPath = path.startsWith("/") ? path : `/${path}`;
  return `${baseUrl}${targetPath}`;
}
