<template>
  <b-button
    variant="outline-extra-light"
    size="lg"
    class="nav-icon rounded-circle text-center border-0 d-flex align-items-center justify-content-center position-relative"
    @click="toggleNotifications"
  >
    <font-awesome-icon
      :icon="notificationsIcon.icon"
      :class="notificationsIcon.class"
    />
    <b-badge
      v-if="unreadCount > 0 && !muted"
      variant="danger"
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

    notificationsIcon () {
      if (this.unreadCount > 0 && !this.muted) {
        return {
          icon: ['fas', 'bell'],
          class: 'text-primary',
        }
      }

      return {
        icon: ['far', 'bell'],
        class: 'text-dark',
      }
    },
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
