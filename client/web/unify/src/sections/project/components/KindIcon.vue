<template>
  <!-- Colored square holding a resource-kind icon. The one place that owns the
       kind icon-square recipe (bg + ring + colored glyph); KindBadge and the
       permission matrix build on it. Pass a `kind` (resolved via kindConfig) or
       a ready `config` object, plus a `size` (sm | md | lg). -->
  <span
    class="inline-flex items-center justify-center shrink-0 ring-1"
    :class="[box, cfg.bg, cfg.ring]"
  >
    <i :class="[cfg.icon, cfg.text, iconText]" />
  </span>
</template>

<script setup>
import { kindConfig } from '@/sections/project/config/kinds'
import { computed } from 'vue'

defineOptions({ name: 'KindIcon' })

const props = defineProps({
  // Resource kind (module, role, …); resolved to a config via kindConfig.
  kind: { type: String, default: '' },
  // Pre-resolved config ({ bg, ring, icon, text }); wins over `kind` when set.
  config: { type: Object, default: null },
  // Square size: sm (w-5), md (w-6), lg (w-8).
  size: { type: String, default: 'sm' },
})

const SIZES = {
  sm: { box: 'w-5 h-5 rounded', iconText: 'text-[10px]' },
  md: { box: 'w-6 h-6 rounded-md', iconText: 'text-xs' },
  lg: { box: 'w-8 h-8 rounded-md', iconText: 'text-sm' },
}

const cfg = computed(() => props.config || kindConfig(props.kind))
const box = computed(() => (SIZES[props.size] || SIZES.sm).box)
const iconText = computed(() => (SIZES[props.size] || SIZES.sm).iconText)
</script>
