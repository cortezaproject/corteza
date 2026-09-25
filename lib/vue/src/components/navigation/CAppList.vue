<template>
  <!-- App items -->
  <CDraggableList
    v-if="areAppsVisible"
    v-model="ordered"
    :class="containerClass"
    :disabled="!canReorder"
    filter="[data-app-disabled], button"
  >
    <!-- Tagless, so the tiles stay direct children of the sortable while
         still animating to their new places once a drop settles. -->
    <TransitionGroup move-class="transition-transform duration-300 ease-in-out">
      <div
        v-for="app in filteredApps"
        :key="app.applicationID"
        data-drag-item
        :data-app-disabled="isOpenable(app) ? undefined : ''"
        class="relative group/tile"
      >
        <a
          :href="isOpenable(app) ? getAppUrl(app) : '#'"
          target="_self"
          :class="[
            itemClass,
            {
              'cursor-pointer': isOpenable(app),
              'cursor-not-allowed opacity-50': !isOpenable(app),
              '!border-primary bg-primary/5': variant === 'list' && isActiveApp(app),
            },
          ]"
          @click="onItemClick($event, app)"
        >
          <!-- List variant -->
          <template v-if="variant === 'list'">
            <img
              :src="getAppLogoUrl(app)"
              :alt="app.unify?.name || app.name"
              class="w-16 h-16 object-contain rounded-md shrink-0"
              loading="lazy"
            />
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="font-medium text-[0.9375rem] leading-snug break-words whitespace-normal">
                <template v-if="app.applicationID === applicationsStore.homeID">
                  {{ nameParts(app).head }}
                  <span class="whitespace-nowrap">
                    {{ nameParts(app).tail }}
                    <i
                      class="pi pi-home text-[0.85em] text-muted-color"
                      data-test-id="app-tile-home"
                      role="img"
                      :title="homeBadge"
                      :aria-label="homeBadge"
                    />
                  </span>
                </template>
                <template v-else>{{ app.unify?.name || app.name }}</template>
              </span>
              <span
                v-if="app.meta?.description"
                class="text-xs text-muted-color break-words whitespace-normal"
                data-test-id="app-tile-description"
              >
                {{ app.meta.description }}
              </span>
            </div>
          </template>

          <!-- Grid variant -->
          <template v-else>
            <Card
              :pt="{
                body: { class: 'grow justify-center gap-0 py-1' },
                title: { class: 'text-center line-clamp-2 group-hover:line-clamp-none' },
              }"
              :class="[
                'group w-80 min-h-72 overflow-hidden transition-all duration-100',
                isOpenable(app)
                  ? 'cursor-pointer hover:shadow-lg hover:scale-105 hover:text-primary'
                  : 'cursor-not-allowed',
              ]"
            >
              <template #header>
                <img
                  :src="getAppLogoUrl(app)"
                  :alt="app.unify?.name || app.name"
                  class="w-full h-full object-contain"
                  loading="lazy"
                  decoding="async"
                />
              </template>
              <template #title>
                <template v-if="app.applicationID === applicationsStore.homeID">
                  {{ nameParts(app).head }}
                  <span class="whitespace-nowrap">
                    {{ nameParts(app).tail }}
                    <i
                      class="pi pi-home text-[0.85em] text-muted-color"
                      data-test-id="app-tile-home"
                      role="img"
                      :title="homeBadge"
                      :aria-label="homeBadge"
                    />
                  </span>
                </template>
                <template v-else>{{ app.unify?.name || app.name }}</template>
              </template>
              <template v-if="app.meta?.description" #subtitle>
                <span class="block text-center break-words" data-test-id="app-tile-description">
                  {{ app.meta.description }}
                </span>
              </template>
            </Card>
          </template>
        </a>
        <Button
          v-if="canPickHome"
          icon="pi pi-ellipsis-v"
          text
          size="small"
          severity="secondary"
          :aria-label="labels.tileMenu"
          data-test-id="app-tile-menu"
          class="!absolute top-1 right-1 opacity-50 transition-opacity group-hover/tile:opacity-100 focus-visible:opacity-100"
          @click="openTileMenu($event, app)"
        />
      </div>
    </TransitionGroup>
  </CDraggableList>

  <!-- Empty state -->
  <div v-else :class="emptyClass">
    <span :class="variant === 'list' ? 'text-muted-color text-sm' : 'text-muted-color text-lg'">
      {{ query ? noResultsText : noAppsText }}
    </span>
  </div>

  <Menu ref="tileMenu" :model="tileMenuItems" popup />
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useRoute } from 'vue-router'
import Button from 'primevue/button'
import Menu from 'primevue/menu'
import { useToast } from 'primevue/usetoast'
import CDraggableList from '../drag/CDraggableList.vue'
import { useRBACStore } from '../../composables/useRBAC'
import { useInternalLink } from '../../composables/useInternalLink'
import { useApplicationsStore } from '../../stores/useApplicationsStore'
import { resolveAppLogoUrl } from '../../utils/appIcons'
import { resolveAppUrl } from '../../utils/appUrl'

const $appIconMap = inject('$appIconMap', {})

// The host app decides which entries it can actually open; without one every
// entry is openable, which is what a registry with no access control means.
const $appReachable = inject('$appReachable', null)

// An entry the user cannot open reads the same as a switched-off one: shown,
// greyed, inert. Hiding it instead would leave the menu silently shorter than
// the registry, with nothing saying why.
const isOpenable = app => app.enabled && ($appReachable ? $appReachable(app) : true)

const props = defineProps({
  query: {
    type: String,
    default: '',
  },
  variant: {
    type: String,
    default: 'list', // 'list' | 'grid'
  },
  noAppsText: {
    type: String,
    default: 'No applications',
  },
  noResultsText: {
    type: String,
    default: 'No applications found',
  },
  labels: {
    type: Object,
    default: () => ({
      tileMenu: 'Application options',
      setHome: 'Set as Home',
      clearHome: 'Remove as Home',
      homeSetTitle: 'Home page set',
      homeSet: '{app} now opens when you sign in or select Home.',
      homeClearedTitle: 'Home page removed',
      homeCleared: '{app} no longer opens when you sign in or select Home.',
      homeErrorTitle: 'Your home page could not be changed',
      homeOwn: 'Your home page',
      homeGlobal: 'Home page for everyone',
    }),
  },
})

const $SystemAPI = inject('$SystemAPI')
const applicationsStore = useApplicationsStore()
const route = useRoute()

const isActiveApp = app => {
  const url = app.unify?.url
  if (!url) return false
  const base = '/' + url.toLowerCase().replace(/^\/|\/$/g, '')
  const path = route.path.toLowerCase()
  return path === base || path.startsWith(base + '/')
}

const { onAnchorClick } = useInternalLink()

// Left-click on an app the current SPA hosts internally routes client-side
// (no reload, no splash). Anything else (other app, external, modified click)
// falls through to the anchor's native navigation, preserving prior behavior.
const onItemClick = (event, app) => {
  if (!isOpenable(app)) {
    event.preventDefault()
    return
  }
  onAnchorClick(event, getAppUrl(app))
}

const apps = computed(() => applicationsStore.unifyOnly)
const canReorder = computed(() => applicationsStore.apps.some(a => a.canUpdateApplication))

const normalizedQuery = computed(() => (props.query || '').trim().toUpperCase())

const isAppVisible = app => {
  const q = normalizedQuery.value
  if (!q) return true
  return (
    (app.name?.toUpperCase() || '').includes(q) ||
    (app.unify?.name?.toUpperCase() || '').includes(q)
  )
}

const filteredApps = computed(() => apps.value.filter(isAppVisible))
const areAppsVisible = computed(() => filteredApps.value.length > 0)

const getAppLogoUrl = app => resolveAppLogoUrl(app, $SystemAPI.baseURL, $appIconMap)

const getAppUrl = app => resolveAppUrl(app.unify?.url)

// The per-tile menu picks the user's own home application.
const toast = useToast()
const rbac = useRBACStore()
const canPickHome = computed(() => rbac.can('system/', 'application.flag.self'))

const tileMenu = ref(null)
const tileMenuApp = ref(null)

const openTileMenu = (event, app) => {
  tileMenuApp.value = app
  tileMenu.value.toggle(event)
}

// The home badge rides with the name's last word, so a wrapping name never
// leaves it alone on a line.
const nameParts = app => {
  const name = (app.unify?.name || app.name || '').trim()
  const at = name.lastIndexOf(' ')
  return { head: at === -1 ? '' : name.slice(0, at), tail: name.slice(at + 1) }
}

// The badge names whose pick the home application is.
const homeBadge = computed(() =>
  applicationsStore.ownHomeID ? props.labels.homeOwn : props.labels.homeGlobal,
)

const pickHome = async (app, set) => {
  const name = app.unify?.name || app.name
  try {
    await applicationsStore.setOwnHome(set ? app.applicationID : '')
    toast.add({
      severity: 'success',
      summary: set ? props.labels.homeSetTitle : props.labels.homeClearedTitle,
      detail: (set ? props.labels.homeSet : props.labels.homeCleared).replace('{app}', name),
      life: 4000,
    })
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: props.labels.homeErrorTitle,
      detail: error?.message || String(error),
      life: 6000,
    })
  }
}

const tileMenuItems = computed(() => {
  const app = tileMenuApp.value
  if (!app) return []

  const isHome = applicationsStore.ownHomeID === app.applicationID
  return [
    {
      label: isHome ? props.labels.clearHome : props.labels.setHome,
      icon: isHome ? 'pi pi-times' : 'pi pi-home',
      command: () => pickHome(app, !isHome),
    },
  ]
})

// The list being dragged is the filtered one, but the order being saved is the
// whole one. The visible apps are put back into the slots they already occupy,
// in their new sequence, so an app hidden by the search keeps its place instead
// of being displaced by an index that never referred to it.
const ordered = computed({
  get: () => filteredApps.value,
  set: next => {
    const slots = []
    apps.value.forEach((app, i) => {
      if (isAppVisible(app)) slots.push(i)
    })

    const full = [...apps.value]
    slots.forEach((slot, i) => {
      full[slot] = next[i]
    })

    applicationsStore.reorder(full)
  },
})

// Layout classes
const containerClass = computed(() =>
  props.variant === 'list' ? 'flex flex-col gap-2' : 'flex flex-wrap justify-center gap-7',
)

const itemClass = computed(() =>
  props.variant === 'list'
    ? 'flex items-center gap-3 p-3 pr-10 rounded-lg border border-surface hover:bg-emphasis hover:border-primary transition-all duration-150 no-underline text-color'
    : 'block',
)

const emptyClass = computed(() =>
  props.variant === 'list'
    ? 'flex items-center justify-center'
    : 'flex justify-center items-center mt-20 w-full',
)
</script>
