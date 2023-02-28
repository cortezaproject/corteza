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
      :my-i-d="$auth.user.userID"
      :users="users"
      class="flex-fill"
      @dismiss="onDismiss"
      @back="edit = undefined"
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
    }
  },

  computed: {
    ...mapGetters({
      users: 'user/set',
    }),
  },

  created () {
    this.fetchReminders()
    // @todo remove this, when sockets get implemented
    this.$root.$on('reminders.pull', this.fetchReminders)
    this.$root.$on('reminder.updated', this.fetchReminders)
    this.$root.$on('reminder.create', this.onEdit)
  },

  beforeDestroy () {
    this.$root.$off('reminders.pull', this.fetchReminders)
    this.$root.$off('reminder.updated', this.fetchReminders)
    this.$root.$off('reminder.create', this.onEdit)
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
      const endpoint = r.reminderID && r.reminderID !== NoID ? 'reminderUpdate' : 'reminderCreate'
      this.$SystemAPI[endpoint](r).then(r => {
        this.edit = undefined
        this.fetchReminders()
        this.$Reminder.prefetch()
      })

      this.onCancel()
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

    fetchReminders () {
      this.$SystemAPI.reminderList({
        assignedTo: this.$auth.user.userID,
        limit: 0,
      }).then(({ set: reminders = [] }) => {
        this.reminders = reminders.map(r => new system.Reminder(r))
      })
    },
  },

}
</script>
