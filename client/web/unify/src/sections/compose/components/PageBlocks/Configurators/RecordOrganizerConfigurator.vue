<template>
  <div class="flex flex-col gap-3">
    <CFormGroup :label="$t('block.general.module')" required>
      <CInputModule
        v-model="moduleID"
        :namespaceID="namespace?.namespaceID"
        :placeholder="$t('block.recordOrganizer.module.placeholder')"
      />
    </CFormGroup>

    <template v-if="moduleID">
      <Divider />

      <Fieldset :legend="$t('block.recordOrganizer.cardFields')">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <CFormGroup
            :label="$t('block.recordOrganizer.labelField.label')"
            :description="$t('block.recordOrganizer.labelField.footnote')"
          >
            <CInputModuleField v-model="labelField" :module="selectedModule" />
          </CFormGroup>

          <CFormGroup
            :label="$t('block.recordOrganizer.descriptionField.label')"
            :description="$t('block.recordOrganizer.descriptionField.footnote')"
          >
            <CInputModuleField v-model="descriptionField" :module="selectedModule" />
          </CFormGroup>

          <CFormGroup
            :label="$t('block.recordOrganizer.positionField.label')"
            :description="$t('block.recordOrganizer.positionField.footnote')"
          >
            <CInputModuleField
              v-model="positionField"
              :module="selectedModule"
              :kinds="['Number']"
              exclude-multi
              :placeholder="$t('block.recordOrganizer.positionField.placeholder')"
            />
          </CFormGroup>
        </div>
      </Fieldset>

      <Divider />

      <Fieldset :legend="$t('block.recordOrganizer.behaviour')">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <CFormGroup
            :label="$t('block.recordOrganizer.groupField.label')"
            :description="$t('block.recordOrganizer.groupField.footnote')"
          >
            <CInputModuleField v-model="groupField" :module="selectedModule" exclude-multi />
          </CFormGroup>

          <CFormGroup
            :label="$t('block.recordOrganizer.group.label')"
            :description="$t('block.recordOrganizer.group.footnote')"
          >
            <CFieldEditor
              v-if="groupFieldDef"
              :field="groupFieldDef"
              :namespace="namespace"
              :model-value="group"
              @update:model-value="group = $event"
            />
            <InputText v-else disabled class="w-full" />
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
            class="md:col-span-2"
          >
            <CInputExpression
              ref="prefilterInput"
              v-model="prefilter"
              dialect="ql"
              :scope="scope"
              :query-fields="queryFields"
              :placeholder="$t('block.recordOrganizer.prefilter.placeholder')"
            />
            <CExpressionHint :scope="scope" @insert="prefilterInput?.insert($event)" />
          </CFormGroup>
        </div>
      </Fieldset>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, useModuleStore } from '@planetcrust/human-vue'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const { CFieldEditor } = components

const { t } = useI18n()
const moduleStore = useModuleStore()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const selectedModule = computed(() => {
  if (!moduleID.value) return null
  return moduleStore.getByID(moduleID.value) || null
})

const prefilterInput = ref(null)
const { scope, queryFields } = useExpressionScope({
  page: computed(() => props.page),
  queryModule: computed(() => selectedModule.value),
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
  // The key value is one value of one field, and the block sends it as a
  // string; an editor that hands back an array contributes its first entry.
  set: v => updateOptions('group', (Array.isArray(v) ? v[0] : v) ?? ''),
})

// The key value is edited with the key field's own editor, so a Select offers
// its options and a Number its spinner instead of everything being typed as
// text. Single-value whatever the field says: one column holds one value, and
// the server refuses a multi-value key field outright.
const groupFieldDef = computed(() => {
  if (!groupField.value) return null

  const field = (selectedModule.value?.fields || []).find(f => f.name === groupField.value)
  if (!field) return null

  const selectType = field.options?.selectType

  return {
    ...field,
    isMulti: false,
    options: {
      ...field.options,
      selectType: selectType === 'multiple' ? 'default' : selectType,
    },
  }
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

// A key value belongs to the field it was entered for: a Select option code
// means nothing to a date picker, and left behind it would filter the column
// down to nothing. Not immediate — opening a saved block must not clear it.
watch(groupField, () => {
  if (block.value.options?.group) updateOptions('group', '')
})
</script>
