<template>
  <div
    data-test-id="role-picker"
  >
    <c-input-select
      ref="picker"
      data-test-id="input-role-picker"
      :options="filtered"
      :get-option-key="r => r.value"
      :get-option-label="r => getRoleLabel(r)"
      :placeholder="$t('admin:picker.role.placeholder')"
      @search="search"
      @input="updateValue($event)"
    />
    <b-form-text
      v-if="$slots['description']"
    >
      <slot name="description" />
    </b-form-text>

    <b-container
      v-if="selected"
      class="p-1"
    >
      <b-row
        v-for="role in selected"
        :key="role.roleID"
        data-test-id="selected-row-list"
      >
        <b-col>{{ getRoleLabel(role) }}</b-col>
        <b-col class="text-right">
          <c-input-confirm
            data-test-id="button-remove-role"
            show-icon
            @confirmed="removeRole(role)"
          />
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import { debounce } from 'lodash'

function roleSorter (a, b) {
  return `${a.name} ${a.handle} ${a.roleID}`.localeCompare(`${b.name} ${b.handle} ${b.roleID}`)
}

export default {
  props: {
    label: {
      type: String,
      default: 'count',
    },

    // list of role IDs
    value: {
      type: Array,
      default: () => ([]),
    },
  },

  data () {
    return {
      roles: [],
      filter: '',
    }
  },

  computed: {
    selected () {
      return this.roles
        .filter(({ roleID }) => (this.value || []).includes(roleID))
        .sort(roleSorter)
    },

    filtered () {
      const match = ({ name = '', handle = '', roleID = '' }) => {
        return `${name} ${handle} ${roleID}`.toLocaleLowerCase().indexOf(this.filter.toLocaleLowerCase()) > -1
      }

      const fits = ({ isClosed, meta = {} }) => {
        return !(isClosed || (meta.context && meta.context.resourceTypes))
      }

      return this.roles.filter(r => !(this.value || []).includes(r.roleID) && fits(r) && match(r))
    },
  },

  mounted () {
    this.preload()
  },

  methods: {
    addRole (role) {
      if (!this.value.includes(role.roleID)) {
        this.value.push(role.roleID)
        this.$emit('input', this.value)
      }
    },

    removeRole (r) {
      this.value.splice(this.value.indexOf(r.roleID), 1)
      this.filter = ''
    },

    preload () {
      return this.$SystemAPI.roleList({ query: this.filter })
        .then(({ set }) => { this.roles = set || [] })
        .catch(this.toastErrorHandler(this.$t('notification:role.fetch.error')))
    },

    search: debounce(function (query = '') {
      if (query !== this.filter) {
        this.filter = query
      }

      this.preload()
    }, 300),

    updateValue (role) {
      // reset picker value for better value presentation
      if (this.$refs.picker) {
        this.$refs.picker._data._value = undefined
      }

      this.addRole(role)
    },

    getRoleLabel ({ name, handle, roleID }) {
      return name || handle || roleID
    },
  },
}
</script>
<style lang="scss">
.results {
  z-index: 100;
  .filtered-role {
    cursor: pointer;
    &:hover {
      background-color: var(--light);
    }
  }
}

</style>
