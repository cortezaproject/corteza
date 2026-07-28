<template>
  <!-- Selectable card: icon + name + description, toggled by clicking
       anywhere on it. Shared by FriaTriggerSection (trigger conditions) and
       FriaVectorsSection (AI harm vectors) — the two scenario-editor
       sections that share this exact card shape. -->
  <button
    type="button"
    class="text-left rounded-lg border-2 p-4 transition-colors"
    :class="[
      selected ? 'border-primary bg-primary/5' : 'border-surface hover:border-primary/60',
      disabled ? 'cursor-not-allowed opacity-70' : 'cursor-pointer',
    ]"
    :disabled="disabled"
    :aria-pressed="selected"
    @click="$emit('toggle')"
  >
    <div class="flex items-center gap-2 mb-1.5">
      <span v-if="icon" class="text-base leading-none" aria-hidden="true">{{ icon }}</span>
      <span class="font-medium text-sm" :class="selected ? 'text-primary' : ''">{{ label }}</span>
    </div>
    <p class="text-xs text-muted-color leading-relaxed">{{ description }}</p>
  </button>
</template>

<script setup>
defineProps({
  label: { type: String, required: true },
  description: { type: String, default: '' },
  icon: { type: String, default: '' },
  selected: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})
defineEmits(['toggle'])
</script>
