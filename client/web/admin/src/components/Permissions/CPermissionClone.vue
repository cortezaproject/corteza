<template>
  <div
    class="d-inline-block"
  >
    <b-button
      data-test-id="button-clone"
      variant="secondary"
      class="mr-2"
      @click="showModal = true"
    >
      {{ $t('ui.clone.label') }}
    </b-button>

    <b-modal
      v-model="showModal"
      data-test-id="modal-clone-permission"
      ok-variant="primary"
      cancel-variant="link"
      centered
      :title="$t('ui.clone.title')"
      :ok-title="$t('ui.clone.clone')"
      :ok-disabled="!selectedRoles.length || processingRoles || processingSubmit"
      @ok="clonePermissions()"
    >
      <b-form-group
        :description="$t('ui.clone.description')"
        class="mb-0"
      >
        <vue-select
          v-model="selectedRoles"
          data-test-id="select-role-list"
          label="name"
          :options="roles"
          :get-option-key="getOptionKey"
          :reduce="role => role.roleID"
          :loading="processingRoles"
          multiple
          :placeholder="$t('ui.clone.pick-role')"
          :calculate-position="calculateDropdownPosition"
          class="bg-white"
        />
      </b-form-group>
    </b-modal>
  </div>
</template>

<script>
import VueSelect from 'vue-select'

export default {
  i18nOptions: {
    namespaces: 'permissions',
  },

  components: {
    VueSelect,
  },

  props: {
    roleId: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      showModal: false,

      roles: [],
      selectedRoles: [],

      processingSubmit: false,
      processingRoles: false,
    }
  },

  mounted () {
    this.processingRoles = true

    this.$SystemAPI.roleList()
      .then(({ set: roles = [] }) => {
        this.roles = roles
      })
      .catch(this.toastErrorHandler(this.$t('notification:role.fetch.error')))
      .finally(() => {
        this.processingRoles = false
      })
  },

  methods: {
    clonePermissions () {
      this.processingSubmit = true

      this.$SystemAPI.roleCloneRules({ roleID: this.roleId, cloneToRoleID: this.selectedRoles })
        .then(() => {
          this.selectedRoles = []
          this.toastSuccess(this.$t('notification:permissions.clone.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:permissions.clone.error')))
        .finally(() => {
          this.processingSubmit = false
          this.showModal = false
        })
    },

    getOptionKey ({ roleID }) {
      return roleID
    },
  },
}
</script>
