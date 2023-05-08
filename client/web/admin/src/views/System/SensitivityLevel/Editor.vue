<template>
  <b-container
    v-if="sensitivityLevel"
    class="py-3"
  >
    <c-content-header
      :title="title"
    >
      <span
        class="text-nowrap"
      >
        <b-button
          v-if="sensitivityLevelID && canCreate"
          variant="primary"
          class="mr-2"
          :to="{ name: 'system.sensitivityLevel.new' }"
        >
          {{ $t('new') }}
        </b-button>
      </span>
    </c-content-header>

    <c-sensitivity-level-editor-info
      :sensitivity-level="sensitivityLevel"
      :processing="info.processing"
      :success="info.success"
      :can-delete="canDelete"
      :can-create="canCreate"
      @submit="onSubmit($event)"
      @delete="onDelete($event)"
    />
  </b-container>
</template>
<script>
import { isEqual, cloneDeep } from 'lodash'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CSensitivityLevelEditorInfo from 'corteza-webapp-admin/src/components/SensitivityLevel/CSensitivityLevelEditorInfo'
import { mapGetters } from 'vuex'

export default {
  components: {
    CSensitivityLevelEditorInfo,
  },

  i18nOptions: {
    namespaces: 'system.sensitivityLevel',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    sensitivityLevelID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      sensitivityLevel: undefined,
      initialSensitivityLevelState: undefined,

      info: {
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
      return this.can('system/', 'dal-sensitivity-level.manage')
    },

    canDelete () {
      return this.sensitivityLevel && this.sensitivityLevel.sensitivityLevelID && this.canCreate
    },

    title () {
      return this.sensitivityLevelID ? this.$t('title.edit') : this.$t('title.create')
    },
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  watch: {
    sensitivityLevelID: {
      immediate: true,
      handler () {
        if (this.sensitivityLevelID) {
          this.fetchSensitivityLevel()
        } else {
          this.sensitivityLevel = {
            handle: '',
            level: 1,
            meta: {
              name: '',
              description: '',
            },
          }

          this.initialSensitivityLevelState = cloneDeep(this.sensitivityLevel)
        }
      },
    },
  },

  methods: {
    fetchSensitivityLevel (sensitivityLevelID = this.sensitivityLevelID) {
      this.incLoader()

      this.$SystemAPI.dalSensitivityLevelRead({ sensitivityLevelID })
        .then(sensitivityLevel => {
          this.sensitivityLevel = sensitivityLevel
          this.initialSensitivityLevelState = cloneDeep(sensitivityLevel)
        })
        .catch(this.toastErrorHandler(this.$t('notification:sensitivityLevel.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onSubmit (sensitivityLevel) {
      this.info.processing = true

      if (this.sensitivityLevelID) {
        this.$SystemAPI.dalSensitivityLevelUpdate(sensitivityLevel)
          .then(sensitivityLevel => {
            this.sensitivityLevel = sensitivityLevel
            this.initialSensitivityLevelState = cloneDeep(sensitivityLevel)

            this.toastSuccess(this.$t('notification:sensitivityLevel.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:sensitivityLevel.update.error')))
          .finally(() => {
            this.info.processing = false
          })
      } else {
        this.$SystemAPI.dalSensitivityLevelCreate(sensitivityLevel)
          .then(sensitivityLevel => {
            this.sensitivityLevel = sensitivityLevel
            this.initialSensitivityLevelState = cloneDeep(sensitivityLevel)

            const { sensitivityLevelID } = sensitivityLevel
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:sensitivityLevel.create.success'))

            this.$router.push({ name: 'system.sensitivityLevel.edit', params: { sensitivityLevelID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:sensitivityLevel.create.error')))
          .finally(() => {
            this.info.processing = false
          })
      }
    },

    onDelete (sensitivityLevelID = this.sensitivityLevelID) {
      this.incLoader()

      if (this.sensitivityLevel.deletedAt) {
        // Sensitivity level is currently deleted -- undelete
        this.$SystemAPI.dalSensitivityLevelUndelete({ sensitivityLevelID })
          .then(() => {
            this.fetchSensitivityLevel()

            this.toastSuccess(this.$t('notification:sensitivityLevel.undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:sensitivityLevel.undelete.error')))
          .finally(() => this.decLoader())
      } else {
        // Sensitivity level is currently not deleted -- delete
        this.$SystemAPI.dalSensitivityLevelDelete({ sensitivityLevelID })
          .then(() => {
            this.fetchSensitivityLevel()

            this.toastSuccess(this.$t('notification:sensitivityLevel.delete.success'))
            this.$router.push({ name: 'system.sensitivityLevel' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:sensitivityLevel.delete.error')))
          .finally(() => this.decLoader())
      }
    },

    checkUnsavedChanges (next, to) {
      const isNewPage = this.$route.path.includes('/new') && to.name.includes('edit')

      if (isNewPage) {
        next(true)
      } else if (!to.name.includes('edit')) {
        next(!isEqual(this.sensitivityLevel, this.initialSensitivityLevelState) ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },
  },
}
</script>
