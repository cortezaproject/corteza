<template>
  <div
    v-if="store.toasts.length"
    class="pointer-events-none fixed right-4 top-[calc(var(--topbar-height)+1rem)] z-[1250] flex w-[min(28rem,calc(100vw-2rem))] flex-col gap-3"
  >
    <div
      v-for="reminder in store.toasts"
      :key="reminder.reminderID"
      class="pointer-events-auto rounded-xl border bg-surface p-4 shadow-lg"
    >
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0">
          <div class="truncate font-semibold">
            {{ reminder.payload?.title || reminder.payload?.link?.label || reminder.resource }}
          </div>
          <div v-if="reminder.payload?.notes" class="mt-1 whitespace-pre-wrap text-sm text-muted-color">
            {{ reminder.payload.notes }}
          </div>
          <div v-if="reminder.remindAt" class="mt-2 text-xs text-muted-color">
            {{ locFullDateTime(reminder.remindAt) }}
          </div>
        </div>

        <Button
          icon="pi pi-times"
          severity="secondary"
          variant="text"
          rounded
          @click="store.hideToast(reminder.reminderID)"
        />
      </div>

      <div class="mt-4 flex flex-wrap items-center gap-2">
        <Button
          v-if="recordRoute(reminder)"
          size="small"
          severity="secondary"
          outlined
          icon="pi pi-external-link"
          :label="$t('reminder.recordPageLink')"
          @click="router.push(recordRoute(reminder))"
        />
        <Button
          size="small"
          severity="warning"
          :label="$t('reminder.dismiss')"
          @click="store.setDismissed(reminder, true)"
        />
        <Button
          v-for="option in snoozeOptions"
          :key="option.label"
          size="small"
          severity="secondary"
          outlined
          :label="option.label"
          @click="store.snoozeReminder(reminder, option.duration)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { filters } from '@cortezaproject/corteza-vue-next'
import { useReminderStore } from '@/stores/reminder'

const router = useRouter()
const store = useReminderStore()
const { locFullDateTime } = filters

const snoozeOptions = [
  { label: '5m', duration: 1000 * 60 * 5 },
  { label: '15m', duration: 1000 * 60 * 15 },
  { label: '1h', duration: 1000 * 60 * 60 },
  { label: '1d', duration: 1000 * 60 * 60 * 24 },
]

function recordRoute (reminder) {
  const link = reminder?.payload?.link
  return link?.params ? { name: link.name || 'page.record', params: link.params } : null
}
</script>
