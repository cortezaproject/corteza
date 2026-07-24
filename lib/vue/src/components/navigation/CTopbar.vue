<template>
  <div class="header-navigation flex flex-wrap items-center p-2 !pr-3">
    <!-- Sidebar toggle + small logo -->
    <div v-if="!hideLogo && !sidebarExpanded">
      <!-- When sidebar is disabled on this route, show only icon logo -->
      <img v-if="sidebarDisabled" :src="iconLogo" class="h-10 w-10 p-2 object-contain" />

      <!-- Normal mode: hamburger toggle only when collapsed -->
      <template v-else>
        <Button
          data-test-id="button-sidebar-toggle"
          severity="secondary"
          icon="pi pi-bars"
          variant="text"
          rounded
          @click="sidebarExpanded = true"
        />
      </template>
    </div>

    <!-- grow/shrink-0: the title keeps its natural width so the tools + right
         cluster wrap below it (root flex-wrap) instead of squeezing it into
         truncation; max-w-full still caps pathological titles (truncate). -->
    <div class="topbar-title-area flex grow shrink-0 basis-auto max-w-full items-center gap-2 ml-2">
      <div
        id="topbar-title"
        class="topbar-title flex items-center min-w-0 overflow-hidden whitespace-nowrap truncate text-2xl font-medium text-color mb-0"
      />

      <div v-if="visiblePageButtons.length" class="flex items-center gap-2">
        <a
          v-for="(btn, i) in visiblePageButtons"
          :key="i"
          :href="btn.url"
          :target="btn.newTab ? '_blank' : '_self'"
          rel="noopener noreferrer"
          class="no-underline shrink-0"
          @click="!btn.newTab && onAnchorClick($event, btn.url)"
        >
          <Button
            v-tooltip.bottom="btn.description || undefined"
            :label="btn.label"
            severity="secondary"
            outlined
            size="small"
          />
        </a>
      </div>
    </div>

    <!-- One wrapper so tools + right icons wrap below the title as a single
         right-aligned unit when the row runs out of room. -->
    <div class="ml-auto flex items-center gap-2">
      <div id="topbar-tools" class="topbar-tools tools-wrapper flex items-center gap-2">
        <slot name="tools" />
      </div>

      <div class="topbar-right flex items-center gap-1">
        <a
          v-if="!settings?.hideHomeButton"
          :href="homeURL"
          class="no-underline"
          @click="onAnchorClick($event, homeURL)"
        >
          <Button
            v-tooltip.bottom="labels.home || 'Home'"
            icon="pi pi-home"
            severity="secondary"
            variant="text"
            rounded
          />
        </a>

        <Button
          v-if="!hideAppSelector && !settings?.hideAppSelector"
          v-tooltip.bottom="labels.appMenu"
          data-test-id="app-selector"
          icon="pi pi-th-large"
          severity="secondary"
          variant="text"
          rounded
          @click="onAppMenuClick"
        />

        <slot name="right-tools" />

        <CAgentSidebarButton v-if="!settings?.hideAgentSidebar" />
        <CNotificationButton v-if="!settings?.hideNotifications" />

        <div v-if="!settings?.hideHelp" class="help-dropdown">
          <Button
            ref="helpMenuRef"
            data-test-id="dropdown-helper"
            icon="pi pi-dollar"
            severity="success"
            variant="text"
            rounded
            @click="toggleHelpMenu"
          />

          <Menu ref="helpMenu" :model="helpMenuItems" :popup="true" class="mt-2" />
        </div>

        <div v-if="!settings?.hideProfile" class="flex">
          <Button
            ref="profileMenuRef"
            data-test-id="dropdown-profile"
            rounded
            variant="outlined"
            severity="secondary"
            size="large"
            class="text-color !p-0 !w-10 !h-10"
            @click="toggleProfileMenu"
          >
            <template #default>
              <Avatar
                :image="avatar || undefined"
                :label="!avatar ? userInitials : undefined"
                :icon="!avatar && !userInitials ? 'pi pi-user' : undefined"
                shape="circle"
                class="!w-full !h-full !text-sm"
              />
            </template>
          </Button>

          <TieredMenu ref="profileMenu" :model="profileMenuItems" popup />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useInternalLink } from '../../composables/useInternalLink'
import CNotificationButton from '../notifications/CNotificationButton.vue'
import CAgentSidebarButton from '../agent/CAgentSidebarButton.vue'

const { onAnchorClick } = useInternalLink()

const sidebarExpanded = defineModel('sidebarExpanded', {
  type: Boolean,
  default: false,
})

const props = defineProps({
  settings: {
    type: Object,
    default: () => ({}),
  },
  hideLogo: {
    type: Boolean,
    default: false,
  },
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
    default: '/',
  },
  labels: {
    type: Object,
    required: true,
  },
  customProfileItems: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['app-menu-click'])

const onAppMenuClick = () => {
  emit('app-menu-click')
}

const $Auth = inject('$Auth')
const $Settings = inject('$Settings')

const settings = computed(() => {
  return {
    ...$Settings.get('ui.topbar', {}),
    ...props.settings,
  }
})

// Page buttons: filter by URL substring match, kept reactive for SPA navigation
const currentHref = ref(window.location.href)

const updateHref = () => {
  currentHref.value = window.location.href
}

// Listen for browser back/forward
onMounted(() => {
  window.addEventListener('popstate', updateHref)

  // Patch pushState/replaceState to detect SPA route changes
  const origPush = history.pushState.bind(history)
  const origReplace = history.replaceState.bind(history)
  history.pushState = (...args) => {
    origPush(...args)
    updateHref()
  }
  history.replaceState = (...args) => {
    origReplace(...args)
    updateHref()
  }

  // Store originals so we can restore later
  window.__pageButtonsOrigPush = origPush
  window.__pageButtonsOrigReplace = origReplace
})

onBeforeUnmount(() => {
  window.removeEventListener('popstate', updateHref)
  // Restore original history methods
  if (window.__pageButtonsOrigPush) history.pushState = window.__pageButtonsOrigPush
  if (window.__pageButtonsOrigReplace) history.replaceState = window.__pageButtonsOrigReplace
})

const homeURL = computed(() => window.location.origin)

const visiblePageButtons = computed(() => {
  const buttons = settings.value?.pageButtons || []
  const href = currentHref.value
  let pathname = ''
  try {
    pathname = new URL(href).pathname
  } catch {
    /* ignore */
  }
  const isRoot = pathname === '' || pathname === '/'
  return buttons.filter(btn => {
    if (!btn.label || !btn.url) return false
    if (!btn.urlMatch) return isRoot
    return href.includes(btn.urlMatch)
  })
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

const profileLinks = computed(() => {
  const { profileLinks = [] } = settings.value || {}
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

  items.push({
    label: props.labels.helpBuyHuman,
    url: 'https://buy.polar.sh/polar_cl_MKycmV54KogEiCMy7Oqf4zXIwWZ7wOPhtf5xg3dTn21',
    target: '_blank',
  })

  items.push({
    label: props.labels.helpManageSubscription,
    url: 'https://polar.sh/planet-crust/portal/',
    target: '_blank',
  })

  items.push({
    label: props.labels.helpPackageDetails,
    url: 'https://docs.planetcrust.com/price-and-terms',
    target: '_blank',
  })

  items.push({ separator: true })

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

  const hasBuiltInProfileItem =
    !settings.value?.hideProfileLink ||
    !settings.value?.hideChangePasswordLink ||
    !settings.value?.hideThemeSelector

  if (profileLinks.value.length && hasBuiltInProfileItem) {
    items.push({ separator: true })
  }

  if (!settings.value?.hideProfileLink) {
    items.push({
      label: props.labels.userSettingsProfile,
      url: $Auth.authURL,
      target: '_blank',
      icon: 'pi pi-user',
    })
  }

  if (!settings.value?.hideChangePasswordLink) {
    items.push({
      label: props.labels.userSettingsChangePassword,
      url: `${$Auth.authURL}/change-password`,
      target: '_blank',
      icon: 'pi pi-key',
    })
  }

  if (!settings.value?.hideThemeSelector) {
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

  if (props.customProfileItems.length) {
    items.push({ separator: true })
    items.push(...props.customProfileItems)
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
  const avatarID = $Auth.user?.meta?.avatarID
  if (!avatarID || avatarID === '0') return ''
  return `${$SystemAPI.baseURL}/attachment/avatar/${avatarID}/original/profile-photo-avatar`
})

const userInitials = computed(() => {
  const name = $Auth.user?.name
  if (name) {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }
  const email = $Auth.user?.email
  if (email) return email[0].toUpperCase()
  return ''
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
  container-type: inline-size;
}

.tools-wrapper > :deep(*) {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  flex-wrap: wrap;
}

/* Mobile trims only — the title is never force-moved; it wraps content-driven
   like every other row (see the template classes). Note the query measures the
   topbar's own width (window minus sidebar), not the viewport. */
@container (max-width: 550px) {
  .header-navigation {
    row-gap: 0.25rem;
  }

  .topbar-title {
    font-size: 1.25rem;
  }

  .topbar-tools {
    display: none !important;
  }
}
</style>
