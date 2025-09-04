<template>
  <c-input-select
    v-model="userGroup.value"
    data-test-id="select-user-group"
    :options="userGroup.options"
    :get-option-label="getOptionLabel"
    :get-option-key="getOptionKey"
    :placeholder="placeholder"
    :loading="processing"
    :filterable="false"
    v-bind="$attrs"
    @search="search"
    @input="onUserGroupUpdate"
  />
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { debounce } from 'lodash'

export default {
  name: 'CInputUserGroup',

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

      userGroup: {
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
    this.fetchUserGroups().then(() => {
      this.getUserGroupByID(this.value)
    })
  },

  methods: {
    search: debounce(function (query) {
      if (query !== this.userGroup.filter.query) {
        this.userGroup.filter.query = query
        this.userGroup.filter.page = 1
      }

      this.fetchUserGroups()
    }, 300),

    fetchUserGroups () {
      this.processing = true

      return this.$SystemAPI.userGroupList(this.userGroup.filter).then(({ set }) => {
        this.userGroup.options = set.map(m => Object.freeze(m))
      }).finally(() => {
        setTimeout(() => {
          this.processing = false
        }, 500)
      })
    },

    getUserGroupByID (userGroupID) {
      if (!userGroupID || userGroupID === NoID) {
        this.userGroup.value = undefined
        return
      }

      const userGroup = this.userGroup.options.find(o => o.userGroupID === userGroupID)

      if (userGroup) {
        this.userGroup.value = userGroup
      } else {
        return this.$SystemAPI.userGroupRead({ userGroupID }).then(userGroup => {
          this.userGroup.value = userGroup
          this.userGroup.options.push(Object.freeze(userGroup))
        })
      }
    },

    onUserGroupUpdate ({ userGroupID }) {
      this.$emit('input', userGroupID)
    },

    getOptionKey ({ userGroupID }) {
      return userGroupID
    },

    getOptionLabel ({ handle, meta, userGroupID }) {
      return handle || meta.short || userGroupID
    },
  },
}
</script>
