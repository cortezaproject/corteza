<template>
  <b-button
    v-b-tooltip.hover="{ title: labels.search || 'Search', delay: { show: 500, hide: 0 } }"
    variant="outline-extra-light"
    size="lg"
    class="nav-icon rounded-circle text-center border-0 d-flex align-items-center justify-content-center"
    @click="openSearch"
  >
    <font-awesome-icon
      :icon="['fas', 'search']"
      class="text-dark"
    />
  </b-button>
</template>

<script>
export default {
  name: 'CSearchButton',

  props: {
    labels: {
      type: Object,
      default: () => ({}),
    },
  },

  methods: {
    openSearch () {
      this.$root.$emit('search:open')
    },
  },
}
</script>

<style lang="scss" scoped>
$nav-icon-size: calc(var(--topbar-height) - 24px);

.nav-icon {
  width: $nav-icon-size !important;
  height: $nav-icon-size !important;
  padding: 0 !important;
}
</style>
