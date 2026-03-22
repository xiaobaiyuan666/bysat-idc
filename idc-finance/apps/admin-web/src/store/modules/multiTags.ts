import { defineStore } from "pinia";

export interface TagItem {
  title: string;
  titleEn?: string;
  path: string;
}

export const useMultiTagsStore = defineStore("multiTags", {
  state: () => ({
    items: [{ title: "工作台", titleEn: "Workbench", path: "/workbench" }] as TagItem[]
  }),
  actions: {
    push(tag: TagItem) {
      const existing = this.items.find(item => item.path === tag.path);
      if (existing) {
        existing.title = tag.title;
        existing.titleEn = tag.titleEn;
        return;
      }
      this.items.push(tag);
    }
  }
});
