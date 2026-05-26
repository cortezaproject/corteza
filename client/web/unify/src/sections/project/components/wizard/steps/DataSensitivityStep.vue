<template>
  <div class="h-full overflow-auto p-6">
    <div class="max-w-3xl mx-auto flex flex-col gap-5">
      <!-- Summary: count of data points per classification -->
      <div class="grid grid-cols-[repeat(auto-fill,minmax(8rem,1fr))] gap-3">
        <div
          v-for="c in summary"
          :key="c.id"
          class="rounded-lg border border-surface p-3 flex flex-col gap-2"
        >
          <Tag :value="c.label" :severity="c.severity" class="!text-[10px] !py-0 self-start" />
          <span class="text-lg font-semibold leading-none">{{ c.count }}</span>
        </div>
      </div>

      <!-- One section per level (most sensitive first), then unclassified gaps -->
      <div v-for="g in groups" :key="g.id">
        <div class="flex items-center gap-2 mb-2">
          <Tag :value="g.label" :severity="g.severity" />
          <span class="text-xs text-muted-color">
            {{ g.count }} {{ g.count === 1 ? 'item' : 'items' }}
          </span>
        </div>
        <div class="rounded-lg border border-surface divide-y divide-surface">
          <!-- Connections & modules; a module's same-level fields nest beneath it -->
          <template v-for="it in g.rows" :key="it.key">
            <div class="flex items-center gap-3 px-3 py-2">
              <i :class="[itemIcon(it), itemColor(it), 'shrink-0']" />
              <div class="min-w-0 flex-1">
                <div class="text-sm truncate">{{ it.name }}</div>
                <div class="text-xs text-muted-color truncate">{{ subtitle(it) }}</div>
              </div>
              <Select
                :model-value="it.level"
                :options="SENSITIVITY_OPTIONS"
                option-label="label"
                option-value="id"
                size="small"
                class="w-44 shrink-0"
                :disabled="disabled"
                @update:model-value="v => setLevel(it, v)"
              />
            </div>
            <div
              v-for="cf in it.childFields"
              :key="cf.key"
              class="flex items-center gap-3 px-3 py-2 pl-9"
            >
              <i class="pi pi-bars text-muted-color shrink-0" />
              <div class="min-w-0 flex-1">
                <div class="text-sm truncate">{{ cf.name }}</div>
              </div>
              <Select
                :model-value="cf.level"
                :options="SENSITIVITY_OPTIONS"
                option-label="label"
                option-value="id"
                size="small"
                class="w-44 shrink-0"
                :disabled="disabled"
                @update:model-value="v => setLevel(cf, v)"
              />
            </div>
          </template>

          <!-- Fields whose module sits in a different category -->
          <template v-for="fg in g.fieldGroups" :key="fg.moduleId">
            <div class="flex items-center gap-2 px-3 py-1.5 bg-emphasis text-xs text-muted-color">
              <i :class="kindConfig('module').icon" />
              <span class="font-medium">{{ fg.moduleName }}</span>
              <span>· fields</span>
            </div>
            <div v-for="it in fg.fields" :key="it.key" class="flex items-center gap-3 px-3 py-2 pl-9">
              <i class="pi pi-bars text-muted-color shrink-0" />
              <div class="min-w-0 flex-1">
                <div class="text-sm truncate">{{ it.name }}</div>
              </div>
              <Select
                :model-value="it.level"
                :options="SENSITIVITY_OPTIONS"
                option-label="label"
                option-value="id"
                size="small"
                class="w-44 shrink-0"
                :disabled="disabled"
                @update:model-value="v => setLevel(it, v)"
              />
            </div>
          </template>
        </div>
      </div>

      <div v-if="!items.length" class="text-sm text-muted-color italic">
        No modules, fields, or connections to classify yet.
      </div>
    </div>
  </div>
</template>

<script setup>
import { kindConfig } from '@/sections/project/config/kinds'
import { SENSITIVITY_OPTIONS, sensitivity } from '@/sections/project/config/sensitivity'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()

const LEVEL_ORDER = ['restricted', 'confidential', 'internal', 'public']

// Every classifiable data point: connections, modules, and module fields.
const items = computed(() => {
  const out = []
  for (const r of props.project.resources || []) {
    if (r.kind === 'connection') {
      out.push({ key: r.id, kind: 'connection', name: r.name, level: r.sensitivity || null, resourceId: r.id })
    } else if (r.kind === 'module') {
      out.push({ key: r.id, kind: 'module', name: r.name, level: r.sensitivity || null, resourceId: r.id })
      for (const f of r.fields || []) {
        out.push({
          key: `${r.id}:${f.id}`,
          kind: 'field',
          name: f.name || 'Untitled field',
          where: r.name,
          level: f.sensitivity || null,
          moduleId: r.id,
          fieldId: f.id,
        })
      }
    }
  }
  return out
})

const summary = computed(() => {
  const counts = {}
  for (const i of items.value) counts[i.level || 'unclassified'] = (counts[i.level || 'unclassified'] || 0) + 1
  return [
    ...LEVEL_ORDER.map(id => ({ id, label: sensitivity(id).label, severity: sensitivity(id).severity, count: counts[id] || 0 })),
    { id: 'unclassified', label: 'Unclassified', severity: 'secondary', count: counts.unclassified || 0 },
  ]
})

// Split a level's items into resource rows (connections/modules) and fields.
// A module's same-level fields nest directly under its row (no header); fields
// whose module sits in a different category get their own "Module · fields"
// header for context.
function buildGroup(id, label, severity, levelItems) {
  const modules = levelItems.filter(i => i.kind === 'module')
  const moduleIdsHere = new Set(modules.map(m => m.resourceId))
  const fields = levelItems.filter(i => i.kind === 'field')

  const childByModule = new Map()
  const orphanByModule = new Map()
  for (const f of fields) {
    if (moduleIdsHere.has(f.moduleId)) {
      if (!childByModule.has(f.moduleId)) childByModule.set(f.moduleId, [])
      childByModule.get(f.moduleId).push(f)
    } else {
      if (!orphanByModule.has(f.moduleId)) {
        orphanByModule.set(f.moduleId, { moduleId: f.moduleId, moduleName: f.where, fields: [] })
      }
      orphanByModule.get(f.moduleId).fields.push(f)
    }
  }

  const rows = levelItems
    .filter(i => i.kind !== 'field')
    .map(i => ({ ...i, childFields: i.kind === 'module' ? childByModule.get(i.resourceId) || [] : [] }))

  return { id, label, severity, count: levelItems.length, rows, fieldGroups: [...orphanByModule.values()] }
}

const groups = computed(() => {
  const g = []
  for (const id of LEVEL_ORDER) {
    const its = items.value.filter(i => i.level === id)
    if (its.length) g.push(buildGroup(id, sensitivity(id).label, sensitivity(id).severity, its))
  }
  const unclassified = items.value.filter(i => !i.level)
  if (unclassified.length) g.push(buildGroup('unclassified', 'Unclassified', 'secondary', unclassified))
  return g
})

const itemIcon = it => (it.kind === 'field' ? 'pi pi-bars' : kindConfig(it.kind).icon)
const itemColor = it => (it.kind === 'field' ? 'text-muted-color' : kindConfig(it.kind).text)
const subtitle = it =>
  it.kind === 'field' ? `Field · ${it.where}` : kindConfig(it.kind).label.replace(/s$/, '')

function setLevel(it, level) {
  if (it.kind === 'field') {
    store.updateField(props.project.id, it.moduleId, it.fieldId, { sensitivity: level || null })
  } else {
    store.updateResource(props.project.id, it.resourceId, { sensitivity: level || '' })
  }
}
</script>
