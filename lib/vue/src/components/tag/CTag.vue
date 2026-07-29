<template>
  <span
    class="inline-flex items-center gap-1.5 rounded-md ring-1 whitespace-nowrap"
    :class="[bg, ring, text, small ? 'px-2 py-1' : 'px-2.5 py-1.5']"
  >
    <i :class="[icon, small ? 'text-xs' : 'text-sm']" />
    <span class="font-medium leading-none" :class="small ? 'text-xs' : 'text-sm'">
      {{ label }}
    </span>
  </span>
</template>

<script setup>
// A status tag: CChip's tinted badge widened to hold its own label, so the
// icon and the text sit on one coloured surface and share one colour.
//
// Same colouring contract as CChip — the caller passes the tint, the ring and
// the icon as classes, because the palettes are per-domain (a project status,
// a queue state, a severity) and the component has no business knowing them.
// The one addition is `text`: it applies to the wrapper, so BOTH the icon and
// the label inherit it. That is the whole difference from CChip, where the
// label is always muted and only the badge carries colour.
defineProps({
  icon: { type: String, default: 'pi pi-circle-fill' },
  // Tailwind classes, e.g. 'bg-green-500/10' / 'ring-green-500/30' /
  // 'text-green-500'. Alpha tints are deliberate in the defaults: they read on
  // both themes, where a flat -100 fill only works on light.
  bg: { type: String, default: 'bg-emphasis' },
  ring: { type: String, default: 'ring-surface' },
  text: { type: String, default: 'text-muted-color' },
  label: { type: String, required: true },
  // Tighter padding and type, for tags inside dropdown items or table cells
  // where a full-size one dominates the row.
  small: { type: Boolean, default: false },
})
</script>
