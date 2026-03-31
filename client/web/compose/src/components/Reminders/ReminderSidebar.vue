<template>
  <Transition
    enter-active-class="transition-transform duration-300 ease-in-out"
    enter-from-class="translate-x-full"
    enter-to-class="translate-x-0"
    leave-active-class="transition-transform duration-300 ease-in-out"
    leave-from-class="translate-x-0"
    leave-to-class="translate-x-full"
  >
    <div v-if="isVisible" class="right-sidebar flex flex-col">
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
  </Transition>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, watch } from 'vue'
import { useReminderStore } from '@/stores/reminder'
import ReminderManager from './ReminderManager.vue'

const store = useReminderStore()
const $eventBus = inject('$eventBus', null)

const isVisible = computed({
  get: () => store.visible,
  set: value => store.setVisible(value),
})

const offSidebar = $eventBus?.on('right-sidebar:opened', name => {
  if (name !== 'reminders') {
    store.setVisible(false)
  }
})

watch(isVisible, async visible => {
  if (visible) {
    $eventBus?.emit('right-sidebar:opened', 'reminders')
    await store.fetchReminders()
  } else {
    store.clearEdit()
  }
})

onBeforeUnmount(() => {
  offSidebar?.()
})
</script>
