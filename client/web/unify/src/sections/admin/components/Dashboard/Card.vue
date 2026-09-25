<template>
  <section class="flex flex-col gap-3 rounded-border border border-surface bg-surface p-4 min-w-0">
    <header v-if="title || $slots.actions" class="flex items-baseline justify-between gap-3">
      <div class="min-w-0">
        <h2 class="text-base font-semibold text-color truncate">{{ title }}</h2>
        <p v-if="subtitle" class="text-sm text-muted-color">{{ subtitle }}</p>
      </div>
      <slot name="actions" />
    </header>
    <slot />
  </section>
</template>

<script setup>
defineProps({
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
})
</script>
