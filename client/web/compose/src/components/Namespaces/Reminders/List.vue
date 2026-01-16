<template>
  <div
    class="d-flex flex-column h-100"
  >
    <div
      class="text-center bg-white sticky-top pt-3"
    >
      <b-button
        data-test-id="button-add-reminder"
        size="sm"
        variant="outline-primary"
        @click="$emit('edit')"
      >
        + {{ $t('reminder.add') }}
      </b-button>
    </div>

    <div class="d-flex flex-column flex-fill overflow-auto">
      <div
        v-for="(r, i) in sortedReminders"
        :key="r.reminderID"
        class="reminder-item-container px-3 pt-3 pb-2 border-bottom"
      >
        <hr
          v-if="r.dismissedAt && sortedReminders[i - 1] ? !sortedReminders[i - 1].dismissedAt : false "
          class="mt-0"
        >

        <div
          class="reminder-item border rounded bg-white p-3 position-relative"
          :class="{ 'dismissed': r.dismissedAt, 'mb-2': !r.remindAt }"
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
                v-if="r.payload.link"
                :to="recordViewer(r.payload.link)"
              >
                <font-awesome-icon
                  :icon="['far', 'file-alt']"
                  class="text-primary mr-2"
                />
                {{ $t('reminder.recordPageLink') }}
              </b-dropdown-item>

              <b-dropdown-item
                @click.stop="$emit('edit', r)"
              >
                <font-awesome-icon
                  :icon="['far', 'edit']"
                  class="text-primary mr-2"
                />
                {{ $t('reminder.edit.label') }}
              </b-dropdown-item>

              <b-dropdown-divider />

              <c-input-confirm
                :text="$t('reminder.delete')"
                show-icon
                borderless
                variant="link"
                size="md"
                button-class="dropdown-item"
                icon-class="text-danger"
                class="w-100"
                @confirmed="$emit('delete', r)"
              />
            </b-dropdown>
          </div>

          <div class="reminder-item-content">
            <div class="d-flex flex-row flex-nowrap align-items-center mb-2">
              <b-form-checkbox
                data-test-id="checkbox-dismiss-reminder"
                :checked="!!r.dismissedAt"
                @change="$emit('dismiss', r, $event)"
              />
              <h6
                data-test-id="span-reminder-title"
                class="text-break text-truncate mb-0"
                :style="`${!!r.dismissedAt ? 'text-decoration: line-through;' : ''}`"
              >
                {{ r.payload.title || r.link || rlLabel(r) || r.linkLabel }}
              </h6>
            </div>

            <div
              v-if="r.payload.notes"
              class="text-secondary mb-1 text-break"
            >
              {{ r.payload.notes }}
            </div>
          </div>
        </div>

        <div
          v-if="r.remindAt"
          class="d-flex justify-content-end mt-2"
        >
          <div
            data-test-id="icon-remind-at"
            class="text-muted small cursor-pointer"
          >
            <font-awesome-icon
              :icon="['far', 'bell']"
              class="text-primary mr-1"
            />
            {{ r.remindAt | locFullDateTime }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { fmt } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  props: {
    reminders: {
      type: Array,
      required: true,
      default: () => [],
    },
  },

  computed: {
    sortedReminders () {
      return [...this.reminders].sort(this.stdSort)
    },
  },

  methods: {
    // Determine abs. link for given router-link
    rlLabel (r) {
      const rl = r.routerLink
      if (!rl) {
        return
      }
      return `${document.location.origin}${this.$router.resolve(rl).href}`
    },

    stdSort (a, b) {
      if (!a.dismissedAt) {
        return -1
      }
      if (!b.dismissedAt) {
        return 0
      }

      return b.dismissedAt - a.dismissedAt
    },

    makeTooltip ({ remindAt }) {
      return fmt.fullDateTime(remindAt)
    },

    recordViewer ({ params } = {}) {
      return params ? { name: 'page.record', params } : undefined
    },
  },
}
</script>

<style lang="scss" scoped>
.reminder-item-container {
  &:hover {
    background-color: var(--light) !important;
  }

  .reminder-item {
    transition: background-color 0.2s ease;

    &.dismissed {
      .reminder-item-content {
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
