<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <div
      class="d-flex flex-column h-100"
    >
      <c-settings-editor
        :settings="apigwSettings"
        class="mb-3"
        @submit="onSettingsSubmit"
      />

      <c-route-list
        class="mb-4 flex-fill"
      />
    </div>
  </b-container>
</template>

<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CSettingsEditor from 'corteza-webapp-admin/src/components/Apigw/CSettingsEditor'
import CRouteList from 'corteza-webapp-admin/src/components/Apigw/CRouteList'

export default {
  components: {
    CRouteList,
    CSettingsEditor,
  },

  mixins: [
    editorHelpers,
  ],

  i18nOptions: {
    namespaces: [ 'system.apigw' ],
  },

  data () {
    return {
      settings: {
        processing: false,
        success: false,

        items: [],
      },
    }
  },

  computed: {
    apigwSettings () {
      if (this.settings.items.length > 0) {
        return this.settings.items.reduce((map, obj) => {
          const { name, value } = obj
          const split = name.split('.')

          if (split[0] === 'apigw') {
            map[name] = value
          }

          return map
        }, {})
      }
      return {}
    },
  },

  created () {
    this.fetchSettings()
  },

  methods: {
    onSettingsSubmit (settings) {
      this.settings.processing = true

      const values = Object.entries(settings).map(([name, value]) => {
        return { name, value }
      })

      this.$SystemAPI.settingsUpdate({ values })
        .then(() => {
          this.animateSuccess('settings')
          this.toastSuccess(this.$t('notification:settings.system.apigw.success'))
          this.$Settings.fetch()
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.system.apigw.error')))
        .finally(() => {
          this.settings.processing = false
        })
    },

    fetchSettings () {
      this.incLoader()

      this.$SystemAPI.settingsList()
        .then((settings = []) => {
          this.settings.items = settings
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.system.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },
  },
}
</script>
