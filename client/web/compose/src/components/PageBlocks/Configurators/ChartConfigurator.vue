<template>
  <div class="flex flex-col gap-3">
    <!-- Chart selection -->
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.chart.display') }}</label>
      <div class="flex gap-2">
        <CInputChart
          v-model="chartID"
          :namespaceID="namespace?.namespaceID"
          :placeholder="$t('block.chart.pick')"
          class="flex-1"
        />
        <Button
          v-if="chartID"
          icon="pi pi-external-link"
          severity="secondary"
          :title="$t('block.chart.openInBuilder')"
          @click="goToChart"
        />
      </div>
    </div>

    <template v-if="chartID">
      <!-- Live filter -->
      <CInputSwitch v-model="liveFilterEnabled" :label="$t('block.chart.enableLiveFilter')" />

      <Divider v-if="isDrillDownAvailable" />

      <div v-if="isDrillDownAvailable" class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.chart.drillDown.label') }}</h5>

        <CInputSwitch v-model="drillDownEnabled" :label="$t('block.chart.drillDown.enabled')" />

        <template v-if="drillDownEnabled">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.chart.drillDown.blockID') }}</label>
            <Select
              v-model="drillDownBlockID"
              :options="drillDownOptions"
              option-label="label"
              option-value="value"
              :placeholder="$t('block.chart.drillDown.openInModal', 'Open in modal')"
              class="w-full"
              show-clear
            />
            <small class="text-muted-color">{{ $t('block.chart.drillDown.blockIDFootnote') }}</small>
          </div>

          <div v-if="!drillDownBlockID" class="flex flex-col gap-1 mt-2">
            <label class="text-primary font-medium text-sm">{{ $t('block.chart.drillDown.fields', 'Fields') }}</label>
            <CFieldPicker
              :all-fields="allModuleFields"
              :model-value="drillDownFields"
              :available-label="$t('field.selector.available')"
              :selected-label="$t('field.selector.selected')"
              :select-all-label="$t('field.selector.selectAll')"
              :unselect-all-label="$t('field.selector.unselectAll')"
              :search-placeholder="$t('field.selector.search')"
              :no-items-label="$t('field.no-items-found')"
              @update:model-value="onFieldPickerUpdate"
            />
          </div>
        </template>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import CInputChart from '@cortezaproject/corteza-vue-next/src/components/input/CInputChart.vue'
import { useModuleStore } from '@/stores/module'
import { useComposeResourceStore } from '@cortezaproject/corteza-vue-next'

const router = useRouter()
const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const moduleStore = useModuleStore()
const composeResourceStore = useComposeResourceStore()

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const chartID = computed({
  get: () => props.block.options?.chartID,
  set: v => updateOptions('chartID', v),
})

const liveFilterEnabled = computed({
  get: () => !!props.block.options?.liveFilterEnabled,
  set: v => updateOptions('liveFilterEnabled', v),
})

const drillDownEnabled = computed({
  get: () => !!props.block.options?.drillDown?.enabled,
  set: v => updateOptions('drillDown', { ...props.block.options?.drillDown, enabled: v }),
})

const drillDownBlockID = computed({
  get: () => props.block.options?.drillDown?.blockID || '',
  set: v => updateOptions('drillDown', { ...props.block.options?.drillDown, blockID: v }),
})

const drillDownFields = computed({
  get: () => {
    const fields = props.block.options?.drillDown?.recordListOptions?.fields || []
    return fields.map(f => f.name ?? f)
  },
  set: v => {
    const drillDown = props.block.options?.drillDown || {}
    const recordListOptions = drillDown.recordListOptions || {}
    updateOptions('drillDown', {
      ...drillDown,
      recordListOptions: { ...recordListOptions, fields: v }
    })
  }
})

function onFieldPickerUpdate(names) {
  drillDownFields.value = names
}

const selectedChart = ref(null)

watch(() => props.block.options?.chartID, async (id) => {
  if (id && props.namespace?.namespaceID) {
    try {
      selectedChart.value = await composeResourceStore.resolveChart(props.namespace.namespaceID, id)
    } catch {
      selectedChart.value = null
    }
  } else {
    selectedChart.value = null
  }
}, { immediate: true })

const selectedChartModuleID = computed(() => {
  return selectedChart.value?.config?.reports?.[0]?.moduleID
})

const selectedChartModule = computed(() => {
  if (!selectedChartModuleID.value) return null
  return moduleStore.getByID(selectedChartModuleID.value) || null
})

const isDrillDownAvailable = computed(() => {
  if (!selectedChart.value) return false
  const metrics = selectedChart.value?.config?.reports?.[0]?.metrics || []
  return !metrics.some(({ type }) => type === 'gauge' || type === 'radar')
})

const drillDownOptions = computed(() => {
  if (!props.page?.blocks) return []
  return props.page.blocks
    .filter(b => b.kind === 'RecordList' && b.blockID && b.options?.moduleID === selectedChartModuleID.value)
    .map(b => ({
      label: b.title || b.kind,
      value: b.blockID
    }))
})

const allModuleFields = computed(() => {
  if (!selectedChartModule.value) return []
  const regular = selectedChartModule.value.fields || []
  const system = (selectedChartModule.value.systemFields?.() || []).map(f => ({
    ...f,
    label: t(`field.system.${f.name}`, f.label || f.name),
    isSystem: true,
  }))
  return [...regular, ...system]
})

function goToChart() {
  if (!chartID.value) return
  router.push({
    name: 'admin.charts.edit',
    params: { chartID: chartID.value },
  })
}
</script>
