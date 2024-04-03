<template>
  <div
    class="d-inline-block"
  >
    <b-button
      data-test-id="button-clone"
      variant="light"
      class="mr-2"
      @click="showModal = true"
    >
      {{ $t('ui.clone.label') }}
    </b-button>

    <b-modal
      v-model="showModal"
      data-test-id="modal-clone-permission"
      ok-variant="primary"
      cancel-variant="light"
      centered
      :title="$t('ui.clone.title')"
      :ok-title="$t('ui.clone.clone')"
      :ok-disabled="!selectedRoles.length || processingSubmit"
      @ok="clonePermissions"
    >
      <b-form-group
        :description="$t('ui.clone.description')"
        class="mb-0"
      >
        <c-input-role
          v-model="selectedRoles"
          data-test-id="select-role-list"
          :selectable="r => !selectedRoles.some(rr => rr.roleID === r.roleID)"
          :placeholder="$t('ui.clone.pick-role')"
          multiple
        />
      </b-form-group>
    </b-modal>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CInputRole } = components

export default {
  i18nOptions: {
    namespaces: 'permissions',
  },

  components: {
    CInputRole,
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

      selectedRoles: [],

      processingSubmit: false,
    }
  },

  methods: {
    clonePermissions () {
      this.processingSubmit = true

      const cloneToRoleID = this.selectedRoles.map(({ roleID }) => roleID)

      this.$SystemAPI.roleCloneRules({ roleID: this.roleId, cloneToRoleID })
        .then(() => {
          this.selectedRoles = []
          this.toastSuccess(this.$t('notification:permissions.clone.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:permissions.clone.error')))
        .finally(() => {
          this.processingSubmit = false
        })
    },
  },
}
</script>
