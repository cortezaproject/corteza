<template>
  <PageBlock :block="block" @refreshBlock="refresh">
    <div class="flex flex-col h-full p-2 calendar-container">
      <!-- Custom Header -->
      <div v-if="!header.hide">
        <div
          v-if="!header.hidePrevNext || !header.hideTitle"
          class="flex items-baseline justify-center mb-2 gap-2"
        >
          <Button
            v-if="!header.hidePrevNext"
            icon="pi pi-angle-left"
            text
            size="small"
            @click="calendarApi?.prev()"
          />
          <span
            v-if="!header.hideTitle"
            class="text-xl font-semibold"
          >
            {{ title }}
          </span>
          <Button
            v-if="!header.hidePrevNext"
            icon="pi pi-angle-right"
            text
            size="small"
            @click="calendarApi?.next()"
          />
        </div>

        <div class="grid grid-cols-12 gap-2 mb-2">
          <div class="col-span-12 sm:col-span-9 flex flex-wrap gap-1">
            <Button
              v-for="view in views"
              :key="view"
              :label="$t(`block.calendar.view.${view}`)"
              size="small"
              severity="secondary"
              @click="changeView(view)"
            />
          </div>
          <div
            v-if="!header.hideToday"
            class="col-span-12 sm:col-span-3 flex justify-end"
          >
            <Button
              :label="$t('block.calendar.today')"
              size="small"
              severity="secondary"
              class="w-full"
              @click="calendarApi?.today()"
            />
          </div>
        </div>
      </div>

      <!-- Calendar -->
      <div ref="containerRef" class="flex-1 min-h-0 overflow-hidden relative">
        <div
          v-if="processing"
          class="absolute inset-0 flex items-center justify-center z-10 bg-surface-0/50"
        >
          <ProgressSpinner style="width: 24px; height: 24px" />
        </div>

        <FullCalendar
          ref="calendarRef"
          class="h-full"
          :options="calendarOptions"
        />
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, inject, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { compose } from '@planetcrust/human-js'
import { useModuleStore } from '@/stores/module'
import { usePageStore } from '@/stores/page'
import PageBlock from './PageBlock.vue'
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import timeGridPlugin from '@fullcalendar/timegrid'
import listPlugin from '@fullcalendar/list'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI')
const $eventBus = inject('$eventBus', null)
const route = useRoute()
const router = useRouter()
const moduleStore = useModuleStore()
const pageStore = usePageStore()

const calendarRef = ref(null)
const containerRef = ref(null)
const processing = ref(false)
const title = ref('')
const currentView = ref('')
const events = ref([])

const loaded = ref({ start: null, end: null })
const refreshing = ref(false)

const options = computed(() => props.block.options || {})

const header = computed(() => options.value.header || {})

const views = computed(() => {
  if (header.value.hide) return []
  const h = header.value
  return props.block.reorderViews
    ? props.block.reorderViews(h.views)
    : (h.views || ['dayGridMonth', 'timeGridWeek', 'timeGridDay', 'listMonth'])
})

const calendarApi = computed(() => {
  return calendarRef.value?.getApi?.()
})

const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, timeGridPlugin, listPlugin],
  initialView: options.value.defaultView || 'dayGridMonth',
  headerToolbar: false,
  editable: false,
  dayMaxEvents: true,
  events: events.value,
  height: '100%',
  expandRows: true,
  handleWindowResize: true,
  datesSet: onDatesSet,
  eventClick: handleEventClick,
}))

/**
 * Called when the calendar date range changes (navigation, view change).
 */
function onDatesSet (info) {
  title.value = info.view.title
  currentView.value = info.view.type
  loadEvents(info.start, info.end)
}

/**
 * Changes the current calendar view.
 */
function changeView (view) {
  calendarApi.value?.changeView(view)
}

/**
 * Loads events for all feeds within the given date range.
 */
async function loadEvents (start, end) {
  if (!start || !end) return
  if (!$ComposeAPI) return

  // Skip if same range and not refreshing
  if (
    loaded.value.start &&
    loaded.value.end &&
    start.getTime() === loaded.value.start.getTime() &&
    end.getTime() === loaded.value.end.getTime() &&
    !refreshing.value
  ) {
    return
  }

  loaded.value = { start, end }
  processing.value = true

  try {
    const allEvents = []
    const feeds = options.value.feeds || []

    for (const feed of feeds) {
      if (feed.resource === 'compose:record' && feed.options?.moduleID) {
        const mod = moduleStore.getByID(feed.options.moduleID)
        if (!mod) {
          // Try to load the module
          try {
            await moduleStore.findByID({
              namespaceID: props.namespace.namespaceID,
              moduleID: feed.options.moduleID,
            })
          } catch {
            continue
          }
        }

        const module = moduleStore.getByID(feed.options.moduleID)
        if (!module) continue

        try {
          // Clone feed to avoid mutating the original (prefilter interpolation)
          const feedClone = compose.PageBlockCalendar.makeFeed(feed)
          const feedEvents = await compose.PageBlockCalendar.RecordFeed(
            $ComposeAPI,
            module,
            props.namespace,
            feedClone,
            { start, end },
          )
          allEvents.push(...feedEvents)
        } catch (e) {
          console.error('Failed to load calendar feed:', e)
        }
      }
    }

    events.value = allEvents
  } catch (e) {
    console.error('Calendar load error:', e)
  } finally {
    setTimeout(() => {
      processing.value = false
      refreshing.value = false
    }, 300)
  }
}

/**
 * Handles click on a calendar event — navigates to the record page.
 */
function handleEventClick ({ event }) {
  const { recordID, moduleID } = event.extendedProps || {}
  if (!moduleID || !recordID) return

  // Find the page for this module
  const pages = pageStore.set || []
  const recordPage = pages.find(p => p.moduleID === moduleID)
  if (!recordPage) return

  const displayOption = options.value.eventDisplayOption || 'sameTab'

  if (displayOption === 'modal') {
    router.push({
      query: {
        ...route.query,
        recordPageID: recordPage.pageID,
        recordID: recordID,
      }
    })
    return
  }

  const routeObj = { name: 'page.record', params: { recordID, pageID: recordPage.pageID } }

  if (displayOption === 'newTab') {
    window.open(router.resolve(routeObj).href)
  } else {
    router.push(routeObj)
  }
}

// Refresh when options change (e.g. feeds updated in configurator)
watch(
  () => options.value,
  () => {
    refreshing.value = true
    const api = calendarApi.value
    if (api) {
      loadEvents(api.view.activeStart, api.view.activeEnd)
    }
  },
  { deep: true },
)

// Refresh when record changes (record page context)
watch(
  () => props.record?.recordID,
  () => {
    refreshing.value = true
    const api = calendarApi.value
    if (api) {
      loadEvents(api.view.activeStart, api.view.activeEnd)
    }
  },
)

let resizeObserver = null

onMounted(() => {
  nextTick(() => {
    const api = calendarApi.value
    if (api) {
      title.value = api.view.title
      currentView.value = api.view.type
    }
  })

  // grid-layout-plus resizes the block via CSS transforms — no window resize fires.
  // Observe container size so FullCalendar can recompute its internal layout.
  nextTick(() => {
    if (!containerRef.value) return
    resizeObserver = new ResizeObserver(() => {
      calendarApi.value?.updateSize()
    })
    resizeObserver.observe(containerRef.value)
  })
})

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
  resizeObserver = null
  events.value = []
  loaded.value = { start: null, end: null }
  offRefetch?.()
})

function refresh() {
  refreshing.value = true
  const api = calendarApi.value
  if (api) {
    loadEvents(api.view.activeStart, api.view.activeEnd)
  }
}

const offRefetch = $eventBus?.on('refetch-records', refresh)
</script>

<style>
.calendar-container .fc {
  font-size: 0.875rem;
}

/* Borders and background using PrimeVue surface colors */
.calendar-container .fc th {
  border-color: var(--p-content-border-color);
}

.calendar-container .fc td {
  border-color: var(--p-content-border-color);
}

.calendar-container .fc .fc-scrollgrid {
  border-color: var(--p-content-border-color);
}

/* Column header styling */
.calendar-container .fc .fc-col-header-cell {
  white-space: pre-wrap;
  background: var(--p-content-hover-background);
  font-weight: 600;
  padding: 0.5rem 0;
}

/* Today highlight matches Human's light blue */
.calendar-container .fc .fc-day-today {
  background: var(--p-highlight-background) !important;
}

/* Event cursor */
.calendar-container .fc .fc-daygrid-event,
.calendar-container .fc .fc-timegrid-event,
.calendar-container .fc .fc-list-event {
  cursor: pointer;
}

/* "More" events popover */
.fc-popover {
  border-color: var(--p-content-border-color);
  box-shadow: var(--p-overlay-popover-shadow);
}

.fc-popover .fc-popover-header {
  padding: 0.5rem;
  background: var(--p-content-hover-background);
}

.fc-popover .fc-popover-body {
  padding: 0.25rem;
}
</style>
