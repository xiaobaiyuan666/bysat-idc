import type { Component } from "vue";

export interface MenuNode {
  index: string;
  path: string;
  title: string;
  subtitle?: string;
  icon?: Component;
  activePath?: string;
  children?: MenuNode[];
}

export interface TagItem {
  path: string;
  title: string;
  activePath?: string;
  closable: boolean;
}
