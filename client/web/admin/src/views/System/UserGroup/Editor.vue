<template>
  <b-container
    v-if="userGroup"
    class="pt-2 pb-3"
  >
    <c-content-header :title="title">
      <b-button
        v-if="userGroupID && canCreate"
        data-test-id="button-new-user-group"
        variant="primary"
        :to="{ name: 'system.userGroup.new' }"
      >
        {{ $t("new") }}
      </b-button>
      <c-permissions-button
        v-if="userGroupID && canGrant"
        :title="userGroup.name || userGroup.handle || userGroup.userGroupID"
        :target="userGroup.name || userGroup.handle || userGroup.userGroupID"
        :resource="`corteza::system:user-group/${userGroupID}`"
      >
        <font-awesome-icon :icon="['fas', 'lock']" />
        {{ $t("permissions") }}
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
      v-if="canManageMembers && groupMembers.active"
      v-model="groupMembers.active"
      class="mt-3"
      :processing="groupMembers.processing"
      :success="groupMembers.success"
      @submit="onMembersSubmit"
    />

    <c-user-group-editor-roles
      v-if="userGroup && userGroupID && groupRoles.active"
      v-model="groupRoles.active"
      class="mt-3"
      :processing="groupRoles.processing"
      :success="groupRoles.success"
      @submit="onUserGroupRolesSubmit"
    />
  </b-container>
</template>

<script>
import { system } from '@cortezaproject/corteza-js'
import CUserGroupEditorInfo from 'corteza-webapp-admin/src/components/UserGroup/CUserGroupEditorInfo'
import CUserGroupEditorMembers from 'corteza-webapp-admin/src/components/UserGroup/CUserGroupEditorMembers'
import CUserGroupEditorRoles from 'corteza-webapp-admin/src/components/UserGroup/CUserGroupEditorRoles'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import { isEqual } from 'lodash'
import { mapGetters } from 'vuex'

export default {
  components: {
    CUserGroupEditorInfo,
    CUserGroupEditorMembers,
    CUserGroupEditorRoles,
  },

  i18nOptions: {
    namespaces: 'system.user-groups',
    keyPrefix: 'editor',
  },

  mixins: [editorHelpers],

  // beforeRouteUpdate (to, from, next) {
  //   this.checkUnsavedChanges(next, to)
  // },

  // beforeRouteLeave (to, from, next) {
  //   this.checkUnsavedChanges(next, to)
  // },

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

      info: {
        processing: false,
        success: false,
      },

      groupMembers: {
        active: undefined,
        initial: undefined,

        processing: false,
        success: false,
      },

      groupRoles: {
        active: undefined,
        initial: undefined,

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
      return (
        this.userGroup &&
        this.userGroup.canManageMembersOnUserGroup &&
        this.userGroup.userGroupID
      )
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
          this.fetchUserGroupRoles()
          this.fetchUserGroupMembers()
        } else {
          this.userGroup = new system.UserGroup()
          this.initialUserGroupState = this.userGroup.clone()
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

    async fetchUserGroupRoles () {
      this.incLoader()
      return this.$SystemAPI
        .roleList({ userGroupID: this.userGroupID })
        .then(({ set = [] }) => {
          this.groupRoles.active = [...set.map(({ roleID }) => roleID)]
          this.groupRoles.initial = [...set.map(({ roleID }) => roleID)]
        })
        .catch(
          this.toastErrorHandler(this.$t('notification:userGroup.roles.error')),
        )
        .finally(() => {
          this.decLoader()
        })
    },

    fetchUserGroupMembers () {
      this.incLoader()
      return this.$SystemAPI.userGroupMemberList({ userGroupID: this.userGroupID })
        .then((set = []) => {
          console.log({ set })
          this.groupMembers.active = [...set]
          this.groupMembers.initial = [...set]
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.roles.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    fetchUserGroup () {
      this.incLoader()

      this.$SystemAPI
        .userGroupRead({ userGroupID: this.userGroupID })
        .then((r) => {
          this.userGroup = new system.UserGroup(r)
          this.initialUserGroupState = this.userGroup.clone()
        })
        .catch(
          this.toastErrorHandler(this.$t('notification:userGroup.fetch.error')),
        )
        .finally(() => {
          this.decLoader()
        })
    },

    /**
     * Handles user role submit event, calls membership add or remove API endpoint
     * and handles response & errors
     */
    onUserGroupRolesSubmit () {
      this.groupRoles.processing = true

      const userGroupID = this.userGroupID
      const { active, initial } = this.groupRoles

      Promise.all([
        // all removed memberships
        ...initial
          .filter((roleID) => !active.includes(roleID))
          .map((roleID) => {
            return this.$SystemAPI.roleMemberRemoveGroup({
              roleID,
              userGroupID,
            })
          }),
        // all new memberships
        ...active
          .filter((roleID) => !initial.includes(roleID))
          .map((roleID) => {
            return this.$SystemAPI.roleMemberAddGroup({ roleID, userGroupID })
          }),
      ])
        .then(async () => {
          this.animateSuccess('groupRoles')
          await this.fetchUserGroupRoles()

          this.toastSuccess(
            this.$t('notification:userGroup.membershipUpdate.success'),
          )
        })
        .catch(
          this.toastErrorHandler(
            this.$t('notification:userGroup.membershipUpdate.error'),
          ),
        )
        .finally(() => {
          this.groupRoles.processing = false
        })
    },

    onDelete () {
      this.incLoader()

      if (this.userGroup.deletedAt) {
        this.$SystemAPI
          .userGroupUndelete({ userGroupID: this.userGroupID })
          .then(() => {
            this.fetchUserGroup()

            this.toastSuccess(
              this.$t('notification:userGroup.undelete.success'),
            )
          })
          .catch(
            this.toastErrorHandler(
              this.$t('notification:userGroup.undelete.error'),
            ),
          )
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI
          .userGroupDelete({ userGroupID: this.userGroupID })
          .then(() => {
            this.fetchUserGroup()

            this.userGroup.deletedAt = new Date()
            this.toastSuccess(this.$t('notification:userGroup.delete.success'))
            this.$router.push({ name: 'system.user-group' })
          })
          .catch(
            this.toastErrorHandler(
              this.$t('notification:userGroup.delete.error'),
            ),
          )
          .finally(() => {
            this.decLoader()
          })
      }
    },

    onInfoSubmit (userGroup) {
      this.info.processing = true

      if (this.userGroupID) {
        this.$SystemAPI
          .userGroupUpdate(userGroup)
          .then(() => {
            this.fetchUserGroup()

            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:userGroup.update.success'))
          })
          .catch(
            this.toastErrorHandler(
              this.$t('notification:userGroup.update.error'),
            ),
          )
          .finally(() => {
            this.info.processing = false
          })
      } else {
        this.$SystemAPI
          .userGroupCreate(userGroup)
          .then(({ userGroupID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:userGroup.create.success'))

            this.$router.push({
              name: 'system.userGroup.edit',
              params: { userGroupID },
            })
          })
          .catch(
            this.toastErrorHandler(
              this.$t('notification:userGroup.create.error'),
            ),
          )
          .finally(() => {
            this.info.processing = false
          })
      }
    },

    onMembersSubmit () {
      this.groupMembers.processing = true

      const userGroupID = this.userGroupID

      const { active, initial } = this.groupMembers

      Promise.all([
        // @note there are no removed memberships

        // all new memberships
        ...active.filter(userID => !initial.includes(userID)).map(userID => {
          return this.$SystemAPI.userGroupMemberAdd({ userGroupID, userID })
        }),
      ])
        .then(() => {
          this.animateSuccess('groupMembers')
          this.fetchUserGroupMembers()

          this.toastSuccess(this.$t('notification:user.membershipUpdate.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.membershipUpdate.error')))
        .finally(() => {
          this.groupMembers.processing = false
        })
    },

    // checkUnsavedChanges (next, to) {
    //   const isNewPage =
    //     this.$route.path.includes('/new') && to.name.includes('edit')
    //   const { deletedAt } = this.userGroup || {}

    //   if (isNewPage || deletedAt) {
    //     next(true)
    //   } else if (!to.name.includes('edit')) {
    //     const isDirty =
    //       (this.userGroupMembers || []).some((m) => m.dirty !== m.current) ||
    //       !isEqual(this.userGroup, this.initialUserGroupState)
    //     next(
    //       isDirty
    //         ? window.confirm(this.$t('general:editor.unsavedChanges'))
    //         : true,
    //     )
    //   }
    // },
  },
}
</script>
