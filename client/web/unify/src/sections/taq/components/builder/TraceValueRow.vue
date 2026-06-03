<template>
  <div class="flex flex-col">
    <div
      :class="[
        'flex items-start gap-1.5 py-0.5 rounded-sm',
        row.isContainer ? 'cursor-pointer hover:bg-emphasis' : '',
      ]"
      :style="{ paddingLeft: `${depth * 12}px` }"
      @click="row.isContainer && toggle()"
    >
      <i
        v-if="row.isContainer"
        :class="[
          'pi text-xs mt-1 w-3 text-muted-color shrink-0',
          expanded ? 'pi-chevron-down' : 'pi-chevron-right',
        ]"
      />

      <span class="trace-key text-sm text-color font-medium shrink-0 whitespace-nowrap">{{ row.key }}</span>

      <span class="trace-pill text-muted-color bg-emphasis shrink-0">{{ row.type }}</span>

      <span class="text-sm whitespace-nowrap">
        <template v-if="row.isContainer">
          <Badge :value="String(row.count)" severity="secondary" />
        </template>
        <template v-else>
          <span v-if="row.value === null || row.value === undefined" class="text-muted-color italic">null</span>
          <span v-else-if="typeof row.value === 'string'" class="trace-value text-color">"{{ row.value }}"</span>
          <span v-else class="trace-value text-color">{{ String(row.value) }}</span>
        </template>
      </span>
    </div>

    <div v-if="row.isContainer && expanded" class="flex flex-col gap-0.5">
      <TraceValueTree :value="row.raw" :search="search" :depth="depth + 1" />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import TraceValueTree from './TraceValueTree.vue'

const props = defineProps({
  row: { type: Object, required: true },
  search: { type: String, default: '' },
  depth: { type: Number, default: 0 },
})

const manualExpanded = ref(false)

// Auto-expand when search is active and something inside matches.
const expanded = computed(() => {
  if (props.search) return true
  return manualExpanded.value
})

watch(() => props.search, (s) => {
  if (!s) manualExpanded.value = false
})

function toggle() {
  manualExpanded.value = !manualExpanded.value
}

</script>

<style scoped>
.trace-pill {
  font-size: 11px;
  line-height: 1;
  padding: 3px 5px;
  border-radius: 3px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}
.trace-key,
.trace-value {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}
</style>
