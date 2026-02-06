<template>
  <div
    class="header-navigation flex flex-wrap items-center py-2 px-3 gap-2 border-b border-surface"
  >
    <!-- Sidebar toggle + small logo -->
    <div class="flex items-center gap-1">
      <!-- When sidebar is disabled on this route, show only icon logo -->
      <img v-if="sidebarDisabled" :src="iconLogo" class="h-8 w-auto object-contain p-2" />

      <!-- Normal mode: hamburger toggle only when collapsed -->
      <template v-else>
        <Button
          v-if="!sidebarExpanded"
          data-test-id="button-sidebar-toggle"
          icon="pi pi-bars"
          variant="text"
          size="large"
          @click="sidebarExpanded = true"
        />
      </template>
    </div>

    <div
      id="topbar-title"
      class="title flex text-truncate items-center text-2xl font-medium mb-0"
    />

    <div id="topbar-tools" class="tools-wrapper ml-auto flex items-center gap-2">
      <slot name="tools" />
    </div>

    <div class="flex items-center gap-1">
      <Button
        v-if="!hideAppSelector && !settings?.hideAppSelector"
        data-test-id="app-selector"
        :href="appSelectorURL"
        severity="secondary"
        variant="text"
        size="small"
        class="mr-[-0.5rem]"
      >
        {{ labels.appMenu }}
      </Button>

      <slot name="right-tools" />

      <div v-if="!settings?.hideHelp" class="help-dropdown">
        <Button
          ref="helpMenuRef"
          data-test-id="dropdown-helper"
          icon="pi pi-question-circle"
          size="large"
          severity="secondary"
          variant="text"
          rounded
          @click="toggleHelpMenu"
        />

        <Menu ref="helpMenu" :model="helpMenuItems" :popup="true" class="mt-2" />
      </div>

      <div v-if="!settings?.hideProfile" class="profile-dropdown">
        <Button
          ref="profileMenuRef"
          data-test-id="dropdown-profile"
          rounded
          variant="outlined"
          icon="pi pi-user"
          severity="secondary"
          size="large"
          class="text-color"
          :class="{ 'p-0': avatar }"
          @click="toggleProfileMenu"
        >
          <template v-if="avatar" #default>
            <Avatar :image="avatar" shape="circle" class="w-full h-full" />
          </template>
        </Button>

        <TieredMenu ref="profileMenu" :model="profileMenuItems" popup />
      </div>
    </div>
  </div>
</template>

<script setup>
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'
import Menu from 'primevue/menu'
import TieredMenu from 'primevue/tieredmenu'
import { computed, inject, ref, watch } from 'vue'

const sidebarExpanded = defineModel('sidebarExpanded', {
  type: Boolean,
  default: false,
})

const props = defineProps({
  sidebarDisabled: {
    type: Boolean,
    default: false,
  },
  hideAppSelector: {
    type: Boolean,
    default: false,
  },
  appSelectorURL: {
    type: String,
    default: '../',
  },
  labels: {
    type: Object,
    required: true,
  },
})

const $Auth = inject('$Auth')
const $Settings = inject('$Settings')

const settings = computed(() => {
  return $Settings.get('ui.topbar', {})
})

const iconLogo = computed(() => {
  return $Settings.attachment('ui.iconLogo')
})

const helpMenuRef = ref()
const helpMenu = ref()
const profileMenuRef = ref()
const profileMenu = ref()

// Make theme reactive by using a ref
const currentTheme = ref($Auth.user?.meta?.theme || '')

// Watch for changes in the auth user's theme
watch(
  () => $Auth.user?.meta?.theme,
  newTheme => {
    currentTheme.value = newTheme || ''
  },
  { immediate: true },
)

const documentationURL = computed(() => {
  // eslint-disable-next-line no-undef
  const [year, month] = VERSION.split('.')
  return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/index.html`
})

const helpLinks = computed(() => {
  const { helpLinks = [] } = props.settings || {}
  return (helpLinks || []).filter(({ handle, url }) => handle && url)
})

const profileLinks = computed(() => {
  const { profileLinks = [] } = props.settings || {}
  return (profileLinks || []).filter(({ handle, url }) => handle && url)
})

const buildVersion = computed(() => {
  // eslint-disable-next-line no-undef
  return VERSION
})

const themes = computed(() => [
  {
    id: 'light',
    label: props.labels.lightTheme,
  },
  {
    id: 'dark',
    label: props.labels.darkTheme,
  },
])

const helpMenuItems = computed(() => {
  const items = []

  helpLinks.value.forEach(helpLink => {
    items.push({
      label: helpLink.handle,
      url: helpLink.url,
      target: helpLink.newTab ? '_blank' : '',
    })
  })

  if (!props.settings?.hideForumLink) {
    items.push({
      label: props.labels.helpForum,
      url: 'https://forum.cortezaproject.org/',
      target: '_blank',
    })
  }

  if (!props.settings?.hideDocumentationLink) {
    items.push({
      label: props.labels.helpDocumentation,
      url: documentationURL.value,
      target: '_blank',
    })
  }

  if (!props.settings?.hideFeedbackLink) {
    items.push({
      label: props.labels.helpFeedback,
      url: 'mailto:info@cortezaproject.org',
      target: '_blank',
    })
  }

  if (items.length > 0) {
    items.push({ separator: true })
  }

  items.push({
    label: buildVersion.value,
    disabled: true,
    class: 'text-sm',
  })

  return items
})

const profileMenuItems = computed(() => {
  const items = []

  if ($Auth.user.name) {
    items.push({
      label: $Auth.user.name,
      disabled: true,
      class: 'font-bold',
    })
  }

  if ($Auth.user.email) {
    items.push({
      label: $Auth.user.email,
      disabled: true,
      class: 'text-sm text-muted-color mb-2 -mt-2',
    })
  }

  profileLinks.value.forEach(profileLink => {
    items.push({
      label: profileLink.handle,
      url: profileLink.url,
      target: profileLink.newTab ? '_blank' : '',
    })
  })

  if (!props.settings?.hideProfileLink) {
    items.push({
      label: props.labels.userSettingsProfile,
      url: $Auth.cortezaAuthURL,
      target: '_blank',
      icon: 'pi pi-user',
    })
  }

  if (!props.settings?.hideChangePasswordLink) {
    items.push({
      label: props.labels.userSettingsChangePassword,
      url: `${$Auth.cortezaAuthURL}/change-password`,
      target: '_blank',
      icon: 'pi pi-key',
    })
  }

  if (!props.settings?.hideThemeSelector) {
    items.push({
      label: props.labels.userSettingsTheme,
      items: themes.value.map(theme => ({
        label: theme.label,
        disabled: currentTheme.value === theme.id,
        icon: `pi pi-${theme.id === 'light' ? 'sun' : 'moon'}`,
        command: () => changeTheme(theme.id),
      })),
      icon: 'pi pi-palette',
    })
  }

  items.push({ separator: true })

  items.push({
    label: props.labels.userSettingsLogout,
    icon: 'pi pi-sign-out',
    command: () => logout(),
  })

  return items
})

const avatar = computed(() => {
  return `${$SystemAPI.baseURL}/attachment/avatar/${$Auth.user.meta.avatarID}/original/profile-photo-avatar`
})

const toggleHelpMenu = event => {
  helpMenu.value.toggle(event)
}

const toggleProfileMenu = event => {
  profileMenu.value.toggle(event)
}

const $SystemAPI = inject('$SystemAPI')

import { useTheme } from '../../composables/useTheme'

const changeTheme = theme => {
  $Auth.user.meta.theme = theme
  currentTheme.value = theme
  useTheme(theme)
  $SystemAPI.userUpdate($Auth.user)
}

const logout = () => {
  $Auth.logout()
}
</script>

<style scoped>
.header-navigation {
  width: 100%;
  min-height: var(--topbar-height);
  background-color: var(--topbar-bg);
}

.nav-icon {
  width: calc(var(--topbar-height) - 24px);
  height: calc(var(--topbar-height) - 24px);
}

.title {
  min-height: calc(var(--topbar-height) - 15px);
}

.tools-wrapper > :deep(*) {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  flex-wrap: wrap;
}
</style>
