<template>
  <div class="d-flex flex-column w-100 vh-100 overflow-hidden">
    <header
      v-show="loaded"
    >
      <c-topbar
        hide-app-selector
        :sidebar-pinned="pinned"
        :settings="$Settings.get('ui.topbar', {})"
        :labels="{
          helpForum: $t('help.forum'),
          helpDocumentation: $t('help.documentation'),
          helpFeedback: $t('help.feedback'),
          helpVersion: $t('help.version'),
          userSettingsLoggedInAs: $t('userSettings.loggedInAs', { user }),
          userSettingsProfile: $t('userSettings.profile'),
          userSettingsChangePassword: $t('userSettings.changePassword'),
          userSettingsLogout: $t('userSettings.logout'),
        }"
      >
        <template #help-dropdown>
          <portal-target name="topbar-help-dropdown" />
        </template>
      </c-topbar>
    </header>

    <main
      v-show="loaded"
      class="flex-fill overflow-hidden"
    >
      <c-app-selector
        :logo="logo"
      />
    </main>

    <c-loader-logo
      v-if="!loaded"
      :logo="logo"
    />

    <c-prompts
      :hide-toasts="!loaded"
    />
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import CAppSelector from '../components/CAppSelector'
import { components } from '@cortezaproject/corteza-vue'

const { CTopbar, CLoaderLogo, CPrompts } = components

export default {
  i18nOptions: {
    namespaces: 'navigation',
  },

  components: {
    CAppSelector,
    CTopbar,
    CLoaderLogo,
    CPrompts,
  },

  data () {
    return {
      loaded: false,

      pinned: false,
    }
  },

  computed: {
    icon () {
      return this.$Settings.attachment('ui.iconLogo')
    },

    logo () {
      return this.$Settings.attachment('ui.mainLogo')
    },

    user () {
      const { user } = this.$auth
      return user.name || user.handle || user.email || ''
    },
  },

  watch: {
    icon: {
      immediate: true,
      handler (icon) {
        if (icon) {
          const favicon = document.getElementById('favicon')
          favicon.href = icon
        }
      },
    },
  },

  created () {
    this.preloadApplications()
      .then(() => {
        setTimeout(() => {
          this.loaded = true
        }, 500)
      })
  },

  methods: {
    ...mapActions({
      preloadApplications: 'applications/load',
    }),
  },
}

</script>
<style lang="scss" scoped>

</style>
