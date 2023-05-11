<template>
  <b-card
    data-test-id="card-role-edit-members"
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="$emit('submit')"
    >
      <b-form-group
        :label="$t('count', { count: members.filter(({ dirty }) => dirty).length })"
        class="mb-0"
      >
        <vue-select
          ref="picker"
          data-test-id="input-role-members"
          :options="options"
          :get-option-key="u => u.value"
          :get-option-label="u => getUserLabel(u)"
          :calculate-position="calculateDropdownPosition"
          :placeholder="$t('admin:picker.member.placeholder')"
          class="bg-white w-100"
          multiple
          @search="search"
          @input="updateValue($event)"
        />
        <table
          v-if="memberUsers && users"
          class="w-100 p-0 table-hover mb-2"
        >
          <tbody>
            <tr
              v-for="user in memberUsers"
              :key="user.userID"
            >
              <td>{{ getUserLabel(user) }}</td>
              <td class="text-right">
                <b-button
                  data-test-id="button-remove-member"
                  variant="link"
                  class="text-danger pr-0"
                  @click="removeMember(user.userID)"
                >
                  {{ $t('remove') }}
                </b-button>
              </td>
            </tr>
          </tbody>
        </table>
      </b-form-group>
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-submit-button
        class="float-right"
        :processing="processing"
        :success="success"
        @submit="$emit('submit')"
      />
    </template>
  </b-card>
</template>

<script>
import { debounce } from 'lodash'
import { VueSelect } from 'vue-select'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  i18nOptions: {
    namespaces: 'system.roles',
    keyPrefix: 'editor.members',
  },

  components: {
    CSubmitButton,
    VueSelect,
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
    this.fetchUsers()
    this.fetchMembers()
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

        this.memberUsers.push({ value: member.userID, label })
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

    fetchUsers () {
      this.$SystemAPI.userList({ query: this.filter })
        .then(({ set: items = [] }) => {
          this.users = items
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.fetch.error')))
    },

    fetchMembers () {
      const userID = this.members.map(({ userID }) => userID)
      if (userID.length > 0) {
        this.$SystemAPI.userList({ query: this.filter, userID })
          .then(({ set: items = [] }) => {
            this.memberUsers = items
          })
          .catch(this.toastErrorHandler(this.$t('notification:user.fetch.error')))
      }
    },

    search: debounce(function (query = '') {
      if (query !== this.filter) {
        this.filter = query
      }

      this.fetchUsers()
    }, 300),

    updateValue (user, index = -1) {
      // reset picker value for better value presentation
      if (this.$refs.picker) {
        this.$refs.picker._data._value = undefined
      }

      if (user[0]) {
        this.addMember(user[0])
      } else {
        if (index >= 0) {
          this.value.splice(index, 1)
        }
      }
    },

    getUserLabel ({ label, name, handle, username, email }) {
      return label || name || handle || username || email || this.$t('unnamed')
    },
  },
}
</script>
