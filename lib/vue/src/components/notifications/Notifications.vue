<template>
  <div class="h-100 d-flex flex-column">
    <b-tabs
      v-model="activeTab"
      card
      fill
      nav-class="border-bottom"
      content-class="h-100 overflow-hidden"
      class="h-100 d-flex flex-column"
    >
      <template #tabs-end>
        <div
          class="d-flex align-items-center justify-content-end"
          style="min-width: 6rem;"
        >
          <b-button
            v-if="hasUnread"
            v-b-tooltip.hover
            variant="outline-light"
            class="p-2 border-0 d-flex align-items-center justify-content-center"
            style="width: 2rem; height: 2rem;"
            :title="$t('markAllAsRead')"
            @click="handleMarkAllAsRead"
          >
            <font-awesome-icon
              :icon="['fas', 'check-double']"
              class="h6 mb-0 text-primary"
            />
          </b-button>

          <b-button
            v-b-tooltip.hover
            variant="outline-light"
            class="p-2 border-0 d-flex align-items-center justify-content-center"
            style="width: 2rem; height: 2rem;"
            :title="$t(muted ? 'unmute' : 'mute')"
            @click="toggleMuted"
          >
            <font-awesome-icon
              :icon="['fas', muted ? 'bell-slash' : 'bell']"
              class="h6 mb-0"
              :class="{ 'text-secondary': muted, 'text-primary': !muted }"
            />
          </b-button>
        </div>
      </template>

      <b-tab
        :title="$t('unread')"
        active
        class="d-flex flex-column h-100 p-0"
      >
        <div class="overflow-auto flex-grow-1 h-100">
          <div
            v-if="loading"
            class="d-flex justify-content-center p-5"
          >
            <b-spinner variant="primary" />
          </div>

          <b-list-group v-else-if="notifications.length > 0">
            <notification-item
              v-for="notification in notifications"
              :key="notification.notificationID"
              :notification="notification"
              @click="onNotificationClick(notification)"
              @mark-read="onMarkAsRead"
              @delete="onDeleteNotification"
            />

            <div
              v-if="hasMorePages"
              class="text-center my-3"
            >
              <b-button
                variant="outline-primary"
                size="sm"
                :disabled="loadingMore"
                @click="loadMore()"
              >
                <b-spinner
                  v-if="loadingMore"
                  small
                />
                <span v-else>{{ $t('loadMore') }}</span>
              </b-button>
            </div>
          </b-list-group>

          <div
            v-else
            class="text-center p-5"
          >
            <font-awesome-icon
              :icon="['far', 'bell']"
              class="text-secondary mb-3"
              size="3x"
            />
            <p class="text-secondary">
              {{ $t('empty') }}
            </p>
            <p class="text-muted small">
              {{ $t('emptyDescription') }}
            </p>
          </div>
        </div>
      </b-tab>

      <b-tab
        :title="$t('all')"
        class="d-flex flex-column h-100 p-0"
      >
        <div class="overflow-auto flex-grow-1 h-100">
          <div
            v-if="loading"
            class="d-flex justify-content-center p-5"
          >
            <b-spinner variant="primary" />
          </div>

          <b-list-group v-else-if="notifications.length > 0">
            <notification-item
              v-for="notification in notifications"
              :key="notification.notificationID"
              :notification="notification"
              @click="onNotificationClick(notification)"
              @mark-read="onMarkAsRead"
              @delete="onDeleteNotification"
            />

            <div
              v-if="hasMorePages"
              class="text-center my-3"
            >
              <b-button
                variant="outline-primary"
                size="sm"
                :disabled="loadingMore"
                @click="loadMore()"
              >
                <b-spinner
                  v-if="loadingMore"
                  small
                />
                <span v-else>{{ $t('loadMore') }}</span>
              </b-button>
            </div>
          </b-list-group>

          <div
            v-else
            class="text-center p-5"
          >
            <font-awesome-icon
              :icon="['far', 'bell']"
              class="text-secondary mb-3"
              size="3x"
            />
            <p class="text-secondary">
              {{ $t('empty') }}
            </p>
            <p class="text-muted small">
              {{ $t('emptyDescription') }}
            </p>
          </div>
        </div>
      </b-tab>
    </b-tabs>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import NotificationItem from './NotificationItem.vue'

export default {
  i18nOptions: {
    namespaces: 'notifications',
  },

  components: {
    NotificationItem,
  },

  data () {
    return {
      activeTab: 0,

      loading: false,

      loadingMore: false,
    }
  },

  computed: {
    ...mapGetters({
      notifications: 'notifications/notifications',
      hasUnread: 'notifications/hasUnread',
      hasMorePages: 'notifications/hasMorePages',
      muted: 'notifications/muted',
    }),
  },

  watch: {
    activeTab: {
      handler () {
        this.loading = true

        this.setPageCursor(null)

        this.loadNotifications()
          .finally(() => {
            setTimeout(() => {
              this.loading = false
            }, 300)
          })
      },
    },
  },

  methods: {
    ...mapActions({
      fetchNotifications: 'notifications/fetchNotifications',
      markAsRead: 'notifications/markAsRead',
      markAllAsRead: 'notifications/markAllAsRead',
      deleteNotification: 'notifications/deleteNotification',
      setPageCursor: 'notifications/setPageCursor',
      toggleMuted: 'notifications/toggleMuted',
    }),

    onNotificationClick ({ notificationID, readAt }) {
      if (!readAt) {
        this.markAsRead(notificationID)
      }
    },

    onMarkAsRead ({ notificationID }) {
      this.markAsRead(notificationID)
    },

    onDeleteNotification ({ notificationID }) {
      this.deleteNotification(notificationID)
        .then(() => {
          this.toastSuccess(this.$t('notificationDeleted'))
        })
        .catch(() => {
          this.toastError(this.$t('notificationDeletedError'))
        })
    },

    handleMarkAllAsRead () {
      this.markAllAsRead()
        .then(() => {
          this.toastSuccess(this.$t('allMarkedAsRead'))
        })
        .catch(() => {
          this.toastError(this.$t('markAllAsReadError'))
        })
    },

    loadNotifications () {
      return this.fetchNotifications({ unreadOnly: this.activeTab === 0 })
    },

    loadMore () {
      this.loadingMore = true

      return this.loadNotifications()
        .finally(() => {
          setTimeout(() => {
            this.loadingMore = false
          }, 300)
        })
    },
  },
}
</script>
