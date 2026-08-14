<template>
  <!-- App items -->
  <CDraggableList
    v-if="areAppsVisible"
    v-model="ordered"
    :class="containerClass"
    :disabled="!canReorder"
    filter="[data-app-disabled]"
  >
    <!-- Tagless, so the anchors stay direct children of the sortable while
         still animating to their new places once a drop settles. -->
    <TransitionGroup move-class="transition-transform duration-300 ease-in-out">
      <a
        v-for="app in filteredApps"
        :key="app.applicationID"
        data-drag-item
        :data-app-disabled="app.enabled ? undefined : ''"
        :href="app.enabled ? getAppUrl(app) : '#'"
        target="_self"
        :class="[
          itemClass,
          {
            'cursor-grab': canReorder && app.enabled,
            'cursor-not-allowed opacity-50': !app.enabled,
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
          <span class="font-medium text-sm truncate">
            {{ app.unify?.name || app.name }}
          </span>
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
              app.enabled
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
              {{ app.unify?.name || app.name }}
            </template>
          </Card>
        </template>
      </a>
    </TransitionGroup>
  </CDraggableList>

  <!-- Empty state -->
  <div v-else :class="emptyClass">
    <span :class="variant === 'list' ? 'text-muted-color text-sm' : 'text-muted-color text-lg'">
      {{ query ? noResultsText : noAppsText }}
    </span>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useRoute } from 'vue-router'
import CDraggableList from '../drag/CDraggableList.vue'
import { useInternalLink } from '../../composables/useInternalLink'
import { useApplicationsStore } from '../../stores/useApplicationsStore'
import { resolveAppLogoUrl } from '../../utils/appIcons'

const $appIconMap = inject('$appIconMap', {})

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
  if (!app.enabled) {
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

const getAppUrl = app => {
  const url = app.unify?.url || ''
  if (!url || url.startsWith('/') || url.startsWith('http')) return url
  return '/' + url
}

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
    ? 'flex items-center gap-3 p-3 rounded-lg border border-surface hover:bg-emphasis hover:border-primary transition-all duration-150 no-underline text-color'
    : 'block',
)

const emptyClass = computed(() =>
  props.variant === 'list'
    ? 'flex items-center justify-center'
    : 'flex justify-center items-center mt-20 w-full',
)
</script>
