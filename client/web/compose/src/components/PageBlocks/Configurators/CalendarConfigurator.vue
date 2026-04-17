<template>
  <div class="flex flex-col gap-4">
    <!-- Display Settings -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">
        {{ $t('block.calendar.viewLabel') }}
      </h5>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- Enabled views -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.calendar.view.enabled') }}
          </label>
          <div class="flex flex-wrap gap-2">
            <div v-for="view in availableViews" :key="view" class="flex items-center gap-1">
              <Checkbox
                :model-value="isViewEnabled(view)"
                :input-id="`view-${view}`"
                binary
                :disabled="isHeaderHidden"
                @update:model-value="toggleView(view, $event)"
              />
              <label :for="`view-${view}`" class="text-sm">
                {{ $t(`block.calendar.view.${view}`) }}
              </label>
            </div>
          </div>
        </div>

        <!-- Default view -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.calendar.view.default') }}
          </label>
          <Select
            :model-value="localOptions.defaultView"
            :options="viewOptions"
            option-label="label"
            option-value="value"
            class="w-full"
            @update:model-value="updateOption('defaultView', $event)"
          />
          <small class="text-muted-color">{{ $t('block.calendar.view.footnote') }}</small>
        </div>

        <!-- Event click behavior -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.calendar.view.onEventClick') }}
          </label>
          <Select
            :model-value="localOptions.eventDisplayOption || 'sameTab'"
            :options="eventDisplayOptions"
            option-label="label"
            option-value="value"
            class="w-full"
            @update:model-value="updateOption('eventDisplayOption', $event)"
          />
        </div>

        <!-- Header toggles -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.calendar.calendarHeader') }}
          </label>
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <Checkbox v-model="headerHide" input-id="hide-header" binary />
              <label for="hide-header" class="text-sm">{{ $t('block.calendar.hideHeader') }}</label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="headerHidePrevNext"
                input-id="hide-nav"
                binary
                :disabled="isHeaderHidden"
              />
              <label for="hide-nav" class="text-sm">
                {{ $t('block.calendar.hideNavigation') }}
              </label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="headerHideToday"
                input-id="hide-today"
                binary
                :disabled="isHeaderHidden"
              />
              <label for="hide-today" class="text-sm">{{ $t('block.calendar.hideToday') }}</label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="headerHideTitle"
                input-id="hide-title"
                binary
                :disabled="isHeaderHidden"
              />
              <label for="hide-title" class="text-sm">{{ $t('block.calendar.hideTitle') }}</label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <Divider />

    <!-- Feed Sources -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">
        {{ $t('block.calendar.feedLabel') }}
      </h5>

      <div v-for="(feed, i) in feeds" :key="i" class="flex flex-col gap-2 p-3 border rounded-lg">
        <div class="flex items-center justify-between">
          <span class="font-semibold text-sm">
            {{ $t('block.calendar.source.label') }} {{ i + 1 }}
          </span>
          <Button icon="pi pi-trash" severity="danger" text size="small" @click="removeFeed(i)" />
        </div>

        <!-- Module -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <!-- Resource type -->
          <div class="flex flex-col gap-1 md:col-span-2">
            <label class="text-primary font-medium text-sm">
              {{ $t('block.calendar.recordFeed.resourceType') }}
            </label>
            <Select
              :model-value="feed.resourceType || 'record'"
              :options="resourceTypeOptions"
              option-label="label"
              option-value="value"
              class="w-full"
              @update:model-value="updateFeedResourceType(feed, $event)"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('block.calendar.recordFeed.moduleLabel') }}
            </label>
            <Select
              :model-value="feed.options?.moduleID"
              :options="modules"
              option-label="name"
              option-value="moduleID"
              :placeholder="$t('block.calendar.recordFeed.modulePlaceholder')"
              class="w-full"
              filter
              @update:model-value="onModuleChange(feed, $event)"
            />
          </div>

          <!-- Title field -->
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('block.calendar.recordFeed.titleLabel') }}
            </label>
            <Select
              v-model="feed.titleField"
              :options="getTitleFields(feed.options?.moduleID)"
              option-label="label"
              option-value="name"
              :placeholder="$t('block.calendar.recordFeed.titlePlaceholder')"
              class="w-full"
              :disabled="!feed.options?.moduleID"
            />
          </div>

          <!-- Start field -->
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('block.calendar.recordFeed.eventStartFieldLabel') }}
            </label>
            <Select
              v-model="feed.startField"
              :options="getDateFields(feed.options?.moduleID)"
              option-label="label"
              option-value="name"
              :placeholder="$t('block.calendar.recordFeed.eventStartFieldPlaceholder')"
              class="w-full"
              :disabled="!feed.options?.moduleID"
            />
          </div>

          <!-- End field -->
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('block.calendar.recordFeed.eventEndFieldLabel') }}
            </label>
            <Select
              v-model="feed.endField"
              :options="getDateFields(feed.options?.moduleID)"
              option-label="label"
              option-value="name"
              :placeholder="$t('block.calendar.recordFeed.eventEndFieldPlaceholder')"
              class="w-full"
              :disabled="!feed.options?.moduleID || feed.allDay"
            />
            <div class="flex items-center gap-2 mt-1">
              <Checkbox v-model="feed.allDay" :input-id="`all-day-${i}`" binary />
              <label :for="`all-day-${i}`" class="text-sm">
                {{ $t('block.calendar.recordFeed.eventAllDay') }}
              </label>
            </div>
          </div>

          <!-- Prefilter -->
          <div class="flex flex-col gap-1 md:col-span-2">
            <label class="text-primary font-medium text-sm">
              {{ $t('block.calendar.recordFeed.prefilterLabel') }}
            </label>
            <InputText
              v-model="feed.options.prefilter"
              :placeholder="$t('block.calendar.recordFeed.prefilterPlaceholder')"
              class="w-full"
            />
          </div>

          <!-- Event color -->
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('block.calendar.colorLabel') }}
            </label>
            <CInputColorPicker
              :model-value="feed.options?.color || '#09344E'"
              show-text
              @update:model-value="feed.options.color = $event"
            />
          </div>
        </div>
      </div>

      <Button
        :label="$t('block.calendar.addEventsSource')"
        icon="pi pi-plus"
        severity="secondary"
        size="small"
        class="self-start"
        @click="addFeed"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { compose } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'
import { useModuleStore } from '@/stores/module'
import { useI18n } from 'vue-i18n'

const { CInputColorPicker } = components

const { t } = useI18n()
const moduleStore = useModuleStore()

defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const modules = computed(() => moduleStore.set || [])

const localOptions = computed(() => block.value.options || {})

const availableViews = computed(() => compose.PageBlockCalendar.availableViews())

const viewOptions = computed(() =>
  availableViews.value.map(v => ({
    value: v,
    label: t(`block.calendar.view.${v}`),
  })),
)

const eventDisplayOptions = computed(() => [
  { value: 'sameTab', label: t('block.calendar.view.openInSameTab') },
  { value: 'newTab', label: t('block.calendar.view.openInNewTab') },
  { value: 'modal', label: t('block.calendar.view.openInModal') },
])

const feeds = computed(() => localOptions.value.feeds || [])

const isHeaderHidden = computed(() => localOptions.value.header?.hide || false)

// Header toggle v-model helpers
const headerHide = computed({
  get: () => localOptions.value.header?.hide || false,
  set: v => updateHeader('hide', v),
})
const headerHidePrevNext = computed({
  get: () => localOptions.value.header?.hidePrevNext || false,
  set: v => updateHeader('hidePrevNext', v),
})
const headerHideToday = computed({
  get: () => localOptions.value.header?.hideToday || false,
  set: v => updateHeader('hideToday', v),
})
const headerHideTitle = computed({
  get: () => localOptions.value.header?.hideTitle || false,
  set: v => updateHeader('hideTitle', v),
})

function isViewEnabled(view) {
  return (localOptions.value.header?.views || []).includes(view)
}

function toggleView(view, enabled) {
  const current = [...(localOptions.value.header?.views || [])]
  if (enabled && !current.includes(view)) {
    current.push(view)
  } else if (!enabled) {
    const idx = current.indexOf(view)
    if (idx > -1) current.splice(idx, 1)
  }
  updateHeader('views', current)
}

function updateHeader(key, value) {
  const currentHeader = { ...(localOptions.value.header || {}) }
  currentHeader[key] = value
  updateOption('header', currentHeader)
}

function updateOption(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

/**
 * Returns String/Email/Url fields for event title.
 */
function getTitleFields(moduleID) {
  if (!moduleID) return []
  const mod = moduleStore.getByID(moduleID)
  if (!mod) return []

  return mod.fields
    .filter(f => ['String', 'Email', 'Url'].includes(f.kind))
    .map(f => ({ name: f.name, label: f.label || f.name }))
    .sort((a, b) => a.label.localeCompare(b.label))
}

/**
 * Returns DateTime fields (non-multi) for start/end date.
 * Also includes system date fields (createdAt, updatedAt).
 */
function getDateFields(moduleID) {
  if (!moduleID) return []
  const mod = moduleStore.getByID(moduleID)
  if (!mod) return []

  const moduleFields = mod.fields
    .filter(f => f.kind === 'DateTime' && !f.isMulti)
    .map(f => ({ name: f.name, label: f.label || f.name }))
    .sort((a, b) => a.label.localeCompare(b.label))

  const systemFields = [
    { name: 'createdAt', label: t('block.calendar.recordFeed.systemField.createdAt') },
    { name: 'updatedAt', label: t('block.calendar.recordFeed.systemField.updatedAt') },
    { name: 'deletedAt', label: t('block.calendar.recordFeed.systemField.deletedAt') },
  ]

  return [...moduleFields, ...systemFields]
}

function addFeed() {
  const newFeed = compose.PageBlockCalendar.makeFeed()
  const currentFeeds = [...feeds.value, newFeed]
  updateOption('feeds', currentFeeds)
}

function removeFeed(index) {
  const currentFeeds = [...feeds.value]
  currentFeeds.splice(index, 1)
  updateOption('feeds', currentFeeds)
}

function onModuleChange(feed, moduleID) {
  feed.options = { ...feed.options, moduleID }
  feed.titleField = ''
  feed.startField = ''
  feed.endField = ''
}

const resourceTypeOptions = [
  { value: 'record', label: t('block.calendar.recordFeed.resourceTypeRecord') },
  { value: 'reminder', label: t('block.calendar.recordFeed.resourceTypeReminder') },
]

function updateFeedResourceType(feed, resourceType) {
  feed.resourceType = resourceType
}
</script>
