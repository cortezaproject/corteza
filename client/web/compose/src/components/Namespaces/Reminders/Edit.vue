<template>
  <div>
    <b-list-group-item class="flex-column align-items-start px-2 py-2 border-0">
      <b-form
        class="import-form"
        @submit.prevent
      >
        <b-form-group :label="$t('reminder.edit.titleLabel')">
          <b-form-input
            v-model="title"
            data-test-id="input-title"
            required
            type="text"
            :placeholder="$t('reminder.edit.titlePlaceholder')"
          />
        </b-form-group>

        <b-form-group :label="$t('reminder.edit.notesLabel')">
          <b-form-textarea
            v-model="notes"
            data-test-id="textarea-notes"
            :placeholder="$t('reminder.edit.notesPlaceholder')"
            rows="6"
            max-rows="10"
          />
        </b-form-group>

        <b-form-group :label="$t('reminder.edit.remindAtLabel')">
          <b-form-select
            v-model="remindAt"
            data-test-id="select-remind-at"
            :options="remindAtPresets"
          />
        </b-form-group>

        <b-form-group :label="$t('reminder.edit.assigneeLabel')">
          <vue-select
            v-model="assignedTo"
            data-test-id="select-assignee"
            :options="assignees"
            option-value="userID"
            option-text="label"
            :placeholder="$t('field.kind.user.suggestionPlaceholder')"
            class="bg-white"
            @search="searchAssignees"
          />
        </b-form-group>

        <b-form-group
          v-if="reminder.payload.link"
          :label="$t('reminder.routesTo')"
        >
          <b-form-input
            v-model="reminder.payload.link.label"
            data-test-id="input-link"
          />
        </b-form-group>
      </b-form>
    </b-list-group-item>
    <div class="position-sticky text-center bg-white py-1 fixed-bottom">
      <b-button
        data-test-id="button-save"
        variant="outline-primary"
        class="px-2"
        @click="save"
      >
        {{ $t('label.save') }}
      </b-button>
    </div>
  </div>
</template>

<script>
import _ from 'lodash'
import moment from 'moment'
import { VueSelect } from 'vue-select'
import { system } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    VueSelect,
  },
  props: {
    edit: {
      type: Object,
      required: false,
      default: () => ({}),
    },

    myID: {
      type: String,
      required: true,
    },

    users: {
      type: Array,
      required: true,
      default: () => [],
    },
  },

  data () {
    return {
      // Do this, so we don't edit the original object
      reminder: {},
      assignees: [{ userID: this.myID, label: this.$t('reminder.edit.assigneePlaceholder') }],
    }
  },

  computed: {
    title: {
      get: function () {
        return (this.reminder.payload || {}).title
      },
      set: function (v) {
        this.updPayload('title', v)
      },
    },

    notes: {
      get: function () {
        return (this.reminder.payload || {}).notes
      },
      set: function (v) {
        this.updPayload('notes', v)
      },
    },

    remindAt: {
      get: function () {
        return (this.reminder.payload || {}).remindAt || null
      },
      set: function (v) {
        this.updPayload('remindAt', v)
      },
    },

    assignedTo: {
      get: function () {
        return this.assignees.find(({ userID }) => userID === this.reminder.assignedTo)
      },
      set: function (user) {
        let userID
        if (user) {
          userID = user.userID
        }
        this.$set(this.reminder, 'assignedTo', userID)
      },
    },

    remindAtPresets () {
      return [
        { value: null, text: this.$t('reminder.edit.remindAtNone') },
        { value: 1000 * 60 * 1, text: this.$t('label.timeMinute', { t: 1 }) },
        { value: 1000 * 60 * 5, text: this.$t('label.timeMinute', { t: 5 }) },
        { value: 1000 * 60 * 15, text: this.$t('label.timeMinute', { t: 15 }) },
        { value: 1000 * 60 * 30, text: this.$t('label.timeMinute', { t: 30 }) },
        { value: 1000 * 60 * 60 * 1, text: this.$t('label.timeHour', { t: 1 }) },
        { value: 1000 * 60 * 60 * 2, text: this.$t('label.timeHour', { t: 2 }) },
        { value: 1000 * 60 * 60 * 24, text: this.$t('label.timeHour', { t: 24 }) },
      ]
    },
  },

  watch: {
    edit: {
      handler: function () {
        this.reminder = new system.Reminder(this.edit)
      },
      deep: true,
      immediate: true,
    },
  },

  methods: {
    save () {
      // @todo support for updating times
      let r = {}
      if (this.remindAt) {
        r.remindAt = moment().add(this.remindAt, 'ms').toISOString()
      }

      r = {
        ...this.reminder,
        ...r,
      }

      this.$emit('save', r)
    },

    // Helper to handle undefined fields
    updPayload (f, v) {
      if (!this.reminder.payload) {
        this.$set(this.reminder, 'payload', {})
      }
      this.$set(this.reminder.payload, f, v)
    },

    searchAssignees (query) {
      if (query) {
        this.debouncedSearch(this, query)
      }
    },

    debouncedSearch: _.debounce((vm, query) => {
      vm.$SystemAPI.userList({ query }).then(({ set }) => {
        vm.assignees = set.map(({ userID, name: label }) => {
          if (userID === vm.myID) {
            return { userID, label: vm.$t('reminder.edit.assigneePlaceholder') }
          }
          return { userID, label }
        })
      })
    }, 300),
  },
}
</script>
