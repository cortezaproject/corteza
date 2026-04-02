<template>
  <PageBlock :block="block" @refreshBlock="loadRevisions">
    <div class="flex flex-col h-full overflow-hidden">
      <!-- Revisions disabled on module -->
      <div v-if="revisionsDisabled" class="flex items-center justify-center h-full p-3 text-muted-color italic">
        {{ $t('block.recordRevisions.viewer.errors.disabled-on-module') }}
      </div>

      <!-- Loading -->
      <div v-else-if="loading" class="flex items-center justify-center h-full">
        <ProgressSpinner style="width: 28px; height: 28px" />
      </div>

      <!-- Load button (when not preloading) -->
      <div v-else-if="!preloadRevisions && !loadedRevisions" class="flex items-center justify-center h-full">
        <Button
          :label="$t('block.recordRevisions.viewer.show-revisions', { revision: record?.revision || 0 })"
          severity="secondary"
          @click="loadRevisions"
        />
      </div>

      <!-- Revisions table -->
      <template v-else>
        <div v-if="!revisions.length" class="flex items-center justify-center h-full p-3 text-muted-color italic">
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
            field="timestamp"
            :header="$t('block.recordRevisions.viewer.revisions.columns.timestamp.label')"
          >
            <template #body="{ data }">
              {{ formatTimestamp(data.timestamp) }}
            </template>
          </Column>
          <Column
            :header="''"
            style="width: 10rem"
          >
            <template #body="{ data }">
              <span v-if="data.changes && data.changes.length" class="text-primary text-sm">
                {{ $t('block.recordRevisions.viewer.show-changes', { count: data.changes.length }) }}
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
                {{ data.old !== undefined ? (Array.isArray(data.old) ? data.old.join(', ') : data.old) : '-' }}
              </template>
            </Column>
            <Column :header="$t('block.recordRevisions.viewer.changes.columns.new-value.label')">
              <template #body="{ data }">
                {{ data.new !== undefined ? (Array.isArray(data.new) ? data.new.join(', ') : data.new) : '-' }}
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
import { useI18n } from 'vue-i18n'
import PageBlock from './PageBlock.vue'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI', null)
const $eventBus = inject('$eventBus', null)

const loading = ref(false)
const loadedRevisions = ref(false)
const revisions = ref([])
const showChangesDialog = ref(false)
const selectedRevision = ref(null)

const options = computed(() => props.block.options || {})
const preloadRevisions = computed(() => options.value.preload)
const revisionsDisabled = computed(() => false) // Would check module config

function formatTimestamp(ts) {
  if (!ts) return '-'
  try {
    return new Date(ts).toLocaleString()
  } catch {
    return ts
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

  loading.value = true
  loadedRevisions.value = true

  try {
    const result = await props.block.fetch($ComposeAPI, props.record, options.value.sortDirection)

    revisions.value = (result || []).map(r => {
      const changes = (r.changes || []).map(c => ({
        key: c.key,
        label: c.key, // Would resolve field label from module
        old: c.old,
        new: c.new,
      }))

      return {
        revision: r.revision,
        operation: r.operation,
        timestamp: r.timestamp,
        userID: r.userID,
        changes,
      }
    })
  } catch (e) {
    console.error('Failed to load revisions:', e)
    revisions.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (preloadRevisions.value) {
    loadRevisions()
  }
})

watch(() => props.record?.recordID, () => {
  if (preloadRevisions.value) loadRevisions()
})

const offRefetch = $eventBus?.on('refetch-records', () => {
  if (preloadRevisions.value) {
    loadRevisions()
  }
})

onBeforeUnmount(() => {
  offRefetch?.()
})
</script>
