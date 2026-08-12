<template>
  <div v-if="!rows.length" class="text-sm text-muted-color italic">—</div>
  <div v-else class="flex flex-col gap-0.5 w-max min-w-full">
    <TraceValueRow
      v-for="row in visibleRows"
      :key="row.key"
      :row="row"
      :search="searchLower"
      :depth="depth"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import TraceValueRow from './TraceValueRow.vue'

const props = defineProps({
  value: { type: null, default: null },
  search: { type: String, default: '' },
  depth: { type: Number, default: 0 },
})

const searchLower = computed(() => (props.search || '').trim().toLowerCase())

const rows = computed(() => toRows(props.value))

const visibleRows = computed(() => {
  if (!searchLower.value) return rows.value
  return rows.value.filter(r => rowMatches(r, searchLower.value))
})

function unwrapTyped(v) {
  if (v && typeof v === 'object' && !Array.isArray(v) && '@type' in v && '@value' in v) {
    return { type: String(v['@type']), value: v['@value'], wrapped: true }
  }
  return { type: null, value: v, wrapped: false }
}

function detectType(value) {
  if (value === null) return 'Null'
  if (Array.isArray(value)) return 'Array'
  switch (typeof value) {
    case 'string':
      return 'String'
    case 'number':
      return 'Number'
    case 'boolean':
      return 'Boolean'
    case 'undefined':
      return 'Null'
    case 'object':
      return 'Object'
    default:
      return 'Any'
  }
}

function toRows(container) {
  const { value } = unwrapTyped(container)
  if (value === null || value === undefined) return []
  if (Array.isArray(value)) {
    return value.map((v, i) => makeRow(`[${i}]`, v))
  }
  if (typeof value === 'object') {
    return Object.keys(value).map(k => makeRow(k, value[k]))
  }
  // Scalar at the top level — present as a single row
  return [makeRow('value', value)]
}

function makeRow(key, raw) {
  const unwrapped = unwrapTyped(raw)
  const actualValue = unwrapped.value
  const displayType = unwrapped.type || detectType(actualValue)
  const isArray = Array.isArray(actualValue)
  const isObject = !isArray && actualValue !== null && typeof actualValue === 'object'
  const isContainer = isArray || isObject
  const count = isArray ? actualValue.length : isObject ? Object.keys(actualValue).length : 0
  return {
    key,
    type: displayType,
    value: actualValue,
    isContainer,
    count,
    raw,
  }
}

function rowMatches(row, term) {
  if (!term) return true
  if (String(row.key).toLowerCase().includes(term)) return true
  if (!row.isContainer) {
    const s = row.value === null || row.value === undefined ? '' : String(row.value)
    return s.toLowerCase().includes(term)
  }
  // Container: check descendants
  const children = toRows(row.raw)
  return children.some(c => rowMatches(c, term))
}
</script>
