<template>
  <c-input-select
    ref="roleSelect"
    :value="value"
    :options="roles"
    :placeholder="placeholder"
    :get-option-key="r => r.roleID"
    :get-option-label="r => getRoleLabel(r)"
    :filterable="false"
    :selectable="selectable"
    :multiple="multiple"
    :clearable="clearable"
    :loading="loading"
    @search="search"
    @input="updateValue"
  />
</template>

<script>
import { debounce } from 'lodash'
import axios from 'axios'

export default {
  props: {
    value: {
      type: [Array, String, Object],
      default: '',
    },

    visible: {
      type: Function,
      default: () => true,
    },

    placeholder: {
      type: String,
      default: 'Start typing to search for roles',
    },

    multiple: {
      type: Boolean,
      default: false,
    },

    clearOnSelect: {
      type: Boolean,
      default: false,
    },

    selectable: {
      type: Function,
      default: () => true,
    },

    clearable: {
      type: Boolean,
      default: true,
    },

    preselect: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      loading: false,
      cancelRequest: null,

      roles: [],
      filter: '',
    }
  },

  mounted () {
    this.fetchRoles(this.preselect)
  },

  methods: {
    fetchRoles (preselect = false) {
      this.loading = true

      if (this.cancelRequest) {
        this.cancelRequest()
        this.cancelRequest = null
      }

      const { response, cancel } = this.$SystemAPI.roleListCancellable({ query: this.filter, limit: 20 })
      this.cancelRequest = cancel

      return Promise.all([response(), new Promise(resolve => setTimeout(resolve, 300))])
        .then(([{ set }]) => {
          this.roles = set.filter(this.visible)

          if (preselect && (!this.value || !this.value.length)) {
            this.updateValue(this.roles[0])
          }
          this.loading = false
        })
        .catch((e) => {
          if (axios.isCancel(e)) return
          this.loading = false
          throw e
        })
    },

    search: debounce(function (query = '') {
      if (query !== this.filter) {
        this.filter = query
      }

      this.fetchRoles()
    }, 400),

    updateValue (role) {
      // reset role-select value for better value presentation
      if (this.clearOnSelect && this.$refs.roleSelect) {
        this.$refs.roleSelect._data._value = undefined
      }

      this.$emit('input', role)
    },

    getRoleLabel ({ name, handle, roleID }) {
      return name || handle || roleID
    },
  },
}
</script>
