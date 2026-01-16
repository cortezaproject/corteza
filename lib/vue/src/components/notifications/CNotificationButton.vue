<template>
  <b-button
    v-b-tooltip.hover="{ title: $t('notifications:title'), delay: { show: 500, hide: 0 } }"
    variant="outline-extra-light"
    size="lg"
    class="nav-icon rounded-circle text-center border-0 d-flex align-items-center justify-content-center position-relative"
    @click="toggleNotifications"
  >
    <font-awesome-icon
      :icon="['far', 'bell']"
      class="text-dark"
    />
    <b-badge
      v-if="unreadCount > 0 && !muted"
      variant="primary"
      pill
      class="position-absolute notification-badge"
    >
      {{ unreadCount > 9 ? '9+' : unreadCount }}
    </b-badge>
  </b-button>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'CNotificationButton',

  computed: {
    ...mapGetters({
      unreadCount: 'notifications/unreadCount',
      muted: 'notifications/muted',
    }),
  },

  methods: {
    ...mapActions({
      toggleVisibility: 'notifications/toggleVisibility',
    }),

    toggleNotifications () {
      this.toggleVisibility()
    },
  },
}
</script>

<style lang="scss" scoped>
.notification-badge {
  top: 0;
  right: 0;
  transform: translate(25%, -25%);
  font-size: 0.7rem;
}
</style>
