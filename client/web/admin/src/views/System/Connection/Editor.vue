<template>
  <b-container
    class="pt-2 pb-3"
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
        :processing="info.processing"
        :success="info.success"
        :can-create="canCreate"
        :disabled="disabled"
        @submit="updateInfo"
        @delete="toggleDelete"
      />

      <c-connection-editor-properties
        v-if="connectionID && connection.meta.properties"
        :properties="connection.meta.properties"
        :processing="properties.processing"
        :success="properties.success"
        class="mt-4"
        @submit="updateProperties"
      />

      <c-connection-editor-dal
        v-if="connectionID && connection.config.dal && canManage"
        :dal="connection.config.dal"
        :issues="connection.issues || []"
        :can-manage="connection.canManageDalConfig"
        class="mt-4"
        @submit="updateDal"
      />
    </b-form>
  </b-container>
</template>

<script>
import { isEqual } from 'lodash'
import { system } from '@cortezaproject/corteza-js'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CConnectionEditorInfo from 'corteza-webapp-admin/src/components/Connection/CConnectionEditorInfo'
import CConnectionEditorProperties from 'corteza-webapp-admin/src/components/Connection/CConnectionEditorProperties'
import CConnectionEditorDal from 'corteza-webapp-admin/src/components/Connection/CConnectionEditorDAL'
import { mapGetters } from 'vuex'

export default {
  components: {
    CConnectionEditorInfo,
    CConnectionEditorDal,
    CConnectionEditorProperties,
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

      properties: {
        processing: false,
        success: false,
      },

      dal: {
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

    canManage () {
      return this.connection.canManageDalConfig
    },

    disabled () {
      return this.info.processing || this.properties.processing || this.dal.processing
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

    updateInfo () {
      const updating = !!this.connectionID
      const op = updating ? 'update' : 'create'
      const fn = updating ? 'dalConnectionUpdate' : 'dalConnectionCreate'

      this.info.processing = true
      this.incLoader()

      const connection = new system.DalConnection(this.initialConnectionState)
      connection.meta.name = this.connection.meta.name
      connection.handle = this.connection.handle
      connection.meta.location.properties.name = this.connection.meta.location.properties.name
      connection.meta.location.geometry.coordinates = this.connection.meta.location.geometry.coordinates
      connection.meta.ownership = this.connection.meta.ownership
      connection.config.privacy.sensitivityLevelID = this.connection.config.privacy.sensitivityLevelID

      return this.$SystemAPI[fn](connection).then(connection => {
        this.animateSuccess('info')
        this.toastSuccess(this.$t(`notification:connection.${op}.success`))

        if (!updating) {
          const { connectionID } = connection
          this.$router.push({ name: `system.connection.edit`, params: { connectionID } })
        } else {
          connection.config.dal = this.connection.config.dal
          connection.meta.properties = this.connection.meta.properties

          this.connection = new system.DalConnection(connection)
          this.initialConnectionState = this.connection.clone()
        }
      }).catch(this.toastErrorHandler(this.$t('notification:connection.update.error')))
        .finally(() => {
          this.info.processing = false
        })
    },

    updateProperties () {
      this.properties.processing = true
      this.incLoader()

      const connection = new system.DalConnection(this.initialConnectionState)
      connection.meta.properties = this.connection.meta.properties

      return this.$SystemAPI.dalConnectionUpdate(connection).then(connection => {
        this.animateSuccess('properties')
        this.toastSuccess(this.$t('notification:connection.update.success'))

        connection.meta.name = this.connection.meta.name
        connection.handle = this.connection.handle
        connection.meta.location.properties.name = this.connection.meta.location.properties.name
        connection.meta.location.geometry.coordinates = this.connection.meta.location.geometry.coordinates
        connection.meta.ownership = this.connection.meta.ownership
        connection.config.privacy.sensitivityLevelID = this.connection.config.privacy.sensitivityLevelID
        connection.config.dal = this.connection.config.dal

        this.connection = new system.DalConnection(connection)
        this.initialConnectionState = this.connection.clone()
      }).catch(this.toastErrorHandler(this.$t('notification:connection.update.error')))
        .finally(() => {
          this.properties.processing = false
        })
    },

    updateDal () {
      this.dal.processing = true
      this.incLoader()

      const connection = new system.DalConnection(this.initialConnectionState)
      connection.config.dal = this.connection.config.dal

      return this.$SystemAPI.dalConnectionUpdate(connection).then(connection => {
        this.animateSuccess('dal')
        this.toastSuccess(this.$t('notification:connection.update.success'))

        connection.meta.name = this.connection.meta.name
        connection.handle = this.connection.handle
        connection.meta.location.properties.name = this.connection.meta.location.properties.name
        connection.meta.location.geometry.coordinates = this.connection.meta.location.geometry.coordinates
        connection.meta.ownership = this.connection.meta.ownership
        connection.config.privacy.sensitivityLevelID = this.connection.config.privacy.sensitivityLevelID
        connection.meta.properties = this.connection.meta.properties

        this.connection = new system.DalConnection(connection)
        this.initialConnectionState = this.connection.clone()
      }).catch(this.toastErrorHandler(this.$t('notification:connection.update.error')))
        .finally(() => {
          this.dal.processing = false
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
            this.initialConnectionState.deletedAt = this.connection.deletedAt

            this.$router.push({ name: `system.connection` })
          } else {
            this.connection.deletedAt = null
            this.initialConnectionState.deletedAt = null
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
