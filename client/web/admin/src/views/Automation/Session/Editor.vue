<template>
  <b-container
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-session-editor-info
      :session="session"
      :user="user"
      :processing="info.processing"
      @cancel="cancelSession()"
    />
  </b-container>
</template>
<script>
import { system } from '@cortezaproject/corteza-js'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CSessionEditorInfo from 'corteza-webapp-admin/src/components/Session/CSessionEditorInfo'

export default {
  components: {
    CSessionEditorInfo,
  },

  i18nOptions: {
    namespaces: 'automation.sessions',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    sessionID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      session: {},
      user: {},

      info: {
        processing: false,
      },
    }
  },

  watch: {
    sessionID: {
      immediate: true,
      handler () {
        if (this.sessionID) {
          this.fetchSession()
        }
      },
    },
  },

  methods: {
    fetchSession () {
      this.incLoader()
      this.info.processing = true

      this.$AutomationAPI.sessionRead({ sessionID: this.sessionID })
        .then(session => {
          this.session = session
          this.fetchUser()
        })
        .catch(this.toastErrorHandler(this.$t('notification:session.fetch.error')))
        .finally(() => {
          this.decLoader()
          this.info.processing = false
        })
    },

    cancelSession (sessionID = this.sessionID) {
      this.incLoader()
      this.info.processing = true

      this.$AutomationAPI.sessionCancel({ sessionID })
        .then(() => {
          this.fetchSession()
          this.toastSuccess(this.$t('notification:session.cancel.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:session.cancel.error')))
        .finally(() => {
          this.decLoader()
          this.info.processing = false
        })
    },

    fetchUser () {
      this.incLoader()

      this.$SystemAPI.userRead({ userID: this.session.createdBy })
        .then(user => {
          this.user = new system.User(user)
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },
  },
}
</script>
