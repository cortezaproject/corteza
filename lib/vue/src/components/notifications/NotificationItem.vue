<template>
  <div class="notification-item-container px-3 pt-3 pb-2 border-bottom">
    <b-list-group-item
      class="notification-item border rounded bg-white p-3 position-relative"
      :class="{ 'read': notification.readAt }"
      @click="$emit('click', notification)"
    >
      <div
        class="action-menu bg-white pb-2 pl-2"
        style="margin-left: -1rem;"
      >
        <b-dropdown
          right
          variant="outline-extra-light"
          toggle-class="text-decoration-none border-0 dropdown-toggle-no-caret"
          no-caret
        >
          <template #button-content>
            <font-awesome-icon
              :icon="['fas', 'ellipsis-v']"
              class="text-secondary"
              style="margin-top: 0.3rem;"
            />
          </template>

          <b-dropdown-item
            v-if="!notification.readAt"
            @click.stop="$emit('mark-read', notification)"
          >
            <font-awesome-icon
              :icon="['far', 'envelope-open']"
              class="text-primary"
            />
            {{ $t('markAsRead') }}
          </b-dropdown-item>

          <b-dropdown-item
            v-else
            @click.stop="$emit('mark-unread', notification)"
          >
            <font-awesome-icon
              :icon="['far', 'envelope']"
              class="text-primary"
            />
            {{ $t('markAsUnread') }}
          </b-dropdown-item>

          <b-dropdown-divider />

          <c-input-confirm
            :text="$t('delete')"
            show-icon
            borderless
            variant="link"
            size="md"
            button-class="dropdown-item"
            icon-class="text-danger"
            class="w-100"
            @confirmed="$emit('delete', notification)"
          />
        </b-dropdown>
      </div>

      <component
        :is="notificationComponent"
        :notification="notification"
        class="notification-item-content"
      />
    </b-list-group-item>

    <div class="d-flex justify-content-end mt-2">
      <div
        :title="notification.createdAt"
        class="text-muted small cursor-pointer"
        @click="$emit('click', notification)"
      >
        {{ notification.createdAt | locFullDateTime }}
      </div>
    </div>
  </div>
</template>

<script>
import NotificationTypes from './types/index.js'
import CInputConfirm from '../input/CInputConfirm.vue'

export default {
  i18nOptions: {
    namespaces: 'notifications',
  },

  components: {
    CInputConfirm,
  },

  props: {
    notification: {
      type: Object,
      required: true,
    },
  },

  computed: {
    notificationComponent () {
      let { kind } = this.notification
      kind = kind.charAt(0).toUpperCase() + kind.slice(1)

      return NotificationTypes[`Notification${kind}`]
    },
  },
}
</script>

<style lang="scss" scoped>
.notification-item-container {
  &:hover {
    background-color: var(--light) !important;
  }

  .notification-item {
    transition: background-color 0.2s ease;
    cursor: pointer;

    &.read {
      .notification-item-content {
        opacity: 0.5 !important;
      }
    }
  }

  &:hover {
    .action-menu {
      opacity: 1 !important;
      pointer-events: auto;
    }
  }
  .action-menu {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    z-index: 2;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
  }
}
</style>
