<script setup lang="ts">
import { computed } from "vue";
import { usePermissionStore } from "@/store";

type MenuNode = {
  id: number;
  parentId: number;
  title: string;
  name: string;
  path: string;
  icon: string;
  permission: string;
  children?: MenuNode[];
};

const permissionStore = usePermissionStore();
const items = computed(() => permissionStore.menus as MenuNode[]);

function countNodes(nodes: MenuNode[]): number {
  return nodes.reduce((sum, item) => sum + 1 + countNodes(item.children ?? []), 0);
}

function countLeaves(nodes: MenuNode[]): number {
  return nodes.reduce((sum, item) => {
    if (!item.children?.length) return sum + 1;
    return sum + countLeaves(item.children);
  }, 0);
}

function countPermissions(nodes: MenuNode[]): number {
  return nodes.reduce((sum, item) => {
    const childCount = item.children?.length ? countPermissions(item.children) : 0;
    return sum + (item.permission ? 1 : 0) + childCount;
  }, 0);
}

const rootCount = computed(() => items.value.length);
const nodeCount = computed(() => countNodes(items.value));
const leafCount = computed(() => countLeaves(items.value));
const permissionCount = computed(() => countPermissions(items.value));
</script>

<template>
  <div class="page-card">
    <div class="page-header">
      <div>
        <h1 class="page-title">菜单管理</h1>
        <p class="page-subtitle">
          当前菜单来自后端菜单接口，用于验证动态菜单链路。这里展示菜单树、路由路径和权限点。
        </p>
      </div>
    </div>

    <div class="summary-strip">
      <div class="summary-pill"><span>一级菜单</span><strong>{{ rootCount }}</strong></div>
      <div class="summary-pill"><span>菜单节点</span><strong>{{ nodeCount }}</strong></div>
      <div class="summary-pill"><span>叶子节点</span><strong>{{ leafCount }}</strong></div>
      <div class="summary-pill"><span>权限点</span><strong>{{ permissionCount }}</strong></div>
    </div>

    <el-table :data="items" row-key="id" border default-expand-all>
      <el-table-column prop="title" label="菜单标题" min-width="180" />
      <el-table-column prop="path" label="路由路径" min-width="220" />
      <el-table-column prop="permission" label="权限点" min-width="220" />
      <el-table-column prop="icon" label="图标" min-width="140" />
    </el-table>
  </div>
</template>
