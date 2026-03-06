<template>
  <div class="flex flex-col gap-5">
    <!-- General -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">
        {{ $t('block.recordList.record.generalLabel') }}
      </h5>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- Module (read-only since it comes from the page) -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.general.module') }}</label>
          <InputText :model-value="moduleName" disabled class="w-full" />
        </div>

        <!-- Inline record edit -->
        <CInputSwitch v-model="inlineEditEnabled" :label="$t('block.record.inlineEdit.enabled')" />

        <!-- Horizontal form layout -->
        <CInputSwitch v-model="horizontalLayout" :label="$t('block.record.horizontalFormLayout')" :disabled="layoutMode === 'noWrap'" />

        <!-- Field layout mode -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.record.fieldsLayoutMode.label') }}</label>
          <Select
            v-model="layoutMode"
            :options="fieldLayoutOptions"
            option-label="text"
            option-value="value"
            class="w-full"
          />
        </div>
      </div>
    </div>

    <template v-if="selectedModule">
      <Divider />

      <!-- Fields -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.general.fields') }}
        </h5>

        <PickList
          v-model="fieldPickerModel"
          data-key="name"
          breakpoint="768px"
          :pt="{
            list: { style: 'height: 300px' },
          }"
        >
          <template #option="{ option }">
            {{ option.label || option.name }}
          </template>
        </PickList>
      </div>

      <Divider />

      <!-- Field conditions -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.record.fieldConditions.label') }}
        </h5>

        <div class="flex flex-col gap-2">
          <div
            v-for="(condition, i) in fieldConditions"
            :key="i"
            class="flex items-center gap-2"
          >
            <Select
              v-model="condition.field"
              :options="conditionFieldOptions"
              option-label="text"
              option-value="value"
              :placeholder="$t('block.record.fieldConditions.selectPlaceholder')"
              class="flex-1"
              filter
            />
            <InputText
              v-model="condition.condition"
              :placeholder="$t('block.record.fieldConditions.placeholder')"
              class="flex-1"
            />
            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              @click="removeCondition(i)"
            />
          </div>
        </div>

        <Button
          :label="$t('general.label.add')"
          icon="pi pi-plus"
          severity="secondary"
          size="small"
          class="self-start"
          @click="addCondition"
        />

        <small class="text-muted-color">
          {{ $t('block.record.fieldConditions.description') }}
        </small>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/stores/module'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const moduleStore = useModuleStore()

const selectedModule = computed(() => {
  const moduleID = props.block.options?.referenceModuleID || props.page?.moduleID
  if (!moduleID) return null
  return moduleStore.getByID(moduleID) || null
})

const moduleName = computed(() => selectedModule.value?.name || t('block.record.noModule'))

const fieldLayoutOptions = [
  { value: 'default', text: t('block.record.fieldsLayoutMode.default') },
  { value: 'noWrap', text: t('block.record.fieldsLayoutMode.noWrap') },
  { value: 'wrap', text: t('block.record.fieldsLayoutMode.wrap') },
]

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const inlineEditEnabled = computed({
  get: () => !!props.block.options?.inlineRecordEditEnabled,
  set: v => updateOptions('inlineRecordEditEnabled', v),
})

const horizontalLayout = computed({
  get: () => !!props.block.options?.horizontalFieldLayoutEnabled,
  set: v => updateOptions('horizontalFieldLayoutEnabled', v),
})

const layoutMode = computed({
  get: () => props.block.options?.recordFieldLayoutOption || 'default',
  set: v => {
    updateOptions('recordFieldLayoutOption', v)
    if (v === 'noWrap') {
      updateOptions('horizontalFieldLayoutEnabled', false)
    }
  },
})

// --- Field picker ---

const selectedFieldNames = ref([])

watch(() => props.block.options?.fields, (fields) => {
  if (fields?.length) {
    selectedFieldNames.value = fields.map(f => f.name ?? f)
  }
}, { immediate: true })

const availableFields = computed(() => {
  if (!selectedModule.value) return []
  const selected = new Set(selectedFieldNames.value)
  return (selectedModule.value.fields || []).filter(f => !selected.has(f.name))
})

const selectedFields = computed(() => {
  if (!selectedModule.value) return []
  return selectedFieldNames.value
    .map(name => (selectedModule.value.fields || []).find(f => f.name === name))
    .filter(Boolean)
})

const fieldPickerModel = computed({
  get: () => [availableFields.value, selectedFields.value],
  set: (val) => {
    const [, selected] = val
    selectedFieldNames.value = selected.map(f => f.name)
    updateOptions('fields', selected.map(f => f.name))
  },
})

// --- Field conditions ---

const fieldConditions = computed(() => {
  return props.block.options?.fieldConditions || []
})

const conditionFieldOptions = computed(() => {
  if (!selectedModule.value) return []
  return (selectedModule.value.fields || []).map(f => ({
    value: f.fieldID || f.name,
    text: f.label || f.name,
  }))
})

function addCondition() {
  const conditions = [...fieldConditions.value, { field: undefined, condition: '', clearOnHide: false }]
  updateOptions('fieldConditions', conditions)
}

function removeCondition(index) {
  const conditions = [...fieldConditions.value]
  conditions.splice(index, 1)
  updateOptions('fieldConditions', conditions)
}
</script>
