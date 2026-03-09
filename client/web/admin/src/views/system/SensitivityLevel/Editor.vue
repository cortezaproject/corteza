<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="sensitivityLevel"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <Panel
        :header="$t('system.sensitivityLevel.editor.info.title')"
        toggleable
        :collapsed="false"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField name="name" class="flex flex-col gap-2">
            <label for="name" class="font-medium text-primary">
              {{ $t('system.sensitivityLevel.editor.info.name') }} *
            </label>
            <InputText id="name" name="name" v-model="sensitivityLevel.name" />
            <Message v-if="$form.name?.invalid" severity="error" size="small" variant="simple">
              {{ $form.name.error?.message }}
            </Message>
          </FormField>

          <FormField name="handle" class="flex flex-col gap-2">
            <label for="handle" class="font-medium text-primary">
              {{ $t('system.sensitivityLevel.editor.info.handle.label') }}
            </label>
            <InputText id="handle" name="handle" v-model="sensitivityLevel.handle" />
            <Message v-if="$form.handle?.invalid" severity="error" size="small" variant="simple">
              {{ $form.handle.error?.message }}
            </Message>
          </FormField>

          <div class="flex flex-col gap-2">
            <label for="level" class="font-medium text-primary">
              {{ $t('system.sensitivityLevel.editor.info.level') }}
            </label>
            <InputNumber id="level" v-model="sensitivityLevel.level" :min="0" />
            <small class="text-surface-500">
              {{ $t('system.sensitivityLevel.editor.info.level.hint') }}
            </small>
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
          @click="$router.push({ name: 'system.sensitivityLevels' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && sensitivityLevel.canDeleteDalSensitivityLevel"
            :label="$t('system.sensitivityLevel.editor.info.delete')"
            :message="$t('general.confirm.delete')"
            :header="sensitivityLevel.name || sensitivityLevel.handle"
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
const sensitivityLevel = ref(null)

const isEdit = computed(() => !!route.params.sensitivityLevelID)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.sensitivityLevel.editor.title.edit')
    : t('system.sensitivityLevel.editor.title.create'),
)

const initialValues = computed(() => ({
  name: sensitivityLevel.value?.name || '',
  handle: sensitivityLevel.value?.handle || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [
      { message: t('system.sensitivityLevel.editor.info.handle.invalid-characters') },
    ]
  }

  return { errors }
})

async function loadSensitivityLevel() {
  const sensitivityLevelID = route.params.sensitivityLevelID
  if (!sensitivityLevelID) {
    sensitivityLevel.value = { name: '', handle: '', level: 0 }
    return
  }

  loading.value = true
  try {
    sensitivityLevel.value = await $SystemAPI.dalSensitivityLevelRead({ sensitivityLevelID })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.sensitivityLevel.fetch.error'))(e)
    router.push({ name: 'system.sensitivityLevels' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return

  saving.value = true
  try {
    const payload = {
      name: sensitivityLevel.value.name,
      handle: sensitivityLevel.value.handle,
      level: sensitivityLevel.value.level ?? 0,
    }

    if (isEdit.value) {
      payload.sensitivityLevelID = sensitivityLevel.value.sensitivityLevelID
      sensitivityLevel.value = await $SystemAPI.dalSensitivityLevelUpdate(payload)
      $toast.toastSuccess(t('notification.sensitivityLevel.update.success'))
    } else {
      const created = await $SystemAPI.dalSensitivityLevelCreate(payload)
      $toast.toastSuccess(t('notification.sensitivityLevel.create.success'))
      router.push({
        name: 'system.sensitivityLevels.edit',
        params: { sensitivityLevelID: created.sensitivityLevelID },
      })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(
        `notification.sensitivityLevel.${isEdit.value ? 'update' : 'create'}.error`,
        'Failed to save sensitivity level',
      ),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.dalSensitivityLevelDelete({
      sensitivityLevelID: sensitivityLevel.value.sensitivityLevelID,
    })
    $toast.toastSuccess(t('notification.sensitivityLevel.delete.success'))
    router.push({ name: 'system.sensitivityLevels' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.sensitivityLevel.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

onMounted(() => loadSensitivityLevel())
watch(
  () => route.params.sensitivityLevelID,
  () => loadSensitivityLevel(),
)
</script>
