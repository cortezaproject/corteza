<script setup lang="ts">
import type { App as McpApp } from '@modelcontextprotocol/ext-apps'
import CFieldBoolViewer from '@planetcrust/human-vue/src/components/field/viewers/CFieldBoolViewer.vue'
import CFieldDateTimeViewer from '@planetcrust/human-vue/src/components/field/viewers/CFieldDateTimeViewer.vue'
import CFieldEmailViewer from '@planetcrust/human-vue/src/components/field/viewers/CFieldEmailViewer.vue'
import CFieldNumberViewer from '@planetcrust/human-vue/src/components/field/viewers/CFieldNumberViewer.vue'
import CFieldSelectViewer from '@planetcrust/human-vue/src/components/field/viewers/CFieldSelectViewer.vue'
import CFieldStringViewer from '@planetcrust/human-vue/src/components/field/viewers/CFieldStringViewer.vue'
import CFieldUrlViewer from '@planetcrust/human-vue/src/components/field/viewers/CFieldUrlViewer.vue'
import CResourceTable from '@planetcrust/human-vue/src/components/resource-table/CResourceTable.vue'
import Button from 'primevue/button'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { mcpAppKey } from '../shared/mount'
import {
  type Page,
  type Row,
  type ViewField,
  appendPage,
  columnsOf,
  isRefKind,
  parseResult,
  refLabel,
} from './recordLookup'

const { t } = useI18n()
const bridge = inject<McpApp>(mcpAppKey)!

// The webapp's viewers for the kinds whose value renders without reaching
// Human's API, which the sandboxed view has no access to.
const viewers: Record<string, unknown> = {
  Bool: CFieldBoolViewer,
  DateTime: CFieldDateTimeViewer,
  Email: CFieldEmailViewer,
  Number: CFieldNumberViewer,
  Select: CFieldSelectViewer,
  String: CFieldStringViewer,
  Url: CFieldUrlViewer,
}

const page = ref<Page>()
const args = ref<Record<string, unknown>>({})
const loading = ref(false)
const loadError = ref('')

bridge.ontoolinput = p => {
  args.value = p.arguments ?? {}
}

bridge.ontoolresult = r => {
  page.value = parseResult(r)
}

const columns = computed(() => (page.value ? columnsOf(page.value) : []))
const tableFields = computed(() =>
  columns.value.map(f => ({ key: f.name, header: f.label || f.name })),
)

// The shape the webapp's viewers read a field in.
function viewerField(f: ViewField) {
  return {
    name: f.name,
    kind: f.kind,
    isMulti: !!f.multi,
    options: f.kind === 'Select' ? { options: f.options ?? [] } : {},
  }
}

function linkOf(row: Row) {
  return page.value?.view?.links?.[row.recordID]
}

function open(row: Row) {
  const url = linkOf(row)
  if (url) bridge.openLink({ url })
}

async function loadMore() {
  if (!page.value?.cursor || loading.value) return

  loading.value = true
  loadError.value = ''
  try {
    const next = await bridge.callServerTool({
      name: 'compose_record_lookup',
      arguments: { ...args.value, recordID: undefined, pageCursor: page.value.cursor },
    })
    if (next.isError) {
      throw new Error(next.content?.map(c => ('text' in c ? c.text : '')).join(' ') || 'error')
    }
    page.value = appendPage(page.value, parseResult(next))
  } catch (e) {
    loadError.value = t('mcpApp.recordLookup.loadFailed', {
      error: e instanceof Error ? e.message : String(e),
    })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="p-3 text-sm">
    <div v-if="!page" id="status" class="text-muted-color">
      {{ t('mcpApp.recordLookup.waiting') }}
    </div>

    <template v-else>
      <div class="flex items-baseline justify-between gap-3 mb-2">
        <div class="font-semibold truncate">{{ page.view?.module.name }}</div>
        <div id="status" class="text-muted-color text-xs shrink-0">
          {{ t('mcpApp.recordLookup.count', page.rows.length) }}
        </div>
      </div>

      <CResourceTable
        :items="page.rows"
        :fields="tableFields"
        size="small"
        scrollable
        scroll-height="24rem"
        :row-class="() => 'cursor-pointer'"
        @row-click="open($event.data)"
      >
        <template v-for="col in columns" :key="col.name" #[`body-${col.name}`]="{ data }">
          <span v-if="isRefKind(col.kind)">{{ refLabel(data.values[col.name], page.refs) }}</span>
          <span v-else-if="col.kind === 'File'">
            {{ [data.values[col.name]].flat().filter(Boolean).length || '' }}
          </span>
          <component
            :is="viewers[col.kind]"
            v-else-if="viewers[col.kind]"
            :field="viewerField(col)"
            :record="data"
            value-only
          />
          <span v-else>{{ [data.values[col.name]].flat().join(', ') }}</span>
        </template>
      </CResourceTable>

      <div v-if="page.cursor || loadError" class="flex items-center gap-3 mt-2">
        <Button
          v-if="page.cursor"
          :label="loading ? t('mcpApp.recordLookup.loading') : t('mcpApp.recordLookup.loadMore')"
          :loading="loading"
          size="small"
          severity="secondary"
          @click="loadMore"
        />
        <span v-if="loadError" class="text-red-500 text-xs">{{ loadError }}</span>
      </div>
    </template>
  </div>
</template>
