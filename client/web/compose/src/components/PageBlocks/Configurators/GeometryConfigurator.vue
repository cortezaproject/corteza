<template>
  <div class="flex flex-col gap-5">
    <div class="flex flex-col gap-3">
      <!-- Starting View: pan/zoom to set starting center and zoom -->
      <div class="flex flex-col gap-2">
        <label class="text-primary font-medium text-sm">
          {{ $t('block.geometry.startingView') }}
        </label>
        <CMap
          :center="center"
          :zoom="zoomStarting"
          :min-zoom="zoomMin"
          :max-zoom="zoomMax"
          :max-bounds="lockBounds ? lockedBounds : null"
          :hide-geo-search="hideGeoSearch"
          hide-current-location-button
          style="height: 40vh;"
          @update:center="onMapCenter"
          @update:zoom="onMapZoom"
          @update:bounds="onMapBounds"
        />
      </div>

      <!-- Zoom min/max constraints -->
      <div class="grid grid-cols-2 gap-2">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.geometry.zoomMin') }}
          </label>
          <InputNumber
            :model-value="zoomMin"
            :min="1"
            :max="20"
            show-buttons
            class="w-full"
            @input="e => (zoomMin = e.value)"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.geometry.zoomMax') }}
          </label>
          <InputNumber
            :model-value="zoomMax"
            :min="1"
            :max="20"
            show-buttons
            class="w-full"
            @input="e => (zoomMax = e.value)"
          />
        </div>
      </div>

      <!-- Toggles -->
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

    <Divider />

    <!-- Feeds Section -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.geometry.feeds') }}</h5>

      <div class="grid grid-cols-12 gap-3">
        <div class="flex flex-col gap-1 col-span-6">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.geometry.displayOption.label') }}
          </label>
          <Select
            v-model="displayOption"
            :options="displayOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>
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
          <!-- Module -->
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.feedModule') }}</label>
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
          </div>

          <!-- Geometry field -->
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">
              {{ $t('block.geometry.feedGeometryField') }}
            </label>
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
          </div>

          <!-- Title field -->
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">
              {{ $t('block.geometry.feedTitleField') }}
            </label>
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
          </div>

          <!-- Color -->
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.feedColor') }}</label>
            <CInputColorPicker
              :model-value="feed.options?.color || '#09344E'"
              show-text
              @update:model-value="updateFeedOption(i, 'color', $event)"
            />
          </div>
        </div>

        <!-- Prefilter -->
        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('block.geometry.feedPrefilter') }}</label>
          <Textarea
            :model-value="feed.options?.prefilter || ''"
            :placeholder="$t('block.geometry.feedPrefilterPlaceholder')"
            rows="2"
            class="w-full"
            @update:model-value="updateFeedOption(i, 'prefilter', $event || '')"
          />
        </div>

        <!-- Toggles -->
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
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { components } from '@planetcrust/human-vue'
import CMap from '@planetcrust/human-vue/src/components/map/CMap.vue'
import { useModuleStore } from '@/stores/module'
import { useI18n } from 'vue-i18n'

const { CInputColorPicker, CInputToggleCard } = components

const { t } = useI18n()
const moduleStore = useModuleStore()

defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const modules = computed(() => moduleStore.set || [])

const displayOptions = [
  { value: 'sameTab', label: t('block.geometry.displayOption.sameTab') },
  { value: 'newTab', label: t('block.geometry.displayOption.newTab') },
  { value: 'modal', label: t('block.geometry.displayOption.modal') },
]

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

function onMapCenter([lat, lng]) {
  updateOptions({
    center: [
      Math.round(lat * 1e6) / 1e6,
      Math.round(lng * 1e6) / 1e6,
    ],
  })
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

// bounds shape stored in options: [[swLat, swLng], [neLat, neLng]]
const lockedBounds = computed(() => {
  const b = block.value.options?.bounds
  if (Array.isArray(b) && b.length === 2 && b.every(p => Array.isArray(p) && p.length === 2)) {
    return b
  }
  return null
})

// Track the map's current viewport bounds (emitted by CMap on move/zoom)
const currentMapBounds = ref(null)

function onMapBounds(b) {
  currentMapBounds.value = b
}

function onLockBoundsToggle(v) {
  if (v) {
    const b = currentMapBounds.value || lockedBounds.value
    if (b) updateOptions({ lockBounds: true, bounds: b })
    else updateOptions({ lockBounds: true })
  } else {
    updateOptions({ lockBounds: false })
  }
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
        color: '#09344E',
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
