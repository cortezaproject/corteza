<template>
  <div class="flex flex-col gap-5">
    <Fieldset :legend="$t('block.geometry.viewLabel')">
      <div class="flex flex-col gap-3">
        <CFormGroup
          :label="$t('block.geometry.startingView')"
          :description="
            lockBounds ? $t('block.geometry.boundsLockedHint') : $t('block.geometry.mapHelpText')
          "
        >
          <CMap
            :center="center"
            :zoom="zoomStarting"
            :min-zoom="zoomMin"
            :max-zoom="zoomMax"
            :max-bounds="lockedBounds"
            :polygons="boundsOutline"
            :disable-pan="lockBounds"
            :hide-geo-search="hideGeoSearch"
            hide-current-location-button
            style="height: 40vh"
            @update:center="onMapCenter"
            @update:zoom="onMapZoom"
            @update:bounds="onMapBounds"
          />
        </CFormGroup>

        <div class="grid grid-cols-2 gap-2">
          <CFormGroup :label="$t('block.geometry.zoomMin')">
            <InputNumber
              :model-value="zoomMin"
              :min="1"
              :max="20"
              show-buttons
              class="w-full"
              @input="e => (zoomMin = e.value)"
            />
          </CFormGroup>
          <CFormGroup :label="$t('block.geometry.zoomMax')">
            <InputNumber
              :model-value="zoomMax"
              :min="1"
              :max="20"
              show-buttons
              class="w-full"
              @input="e => (zoomMax = e.value)"
            />
          </CFormGroup>
        </div>

        <div class="grid grid-cols-2 gap-3 items-start">
          <CInputToggleCard
            v-model="hideGeoSearch"
            :label="$t('block.geometry.hideGeoSearch')"
            :description="$t('block.geometry.hideGeoSearchDescription')"
          />
          <CInputToggleCard
            :model-value="lockBounds"
            :label="$t('block.geometry.lockBounds')"
            :description="$t('block.geometry.lockBoundsDescription')"
            @update:model-value="onLockBoundsToggle"
          />
        </div>
      </div>
    </Fieldset>

    <Divider />

    <Fieldset :legend="$t('block.geometry.feeds')">
      <div class="flex flex-col gap-3">
        <div class="grid grid-cols-12 gap-3">
          <CFormGroup :label="$t('block.geometry.displayOption.label')" class="col-span-6">
            <Select
              v-model="displayOption"
              :options="displayOptions"
              option-label="label"
              option-value="value"
              class="w-full"
            />
          </CFormGroup>
        </div>

        <Button
          :label="$t('block.geometry.addFeed')"
          icon="pi pi-plus"
          size="small"
          severity="secondary"
          class="self-start"
          @click="addFeed"
        />

        <div
          v-for="(feed, i) in feeds"
          :key="i"
          class="flex flex-col gap-3 p-3 border border-surface rounded-border"
        >
          <div class="flex items-center justify-between">
            <span class="font-semibold text-sm">{{ $t('block.geometry.feed') }} {{ i + 1 }}</span>
            <Button icon="pi pi-trash" severity="danger" text size="small" @click="removeFeed(i)" />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <CFormGroup :label="$t('block.geometry.feedModule')">
              <Select
                :model-value="feed.options?.moduleID || ''"
                :options="modules"
                option-label="name"
                option-value="moduleID"
                :placeholder="$t('block.geometry.feedModulePlaceholder')"
                class="w-full"
                filter
                show-clear
                @update:model-value="updateFeedOption(i, 'moduleID', $event || '')"
              />
            </CFormGroup>

            <CFormGroup :label="$t('block.geometry.feedGeometryField')">
              <Select
                :model-value="feed.geometryField || ''"
                :options="getGeometryFields(feed.options?.moduleID)"
                option-label="label"
                option-value="name"
                class="w-full"
                :disabled="!feed.options?.moduleID"
                show-clear
                @update:model-value="updateFeed(i, 'geometryField', $event || '')"
              />
            </CFormGroup>

            <CFormGroup :label="$t('block.geometry.feedTitleField')">
              <Select
                :model-value="feed.titleField || ''"
                :options="getStringFields(feed.options?.moduleID)"
                option-label="label"
                option-value="name"
                class="w-full"
                :disabled="!feed.options?.moduleID"
                show-clear
                @update:model-value="updateFeed(i, 'titleField', $event || '')"
              />
            </CFormGroup>

            <CFormGroup :label="$t('block.geometry.feedColor')">
              <CInputColorPicker
                :model-value="feed.options?.color || getPrimaryColor()"
                show-text
                @update:model-value="updateFeedOption(i, 'color', $event)"
              />
            </CFormGroup>
          </div>

          <CFormGroup :label="$t('block.geometry.feedPrefilter')">
            <CInputExpression
              :ref="el => (prefilterInputs[i] = el)"
              :model-value="feed.options?.prefilter || ''"
              dialect="ql"
              :scope="scope"
              :query-fields="getFields(feed.options?.moduleID)"
              :placeholder="$t('block.geometry.feedPrefilterPlaceholder')"
              @update:model-value="updateFeedOption(i, 'prefilter', $event || '')"
            />
            <CExpressionHint :scope="scope" @insert="prefilterInputs[i]?.insert($event)" />
          </CFormGroup>

          <div class="flex gap-4">
            <div class="flex items-center gap-2">
              <Checkbox
                :model-value="feed.displayMarker !== false"
                binary
                :input-id="`marker-${i}`"
                @update:model-value="updateFeed(i, 'displayMarker', $event)"
              />
              <label :for="`marker-${i}`" class="text-sm">
                {{ $t('block.geometry.displayMarker') }}
              </label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox
                :model-value="!!feed.displayPolygon"
                binary
                :input-id="`polygon-${i}`"
                @update:model-value="updateFeed(i, 'displayPolygon', $event)"
              />
              <label :for="`polygon-${i}`" class="text-sm">
                {{ $t('block.geometry.displayPolygon') }}
              </label>
            </div>
          </div>
        </div>
      </div>
    </Fieldset>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { components, mapGeo, useModuleStore } from '@planetcrust/human-vue'
import { useI18n } from 'vue-i18n'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const { CInputColorPicker, CInputToggleCard, CMap } = components

const { t } = useI18n()
const moduleStore = useModuleStore()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const prefilterInputs = ref([])
const { scope } = useExpressionScope({ page: computed(() => props.page) })

const block = inject('blockDraft')

const modules = computed(() => moduleStore.set || [])

const displayOptions = [
  { value: 'sameTab', label: t('block.geometry.displayOption.sameTab') },
  { value: 'newTab', label: t('block.geometry.displayOption.newTab') },
  { value: 'modal', label: t('block.geometry.displayOption.modal') },
]

// Default feed colour follows the live theme primary so it doesn't drift
// from a customized brand colour; the brand hex is kept only as the
// defensive fallback for when the CSS var isn't available.
function getPrimaryColor() {
  return (
    getComputedStyle(document.documentElement).getPropertyValue('--p-primary-color').trim() ||
    '#09344E'
  )
}

function getFields(moduleID) {
  if (!moduleID) return []
  const mod = moduleStore.getByID(moduleID)
  return mod?.fields || []
}

function getGeometryFields(moduleID) {
  return getFields(moduleID)
    .filter(f => f.kind === 'Geometry')
    .map(f => ({ name: f.name, label: f.label || f.name }))
}

function getStringFields(moduleID) {
  return getFields(moduleID)
    .filter(f => ['String', 'Email', 'Url'].includes(f.kind))
    .map(f => ({ name: f.name, label: f.label || f.name }))
}

function updateOptions(patch) {
  if (!block.value.options) block.value.options = {}
  Object.assign(block.value.options, patch)
}

const hideGeoSearch = computed({
  get: () => block.value.options?.hideGeoSearch ?? true,
  set: v => updateOptions({ hideGeoSearch: v }),
})

const displayOption = computed({
  get: () => block.value.options?.displayOption || 'sameTab',
  set: v => updateOptions({ displayOption: v }),
})

// center: [lat, lng]
const center = computed(() => {
  const c = block.value.options?.center
  return Array.isArray(c) && c.length === 2 ? c : [0, 0]
})

const zoomStarting = computed(() => block.value.options?.zoomStarting ?? 2)

function onMapCenter(center) {
  const rounded = mapGeo.roundLatLng(center)
  if (rounded) updateOptions({ center: rounded })
}

function onMapZoom(z) {
  updateOptions({ zoomStarting: z })
}
const zoomMin = computed({
  get: () => block.value.options?.zoomMin ?? 1,
  set: v => updateOptions({ zoomMin: v }),
})
const zoomMax = computed({
  get: () => block.value.options?.zoomMax ?? 18,
  set: v => updateOptions({ zoomMax: v }),
})

const lockBounds = computed(() => !!block.value.options?.lockBounds)

// options.bounds is [[swLat, swLng], [neLat, neLng]], and holds the locked area
// only — unlocking clears it.
const lockedBounds = computed(() =>
  lockBounds.value ? mapGeo.parseBounds(block.value.options?.bounds) : null,
)

// Outlines the locked area, which stays visible when zooming out past it.
const boundsOutline = computed(() => {
  const ring = mapGeo.boundsRing(lockedBounds.value)
  if (!ring) return []
  return [
    {
      latLngs: ring,
      color: getPrimaryColor(),
      weight: 2,
      dashArray: '6 4',
      fillOpacity: 0.05,
      interactive: false,
    },
  ]
})

// The map's live viewport, reported by CMap on ready and on every move.
const currentMapBounds = ref(null)

function onMapBounds(b) {
  currentMapBounds.value = b
}

// Locking captures the area on screen right now and freezes the preview on it;
// unlocking gives the map back and drops the area, since nothing is bounded any
// more.
function onLockBoundsToggle(v) {
  if (!v) {
    updateOptions({ lockBounds: false, bounds: null })
    return
  }

  const bounds = mapGeo.roundBounds(currentMapBounds.value)
  updateOptions({ lockBounds: true, ...(bounds ? { bounds } : {}) })
}

// --- Feeds ---
const feeds = computed(() => block.value.options?.feeds || [])

function writeFeeds(next) {
  updateOptions({ feeds: next })
}

function addFeed() {
  writeFeeds([
    ...feeds.value,
    {
      resource: 'compose:record',
      titleField: '',
      geometryField: '',
      displayMarker: true,
      displayPolygon: false,
      options: {
        moduleID: '',
        color: getPrimaryColor(),
        prefilter: '',
      },
    },
  ])
}

function removeFeed(i) {
  const next = [...feeds.value]
  next.splice(i, 1)
  writeFeeds(next)
}

function updateFeed(i, key, value) {
  const next = [...feeds.value]
  next[i] = { ...next[i], [key]: value }
  writeFeeds(next)
}

function updateFeedOption(i, key, value) {
  const next = [...feeds.value]
  const feed = next[i]
  next[i] = {
    ...feed,
    options: { ...(feed.options || {}), [key]: value },
  }
  writeFeeds(next)
}
</script>
