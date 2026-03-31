<template>
  <Button
    v-tooltip.bottom="$t('reminder.listLabel')"
    icon="pi pi-clock"
    severity="secondary"
    variant="text"
    rounded
    class="relative"
    @click="handleClick"
  >
    <Badge
      v-if="store.activeCount > 0"
      :value="store.activeCount > 9 ? '9+' : store.activeCount"
      class="!absolute !-top-1 !-right-1"
      severity="contrast"
    />
  </Button>
</template>

<script setup>
import { useReminderStore } from '@/stores/reminder'

const store = useReminderStore()

async function handleClick () {
  if (!store.visible) {
    await store.fetchReminders()
  }

  store.toggleVisibility()
}
</script>
