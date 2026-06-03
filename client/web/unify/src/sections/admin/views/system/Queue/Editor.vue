<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="queue"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:queue/${queue.queueID}`"
          :title="queue.queue"
          :target="queue.queue"
        />
      </div>
      <Panel :header="$t('system.queues.editor.info.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup name="queue" :label="$t('system.queues.editor.info.name')" required>
            <InputText id="queue" name="queue" v-model="queue.queue" />
          </CFormGroup>

          <CFormGroup name="consumer" :label="$t('system.queues.editor.info.consumer')" required>
            <Select
              id="consumer"
              v-model="queue.consumer"
              :options="consumerOptions"
              option-label="label"
              option-value="value"
            />
          </CFormGroup>

          <CFormGroup
            name="pollDelay"
            :label="$t('system.queues.editor.info.poll_delay')"
            :description="queue.meta.poll_delay
              ? $t('system.queues.editor.info.poll_delay_set')
              : $t('system.queues.editor.info.poll_delay_empty')
            "
          >
            <InputText
              id="pollDelay"
              name="pollDelay"
              v-model="queue.meta.poll_delay"
              placeholder="1h / 1m15s / 1h90s"
            />
          </CFormGroup>
        </div>
      </Panel>
    </div>

    <CEditorActions :back-to="{ name: 'system.queues' }">
      <CInputDelete
        v-if="isEdit && queue.canDeleteQueue"
        :label="$t('system.queues.editor.delete')"
        :message="$t('general.confirm.delete')"
        :header="queue.queue || queue.queueID"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button
        type="submit"
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
      />
    </CEditorActions>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'

const { CInputDelete } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const queue = ref(null)
const initialQueue = ref(null)

const consumerOptions = computed(() => [
  { label: t('system.queues.editor.info.consumerOptions.store'), value: 'store' },
  { label: t('system.queues.editor.info.consumerOptions.eventbus'), value: 'eventbus' },
  { label: t('system.queues.editor.info.consumerOptions.human'), value: 'corteza' },
  { label: t('system.queues.editor.info.consumerOptions.redis'), value: 'redis' },
])

const isEdit = computed(() => !!route.params.queueID)

const pageTitle = computed(() =>
  isEdit.value ? t('system.queues.editor.title.edit') : t('system.queues.editor.title.create'),
)

const initialValues = computed(() => ({
  queue: queue.value?.queue || '',
  consumer: queue.value?.consumer || '',
  pollDelay: queue.value?.meta?.poll_delay || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.queue || values.queue.trim().length === 0) {
    errors.queue = [{ message: t('general.label.required') }]
  }

  if (!values.consumer || values.consumer.trim().length === 0) {
    errors.consumer = [{ message: t('general.label.required') }]
  }

  if (values.pollDelay && values.pollDelay.trim().length > 0) {
    const durationRegex = /^((\d+h)?(\d+m)?(\d+s)?)$/
    const match = values.pollDelay.trim().match(durationRegex)
    if (!match || match[0] !== values.pollDelay.trim()) {
      errors.pollDelay = [{ message: t('system.queues.editor.info.poll_delay_invalid') }]
    }
  }

  return { errors }
})

function newQueue() {
  return {
    queue: '',
    consumer: 'corteza',
    meta: {
      poll_delay: '',
    },
  }
}

function normalizeQueue(raw) {
  return {
    ...raw,
    meta: {
      poll_delay: raw.meta?.poll_delay || '',
    },
  }
}

async function loadQueue() {
  const queueID = route.params.queueID
  if (!queueID) {
    queue.value = newQueue()
    initialQueue.value = cloneDeep(queue.value)
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.queuesRead({ queueID })
    queue.value = normalizeQueue(raw)
    initialQueue.value = cloneDeep(queue.value)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.queue.fetch.error'))(e)
    router.push({ name: 'system.queues' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document.querySelector('.p-message-error')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  saving.value = true
  try {
    const payload = {
      queue: queue.value.queue,
      consumer: queue.value.consumer,
      meta: queue.value.meta,
    }

    if (isEdit.value) {
      payload.queueID = queue.value.queueID
      const raw = await $SystemAPI.queuesUpdate(payload)
      queue.value = normalizeQueue(raw)
      initialQueue.value = cloneDeep(queue.value)
      $toast.toastSuccess(t('notification.queue.update.success'))
    } else {
      const created = await $SystemAPI.queuesCreate(payload)
      $toast.toastSuccess(t('notification.queue.create.success'))
      markSaved()
      router.push({ name: 'system.queues.edit', params: { queueID: created.queueID } })
    }
  } catch (e) {
    $toast.toastErrorHandler(t(`notification.queue.${isEdit.value ? 'update' : 'create'}.error`))(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.queuesDelete({ queueID: queue.value.queueID })
    $toast.toastSuccess(t('notification.queue.delete.success'))
    router.push({ name: 'system.queues' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.queue.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!queue.value && !!initialQueue.value && !isEqual(queue.value, initialQueue.value),
  messageKey: 'general.editor.unsavedChanges',
})

onMounted(() => loadQueue())
watch(
  () => route.params.queueID,
  () => loadQueue(),
)
</script>
