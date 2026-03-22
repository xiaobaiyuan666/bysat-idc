import { defineStore } from "pinia";

import type { TagItem } from "@/layout/types";

const fixedPaths = new Set(["/workbench"]);

export const useMultiTagsStore = defineStore("pure-multi-tags", {
  state: () => ({
    tags: [] as TagItem[],
  }),
  actions: {
    addTag(tag: TagItem) {
      if (!tag.title) {
        return;
      }

      const exists = this.tags.some((item) => item.path === tag.path);
      if (exists) {
        return;
      }

      this.tags.push({
        ...tag,
        closable: !fixedPaths.has(tag.path),
      });
    },
    removeTag(path: string) {
      this.tags = this.tags.filter((item) => item.path !== path || fixedPaths.has(item.path));
    },
    resetToFixed() {
      this.tags = this.tags.filter((item) => fixedPaths.has(item.path));
    },
  },
});

export function useMultiTagsStoreHook() {
  return useMultiTagsStore();
}
