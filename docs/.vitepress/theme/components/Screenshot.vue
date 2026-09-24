<template>
  <figure class="screenshot">
    <div class="frame">
      <img
        class="shot light"
        :src="withBase(`/screenshots/${name}-light.webp`)"
        :alt="alt"
        loading="lazy"
        decoding="async"
        width="1440"
        height="900"
      />
      <img
        class="shot dark"
        :src="withBase(`/screenshots/${name}-dark.webp`)"
        :alt="alt"
        loading="lazy"
        decoding="async"
        width="1440"
        height="900"
      />
    </div>
    <figcaption v-if="caption">{{ caption }}</figcaption>
  </figure>
</template>

<script setup>
import { withBase } from 'vitepress'

defineProps({
  name: { type: String, required: true },
  alt: { type: String, required: true },
  caption: { type: String, default: '' },
})
</script>

<style scoped>
.screenshot {
  margin: 24px 0;
}

.frame {
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  overflow: hidden;
  background: var(--vp-c-bg-soft);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
}

.shot {
  display: block;
  width: 100%;
  height: auto;
}

.shot.dark {
  display: none;
}

figcaption {
  margin-top: 8px;
  font-size: 0.8125rem;
  color: var(--vp-c-text-2);
  text-align: center;
}
</style>

<style>
.dark .screenshot .shot.light {
  display: none;
}

.dark .screenshot .shot.dark {
  display: block;
}

.dark .screenshot .frame {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.35);
}
</style>
