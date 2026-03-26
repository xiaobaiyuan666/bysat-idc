<script setup lang="ts">
import { computed, onMounted, reactive } from "vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
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

const filters = reactive({
  keyword: ""
});

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

function filterTree(nodes: MenuNode[], keyword: string): MenuNode[] {
  if (!keyword) return nodes;
  return nodes
    .map(item => {
      const children = filterTree(item.children ?? [], keyword);
      const selfMatched =
        item.title.toLowerCase().includes(keyword) ||
        item.path.toLowerCase().includes(keyword) ||
        item.permission.toLowerCase().includes(keyword);
      if (!selfMatched && children.length === 0) return null;
      return {
        ...item,
        children
      };
    })
    .filter(Boolean) as MenuNode[];
}

const filteredItems = computed(() => filterTree(items.value, filters.keyword.trim().toLowerCase()));
const rootCount = computed(() => items.value.length);
const nodeCount = computed(() => countNodes(items.value));
const leafCount = computed(() => countLeaves(items.value));
const permissionCount = computed(() => countPermissions(items.value));

onMounted(() => {
  if (!permissionStore.menus.length || !permissionStore.permissions.length) {
    void permissionStore.load();
  }
});
</script>

<template>
  <PageWorkbench
    eyebrow="系统 / 菜单"
    title="菜单权限工作台"
    subtitle="查看后台菜单树、路由路径和权限点，便于角色授权时核对。"
  >
    <template #metrics>
      <div class="summary-strip">
        <div class="summary-pill"><span>一级菜单</span><strong>{{ rootCount }}</strong></div>
        <div class="summary-pill"><span>菜单节点</span><strong>{{ nodeCount }}</strong></div>
        <div class="summary-pill"><span>叶子节点</span><strong>{{ leafCount }}</strong></div>
        <div class="summary-pill"><span>权限点</span><strong>{{ permissionCount }}</strong></div>
      </div>
    </template>

    <template #filters>
      <div class="filter-bar">
        <el-input v-model="filters.keyword" clearable placeholder="搜索菜单标题、路由或权限点" />
      </div>
    </template>

    <div class="panel-card">
      <div class="section-card__head">
        <strong>菜单树</strong>
        <span class="section-card__meta">当前匹配 {{ countNodes(filteredItems) }} 个节点</span>
      </div>
      <el-table :data="filteredItems" row-key="id" border default-expand-all>
        <el-table-column prop="title" label="菜单标题" min-width="180" />
        <el-table-column prop="path" label="路由路径" min-width="220" />
        <el-table-column prop="permission" label="权限点" min-width="220" />
        <el-table-column prop="icon" label="图标" min-width="120" />
      </el-table>
    </div>
  </PageWorkbench>
</template>
