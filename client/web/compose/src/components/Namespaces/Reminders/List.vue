<template>
  <div
    class="h-100"
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
      v-for="(r, i) in sortedReminders"
      :key="r.reminderID"
    >
      <hr v-if="r.dismissedAt && sortedReminders[i - 1] ? !sortedReminders[i - 1].dismissedAt : false ">

      <div
        class="d-flex flex-row flex-nowrap align-items-center mb-2 overflow-auto border card"
      >
        <b-form-checkbox
          v-b-tooltip.hover.left.350
          data-test-id="checkbox-dismiss-reminder"
          :checked="!!r.dismissedAt"
          switch
          :title="$t(`reminder.${!!r.dismissedAt ? 'undismiss' : 'dismiss'}`)"
          class="my-2 ml-2"
          @change="$emit('dismiss', r, $event)"
        />

        <div
          data-test-id="span-reminder-title"
          class="text-break"
          :style="`${!!r.dismissedAt ? 'text-decoration: line-through;' : ''}`"
        >
          {{ r.payload.title || r.link || rlLabel(r) || r.linkLabel }}
        </div>

        <div
          class="d-flex align-items-center text-primary ml-auto px-2"
        >
          <font-awesome-icon
            v-if="r.snoozeCount"
            data-test-id="icon-snoozed-reminder"
            :title="makeTooltip(r)"
            :icon="['far', 'clock']"
            class="m-2"
          />

          <font-awesome-icon
            v-if="r.remindAt"
            data-test-id="icon-remind-at"
            :title="makeTooltip(r)"
            :icon="['far', 'bell']"
            class="m-2"
          />

          <b-button-group
            size="sm"
          >
            <b-button
              v-if="r.payload.link"
              :to="recordViewer(r.payload.link)"
              :title="$t('reminder.recordPageLink')"
              variant="outline-light"
              class="d-flex align-items-center py-2 text-primary border-0"
            >
              <font-awesome-icon :icon="['far', 'file-alt']" />
            </b-button>

            <b-button
              data-test-id="button-edit-reminder"
              variant="outline-light"
              :title="$t('reminder.edit.label')"
              class="d-flex align-items-center py-2 text-primary border-0"
              @click="$emit('edit', r)"
            >
              <font-awesome-icon :icon="['far', 'edit']" />
            </b-button>

            <b-button
              data-test-id="button-delete-reminder"
              variant="outline-light"
              :title="$t('reminder.delete')"
              class="d-flex align-items-center py-2 text-danger border-0"
              @click.prevent="$emit('delete', r)"
            >
              <font-awesome-icon
                :icon="['far', 'trash-alt']"
              />
            </b-button>
          </b-button-group>
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

      return a.dismissedAt - b.dismissedAt
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
