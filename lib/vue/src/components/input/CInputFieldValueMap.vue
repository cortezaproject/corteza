<template>
  <div class="flex flex-col gap-2">
    <div
      v-for="(row, index) in rows"
      :key="index"
      class="flex items-center gap-2"
    >
      <Select
        v-model="row.field"
        :options="getAvailableFields(index)"
        option-label="label"
        option-value="name"
        :placeholder="t('builder.fieldValueMap.selectField')"
        :disabled="disabled || !moduleID"
        class="flex-1"
        @update:model-value="onRowChange"
      />
      <InputText
        v-model="row.value"
        :placeholder="t('builder.fieldValueMap.enterValue')"
        :disabled="disabled || !row.field"
        class="flex-1"
        @update:model-value="onRowChange"
      />
      <Button
        icon="pi pi-trash"
        text
        rounded
        severity="danger"
        size="small"
        @click="removeRow(index)"
      />
    </div>

    <div>
      <Button
        :label="t('builder.fieldValueMap.addRow')"
        icon="pi pi-plus"
        severity="secondary"
        outlined
        size="small"
        :disabled="disabled || !moduleID"
        @click="addRow"
      />
    </div>
  </div>
</template>

<script setup>
import { inject, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  namespaceID: {
    type: [String, Number],
    default: null,
  },
  moduleID: {
    type: [String, Number],
    default: null,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])
const { t } = useI18n()
const $ComposeAPI = inject('$ComposeAPI')

const fields = ref([])
const rows = ref([])
const loading = ref(false)

let cancelCurrentRequest = null

function modelValueToRows(val) {
  if (!val || typeof val !== 'object') return []
  return Object.entries(val).map(([field, value]) => ({ field, value }))
}

function rowsToModelValue() {
  const result = {}
  for (const row of rows.value) {
    if (row.field) {
      result[row.field] = row.value ?? ''
    }
  }
  return result
}

function getAvailableFields(currentIndex) {
  const selectedFields = new Set(
    rows.value
      .filter((_, i) => i !== currentIndex)
      .map(r => r.field)
      .filter(Boolean),
  )
  return fields.value.filter(f => !selectedFields.has(f.name))
}

function addRow() {
  rows.value.push({ field: null, value: '' })
}

function removeRow(index) {
  rows.value.splice(index, 1)
  emitValue()
}

function onRowChange() {
  emitValue()
}

function emitValue() {
  emit('update:modelValue', rowsToModelValue())
}

async function fetchModuleFields() {
  if (!props.namespaceID || !props.moduleID || !$ComposeAPI) {
    fields.value = []
    return
  }

  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $ComposeAPI.moduleReadCancellable({
      namespaceID: props.namespaceID,
      moduleID: props.moduleID,
    })
    cancelCurrentRequest = cancel

    const mod = await response()
    fields.value = (mod.fields || [])
      .filter(f => !f.isSystem)
      .map(f => ({
        name: f.name,
        label: f.label || f.name,
        kind: f.kind,
      }))
  } catch (e) {
    if (e?.message !== 'canceled') {
      fields.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

// Initialize rows from modelValue
watch(
  () => props.modelValue,
  (val) => {
    const incoming = modelValueToRows(val)
    const current = rowsToModelValue()
    // Only update if externally changed to avoid cursor jumps
    if (JSON.stringify(val) !== JSON.stringify(current)) {
      rows.value = incoming.length ? incoming : []
    }
  },
  { immediate: true },
)

// Re-fetch fields when module changes; clear rows
watch(
  () => props.moduleID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      rows.value = []
      emitValue()
    }
    if (newVal && props.namespaceID) {
      fetchModuleFields()
    } else {
      fields.value = []
    }
  },
  { immediate: true },
)

// Re-fetch when namespace changes
watch(
  () => props.namespaceID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      fields.value = []
      rows.value = []
      emitValue()
    }
  },
)

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
