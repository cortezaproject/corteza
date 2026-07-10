<template>
  <!-- Bordered chip with a colored kind-icon square + label, used for the role
       and module badges across the wizard so they can't drift apart. Pass
       `interactive` to render a togglable <button> (emits `click`); pass `muted`
       for the de-emphasized "off" state (smaller, muted label, dimmed icon). -->
  <!-- Always a <span> (never a native <button>, which picks up the app's default
       button background). When `interactive`, it behaves like a toggle via
       role/tabindex + keyboard handling. -->
  <span
    class="inline-flex items-center gap-2 pl-1.5 pr-3 py-1.5 rounded-lg border border-surface font-medium shrink-0 transition-all"
    :class="[
      muted
        ? 'text-xs text-muted-color hover:text-muted-color'
        : 'text-sm text-color hover:text-color',
      interactive ? 'cursor-pointer hover:bg-emphasis' : '',
    ]"
    :role="interactive ? 'button' : undefined"
    :tabindex="interactive ? 0 : undefined"
    @click="interactive && $emit('click')"
    @keydown.enter.prevent="interactive && $emit('click')"
    @keydown.space.prevent="interactive && $emit('click')"
  >
    <KindIcon
      :kind="kind"
      class="transition-opacity"
      :class="{ 'opacity-50': muted }"
    />
    <span>{{ label }}</span>
  </span>
</template>

<script setup>
import KindIcon from '@/sections/project/components/KindIcon.vue'

defineOptions({ name: 'KindBadge' })

defineProps({
  // Resource kind (module, role, …) — resolves the icon + colors.
  kind: { type: String, required: true },
  label: { type: String, default: '' },
  // De-emphasized "off" state: smaller, muted label + dimmed icon.
  muted: { type: Boolean, default: false },
  // Render as a togglable <button> with a hover background; emits `click`.
  interactive: { type: Boolean, default: false },
})

defineEmits(['click'])
</script>
