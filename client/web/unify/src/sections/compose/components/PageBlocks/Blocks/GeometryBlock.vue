<template>
  <PageBlock :block="block" :record="record" @refreshBlock="refresh">
    <div class="flex flex-col h-full relative">
      <div
        v-if="processing"
        class="absolute inset-0 flex items-center justify-center z-10 loading-overlay"
      >
        <ProgressSpinner style="width: 24px; height: 24px" />
      </div>

      <CMap
        :center="mapCenter"
        :zoom="zoomStarting"
        :min-zoom="zoomMin"
        :max-zoom="zoomMax"
        :max-bounds="lockedBounds"
        :markers="markers"
        :polygons="polygons"
        :hide-geo-search="hideGeoSearch"
        class="h-full w-full"
        @marker-click="onMarkerClick"
      />
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, inject, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { compose } from '@planetcrust/human-js'
import { components, mapGeo, useModuleStore, usePageStore } from '@planetcrust/human-vue'
import PageBlock from './PageBlock.vue'
import { evaluatePrefilter, usesRecordVariables } from '../../../lib/record-filter'

const { CMap } = components

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI')
const $Auth = inject('$Auth', {})
const $eventBus = inject('$eventBus', null)
const route = useRoute()
const router = useRouter()
const moduleStore = useModuleStore()
const pageStore = usePageStore()

const processing = ref(false)
const markers = ref([])
const polygons = ref([])

const options = computed(() => props.block.options || {})
const zoomStarting = computed(() => options.value.zoomStarting || 2)
const zoomMin = computed(() => options.value.zoomMin || 0)
const zoomMax = computed(() => options.value.zoomMax || 0)
const hideGeoSearch = computed(() => options.value.hideGeoSearch !== false)

const mapCenter = computed(() => {
  const c = options.value.center
  return Array.isArray(c) && c.length === 2 ? c : [30, 30]
})

// options.bounds is [[swLat, swLng], [neLat, neLng]], and holds the locked area
// the configurator captured — it caps panning, the starting view stays
// center/zoom.
const lockedBounds = computed(() =>
  options.value.lockBounds ? mapGeo.parseBounds(options.value.bounds) : null,
)

function parseCoords(raw) {
  if (!raw) return null
  try {
    const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (parsed?.coordinates?.length === 2) return parsed.coordinates
  } catch {
    // ignore
  }
  return null
}

/**
 * Extract all valid [lat,lng] pairs from a record's geometry field
 * (handles multi-value fields).
 */
function recordPoints(record, feed, module) {
  const field = module.fields.find(f => f.name === feed.geometryField)
  if (!field) return []

  const raw = record.values?.[feed.geometryField]
  const values = field.isMulti && Array.isArray(raw) ? raw : [raw]
  return values.map(parseCoords).filter(Boolean)
}

async function loadFeeds() {
  if (!$ComposeAPI) return

  processing.value = true
  try {
    const outMarkers = []
    const outPolygons = []
    const feeds = options.value.feeds || []

    for (const feed of feeds) {
      if (feed.resource !== 'compose:record') continue
      if (!feed.options?.moduleID || !feed.geometryField) continue

      let module = moduleStore.getByID(feed.options.moduleID)
      if (!module) {
        try {
          await moduleStore.findByID({
            namespaceID: props.namespace.namespaceID,
            moduleID: feed.options.moduleID,
          })
          module = moduleStore.getByID(feed.options.moduleID)
        } catch {
          continue
        }
      }
      if (!module) continue

      let records = []
      try {
        const feedClone = compose.PageBlockGeometry.makeFeed(feed)
        const prefilter = feedClone.options.prefilter

        if (prefilter) {
          const record = props.record
          const user = $Auth?.user || {}

          if (!record && usesRecordVariables(prefilter)) {
            console.warn(
              'Skipping geometry feed: prefilter uses record variables outside a record page',
            )
            continue
          }

          feedClone.options.prefilter = evaluatePrefilter(prefilter, {
            record,
            user,
            recordID: record?.recordID || '0',
            ownerID: record?.ownedBy || '0',
            userID: user?.userID || '0',
          })
        }

        records = await compose.PageBlockGeometry.RecordFeed(
          $ComposeAPI,
          module,
          props.namespace,
          feedClone,
        )
      } catch (e) {
        console.error('Failed to load geometry feed:', e)
        continue
      }

      const color = feed.options?.color
      const feedPoints = []

      for (const r of records) {
        const pts = recordPoints(r, feed, module)
        feedPoints.push(...pts)

        if (feed.displayMarker !== false) {
          for (const coords of pts) {
            outMarkers.push({
              value: coords,
              color,
              title: feed.titleField ? r.values?.[feed.titleField] : undefined,
              recordID: r.recordID,
              moduleID: module.moduleID,
            })
          }
        }
      }

      if (feed.displayPolygon && feedPoints.length >= 3) {
        outPolygons.push({ latLngs: feedPoints, color })
      }
    }

    markers.value = outMarkers
    polygons.value = outPolygons
  } finally {
    processing.value = false
  }
}

/**
 * Marker click → navigate to the record page (same pattern as CalendarBlock).
 */
function onMarkerClick({ marker }) {
  const { recordID, moduleID } = marker || {}
  if (!recordID || !moduleID) return

  const recordPage = (pageStore.set || []).find(p => p.moduleID === moduleID)
  if (!recordPage) return

  const displayOption = options.value.displayOption || 'sameTab'

  if (displayOption === 'modal') {
    router.push({
      query: {
        ...route.query,
        recordPageID: recordPage.pageID,
        recordID,
      },
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

function refresh() {
  loadFeeds()
}

watch(() => options.value, refresh, { deep: true, immediate: true })
watch(() => props.record?.recordID, refresh)

const offRefetch = $eventBus?.on('refetch-records', refresh)

onBeforeUnmount(() => {
  markers.value = []
  polygons.value = []
  offRefetch?.()
})
</script>

<style scoped>
.loading-overlay {
  background: var(--p-mask-background);
}
</style>
