<template>
  <div
    data-test-id="role-picker"
  >
    <vue-select
      ref="picker"
      data-test-id="input-role-picker"
      :options="filtered"
      :get-option-key="r => r.value"
      :get-option-label="r => getRoleLabel(r)"
      :calculate-position="calculateDropdownPosition"
      :placeholder="$t('admin:picker.role.placeholder')"
      class="bg-white w-100"
      multiple
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
          <b-button
            data-test-id="button-remove-role"
            variant="link"
            @click="removeRole(role)"
          >
            <font-awesome-icon :icon="['far', 'trash-alt']" />
          </b-button>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import { debounce } from 'lodash'
import { VueSelect } from 'vue-select'

function roleSorter (a, b) {
  return `${a.name} ${a.handle} ${a.roleID}`.localeCompare(`${b.name} ${b.handle} ${b.roleID}`)
}

export default {
  components: {
    VueSelect,
  },

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

    updateValue (role, index = -1) {
      // reset picker value for better value presentation
      if (this.$refs.picker) {
        this.$refs.picker._data._value = undefined
      }

      if (role[0]) {
        this.addRole(role[0])
      } else {
        if (index >= 0) {
          this.value.splice(index, 1)
        }
      }
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
      background-color: $light;
    }
  }
}

</style>
