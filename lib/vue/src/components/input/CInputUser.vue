<template>
  <c-input-select
    v-model="user.value"
    data-test-id="select-user"
    :options="user.options"
    :get-option-label="getOptionLabel"
    :get-option-key="getOptionKey"
    :placeholder="placeholder"
    :loading="processing"
    :filterable="false"
    v-bind="$attrs"
    @search="search"
    @input="onUserUpdate"
  />
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { debounce } from 'lodash'

export default {
  name: 'CInputUser',

  props: {
    value: {
      type: String,
      default: null,
    },

    placeholder: {
      type: String,
      default: '',
    },
  },

  data () {
    return {
      processing: false,

      user: {
        options: [],
        value: undefined,

        filter: {
          query: null,
          limit: 20,
        },
      },
    }
  },

  created () {
    this.fetchUsers().then(() => {
      this.getUserByID(this.value)
    })
  },

  methods: {
    search: debounce(function (query) {
      if (query !== this.user.filter.query) {
        this.user.filter.query = query
        this.user.filter.page = 1
      }

      this.fetchUsers()
    }, 300),

    fetchUsers () {
      this.processing = true

      return this.$SystemAPI.userList(this.user.filter).then(({ set }) => {
        this.user.options = set.map(m => Object.freeze(m))
      }).finally(() => {
        setTimeout(() => {
          this.processing = false
        }, 500)
      })
    },

    getUserByID (userID) {
      if (!userID || userID === NoID) {
        this.user.value = undefined
        return
      }

      const user = this.user.options.find(o => o.userID === userID)

      if (user) {
        this.user.value = user
      } else {
        return this.$SystemAPI.userRead({ userID }).then(user => {
          this.user.value = user
          this.user.options.push(Object.freeze(user))
        })
      }
    },

    onUserUpdate (u) {
      this.$emit('input', u.userID)
      this.$emit('input-object', u)
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
