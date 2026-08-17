<template>
  <small v-if="chips.length" class="text-muted-color">
    {{
      hasRecord
        ? $t('block.content.interpolationFootnote')
        : $t('block.content.interpolationFootnoteNonRecord')
    }}
    <!-- The separator is a real space, not a margin: whitespace-only text
         between two elements is dropped at compile time, and these are snippets
         an author copies — a gap drawn in CSS would still hand a reader, a
         clipboard and a screen reader one run-on token. -->
    <template v-for="chip in chips" :key="chip">
      <component
        :is="insertable ? 'button' : 'span'"
        :type="insertable ? 'button' : undefined"
        class="c-expression-hint__chip"
        :class="{ 'c-expression-hint__chip--static': !insertable }"
        :title="insertable ? $t('block.content.interpolationFootnoteInsert') : undefined"
        @click="insertable && emit('insert', chip)"
      >
        <code>{{ chip }}</code>
      </component>
      {{ ' ' }}
    </template>
    <template v-if="dependsOnPlacement">
      {{ $t('block.content.interpolationFootnotePlacement') }}
    </template>
  </small>
</template>

<script setup>
import { computed } from 'vue'
import { resolvePath } from './catalog'

const props = defineProps({
  // ScopeEntry[] from buildScope() — the same value given to CInputExpression.
  scope: { type: Array, default: () => [] },
  // Members of record.values to name before falling back to an ellipsis.
  maxFields: { type: Number, default: 4 },
  // Set where the value is authored away from the page it renders on (charts
  // are namespace-level and can be placed on either kind of page).
  dependsOnPlacement: { type: Boolean, default: false },
  // Cleared where the hint has no cursor to insert into, e.g. beside a rich
  // text body — the chips then read as labels rather than offering an action
  // that would do nothing.
  insertable: { type: Boolean, default: true },
})

const emit = defineEmits(['insert'])

// Paths worth naming, in the order an author meets them. Each is shown only
// when it actually resolves against the scope, so a non-record page never
// advertises a record variable.
const CANDIDATES = ['recordID', 'ownerID', 'userID', 'record.values.*', 'user.name', 'user.email']

const hasRecord = computed(() => resolvePath(props.scope, ['recordID']).status === 'found')

const chips = computed(() => {
  const out = []

  for (const candidate of CANDIDATES) {
    if (candidate === 'record.values.*') {
      const res = resolvePath(props.scope, ['record', 'values'])
      if (res.status !== 'found') continue

      const fields = res.entry.fields
      if (!fields) {
        // Module not loaded — name the shape rather than a field that may not exist.
        out.push('${record.values.fieldName}')
        continue
      }

      for (const f of fields.slice(0, props.maxFields)) out.push(`\${record.values.${f.name}}`)
      if (fields.length > props.maxFields) out.push('${record.values.…}')
      continue
    }

    if (resolvePath(props.scope, candidate.split('.')).status === 'found') {
      out.push(`\${${candidate}}`)
    }
  }

  return out
})
</script>

<style scoped>
.c-expression-hint__chip {
  padding: 0;
  background: none;
  border: 0;
  cursor: pointer;
  vertical-align: baseline;
}

.c-expression-hint__chip code {
  padding: 0.05rem 0.3rem;
  border-radius: 4px;
  background: color-mix(in srgb, var(--p-primary-color) 12%, transparent);
  color: var(--p-primary-color);
  font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
  font-size: 0.85em;
}

.c-expression-hint__chip:hover code {
  background: color-mix(in srgb, var(--p-primary-color) 24%, transparent);
}

.c-expression-hint__chip--static {
  cursor: default;
}

.c-expression-hint__chip--static:hover code {
  background: color-mix(in srgb, var(--p-primary-color) 12%, transparent);
}
</style>
