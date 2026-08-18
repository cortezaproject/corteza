<template>
  <MultiSelect
    v-if="multiple"
    :model-value="modelValue || []"
    :options="options"
    option-label="label"
    :option-value="valueKey"
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
    :option-value="valueKey"
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
  // What the v-model holds. Names address a field in block options; the
  // record block addresses one by `fieldID` instead.
  valueKey: {
    type: String,
    default: 'name',
  },
  // Extra predicate for a narrowing the kind/multi/queryable props cannot
  // express — a Record field pointing at one particular module, say.
  filter: {
    type: Function,
    default: null,
  },
  // Append the technical field name to each label — sorting and filtering are
  // by column, so two fields sharing a label still have to be told apart.
  showName: {
    type: Boolean,
    default: false,
  },
  // Synthetic entries offered above the module's own fields, as `{ name, label }`
  // — an aggregate like "count" is picked here but is not a field.
  extraOptions: {
    type: Array,
    default: () => [],
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

function decorate(field) {
  const label = field.isSystem ? t(`field.system.${field.name}`, labelOf(field)) : labelOf(field)
  return props.showName && label !== field.name ? `${label} (${field.name})` : label
}

const options = computed(() => {
  const extra = props.extraOptions.map(o => ({ ...o, isSystem: false, isExtra: true }))

  const mod = resolvedModule.value
  if (!mod) return extra

  const own = mod.fields || []
  // systemFields() is the one list of them; a hand-written copy drifts.
  const system = props.includeSystem
    ? (mod.systemFields?.() || []).map(f => ({ ...f, isSystem: true }))
    : []

  const fields = [...own, ...system]
    .filter(f => !props.kinds.length || props.kinds.includes(f.kind))
    .filter(f => !props.excludeMulti || !f.isMulti)
    .filter(f => !props.queryableOnly || f.isQueryable)
    .filter(f => !props.filter || props.filter(f))
    .map(f => ({
      name: f.name,
      fieldID: f.fieldID,
      kind: f.kind,
      isSystem: !!f.isSystem,
      label: decorate(f),
    }))
    .sort((a, b) => Number(a.isSystem) - Number(b.isSystem) || a.label.localeCompare(b.label))

  return [...extra, ...fields]
})
</script>
