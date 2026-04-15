<template>
  <div class="flex h-full p-3 gap-0 overflow-hidden bg-surface-ground">

    <!-- ── Apps column (left) ──────────────────────────────── -->
    <div
      class="column-panel flex flex-col shrink-0"
      :style="{ width: menuWidth + 'px' }"
    >
      <div class="flex items-center gap-2 px-3 py-2 border-b border-surface shrink-0">
        <i class="pi pi-th-large text-primary text-sm" />
        <span class="font-semibold text-base text-color">{{ $t('home.column.apps') }}</span>
      </div>

      <div class="px-3 pt-3 pb-2 shrink-0">
        <CInputSearch v-model="appsQuery" size="small" class="w-full" />
      </div>

      <div class="flex-1 overflow-y-auto px-3 pb-3">
        <CAppList
          variant="list"
          :query="appsQuery"
          :no-apps-text="$t('home.apps.empty')"
          :no-results-text="$t('home.apps.noResults')"
        />
      </div>
    </div>

    <!-- ── Resize handle: menu / agent ──────────────────────── -->
    <div
      class="resize-handle group"
      @mousedown="startMenuResize"
    >
      <div class="resize-handle-bar group-hover:opacity-100" />
    </div>

    <!-- ── Agent column (middle) ────────────────────────────── -->
    <div class="column-panel flex flex-col flex-1 min-w-0">
      <CAgentChat />
    </div>

    <!-- ── Resize handle: agent / notifications ──────────────── -->
    <div
      class="resize-handle group"
      @mousedown="startNotificationsResize"
    >
      <div class="resize-handle-bar group-hover:opacity-100" />
    </div>

    <!-- ── Notifications column (right) ─────────────────────── -->
    <div
      class="column-panel flex flex-col shrink-0"
      :style="{ width: notificationsWidth + 'px' }"
    >
      <CNotificationsPanel />
    </div>

  </div>
</template>

<script setup>
import {
  components,
  useNotificationsStore,
} from '@cortezaproject/corteza-vue-next'
import { ref } from 'vue'
import { useColumnResize } from '../composables/useColumnResize'

const { CInputSearch, CAppList, CAgentChat, CNotificationsPanel } = components

const appsQuery = ref('')

// ── Notifications column ─────────────────────────────────────
useNotificationsStore()

// ── Column resize ────────────────────────────────────────────
const { menuWidth, notificationsWidth, startMenuResize, startNotificationsResize } = useColumnResize()
</script>

<style scoped>
/* Matches .right-sidebar visual style from useTheme.ts */
.column-panel {
  background-color: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: var(--p-border-radius-xl);
  box-shadow: var(--p-overlay-popover-shadow);
  overflow: hidden;
}

/* Slim drag handle between columns */
.resize-handle {
  width: 12px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: ew-resize;
}

.resize-handle-bar {
  width: 2px;
  height: 2rem;
  border-radius: 9999px;
  background-color: var(--p-content-border-color);
  opacity: 0.4;
  transition: opacity 150ms;
}
</style>
