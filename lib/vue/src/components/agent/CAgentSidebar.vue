<template>
  <Transition
    enter-active-class="transition-transform duration-300 ease-in-out"
    enter-from-class="translate-x-full"
    enter-to-class="translate-x-0"
    leave-active-class="transition-transform duration-300 ease-in-out"
    leave-from-class="translate-x-0"
    leave-to-class="translate-x-full"
  >
    <div
      v-if="isVisible"
      class="right-sidebar flex flex-row"
      :style="{ width: `${drawerWidth}px` }"
    >
      <div
        class="resize-handle w-1 h-full cursor-ew-resize hover:bg-primary/20 transition-colors shrink-0"
        @mousedown="startDrawerResize"
      />
      <div class="flex-1 flex flex-col min-w-0">
        <CAgentChat>
          <template #actions>
            <Button
              icon="pi pi-times"
              severity="secondary"
              variant="text"
              rounded
              size="small"
              @click="isVisible = false"
            />
          </template>
        </CAgentChat>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, inject, watch, onBeforeUnmount } from 'vue'
import { useAgentSidebarStore } from '../../stores/useAgentSidebarStore'
import { useRightSidebarResize } from '../../composables/useRightSidebarResize'
import CAgentChat from './CAgentChat.vue'

const agentStore = useAgentSidebarStore()
const $eventBus = inject<any>('$eventBus', null)
const { drawerWidth, startDrawerResize } = useRightSidebarResize()

const isVisible = computed({
  get: () => agentStore.visible,
  set: value => agentStore.setVisible(value),
})

watch(isVisible, visible => {
  if (visible) {
    $eventBus?.emit('right-sidebar:opened', 'agent')
  }
})

const offSidebar = $eventBus?.on('right-sidebar:opened', (name: string) => {
  if (name !== 'agent') {
    agentStore.setVisible(false)
  }
})

onBeforeUnmount(() => {
  offSidebar?.()
})
</script>
