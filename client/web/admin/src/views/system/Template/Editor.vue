<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="template"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-if="template.canGrant"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:template/${template.templateID}`"
          :title="template.meta?.short || template.handle || template.templateID"
          :target="template.meta?.short || template.handle || template.templateID"
        />
      </div>
      <Card
        :pt="{
          body: { class: 'p-0 flex flex-col h-full min-h-0' },
          content: { class: 'p-0 flex flex-col h-full min-h-0' },
        }"
        class="overflow-hidden flex-1 min-h-0 flex flex-col"
      >
        <template #content>
          <Tabs v-model:value="activeTab" class="flex flex-col h-full min-h-0">
            <TabList class="rounded-t-lg shrink-0">
              <Tab value="basic">{{ $t('system.templates.editor.tabs.basic') }}</Tab>
              <Tab value="content">{{ $t('system.templates.editor.tabs.content') }}</Tab>
              <Tab v-if="isEdit && !template.partial" value="preview">
                {{ $t('system.templates.editor.tabs.preview') }}
              </Tab>
            </TabList>

            <TabPanels class="flex-1 overflow-y-auto min-h-0">
              <TabPanel value="basic">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField name="name" class="flex flex-col gap-2">
                    <label for="name" class="font-medium text-primary">
                      {{ $t('system.templates.editor.info.meta.short') }}
                      <span class="text-red-500">*</span>
                    </label>
                    <InputText id="name" name="name" v-model="template.meta.short" />
                    <Message
                      v-if="$form.name?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.name.error?.message }}
                    </Message>
                  </FormField>

                  <FormField name="handle" class="flex flex-col gap-2">
                    <label for="handle" class="font-medium text-primary">
                      {{ $t('system.templates.editor.info.handle') }}
                    </label>
                    <InputText id="handle" name="handle" v-model="template.handle" />
                    <Message
                      v-if="$form.handle?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.handle.error?.message }}
                    </Message>
                  </FormField>

                  <div class="flex flex-col gap-2">
                    <label for="type" class="font-medium text-primary">
                      {{ $t('system.templates.editor.info.type') }}
                    </label>
                    <Select
                      id="type"
                      v-model="template.type"
                      :options="typeOptions"
                      option-label="label"
                      option-value="value"
                    />
                  </div>

                  <div class="flex flex-col gap-2">
                    <label for="language" class="font-medium text-primary">
                      {{ $t('system.templates.editor.info.language') }}
                    </label>
                    <InputText id="language" v-model="template.language" placeholder="en" />
                  </div>

                  <div class="flex flex-col gap-2 md:col-span-2">
                    <label for="description" class="font-medium text-primary">
                      {{ $t('system.templates.editor.info.meta.description') }}
                    </label>
                    <Textarea id="description" v-model="template.meta.description" rows="2" />
                  </div>

                  <div class="flex items-center gap-3">
                    <ToggleSwitch id="partial" v-model="template.partial" />
                    <label for="partial" class="font-medium text-primary cursor-pointer">
                      {{ $t('system.templates.editor.info.partial') }}
                    </label>
                  </div>
                </div>
              </TabPanel>

              <TabPanel value="content" class="h-full p-0">
                <div class="flex h-full min-h-0 gap-3">
                  <!-- Toolbox sidebar -->
                  <div class="w-64 shrink-0 overflow-y-auto">
                    <CTemplateToolbox :partials="partials" />
                  </div>

                  <!-- Code editor -->
                  <div class="flex-1 min-w-0 flex flex-col">
                    <CCodeEditor
                      v-model="template.template"
                      :language="editorLanguage"
                      min-height="500px"
                      :border="false"
                    />
                  </div>
                </div>
              </TabPanel>

              <TabPanel v-if="isEdit && !template.partial" value="preview">
                <CTemplatePreview :template="template" />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.templates' })"
        />
        <div class="flex gap-2">
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
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
import { components, useUnsavedGuard } from '@cortezaproject/corteza-vue-next'
import { cloneDeep, isEqual } from 'lodash-es'
import CCodeEditor from '@/components/Template/CCodeEditor.vue'
import CTemplateToolbox from '@/components/Template/CTemplateToolbox.vue'
import CTemplatePreview from '@/components/Template/CTemplatePreview.vue'

const { CInputDelete } = components

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
const activeTab = ref('basic')
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
    activeTab.value = 'basic'
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
