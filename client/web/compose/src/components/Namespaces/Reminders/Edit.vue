<template>
  <div
    class="d-flex flex-column h-100"
  >
    <b-form
      v-if="reminder"
      class="flex-fill overflow-auto p-2 text-primary"
      @submit.prevent
    >
      <b-form-group
        v-if="reminder.reminderID !== '0'"
      >
        <b-form-checkbox
          :checked="!!reminder.dismissedAt"
          class="mt-2"
          @change="$emit('dismiss', reminder, $event)"
        >
          {{ $t('reminder.dismissed') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        :label="$t('reminder.edit.titleLabel')"
        label-class="text-primary"
      >
        <b-form-input
          v-model="reminder.payload.title"
          data-test-id="input-title"
          required
          :placeholder="$t('reminder.edit.titlePlaceholder')"
        />
      </b-form-group>

      <b-form-group
        :label="$t('reminder.edit.notesLabel')"
        label-class="text-primary"
      >
        <b-form-textarea
          v-model="reminder.payload.notes"
          data-test-id="textarea-notes"
          :placeholder="$t('reminder.edit.notesPlaceholder')"
          rows="6"
          max-rows="10"
        />
      </b-form-group>

      <b-form-group
        :label="$t('reminder.edit.remindAtLabel')"
        label-class="text-primary"
      >
        <c-input-date-time
          v-model="reminder.remindAt"
          data-test-id="select-remind-at"
          only-future
          :labels="{
            clear: $t('label.clear'),
            none: $t('label.none'),
            now: $t('label.now'),
            today: $t('label.today'),
          }"
        />
      </b-form-group>

      <b-form-group
        :label="$t('reminder.edit.assigneeLabel')"
        label-class="text-primary"
      >
        <c-input-select
          v-model="reminder.assignedTo"
          data-test-id="select-assignee"
          :options="assignees"
          :get-option-label="getUserLabel"
          :get-option-key="getOptionKey"
          :loading="processingUsers"
          :placeholder="$t('field.kind.user.suggestionPlaceholder')"
          :reduce="user => user.userID"
          option-value="userID"
          @search="searchUsers"
        />
      </b-form-group>

      <b-form-group
        v-if="reminder.payload.link"
        :label="$t('reminder.routesTo')"
        label-class="text-primary"
      >
        <b-input-group>
          <b-form-input
            v-model="reminder.payload.link.label"
            data-test-id="input-link"
          />

          <b-input-group-append>
            <b-button
              v-b-tooltip.hover="{ title: $t('reminder.recordPageLink'), container: '#body' }"
              :disabled="!recordViewer"
              :to="recordViewer"
              variant="light"
              class="d-flex align-items-center text-primary"
            >
              <font-awesome-icon :icon="['far', 'file-alt']" />
            </b-button>
          </b-input-group-append>
        </b-input-group>
      </b-form-group>

      <b-form-group
        v-if="reminder.dismissedAt"
        :label="$t('reminder.dismissedAt')"
        label-class="text-primary"
      >
        {{ reminder.dismissedAt | locFullDateTime }}
      </b-form-group>

      <b-form-group
        v-if="reminder.snoozeCount"
        :label="$t('reminder.snooze.count')"
        label-class="text-primary"
      >
        {{ reminder.snoozeCount }}
      </b-form-group>
    </b-form>

    <div class="d-flex align-items-center justify-content-around py-3">
      <b-button
        data-test-id="button-back"
        variant="outline-light"
        class="text-primary border-0"
        @click="$emit('back')"
      >
        <font-awesome-icon
          :icon="['fas', 'chevron-left']"
          class="back-icon"
        />
        {{ $t('label.back') }}
      </b-button>

      <c-button-submit
        data-test-id="button-save"
        :disabled="disableSave"
        :processing="processingSave"
        :text="$t('label.save')"
        @submit="$emit('save', reminder)"
      />
    </div>
  </div>
</template>

<script>
import _ from 'lodash'
import { system } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
const { CInputDateTime } = components

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    CInputDateTime,
  },

  props: {
    edit: {
      type: Object,
      required: false,
      default: () => ({}),
    },

    users: {
      type: Array,
      required: true,
      default: () => [],
    },

    disableSave: {
      type: Boolean,
      default: false,
    },

    processingSave: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      processingUsers: false,

      // Do this, so we don't edit the original object
      reminder: undefined,
      assignees: [{ userID: this.$auth.user.userID }],
    }
  },

  computed: {
    recordViewer () {
      const { params } = this.reminder.payload.link || {}
      return params ? { name: 'page.record', params } : undefined
    },
  },

  watch: {
    edit: {
      immediate: true,
      deep: true,
      handler (edit) {
        this.reminder = new system.Reminder(edit)
        this.searchUsers()
      },
    },
  },

  methods: {
    searchUsers: _.debounce(function (query) {
      this.processingUsers = true

      this.$SystemAPI.userList({ query, limit: 15 }).then(({ set = [] }) => {
        this.assignees = set
      }).finally(() => {
        this.processingUsers = false
      })
    }, 300),

    getUserLabel ({ userID, email, name, username }) {
      if (userID === this.$auth.user.userID) {
        return this.$t('reminder.edit.assigneePlaceholder')
      }

      return name || username || email || `<@${userID}>`
    },

    getOptionKey ({ userID }) {
      return userID
    },
  },
}
</script>
