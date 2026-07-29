<template>
  <div>
    <Message v-if="!blockers.length" severity="success" :closable="false" class="mb-3">
      {{
        modules.length
          ? $t('project.publish.migration.allAccountedFor')
          : $t('project.publish.migration.nothingToMigrate')
      }}
    </Message>
    <Message v-else severity="warn" :closable="false" class="mb-3">
      {{ $t('project.publish.migration.intro') }}
    </Message>

    <div
      v-for="mod in modules"
      :key="mod.handle"
      class="rounded-lg border border-surface overflow-hidden mb-2.5"
    >
      <div class="flex items-center gap-2 px-3 py-2 bg-emphasis border-b border-surface">
        <i :class="[kindConfig('module').icon, kindConfig('module').text, 'text-xs']" />
        <b class="text-sm font-medium">{{ mod.handle }}</b>
        <span v-if="mod.records" class="ml-auto text-xs text-muted-color">
          {{ $t('project.publish.migration.moduleRecords', { count: $n(mod.records) }) }}
        </span>
      </div>

      <!-- Undecided first: these are the only reason this stage can block. -->
      <div
        v-for="lost in mod.undecided"
        :key="lost.path"
        class="flex flex-wrap items-center gap-2 px-3 py-2.5 border-b border-surface last:border-b-0 border-l-2"
        :class="
          decisions[lost.path]
            ? 'border-l-green-500 bg-green-500/5'
            : 'border-l-amber-500 bg-amber-500/5'
        "
      >
        <span class="text-sm">
          <b class="font-medium">{{ lost.name || lost.module }}</b>
          <span class="text-muted-color">
            {{
              lost.kind === 'field'
                ? $t('project.publish.migration.fieldGone')
                : $t('project.publish.migration.moduleGone')
            }}
          </span>
        </span>

        <div v-if="decisions[lost.path]" class="ml-auto flex items-center gap-2">
          <span class="text-xs text-green-500 font-medium">{{ outcomeLabel(lost) }}</span>
          <Button
            :label="$t('project.publish.migration.change')"
            severity="secondary"
            text
            size="small"
            :disabled="disabled"
            @click="decide(lost.path, null)"
          />
        </div>

        <div v-else class="ml-auto flex flex-wrap items-center gap-2">
          <Select
            v-if="lost.kind === 'field' && mod.targets.length"
            :model-value="null"
            :options="mod.targets"
            :placeholder="$t('project.publish.migration.moveInto')"
            :disabled="disabled"
            size="small"
            class="w-56"
            @update:model-value="target => decide(lost.path, `move:${target}`)"
          />
          <Button
            :label="$t('project.publish.migration.discard')"
            severity="danger"
            outlined
            size="small"
            :disabled="disabled"
            @click="decide(lost.path, 'drop')"
          />
        </div>
      </div>

      <!-- Everything the system can settle on its own, stated with its
           consequence so nothing is silently converted. -->
      <div
        v-for="field in mod.settled"
        :key="field.targetField"
        class="flex items-center gap-2 px-3 py-2 text-sm border-b border-surface last:border-b-0"
      >
        <span class="font-mono text-xs">
          {{ field.sourceField || $t('project.publish.migration.newField') }}
        </span>
        <i class="pi pi-arrow-right text-muted-color text-xs" />
        <span class="font-mono text-xs">{{ field.targetField }}</span>
        <span class="ml-auto text-xs text-muted-color">{{ settledLabel(field) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
// Publish stage 2 — where records go. The backend's suggested mapping already
// settles everything with a defensible default (same-name fields copy, retyped
// fields cast, new fields start empty); this stage exists for what it CANNOT
// settle: data whose destination no longer exists. Those are the only rows that
// ask the user anything, which is what makes the question worth reading.
import { kindConfig } from '@/sections/project/config/kinds'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  // Suggested per-module mappings from the deployment plan.
  mappings: { type: Array, default: () => [] },
  changes: { type: Array, default: () => [] },
  // path → 'drop' | 'move:<targetField>'
  decisions: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:decisions'])

const { t } = useI18n()

// Destructive changes needing a human call: data exists and its destination is
// gone. A retyped field is NOT here — casting is a defensible default, so it is
// reported rather than asked about.
const blockers = computed(() =>
  props.changes.filter(
    c => c.risk === 'dangerous' && c.op === 'removed' && !props.decisions[c.path],
  ),
)

const modules = computed(() => {
  const removed = props.changes.filter(c => c.risk === 'dangerous' && c.op === 'removed')

  // A removed module has no mapping entry of its own (there is nothing left to
  // map into), so it gets a block built from the change alone.
  const fromMappings = props.mappings.map(m => {
    const fields = m.fields || []
    return {
      handle: m.module,
      records: recordsFor(m.module),
      // Only newly added fields are offered as a destination: moving data into
      // a field that already receives its own is how you lose the other column.
      targets: fields.filter(f => f.op === 'default').map(f => f.targetField),
      undecided: removed.filter(c => c.kind === 'field' && c.module === m.module),
      settled: fields,
    }
  })

  const removedModules = removed
    .filter(c => c.kind === 'module')
    .map(c => ({
      handle: c.name,
      records: c.records || 0,
      targets: [],
      undecided: [c],
      settled: [],
    }))

  return [...removedModules, ...fromMappings].filter(m => m.undecided.length || m.settled.length)
})

function recordsFor(moduleHandle) {
  const hit = props.changes.find(c => c.module === moduleHandle && c.records)
  return hit?.records || 0
}

function decide(path, value) {
  const next = { ...props.decisions }
  if (value === null) delete next[path]
  else next[path] = value
  emit('update:decisions', next)
}

function outcomeLabel(lost) {
  const decision = props.decisions[lost.path]
  if (!decision) return ''
  if (decision === 'drop') return t('project.publish.migration.outcome.discarded')
  return t('project.publish.migration.outcome.moved', { field: decision.slice('move:'.length) })
}

const settledLabel = field =>
  ({
    copy: t('project.publish.migration.outcome.copied'),
    cast: t('project.publish.migration.outcome.converted'),
    default: t('project.publish.migration.outcome.startsEmpty'),
  })[field.op] || field.op
</script>
