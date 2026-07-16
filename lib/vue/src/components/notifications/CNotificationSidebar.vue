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
      <CResizeHandle @mousedown="startDrawerResize" />
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
import { computed } from 'vue'
import { useRightSidebarStore } from '../../stores/useRightSidebarStore'
import { useRightSidebarResize } from '../../composables/useRightSidebarResize'
import CNotificationsPanel from './CNotificationsPanel.vue'
import CResizeHandle from '../CResizeHandle.vue'

const rightSidebarStore = useRightSidebarStore()
const { drawerWidth, startDrawerResize } = useRightSidebarResize()

const isVisible = computed({
  get: () => rightSidebarStore.isOpen('notifications'),
  set: value => value ? rightSidebarStore.open('notifications') : rightSidebarStore.close('notifications'),
})
</script>
