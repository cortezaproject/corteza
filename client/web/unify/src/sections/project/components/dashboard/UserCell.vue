<template>
  <!-- Owner cell: initials avatar + name (mirrors the demo's user-cell). Avatar
       color is derived deterministically from the name. Blank → em-dash. -->
  <span v-if="name" class="inline-flex items-center gap-2 whitespace-nowrap">
    <span
      class="w-5 h-5 rounded-full text-[10px] font-semibold flex items-center justify-center shrink-0"
      :class="avatarCls"
    >
      {{ initials }}
    </span>
    <span class="text-sm">{{ name }}</span>
  </span>
  <span v-else class="text-muted-color">—</span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ name: { type: String, default: '' } })

const PALETTE = [
  'bg-blue-100 text-blue-700 dark:bg-blue-500/20 dark:text-blue-300',
  'bg-purple-100 text-purple-700 dark:bg-purple-500/20 dark:text-purple-300',
  'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-300',
  'bg-amber-100 text-amber-700 dark:bg-amber-500/20 dark:text-amber-300',
  'bg-rose-100 text-rose-700 dark:bg-rose-500/20 dark:text-rose-300',
]

const initials = computed(() =>
  (props.name || '')
    .split(/\s+/)
    .map(w => w.replace(/[^A-Za-z]/g, '')[0])
    .filter(Boolean)
    .slice(0, 2)
    .join('')
    .toUpperCase(),
)
const avatarCls = computed(() => {
  const sum = (props.name || '').split('').reduce((a, c) => a + c.charCodeAt(0), 0)
  return PALETTE[sum % PALETTE.length]
})
</script>
