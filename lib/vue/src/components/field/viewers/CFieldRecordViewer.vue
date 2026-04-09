<template>
  <div>
    <span
      v-for="(rec, index) in resolvedRecords"
      :key="getRecordID(rec) || index"
      :class="[
        { block: isNewlineDelimiter, 'mt-1': isNewlineDelimiter && index !== 0 },
        canNavigate(rec) ? 'record-link' : '',
      ]"
      @click="navigateToRecord(rec)"
    >
      {{ getRecordLabel(rec) }}{{ index !== resolvedRecords.length - 1 && !isNewlineDelimiter ? delimiter : '' }}
    </span>
  </div>
</template>

<script setup>
import { computed, inject, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  record: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  valueOnly: {
    type: Boolean,
    default: false,
  },
  extraOptions: {
    type: Object,
    default: () => ({}),
  },
  disableClick: {
    type: Boolean,
    default: false,
  },
})

const router = useRouter()
const route = useRoute()
const $recordStore = inject('$recordStore', null)
const $recordRoutes = inject('$recordRoutes', null)
const $pageStore = inject('$pageStore', null)

const recordIDs = computed(() => {
  const v = props.field.isSystem
    ? props.record[props.field.name]
    : props.record?.values?.[props.field.name]

  if (!v) return []
  if (props.field.isMulti && Array.isArray(v)) return v.filter(Boolean)
  return [v].filter(Boolean)
})

const delimiter = computed(() => props.field.options?.multiDelimiter || ', ')
const isNewlineDelimiter = computed(() => delimiter.value === '\n')
const labelField = computed(() => props.field.options?.labelField || '')

function getRecordID(rec) {
  return rec?.recordID || ''
}

function getRecordLabel(rec) {
  if (!rec) return ''

  const lf = labelField.value

  // compose.Record stores values as an object { fieldName: value }
  if (rec.values && typeof rec.values === 'object' && !Array.isArray(rec.values)) {
    if (lf && rec.values[lf]) return rec.values[lf]
    const firstVal = Object.values(rec.values).find(v => v)
    if (firstVal) return firstVal
  }

  // Raw API record stores values as array [{ name, value }]
  if (Array.isArray(rec.values)) {
    if (lf) {
      const entry = rec.values.find(v => v.name === lf)
      if (entry?.value) return entry.value
    }
    const first = rec.values.find(v => v.value)
    if (first?.value) return first.value
  }

  return rec.recordID || ''
}

const resolvedRecords = computed(() => {
  return recordIDs.value.map(id => {
    if ($recordStore) {
      return $recordStore.getByID(id) || { recordID: id }
    }
    return { recordID: id }
  })
})

function canNavigate(rec) {
  if (props.disableClick) return false
  return !!getRecordID(rec)
}

function navigateToRecord(rec) {
  if (!canNavigate(rec)) return

  const recordID = getRecordID(rec)
  const moduleID = props.field.options?.moduleID

  // Use injected routes if available (admin context)
  if ($recordRoutes && moduleID) {
    router.push($recordRoutes.view(moduleID, recordID))
    return
  }

  // Find the record page for this module via injected page store
  if ($pageStore && moduleID) {
    const pages = $pageStore.set || []
    const page = pages.find(p => p.moduleID === moduleID)
    if (page) {
      const displayOption = props.extraOptions?.recordSelectorDisplayOption || 'sameTab'

      if (displayOption === 'modal') {
        router.push({
          query: {
            ...route.query,
            recordPageID: page.pageID,
            recordID: recordID,
          }
        })
        return
      }

      const routeObj = { name: 'page.record', params: { pageID: page.pageID, recordID } }

      if (displayOption === 'newTab') {
        window.open(router.resolve(routeObj).href)
      } else {
        router.push(routeObj)
      }
      return
    }
  }
}

// Resolve records not yet in the store
watch(
  () => [recordIDs.value, props.field.options?.moduleID, props.namespace?.namespaceID],
  ([ids]) => {
    if (!ids.length || !$recordStore) return

    const moduleID = props.field.options?.moduleID
    const namespaceID = props.namespace?.namespaceID
    if (!moduleID || !namespaceID) return

    $recordStore.resolveRecordLabels({ namespaceID, moduleID, recordIDs: ids })
  },
  { immediate: true },
)
</script>

<style scoped>
.record-link {
  color: var(--p-primary-color);
  cursor: pointer;
}

.record-link:hover {
  text-decoration: underline;
}
</style>

