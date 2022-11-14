<template>
  <vue-select
    data-test-id="select-user"
    :options="user.options"
    :get-option-label="getOptionLabel"
    :get-option-key="getOptionKey"
    :value="user.value"
    class="bg-white"
    @search="search"
    @input="updateRunAs"
  />
</template>

<script>
import { debounce } from 'lodash'
import { VueSelect } from 'vue-select'

export default {
  components: {
    VueSelect,
  },

  props: {
    userID: {
      type: String,
      default: null,
    },

  },

  data () {
    return {
      user: {
        options: [],
        value: undefined,

        filter: {
          query: null,
          limit: 10,
        },
      },
    }
  },

  created () {
    this.fetchUsers()
    this.getUserByID()
  },

  methods: {
    search: debounce(function (query) {
      if (query !== this.user.filter.query) {
        this.user.filter.query = query
        this.user.filter.page = 1
      }

      if (query) {
        this.fetchUsers()
      }
    }, 300),

    fetchUsers () {
      this.$SystemAPI.userList(this.user.filter)
        .then(({ set }) => {
          this.user.options = set.map(m => Object.freeze(m))
        })
    },

    async getUserByID () {
      this.$SystemAPI.userRead({ userID: this.userID })
        .then(user => {
          this.user.value = user
        }).catch(() => {
          return {}
        })
    },

    updateRunAs (user) {
      if (user && user.userID) {
        this.user.value = user
      } else {
        this.user.value = null
      }
      this.$emit('updateUser', user)
    },

    getOptionKey ({ userID }) {
      return userID
    },

    getOptionLabel ({ userID, email, name, username }) {
      return name || username || email || `<@${userID}>`
    },
  },
}
</script>
