<template>
  <div class="flex flex-col gap-3">
    <CFormGroup :label="$t('block.recordOrganizer.module.placeholder')">
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
    </CFormGroup>

    <template v-if="moduleID">
      <CFormGroup
        :label="$t('block.recordOrganizer.labelField.label')"
        :description="$t('block.recordOrganizer.labelField.footnote')"
      >
        <Select
          v-model="labelField"
          :options="fieldOptions"
          option-label="label"
          option-value="name"
          class="w-full"
          filter
          show-clear
        />
      </CFormGroup>

      <CFormGroup
        :label="$t('block.recordOrganizer.descriptionField.label')"
        :description="$t('block.recordOrganizer.descriptionField.footnote')"
      >
        <Select
          v-model="descriptionField"
          :options="fieldOptions"
          option-label="label"
          option-value="name"
          class="w-full"
          filter
          show-clear
        />
      </CFormGroup>

      <CFormGroup
        :label="$t('block.recordOrganizer.positionField.label')"
        :description="$t('block.recordOrganizer.positionField.footnote')"
      >
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
      </CFormGroup>

      <CFormGroup
        :label="$t('block.recordOrganizer.groupField.label')"
        :description="$t('block.recordOrganizer.groupField.footnote')"
      >
        <Select
          v-model="groupField"
          :options="fieldOptions"
          option-label="label"
          option-value="name"
          class="w-full"
          filter
          show-clear
        />
      </CFormGroup>

      <CFormGroup
        v-if="groupField"
        :label="$t('block.recordOrganizer.group.label')"
        :description="$t('block.recordOrganizer.group.footnote')"
      >
        <InputText v-model="group" class="w-full" />
      </CFormGroup>

      <CFormGroup :label="$t('block.recordOrganizer.displayOption.label')">
        <Select
          v-model="displayOption"
          :options="displayOptions"
          option-label="label"
          option-value="value"
          class="w-full"
        />
      </CFormGroup>

      <CFormGroup
        :label="$t('block.recordOrganizer.prefilter.label')"
        :description="$t('block.recordOrganizer.prefilter.footnote')"
      >
        <Textarea
          v-model="prefilter"
          :placeholder="$t('block.recordOrganizer.prefilter.placeholder')"
          rows="3"
          class="w-full"
        />
      </CFormGroup>
    </template>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@planetcrust/human-vue'

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
