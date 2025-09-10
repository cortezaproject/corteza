<template>
  <div
    data-test-id="role-picker"
    class="d-flex flex-column"
  >
    <c-input-user
      data-test-id="input-user-picker"
      :selectable="r => !value.includes(r.userID)"
      :placeholder="$t('admin:picker.user.placeholder')"
      :visible="isUserVisible"
      clear-on-select
      @input-object="addUser($event)"
    />

    <b-spinner
      v-if="preloading"
      class="mx-auto my-4"
    />

    <b-table-simple
      v-else-if="getSelectedUsers.length"
      responsive
      small
      hover
      class="w-100 p-0 mb-0 mt-1"
    >
      <tbody>
        <tr
          v-for="user in getSelectedUsers"
          :key="user.userID"
          data-test-id="selected-row-list"
        >
          <td class="align-middle">
            {{ getUserLabel(user) }}
          </td>
          <td
            v-if="!noRemove"
            class="text-right"
          >
            <c-input-confirm
              data-test-id="button-remove-user"
              show-icon
              @confirmed="removeUser(user.userID)"
            />
          </td>
        </tr>
      </tbody>
    </b-table-simple>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CInputUser } = components

export default {
  components: {
    CInputUser,
  },

  props: {
    // list of user IDs
    value: {
      type: Array,
      default: () => ([]),
    },

    noRemove: {
      type: Boolean,
      required: false,
      default: false,
    },
  },

  data () {
    return {
      fetching: false,
      preloading: false,

      filter: '',

      selectedUsers: [],
    }
  },

  computed: {
    getSelectedUsers () {
      return this.selectedUsers.filter(({ userID }) => this.value.includes(userID))
    },
  },

  mounted () {
    this.preloadSelected()
  },

  methods: {
    addUser (user) {
      if (!this.value.includes(user.userID)) {
        this.selectedUsers.push(user)
        this.$emit('input', [...this.value, user.userID])
      }
    },

    removeUser (userID) {
      this.selectedUsers = this.selectedUsers.filter(({ userID: rID }) => rID !== userID)
      this.$emit('input', this.value.filter(v => v !== userID))
    },

    preloadSelected () {
      if (!this.value.length) {
        return
      }

      this.preloading = true

      return this.$SystemAPI.userList({ userID: this.value })
        .then(({ set }) => {
          this.selectedUsers = set || []
        })
        .finally(() => {
          this.preloading = false
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.fetch.error')))
    },

    getUserLabel ({ name, handle, userID }) {
      return name || handle || userID
    },

    isUserVisible ({ isClosed, meta = {} }) {
      return !(isClosed || (meta.context && meta.context.resourceTypes))
    },
  },
}
</script>

<style lang="scss">
.results {
  z-index: 100;
  .filtered-user {
    cursor: pointer;
    &:hover {
      background-color: var(--light);
    }
  }
}

</style>
