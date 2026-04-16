<template>
  <div class="flex flex-col gap-2">
    <div
      v-if="!triggers || triggers.length === 0"
      class="text-muted-color p-4 border rounded-lg bg-highlight text-center text-sm"
    >
      {{ $t('automation.workflows.editor.triggers.empty') }}
    </div>
    <CResourceTable
      v-else
      :items="triggers"
      :fields="triggerFields"
    >
      <template #body-resourceType="{ data }">
        {{ formatResourceType(data.resourceType) }}
      </template>
      <template #body-constraints="{ data }">
        <span v-for="(c, i) in data.constraints" :key="i" class="text-xs font-mono">
          <template v-if="c.name">
            {{ capitalize(c.name) }} {{ c.op }} "{{ (c.values || []).join(' or ') }}"
            <code v-if="i < data.constraints.length - 1">and </code>
          </template>
        </span>
      </template>
    </CResourceTable>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'

const { CResourceTable } = components
const { t } = useI18n()

const props = defineProps({
  triggers: { type: Array, default: () => [] },
})

const triggerFields = [
  { key: 'resourceType', header: t('automation.workflows.editor.triggers.columns.resourceType') },
  { key: 'eventType', header: t('automation.workflows.editor.triggers.columns.eventType') },
  { key: 'constraints', header: t('automation.workflows.editor.triggers.columns.constraints') },
]

function formatResourceType(rt) {
  if (!rt) return ''
  return rt.split(':').map(s => s.charAt(0).toUpperCase() + s.slice(1).toLowerCase()).join(' ')
}

function capitalize(s) {
  return s ? s.charAt(0).toUpperCase() + s.slice(1).toLowerCase() : ''
}
</script>
