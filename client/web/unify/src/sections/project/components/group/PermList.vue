<template>
  <div class="flex flex-col divide-y divide-surface rounded-lg border border-surface max-h-72 overflow-auto">
    <div v-for="row in rows" :key="row.id" class="flex items-center gap-2 px-3 py-2">
      <i :class="[cfg(row.kind).icon, cfg(row.kind).text, 'text-sm shrink-0']" />
      <span class="text-sm truncate flex-1 min-w-0">{{ row.name }}</span>
      <span v-if="baseline(row)" class="text-[10px] uppercase tracking-wider text-muted-color shrink-0">
        {{ t('project.groups.permList.fromGroup') }}
      </span>
      <SelectButton
        :model-value="getLevel(row)"
        :options="OPTIONS"
        option-label="label"
        option-value="value"
        :allow-empty="false"
        :disabled="disabled"
        size="small"
        @update:model-value="v => v && emit('set', row, v)"
      />
    </div>
    <div v-if="!rows.length" class="px-3 py-3 text-sm text-muted-color italic text-center">
      {{ t('project.groups.permList.empty') }}
    </div>
  </div>
</template>

<script setup>
import { kindConfig } from '@/sections/project/config/kinds'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  rows: { type: Array, default: () => [] },
  // (row) => 'none' | 'read' | 'write'
  getLevel: { type: Function, required: true },
  // (row) => boolean — level comes from the group baseline, not an explicit binding
  baseline: { type: Function, required: true },
  // Read-only: show levels but don't allow changing them.
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['set'])

// Display labels only — `value` is the persisted/compared level (see
// utils/rbac.js LEVELS) and must stay the literal lowercase string.
const OPTIONS = computed(() => [
  { label: t('project.groups.permList.levels.none'), value: 'none' },
  { label: t('project.groups.permList.levels.read'), value: 'read' },
  { label: t('project.groups.permList.levels.write'), value: 'write' },
])

const cfg = kindConfig
</script>
