<template>
  <div
    v-if="errors.length"
    ref="summary"
    class="alert alert-danger"
    role="alert"
    tabindex="-1"
    :aria-label="title"
    data-c311-error-summary
  >
    <h2 class="h5 mb-2">
      {{ title }}
    </h2>
    <ul class="mb-0 pl-3">
      <li v-for="(error, index) in errors" :key="`${error.field}-${index}`">
        <a :href="`#${fieldID(error.field)}`" @click.prevent="focusField(error.field)">
          {{ error.message || error.code || error.field }}
        </a>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'C311ErrorSummary',
  props: {
    errors: {
      type: Array,
      default: () => [],
    },
    title: {
      type: String,
      default: 'There is a problem',
    },
    focusOnUpdate: {
      type: Boolean,
      default: true,
    },
  },
  watch: {
    errors: {
      deep: true,
      handler (value) {
        if (!this.focusOnUpdate || !value.length) return
        this.$nextTick(() => this.$refs.summary?.focus())
      },
    },
  },
  methods: {
    fieldID (field) {
      const normalized = String(field || '')
        .replace(/^#/, '')
        .replace(/^\//, '')
        .replace(/[^a-zA-Z0-9_\-/]/g, '-')
        .replace(/\//g, '-')
      return normalized.startsWith('c311-') ? normalized : `c311-${normalized}`
    },
    focusField (field) {
      const raw = String(field || '').replace(/^#/, '').replace(/^\//, '').replace(/\//g, '-')
      const candidates = [this.fieldID(field), raw, raw.split('-').pop()]
      const element = candidates.map(id => document.getElementById(id)).find(Boolean)
      if (element) element.focus()
    },
  },
}
</script>
