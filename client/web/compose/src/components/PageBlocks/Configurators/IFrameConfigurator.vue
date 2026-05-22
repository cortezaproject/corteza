<template>
  <div class="flex flex-col gap-3">
    <CFormGroup
      v-if="isRecordPage"
      :label="$t('block.iframe.srcFieldLabel')"
      :description="$t('block.iframe.srcFieldDesc')"
    >
      <Select
        v-model="srcField"
        :options="urlFields"
        option-label="text"
        option-value="value"
        :placeholder="$t('block.iframe.pickURLField')"
        class="w-full"
        show-clear
        :disabled="!urlFields.length"
      />
    </CFormGroup>

    <CFormGroup
      :label="$t('block.iframe.srcLabel')"
      :description="isRecordPage ? $t('block.iframe.srcDesc') : ''"
    >
      <InputText
        v-model="srcUrl"
        :placeholder="$t('block.content.urlPlaceholder')"
        class="w-full"
      />
    </CFormGroup>

    <!-- Interpolation footnote -->
    <small class="text-muted-color">
      {{ $t('block.content.interpolationFootnote') }}
      <code>${record.values.fieldName}</code>,
      <code>${recordID}</code>,
      <code>${ownerID}</code>,
      <code>${userID}</code>,
      <code>${user.name}</code>
    </small>


  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useModuleStore } from '@/stores/module'

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const moduleStore = useModuleStore()

const isRecordPage = computed(() => !!props.page?.moduleID && props.page.moduleID !== '0')

const urlFields = computed(() => {
  if (!props.page?.moduleID) return []
  const mod = moduleStore.getByID(props.page.moduleID)
  if (!mod) return []
  return mod.fields
    .filter(f => f.kind === 'Url')
    .map(({ label, name }) => ({ value: name, text: label || name }))
    .sort((a, b) => a.text.localeCompare(b.text))
})

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

const srcUrl = computed({
  get: () => block.value.options?.src || block.value.options?.url || '',
  set: v => updateOptions('src', v),
})

const srcField = computed({
  get: () => block.value.options?.srcField || '',
  set: v => updateOptions('srcField', v),
})


</script>
