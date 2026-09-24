<template>
  <figure class="dg data-model" aria-label="What a namespace contains">
    <div class="ns">
      <span class="title">Namespace</span>
      <span class="sub">one app</span>
    </div>
    <div class="parts">
      <div class="part">
        <span class="title">Modules</span>
        <span class="sub">the kinds of things you keep</span>
        <div class="chips">
          <span class="chip">Fields</span>
          <span class="chip">Records</span>
        </div>
      </div>
      <div class="part">
        <span class="title">Pages</span>
        <span class="sub">the screens people work in</span>
        <div class="chips">
          <span v-if="layouts" class="chip">Layouts</span>
          <span class="chip">Blocks</span>
        </div>
      </div>
      <div class="part">
        <span class="title">Charts</span>
        <span class="sub">summaries of module data</span>
        <div class="chips">
          <span class="chip muted">shown in Chart blocks</span>
        </div>
      </div>
    </div>
  </figure>
</template>

<script setup>
defineProps({ layouts: { type: Boolean, default: false } })
</script>

<style scoped>
.dg {
  margin: 24px 0;
}

.data-model {
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  padding: 14px;
  background: var(--vp-c-bg-soft);
}

.ns {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin: 2px 4px 12px;
}

.parts {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.part {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg);
}

.title {
  font-family: var(--vp-font-family-heading);
  font-weight: 500;
  font-size: 0.9375rem;
  color: var(--vp-c-text-1);
}

.ns .title {
  font-size: 1rem;
}

.sub {
  font-size: 0.8125rem;
  color: var(--vp-c-text-2);
}

.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

.chip {
  font-size: 0.8125rem;
  padding: 2px 10px;
  border-radius: 999px;
  background: var(--vp-c-brand-soft);
  color: var(--vp-c-brand-1);
}

.chip.muted {
  background: transparent;
  border: 1px solid var(--vp-c-divider);
  color: var(--vp-c-text-2);
}

@media (max-width: 640px) {
  .parts {
    grid-template-columns: 1fr;
  }
}
</style>
