<template>
  <nav class="c311-main-nav d-flex flex-wrap align-items-center" :aria-label="label" data-c311-main-nav>
    <router-link
      v-for="item in visibleItems"
      :key="item.route || item.href"
      class="btn btn-link btn-sm"
      :to="item.route"
      :data-c311-route="item.route"
      :aria-current="$route && $route.path === item.route ? 'page' : undefined"
    >
      {{ item.label }}
    </router-link>
  </nav>
</template>

<script>
export default {
  name: 'C311MainNav',
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: 'Primary navigation',
    },
  },
  computed: {
    visibleItems () {
      return this.items.filter(item => {
        if (!item.capability) return true
        if (!this.$C311 || !this.$C311.can(item.capability)) return false
        return !item.scope || this.$C311.hasScope(item.scope)
      })
    },
  },
}
</script>
