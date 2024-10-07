<template>
  <div
    data-test-id="role-picker"
    class="d-flex flex-column"
  >
    <c-input-role
      data-test-id="input-role-picker"
      :selectable="r => !value.includes(r.roleID)"
      :placeholder="$t('admin:picker.role.placeholder')"
      :visible="isRoleVisible"
      clear-on-select
      @input="addRole($event)"
    />

    <b-spinner
      v-if="preloading"
      class="mx-auto my-4"
    />

    <b-table-simple
      v-else-if="getSelectedRoles.length"
      responsive
      small
      hover
      class="w-100 p-0 mb-0 mt-1"
    >
      <tbody>
        <tr
          v-for="role in getSelectedRoles"
          :key="role.roleID"
          data-test-id="selected-row-list"
        >
          <td class="align-middle">
            {{ getRoleLabel(role) }}
          </td>
          <td class="text-right">
            <c-input-confirm
              data-test-id="button-remove-role"
              show-icon
              @confirmed="removeRole(role.roleID)"
            />
          </td>
        </tr>
      </tbody>
    </b-table-simple>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CInputRole } = components

export default {
  components: {
    CInputRole,
  },

  props: {
    // list of role IDs
    value: {
      type: Array,
      default: () => ([]),
    },
  },

  data () {
    return {
      fetching: false,
      preloading: false,

      filter: '',

      selectedRoles: [],
    }
  },

  computed: {
    getSelectedRoles () {
      return this.selectedRoles.filter(({ roleID }) => this.value.includes(roleID))
    },
  },

  mounted () {
    this.preloadSelected()
  },

  methods: {
    addRole (role) {
      if (!this.value.includes(role.roleID)) {
        this.selectedRoles.push(role)
        this.$emit('input', [...this.value, role.roleID])
      }
    },

    removeRole (roleID) {
      this.selectedRoles = this.selectedRoles.filter(({ roleID: rID }) => rID !== roleID)
      this.$emit('input', this.value.filter(v => v !== roleID))
    },

    preloadSelected () {
      if (!this.value.length) {
        return
      }

      this.preloading = true

      return this.$SystemAPI.roleList({ roleID: this.value })
        .then(({ set }) => {
          this.selectedRoles = set || []
        })
        .finally(() => {
          this.preloading = false
        })
        .catch(this.toastErrorHandler(this.$t('notification:role.fetch.error')))
    },

    getRoleLabel ({ name, handle, roleID }) {
      return name || handle || roleID
    },

    isRoleVisible ({ isClosed, meta = {} }) {
      return !(isClosed || (meta.context && meta.context.resourceTypes))
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
