<template>
  <div class="d-flex flex-column w-100 vh-100 overflow-hidden">
    <header>
      <c-topbar
        :sidebar-pinned="pinned"
        :settings="$Settings.get('ui.topbar', {})"
        :labels="{
          appMenu: $t('appMenu'),
          helpForum: $t('help.forum'),
          helpDocumentation: $t('help.documentation'),
          helpFeedback: $t('help.feedback'),
          helpVersion: $t('help.version'),
          userSettingsLoggedInAs: $t('userSettings.loggedInAs', { user }),
          userSettingsProfile: $t('userSettings.profile'),
          userSettingsChangePassword: $t('userSettings.changePassword'),
          userSettingsLogout: $t('userSettings.logout'),
          lightTheme: $t('general:themes.labels.light'),
          darkTheme: $t('general:themes.labels.dark'),
        }"
      >
        <template #title>
          {{ $t('discovery') }}
        </template>
      </c-topbar>
    </header>

    <aside>
      <c-sidebar
        :expanded.sync="expanded"
        :pinned.sync="pinned"
        :icon="icon"
        :logo="logo"
        expand-on-hover
      >
        <template #body-expanded>
          <filters />
        </template>
      </c-sidebar>
    </aside>

    <main class="d-inline-flex h-100">
      <!--
        Content spacer
        Large and xl screens should push in content when the nav is expanded
      -->
      <template>
        <div
          class="sidebar-spacer d-print-none"
          :class="{
            'expanded': expanded && pinned,
          }"
        />
      </template>

      <search />
    </main>
  </div>
</template>

<script>
import Search from '../components/Search.vue'
import Filters from '../components/Filters.vue'
import { components } from '@cortezaproject/corteza-vue'
const { CTopbar, CSidebar } = components

export default {
  i18nOptions: {
    namespaces: 'navigation',
  },

  components: {
    CTopbar,
    CSidebar,
    Search,
    Filters,
  },

  data () {
    return {
      expanded: undefined,
      pinned: undefined,
    }
  },

  computed: {
    user () {
      const { user } = this.$auth
      return user.name || user.handle || user.email || ''
    },

    icon () {
      return this.$Settings.attachment('ui.iconLogo')
    },

    logo () {
      return this.$Settings.attachment('ui.mainLogo')
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
    this.$root.$on('alert', this.displayToast)
  },

  methods: {
    displayToast ({ title, message, variant, countdown }) {
      this.$bvToast.toast(message, {
        title,
        variant,
        solid: true,
        autoHideDelay: countdown,
        toaster: 'b-toaster-bottom-right',
      })
    },
  },
}
</script>
