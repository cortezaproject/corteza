<template>
  <div class="flex flex-col gap-3">
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.report.label') }}</label>
      <InputText
        v-model="reportID"
        :placeholder="$t('block.report.reportIDPlaceholder')"
        class="w-full"
      />
      <small class="text-muted-color">{{ $t('block.report.reportIDFootnote') }}</small>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.report.scenario.label') }}</label>
      <InputText
        v-model="scenarioID"
        :placeholder="$t('block.report.scenarioIDPlaceholder')"
        class="w-full"
      />
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.report.element.label') }}</label>
      <InputText
        v-model="elementID"
        :placeholder="$t('block.report.elementIDPlaceholder')"
        class="w-full"
      />
    </div>

    <Divider />

    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.general.refreshRate') }}</label>
      <InputNumber
        v-model="refreshRate"
        :min="0"
        suffix=" s"
        class="w-full"
      />
      <small class="text-muted-color">{{ $t('block.general.refreshRateFootnote') }}</small>
    </div>

    <div class="flex items-center gap-2">
      <Checkbox v-model="showRefresh" binary input-id="showRefresh" />
      <label for="showRefresh" class="text-sm">{{ $t('block.general.showRefresh') }}</label>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const reportID = computed({
  get: () => props.block.options?.reportID || '',
  set: v => updateOptions('reportID', v),
})

const elementID = computed({
  get: () => props.block.options?.elementID || '',
  set: v => updateOptions('elementID', v),
})

const scenarioID = computed({
  get: () => props.block.options?.scenarioID || '',
  set: v => updateOptions('scenarioID', v),
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
