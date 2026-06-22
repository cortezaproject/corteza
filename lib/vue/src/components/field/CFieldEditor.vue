<template>
  <!-- Multi-value fields where the editor absorbs the whole array (e.g. MultiSelect) -->
  <component
    v-if="multiAbsorbing"
    :is="editorComponent"
    :field="field"
    :namespace="namespace"
    :model-value="modelValue"
    :disabled="disabled"
    @update:model-value="$emit('update:modelValue', $event)"
  />

  <div v-else-if="field.isMulti" class="flex flex-col gap-2 w-full">
    <div v-for="entry in entries" :key="entry.id" class="flex items-center gap-2">
      <component
        :is="editorComponent"
        :field="field"
        :namespace="namespace"
        :model-value="entry.value"
        :disabled="disabled"
        class="flex-1"
        @update:model-value="updateValue(entry.id, $event)"
      />
      <Button
        icon="pi pi-trash"
        text
        rounded
        severity="danger"
        size="small"
        :disabled="disabled"
        @click="removeValue(entry.id)"
      />
    </div>
    <Button
      :label="addLabel"
      icon="pi pi-plus"
      severity="secondary"
      outlined
      size="small"
      :disabled="disabled"
      @click="addValue"
    />
  </div>

  <component
    v-else
    :is="editorComponent"
    :field="field"
    :namespace="namespace"
    :model-value="singleValue"
    :disabled="disabled"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { computed, provide, ref, watch } from 'vue'
import { resolveFieldEditor } from './registry'

// PrimeVue editable inputs (InputText, Textarea, InputNumber, …) auto-bind to
// the nearest @primevue/forms <FormField> via injected $pcFormField/$pcForm.
// An input without its own `name` adopts the FormField's name, so every input
// CFieldEditor renders for a multi-value field would share that field's single
// value — typing in one would change all of them. CFieldEditor manages its own
// model explicitly, so sever the injection for its sub-editors.
provide('$pcFormField', undefined)
provide('$pcForm', undefined)

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: [String, Array],
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  addLabel: {
    type: String,
    default: '',
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  allowEmpty: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const editorComponent = computed(() => resolveFieldEditor(props.field.kind))

// When the editor absorbs the full multi-value array itself (not per-entry)
const multiAbsorbing = computed(() => {
  // File always manages an array of attachment IDs regardless of isMulti
  if (props.field.kind === 'File') return true
  if (!props.field.isMulti) return false
  const selectType = props.field.options?.selectType
  if ((props.field.kind === 'Select' || props.field.kind === 'User' || props.field.kind === 'Record') && selectType === 'multiple') return true
  return false
})

const singleValue = computed(() => {
  if (Array.isArray(props.modelValue)) return props.modelValue[0] ?? ''
  return props.modelValue ?? ''
})

// Multi-value: maintain { id, value } pairs as single source of truth
let nextId = 0
const entries = ref([])

function valuesToEntries(vals) {
  return vals.map(v => ({ id: nextId++, value: v }))
}

function getValuesArray() {
  if (Array.isArray(props.modelValue)) return props.modelValue.length ? props.modelValue : (props.allowEmpty ? [] : [''])
  if (props.modelValue) return [props.modelValue]
  return props.allowEmpty ? [] : ['']
}

// Sync entries when modelValue changes externally
watch(
  () => props.modelValue,
  () => {
    const incoming = getValuesArray()
    const current = entries.value.map(e => e.value)
    // Only rebuild entries if the values actually differ to preserve stable keys
    if (
      incoming.length !== current.length ||
      incoming.some((v, i) => v !== current[i])
    ) {
      entries.value = valuesToEntries(incoming)
    }
  },
  { immediate: true },
)

function updateValue(id, value) {
  const index = entries.value.findIndex(e => e.id === id)
  if (index === -1) return
  entries.value[index].value = value
  emit('update:modelValue', entries.value.map(e => e.value))
}

function removeValue(id) {
  const index = entries.value.findIndex(e => e.id === id)
  if (index === -1) return
  entries.value.splice(index, 1)
  if (!entries.value.length && !props.allowEmpty) {
    entries.value = [{ id: nextId++, value: '' }]
  }
  emit('update:modelValue', entries.value.map(e => e.value))
}

function addValue() {
  entries.value.push({ id: nextId++, value: '' })
  emit('update:modelValue', entries.value.map(e => e.value))
}
</script>
