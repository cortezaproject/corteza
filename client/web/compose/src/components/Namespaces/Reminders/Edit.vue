<template>
  <div
    class="d-flex flex-column w-100"
  >
    <b-form
      v-if="reminder"
      class="p-2 text-primary"
      @submit.prevent
    >
      <b-form-group v-if="reminder.reminderID !== '0'">
        <b-form-checkbox
          :checked="!!reminder.dismissedAt"
          switch
          class="mt-2"
          @change="$emit('dismiss', reminder, $event)"
        >
          {{ $t('reminder.dismissed') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group :label="$t('reminder.edit.titleLabel')">
        <b-form-input
          v-model="reminder.payload.title"
          data-test-id="input-title"
          required
          :placeholder="$t('reminder.edit.titlePlaceholder')"
        />
      </b-form-group>

      <b-form-group :label="$t('reminder.edit.notesLabel')">
        <b-form-textarea
          v-model="reminder.payload.notes"
          data-test-id="textarea-notes"
          :placeholder="$t('reminder.edit.notesPlaceholder')"
          rows="6"
          max-rows="10"
        />
      </b-form-group>

      <b-form-group :label="$t('reminder.edit.remindAtLabel')">
        <c-input-date-time
          v-model="reminder.remindAt"
          data-test-id="select-remind-at"
          :labels="{
            clear: $t('label.clear'),
            none: $t('label.none'),
            now: $t('label.now'),
            today: $t('label.today'),
          }"
        />
      </b-form-group>

      <b-form-group :label="$t('reminder.edit.assigneeLabel')">
        <vue-select
          v-model="reminder.assignedTo"
          data-test-id="select-assignee"
          :options="assignees"
          :get-option-label="getUserLabel"
          :get-option-key="getOptionKey"
          :loading="processingUsers"
          :placeholder="$t('field.kind.user.suggestionPlaceholder')"
          :calculate-position="calculateDropdownPosition"
          :reduce="user => user.userID"
          option-value="userID"
          class="bg-white"
          @search="searchUsers"
        />
      </b-form-group>

      <b-form-group
        v-if="reminder.payload.link"
        :label="$t('reminder.routesTo')"
      >
        <b-input-group>
          <b-form-input
            v-model="reminder.payload.link.label"
            data-test-id="input-link"
          />

          <b-input-group-append>
            <b-button
              :disabled="!recordViewer"
              :to="recordViewer"
              :title="$t('reminder.recordPageLink')"
              variant="light"
              class="d-flex align-items-center text-primary"
            >
              <font-awesome-icon :icon="['far', 'file-alt']" />
            </b-button>
          </b-input-group-append>
        </b-input-group>
      </b-form-group>

      <div class="text-center py-1">
        <b-button
          data-test-id="button-save"
          variant="outline-primary"
          @click="$emit('save', reminder)"
        >
          {{ $t('label.save') }}
        </b-button>
      </div>
    </b-form>

    <b-button
      variant="outline-light"
      class="text-primary mt-auto border-0"
      @click="$emit('back')"
    >
      {{ $t('label.back') }}
    </b-button>
  </div>
</template>

<script>
import _ from 'lodash'
import { VueSelect } from 'vue-select'
import { system } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
const { CInputDateTime } = components

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    VueSelect,
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
