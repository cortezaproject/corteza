<template>
  <div class="flex flex-col gap-4 h-full">
    <Message v-if="template.partial" severity="info" :closable="false">
      {{ $t('system.templates.editor.content.preview.partialNotice') }}
    </Message>

    <template v-else>
      <!-- Preview variables/options JSON editor -->
      <CFormGroup
        :label="$t('system.templates.editor.content.preview.title')"
        :description="$t('system.templates.editor.content.preview.description')"
      >
        <CCodeEditor
          v-model="previewData"
          language="json"
          min-height="200px"
        />
      </CFormGroup>

      <!-- Preview buttons -->
      <div class="flex gap-2 flex-wrap">
        <Button
          v-if="canPreviewHTML"
          :label="$t('system.templates.editor.content.preview.html')"
          icon="pi pi-eye"
          severity="secondary"
          outlined
          :loading="previewLoading === 'html'"
          @click="openPreview('html')"
        />
        <Button
          v-if="canPreviewPDF"
          :label="$t('system.templates.editor.content.preview.pdf')"
          icon="pi pi-file-pdf"
          severity="secondary"
          outlined
          :loading="previewLoading === 'pdf'"
          @click="openPreview('pdf')"
        />
        <span v-if="!canPreviewHTML && !canPreviewPDF" class="text-muted-color text-sm">
          {{ $t('system.templates.editor.content.preview.noDrivers') }}
        </span>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import CCodeEditor from './CCodeEditor.vue'

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const props = defineProps({
  template: {
    type: Object,
    required: true,
  },
})

const previewLoading = ref(null)
const availableDrivers = ref([])

const previewData = ref(JSON.stringify({
  variables: {
    param1: 'value1',
    param2: {
      nestedParam1: 'value2',
    },
  },
  options: {
    documentSize: 'A4',
    contentScale: '1',
    orientation: 'portrait',
    margin: '0.3',
  },
}, null, 2))

const canPreviewHTML = computed(() =>
  availableDrivers.value.some(d => (d.outputTypes || []).includes('text/html')),
)

const canPreviewPDF = computed(() =>
  availableDrivers.value.some(d => (d.outputTypes || []).includes('application/pdf')),
)

async function loadDrivers() {
  try {
    const result = await $SystemAPI.templateRenderDrivers()
    availableDrivers.value = result?.set || []
  } catch (e) {
    console.error('Failed to load render drivers:', e)
  }
}

async function openPreview(ext) {
  if (!props.template.templateID) {
    $toast.toastWarning(t('system.templates.editor.content.preview.saveFirst'))
    return
  }

  previewLoading.value = ext
  try {
    const data = JSON.parse(previewData.value)

    const cfg = {
      method: 'post',
      responseType: 'blob',
      url: $SystemAPI.templateRenderEndpoint({
        templateID: props.template.templateID,
        filename: 'preview',
        ext,
      }),
      data,
    }

    const response = await $SystemAPI.api().request(cfg)
    const blob = window.URL.createObjectURL(response.data)
    window.open(blob, '_newtab')
  } catch (e) {
    if (e instanceof SyntaxError) {
      $toast.toastErrorHandler(t('system.templates.editor.content.preview.invalidJSON'))(e)
    } else {
      $toast.toastErrorHandler(t('system.templates.editor.content.preview.error'))(e)
    }
  } finally {
    previewLoading.value = null
  }
}

onMounted(() => loadDrivers())
</script>
