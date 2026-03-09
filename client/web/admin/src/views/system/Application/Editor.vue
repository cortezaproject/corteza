<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="application"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <Panel
        :header="$t('system.applications.editor.info.title', 'Basic information')"
        toggleable
        :collapsed="false"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField name="name" class="flex flex-col gap-2 md:col-span-2">
            <label for="name" class="font-medium text-primary">
              {{ $t('system.applications.editor.info.name', 'Name') }} *
            </label>
            <InputText id="name" name="name" v-model="application.name" />
            <Message v-if="$form.name?.invalid" severity="error" size="small" variant="simple">
              {{ $form.name.error?.message }}
            </Message>
          </FormField>

          <div class="flex flex-col gap-2">
            <label for="weight" class="font-medium text-primary">
              {{ $t('system.applications.editor.info.weight', 'Weight') }}
            </label>
            <InputNumber id="weight" v-model="application.weight" :min="0" />
          </div>

          <div class="flex items-center gap-3">
            <ToggleSwitch id="enabled" v-model="application.enabled" />
            <label for="enabled" class="font-medium text-primary cursor-pointer">
              {{ $t('system.applications.editor.info.enabled', 'Enabled') }}
            </label>
          </div>
        </div>
      </Panel>

      <Panel
        :header="$t('system.applications.editor.unify.title', 'Unify integration')"
        toggleable
        :collapsed="false"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-2">
            <label for="unifyName" class="font-medium text-primary">
              {{ $t('system.applications.editor.unify.name.label', 'Name') }}
            </label>
            <InputText id="unifyName" v-model="application.unify.name" />
          </div>

          <div class="flex flex-col gap-2">
            <label for="unifyUrl" class="font-medium text-primary">
              {{ $t('system.applications.editor.unify.url.label', 'URL') }}
            </label>
            <InputText id="unifyUrl" v-model="application.unify.url" />
          </div>

          <div class="flex items-center gap-3">
            <ToggleSwitch id="unifyListed" v-model="application.unify.listed" />
            <label for="unifyListed" class="font-medium text-primary cursor-pointer">
              {{ $t('system.applications.editor.unify.listed', 'Show in app list') }}
            </label>
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
          @click="$router.push({ name: 'system.applications' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && application.canDeleteApplication"
            :label="$t('system.applications.editor.info.delete', 'Delete')"
            :message="$t('general.confirm.delete')"
            :header="application.name"
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
import { system } from '@cortezaproject/corteza-js-next'
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
const application = ref(null)

const isEdit = computed(() => !!route.params.applicationID)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.applications.editor.title.edit', 'Edit Application')
    : t('system.applications.editor.title.create', 'New Application'),
)

const initialValues = computed(() => ({
  name: application.value?.name || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  return { errors }
})

async function loadApplication() {
  const applicationID = route.params.applicationID
  if (!applicationID) {
    application.value = new system.Application({ enabled: true })
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.applicationRead({ applicationID })
    application.value = new system.Application(raw)
  } catch (e) {
    $toast.toastErrorHandler(
      t('notification.application.fetch.error', 'Failed to load application'),
    )(e)
    router.push({ name: 'system.applications' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return

  saving.value = true
  try {
    const payload = {
      name: application.value.name,
      enabled: application.value.enabled,
      weight: application.value.weight,
      unify: application.value.unify,
    }

    if (isEdit.value) {
      payload.applicationID = application.value.applicationID
      const raw = await $SystemAPI.applicationUpdate(payload)
      application.value = new system.Application(raw)
      $toast.toastSuccess(t('notification.application.update.success', 'Application updated'))
    } else {
      const created = await $SystemAPI.applicationCreate(payload)
      $toast.toastSuccess(t('notification.application.create.success', 'Application created'))
      router.push({
        name: 'system.applications.edit',
        params: { applicationID: created.applicationID },
      })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(
        `notification.application.${isEdit.value ? 'update' : 'create'}.error`,
        'Failed to save application',
      ),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.applicationDelete({ applicationID: application.value.applicationID })
    $toast.toastSuccess(t('notification.application.delete.success', 'Application deleted'))
    router.push({ name: 'system.applications' })
  } catch (e) {
    $toast.toastErrorHandler(
      t('notification.application.delete.error', 'Failed to delete application'),
    )(e)
  } finally {
    deleting.value = false
  }
}

onMounted(() => loadApplication())
watch(
  () => route.params.applicationID,
  () => loadApplication(),
)
</script>
