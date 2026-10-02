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
    :clearable="false"
    v-bind="$attrs"
    @search="search"
    @input="onUserGroupUpdate"
  />
</template>

<script>
import { NoID, system } from '@cortezaproject/corteza-js'
import { debounce } from 'lodash'
import axios from 'axios'

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
      cancelRequest: null,

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

    fetchUserGroups() {
      this.processing = true

      if (this.cancelRequest) {
        this.cancelRequest()
        this.cancelRequest = null
      }

      const { response, cancel } = this.$SystemAPI.userGroupListCancellable(this.userGroup.filter)
      this.cancelRequest = cancel

      return Promise.all([response(), new Promise(resolve => setTimeout(resolve, 300))])
        .then(([{ set }]) => {
          this.userGroup.options = set.map((m) => new system.UserGroup(m))
          this.processing = false
        })
        .catch((e) => {
          if (axios.isCancel(e)) return
          this.processing = false
          throw e
        })
    },

    getUserGroupByID (userGroupID) {
      if (!userGroupID || userGroupID === NoID) {
        this.userGroup.value = this.userGroup.options.find(({ isRoot }) => !!isRoot)

        this.$emit('input', this.userGroup.value?.userGroupID)
        return
      }

      const userGroup = this.userGroup.options.find(o => o.userGroupID === userGroupID)

      if (userGroup) {
        this.userGroup.value = userGroup
      }
    },

    onUserGroupUpdate ({ userGroupID }) {
      this.$emit('input', userGroupID)
    },

    getOptionKey ({ userGroupID }) {
      return userGroupID
    },

    getOptionLabel ({ handle, meta = {}, userGroupID }) {
      return meta.short || handle || userGroupID
    },
  },
}
</script>
