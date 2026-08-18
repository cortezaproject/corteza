<template>
  <PageBlock :block="block" :record="record" @refreshBlock="pullRecords">
    <div
      v-if="!isConfigured"
      class="flex items-center justify-center h-full p-3 text-muted-color italic"
    >
      {{ $t('block.recordOrganizer.notConfigured') }}
    </div>

    <div v-else-if="loading" class="flex items-center justify-center h-full">
      <ProgressSpinner style="width: 28px; height: 28px" />
    </div>

    <div v-else class="h-full flex flex-col">
      <!-- Add record button -->
      <div v-if="canAddRecord" class="p-3 border-b border-surface">
        <span v-tooltip.bottom="addRecordDisabled ? $t('block.noRecordPage') : ''">
          <Button
            :label="$t('block.recordOrganizer.addNewRecord')"
            severity="primary"
            size="small"
            :disabled="addRecordDisabled"
            @click="createNewRecord"
          />
        </span>
      </div>

      <!-- Records -->
      <CDraggableList
        v-model="records"
        class="flex-1 overflow-auto p-3"
        :drag-key="blockKey"
        :group="dragGroup"
        :disabled="!canDrag"
        @add="onAdd"
        @update="onUpdate"
      >
        <div v-if="!records.length" class="text-muted-color text-sm">
          {{ $t('block.recordOrganizer.noRecords') }}
        </div>

        <div
          v-for="record in records"
          :key="record.recordID"
          data-drag-item
          data-organizer-card
          class="record-card p-3 mb-3 border border-surface rounded-border cursor-pointer hover:bg-emphasis transition-colors"
          :class="{ grab: canDrag }"
          @click="handleRecordClick(record)"
        >
          <h6 v-if="labelFieldName" class="font-medium mb-1 text-color">
            <CFieldViewer
              v-if="labelFieldDef"
              :field="labelFieldDef"
              :record="record"
              :namespace="namespace"
              :disable-click="true"
              :value-only="true"
            />
            <template v-else>
              {{ getFieldValue(record, labelFieldName) || $t('block.record.preview.untitled') }}
            </template>
          </h6>
          <p v-if="descriptionFieldName" class="text-sm text-muted-color mb-0">
            <CFieldViewer
              v-if="descriptionFieldDef"
              :field="descriptionFieldDef"
              :record="record"
              :namespace="namespace"
              :disable-click="true"
              :value-only="true"
            />
            <template v-else>{{ getFieldValue(record, descriptionFieldName) }}</template>
          </p>
        </div>
      </CDraggableList>
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, inject } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components, useRecordStore, useModuleStore, usePageStore } from '@planetcrust/human-vue'
import PageBlock from './PageBlock.vue'
import { evaluatePrefilter, getFieldFilter } from '../../../lib/record-filter'

const { CFieldViewer, CDraggableList } = components

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $Auth = inject('$Auth', {})
const $eventBus = inject('$eventBus', null)
const $ComposeAPI = inject('$ComposeAPI', null)
const $toast = inject('$toast', null)
const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const recordStore = useRecordStore()
const moduleStore = useModuleStore()
const pageStore = usePageStore()
const $recordRoutes = inject('$recordRoutes', null)

const loading = ref(false)
const records = ref([])

const options = computed(() => props.block.options || {})
const labelFieldName = computed(() => options.value.labelField || '')
const descriptionFieldName = computed(() => options.value.descriptionField || '')
const isConfigured = computed(() => !!options.value.moduleID)
const canAddRecord = computed(() => !!options.value.moduleID)

const organizerModule = computed(() => moduleStore.getByID(options.value.moduleID))

// The record page for THIS block's module — not the page the block sits on.
//
// Using props.page.pageID sent every "add" to whatever page was hosting the
// organizer, which is only ever right by coincidence. RecordListBlock resolves
// it this way and that is the behaviour being matched.
const recordPageID = computed(() => {
  // Custom routes carry the module themselves, so no public page is needed.
  if ($recordRoutes) return 'admin'

  const moduleID = organizerModule.value?.moduleID
  if (!moduleID) return null

  return (pageStore.set || []).find(p => p.moduleID === moduleID)?.pageID || null
})
// Shown but disabled when there is nowhere to go, with the reason in a
// tooltip: a button that silently does nothing is the complaint that opened
// this issue, and hiding it leaves the configurator no clue either.
const addRecordDisabled = computed(() => !recordPageID.value)

function fieldDef(fieldName) {
  if (!fieldName) return null
  return organizerModule.value?.fields?.find(f => f.name === fieldName) || null
}

const labelFieldDef = computed(() => fieldDef(labelFieldName.value))
const descriptionFieldDef = computed(() => fieldDef(descriptionFieldName.value))

function getFieldValue(record, fieldName) {
  if (!record || !fieldName) return ''
  const val = record.values?.[fieldName]
  if (Array.isArray(val)) return val[0] || ''
  return val || ''
}

// The records this column holds: the block's own prefilter AND its group.
// Shared by the listing and by organize(), which scopes its repositioning to
// the same set.
function buildQuery() {
  const { filter: prefilter, groupField, group } = options.value
  const filterParts = []

  if (prefilter) {
    const record = props.record
    const user = $Auth?.user || {}
    filterParts.push(
      `(${evaluatePrefilter(prefilter, {
        record,
        user,
        recordID: record?.recordID || '0',
        ownerID: record?.ownedBy || '0',
        userID: user?.userID || '0',
      })})`,
    )
  }

  // Through the shared filter helper, never string interpolation: it escapes
  // the value, and reads an empty group as IS NULL — the ungrouped column —
  // where a bare `= ''` finds nothing on a text column and is a hard postgres
  // error on a numeric one.
  if (groupField && group !== undefined) {
    const kind = fieldDef(groupField)?.kind || 'String'
    const condition = getFieldFilter(groupField, kind, group, '=')
    if (condition) filterParts.push(`(${condition})`)
  }

  return filterParts.join(' AND ')
}

// silent skips the spinner: a column refreshing because something moved
// somewhere else swaps its cards in place rather than blanking to a spinner.
async function pullRecords({ silent = false } = {}) {
  if (!options.value.moduleID) return

  if (!silent) loading.value = true

  try {
    const { namespaceID } = props.namespace
    const { moduleID, positionField } = options.value

    // The shared record store requires the module in the module store; ensure
    // it's loaded before listing (organizer modules may differ from the page's).
    await moduleStore.findByID({ namespaceID, moduleID })

    const query = buildQuery()
    const sort = positionField || 'updatedAt'

    const { set = [] } = await recordStore.list({ namespaceID, moduleID, query, sort })
    records.value = set
  } catch (e) {
    console.error('Failed to load records for organizer:', e)
    records.value = []
  } finally {
    loading.value = false
  }
}

// Same destination rule as createNewRecord: the record belongs to the
// organizer's module, so it opens on that module's record page. This call site
// had the identical defect and was fixed alongside it, though the reported
// symptom was only ever the add button.
function handleRecordClick(record) {
  if (!recordPageID.value) return

  const { displayOption } = options.value

  if (displayOption === 'modal' && !$recordRoutes) {
    router.push({
      query: {
        ...route.query,
        recordPageID: recordPageID.value,
        recordID: record.recordID,
      },
    })
    return
  }

  const recordRoute = $recordRoutes
    ? $recordRoutes.view(organizerModule.value.moduleID, record.recordID)
    : {
        name: 'page.record',
        params: { pageID: recordPageID.value, recordID: record.recordID },
      }

  if (displayOption === 'newTab') {
    window.open(router.resolve(recordRoute).href)
  } else {
    router.push(recordRoute)
  }
}

// prefillQuery carries the organizer's own bucket into the new record.
//
// The block shows records where groupField equals group, so a record created
// from it belongs in that bucket — RecordView and the admin create view both
// read refField/refValue and set the value for us.
function prefillQuery() {
  const { groupField, group } = options.value
  if (!groupField || group === undefined || group === '') return {}

  return { refField: groupField, refValue: group }
}

function createNewRecord() {
  if (!recordPageID.value) return

  const { addRecordDisplayOption, displayOption } = options.value
  const displayMode = addRecordDisplayOption || displayOption || 'sameTab'
  const refQuery = prefillQuery()

  if (displayMode === 'modal' && !$recordRoutes) {
    router.push({
      query: {
        ...route.query,
        recordPageID: recordPageID.value,
        recordID: '0',
        ...refQuery,
      },
    })
    return
  }

  const recordRoute = $recordRoutes
    ? $recordRoutes.create(organizerModule.value.moduleID)
    : {
        name: 'page.record',
        params: { pageID: recordPageID.value, recordID: '0' },
      }

  if (Object.keys(refQuery).length > 0) {
    recordRoute.query = { ...(recordRoute.query || {}), ...refQuery }
  }

  if (displayMode === 'newTab') {
    window.open(router.resolve(recordRoute).href)
  } else {
    router.push(recordRoute)
  }
}

// ---- Board drag and drop ----
//
// Columns are separate block instances of one board, joined by the shared
// draggable's group: the name carries the moduleID, so a column refuses a card
// from a board over another module before the pointer ever lands.
//
// The move itself is local — the card leaves one column's list and enters the
// other's — so neither of the two columns you are watching refetches anything.
// The server is told afterwards, and only a rejection moves the card back.
const blockKey = computed(() => String(props.block.blockID ?? ''))

// Repositioning needs a position field; regrouping needs a group field. With
// neither there is nothing a drop could change.
const canDrag = computed(() => !!(options.value.positionField || options.value.groupField))

const dragGroup = computed(() => ({
  name: `record-organizer-${options.value.moduleID}`,
  pull: canDrag.value,
  put: canDrag.value,
}))

// Corteza's calcNewPosition: sit one past the card you were dropped behind.
// The card is already in the list at `index`, so index - 1 is that card.
function positionFor(index) {
  const { positionField } = options.value
  if (!positionField || index <= 0) return 0
  const before = records.value[index - 1]
  return parseInt(before?.values?.[positionField] || 0, 10) + 1
}

// A card dropped in from another column of the board.
function onAdd({ item, index, fromKey }) {
  return organize(item, index, fromKey)
}

// A card moved within this column: it never left, so this column is both ends.
function onUpdate({ item, index }) {
  return organize(item, index, blockKey.value)
}

async function organize(record, index, fromKey) {
  if (!record?.recordID) return

  const { moduleID, positionField, groupField, group } = options.value
  const args = [
    { name: 'recordID', value: String(record.recordID) },
    { name: 'filter', value: buildQuery() },
    { name: 'positionField', value: positionField || '' },
    { name: 'position', value: String(positionFor(index)) },
  ]

  // Dropping into a column is what sets the record's group — an empty group is
  // meaningful (the ungrouped column), so it is sent as an empty string.
  if (groupField) {
    args.push({ name: 'groupField', value: groupField })
    args.push({ name: 'group', value: group ?? '' })
  }

  // What the two ends of the move are, so they can ignore an event that only
  // tells them what they already show.
  const organized = {
    recordID: String(record.recordID),
    moduleID,
    fromKey,
    toKey: blockKey.value,
  }

  try {
    await $ComposeAPI.recordExec({
      procedure: 'organize',
      namespaceID: props.namespace.namespaceID,
      moduleID,
      args,
    })
    $eventBus?.emit('refetch-records', { organized })
  } catch (err) {
    console.error('Failed to organize record:', err)
    $toast?.toastErrorHandler(t('block.recordOrganizer.moveFailed'))(err)
    // The card sits where the server refused to put it; both ends re-read and
    // it lands back where it came from.
    $eventBus?.emit('refetch-records', { organized: { ...organized, failed: true } })
  }
}

onMounted(() => pullRecords())
watch(
  () => props.record?.recordID,
  () => pullRecords(),
)
watch(
  () => props.block.options,
  () => pullRecords(),
  { deep: true },
)

// A bare refetch — a record was saved somewhere — reloads the column as before.
// One carrying an organized move is narrower: the two columns it moved between
// already show it, everyone else swaps their cards in silently.
const offRefetch = $eventBus?.on('refetch-records', payload => {
  const organized = payload?.organized
  if (!organized) return pullRecords()

  const mine = organized.fromKey === blockKey.value || organized.toKey === blockKey.value

  if (organized.failed) {
    if (mine) pullRecords({ silent: true })
    return
  }

  if (mine) return
  pullRecords({ silent: true })
})

onBeforeUnmount(() => {
  offRefetch?.()
})
</script>

<style scoped>
.grab {
  cursor: grab;
}

.grab:active {
  cursor: grabbing;
}
</style>
