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
        <CNotificationsPanel>
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
        </CNotificationsPanel>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, inject, onBeforeUnmount, watch } from 'vue'
import { useNotificationsStore } from '../../stores/useNotificationsStore'
import { useRightSidebarResize } from '../../composables/useRightSidebarResize'
import CNotificationsPanel from './CNotificationsPanel.vue'

const notifications = useNotificationsStore()
const $eventBus = inject('$eventBus', null)
const { drawerWidth, startDrawerResize } = useRightSidebarResize()

const isVisible = computed({
  get: () => notifications.visible,
  set: value => notifications.setVisible(value),
})

const offSidebar = $eventBus?.on('right-sidebar:opened', name => {
  if (name !== 'notifications') {
    notifications.setVisible(false)
  }
})

watch(isVisible, visible => {
  if (visible) {
    $eventBus?.emit('right-sidebar:opened', 'notifications')
  }
})

onBeforeUnmount(() => {
  offSidebar?.()
})
</script>
