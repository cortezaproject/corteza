<template>
  <b-container
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-compose-editor-basic
      :basic="settings"
      :processing="basic.processing"
      :success="basic.success"
      :can-manage="canManage"
      @submit="onSubmit($event, 'basic')"
    />

    <c-compose-editor-ui
      :settings="settings"
      :processing="ui.processing"
      :success="ui.success"
      :can-manage="canManage"
      class="mt-3"
      @submit="onSubmit($event, 'ui')"
    />
  </b-container>
</template>

<script>
import { isEqual, cloneDeep } from 'lodash'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CComposeEditorBasic from 'corteza-webapp-admin/src/components/Settings/Compose/CComposeEditorBasic'
import CComposeEditorUI from 'corteza-webapp-admin/src/components/Settings/Compose/CComposeEditorUI'
import { mapGetters } from 'vuex'

const prefix = 'compose.'

export default {
  i18nOptions: {
    namespaces: 'compose.settings',
    keyPrefix: 'editor',
  },

  components: {
    CComposeEditorBasic,
    'c-compose-editor-ui': CComposeEditorUI,
  },

  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      settings: {},
      initialSettingsState: {},

      basic: {
        processing: false,
        success: false,
      },

      ui: {
        processing: false,
        success: false,
      },
    }
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canManage () {
      return this.can('system/', 'settings.manage')
    },
  },

  created () {
    this.fetchSettings()
  },

  methods: {
    onSubmit (settings, type) {
      this[type].processing = true

      const values = Object.entries(settings).map(([name, value]) => {
        return { name, value }
      })

      this.$SystemAPI.settingsUpdate({ values })
        .then(() => {
          this.animateSuccess(type)
          this.toastSuccess(this.$t('notification:settings.compose.update.success'))
          this.initialSettingsState = cloneDeep(this.settings)
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.compose.update.error')))
        .finally(() => {
          this[type].processing = false
        })
    },

    fetchSettings () {
      this.incLoader()

      this.$SystemAPI.settingsList({ prefix })
        .then(settings => {
          settings.forEach(({ name, value }) => {
            this.$set(this.settings, name, value)
            this.$set(this.initialSettingsState, name, cloneDeep(value))
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.compose.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    checkUnsavedChanges (next, to) {
      const isNewPage = this.$route.path.includes('/new') && to.name.includes('edit')

      if (isNewPage) {
        next(true)
      } else if (!to.name.includes('edit')) {
        next(!isEqual(this.settings, this.initialSettingsState) ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },
  },
}
</script>
