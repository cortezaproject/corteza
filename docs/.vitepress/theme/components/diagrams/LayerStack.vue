<template>
  <figure class="dg layer-stack" aria-label="The layers you build with in Human">
    <div class="access">
      <span class="name">People and access</span>
      <span class="what">roles and permissions, applied to every layer below</span>
    </div>
    <div v-for="l in layers" :key="l.name" class="layer" :class="l.key">
      <span class="name">{{ l.name }}</span>
      <span class="what">{{ l.what }}</span>
    </div>
  </figure>
</template>

<script setup>
const layers = [
  { key: 'ai', name: 'AI', what: 'agents that read and act on your data' },
  {
    key: 'auto',
    name: 'Automation',
    what: 'TAQs that react to changes, run on schedules or on demand',
  },
  { key: 'ui', name: 'Interface', what: 'pages, blocks and charts people work in' },
  { key: 'data', name: 'Data', what: 'namespaces, modules, fields and records' },
]
</script>

<style scoped>
.dg {
  margin: 24px 0;
}

.layer-stack {
  display: grid;
  gap: 8px;
}

.layer,
.access {
  display: grid;
  grid-template-columns: 11rem 1fr;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 10px;
}

.layer {
  border: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg-soft);
}

.layer.data {
  background: var(--vp-c-brand-soft);
  border-color: color-mix(in srgb, var(--vp-c-brand-1) 30%, transparent);
}

.access {
  border: 1px dashed var(--vp-c-accent);
  background: color-mix(in srgb, var(--vp-c-accent) 8%, transparent);
}

.name {
  font-family: var(--vp-font-family-heading);
  font-weight: 500;
  font-size: 0.9375rem;
  color: var(--vp-c-text-1);
}

.what {
  font-size: 0.875rem;
  color: var(--vp-c-text-2);
}

@media (max-width: 560px) {
  .layer,
  .access {
    grid-template-columns: 1fr;
    gap: 2px;
  }
}
</style>
