<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-2">
      <Checkbox v-model="preload" binary input-id="preload" />
      <label for="preload" class="text-sm">{{ $t('block.recordRevisions.configurator.preload') }}</label>
    </div>

    <CFormGroup
      :label="$t('block.recordRevisions.configurator.sortDirection.label')"
      :description="$t('block.recordRevisions.configurator.sortDirection.footnote')"
    >
      <Select
        v-model="sortDirection"
        :options="sortOptions"
        option-label="label"
        option-value="value"
        class="w-full"
      />
    </CFormGroup>

    <template v-if="selectedModule">
      <Divider />

      <Fieldset :legend="$t('block.recordRevisions.configurator.displayedFields')">
        <CFieldPicker
          :all-fields="allModuleFields"
          :model-value="selectedFieldNames"
          :available-label="$t('field.selector.available')"
          :selected-label="$t('field.selector.selected')"
          :select-all-label="$t('field.selector.selectAll')"
          :unselect-all-label="$t('field.selector.unselectAll')"
          :search-placeholder="$t('field.selector.search')"
          :no-items-label="$t('field.no-items-found')"
          @update:model-value="onFieldPickerUpdate"
        />
      </Fieldset>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/sections/compose/stores/module'

const { t } = useI18n()
const moduleStore = useModuleStore()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const sortOptions = [
  { value: 'desc', label: t('block.recordRevisions.configurator.sortDirection.desc') },
  { value: 'asc', label: t('block.recordRevisions.configurator.sortDirection.asc') },
]

const selectedModule = computed(() => {
  const moduleID = props.page?.moduleID
  if (!moduleID || moduleID === '0') return null
  return moduleStore.getByID(moduleID) || null
})

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

const preload = computed({
  get: () => block.value.options?.preload || false,
  set: v => updateOptions('preload', v),
})

const sortDirection = computed({
  get: () => (block.value.options?.sortDirection || 'desc').toLowerCase(),
  set: v => updateOptions('sortDirection', v),
})

// --- Field picker for displayedFields ---
const selectedFieldNames = ref([])

watch(() => block.value.options?.displayedFields, (fields) => {
  if (fields?.length) {
    selectedFieldNames.value = fields.map(f => (typeof f === 'string' ? f : f.name))
  }
}, { immediate: true })

const allModuleFields = computed(() => {
  if (!selectedModule.value) return []
  return selectedModule.value.fields || []
})

function onFieldPickerUpdate(names) {
  selectedFieldNames.value = names
  updateOptions('displayedFields', names)
}
</script>
