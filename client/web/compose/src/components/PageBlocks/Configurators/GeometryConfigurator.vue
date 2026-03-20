<template>
  <div class="flex flex-col gap-5">
    <Message severity="info" variant="simple">
      {{ $t('block.geometry.feedLabel') }}
    </Message>

    <!-- Map Settings -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.geometry.mapSettings') }}</h5>

      <div class="flex items-center gap-2">
        <Checkbox v-model="hideGeoSearch" binary input-id="hideGeoSearch" />
        <label for="hideGeoSearch" class="text-sm">{{ $t('block.geometry.hideGeoSearch') }}</label>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.geometry.displayOption.label') }}</label>
        <Select
          v-model="displayOption"
          :options="displayOptions"
          option-label="label"
          option-value="value"
          class="w-full"
        />
      </div>

      <!-- Center Coordinates -->
      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.geometry.center') }}</label>
        <div class="grid grid-cols-2 gap-2">
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.centerLat') }}</label>
            <InputNumber v-model="centerLat" :min="-90" :max="90" :max-fraction-digits="6" class="w-full" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.centerLng') }}</label>
            <InputNumber v-model="centerLng" :min="-180" :max="180" :max-fraction-digits="6" class="w-full" />
          </div>
        </div>
      </div>

      <!-- Zoom -->
      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.geometry.zoom') }}</label>
        <div class="grid grid-cols-3 gap-2">
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.zoomStarting') }}</label>
            <InputNumber v-model="zoomStarting" :min="1" :max="20" class="w-full" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.zoomMin') }}</label>
            <InputNumber v-model="zoomMin" :min="1" :max="20" class="w-full" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.zoomMax') }}</label>
            <InputNumber v-model="zoomMax" :min="1" :max="20" class="w-full" />
          </div>
        </div>
      </div>

      <!-- Lock Bounds -->
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <Checkbox v-model="lockBounds" binary input-id="lockBounds" />
          <label for="lockBounds" class="text-sm">{{ $t('block.geometry.lockBounds') }}</label>
        </div>

        <template v-if="lockBounds">
          <div class="grid grid-cols-2 gap-2">
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.geometry.boundsMinLat') }}</label>
              <InputNumber v-model="boundsMinLat" :min="-90" :max="90" :max-fraction-digits="6" class="w-full" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.geometry.boundsMaxLat') }}</label>
              <InputNumber v-model="boundsMaxLat" :min="-90" :max="90" :max-fraction-digits="6" class="w-full" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.geometry.boundsMinLng') }}</label>
              <InputNumber v-model="boundsMinLng" :min="-180" :max="180" :max-fraction-digits="6" class="w-full" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.geometry.boundsMaxLng') }}</label>
              <InputNumber v-model="boundsMaxLng" :min="-180" :max="180" :max-fraction-digits="6" class="w-full" />
            </div>
          </div>
        </template>
      </div>
    </div>

    <Divider />

    <!-- Feeds Section -->
    <div class="flex flex-col gap-3">
      <div class="flex items-center justify-between">
        <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.geometry.feeds') }}</h5>
        <Button
          :label="$t('general.label.add')"
          icon="pi pi-plus"
          size="small"
          severity="secondary"
          @click="addFeed"
        />
      </div>

      <div v-for="(feed, i) in feeds" :key="i" class="flex flex-col gap-3 p-3 border border-surface rounded-border">
        <div class="flex items-center justify-between">
          <span class="font-semibold text-sm">{{ $t('block.geometry.feed') }} {{ i + 1 }}</span>
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            size="small"
            @click="removeFeed(i)"
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <!-- Module -->
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.feedModule') }}</label>
            <Select
              :model-value="feed.moduleID || ''"
              :options="modules"
              option-label="name"
              option-value="moduleID"
              :placeholder="$t('block.geometry.feedModulePlaceholder')"
              class="w-full"
              filter
              show-clear
              @update:model-value="updateFeed(i, 'moduleID', $event)"
            />
          </div>

          <!-- Geometry field -->
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.feedGeometryField') }}</label>
            <Select
              :model-value="feed.geometryField || ''"
              :options="getAllFields(feed.moduleID)"
              option-label="label"
              option-value="name"
              class="w-full"
              :disabled="!feed.moduleID"
              show-clear
              @update:model-value="updateFeed(i, 'geometryField', $event)"
            />
          </div>

          <!-- Title field -->
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.feedTitleField') }}</label>
            <Select
              :model-value="feed.titleField || ''"
              :options="getStringFields(feed.moduleID)"
              option-label="label"
              option-value="name"
              class="w-full"
              :disabled="!feed.moduleID"
              show-clear
              @update:model-value="updateFeed(i, 'titleField', $event)"
            />
          </div>

          <!-- Color -->
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.geometry.feedColor') }}</label>
            <CInputColorPicker
              :model-value="feed.color || '#09344E'"
              show-text
              @update:model-value="updateFeed(i, 'color', $event)"
            />
          </div>
        </div>

        <!-- Prefilter -->
        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('block.geometry.feedPrefilter') }}</label>
          <Textarea
            :model-value="feed.prefilter || ''"
            :placeholder="$t('block.geometry.feedPrefilterPlaceholder')"
            rows="2"
            class="w-full"
            @update:model-value="updateFeed(i, 'prefilter', $event)"
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
            <label :for="`marker-${i}`" class="text-sm">{{ $t('block.geometry.displayMarker') }}</label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox
              :model-value="!!feed.displayPolygon"
              binary
              :input-id="`polygon-${i}`"
              @update:model-value="updateFeed(i, 'displayPolygon', $event)"
            />
            <label :for="`polygon-${i}`" class="text-sm">{{ $t('block.geometry.displayPolygon') }}</label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { components } from '@cortezaproject/corteza-vue-next'
import { useModuleStore } from '@/stores/module'
import { useI18n } from 'vue-i18n'

const { CInputColorPicker } = components

const { t } = useI18n()
const moduleStore = useModuleStore()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const modules = computed(() => moduleStore.set || [])

const displayOptions = [
  { value: 'sameTab', label: t('block.geometry.displayOption.sameTab') },
  { value: 'newTab', label: t('block.geometry.displayOption.newTab') },
  { value: 'modal', label: t('block.geometry.displayOption.modal') },
]

function getAllFields(moduleID) {
  if (!moduleID) return []
  const mod = moduleStore.getByID(moduleID)
  if (!mod) return []
  return (mod.fields || []).map(f => ({ name: f.name, label: f.label || f.name }))
}

function getStringFields(moduleID) {
  if (!moduleID) return []
  const mod = moduleStore.getByID(moduleID)
  if (!mod) return []
  return (mod.fields || [])
    .filter(f => ['String', 'Email', 'Url'].includes(f.kind))
    .map(f => ({ name: f.name, label: f.label || f.name }))
}

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

function updateMapOption(key, value) {
  const map = { ...(props.block.options?.map || {}), [key]: value }
  updateOptions('map', map)
}

const hideGeoSearch = computed({
  get: () => props.block.options?.hideGeoSearch || false,
  set: v => updateOptions('hideGeoSearch', v),
})

const displayOption = computed({
  get: () => props.block.options?.displayOption || 'sameTab',
  set: v => updateOptions('displayOption', v),
})

const centerLat = computed({
  get: () => props.block.options?.map?.center?.[0] ?? 0,
  set: v => {
    const center = [...(props.block.options?.map?.center || [0, 0])]
    center[0] = v
    updateMapOption('center', center)
  },
})

const centerLng = computed({
  get: () => props.block.options?.map?.center?.[1] ?? 0,
  set: v => {
    const center = [...(props.block.options?.map?.center || [0, 0])]
    center[1] = v
    updateMapOption('center', center)
  },
})

const zoomStarting = computed({
  get: () => props.block.options?.map?.zoomStarting ?? 10,
  set: v => updateMapOption('zoomStarting', v),
})

const zoomMin = computed({
  get: () => props.block.options?.map?.zoomMin ?? 1,
  set: v => updateMapOption('zoomMin', v),
})

const zoomMax = computed({
  get: () => props.block.options?.map?.zoomMax ?? 18,
  set: v => updateMapOption('zoomMax', v),
})

const lockBounds = computed({
  get: () => props.block.options?.map?.lockBounds || false,
  set: v => updateMapOption('lockBounds', v),
})

const boundsMinLat = computed({
  get: () => props.block.options?.map?.bounds?.[0] ?? -90,
  set: v => {
    const bounds = [...(props.block.options?.map?.bounds || [-90, 90, -180, 180])]
    bounds[0] = v
    updateMapOption('bounds', bounds)
  },
})

const boundsMaxLat = computed({
  get: () => props.block.options?.map?.bounds?.[1] ?? 90,
  set: v => {
    const bounds = [...(props.block.options?.map?.bounds || [-90, 90, -180, 180])]
    bounds[1] = v
    updateMapOption('bounds', bounds)
  },
})

const boundsMinLng = computed({
  get: () => props.block.options?.map?.bounds?.[2] ?? -180,
  set: v => {
    const bounds = [...(props.block.options?.map?.bounds || [-90, 90, -180, 180])]
    bounds[2] = v
    updateMapOption('bounds', bounds)
  },
})

const boundsMaxLng = computed({
  get: () => props.block.options?.map?.bounds?.[3] ?? 180,
  set: v => {
    const bounds = [...(props.block.options?.map?.bounds || [-90, 90, -180, 180])]
    bounds[3] = v
    updateMapOption('bounds', bounds)
  },
})

// --- Feeds ---
const feeds = computed(() => props.block.options?.feeds || [])

function addFeed() {
  updateOptions('feeds', [...feeds.value, {
    moduleID: '',
    geometryField: '',
    titleField: '',
    color: '#09344E',
    prefilter: '',
    displayMarker: true,
    displayPolygon: false,
  }])
}

function removeFeed(i) {
  const updated = [...feeds.value]
  updated.splice(i, 1)
  updateOptions('feeds', updated)
}

function updateFeed(i, key, value) {
  const updated = [...feeds.value]
  updated[i] = { ...updated[i], [key]: value }
  updateOptions('feeds', updated)
}
</script>
