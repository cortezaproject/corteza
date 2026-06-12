<template>
  <div class="flex flex-col gap-5">
    <!-- Summary: field count per classification, styled like the resource
         metric cards in the graph panel. -->
    <div class="flex flex-wrap gap-2">
      <div
        v-for="c in summary"
        :key="c.id"
        class="rounded-lg border border-surface px-2.5 py-1.5 flex items-center gap-2 whitespace-nowrap"
      >
        <Tag :value="c.label" :severity="c.severity" class="!text-[10px] !py-0" />
        <span class="text-base font-semibold leading-none">{{ c.count }}</span>
      </div>
    </div>

    <!-- One card per module; classify each of its fields. -->
    <div
      v-for="m in moduleCards"
      :key="m.id"
      class="rounded-lg border border-surface overflow-hidden"
    >
      <div class="flex items-center gap-2 px-3 py-2 bg-emphasis border-b border-surface">
        <span
          class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
          :class="[moduleCfg.bg, moduleCfg.ring]"
        >
          <i :class="[moduleCfg.icon, moduleCfg.text, 'text-xs']" />
        </span>
        <span class="font-medium text-sm truncate">{{ m.name }}</span>
      </div>
      <div class="divide-y divide-surface">
        <div
          v-for="f in m.fields"
          :key="f.fieldId"
          class="flex items-center gap-3 px-3 py-2"
        >
          <div class="min-w-0 flex-1">
            <div class="text-sm truncate">{{ f.name }}</div>
          </div>
          <Select
            :model-value="f.level"
            :options="SENSITIVITY_OPTIONS"
            option-label="label"
            option-value="id"
            size="small"
            class="w-44 shrink-0"
            placeholder="Unclassified"
            :disabled="disabled"
            @update:model-value="v => setLevel(f, v)"
          />
        </div>
        <div v-if="!m.fields.length" class="px-3 py-2 text-sm text-muted-color italic">
          No fields yet.
        </div>
      </div>
    </div>

    <div v-if="!moduleCards.length" class="text-sm text-muted-color italic">
      No modules to classify yet. Add modules and fields in the Data Model step first.
    </div>
  </div>
</template>

<script setup>
import { kindConfig } from '@/sections/project/config/kinds'
import { SENSITIVITY_OPTIONS, sensitivity } from '@/sections/project/config/sensitivity'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useToast } from 'primevue/usetoast'
import { computed } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const toast = useToast()
const moduleCfg = kindConfig('module')

// Most sensitive first; unclassified is handled separately (level === null).
const LEVEL_ORDER = ['restricted', 'confidential', 'internal', 'public']

// One card per module, each carrying its classifiable fields.
const moduleCards = computed(() =>
  (props.project.resources || [])
    .filter(m => m.kind === 'module')
    .map(m => ({
      id: m.id,
      name: m.name,
      fields: (m.fields || []).map(f => ({
        fieldId: f.id,
        name: f.name || 'Untitled field',
        level: f.sensitivity || null,
        moduleId: m.id,
      })),
    })),
)

// Whole-project field counts per level (drives the summary strip).
const summary = computed(() => {
  const counts = {}
  for (const m of moduleCards.value) {
    for (const f of m.fields) {
      const k = f.level || 'unclassified'
      counts[k] = (counts[k] || 0) + 1
    }
  }
  return [
    ...LEVEL_ORDER.map(id => ({
      id,
      label: sensitivity(id).label,
      severity: sensitivity(id).severity,
      count: counts[id] || 0,
    })),
    { id: 'unclassified', label: 'Unclassified', severity: 'secondary', count: counts.unclassified || 0 },
  ]
})

async function setLevel(f, level) {
  try {
    await store.updateField(props.project.id, f.moduleId, f.fieldId, { sensitivity: level || null })
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: 'Could not update classification',
      detail: err.message,
      life: 4000,
    })
  }
}
</script>
