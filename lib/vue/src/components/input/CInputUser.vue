<template>
  <c-input-select
    ref="userSelect"
    :value="user.value"
    data-test-id="select-user"
    :options="user.options"
    :get-option-label="getOptionLabel"
    :get-option-key="getOptionKey"
    :placeholder="placeholder"
    :loading="processing"
    :filterable="false"
    :clearable="clearable"
    v-bind="$attrs"
    @search="search"
    @input="onUserUpdate"
  />
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { debounce } from 'lodash'
import axios from 'axios'

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

    clearable: {
      type: Boolean,
      default: false,
    },

    clearOnSelect: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      processing: false,
      cancelRequest: null,

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
      }

      this.fetchUsers()
    }, 300),

    fetchUsers () {
      this.processing = true

      if (this.cancelRequest) {
        this.cancelRequest()
        this.cancelRequest = null
      }

      const { response, cancel } = this.$SystemAPI.userListCancellable(this.user.filter)
      this.cancelRequest = cancel

      return Promise.all([response(), new Promise(resolve => setTimeout(resolve, 300))])
        .then(([{ set }]) => {
          this.user.options = set.map(m => Object.freeze(m))
          this.processing = false
        })
        .catch((e) => {
          if (axios.isCancel(e)) return
          this.processing = false
          throw e
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

    onUserUpdate (user) {
      if (this.clearOnSelect && this.$refs.userSelect) {
        this.$refs.userSelect._data._value = undefined
      } else {
        this.user.value = user
      }

      this.$emit('input', user.userID)
      this.$emit('input-object', user)
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
