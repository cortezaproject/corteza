<template>
  <div class="flex flex-col gap-3">
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.recordOrganizer.module.placeholder') }}</label>
      <Select
        v-model="moduleID"
        :options="modules"
        option-label="name"
        option-value="moduleID"
        :placeholder="$t('block.recordOrganizer.module.placeholder')"
        class="w-full"
        filter
        show-clear
      />
    </div>

    <template v-if="moduleID">
      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.recordOrganizer.labelField.label') }}</label>
        <Select
          v-model="labelField"
          :options="fieldOptions"
          option-label="label"
          option-value="name"
          class="w-full"
          filter
          show-clear
        />
        <small class="text-muted-color">{{ $t('block.recordOrganizer.labelField.footnote') }}</small>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.recordOrganizer.descriptionField.label') }}</label>
        <Select
          v-model="descriptionField"
          :options="fieldOptions"
          option-label="label"
          option-value="name"
          class="w-full"
          filter
          show-clear
        />
        <small class="text-muted-color">{{ $t('block.recordOrganizer.descriptionField.footnote') }}</small>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.recordOrganizer.positionField.label') }}</label>
        <Select
          v-model="positionField"
          :options="numberFieldOptions"
          option-label="label"
          option-value="name"
          :placeholder="$t('block.recordOrganizer.positionField.placeholder')"
          class="w-full"
          filter
          show-clear
        />
        <small class="text-muted-color">{{ $t('block.recordOrganizer.positionField.footnote') }}</small>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.recordOrganizer.groupField.label') }}</label>
        <Select
          v-model="groupField"
          :options="fieldOptions"
          option-label="label"
          option-value="name"
          class="w-full"
          filter
          show-clear
        />
        <small class="text-muted-color">{{ $t('block.recordOrganizer.groupField.footnote') }}</small>
      </div>

      <div v-if="groupField" class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.recordOrganizer.group.label') }}</label>
        <InputText v-model="group" class="w-full" />
        <small class="text-muted-color">{{ $t('block.recordOrganizer.group.footnote') }}</small>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.recordOrganizer.displayOption.label') }}</label>
        <Select
          v-model="displayOption"
          :options="displayOptions"
          option-label="label"
          option-value="value"
          class="w-full"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.recordOrganizer.prefilter.label') }}</label>
        <Textarea
          v-model="prefilter"
          :placeholder="$t('block.recordOrganizer.prefilter.placeholder')"
          rows="3"
          class="w-full"
        />
        <small class="text-muted-color">{{ $t('block.recordOrganizer.prefilter.footnote') }}</small>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/stores/module'

const { t } = useI18n()
const moduleStore = useModuleStore()

defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const modules = computed(() => moduleStore.set || [])

const selectedModule = computed(() => {
  if (!moduleID.value) return null
  return moduleStore.getByID(moduleID.value) || null
})

const fieldOptions = computed(() => {
  if (!selectedModule.value) return []
  return (selectedModule.value.fields || []).map(f => ({
    name: f.name,
    label: f.label || f.name,
  }))
})

const numberFieldOptions = computed(() => {
  if (!selectedModule.value) return []
  return (selectedModule.value.fields || [])
    .filter(f => f.kind === 'Number')
    .map(f => ({ name: f.name, label: f.label || f.name }))
})

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

const moduleID = computed({
  get: () => block.value.options?.moduleID || '',
  set: v => updateOptions('moduleID', v),
})

const labelField = computed({
  get: () => block.value.options?.labelField || '',
  set: v => updateOptions('labelField', v),
})

const descriptionField = computed({
  get: () => block.value.options?.descriptionField || '',
  set: v => updateOptions('descriptionField', v),
})

const positionField = computed({
  get: () => block.value.options?.positionField || '',
  set: v => updateOptions('positionField', v),
})

const groupField = computed({
  get: () => block.value.options?.groupField || '',
  set: v => updateOptions('groupField', v),
})

const group = computed({
  get: () => block.value.options?.group || '',
  set: v => updateOptions('group', v),
})

const displayOptions = [
  { value: 'sameTab', label: t('block.record.openInSameTab') },
  { value: 'newTab', label: t('block.record.openInNewTab') },
  { value: 'modal', label: t('block.record.openInModal') },
]

const displayOption = computed({
  get: () => block.value.options?.displayOption || 'sameTab',
  set: v => updateOptions('displayOption', v),
})

const prefilter = computed({
  get: () => block.value.options?.filter || '',
  set: v => updateOptions('filter', v),
})
</script>
