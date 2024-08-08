<template>
  <b-container
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <cui-cdns-editor
      :settings="settings"
      :processing="cdns.processing"
      :success="cdns.success"
      :can-manage="canManage"
      @submit="onSubmit($event)"
    />
  </b-container>
</template>

<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CUICdnsEditor from '../../components/Settings/UI/CUICdnsEditor.vue'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: [ 'ui.cdns' ],
    keyPrefix: 'editor',
  },

  components: {
    'cui-cdns-editor': CUICdnsEditor,
  },

  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      settings: [],

      cdns: {
        processing: false,
        success: false,
      },
    }
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
    fetchSettings () {
      this.incLoader()

      this.$Settings.fetch()
      return this.$SystemAPI.settingsList({ prefix: 'ui.cdn-scripts' })
        .then(settings => {
          this.settings = settings
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.ui.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onSubmit (settings) {
      this['cdns'].processing = true

      const values = Object.entries(settings).map(([name, value]) => {
        return { name, value }
      })

      this.$SystemAPI.settingsUpdate({ values: values })
        .then(() => {
          this.animateSuccess('cdns')
          this.toastSuccess(this.$t('notification:settings.ui.update.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.ui.update.error')))
        .finally(() => {
          this['cdns'].processing = false
        })
    },
  },
}
</script>
