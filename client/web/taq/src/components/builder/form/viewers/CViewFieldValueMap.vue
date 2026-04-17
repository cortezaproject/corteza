<template>
  <div v-if="rows.length" class="w-full rounded-lg border border-surface overflow-hidden">
    <!-- Header -->
    <div class="flex text-xs font-medium text-muted-color bg-emphasis">
      <span class="px-2 py-1 shrink-0 w-[90px]">{{ $t('builder.preview.field') }}</span>
      <span class="px-2 py-1 flex-1">{{ $t('builder.preview.value') }}</span>
    </div>
    <!-- Rows -->
    <div
      v-for="row in rows"
      :key="row.field"
      class="flex text-xs border-t border-surface"
    >
      <span class="px-2 py-1 text-muted-color shrink-0 w-[90px] truncate" :title="row.fieldLabel">
        {{ row.fieldLabel }}
      </span>
      <div v-if="row.isRef" class="px-2 py-1 flex-1">
        <CViewReference
          :source="row.refLabel"
        />
      </div>
      <span
        v-else
        class="px-2 py-1 truncate flex-1"
        :class="row.value ? 'text-color-emphasis' : 'italic text-muted-color'"
        :title="String(row.value || '')"
      >
        {{ row.value || notSet }}
      </span>
    </div>
  </div>
  <span v-else class="italic text-muted-color">{{ notSet }}</span>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useComposeResourceStore } from '@planetcrust/human-vue/src/stores/useComposeResourceStore'
import CViewReference from './CViewReference.vue'

const { t } = useI18n()
const notSet = t('builder.preview.notSet')

const props = defineProps({
  /** Aggregate entries: [{ target, value, scope?, source? }] */
  modelValue: { type: Array, default: () => [] },
  namespaceID: { type: [String, Number], default: null },
  moduleID: { type: [String, Number], default: null },
  nodes: { type: Array, default: () => [] },
})

const store = useComposeResourceStore()
const fieldLabels = ref(new Map())

/**
 * Fetch module fields to resolve field handles → labels
 */
async function fetchFields() {
  if (!props.namespaceID || !props.moduleID) return

  try {
    const mod = await store.resolveModule(String(props.namespaceID), String(props.moduleID))
    if (mod?.fields) {
      const map = new Map()
      for (const f of mod.fields) {
        map.set(f.name, f.label || f.name)
      }
      fieldLabels.value = map
    }
  } catch {
    // Ignore
  }
}

watch(() => [props.namespaceID, props.moduleID], fetchFields)
onMounted(fetchFields)

const rows = computed(() => {
  return props.modelValue
    .filter(entry => entry.target)
    .map(entry => {
      const fieldLabel = fieldLabels.value.get(entry.target) || entry.target
      const isRef = !!(entry.scope && (entry.source || entry.expr))

      let refLabel = ''
      if (isRef) {
        refLabel = entry.source || entry.expr
      }

      return {
        field: entry.target,
        fieldLabel,
        value: entry.value ?? '',
        isRef,
        refLabel,
      }
    })
})
</script>
