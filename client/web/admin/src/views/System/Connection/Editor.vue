<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="connectionID ? $t('title.edit') : $t('title.create')"
    />

    <b-form
      v-if="connection && sensitivityLevels"
      @submit.prevent="onSubmit"
    >
      <c-connection-editor-info
        :connection="connection"
        :sensitivity-levels="sensitivityLevels"
        :disabled="disabled"
        :is-primary="isPrimary"
      />

      <c-connection-editor-properties
        :properties="connection.meta.properties"
        :disabled="disabled"
        class="mt-4"
      />

      <c-connection-editor-dal
        v-if="connection.config.dal"
        :dal="connection.config.dal"
        :issues="connection.issues || []"
        :disabled="disabled"
        :can-manage="connection.canManageDalConfig"
        class="mt-4"
      />

      <b-card
        body-class="d-flex flex-wrap flex-fill-child gap-1"
        class="mt-4"
      >
        <confirmation-toggle
          v-if="connection && connectionID && !isPrimary && !disabled"
          @confirmed="toggleDelete"
        >
          {{ connection.deletedAt ? $t('general:label.undelete') : $t('general:label.delete') }}
        </confirmation-toggle>

        <c-button-submit
          :disabled="disabled || saveDisabled"
          :processing="info.processing"
          :success="info.success"
          :text="$t('admin:general.label.submit')"
          class="ml-auto"
          @submit="onSubmit"
        />
      </b-card>
    </b-form>
  </b-container>
</template>

<script>
import { isEqual } from 'lodash'
import { system, NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CConnectionEditorInfo from 'corteza-webapp-admin/src/components/Connection/CConnectionEditorInfo'
import CConnectionEditorProperties from 'corteza-webapp-admin/src/components/Connection/CConnectionEditorProperties'
import CConnectionEditorDal from 'corteza-webapp-admin/src/components/Connection/CConnectionEditorDAL'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'
import { mapGetters } from 'vuex'

export default {
  components: {
    CConnectionEditorInfo,
    CConnectionEditorDal,
    CConnectionEditorProperties,
    ConfirmationToggle,
  },

  i18nOptions: {
    namespaces: 'system.connections',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    connectionID: {
      type: String,
      default: undefined,
    },
  },

  data () {
    return {
      info: {
        processing: false,
        success: false,
      },

      connection: undefined,
      initialConnectionState: undefined,

      sensitivityLevels: undefined,
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('system/', 'dal-connection.create')
    },

    isPrimary () {
      return this.connection.type === 'corteza::system:primary-dal-connection'
    },

    disabled () {
      return this.info.processing
    },

    fresh () {
      return !this.connection.connectionID || this.connection.connectionID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : true // this.user.canUpdateUser
    },

    nameState () {
      return this.connection.meta.name ? null : false
    },

    handleState () {
      return handle.handleState(this.connection.handle)
    },

    saveDisabled () {
      return !this.editable || [this.nameState, this.handleState].includes(false)
    },

  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  watch: {
    connectionID: {
      immediate: true,
      handler (connectionID) {
        if (connectionID) {
          this.fetchConnection(connectionID)
        } else {
          this.connection = new system.DalConnection()
          this.initialConnectionState = this.connection.clone()
        }
      },
    },
  },

  mounted () {
    this.fetchSensitivityLevels()
  },

  methods: {
    fetchConnection (connectionID) {
      this.incLoader()
      return this.$SystemAPI.dalConnectionRead({ connectionID }).then(connection => {
        this.connection = new system.DalConnection(connection)
        this.initialConnectionState = this.connection.clone()
      }).catch(this.toastErrorHandler(this.$t('notification:connection.fetch.error')))
        .finally(async () => {
          this.decLoader()
        })
    },

    async fetchSensitivityLevels () {
      this.info.processing = true

      return this.$SystemAPI.dalSensitivityLevelList()
        .then(({ set = [] }) => {
          this.sensitivityLevels = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:sensitivityLevel.fetch.error')))
        .finally(() => {
          this.info.processing = false
        })
    },

    onSubmit () {
      const updating = !!this.connectionID
      const op = updating ? 'update' : 'create'
      const fn = updating ? 'dalConnectionUpdate' : 'dalConnectionCreate'

      this.info.processing = true
      this.incLoader()

      return this.$SystemAPI[fn](this.connection)
        .then(connection => {
          const { connectionID } = connection

          this.animateSuccess('info')
          this.toastSuccess(this.$t(`notification:connection.${op}.success`))
          if (!updating) {
            this.$router.push({ name: `system.connection.edit`, params: { connectionID } })
          } else {
            this.connection = new system.DalConnection(connection)
            this.initialConnectionState = this.connection.clone()
          }
        })
        .catch(this.toastErrorHandler(this.$t(`notification:connection.${op}.error`)))
        .finally(() => {
          this.info.processing = false
        })
    },

    toggleDelete () {
      const { deletedAt } = this.connection
      const deleting = !deletedAt
      const op = deleting ? 'delete' : 'undelete'
      const fn = deleting ? 'dalConnectionDelete' : 'dalConnectionUndelete'

      this.info.processing = true
      this.incLoader()

      return this.$SystemAPI[fn](this.connection)
        .then(() => {
          this.toastSuccess(this.$t(`notification:connection.${op}.success`))

          if (deleting) {
            /**
             * Resource deleted, move back to the list
             */
            this.connection.deletedAt = new Date()
            this.$router.push({ name: `system.connection` })
          } else {
            this.connection.deletedAt = null
          }
        })
        .catch(this.toastErrorHandler(this.$t(`notification:connection.${op}.error`)))
        .finally(() => {
          this.info.processing = false
        })
    },

    checkUnsavedChanges (next, to) {
      const isNewPage = this.$route.path.includes('/new') && to.name.includes('edit')

      if (isNewPage || this.connection.deletedAt) {
        next(true)
      } else if (!to.name.includes('edit')) {
        next(!isEqual(this.connection, this.initialConnectionState) ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },
  },
}
</script>
