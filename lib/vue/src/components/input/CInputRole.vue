<template>
  <c-input-select
    ref="picker"
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

      return this.$SystemAPI.roleList({ query: this.filter, limit: 20 })
        .then(({ set }) => {
          this.roles = set.filter(this.visible)

          if (preselect && (!this.value || !this.value.length)) {
            this.updateValue(this.roles[0])
          }
        }).finally(() => {
          const timeout = this.filter ? 300 : 0

          setTimeout(() => {
            this.loading = false
          }, timeout)
        })
    },

    search: debounce(function (query = '') {
      if (query !== this.filter) {
        this.filter = query
      }

      this.fetchRoles()
    }, 400),

    updateValue (role) {
      // reset picker value for better value presentation
      if (this.$refs.picker && this.clearOnSelect) {
        this.$refs.picker._data._value = undefined
      }

      this.$emit('input', role)
    },

    getRoleLabel ({ name, handle, roleID }) {
      return name || handle || roleID
    },
  },
}
</script>
