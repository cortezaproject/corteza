<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="queue"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <Panel :header="$t('system.queues.editor.info.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField name="queue" class="flex flex-col gap-2">
            <label for="queue" class="font-medium text-primary">
              {{ $t('system.queues.editor.info.name') }} *
            </label>
            <InputText id="queue" name="queue" v-model="queue.queue" />
            <Message v-if="$form.queue?.invalid" severity="error" size="small" variant="simple">
              {{ $form.queue.error?.message }}
            </Message>
          </FormField>

          <FormField name="consumer" class="flex flex-col gap-2">
            <label for="consumer" class="font-medium text-primary">
              {{ $t('system.queues.editor.info.consumer') }} *
            </label>
            <Select
              id="consumer"
              v-model="queue.consumer"
              :options="consumerOptions"
              option-label="label"
              option-value="value"
            />
            <Message v-if="$form.consumer?.invalid" severity="error" size="small" variant="simple">
              {{ $form.consumer.error?.message }}
            </Message>
          </FormField>

          <div class="flex flex-col gap-2">
            <label for="handler" class="font-medium text-primary">
              {{ $t('system.queues.editor.info.handler') }}
            </label>
            <InputText id="handler" v-model="queue.meta.handler" />
          </div>

          <div class="flex flex-col gap-2">
            <label for="dispatchTimeout" class="font-medium text-primary">
              {{ $t('system.queues.editor.info.dispatchTimeout') }}
            </label>
            <InputNumber id="dispatchTimeout" v-model="queue.meta.dispatch.timeout" :min="0" />
          </div>
        </div>
      </Panel>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.queues' })"
        />
        <div class="flex gap-2">
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
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'

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

const consumerOptions = computed(() => [
  { label: t('system.queues.editor.info.consumerOptions.store'), value: 'store' },
  {
    label: t('system.queues.editor.info.consumerOptions.eventbus'),
    value: 'eventbus',
  },
])

const isEdit = computed(() => !!route.params.queueID)

const pageTitle = computed(() =>
  isEdit.value ? t('system.queues.editor.title.edit') : t('system.queues.editor.title.create'),
)

const initialValues = computed(() => ({
  queue: queue.value?.queue || '',
  consumer: queue.value?.consumer || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.queue || values.queue.trim().length === 0) {
    errors.queue = [{ message: t('general.label.required') }]
  }

  if (!values.consumer || values.consumer.trim().length === 0) {
    errors.consumer = [{ message: t('general.label.required') }]
  }

  return { errors }
})

function newQueue() {
  return {
    queue: '',
    consumer: 'store',
    meta: {
      handler: '',
      dispatch: { timeout: 0 },
    },
  }
}

async function loadQueue() {
  const queueID = route.params.queueID
  if (!queueID) {
    queue.value = newQueue()
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.queuesRead({ queueID })
    queue.value = {
      ...raw,
      meta: {
        handler: raw.meta?.handler || '',
        dispatch: { timeout: raw.meta?.dispatch?.timeout ?? 0 },
      },
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.queue.fetch.error'))(e)
    router.push({ name: 'system.queues' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return

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
      queue.value = {
        ...raw,
        meta: {
          handler: raw.meta?.handler || '',
          dispatch: { timeout: raw.meta?.dispatch?.timeout ?? 0 },
        },
      }
      $toast.toastSuccess(t('notification.queue.update.success'))
    } else {
      const created = await $SystemAPI.queuesCreate(payload)
      $toast.toastSuccess(t('notification.queue.create.success'))
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

onMounted(() => loadQueue())
watch(
  () => route.params.queueID,
  () => loadQueue(),
)
</script>
