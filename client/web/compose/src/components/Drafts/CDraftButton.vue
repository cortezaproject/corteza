<template>
  <b-button
    v-b-tooltip.hover="{ title: $t('drafts:title'), delay: { show: 500, hide: 0 } }"
    variant="outline-extra-light"
    size="lg"
    class="nav-icon rounded-circle text-center border-0 d-flex align-items-center justify-content-center position-relative"
    @click="toggleDrafts"
  >
    <font-awesome-icon
      :icon="['far', 'file']"
      class="text-dark"
    />
    <b-badge
      v-if="draftCount > 0"
      variant="primary"
      pill
      class="position-absolute draft-badge"
    >
      {{ draftCount > 9 ? '9+' : draftCount }}
    </b-badge>
  </b-button>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'

export default {
  name: 'CDraftButton',

  computed: {
    ...mapGetters({
      drafts: 'drafts/getAllDrafts',
    }),

    draftCount () {
      return (this.drafts || []).length
    },
  },

  methods: {
    ...mapActions({
      toggleVisibility: 'drafts/toggleVisibility',
    }),

    toggleDrafts () {
      this.toggleVisibility()
    },
  },
}
</script>

<style lang="scss" scoped>
.draft-badge {
  top: 0;
  right: 0;
  transform: translate(25%, -25%);
  font-size: 0.7rem;
}
</style>
