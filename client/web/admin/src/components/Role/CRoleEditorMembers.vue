<template>
  <b-card
    data-test-id="card-role-edit-members"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-form
      @submit.prevent="$emit('submit')"
    >
      <c-input-select
        ref="picker"
        data-test-id="input-role-members"
        :options="options"
        :get-option-key="u => u.value"
        :get-option-label="u => getUserLabel(u)"
        :placeholder="$t('admin:picker.member.placeholder')"
        :filterable="false"
        @search="search"
        @input="updateValue($event)"
      />
      <b-table-simple
        v-if="memberUsers && users"
        responsive
        small
        hover
        class="w-100 p-0 mb-0 mt-1"
      >
        <tbody>
          <tr
            v-for="user in memberUsers"
            :key="user.userID"
          >
            <td class="align-middle">
              {{ getUserLabel(user) }}
            </td>
            <td class="text-right">
              <c-input-confirm
                data-test-id="button-remove-member"
                show-icon
                @confirmed="removeMember(user.userID)"
              />
            </td>
          </tr>
        </tbody>
      </b-table-simple>
    </b-form>

    <template #footer>
      <c-button-submit
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit')"
      />
    </template>
  </b-card>
</template>

<script>
import { debounce } from 'lodash'

export default {
  i18nOptions: {
    namespaces: 'system.roles',
    keyPrefix: 'editor.members',
  },

  props: {
    currentMembers: {
      type: Array,
      required: true,
      default: () => [],
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },
  },

  data () {
    return {
      users: [],
      filter: '',
      memberUsers: [],
    }
  },

  computed: {
    members: {
      get () {
        return this.currentMembers
      },

      set (members) {
        this.$emit('update:current-members', members)
      },
    },

    options: {
      get () {
        const memberIDs = this.currentMembers.map(m => m.userID)
        return this.users
          .filter(u => !this.isMember(u) && !memberIDs.includes(u.userID))
      },
    },
  },

  mounted () {
    this.fetchMembers()
    this.fetchUsers()
  },

  methods: {
    memberIndex (u) {
      return this.members.findIndex(({ userID }) => userID === u)
    },

    isMember (u) {
      return this.members.findIndex(({ userID, dirty }) => userID === u && dirty) >= 0
    },

    addMember (member) {
      if (!this.currentMembers.includes(member.userID)) {
        const label = this.getUserLabel(member)
        const i = this.memberIndex(member.userID)
        if (i < 0) {
          this.members.push({ userID: member.userID, label: label, current: false, dirty: true })
        } else {
          this.$set(this.members, i, { ...this.members[i], label: label, dirty: true })
        }

        this.memberUsers.push(member)
      }
    },

    removeMember (userID) {
      const i = this.memberIndex(userID)
      if (i > -1) {
        this.$set(this.members, i, { ...this.members[i], dirty: false })
      }

      this.memberUsers = this.memberUsers.filter(({ userID: uID }) => uID !== userID)
      const removedMember = this.users
        .filter(u => u.userID === userID)
        .map(m => { return { value: m.userID, label: this.getUserLabel(m) } })
      this.options.push(...removedMember)
    },

    fetchMembers () {
      const userID = this.members.map(({ userID }) => userID)

      if (userID.length === 0) {
        return
      }

      this.$SystemAPI.userList({ userID, deleted: 1, suspended: 1 }).then(({ set: users = [] }) => {
        this.memberUsers = users
      })
    },

    fetchUsers (query = this.filter) {
      this.$SystemAPI.userList({ query, limit: 25 })
        .then(({ set: items = [] }) => {
          this.users = items
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.fetch.error')))
    },

    search: debounce(function (query = '') {
      if (query !== this.filter) {
        this.filter = query
      }

      this.fetchUsers()
    }, 300),

    updateValue (user) {
      // reset picker value for better value presentation
      if (this.$refs.picker) {
        this.$refs.picker._data._value = undefined
      }

      this.addMember(user)
    },

    getUserLabel ({ label, name, handle, username, email }) {
      return label || name || handle || username || email || this.$t('unnamed')
    },
  },
}
</script>
