<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-2">
      <Checkbox v-model="preload" binary input-id="preload" />
      <label for="preload" class="text-sm">{{ $t('block.recordRevisions.configurator.preload') }}</label>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.recordRevisions.configurator.sortDirection.label') }}</label>
      <Select
        v-model="sortDirection"
        :options="sortOptions"
        option-label="label"
        option-value="value"
        class="w-full"
      />
      <small class="text-muted-color">
        {{ $t('block.recordRevisions.configurator.sortDirection.footnote') }}
      </small>
    </div>

    <template v-if="selectedModule">
      <Divider />

      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordRevisions.configurator.displayedFields') }}
        </h5>

        <PickList
          v-model="fieldPickerModel"
          data-key="name"
          breakpoint="768px"
          :pt="{
            list: { style: 'height: 250px' },
          }"
        >
          <template #option="{ option }">
            {{ option.label || option.name }}
          </template>
        </PickList>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/stores/module'

const { t } = useI18n()
const moduleStore = useModuleStore()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

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
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const preload = computed({
  get: () => props.block.options?.preload || false,
  set: v => updateOptions('preload', v),
})

const sortDirection = computed({
  get: () => (props.block.options?.sortDirection || 'desc').toLowerCase(),
  set: v => updateOptions('sortDirection', v),
})

// --- Field picker for displayedFields ---
const selectedFieldNames = ref([])

watch(() => props.block.options?.displayedFields, (fields) => {
  if (fields?.length) {
    selectedFieldNames.value = fields.map(f => (typeof f === 'string' ? f : f.name))
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
    updateOptions('displayedFields', selected.map(f => f.name))
  },
})
</script>
