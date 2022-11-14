<template>
  <div class="d-flex flex-column w-100 vh-100">
    <header>
      <c-topbar
        :sidebar-pinned="pinned"
        :settings="$Settings.get('ui.topbar', {})"
        :labels="{
          helpForum: $t('navigation:help.forum'),
          helpDocumentation: $t('navigation:help.documentation'),
          helpFeedback: $t('navigation:help.feedback'),
          helpVersion: $t('navigation:help.version'),
          userSettingsLoggedInAs: $t('navigation:userSettings.loggedInAs', { user }),
          userSettingsProfile: $t('navigation:userSettings.profile'),
          userSettingsChangePassword: $t('navigation:userSettings.changePassword'),
          userSettingsLogout: $t('navigation:userSettings.logout'),
        }"
      >
        <template #title>
          <portal-target
            name="topbar-title"
          />
        </template>

        <template #tools>
          <portal-target
            name="topbar-tools"
          />
        </template>
      </c-topbar>
    </header>

    <aside>
      <c-sidebar
        :expanded.sync="expanded"
        :pinned.sync="pinned"
        :icon="icon"
        :logo="logo"
        :disabled-routes="['dashboard']"
      >
        <template #header-expanded>
          <portal-target name="sidebar-header-expanded" />
        </template>

        <template #body-expanded>
          <portal-target name="sidebar-body-expanded" />
        </template>

        <template #footer-expanded>
          <portal-target name="sidebar-footer-expanded" />
        </template>
      </c-sidebar>
    </aside>

    <main class="d-inline-flex h-100 overflow-auto">
      <!--
        Content spacer
        Large and xl screens should push in content when the nav is expanded
      -->
      <template>
        <div
          class="spacer"
          :class="{
            'expanded': expanded && pinned,
          }"
        />
      </template>

      <div
        class="d-flex flex-column w-100"
      >
        <router-view
          class="flex-grow-1 overflow-auto"
        />

        <portal-target
          name="editor-toolbar"
        />
      </div>
    </main>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CTopbar, CSidebar } = components

export default {
  name: 'Layout',

  components: {
    CTopbar,
    CSidebar,
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
}
</script>

<style lang="scss" scoped>
.spacer {
  min-width: 0;
  -webkit-transition: min-width 0.2s ease-in-out;
  -moz-transition: min-width 0.2s ease-in-out;
  -o-transition: min-width 0.2s ease-in-out;
  transition: min-width 0.2s ease-in-out;

  &.expanded {
    min-width: $sidebar-width;
    -webkit-transition: min-width 0.2s ease-in-out;
    -moz-transition: min-width 0.2s ease-in-out;
    -o-transition: min-width 0.2s ease-in-out;
    transition: min-width 0.2s ease-in-out;
  }
}
</style>
