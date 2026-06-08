<template>
  <!-- App items -->
  <TransitionGroup
    v-if="areAppsVisible"
    :class="containerClass"
    tag="div"
    move-class="transition-transform duration-300 ease-in-out"
  >
    <a
      v-for="(app, index) in filteredApps"
      :key="app.applicationID"
      :href="app.enabled ? getAppUrl(app) : '#'"
      target="_self"
      :draggable="canReorder && app.enabled"
      :class="[itemClass, {
        'cursor-grab': canReorder && app.enabled,
        'cursor-not-allowed opacity-50': !app.enabled,
        'border-t-2 !border-t-primary': canReorder && variant === 'list' && dropTargetIndex === index,
        '!border-primary bg-primary/5': variant === 'list' && isActiveApp(app),
      }]"
      @click="onItemClick($event, app)"
      @dragstart="canReorder && app.enabled && onDragStart(index)"
      @dragover="canReorder && app.enabled && onDragOver($event, index)"
      @dragleave="canReorder && onDragLeave()"
      @drop.prevent="canReorder && app.enabled && onDrop(index)"
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
          :class="['group w-80 min-h-72 overflow-hidden transition-all duration-100', app.enabled ? 'cursor-pointer hover:shadow-lg hover:scale-105 hover:text-primary' : 'cursor-not-allowed', { 'ring-2 ring-primary ring-offset-2': canReorder && dropTargetIndex === index }]"
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

  <!-- Empty state -->
  <div v-else :class="emptyClass">
    <span :class="variant === 'list' ? 'text-muted-color text-sm' : 'text-muted-color text-lg'">
      {{ query ? noResultsText : noAppsText }}
    </span>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useRoute } from 'vue-router'
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

// Drag state
const draggedIndex = ref(null)
const dropTargetIndex = ref(null)

function onDragStart(index) {
  draggedIndex.value = index
}

function onDragOver(e, index) {
  e.preventDefault()
  dropTargetIndex.value = index
}

function onDragLeave() {
  dropTargetIndex.value = null
}

function onDrop(index) {
  if (draggedIndex.value === null || draggedIndex.value === index) {
    draggedIndex.value = null
    dropTargetIndex.value = null
    return
  }
  const reordered = [...apps.value]
  const [moved] = reordered.splice(draggedIndex.value, 1)
  reordered.splice(index, 0, moved)
  draggedIndex.value = null
  dropTargetIndex.value = null
  applicationsStore.reorder(reordered)
}

// Layout classes
const containerClass = computed(() =>
  props.variant === 'list'
    ? 'flex flex-col gap-2'
    : 'flex flex-wrap justify-center gap-7',
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
