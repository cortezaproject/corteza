<template>
  <div class="flex flex-col gap-2">
    <div
      v-if="!triggers || triggers.length === 0"
      class="text-muted-color p-4 border rounded-lg bg-highlight text-center text-sm"
    >
      {{ $t('automation.workflows.editor.triggers.empty') }}
    </div>
    <DataTable v-else :value="triggers" class="text-sm">
      <Column field="resourceType" :header="$t('automation.workflows.editor.triggers.resourceType')">
        <template #body="{ data }">
          {{ formatResourceType(data.resourceType) }}
        </template>
      </Column>
      <Column field="eventType" :header="$t('automation.workflows.editor.triggers.eventType')" />
      <Column field="constraints" :header="$t('automation.workflows.editor.triggers.constraints')">
        <template #body="{ data }">
          <span v-for="(c, i) in data.constraints" :key="i" class="text-xs font-mono">
            <template v-if="c.name">
              {{ capitalize(c.name) }} {{ c.op }} "{{ (c.values || []).join(' or ') }}"
              <code v-if="i < data.constraints.length - 1"> and </code>
            </template>
          </span>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
const props = defineProps({
  triggers: { type: Array, default: () => [] },
})

function formatResourceType(rt) {
  if (!rt) return ''
  return rt.split(':').map(s => s.charAt(0).toUpperCase() + s.slice(1).toLowerCase()).join(' ')
}

function capitalize(s) {
  return s ? s.charAt(0).toUpperCase() + s.slice(1).toLowerCase() : ''
}
</script>
