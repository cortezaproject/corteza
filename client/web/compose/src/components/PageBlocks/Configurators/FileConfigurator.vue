<template>
  <div class="flex flex-col gap-3">
    <!-- View mode -->
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.file.view.modeLabel') }}</label>
      <SelectButton
        v-model="mode"
        :options="modes"
        option-label="text"
        option-value="value"
      />
      <small class="text-muted-color">{{ $t('block.file.view.modeFootnote') }}</small>
    </div>

    <!-- Display options -->
    <div class="flex flex-col gap-2">
      <div v-if="mode === 'gallery'" class="flex items-center gap-2">
        <Checkbox v-model="hideFileName" :binary="true" input-id="hideFileName" />
        <label for="hideFileName" class="text-sm">{{ $t('block.file.view.showName') }}</label>
      </div>

      <div class="flex items-center gap-2">
        <Checkbox v-model="clickToView" :binary="true" input-id="clickToView" />
        <label for="clickToView" class="text-sm">{{ $t('block.file.view.clickToView') }}</label>
      </div>

      <div class="flex items-center gap-2">
        <Checkbox v-model="enableDownload" :binary="true" input-id="enableDownload" />
        <label for="enableDownload" class="text-sm">{{ $t('block.file.view.enableDownload') }}</label>
      </div>
    </div>

    <!-- Gallery preview styling -->
    <template v-if="mode === 'gallery'">
      <Divider />

      <h5 class="text-lg font-semibold text-primary m-0">
        {{ $t('block.file.view.previewStyle') }}
      </h5>
      <small class="text-muted-color">{{ $t('block.file.view.description') }}</small>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.height') }}</label>
          <InputText v-model="height" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.width') }}</label>
          <InputText v-model="width" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.maxHeight') }}</label>
          <InputText v-model="maxHeight" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.maxWidth') }}</label>
          <InputText v-model="maxWidth" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.borderRadius') }}</label>
          <InputText v-model="borderRadius" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.margin') }}</label>
          <InputText v-model="margin" class="w-full" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const modes = [
  { value: 'list', text: t('block.file.view.list') },
  { value: 'gallery', text: t('block.file.view.gallery') },
]

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const mode = computed({
  get: () => props.block.options?.mode || 'list',
  set: v => updateOptions('mode', v),
})

const hideFileName = computed({
  get: () => !!props.block.options?.hideFileName,
  set: v => updateOptions('hideFileName', v),
})

const clickToView = computed({
  get: () => !!props.block.options?.clickToView,
  set: v => updateOptions('clickToView', v),
})

const enableDownload = computed({
  get: () => props.block.options?.enableDownload !== false,
  set: v => updateOptions('enableDownload', v),
})

const height = computed({
  get: () => props.block.options?.height || '',
  set: v => updateOptions('height', v),
})

const width = computed({
  get: () => props.block.options?.width || '',
  set: v => updateOptions('width', v),
})

const maxHeight = computed({
  get: () => props.block.options?.maxHeight || '',
  set: v => updateOptions('maxHeight', v),
})

const maxWidth = computed({
  get: () => props.block.options?.maxWidth || '',
  set: v => updateOptions('maxWidth', v),
})

const borderRadius = computed({
  get: () => props.block.options?.borderRadius || '',
  set: v => updateOptions('borderRadius', v),
})

const margin = computed({
  get: () => props.block.options?.margin || '',
  set: v => updateOptions('margin', v),
})
</script>
