<template>
  <div>
    <template v-if="hasCoordinates">
      <span
        v-for="(coord, index) in coordList"
        :key="index"
        :class="{ 'block': field.options?.multiDelimiter === '\n' }"
      >
        <a
          v-if="!disableClick"
          class="text-primary cursor-pointer text-nowrap"
          @click.stop="openMap(index)"
        >
          {{ coord[0] }}, {{ coord[1] }}{{ index < coordList.length - 1 ? (field.options?.multiDelimiter || ', ') : '' }}
        </a>
        <span v-else class="text-nowrap">
          {{ coord[0] }}, {{ coord[1] }}{{ index < coordList.length - 1 ? (field.options?.multiDelimiter || ', ') : '' }}
        </span>
      </span>
    </template>

    <!-- Read-only map dialog -->
    <Dialog
      v-model:visible="showMapDialog"
      :header="field.label || field.name"
      modal
      :style="{ width: '80vw', height: '80vh' }"
      :pt="{
        root: 'flex flex-col',
        content: 'flex flex-col flex-1 min-h-0 p-0',
      }"
      @after-show="mapComponentRef?.invalidateSize()"
    >
      <CMap
        ref="mapComponentRef"
        :center="mapCenter"
        :zoom="13"
        :markers="mapMarkers"
        :hide-geo-search="field.options?.hideGeoSearch"
        :hide-current-location-button="field.options?.hideCurrentLocationButton"
        class="flex-1 min-h-0 w-full"
      />
    </Dialog>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import CMap from '../../map/CMap.vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  record: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  valueOnly: {
    type: Boolean,
    default: false,
  },
  extraOptions: {
    type: Object,
    default: () => ({}),
  },
  disableClick: {
    type: Boolean,
    default: false,
  },
})

const showMapDialog = ref(false)
const selectedIndex = ref(0)
const mapComponentRef = ref(null)

const rawValue = computed(() => {
  if (props.field.isSystem) {
    return props.record[props.field.name]
  }
  return props.record?.values?.[props.field.name]
})

// Parse coordinates from stored JSON values
const coordList = computed(() => {
  const v = rawValue.value
  if (!v) return []

  const parseCoords = (str) => {
    try {
      const parsed = JSON.parse(str || '{}')
      if (parsed?.coordinates?.length === 2) return parsed.coordinates
    } catch {
      // ignore
    }
    return null
  }

  if (props.field.isMulti && Array.isArray(v)) {
    return v.map(parseCoords).filter(Boolean)
  }

  const coords = parseCoords(v)
  return coords ? [coords] : []
})

const hasCoordinates = computed(() => coordList.value.length > 0)

const mapCenter = computed(() => {
  const coords = coordList.value[selectedIndex.value]
  return coords || props.field.options?.center || [30, 30]
})

const mapMarkers = computed(() =>
  coordList.value.map((coords, i) => ({
    value: coords,
    opacity: selectedIndex.value === i ? 1.0 : 0.6,
  })),
)

function openMap(index) {
  selectedIndex.value = index
  showMapDialog.value = true
}
</script>
