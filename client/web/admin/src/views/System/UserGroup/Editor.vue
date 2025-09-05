<template>
  <b-container
    v-if="userGroup"
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="title"
    >
      <b-button
        v-if="userGroupID && canCreate"
        data-test-id="button-new-user-group"
        variant="primary"
        :to="{ name: 'system.userGroup.new' }"
      >
        {{ $t('new') }}
      </b-button>
      <c-permissions-button
        v-if="userGroupID && canGrant"
        :title="userGroup.name || userGroup.handle || userGroup.userGroupID"
        :target="userGroup.name || userGroup.handle || userGroup.userGroupID"
        :resource="`corteza::system:user-group/${userGroupID}`"
      >
        <font-awesome-icon :icon="['fas', 'lock']" />
        {{ $t('permissions') }}
      </c-permissions-button>
    </c-content-header>

    <c-user-group-editor-info
      :user-group="userGroup"
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      :parent-user-groups="parentUserGroups"
      @submit="onInfoSubmit"
      @delete="onDelete"
    />
    <c-user-group-editor-members
      v-if="canManageMembers"
      class="mt-3"
      :processing="members.processing"
      :success="members.success"
      :current-members.sync="userGroupMembers"
      @submit="onMembersSubmit"
    />
    <c-user-group-editor-roles
      v-if="userGroup && userGroupID && membership.active"
      v-model="membership.active"
      class="mt-3"
      :processing="roles.processing"
      :success="roles.success"
      @submit="onMembershipSubmit"
    />
  </b-container>
</template>

<script>
import { isEqual } from 'lodash'
import { system } from '@cortezaproject/corteza-js'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CUserGroupEditorInfo from 'corteza-webapp-admin/src/components/UserGroup/CUserGroupEditorInfo'
import CUserGroupEditorMembers from 'corteza-webapp-admin/src/components/UserGroup/CUserGroupEditorMembers'
import CUserGroupEditorRoles from 'corteza-webapp-admin/src/components/UserGroup/CUserGroupEditorRoles'
import CPermissionClone from 'corteza-webapp-admin/src/components/Permissions/CPermissionClone'
import { mapGetters } from 'vuex'

export default {
  components: {
    CUserGroupEditorInfo,
    CUserGroupEditorMembers,
    CUserGroupEditorRoles,
    CPermissionClone,
  },

  i18nOptions: {
    namespaces: 'system.user-groups',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  props: {
    userGroupID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      userGroup: undefined,
      initialUserGroupState: undefined,
      isContext: false,

      userGroupMembers: undefined,

      membership: {
        active: undefined,
        initial: undefined,
      },

      info: {
        processing: false,
        success: false,
      },

      members: {
        processing: false,
        success: false,
      },

      roles: {
        processing: false,
        success: false,
      },

      parentUserGroups: [],
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canManageMembers () {
      return this.userGroup &&
        this.userGroup.canManageMembersOnUserGroup &&
        this.userGroup.userGroupID &&
        this.userGroupMembers &&
        !this.userGroup.isClosed &&
        !this.userGroup.isContext
    },

    canCreate () {
      return this.can('system/', 'user-group.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },

    title () {
      return this.userGroupID ? this.$t('title.edit') : this.$t('title.create')
    },
  },

  watch: {
    userGroupID: {
      immediate: true,
      async handler () {
        if (this.userGroupID) {
          this.fetchUserGroup()
          this.fetchMembership()
        } else {
          this.userGroup = new system.UserGroup()
          this.initialUserGroupState = this.userGroup.clone()
          this.isContext = false
          this.userGroupMembers = undefined
        }

        await this.fetchUserGroups()
      },
    },
  },

  methods: {
    async fetchUserGroups () {
      const { set: userGroups } = await this.$SystemAPI.userGroupList()

      this.parentUserGroups = userGroups
        .filter(({ userGroupID }) => userGroupID !== this.userGroupID)
        .map(({ name, handle, userGroupID }) => ({
          value: userGroupID,
          text: name || handle || userGroupID,
        }))
    },

    async fetchMembership () {
      this.incLoader()
      return this.$SystemAPI.roleList({ userGroupID: this.userGroupID })
        .then(({ set = [] }) => {
          this.membership = {
            active: [...set.map(({ roleID }) => roleID)],
            initial: [...set.map(({ roleID }) => roleID)],
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:user-group.roles.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    fetchUserGroup () {
      this.incLoader()

      if (this.userGroupID === '1') {
        // Do not show editor for role everyone
        this.$router.push({ name: 'system.userGroup.list' })
      }

      this.$SystemAPI.userGroupRead({ userGroupID: this.userGroupID })
        .then(r => {
          this.userGroup = new system.UserGroup(r)
          this.initialUserGroupState = this.userGroup.clone()

          this.isContext = !!this.userGroup.isContext

          if (this.userGroup.canManageMembersOnUserGroup && !this.userGroup.isContext && !this.userGroup.isClosed) {
            return this.$SystemAPI.userGroupMemberList(r).then((mm = []) => {
              this.userGroupMembers = mm.map(userID => ({ userID, current: true, dirty: true }))
            })
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:user-group.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    /**
     * Handles user role submit event, calls membership add or remove API endpoint
     * and handles response & errors
     */
    onMembershipSubmit () {
      this.roles.processing = true

      const userGroupID = this.userGroupID

      const { active, initial } = this.membership

      Promise.all([
        // all removed memberships
        ...initial.filter(roleID => !active.includes(roleID)).map(roleID => {
          return this.$SystemAPI.userGroupMembershipRemove({ roleID, userGroupID })
        }),
        // all new memberships
        ...active.filter(roleID => !initial.includes(roleID)).map(roleID => {
          return this.$SystemAPI.roleMemberAddGroup({ roleID, userGroupID })
        }),
      ])
        .then(async () => {
          this.animateSuccess('roles')
          await this.fetchMembership()

          this.toastSuccess(this.$t('notification:user-group.membershipUpdate.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user-group.membershipUpdate.error')))
        .finally(() => {
          this.roles.processing = false
        })
    },

    onDelete () {
      this.incLoader()

      if (this.userGroup.deletedAt) {
        this.$SystemAPI.userGroupUndelete({ userGroupID: this.userGroupID })
          .then(() => {
            this.fetchUserGroup()

            this.toastSuccess(this.$t('notification:user-group.undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:user-group.undelete.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI.userGroupDelete({ userGroupID: this.userGroupID })
          .then(() => {
            this.fetchUserGroup()

            this.userGroup.deletedAt = new Date()
            this.toastSuccess(this.$t('notification:user-group.delete.success'))
            this.$router.push({ name: 'system.user-group' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:user-group.delete.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    onInfoSubmit (userGroup) {
      this.info.processing = true

      if (this.userGroupID) {
        this.$SystemAPI.userGroupUpdate(userGroup)
          .then(() => {
            this.fetchUserGroup()

            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:user-group.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:user-group.update.error')))
          .finally(() => {
            this.info.processing = false
          })
      } else {
        this.$SystemAPI.userGroupCreate(userGroup)
          .then(({ userGroupID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:suer-group.create.success'))

            this.$router.push({ name: 'system.userGroup.edit', params: { userGroupID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:user-group.create.error')))
          .finally(() => {
            this.info.processing = false
          })
      }
    },

    onMembersSubmit () {
      this.members.processing = true

      const { userGroupID } = this.userGroup
      if (userGroupID) {
        Promise.all(this.userGroupMembers.map(async user => {
          const { userID, current, dirty } = user
          if (dirty !== current) {
            if (dirty) {
              return this.$SystemAPI.userGroupMemberAdd({ userGroupID, userID })
            } else {
              return this.$SystemAPI.userGroupMemberRemove({ userGroupID, userID })
            }
          }
        }))
          .then(() => {
            this.fetchUserGroup()
            this.animateSuccess('members')

            this.toastSuccess(this.$t('notification:user-group.membershipUpdate.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:user-group.membershipUpdate.error')))
          .finally(() => {
            this.members.processing = false
          })
      }
    },

    checkUnsavedChanges (next, to) {
      const isNewPage = this.$route.path.includes('/new') && to.name.includes('edit')
      const { deletedAt } = this.userGroup || {}

      if (isNewPage || deletedAt) {
        next(true)
      } else if (!to.name.includes('edit')) {
        const isDirty = (this.userGroupMembers || []).some(m => m.dirty !== m.current) || !isEqual(this.userGroup, this.initialUserGroupState)
        next(isDirty ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },
  },
}
</script>
