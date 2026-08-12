<template>
  <div>
    <span
      v-for="(rec, index) in resolvedRecords"
      :key="getRecordID(rec) || index"
      :class="[
        { block: isNewlineDelimiter, 'mt-1': isNewlineDelimiter && index !== 0 },
        canNavigate(rec) ? 'text-primary font-medium cursor-pointer hover:underline' : '',
      ]"
      @click="navigateToRecord(rec)"
    >
      <CFieldViewer
        v-if="rec.values && labelFieldDef"
        :field="labelFieldDef"
        :record="rec"
        :namespace="namespace"
        :disable-click="true"
        :value-only="true"
      />
      <template v-else>{{ getRecordID(rec) }}</template>
      {{ index !== resolvedRecords.length - 1 && !isNewlineDelimiter ? delimiter : '' }}
    </span>
  </div>
</template>

<script setup>
import { computed, inject, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRecordStore } from '../../../stores/useRecordStore'
import { useModuleStore } from '../../../stores/useModuleStore'
import { usePageStore } from '../../../stores/usePageStore'
import CFieldViewer from '../CFieldViewer.vue'

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
const recordStore = useRecordStore()
const moduleStore = useModuleStore()
const pageStore = usePageStore()
const $recordRoutes = inject('$recordRoutes', null)

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

const labelFieldDef = computed(() => {
  const moduleID = props.field.options?.moduleID
  if (!moduleID) return null

  const mod = moduleStore.getByID(moduleID)
  if (!mod?.fields?.length) return null

  const lf = props.field.options?.labelField
  if (lf) {
    const found = mod.fields.find(f => f.name === lf)
    if (found) return found
  }
  return mod.fields[0] || null
})

function getRecordID(rec) {
  return rec?.recordID || ''
}

const resolvedRecords = computed(() => {
  return recordIDs.value.map(id => recordStore.getByID(id) || { recordID: id })
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

  // Find the record page for this module via page store
  if (moduleID) {
    const pages = pageStore.set || []
    const page = pages.find(p => p.moduleID === moduleID)
    if (page) {
      const displayOption = props.extraOptions?.recordSelectorDisplayOption || 'sameTab'

      if (displayOption === 'modal') {
        router.push({
          query: {
            ...route.query,
            recordPageID: page.pageID,
            recordID: recordID,
          },
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
    if (!ids.length) return

    const moduleID = props.field.options?.moduleID
    const namespaceID = props.namespace?.namespaceID
    if (!moduleID || !namespaceID) return

    recordStore.resolveRecordLabels({ namespaceID, moduleID, recordIDs: ids })
  },
  { immediate: true },
)
</script>
