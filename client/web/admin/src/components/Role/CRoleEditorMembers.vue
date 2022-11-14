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
        <table
          v-if="members && users"
          class="w-100 p-0 table-hover mb-2"
        >
          <tbody>
            <tr
              v-for="u in memberUsers"
              :key="u.userID"
            >
              <td>{{ u.name || u.handle || u.username || u.email || $t('unnamed') }}</td>
              <td class="text-right">
                <b-button
                  data-test-id="button-remove-member"
                  variant="link"
                  class="text-danger pr-0"
                  @click="removeMember(u)"
                >
                  {{ $t('remove') }}
                </b-button>
              </td>
            </tr>
          </tbody>
        </table>
        <b-input-group>
          <b-input-group-prepend>
            <b-input-group-text>{{ $t('searchUsers') }}</b-input-group-text>
          </b-input-group-prepend>
          <b-form-input
            v-model.trim="filter"
            data-test-id="input-search"
          />
        </b-input-group>
        <table
          v-if="filter && users"
          class="w-100 p-0 table-hover mt-2"
        >
          <tbody>
            <tr
              v-for="u in filtered"
              :key="u.userID"
            >
              <td>{{ u.name || u.handle || u.username || u.email || $t('unnamed') }}</td>
              <td class="text-right">
                <b-button
                  v-if="isMember(u)"
                  variant="link"
                  class="text-danger pr-0"
                  @click="removeMember(u)"
                >
                  {{ $t('remove') }}
                </b-button>
                <b-button
                  v-else
                  data-test-id="button-add-member"
                  variant="light"
                  @click="addMember(u)"
                >
                  {{ $t('add') }}
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
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  i18nOptions: {
    namespaces: 'system.roles',
    keyPrefix: 'editor.members',
  },

  components: {
    CSubmitButton,
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
      filter: '',
      users: [],

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

    filtered () {
      return this.users.filter(u => !this.isMember(u))
    },
  },

  watch: {
    filter: {
      handler () {
        const query = this.filter
        this.$SystemAPI.userList({ query })
          .then(({ set: items = [] }) => {
            this.users = items
          }).catch(this.toastErrorHandler(this.$t('notification:user.fetch.error')))
      },
    },
  },

  mounted () {
    const userID = this.members.map(({ userID }) => userID)
    if (userID.length > 0) {
      this.$SystemAPI.userList({ userID })
        .then(({ set: items = [] }) => {
          this.memberUsers = items
        }).catch(this.toastErrorHandler(this.$t('notification:user.fetch.error')))
    }
  },

  methods: {
    memberIndex (u) {
      u = typeof u === 'object' ? u.userID : u
      return this.members.findIndex(({ userID }) => userID === u)
    },

    isMember (u) {
      u = typeof u === 'object' ? u.userID : u
      return this.members.findIndex(({ userID, dirty }) => userID === u && dirty) >= 0
    },

    addMember (u) {
      const i = this.memberIndex(u)
      if (i < 0) {
        this.members.push({ userID: typeof u === 'object' ? u.userID : u, current: false, dirty: true })
      } else {
        this.$set(this.members, i, { ...this.members[i], dirty: true })
      }

      this.memberUsers.push(u)
    },

    removeMember (u) {
      const i = this.memberIndex(u)
      if (i > -1) {
        this.$set(this.members, i, { ...this.members[i], dirty: false })
      }
      u = typeof u === 'object' ? u.userID : u
      this.memberUsers = this.memberUsers.filter(({ userID }) => userID !== u)
    },
  },
}
</script>
