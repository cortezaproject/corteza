<template>
  <!-- Compose section layout: wraps all compose routes and mounts the
       compose-only overlays while the section is active. The shared
       user/record/namespace/module stores are provided + preloaded by the
       shell, so they're not handled here. -->
  <RouterView />

  <ReminderSidebar />
  <ReminderToastHost />
  <CTranslatorDialog />
</template>

<script setup>
import ReminderSidebar from '@/sections/compose/components/Reminders/ReminderSidebar.vue'
import ReminderToastHost from '@/sections/compose/components/Reminders/ReminderToastHost.vue'
import CTranslatorDialog from '@/sections/compose/components/Translator/CTranslatorDialog.vue'
import { useReminderStore } from '@/sections/compose/stores/reminder'
import { inject, onBeforeUnmount, onMounted } from 'vue'
import { RouterView } from 'vue-router'

const $eventBus = inject('$eventBus', null)

const reminderStore = useReminderStore()

let offRealtime
onMounted(() => {
  // Reminders are compose-specific (not preloaded by the shell).
  reminderStore.fetchReminders()

  // The shell owns the single realtime socket and re-emits every message on
  // the event bus; compose reacts to its own `reminder` messages here.
  offRealtime = $eventBus?.on('realtime', msg => {
    if (msg?.['@type'] === 'reminder') {
      reminderStore.handleRealtimeReminder(msg['@value'])
    }
  })
})

onBeforeUnmount(() => {
  offRealtime?.()
  reminderStore.dispose()
})
</script>
