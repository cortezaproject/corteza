<template>
  <div class="flex flex-col gap-3">
    <CFormGroup :label="$t('block.chart.display')" required>
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
    </CFormGroup>

    <template v-if="chartID && isDrillDownAvailable">
      <Divider />

      <Fieldset :legend="$t('block.chart.drillDown.label')">
        <div class="flex flex-col gap-3">
          <CInputToggleCard
            v-model="drillDownEnabled"
            :label="$t('block.chart.drillDown.enabled')"
            :description="$t('block.chart.drillDown.description')"
          />

          <template v-if="drillDownEnabled">
            <CFormGroup
              :label="$t('block.chart.drillDown.blockID')"
              :description="$t('block.chart.drillDown.blockIDFootnote')"
            >
              <Select
                v-model="drillDownBlockID"
                :options="drillDownOptions"
                option-label="label"
                option-value="value"
                :placeholder="$t('block.chart.drillDown.openInModal')"
                class="w-full"
                show-clear
              />
            </CFormGroup>

            <CFormGroup v-if="!drillDownBlockID" :label="$t('block.chart.drillDown.fields')">
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
            </CFormGroup>
          </template>
        </div>
      </Fieldset>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import CInputChart from '@planetcrust/human-vue/src/components/input/CInputChart.vue'
import { useModuleStore } from '@planetcrust/human-vue'
import { useChartStore } from '@planetcrust/human-vue'

const router = useRouter()
const { t } = useI18n()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')
const $ComposeAPI = inject('$ComposeAPI')

const moduleStore = useModuleStore()
const chartStore = useChartStore()

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

const chartID = computed({
  get: () => block.value.options?.chartID,
  set: v => updateOptions('chartID', v),
})

const drillDownEnabled = computed({
  get: () => !!block.value.options?.drillDown?.enabled,
  set: v => updateOptions('drillDown', { ...block.value.options?.drillDown, enabled: v }),
})

const drillDownBlockID = computed({
  get: () => block.value.options?.drillDown?.blockID || '',
  set: v => updateOptions('drillDown', { ...block.value.options?.drillDown, blockID: v }),
})

const drillDownFields = computed({
  get: () => {
    const fields = block.value.options?.drillDown?.recordListOptions?.fields || []
    return fields.map(f => f.name ?? f)
  },
  set: v => {
    const drillDown = block.value.options?.drillDown || {}
    const recordListOptions = drillDown.recordListOptions || {}
    updateOptions('drillDown', {
      ...drillDown,
      recordListOptions: { ...recordListOptions, fields: v },
    })
  },
})

function onFieldPickerUpdate(names) {
  drillDownFields.value = names
}

const selectedChart = ref(null)

watch(
  () => block.value.options?.chartID,
  async id => {
    if (id && props.namespace?.namespaceID) {
      try {
        selectedChart.value = await chartStore.findByID({
          namespaceID: props.namespace.namespaceID,
          chartID: id,
        })
      } catch {
        selectedChart.value = null
      }
    } else {
      selectedChart.value = null
    }
  },
  { immediate: true },
)

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
    .filter(
      b =>
        b.kind === 'RecordList' && b.blockID && b.options?.moduleID === selectedChartModuleID.value,
    )
    .map(b => ({
      label: b.title || b.kind,
      value: b.blockID,
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
