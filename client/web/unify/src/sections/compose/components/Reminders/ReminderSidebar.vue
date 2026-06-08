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
      <div class="flex items-center justify-between pl-3 pt-3 pb-2 pr-1">
        <h3 class="m-0 text-lg font-semibold">
          {{ $t('reminder.listLabel') }}
        </h3>

        <Button
          icon="pi pi-times"
          severity="secondary"
          variant="text"
          rounded
          size="small"
          @click="isVisible = false"
        />
      </div>

      <div class="flex min-h-0 flex-1 flex-col">
        <div class="min-h-0 flex-1">
          <ReminderManager />
        </div>
      </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useReminderStore } from '@/sections/compose/stores/reminder'
import ReminderManager from './ReminderManager.vue'
import { useRightSidebarResize, useRightSidebarStore } from '@planetcrust/human-vue'

const store = useReminderStore()
const rightSidebarStore = useRightSidebarStore()
const { drawerWidth, startDrawerResize } = useRightSidebarResize()

const isVisible = computed({
  get: () => rightSidebarStore.isOpen('reminders'),
  set: value => value ? rightSidebarStore.open('reminders') : rightSidebarStore.close('reminders'),
})

watch(isVisible, async visible => {
  if (visible) {
    await store.fetchReminders()
  } else {
    store.clearEdit()
  }
})
</script>
