<template>
  <div>
    <button
      ref="trigger"
      type="button"
      class="btn btn-link p-0"
      :aria-expanded="open ? 'true' : 'false'"
      :aria-controls="drawerID"
      @click="openDrawer"
    >
      {{ label }}
    </button>
    <aside
      v-if="open"
      :id="drawerID"
      ref="drawer"
      class="c311-help-drawer"
      role="dialog"
      aria-modal="true"
      :aria-label="title"
      data-c311-help-drawer
      @keydown.esc="closeDrawer"
      @keydown.tab="cycleFocus"
    >
      <div class="c311-help-drawer__header d-flex justify-content-between align-items-center">
        <h2 class="h5 mb-0">{{ title }}</h2>
        <button ref="close" type="button" class="btn btn-link" @click="closeDrawer">{{ closeLabel }}</button>
      </div>
      <div class="p-3">
        <slot>{{ content }}</slot>
      </div>
    </aside>
  </div>
</template>

<script>
export default {
  name: 'C311HelpDrawer',
  props: {
    label: { type: String, default: 'Help' },
    title: { type: String, default: 'Help' },
    closeLabel: { type: String, default: 'Close' },
    content: { type: String, default: '' },
    helpKey: { type: String, default: '' },
  },
  data: () => ({ open: false }),
  computed: {
    drawerID () { return `c311-help-${this._uid}` },
  },
  watch: {
    open (value) {
      if (!value) return
      this.$nextTick(() => this.$refs.close?.focus())
      this.$emit('open', this.helpKey)
    },
  },
  methods: {
    openDrawer () { this.open = true },
    closeDrawer () {
      this.open = false
      this.$nextTick(() => this.$refs.trigger?.focus())
      this.$emit('close', this.helpKey)
    },
    focusable () {
      return Array.from(this.$refs.drawer?.querySelectorAll('button, a[href], input, select, textarea, [tabindex]:not([tabindex="-1"])') || []).filter(element => !element.disabled)
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
.c311-help-drawer {
  position: fixed;
  z-index: 1040;
  top: 0;
  right: 0;
  bottom: 0;
  width: min(26rem, 100vw);
  overflow-y: auto;
  background: #fff;
  box-shadow: -0.25rem 0 1rem rgba(0, 0, 0, 0.18);
}

.c311-help-drawer__header {
  padding: 1rem;
  border-bottom: 1px solid #dee2e6;
}
</style>
