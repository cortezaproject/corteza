<template>
  <PageBlock :block="block" @refreshBlock="pullRecords">
    <div v-if="!isConfigured" class="flex items-center justify-center h-full p-3 text-muted-color italic">
      {{ $t('block.recordOrganizer.notConfigured') }}
    </div>

    <div v-else-if="loading" class="flex items-center justify-center h-full">
      <ProgressSpinner style="width: 28px; height: 28px" />
    </div>

    <div v-else class="h-full flex flex-col">
      <!-- Add record button -->
      <div v-if="canAddRecord" class="p-3 border-b border-surface">
        <Button
          :label="$t('block.recordOrganizer.addNewRecord')"
          severity="primary"
          size="small"
          @click="createNewRecord"
        />
      </div>

      <!-- Records -->
      <div class="flex-1 overflow-auto p-3">
        <div v-if="!records.length" class="text-muted-color text-sm">
          {{ $t('block.recordOrganizer.noRecords') }}
        </div>

        <div
          v-for="record in records"
          :key="record.recordID"
          class="record-card p-3 mb-3 border border-surface rounded-border cursor-pointer hover:bg-emphasis transition-colors"
          @click="handleRecordClick(record)"
        >
          <h6 v-if="labelFieldName" class="font-medium mb-1 text-color">
            {{ getFieldValue(record, labelFieldName) || $t('block.record.preview.untitled') }}
          </h6>
          <p v-if="descriptionFieldName" class="text-sm text-muted-color mb-0">
            {{ getFieldValue(record, descriptionFieldName) }}
          </p>
        </div>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, inject } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import PageBlock from './PageBlock.vue'
import { evaluatePrefilter } from '../../../lib/record-filter'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI', null)
const $auth = inject('$auth', {})
const $eventBus = inject('$eventBus', null)
const router = useRouter()
const route = useRoute()

const loading = ref(false)
const records = ref([])

const options = computed(() => props.block.options || {})
const labelFieldName = computed(() => options.value.labelField || '')
const descriptionFieldName = computed(() => options.value.descriptionField || '')
const isConfigured = computed(() => !!options.value.moduleID)
const canAddRecord = computed(() => !!options.value.moduleID)

function getFieldValue(record, fieldName) {
  if (!record || !fieldName) return ''
  const val = record.values?.[fieldName]
  if (Array.isArray(val)) return val[0] || ''
  return val || ''
}

async function pullRecords() {
  if (!$ComposeAPI || !options.value.moduleID) return

  loading.value = true

  try {
    const { namespaceID } = props.namespace
    const { moduleID, positionField, filter: prefilter, groupField, group } = options.value

    const filterParts = []

    if (prefilter) {
      const record = props.record
      const user = $auth?.user || {}
      filterParts.push(`(${evaluatePrefilter(prefilter, {
        record, user,
        recordID: record?.recordID || '0',
        ownerID: record?.ownedBy || '0',
        userID: user?.userID || '0',
      })})`)
    }

    if (groupField && group !== undefined) {
      filterParts.push(`(${groupField} = '${group}')`)
    }

    const query = filterParts.join(' AND ')
    const sort = positionField || 'updatedAt'

    const { set = [] } = await $ComposeAPI.recordList({ namespaceID, moduleID, query, sort })
    records.value = set
  } catch (e) {
    console.error('Failed to load records for organizer:', e)
    records.value = []
  } finally {
    loading.value = false
  }
}

function handleRecordClick(record) {
  const { displayOption } = options.value

  if (displayOption === 'modal') {
    router.push({
      query: {
        ...route.query,
        recordPageID: props.page.pageID,
        recordID: record.recordID,
      }
    })
    return
  }

  const recordRoute = {
    name: 'page.record',
    params: { pageID: props.page.pageID, recordID: record.recordID },
  }

  if (displayOption === 'newTab') {
    window.open(router.resolve(recordRoute).href)
  } else {
    router.push(recordRoute)
  }
}

function createNewRecord() {
  const { addRecordDisplayOption, displayOption } = options.value
  const displayMode = addRecordDisplayOption || displayOption || 'sameTab'

  if (displayMode === 'modal') {
    router.push({
      query: {
        ...route.query,
        recordPageID: props.page.pageID,
        recordID: '0',
      }
    })
    return
  }

  const recordRoute = {
    name: 'page.record',
    params: { pageID: props.page.pageID, recordID: '0' },
  }

  if (displayMode === 'newTab') {
    window.open(router.resolve(recordRoute).href)
  } else {
    router.push(recordRoute)
  }
}

onMounted(() => pullRecords())
watch(() => props.record?.recordID, () => pullRecords())
watch(() => props.block.options, () => pullRecords(), { deep: true })

const offRefetch = $eventBus?.on('refetch-records', () => pullRecords())

onBeforeUnmount(() => {
  offRefetch?.()
})
</script>
