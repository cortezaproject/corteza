<template>
  <div>
    <p v-if="!changes.length" class="text-sm text-muted-color">
      {{ $t('project.publish.changes.none') }}
    </p>

    <ul v-else class="rounded-lg border border-surface overflow-hidden">
      <li
        v-for="change in changes"
        :key="change.path"
        class="flex items-center gap-3 px-3 py-2.5 border-l-2 border-b border-surface last:border-b-0"
        :class="isDangerous(change) ? 'border-l-red-500 bg-red-500/5' : 'border-l-transparent'"
      >
        <!-- Verb glyph: what happened, in one character. The old payload could
             not express this at all — added and removed modules were the same
             path with a different risk. -->
        <span
          class="inline-flex items-center justify-center w-5 h-5 rounded shrink-0 text-sm font-medium leading-none"
          :class="opClass(change.op)"
          :title="$t(`project.publish.changes.op.${change.op}`)"
        >
          {{ opGlyph(change.op) }}
        </span>

        <span
          class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded shrink-0 text-xs ring-1"
          :class="[
            kindConfig(change.kind).bg,
            kindConfig(change.kind).text,
            kindConfig(change.kind).ring,
          ]"
        >
          <i :class="[kindConfig(change.kind).icon, 'text-xs']" />
          {{ kindLabel(change) }}
        </span>

        <span class="flex-1 min-w-0">
          <span class="block text-sm truncate">{{ displayName(change) }}</span>
          <span v-if="change.detail" class="block text-xs text-muted-color truncate">
            {{ change.detail }}
          </span>
        </span>

        <!-- Counts, not adjectives: "1,284 records" is actionable in a way
             that "dangerous" is not. -->
        <span
          class="text-xs shrink-0 text-right"
          :class="isDangerous(change) ? 'text-red-500 font-medium' : 'text-muted-color'"
        >
          {{ riskLabel(change) }}
        </span>
      </li>
    </ul>

    <template v-if="inventory.length">
      <h4 class="text-xs font-medium uppercase tracking-wide text-muted-color mt-5 mb-2">
        {{ $t('project.publish.changes.inventoryHeading') }}
      </h4>
      <div class="flex flex-wrap gap-1.5">
        <span
          v-for="entry in inventory"
          :key="entry.kind"
          class="inline-flex items-center gap-1.5 px-2 py-1 rounded-md border border-surface bg-emphasis text-xs"
        >
          <i :class="[kindConfig(entry.kind).icon, kindConfig(entry.kind).text, 'text-xs']" />
          {{ $t(kindConfig(entry.kind).labelKey) }}
          <b class="font-medium">{{ entry.count }}</b>
          <i v-if="entry.added" class="not-italic text-green-500 font-medium">+{{ entry.added }}</i>
        </span>
      </div>
    </template>
  </div>
</template>

<script setup>
// Publish stage 1 — what this revision changes against its parent, and what it
// contains. Both come from the backend: changes from the deployment plan
// (op/kind/name/risk/records), inventory counts from the same project graph the
// Build canvas draws, so the two surfaces cannot disagree.
import { kindConfig } from '@/sections/project/config/kinds'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  changes: { type: Array, default: () => [] },
  // [{ kind, count, added }] — `added` is how many of them are new in this
  // revision, derived from the plan rather than a second graph fetch.
  inventory: { type: Array, default: () => [] },
})

const { t, n } = useI18n()

const isDangerous = change => change.risk === 'dangerous'

const opGlyph = op => ({ added: '+', removed: '−', changed: '~' })[op] || '•'

const opClass = op =>
  ({
    added: 'bg-green-500/10 text-green-500',
    removed: 'bg-red-500/10 text-red-500',
    changed: 'bg-amber-500/10 text-amber-500',
  })[op] || 'bg-emphasis text-muted-color'

// A field change names its kind as "field", which is not a resource kind in
// config/kinds — it borrows the owning module's colour and says so.
const kindLabel = change =>
  change.kind === 'field'
    ? t('project.publish.changes.field')
    : t(kindConfig(change.kind).singularKey)

const displayName = change => (change.module ? `${change.module}.${change.name}` : change.name)

// Records are only counted for destructive changes, and the count is
// best-effort — 0 means the backend could not read it, never "no records".
function riskLabel(change) {
  if (!isDangerous(change)) return t('project.publish.changes.risk.safe')
  if (!change.records) return t('project.publish.changes.risk.dangerous')
  return t('project.publish.changes.risk.records', { count: n(change.records) })
}
</script>
