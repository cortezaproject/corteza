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
      <!--
        include hidden input to enable
        trigger submit event w/ ENTER
      -->
      <input
        type="submit"
        class="d-none"
        :disabled="saveDisabled"
      >

      <div
        class="d-flex mt-2"
      >
        <confirmation-toggle
          v-if="connection && connectionID && !isPrimary && !disabled"
          @confirmed="toggleDelete"
        >
          {{ connection.deletedAt ? $t('general:label.undelete') : $t('general:label.delete') }}
        </confirmation-toggle>

        <c-submit-button
          :processing="processing"
          :disabled="disabled || saveDisabled"
          class="ml-auto"
          @submit="onSubmit"
        />
      </div>
    </b-form>
  </b-container>
</template>

<script>
import { system, NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CConnectionEditorInfo from 'corteza-webapp-admin/src/components/Connection/CConnectionEditorInfo'
import CConnectionEditorProperties from 'corteza-webapp-admin/src/components/Connection/CConnectionEditorProperties'
import CConnectionEditorDal from 'corteza-webapp-admin/src/components/Connection/CConnectionEditorDAL'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'
import { mapGetters } from 'vuex'

export default {
  components: {
    CConnectionEditorInfo,
    CConnectionEditorDal,
    CConnectionEditorProperties,
    ConfirmationToggle,
    CSubmitButton,
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
      processing: false,
      connection: undefined,

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
      return this.processing
    },

    fresh () {
      return !this.connection.connectionID || this.connection.connectionID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : true // this.user.canUpdateUser
    },

    handleState () {
      return handle.handleState(this.connection.handle)
    },

    saveDisabled () {
      return !this.editable || [this.handleState].includes(false)
    },

  },

  watch: {
    connectionID: {
      immediate: true,
      handler (connectionID) {
        if (connectionID) {
          this.fetchConnection(connectionID)
        } else {
          this.connection = new system.DalConnection()
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
      }).catch(this.toastErrorHandler(this.$t('notification:connection.fetch.error')))
        .finally(async () => {
          this.decLoader()
        })
    },

    async fetchSensitivityLevels () {
      this.processing = true

      return this.$SystemAPI.dalSensitivityLevelList()
        .then(({ set = [] }) => {
          this.sensitivityLevels = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:sensitivityLevel.fetch.error')))
        .finally(() => {
          this.processing = false
        })
    },

    onSubmit () {
      const updating = !!this.connectionID
      const op = updating ? 'update' : 'create'
      const fn = updating ? 'dalConnectionUpdate' : 'dalConnectionCreate'

      this.processing = true
      this.incLoader()

      return this.$SystemAPI[fn](this.connection)
        .then(connection => {
          const { connectionID } = connection

          this.toastSuccess(this.$t(`notification:connection.${op}.success`))
          if (!updating) {
            this.$router.push({ name: `system.connection.edit`, params: { connectionID } })
          } else {
            this.connection = new system.DalConnection(connection)
          }
        })
        .catch(this.toastErrorHandler(this.$t(`notification:connection.${op}.error`)))
        .finally(() => {
          this.processing = false
        })
    },

    toggleDelete () {
      const { deletedAt } = this.connection
      const deleting = !deletedAt
      const op = deleting ? 'delete' : 'undelete'
      const fn = deleting ? 'dalConnectionDelete' : 'dalConnectionUndelete'

      this.processing = true
      this.incLoader()

      return this.$SystemAPI[fn](this.connection)
        .then(connection => {
          this.toastSuccess(this.$t(`notification:connection.${op}.success`))

          if (deleting) {
            /**
             * Resource deleted, move back to the list
             */
            this.$router.push({ name: `system.connection` })
          } else {
            this.connection.deletedAt = null
          }
        })
        .catch(this.toastErrorHandler(this.$t(`notification:connection.${op}.error`)))
        .finally(() => {
          this.processing = false
        })
    },
  },
}
</script>
