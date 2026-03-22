import http from "./http";

export async function login(payload: { username: string; password: string }) {
  const { data } = await http.post("/auth/login", payload);
  return data.data as { token: string; displayName: string; roles: string[] };
}

export async function fetchMenus() {
  const { data } = await http.get("/menus");
  return data.data as Array<{
    id: number;
    parentId: number;
    title: string;
    titleEn?: string;
    name: string;
    path: string;
    icon: string;
    permission: string;
  }>;
}

export async function fetchPermissions() {
  const { data } = await http.get("/permissions");
  return data.data as string[];
}
