<template>
  <PageBlock :block="block" :record="record" @refreshBlock="loadRevisions">
    <div class="flex flex-col h-full overflow-hidden">
      <!-- Reading revisions is its own permission on the record -->
      <div
        v-if="!canSearchRevisions"
        class="flex items-center justify-center h-full p-3 text-muted-color italic"
      >
        {{ $t('block.noPermission') }}
      </div>

      <!-- Revisions disabled on module -->
      <div
        v-else-if="revisionsDisabled"
        class="flex items-center justify-center h-full p-3 text-muted-color italic"
      >
        {{ $t('block.recordRevisions.viewer.errors.disabled-on-module') }}
      </div>

      <!-- Loading -->
      <div v-else-if="loading" class="flex items-center justify-center h-full">
        <ProgressSpinner style="width: 28px; height: 28px" />
      </div>

      <!-- Load button (when not preloading) -->
      <div
        v-else-if="!preloadRevisions && !loadedRevisions"
        class="flex items-center justify-center h-full"
      >
        <Button
          :label="
            $t('block.recordRevisions.viewer.show-revisions', { revision: record?.revision || 0 })
          "
          severity="secondary"
          @click="loadRevisions"
        />
      </div>

      <!-- Revisions table -->
      <template v-else>
        <!-- A failed read and an empty history are different answers -->
        <div
          v-if="loadError"
          :title="loadError"
          class="flex items-center justify-center h-full p-3 text-red-500 italic"
        >
          {{ $t('block.recordRevisions.viewer.errors.load-failed') }}
        </div>

        <div
          v-else-if="!revisions.length"
          class="flex items-center justify-center h-full p-3 text-muted-color italic"
        >
          {{ $t('block.recordRevisions.viewer.errors.no-revisions') }}
        </div>

        <DataTable
          v-else
          :value="revisions"
          scrollable
          scroll-height="flex"
          class="flex-1"
          size="small"
          :row-class="() => 'cursor-pointer'"
          @row-click="onRowClick"
        >
          <Column field="revision" header="#" style="width: 4rem" />
          <Column
            field="operation"
            :header="$t('block.recordRevisions.viewer.revisions.columns.operation.label')"
          >
            <template #body="{ data }">
              {{ $t(`block.recordRevisions.viewer.operations.${data.operation}`) }}
            </template>
          </Column>
          <Column
            field="userID"
            :header="$t('block.recordRevisions.viewer.revisions.columns.user.label')"
          >
            <template #body="{ data }">
              <CFieldViewer
                v-if="data.userID && data.userID !== '0'"
                :field="userField"
                :record="data"
                :namespace="namespace"
                value-only
              />
              <template v-else>-</template>
            </template>
          </Column>
          <Column
            field="timestamp"
            :header="$t('block.recordRevisions.viewer.revisions.columns.timestamp.label')"
          >
            <template #body="{ data }">
              {{ formatTimestamp(data.timestamp) }}
            </template>
          </Column>
          <Column :header="''" style="width: 10rem">
            <template #body="{ data }">
              <span v-if="data.changes && data.changes.length" class="text-primary text-sm">
                {{
                  $t('block.recordRevisions.viewer.show-changes', { count: data.changes.length })
                }}
              </span>
            </template>
          </Column>
        </DataTable>

        <!-- Expanded row changes -->
        <Dialog
          v-model:visible="showChangesDialog"
          :header="`Revision #${selectedRevision?.revision || ''}`"
          modal
          :style="{ width: '600px' }"
        >
          <DataTable
            v-if="selectedRevision?.changes"
            :value="selectedRevision.changes"
            size="small"
          >
            <Column
              field="label"
              :header="$t('block.recordRevisions.viewer.changes.columns.field.label')"
            />
            <Column :header="$t('block.recordRevisions.viewer.changes.columns.old-value.label')">
              <template #body="{ data }">
                <CFieldViewer
                  v-if="data.field && selectedRevision.oldRecord && hasValue(data.old)"
                  :field="data.field"
                  :record="selectedRevision.oldRecord"
                  :namespace="namespace"
                  value-only
                  disable-click
                />
                <template v-else>{{ rawValue(data.old) }}</template>
              </template>
            </Column>
            <Column :header="$t('block.recordRevisions.viewer.changes.columns.new-value.label')">
              <template #body="{ data }">
                <CFieldViewer
                  v-if="data.field && selectedRevision.newRecord && hasValue(data.new)"
                  :field="data.field"
                  :record="selectedRevision.newRecord"
                  :namespace="namespace"
                  value-only
                  disable-click
                />
                <template v-else>{{ rawValue(data.new) }}</template>
              </template>
            </Column>
          </DataTable>
        </Dialog>
      </template>
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, inject } from 'vue'
import { components, useModuleStore } from '@planetcrust/human-vue'
import { compose } from '@planetcrust/human-js'
import PageBlock from './PageBlock.vue'

const { CFieldViewer } = components

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI', null)
const $eventBus = inject('$eventBus', null)

const moduleStore = useModuleStore()

const loading = ref(false)
const loadedRevisions = ref(false)
const loadError = ref('')
const revisions = ref([])
const showChangesDialog = ref(false)
const selectedRevision = ref(null)

const options = computed(() => props.block.options || {})
const preloadRevisions = computed(() => options.value.preload)
const canSearchRevisions = computed(() => props.record?.canSearchRevisions !== false)

// The module the page is built on: it holds the field definitions the changes
// are written against, and the switch that decides whether revisions exist.
const pageModule = computed(() => {
  const moduleID = props.page?.moduleID
  if (!moduleID || moduleID === '0') return null
  return moduleStore.getByID(moduleID) || null
})

// Claimed only where the module is actually known — an unresolved module is not
// evidence that revisions are off.
const revisionsDisabled = computed(
  () => !!pageModule.value && pageModule.value.config?.recordRevisions?.enabled === false,
)

// The revision's author, rendered by the same viewer any User field gets.
const userField = Object.freeze({ isSystem: true, name: 'userID', label: '', kind: 'User' })

function formatTimestamp(ts) {
  if (!ts) return '-'
  try {
    return new Date(ts).toLocaleString()
  } catch {
    return ts
  }
}

function hasValue(v) {
  return Array.isArray(v) && v.length > 0
}

function rawValue(v) {
  return hasValue(v) ? v.join(', ') : '-'
}

// One side of a revision as a record, so field viewers can render the values the
// way the record page renders them — users and references resolved, not raw IDs.
function sideAsRecord(mod, changes, side) {
  if (!mod) return null

  const draft = { values: {} }
  let any = false

  for (const c of changes) {
    if (!c.field || !hasValue(c[side])) continue
    any = true
    const value = c.field.isMulti ? c[side] : c[side][0]
    if (c.field.isSystem) draft[c.key] = value
    else draft.values[c.key] = value
  }

  if (!any) return null

  try {
    return new compose.Record(mod, draft)
  } catch {
    return null
  }
}

function onRowClick({ data }) {
  if (data.changes && data.changes.length) {
    selectedRevision.value = data
    showChangesDialog.value = true
  }
}

async function loadRevisions() {
  if (!$ComposeAPI || !props.record?.recordID || props.record.recordID === '0') {
    return
  }

  if (!canSearchRevisions.value || revisionsDisabled.value) {
    return
  }

  loading.value = true
  loadedRevisions.value = true
  loadError.value = ''

  try {
    const result = await props.block.fetch($ComposeAPI, props.record, options.value.sortDirection)
    const mod = pageModule.value

    revisions.value = (result || []).map(r => {
      const changes = (r.changes || []).map(c => {
        const field = mod?.findField ? mod.findField(c.key) || null : null

        return {
          key: c.key,
          label: field ? field.label || field.name : c.key,
          field,
          old: c.old,
          new: c.new,
        }
      })

      return {
        revision: r.revision,
        operation: r.operation,
        timestamp: r.timestamp,
        userID: r.userID,
        changes,
        oldRecord: sideAsRecord(mod, changes, 'old'),
        newRecord: sideAsRecord(mod, changes, 'new'),
      }
    })
  } catch (e) {
    console.error('Failed to load revisions:', e)
    revisions.value = []
    loadError.value = e?.message || 'unknown error'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (preloadRevisions.value) {
    loadRevisions()
  }
})

watch(
  () => props.record?.recordID,
  () => {
    if (preloadRevisions.value) loadRevisions()
  },
)

const offRefetch = $eventBus?.on('refetch-records', () => {
  if (preloadRevisions.value) {
    loadRevisions()
  }
})

onBeforeUnmount(() => {
  offRefetch?.()
})
</script>
