<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ isEdit ? $t('federation.nodes.editor.title.edit') : $t('federation.nodes.editor.title.create') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="node"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::federation:node/${node.nodeID}`"
          :title="node.name || node.nodeID"
          :target="node.name || node.nodeID"
        />
      </div>
      <Card class="shadow">
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <CFormGroup name="name" :label="$t('federation.nodes.editor.info.name')" required class="md:col-span-2">
              <InputText
                id="name"
                name="name"
                v-model="node.name"
                :placeholder="$t('federation.nodes.editor.info.name')"
              />
            </CFormGroup>

            <CFormGroup name="baseURL" :label="$t('federation.nodes.editor.info.baseURL')" required class="md:col-span-2">
              <InputText
                id="baseURL"
                name="baseURL"
                v-model="node.baseURL"
                placeholder="https://..."
              />
            </CFormGroup>

            <CFormGroup
              :label="$t('federation.nodes.editor.info.contact')"
              input-id="contact"
              class="md:col-span-2"
            >
              <InputText
                id="contact"
                v-model="node.contact"
                type="email"
                placeholder="contact@example.com"
              />
            </CFormGroup>
          </div>
        </template>
      </Card>
    </div>

    <CEditorActions :back-to="{ name: 'federation.nodes' }">
      <CInputDelete
        v-if="isEdit && !node.deletedAt"
        :label="$t('federation.nodes.editor.delete')"
        :message="$t('general.confirm.delete')"
        :header="node.name || node.nodeID"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button
        v-if="isEdit"
        :label="$t('federation.nodes.editor.generateURI')"
        icon="pi pi-qrcode"
        severity="secondary"
        outlined
        @click="handleGenerateURI"
      />
      <Button
        type="submit"
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
      />
    </CEditorActions>
  </Form>

  <!-- Generate URI Dialog -->
  <Dialog
    v-model:visible="uriDialogVisible"
    modal
    :header="$t('federation.nodes.editor.generatedURI.title')"
    :style="{ width: '600px' }"
  >
    <div class="flex flex-col gap-4 p-2">
      <div v-if="generatingURI" class="flex items-center justify-center p-4">
        <ProgressSpinner />
      </div>
      <template v-else>
        <p class="text-sm text-muted-color">{{ $t('federation.nodes.editor.generatedURI.description') }}</p>
        <div class="flex flex-col gap-2">
          <Textarea
            :value="generatedURI"
            readonly
            rows="4"
            class="font-mono text-sm w-full"
          />
        </div>
        <div class="flex justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            outlined
            size="small"
            @click="uriDialogVisible = false"
          />
          <Button
            :label="$t('general.label.copy')"
            icon="pi pi-copy"
            size="small"
            @click="copyURI"
          />
        </div>
      </template>
    </div>
  </Dialog>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'

const { CInputDelete } = components

const vueRoute = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const $FederationAPI = inject('$FederationAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const node = ref(null)
const initialNode = ref(null)

const uriDialogVisible = ref(false)
const generatingURI = ref(false)
const generatedURI = ref('')

const isEdit = computed(() => !!vueRoute.params.nodeID)

const initialValues = computed(() => ({
  name: node.value?.name || '',
  baseURL: node.value?.baseURL || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (!values.baseURL || values.baseURL.trim().length === 0) {
    errors.baseURL = [{ message: t('general.label.required') }]
  } else {
    try {
      new URL(values.baseURL)
    } catch {
      errors.baseURL = [{ message: t('federation.nodes.editor.info.baseURLInvalid') }]
    }
  }

  return { errors }
})

async function loadNode() {
  if (!vueRoute.params.nodeID) {
    node.value = { name: '', baseURL: '', contact: '' }
    initialNode.value = cloneDeep(node.value)
    return
  }

  loading.value = true
  try {
    const raw = await $FederationAPI.nodeRead({ nodeID: vueRoute.params.nodeID })
    node.value = raw
    initialNode.value = cloneDeep(node.value)
  } catch (e) {
    $toast.toastErrorHandler(t('federation.nodes.editor.fetch.error'))(e)
    router.push({ name: 'federation.nodes' })
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
      name: node.value.name,
      baseURL: node.value.baseURL,
      contact: node.value.contact,
    }

    if (isEdit.value) {
      payload.nodeID = node.value.nodeID
      await $FederationAPI.nodeUpdate(payload)
      node.value = { ...node.value, ...payload }
      initialNode.value = cloneDeep(node.value)
      $toast.toastSuccess(t('federation.nodes.editor.update.success'))
    } else {
      const created = await $FederationAPI.nodeCreate(payload)
      $toast.toastSuccess(t('federation.nodes.editor.create.success'))
      markSaved()
      router.push({ name: 'federation.nodes.edit', params: { nodeID: created.nodeID } })
    }
  } catch (e) {
    $toast.toastErrorHandler(t(`federation.nodes.editor.${isEdit.value ? 'update' : 'create'}.error`))(e)
  } finally {
    saving.value = false
  }
}

async function handleGenerateURI() {
  uriDialogVisible.value = true
  generatingURI.value = true
  generatedURI.value = ''

  try {
    const result = await $FederationAPI.nodeGenerateUri({ nodeID: node.value.nodeID })
    generatedURI.value = typeof result === 'string' ? result : (result?.uri || JSON.stringify(result))
  } catch (e) {
    $toast.toastErrorHandler(t('federation.nodes.editor.generateURI.error'))(e)
    uriDialogVisible.value = false
  } finally {
    generatingURI.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $FederationAPI.nodeDelete({ nodeID: node.value.nodeID })
    $toast.toastSuccess(t('federation.nodes.editor.delete.success'))
    router.push({ name: 'federation.nodes' })
  } catch (e) {
    $toast.toastErrorHandler(t('federation.nodes.editor.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

function copyURI() {
  navigator.clipboard.writeText(generatedURI.value).catch(() => {})
  $toast.toastSuccess(t('general.label.copied'))
}

const { markSaved } = useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!node.value && !!initialNode.value && !isEqual(node.value, initialNode.value),
  messageKey: 'general.editor.unsavedChanges',
})

onMounted(() => loadNode())
watch(
  () => vueRoute.params.nodeID,
  () => loadNode(),
)
</script>
