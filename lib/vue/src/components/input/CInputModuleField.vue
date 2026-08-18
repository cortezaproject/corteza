<template>
  <MultiSelect
    v-if="multiple"
    :model-value="modelValue || []"
    :options="options"
    option-label="label"
    option-value="name"
    :placeholder="placeholder"
    :disabled="disabled || !options.length"
    :show-clear="showClear"
    class="w-full"
    display="chip"
    filter
    :filter-fields="['label', 'name']"
    fluid
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #option="{ option }">
      <span :class="option.isSystem ? 'italic' : ''">{{ option.label }}</span>
    </template>
  </MultiSelect>

  <Select
    v-else
    :model-value="modelValue || null"
    :options="options"
    option-label="label"
    option-value="name"
    :placeholder="placeholder"
    :disabled="disabled || !options.length"
    :show-clear="showClear"
    class="w-full"
    filter
    :filter-fields="['label', 'name']"
    fluid
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #option="{ option }">
      <span :class="option.isSystem ? 'italic' : ''">{{ option.label }}</span>
    </template>
  </Select>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '../../stores/useModuleStore'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  // Field name, or names when `multiple`.
  modelValue: {
    type: [String, Array],
    default: null,
  },
  // Either the module itself or its ID; the ID is resolved against the store.
  module: {
    type: Object,
    default: null,
  },
  moduleID: {
    type: [String, Number],
    default: '',
  },
  // Field kinds to offer. Empty offers every kind.
  kinds: {
    type: Array,
    default: () => [],
  },
  // Append the record's system fields (createdAt, ownedBy, …).
  includeSystem: {
    type: Boolean,
    default: false,
  },
  // Drop multi-value fields — sort keys and single-value mappings cannot use them.
  excludeMulti: {
    type: Boolean,
    default: false,
  },
  // Keep only fields the record store can filter on.
  queryableOnly: {
    type: Boolean,
    default: false,
  },
  multiple: {
    type: Boolean,
    default: false,
  },
  placeholder: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  showClear: {
    type: Boolean,
    default: true,
  },
})

defineEmits(['update:modelValue'])

const { t } = useI18n()
const moduleStore = useModuleStore()

const resolvedModule = computed(
  () => props.module || (props.moduleID ? moduleStore.getByID(props.moduleID) : null),
)

// A field's label is authored and can be blank on an imported module, so the
// name is what stands in — a picker row with no text is unpickable.
function labelOf(field) {
  return field.label || field.name
}

const options = computed(() => {
  const mod = resolvedModule.value
  if (!mod) return []

  const own = mod.fields || []
  // systemFields() is the one list of them; a hand-written copy drifts.
  const system = props.includeSystem
    ? (mod.systemFields?.() || []).map(f => ({ ...f, isSystem: true }))
    : []

  return [...own, ...system]
    .filter(f => !props.kinds.length || props.kinds.includes(f.kind))
    .filter(f => !props.excludeMulti || !f.isMulti)
    .filter(f => !props.queryableOnly || f.isQueryable)
    .map(f => ({
      name: f.name,
      kind: f.kind,
      isSystem: !!f.isSystem,
      label: f.isSystem ? t(`field.system.${f.name}`, labelOf(f)) : labelOf(f),
    }))
    .sort((a, b) => Number(a.isSystem) - Number(b.isSystem) || a.label.localeCompare(b.label))
})
</script>
