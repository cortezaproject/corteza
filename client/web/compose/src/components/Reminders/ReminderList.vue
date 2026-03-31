<template>
  <div class="flex h-full flex-col">
    <div class="sticky top-0 z-10 border-b bg-surface p-3">
      <Button
        data-test-id="button-add-reminder"
        size="small"
        outlined
        icon="pi pi-plus"
        :label="$t('reminder.add')"
        @click="$emit('create')"
      />
    </div>

    <div class="flex-1 space-y-3 overflow-auto p-3">
      <div
        v-for="reminder in sortedReminders"
        :key="reminder.reminderID"
        class="rounded-xl border p-4 shadow-sm"
        :class="{ 'opacity-60': reminder.dismissedAt }"
      >
        <div class="flex items-start gap-3">
          <Checkbox
            :model-value="!!reminder.dismissedAt"
            binary
            size="small"
            class="mt-1"
            @update:model-value="$emit('dismiss', { reminder, value: $event })"
          />

          <div class="min-w-0 flex-1">
            <div class="truncate font-medium" :class="{ 'line-through': reminder.dismissedAt }">
              {{ reminder.payload?.title || reminder.payload?.link?.label || reminder.resource }}
            </div>

            <div
              v-if="reminder.payload?.notes"
              class="mt-1 whitespace-pre-wrap text-sm text-muted-color"
            >
              {{ reminder.payload.notes }}
            </div>

            <div class="mt-3 flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-color">
              <span v-if="reminder.remindAt">
                {{ $t('reminder.edit.remindAtLabel') }}: {{ locFullDateTime(reminder.remindAt) }}
              </span>
              <span v-if="reminder.dismissedAt">
                {{ $t('reminder.dismissedAt') }}: {{ locFullDateTime(reminder.dismissedAt) }}
              </span>
            </div>
          </div>

          <div class="flex shrink-0 items-center gap-1">
            <Button
              v-if="recordRoute(reminder)"
              v-tooltip.bottom="$t('reminder.recordPageLink')"
              icon="pi pi-external-link"
              severity="secondary"
              variant="text"
              size="small"
              @click="goToRecord(reminder)"
            />
            <Button
              v-tooltip.bottom="$t('reminder.edit.label')"
              icon="pi pi-pencil"
              severity="secondary"
              variant="text"
              size="small"
              @click="$emit('edit', reminder)"
            />
            <Button
              v-tooltip.bottom="$t('reminder.delete')"
              icon="pi pi-trash"
              severity="danger"
              variant="text"
              size="small"
              @click="$emit('delete', reminder)"
            />
          </div>
        </div>
      </div>

      <div v-if="!sortedReminders.length" class="py-10 text-center text-sm text-muted-color">
        {{ $t('reminder.listLabel') }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { filters } from '@cortezaproject/corteza-vue-next'

const props = defineProps({
  reminders: {
    type: Array,
    default: () => [],
  },
})

defineEmits(['create', 'edit', 'dismiss', 'delete'])

const router = useRouter()
const { locFullDateTime } = filters

const sortedReminders = computed(() => {
  return [...props.reminders].sort((a, b) => {
    if (!!a.dismissedAt !== !!b.dismissedAt) {
      return a.dismissedAt ? 1 : -1
    }

    const aTime = a.remindAt ? new Date(a.remindAt).getTime() : Number.MAX_SAFE_INTEGER
    const bTime = b.remindAt ? new Date(b.remindAt).getTime() : Number.MAX_SAFE_INTEGER
    return aTime - bTime
  })
})

function recordRoute(reminder) {
  const link = reminder?.payload?.link
  return link?.params ? { name: link.name || 'page.record', params: link.params } : null
}

function goToRecord(reminder) {
  const route = recordRoute(reminder)
  if (route) {
    router.push(route)
  }
}
</script>
