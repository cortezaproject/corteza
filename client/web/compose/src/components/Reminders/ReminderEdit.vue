<template>
  <div class="flex h-full flex-col">
    <form
      v-if="localReminder"
      class="flex-1 overflow-auto p-4"
      @submit.prevent="$emit('save', localReminder)"
    >
      <div class="space-y-4">
        <div class="space-y-2">
          <label class="text-sm font-medium text-color">{{ $t('reminder.edit.titleLabel') }}</label>
          <InputText
            v-model="localReminder.payload.title"
            data-test-id="input-title"
            class="w-full"
            :placeholder="$t('reminder.edit.titlePlaceholder')"
          />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium text-color">{{ $t('reminder.edit.notesLabel') }}</label>
          <Textarea
            v-model="localReminder.payload.notes"
            data-test-id="textarea-notes"
            rows="6"
            auto-resize
            class="w-full"
            :placeholder="$t('reminder.edit.notesPlaceholder')"
          />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium text-color">
            {{ $t('reminder.edit.remindAtLabel') }}
          </label>
          <CInputDateTime
            v-model="localReminder.remindAt"
            data-test-id="select-remind-at"
            only-future
          />
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium text-color">
            {{ $t('reminder.edit.assigneeLabel') }}
          </label>
          <CInputUser
            v-model="localReminder.assignedTo"
            data-test-id="select-assignee"
            :placeholder="$t('reminder.edit.assigneePlaceholder')"
          />
        </div>

        <div v-if="localReminder.payload?.link" class="space-y-2">
          <label class="text-sm font-medium text-color">{{ $t('reminder.routesTo') }}</label>
          <div class="flex items-center gap-2">
            <InputText
              v-model="localReminder.payload.link.label"
              data-test-id="input-link"
              class="flex-1"
            />
            <Button
              v-if="recordRoute"
              v-tooltip.bottom="$t('reminder.recordPageLink')"
              icon="pi pi-external-link"
              severity="secondary"
              outlined
              @click.prevent="router.push(recordRoute)"
            />
          </div>
        </div>

        <div v-if="localReminder.reminderID !== NoID" class="space-y-3">
          <div class="flex items-center gap-2">
            <Checkbox
              :model-value="!!localReminder.dismissedAt"
              binary
              input-id="dismissed"
              @update:model-value="$emit('dismiss', { reminder: localReminder, value: $event })"
            />
            <label for="dismissed">{{ $t('reminder.dismissed') }}</label>
          </div>

          <div v-if="localReminder.dismissedAt" class="text-sm text-muted-color">
            {{ $t('reminder.dismissedAt') }}: {{ locFullDateTime(localReminder.dismissedAt) }}
          </div>

          <div v-if="localReminder.snoozeCount" class="text-sm text-muted-color">
            {{ $t('reminder.snooze.count') }}: {{ localReminder.snoozeCount }}
          </div>
        </div>
      </div>
    </form>

    <div class="flex items-center justify-between border-t p-4">
      <Button
        data-test-id="button-back"
        severity="secondary"
        outlined
        icon="pi pi-chevron-left"
        size="small"
        :label="$t('general.label.back')"
        @click="$emit('back')"
      />

      <Button
        data-test-id="button-save"
        icon="pi pi-check"
        size="small"
        :label="$t('general.label.save')"
        :disabled="disableSave"
        :loading="processing"
        @click="$emit('save', localReminder)"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, watch, ref } from 'vue'
import { useRouter } from 'vue-router'
import { system, NoID } from '@cortezaproject/corteza-js-next'
import { components, filters } from '@cortezaproject/corteza-vue-next'

const { CInputDateTime, CInputUser } = components
const { locFullDateTime } = filters

const props = defineProps({
  reminder: {
    type: Object,
    default: () => null,
  },
  processing: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['back', 'save', 'dismiss'])

const router = useRouter()
const localReminder = ref(null)

watch(
  () => props.reminder,
  reminder => {
    localReminder.value = reminder
      ? new system.Reminder({
          ...reminder,
          payload: { ...(reminder.payload || {}) },
        })
      : null
  },
  { immediate: true, deep: true },
)

const disableSave = computed(() => {
  return !localReminder.value?.payload?.title?.trim()
})

const recordRoute = computed(() => {
  const link = localReminder.value?.payload?.link
  return link?.params ? { name: link.name || 'page.record', params: link.params } : null
})
</script>
