<template>
  <b-container
    v-if="role"
    class="py-3"
  >
    <c-content-header
      :title="title"
    >
      <span
        class="text-nowrap"
      >
        <b-button
          v-if="roleID && canCreate"
          data-test-id="button-new-role"
          variant="primary"
          class="mr-2"
          :to="{ name: 'system.role.new' }"
        >
          {{ $t('new') }}
        </b-button>

        <c-permissions-button
          v-if="roleID && canGrant"
          :title="role.name || role.handle || role.roleID"
          :target="role.name || role.handle || role.roleID"
          :resource="`corteza::system:role/${roleID}`"
          button-variant="light"
          class="mr-2"
        >
          <font-awesome-icon :icon="['fas', 'lock']" />
          {{ $t('permissions') }}
        </c-permissions-button>

        <c-permission-clone
          v-if="roleID && canGrant"
          :role-id="roleID"
        />
      </span>

      <c-corredor-manual-buttons
        ui-page="role/editor"
        ui-slot="toolbar"
        resource-type="system:role"
        default-variant="link"
        class="mr-1"
        @click="dispatchCortezaSystemRoleEvent($event, { role })"
      />
    </c-content-header>

    <c-role-editor-info
      :role="role"
      :processing="info.processing"
      :success="info.success"
      :is-context.sync="isContext"
      :can-create="canCreate"
      @submit="onInfoSubmit"
      @delete="onDelete"
      @status="onStatusChange"
    />

    <c-role-editor-members
      v-if="!isContext && canManageMembers"
      class="mt-3"
      :processing="members.processing"
      :success="members.success"
      :current-members.sync="roleMembers"
      @submit="onMembersSubmit"
    />
  </b-container>
</template>

<script>
import { system } from '@cortezaproject/corteza-js'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CRoleEditorInfo from 'corteza-webapp-admin/src/components/Role/CRoleEditorInfo'
import CRoleEditorMembers from 'corteza-webapp-admin/src/components/Role/CRoleEditorMembers'
import CPermissionClone from 'corteza-webapp-admin/src/components/Permissions/CPermissionClone'
import { mapGetters } from 'vuex'

export default {
  components: {
    CRoleEditorInfo,
    CRoleEditorMembers,
    CPermissionClone,
  },

  i18nOptions: {
    namespaces: 'system.roles',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    roleID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      role: undefined,
      isContext: false,

      roleMembers: null,

      info: {
        processing: false,
        success: false,
      },
      members: {
        processing: false,
        success: false,
      },
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canManageMembers () {
      return this.role &&
        this.role.canManageMembersOnRole &&
        this.role.roleID &&
        this.roleMembers &&
        !this.role.isClosed &&
        !this.role.isContext
    },

    canCreate () {
      return this.can('system/', 'role.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },

    title () {
      return this.roleID ? this.$t('title.edit') : this.$t('title.create')
    },
  },

  watch: {
    roleID: {
      immediate: true,
      handler () {
        if (this.roleID) {
          this.fetchRole()
        } else {
          this.role = new system.Role()
          this.isContext = false
        }
      },
    },

    'role.isContext': {
      immediate: true,
      handler (v) {
        if (v) {
          this.isContext = true
        }
      },
    },
  },

  methods: {
    fetchRole () {
      this.incLoader()

      if (this.roleID === '1') {
        // Do not show editor for role everyone
        this.$router.push({ name: 'system.role.list' })
      }

      this.$SystemAPI.roleRead({ roleID: this.roleID })
        .then(r => {
          this.role = new system.Role(r)
          this.isContext = !!this.role.isContext

          if (this.role.canManageMembersOnRole && !this.role.isContext && !this.role.isClosed) {
            return this.$SystemAPI.roleMemberList(r).then((mm = []) => {
              this.roleMembers = mm.map(userID => ({ userID, current: true, dirty: true }))
            })
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:role.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onDelete () {
      this.incLoader()

      if (this.role.deletedAt) {
        this.$SystemAPI.roleUndelete({ roleID: this.roleID })
          .then(() => {
            this.fetchRole()

            this.toastSuccess(this.$t('notification:role.undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:role.undelete.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI.roleDelete({ roleID: this.roleID })
          .then(() => {
            this.fetchRole()

            this.toastSuccess(this.$t('notification:role.delete.success'))
            this.$router.push({ name: 'system.role' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:role.delete.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    onInfoSubmit (role) {
      this.info.processing = true

      if (this.roleID) {
        this.$SystemAPI.roleUpdate(role)
          .then(role => {
            this.fetchRole()

            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:role.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:role.update.error')))
          .finally(() => {
            this.info.processing = false
          })
      } else {
        this.$SystemAPI.roleCreate(role)
          .then(({ roleID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:role.create.success'))

            this.$router.push({ name: 'system.role.edit', params: { roleID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:role.create.error')))
          .finally(() => {
            this.info.processing = false
          })
      }
    },

    /**
     * Handles user status change event, calls suspend or unsuspend API endpoint
     * and handles response & errors
     */
    onStatusChange () {
      this.incLoader()

      const roleID = this.roleID

      if (this.role.archivedAt) {
        this.$SystemAPI.roleUnarchive({ roleID })
          .then(() => {
            this.fetchRole()

            this.toastSuccess(this.$t('notification:role.unarchive.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:role.unarchive.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI.roleArchive({ roleID })
          .then(() => {
            this.fetchRole()

            this.toastSuccess(this.$t('notification:role.archive.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:role.archive.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    onMembersSubmit () {
      this.members.processing = true

      const { roleID } = this.role
      if (roleID) {
        Promise.all(this.roleMembers.map(async user => {
          let { userID, current, dirty } = user
          if (dirty !== current) {
            if (dirty) {
              return this.$SystemAPI.roleMemberAdd({ roleID, userID })
            } else {
              return this.$SystemAPI.roleMemberRemove({ roleID, userID })
            }
          }
        }))
          .then(() => {
            this.fetchRole()
            this.animateSuccess('members')

            this.toastSuccess(this.$t('notification:role.membershipUpdate.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:role.membershipUpdate.error')))
          .finally(() => {
            this.members.processing = false
          })
      }
    },
  },
}
</script>
