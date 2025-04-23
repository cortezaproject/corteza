<template>
  <b-sidebar
    v-model="isVisible"
    header-class="d-flex align-items-center justify-content-between notification-sidebar-header bg-white pl-3 pr-2"
    body-class="d-flex flex-column overflow-hidden bg-white"
    :backdrop="isMobile"
    no-footer
    right
    shadow="sm"
    no-close-on-route-change
    no-close-on-esc
    width="400px"
  >
    <template #header>
      <h5
        class="text-primary mb-0"
      >
        <b>{{ $t('title') }}</b>
      </h5>

      <b-button
        variant="outline-light"
        class="d-flex align-items-center justify-content-center p-2 border-0 text-secondary"
        @click="isVisible = false"
      >
        <font-awesome-icon
          :icon="['fas', 'times']"
          class="h6 mb-0"
        />
      </b-button>
    </template>

    <notifications />
  </b-sidebar>
</template>

<script lang="js">
import Notifications from './Notifications.vue'
import { mapGetters, mapMutations } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'notifications',
  },

  components: {
    Notifications,
  },

  computed: {
    ...mapGetters({
      visible: 'notifications/visible',
    }),

    isVisible: {
      get () {
        return this.visible
      },

      set (visible) {
        this.setVisible(visible)
      },
    },

    isMobile () {
      return window.innerWidth < 576
    },
  },

  methods: {
    ...mapMutations({
      setVisible: 'notifications/setVisible',
    }),
  },
}
</script>

<style lang="scss">
.notification-sidebar-header {
  height: 4rem;
}
</style>
