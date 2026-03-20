<template>
  <div class="flex flex-col gap-3">
    <!-- URL field from record (only on record pages) -->
    <div v-if="isRecordPage" class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.iframe.srcFieldLabel') }}</label>
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
      <small class="text-muted-color">{{ $t('block.iframe.srcFieldDesc') }}</small>
    </div>

    <!-- Static URL -->
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.iframe.srcLabel') }}</label>
      <InputText
        v-model="srcUrl"
        :placeholder="$t('block.content.urlPlaceholder')"
        class="w-full"
      />
      <small v-if="isRecordPage" class="text-muted-color">{{ $t('block.iframe.srcDesc') }}</small>
    </div>

    <!-- Interpolation footnote -->
    <small class="text-muted-color">
      {{ $t('block.content.interpolationFootnote') }}
      <code>${record.values.fieldName}</code>,
      <code>${recordID}</code>,
      <code>${ownerID}</code>,
      <code>${userID}</code>,
      <code>${user.name}</code>
    </small>

    <Divider />

    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.general.refreshRate') }}</label>
      <InputNumber
        v-model="refreshRate"
        :min="0"
        suffix=" s"
        class="w-full"
      />
    </div>

    <div class="flex items-center gap-2">
      <Checkbox v-model="showRefresh" binary input-id="showRefreshIframe" />
      <label for="showRefreshIframe" class="text-sm">{{ $t('block.general.showRefresh') }}</label>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useModuleStore } from '@/stores/module'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

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
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const srcUrl = computed({
  get: () => props.block.options?.src || props.block.options?.url || '',
  set: v => updateOptions('src', v),
})

const srcField = computed({
  get: () => props.block.options?.srcField || '',
  set: v => updateOptions('srcField', v),
})

const refreshRate = computed({
  get: () => props.block.options?.refreshRate ?? 0,
  set: v => updateOptions('refreshRate', v),
})

const showRefresh = computed({
  get: () => !!props.block.options?.showRefresh,
  set: v => updateOptions('showRefresh', v),
})
</script>
