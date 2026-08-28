<template>
  <div class="c311-app-shell" :class="`c311-app-shell--${mode}`" data-c311-app-shell>
    <a href="#c311-main-content" class="sr-only sr-only-focusable c311-skip-link">
      {{ skipLabel }}
    </a>
    <header class="c311-app-shell__header">
      <slot name="header">
        <div class="d-flex align-items-center justify-content-between px-3 py-2 border-bottom">
          <strong>{{ brand }}</strong>
          <slot name="nav" />
          <slot name="actions" />
        </div>
      </slot>
    </header>
    <main
      id="c311-main-content"
      ref="main"
      class="c311-app-shell__main"
      tabindex="-1"
      data-c311-main
      :aria-labelledby="headingID"
    >
      <h1 :id="headingID" ref="heading" class="h2 mb-3" tabindex="-1">
        {{ title }}
      </h1>
      <slot />
    </main>
    <c311-status-announcer :message="statusMessage" />
  </div>
</template>

<script>
import C311StatusAnnouncer from './C311StatusAnnouncer.vue'

export default {
  name: 'C311AppShell',
  components: { C311StatusAnnouncer },
  props: {
    mode: {
      type: String,
      default: 'public',
    },
    brand: {
      type: String,
      default: 'City 311',
    },
    title: {
      type: String,
      default: '',
    },
    skipLabel: {
      type: String,
      default: 'Skip to main content',
    },
    statusMessage: {
      type: String,
      default: '',
    },
  },
  computed: {
    headingID () {
      return `c311-heading-${this._uid}`
    },
  },
  methods: {
    focusMain () {
      this.$nextTick(() => (this.$refs.heading || this.$refs.main)?.focus())
    },
  },
}
</script>

<style scoped>
.c311-app-shell__main {
  min-width: 0;
  padding: 1.5rem;
}

.c311-skip-link:focus {
  position: fixed;
  z-index: 1050;
  top: 0.5rem;
  left: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: #fff;
  color: #111;
  box-shadow: 0 0 0 3px #1a73e8;
}
</style>
