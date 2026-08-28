<template>
  <div
    v-if="open"
    ref="modal"
    class="c311-focus-modal"
    role="dialog"
    aria-modal="true"
    :aria-label="title"
    data-c311-focus-modal
    @keydown.esc="close"
    @keydown.tab="cycleFocus"
  >
    <div class="c311-focus-modal__backdrop" @click="close" />
    <div class="c311-focus-modal__content" tabindex="-1">
      <h2 class="h5">{{ title }}</h2>
      <slot />
      <button ref="closeButton" type="button" class="btn btn-secondary" @click="close">{{ closeLabel }}</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'C311FocusModal',
  props: {
    value: Boolean,
    title: { type: String, default: 'Dialog' },
    closeLabel: { type: String, default: 'Close' },
  },
  data: () => ({ opener: null }),
  computed: {
    open: {
      get () { return this.value },
      set (value) { this.$emit('input', value) },
    },
  },
  watch: {
    open (value) {
      if (value) {
        this.opener = document.activeElement
        this.$nextTick(() => this.$refs.closeButton?.focus())
      } else {
        this.$nextTick(() => this.opener?.focus?.())
      }
    },
  },
  mounted () {
    if (!this.open) return
    this.opener = document.activeElement
    this.$nextTick(() => this.$refs.closeButton?.focus())
  },
  methods: {
    close () { this.open = false },
    focusable () {
      return Array.from(this.$refs.modal?.querySelectorAll('button, a[href], input, select, textarea, [tabindex]:not([tabindex="-1"])') || []).filter(element => !element.disabled)
    },
    cycleFocus (event) {
      const elements = this.focusable()
      if (!elements.length) return
      const first = elements[0]
      const last = elements[elements.length - 1]
      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault()
        last.focus()
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault()
        first.focus()
      }
    },
  },
}
</script>

<style scoped>
.c311-focus-modal {
  position: fixed;
  z-index: 1050;
  inset: 0;
  display: grid;
  place-items: center;
}

.c311-focus-modal__backdrop {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
}

.c311-focus-modal__content {
  position: relative;
  width: min(32rem, calc(100vw - 2rem));
  max-height: calc(100vh - 2rem);
  overflow: auto;
  padding: 1.25rem;
  background: #fff;
  border-radius: 0.25rem;
}
</style>
