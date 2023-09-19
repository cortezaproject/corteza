<template>
  <div class="d-flex flex-column h-100">
    <list
      v-if="!edit"
      :reminders="reminders"
      class="flex-fill"
      @edit="onEdit"
      @dismiss="onDismiss"
      @delete="onDelete"
    />

    <edit
      v-else
      :edit="edit"
      :users="users"
      :disable-save="disableSave"
      :processing-save="processingSave"
      class="flex-fill"
      @dismiss="onDismiss"
      @back="onCancel()"
      @save="onSave"
    />
  </div>
</template>

<script>
import List from './List'
import Edit from './Edit'
import { mapGetters } from 'vuex'
import { system, NoID } from '@cortezaproject/corteza-js'

export default {
  components: {
    List,
    Edit,
  },

  data () {
    return {
      reminders: [],
      edit: null,
      disableSave: false,
      processingSave: false,
    }
  },

  computed: {
    ...mapGetters({
      users: 'user/set',
    }),
  },

  mounted () {
    this.fetchReminders()
    // @todo remove this, when sockets get implemented
    this.$root.$on('reminders.pull', this.fetchReminders)
    this.$root.$on('reminder.updated', this.fetchReminders)
    this.$root.$on('reminder.create', this.onEdit)
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    onEdit ({ reminderID, resource, assignedTo, payload = {} } = {}) {
      if (reminderID) {
        this.edit = this.reminders.find(r => r.reminderID === reminderID)
      } else {
        this.edit = {
          resource: resource || `namespace:${this.namespaceID}`,
          assignedTo: assignedTo || this.$auth.user.userID,
          payload,
        }
      }

      this.$root.$emit('reminders.show')
    },

    onSave (r) {
      this.processingSave = true
      const endpoint = r.reminderID && r.reminderID !== NoID ? 'reminderUpdate' : 'reminderCreate'

      this.$SystemAPI[endpoint](r).then(() => {
        return this.fetchReminders()
      }).then(() => {
        this.onCancel()
        this.$Reminder.prefetch()
      }).finally(() => {
        this.processingSave = false
      })
    },

    onCancel () {
      this.edit = undefined
    },

    onDismiss ({ reminderID }, value) {
      const endpoint = value ? 'reminderDismiss' : 'reminderUndismiss'
      this.$SystemAPI[endpoint]({ reminderID }).then(() => {
        this.fetchReminders()
      })
    },

    onDelete ({ reminderID }) {
      this.$SystemAPI.reminderDelete({ reminderID }).then(() => {
        this.fetchReminders()
      })
    },

    async fetchReminders () {
      return this.$SystemAPI.reminderList({
        assignedTo: this.$auth.user.userID,
        limit: 0,
      }).then(({ set: reminders = [] }) => {
        this.reminders = reminders.map(r => new system.Reminder(r))
      })
    },

    setDefaultValues () {
      this.reminder = []
      this.edit = null
      this.disableSave = false
      this.processingSave = false
    },

    destroyEvents () {
      this.$root.$off('reminders.pull', this.fetchReminders)
      this.$root.$off('reminder.updated', this.fetchReminders)
      this.$root.$off('reminder.create', this.onEdit)
    },
  },

}
</script>
