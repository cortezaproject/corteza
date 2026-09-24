<script setup lang="ts">
import type { App as McpApp } from '@modelcontextprotocol/ext-apps'
import { inject, ref } from 'vue'
import { mcpAppKey } from '../shared/mount'

type LookupResult = {
  records?: unknown[]
  recordID?: string
  nextPageCursor?: string
}

const bridge = inject<McpApp>(mcpAppKey)!
const result = ref<LookupResult>()

bridge.ontoolresult = r => {
  result.value = (r.structuredContent ?? {}) as LookupResult
}

// A list carries 'records'; a lookup by recordID returns the record itself.
function count(r: LookupResult) {
  return Array.isArray(r.records) ? r.records.length : r.recordID ? 1 : 0
}
</script>

<template>
  <div class="p-4 text-sm">
    <div class="font-semibold">Human records</div>
    <div id="status" class="text-muted-color text-xs">
      <template v-if="!result">Waiting for the lookup…</template>
      <template v-else>
        {{ count(result) }} {{ count(result) === 1 ? 'record' : 'records'
        }}{{ result.nextPageCursor ? ', more available' : '' }}
      </template>
    </div>
  </div>
</template>
