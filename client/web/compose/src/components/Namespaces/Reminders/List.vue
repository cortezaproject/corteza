<template>
  <div
    class="h-100 pl-1"
  >
    <div
      class="text-center bg-white py-2 sticky-top"
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

    <div
      v-for="r in sortedReminders"
      :key="r.reminderID"
      class="d-flex align-items-center mb-1 overflow-auto"
    >
      <b-form-checkbox
        data-test-id="checkbox-dismiss-reminder"
        :checked="!!r.dismissedAt"
        :disabled="!!r.dismissedAt"
        @change="$emit('dismiss', r)"
      >
        <span
          data-test-id="span-reminder-title"
          class="text-break"
          :style="`${!!r.dismissedAt ? 'text-decoration: line-through;' : ''}`"
        >
          {{ r.payload.title || r.link || rlLabel(r) || r.linkLabel }}
        </span>
      </b-form-checkbox>

      <div
        class="ml-auto"
      >
        <font-awesome-icon
          v-if="r.snoozeCount"
          data-test-id="icon-snoozed-reminder"
          :icon="['far', 'clock']"
          class="ml-1"
        />

        <font-awesome-icon
          v-else-if="r.remindAt"
          v-b-tooltip.hover
          data-test-id="icon-remind-at"
          :title="makeTooltip(r)"
          :icon="['far', 'bell']"
          class="ml-1"
        />

        <b-button
          data-test-id="button-edit-reminder"
          variant="link"
          class="p-1 ml-2"
          @click="$emit('edit', r)"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
            class="text-primary"
          />
        </b-button>

        <b-button
          data-test-id="button-delete-reminder"
          variant="link"
          class="text-dark p-1"
          @click.prevent="$emit('delete', r)"
        >
          <font-awesome-icon
            :icon="['far', 'trash-alt']"
          />
        </b-button>
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
      if (!a.remindAt) {
        return -1
      }
      if (!b.remindAt) {
        return 0
      }

      return a.remindAt - b.remindAt
    },

    makeTooltip ({ remindAt }) {
      return fmt.fullDateTime(remindAt)
    },
  },
}
</script>
