<template>
  <!-- One numbered, headed region of the single-page scenario editor (see
       FriaScenarioEditor.vue) — mirrors the design mockup's `.form-card` +
       numbered `.section-header` 1:1, shared by all five section components
       so they only own their field content, not this chrome. -->
  <section class="rounded-xl border border-surface bg-surface p-5 flex flex-col gap-5">
    <header class="flex items-start gap-3">
      <span
        class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-primary/10 text-primary text-xs font-semibold shrink-0"
      >
        {{ number }}
      </span>
      <div class="min-w-0">
        <h3 class="text-sm font-semibold">{{ title }}</h3>
        <p v-if="description" class="text-xs text-muted-color leading-relaxed mt-0.5 max-w-3xl">
          {{ description }}
        </p>
      </div>
    </header>

    <slot />
  </section>
</template>

<script setup>
defineProps({
  number: { type: [Number, String], required: true },
  title: { type: String, required: true },
  description: { type: String, default: '' },
})
</script>
