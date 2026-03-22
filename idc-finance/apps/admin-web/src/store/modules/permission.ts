import { defineStore } from "pinia";
import { fetchMenus, fetchPermissions } from "@/api/auth";

export interface MenuItem {
  id: number;
  parentId: number;
  title: string;
  titleEn?: string;
  name: string;
  path: string;
  icon: string;
  permission: string;
  children?: MenuItem[];
}

function buildMenuTree(flatMenus: MenuItem[]) {
  const map = new Map<number, MenuItem>();
  flatMenus.forEach(item => map.set(item.id, { ...item, children: [] }));
  const roots: MenuItem[] = [];
  map.forEach(item => {
    if (item.parentId === 0) {
      roots.push(item);
      return;
    }
    const parent = map.get(item.parentId);
    if (parent) {
      parent.children?.push(item);
    }
  });
  return roots;
}

export const usePermissionStore = defineStore("permission", {
  state: () => ({
    menus: [] as MenuItem[],
    permissions: [] as string[]
  }),
  actions: {
    async load() {
      const [menus, permissions] = await Promise.all([
        fetchMenus(),
        fetchPermissions()
      ]);
      this.menus = buildMenuTree(menus);
      this.permissions = permissions;
    }
  }
});
