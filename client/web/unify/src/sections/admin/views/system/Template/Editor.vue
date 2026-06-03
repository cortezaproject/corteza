<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="template"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2">
        <CPermissionsButton
          v-if="template.canGrant"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:template/${template.templateID}`"
          :title="template.meta?.short || template.handle || template.templateID"
          :target="template.meta?.short || template.handle || template.templateID"
        />
      </div>

      <Panel :header="$t('system.templates.editor.tabs.basic')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup name="name" :label="$t('system.templates.editor.info.meta.short')" required>
            <InputText id="name" name="name" v-model="template.meta.short" />
          </CFormGroup>

          <CFormGroup name="handle" :label="$t('system.templates.editor.info.handle')">
            <InputText id="handle" name="handle" v-model="template.handle" />
          </CFormGroup>

          <CFormGroup :label="$t('system.templates.editor.info.type')" input-id="type">
            <Select
              id="type"
              v-model="template.type"
              :options="typeOptions"
              option-label="label"
              option-value="value"
            />
          </CFormGroup>

          <CFormGroup :label="$t('system.templates.editor.info.language')" input-id="language">
            <InputText id="language" v-model="template.language" placeholder="en" />
          </CFormGroup>

          <CFormGroup
            :label="$t('system.templates.editor.info.meta.description')"
            input-id="description"
            class="md:col-span-2"
          >
            <Textarea id="description" v-model="template.meta.description" rows="2" />
          </CFormGroup>

          <CInputToggleCard
            v-model="template.partial"
            :label="$t('system.templates.editor.info.partial')"
            :description="$t('system.templates.editor.info.partialDescription')"
            class="self-start"
          />
        </div>
      </Panel>

      <Panel :header="$t('system.templates.editor.tabs.content')" toggleable :collapsed="false">
        <div class="flex gap-3">
          <!-- Toolbox sidebar -->
          <div class="w-64 shrink-0">
            <CTemplateToolbox :partials="partials" />
          </div>

          <!-- Code editor -->
          <div class="flex-1 min-w-0 flex flex-col">
            <CCodeEditor
              v-model="template.template"
              :language="editorLanguage"
              min-height="500px"
            />
          </div>
        </div>
      </Panel>

      <Panel
        v-if="isEdit && !template.partial"
        :header="$t('system.templates.editor.tabs.preview')"
        toggleable
        :collapsed="false"
      >
        <CTemplatePreview :template="template" />
      </Panel>
    </div>

    <CEditorActions :back-to="{ name: 'system.templates' }">
      <CInputDelete
        v-if="isEdit && template.canDeleteTemplate"
        :label="$t('system.templates.editor.info.delete')"
        :message="$t('general.confirm.delete')"
        :header="template.meta?.short || template.handle"
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
import { system } from '@planetcrust/human-js'
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'
import CCodeEditor from '@/sections/admin/components/Template/CCodeEditor.vue'
import CTemplateToolbox from '@/sections/admin/components/Template/CTemplateToolbox.vue'
import CTemplatePreview from '@/sections/admin/components/Template/CTemplatePreview.vue'

const { CInputDelete, CInputToggleCard } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const template = ref(null)
const initialTemplate = ref(null)
const partials = ref([])

const typeOptions = computed(() => [
  { label: t('system.templates.editor.info.contentType.text_html'), value: 'text/html' },
  { label: t('system.templates.editor.info.contentType.text_plain'), value: 'text/plain' },
  { label: t('system.templates.editor.info.contentType.text_pdf'), value: 'text/pdf' },
])

const isEdit = computed(() => !!route.params.templateID)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.templates.editor.title.edit')
    : t('system.templates.editor.title.create'),
)

const editorLanguage = computed(() => {
  if (template.value?.type === 'text/html') return 'html'
  return 'text'
})

const initialValues = computed(() => ({
  name: template.value?.meta?.short || '',
  handle: template.value?.handle || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [{ message: t('system.templates.editor.info.invalid-handle-characters') }]
  }

  return { errors }
})

async function loadPartials() {
  try {
    const result = await $SystemAPI.templateList({ partial: true, limit: 0 })
    partials.value = (result?.set || []).filter(tpl => tpl.partial)
  } catch (e) {
    console.error('Failed to load partials:', e)
  }
}

async function loadTemplate() {
  const templateID = route.params.templateID
  if (!templateID) {
    template.value = new system.Template({ type: 'text/html' })
    initialTemplate.value = cloneDeep(template.value)
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.templateRead({ templateID })
    template.value = new system.Template(raw)
    initialTemplate.value = cloneDeep(template.value)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.template.fetch.error'))(e)
    router.push({ name: 'system.templates' })
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
      handle: template.value.handle,
      meta: template.value.meta,
      type: template.value.type,
      language: template.value.language,
      partial: template.value.partial,
      template: template.value.template,
    }

    if (isEdit.value) {
      payload.templateID = template.value.templateID
      const raw = await $SystemAPI.templateUpdate(payload)
      template.value = new system.Template(raw)
      initialTemplate.value = cloneDeep(template.value)
      $toast.toastSuccess(t('notification.template.update.success'))
    } else {
      const created = await $SystemAPI.templateCreate(payload)
      $toast.toastSuccess(t('notification.template.create.success'))
      markSaved()
      router.push({ name: 'system.templates.edit', params: { templateID: created.templateID } })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(
        `notification.template.${isEdit.value ? 'update' : 'create'}.error`,
        'Failed to save template',
      ),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.templateDelete({ templateID: template.value.templateID })
    $toast.toastSuccess(t('notification.template.delete.success'))
    router.push({ name: 'system.templates' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.template.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!template.value && !!initialTemplate.value && !isEqual(template.value, initialTemplate.value),
  messageKey: 'general.editor.unsavedChanges',
})

onMounted(() => {
  loadTemplate()
  loadPartials()
})
watch(
  () => route.params.templateID,
  () => loadTemplate(),
)
</script>
