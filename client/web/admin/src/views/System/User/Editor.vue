<template>
  <b-container
    v-if="user"
    class="py-3"
  >
    <c-content-header
      :title="title"
    >
      <span
        v-if="userID"
        class="text-nowrap"
      >
        <b-button
          v-if="canCreate"
          data-test-id="button-new-user"
          variant="primary"
          class="mr-2"
          :to="{ name: 'system.user.new' }"
        >
          {{ $t('new') }}
        </b-button>
        <c-permissions-button
          v-if="canGrant"
          :title="user.name || user.handle || user.email || userID"
          :target="user.name || user.handle || user.email || userID"
          :resource="`corteza::system:user/${userID}`"
          button-variant="light"
        >
          <font-awesome-icon :icon="['fas', 'lock']" />
          {{ $t('permissions') }}
        </c-permissions-button>
      </span>
      <c-corredor-manual-buttons
        ui-page="user/editor"
        ui-slot="toolbar"
        resource-type="system:user"
        default-variant="link"
        class="mr-1"
        @click="dispatchCortezaSystemUserEvent($event, { user })"
      />
    </c-content-header>

    <c-user-editor-info
      :user="user"
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      @submit="onInfoSubmit"
      @delete="onDelete"
      @status="onStatusChange"
      @patch="onPatch"
      @sessionsRevoke="onSessionsRevoke"
    />

    <c-user-editor-avatar
      v-if="user && userID && $Settings.get('auth.internal.profile-avatar.Enabled', false)"
      :user="user"
      :processing="avatar.processing"
      :success="avatar.success"
      class="mt-3"
      @submit="onAvatarSubmit"
      @onUpload="onAvatarUpload"
      @resetAttachment="onResetAvatar"
    />

    <c-user-editor-roles
      v-if="user && userID"
      v-model="membership.active"
      class="mt-3"
      :processing="roles.processing"
      :success="roles.success"
      @submit="onMembershipSubmit"
    />

    <c-user-editor-mfa
      v-if="user && userID"
      class="mt-3"
      :mfa="user.meta.securityPolicy.mfa"
      :processing="mfa.processing"
      :success="mfa.success"
      @patch="onPatch"
    />

    <c-user-editor-password
      v-if="user && userID"
      class="mt-3"
      :processing="password.processing"
      :success="password.success"
      :user-i-d="userID"
      @submit="onPasswordSubmit"
    />

    <c-user-editor-external-auth-providers
      v-if="user && userID"
      class="mt-3"
      :value="externalAuthProviders"
      @delete="onExternalAuthProviderDelete"
    />
  </b-container>
</template>

<script>
import { isEqual } from 'lodash'
import { NoID, system } from '@cortezaproject/corteza-js'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CUserEditorInfo from 'corteza-webapp-admin/src/components/User/CUserEditorInfo'
import CUserEditorAvatar from '../../../components/User/CUserEditorAvatar'
import CUserEditorPassword from 'corteza-webapp-admin/src/components/User/CUserEditorPassword'
import CUserEditorMfa from 'corteza-webapp-admin/src/components/User/CUserEditorMFA'
import CUserEditorRoles from 'corteza-webapp-admin/src/components/User/CUserEditorRoles'
import CUserEditorExternalAuthProviders from 'corteza-webapp-admin/src/components/User/CUserEditorExternalAuthProviders'
import { mapGetters } from 'vuex'

export default {
  components: {
    CUserEditorRoles,
    CUserEditorPassword,
    CUserEditorInfo,
    CUserEditorAvatar,
    CUserEditorMfa,
    CUserEditorExternalAuthProviders,
  },

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    userID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      user: undefined,
      initialUserState: undefined,

      membership: {
        active: [],
        initial: [],
      },

      externalAuthProviders: [],

      // Processing and success flags for each form
      info: {
        processing: false,
        success: false,
      },
      avatar: {
        processing: false,
        success: false,
      },
      password: {
        processing: false,
        success: false,
      },
      mfa: {
        processing: false,
        success: false,
      },
      roles: {
        processing: false,
        success: false,
      },
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('system/', 'user.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },

    title () {
      return this.userID ? this.$t('title.edit') : this.$t('title.create')
    },
  },

  watch: {
    userID: {
      immediate: true,
      handler (userID) {
        if (userID) {
          this.fetchUser()
          this.fetchMembership()
          this.fetchExternalAuthProviders()
        } else {
          this.user = new system.User()
          this.initialUserState = this.user.clone()
        }
      },
    },
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  methods: {
    makeEvent (res) {
      return system.UserEvent(res)
    },

    async fetchUser () {
      this.incLoader()

      return this.$SystemAPI.userRead({ userID: this.userID })
        .then(user => {
          this.user = new system.User(user)
          this.initialUserState = this.user.clone()
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    fetchMembership () {
      this.incLoader()
      return this.$SystemAPI.userMembershipList({ userID: this.userID })
        .then((set = []) => {
          this.membership = { active: [...set], initial: [...set] }
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.roles.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    fetchExternalAuthProviders () {
      this.incLoader()

      return this.$SystemAPI.userListCredentials({ userID: this.userID })
        .then((providers = []) => {
          this.externalAuthProviders = providers.map(({ credentialsID = '', label = '', kind = '' }) => ({ credentialsID, label, type: kind }))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.external-auth-providers.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    /**
     * Handles user info submit event, calls user update or create API endpoint
     * and handles response & errors
     *
     * @param user {Object}
     */
    onInfoSubmit (user) {
      this.info.processing = true

      const payload = { ...user }

      if (payload.userID !== NoID) {
        // On update, reset the user obj
        this.$SystemAPI.userUpdate(payload)
          .then(user => {
            this.user = new system.User(user)
            this.initialUserState = this.user.clone()

            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:user.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:user.update.error')))
          .finally(() => {
            this.info.processing = false
          })
      } else {
        // On creation, redirect to edit page
        this.$SystemAPI.userCreate(payload)
          .then(({ userID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:user.create.success'))

            this.$router.push({ name: 'system.user.edit', params: { userID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:user.create.error')))
          .finally(() => {
            this.info.processing = false
          })
      }
    },

    onAvatarSubmit (user) {
      this.avatar.processing = true

      const payload = {
        userID: user.userID,
        avatarColor: user.meta.avatarColor,
        avatarBgColor: user.meta.avatarBgColor,
      }

      this.$SystemAPI.userProfileAvatarInitial(payload)
        .then(() => this.fetchUser())
        .then(() => {
          this.animateSuccess('avatar')
          this.toastSuccess(this.$t('notification:user.avatarSettings.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.avatarSettings.error')))
        .finally(() => {
          this.avatar.processing = false
        })
    },

    /**
     * Handles user delete event, calls user delete API endpoint
     * and handles response & errors
     */
    onDelete () {
      this.incLoader()

      if (this.user.deletedAt) {
        this.$SystemAPI.userUndelete({ userID: this.userID })
          .then(() => {
            this.fetchUser()

            this.toastSuccess(this.$t('notification:user.undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:user.undelete.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI.userDelete({ userID: this.userID })
          .then(() => {
            this.fetchUser()

            this.toastSuccess(this.$t('notification:user.delete.success'))
            this.$router.push({ name: 'system.user' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:user.delete.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    onExternalAuthProviderDelete (credentialsID = '') {
      this.incLoader()

      this.$SystemAPI.userDeleteCredentials({ userID: this.userID, credentialsID })
        .then(() => {
          this.fetchExternalAuthProviders()

          this.toastSuccess(this.$t('notification:user.external-auth-providers.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.external-auth-providers.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    /**
     * Handles user password submit event, calls set password API endpoint
     * and handles response & errors
     *
     * @param password {String}
     */
    onPasswordSubmit (password = '') {
      this.password.processing = true

      this.$SystemAPI.userSetPassword({ userID: this.userID, password })
        .then(() => {
          this.fetchExternalAuthProviders()
          this.animateSuccess('password')

          this.toastSuccess(this.$t('notification:user.passwordChange.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.passwordChange.error')))
        .finally(() => {
          this.password.processing = false
        })
    },

    /**
     * Handles user MFA submit event
     *
     */
    onPatch (path, value) {
      const cfg = {
        method: 'patch',
        url: this.$SystemAPI.userPartialUpdateEndpoint({ userID: this.userID }),
        data: [{ path, value, op: 'replace' }],
      }

      return this.$SystemAPI.api().request(cfg).then(response => {
        this.fetchUser()
      })
    },

    /**
     * Handles user role submit event, calls membership add or remove API endpoint
     * and handles response & errors
     */
    onMembershipSubmit () {
      this.roles.processing = true

      const userID = this.userID

      const { active, initial } = this.membership

      Promise.all([
        // all removed memberships
        ...initial.filter(roleID => !active.includes(roleID)).map(roleID => {
          return this.$SystemAPI.userMembershipRemove({ roleID, userID })
        }),
        // all new memberships
        ...active.filter(roleID => !initial.includes(roleID)).map(roleID => {
          return this.$SystemAPI.userMembershipAdd({ roleID, userID })
        }),
      ])
        .then(() => {
          this.animateSuccess('roles')
          this.fetchMembership()

          this.toastSuccess(this.$t('notification:user.membershipUpdate.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.membershipUpdate.error')))
        .finally(() => {
          this.roles.processing = false
        })
    },

    /**
     * Handles user status change event, calls suspend or unsuspend API endpoint
     * and handles response & errors
     */
    onStatusChange () {
      this.incLoader()

      const userID = this.userID

      if (this.user.suspendedAt) {
        this.$SystemAPI.userUnsuspend({ userID })
          .then(() => {
            this.fetchUser()

            this.toastSuccess(this.$t('notification:user.unsuspend.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:user.unsuspend.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI.userSuspend({ userID })
          .then(() => {
            this.fetchUser()

            this.toastSuccess(this.$t('notification:user.suspend.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:user.suspend.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    /**
     * Handles user logout event, calls user logout API endpoint
     * Removes all active auth session and token of user
     */
    onSessionsRevoke () {
      this.incLoader()

      const userID = this.userID

      this.$SystemAPI.userSessionsRemove({ userID })
        .then(() => {
          this.fetchUser()

          this.toastSuccess(this.$t('notification:user.sessionsRevoke.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.sessionsRevoke.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onAvatarUpload () {
      this.fetchUser().then(() => {
        this.toastSuccess(this.$t('notification:user.avatarUpload.success'))
      })
    },

    onResetAvatar () {
      const userID = this.userID

      this.$SystemAPI.userDeleteAvatar({ userID })
        .then(() => this.fetchUser())
        .then(() => {
          this.toastSuccess(this.$t('notification:user.avatarDelete.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.avatarDelete.error')))
    },

    checkUnsavedChanges (next, to) {
      const isNewPage = this.$route.path.includes('/new') && to.name.includes('edit')

      if (isNewPage) {
        next(true)
      } else if (!to.name.includes('edit')) {
        let userChangesStatus = !isEqual(this.user, this.initialUserState)
        let membershipChangesStatus = !isEqual(this.membership.initial, this.membership.active)

        next((userChangesStatus || membershipChangesStatus) ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },
  },
}
</script>
