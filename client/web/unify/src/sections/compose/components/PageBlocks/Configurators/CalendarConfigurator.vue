<template>
  <div class="flex flex-col gap-4">
    <Fieldset :legend="$t('block.calendar.viewLabel')">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <CFormGroup :label="$t('block.calendar.view.enabled')">
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
        </CFormGroup>

        <CFormGroup
          :label="$t('block.calendar.view.default')"
          :description="$t('block.calendar.view.footnote')"
        >
          <Select
            :model-value="localOptions.defaultView"
            :options="viewOptions"
            option-label="label"
            option-value="value"
            class="w-full"
            @update:model-value="updateOption('defaultView', $event)"
          />
        </CFormGroup>

        <CFormGroup :label="$t('block.calendar.view.onEventClick')">
          <Select
            :model-value="localOptions.eventDisplayOption || 'sameTab'"
            :options="eventDisplayOptions"
            option-label="label"
            option-value="value"
            class="w-full"
            @update:model-value="updateOption('eventDisplayOption', $event)"
          />
        </CFormGroup>

        <CFormGroup :label="$t('block.calendar.calendarHeader')">
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
        </CFormGroup>
      </div>
    </Fieldset>

    <Divider />

    <Fieldset :legend="$t('block.calendar.feedLabel')">
      <div class="flex flex-col gap-3">
        <div
          v-for="(feed, i) in feeds"
          :key="i"
          class="flex flex-col gap-2 p-3 border border-surface rounded-border"
        >
          <div class="flex items-center justify-between">
            <span class="font-semibold text-sm">
              {{ $t('block.calendar.source.label') }} {{ i + 1 }}
            </span>
            <Button icon="pi pi-trash" severity="danger" text size="small" @click="removeFeed(i)" />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <CFormGroup :label="$t('block.calendar.eventSource')" class="md:col-span-2">
              <Select
                :model-value="feed.resource || feedResources.record"
                :options="feedSourceOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                @update:model-value="feed.resource = $event"
              />
            </CFormGroup>

            <!-- Module and field mapping describe a record feed; a reminder
                 feed carries its own dates and title. -->
            <template v-if="isRecordFeed(feed)">
              <CFormGroup :label="$t('block.calendar.recordFeed.moduleLabel')">
                <CInputModule
                  :model-value="feed.options?.moduleID"
                  :namespaceID="namespace?.namespaceID"
                  :placeholder="$t('block.calendar.recordFeed.modulePlaceholder')"
                  @update:model-value="onModuleChange(feed, $event)"
                />
              </CFormGroup>

              <CFormGroup :label="$t('block.calendar.recordFeed.titleLabel')">
                <CInputModuleField
                  v-model="feed.titleField"
                  :module-i-d="feed.options?.moduleID"
                  :kinds="['String', 'Email', 'Url']"
                  :placeholder="$t('block.calendar.recordFeed.titlePlaceholder')"
                  :disabled="!feed.options?.moduleID"
                />
              </CFormGroup>

              <CFormGroup :label="$t('block.calendar.recordFeed.eventStartFieldLabel')">
                <CInputModuleField
                  v-model="feed.startField"
                  :module-i-d="feed.options?.moduleID"
                  :kinds="['DateTime']"
                  exclude-multi
                  include-system
                  :placeholder="$t('block.calendar.recordFeed.eventStartFieldPlaceholder')"
                  :disabled="!feed.options?.moduleID"
                />
              </CFormGroup>

              <CFormGroup :label="$t('block.calendar.recordFeed.eventEndFieldLabel')">
                <CInputModuleField
                  v-model="feed.endField"
                  :module-i-d="feed.options?.moduleID"
                  :kinds="['DateTime']"
                  exclude-multi
                  include-system
                  :placeholder="$t('block.calendar.recordFeed.eventEndFieldPlaceholder')"
                  :disabled="!feed.options?.moduleID || feed.allDay"
                />
                <div class="flex items-center gap-2 mt-1">
                  <Checkbox v-model="feed.allDay" :input-id="`all-day-${i}`" binary />
                  <label :for="`all-day-${i}`" class="text-sm">
                    {{ $t('block.calendar.recordFeed.eventAllDay') }}
                  </label>
                </div>
              </CFormGroup>

              <CFormGroup
                :label="$t('block.calendar.recordFeed.prefilterLabel')"
                class="md:col-span-2"
              >
                <CInputExpression
                  :ref="el => (prefilterInputs[i] = el)"
                  v-model="feed.options.prefilter"
                  dialect="ql"
                  :scope="scope"
                  :query-fields="moduleFields(feed.options?.moduleID)"
                  :placeholder="$t('block.calendar.recordFeed.prefilterPlaceholder')"
                />
                <CExpressionHint :scope="scope" @insert="prefilterInputs[i]?.insert($event)" />
              </CFormGroup>
            </template>

            <CFormGroup :label="$t('block.calendar.colorLabel')">
              <CInputColorPicker
                :model-value="feed.options?.color || getPrimaryColor()"
                show-text
                @update:model-value="feed.options.color = $event"
              />
            </CFormGroup>
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
    </Fieldset>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { compose } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'
import { useModuleStore } from '@planetcrust/human-vue'
import { useI18n } from 'vue-i18n'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const { CInputColorPicker } = components

const { t } = useI18n()
const moduleStore = useModuleStore()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const prefilterInputs = ref([])
const { scope } = useExpressionScope({ page: computed(() => props.page) })

// Fields of the module this row queries — bare identifiers in its filter.
function moduleFields(moduleID) {
  return (moduleID && moduleStore.getByID(moduleID)?.fields) || []
}

const block = inject('blockDraft')

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

// Default feed colour follows the live theme primary so it doesn't drift
// from a customized brand colour; the brand hex is kept only as the
// defensive fallback for when the CSS var isn't available.
function getPrimaryColor() {
  return (
    getComputedStyle(document.documentElement).getPropertyValue('--p-primary-color').trim() ||
    '#09344E'
  )
}

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

// A feed names its source with the resource it reads, which is what the block
// runtime dispatches on; the option list follows the model so a resource added
// there needs no change here.
const feedResources = compose.PageBlockCalendar.feedResources

const feedSourceOptions = computed(() =>
  Object.entries(feedResources).map(([key, value]) => ({
    value,
    label: t(`block.calendar.${key}Feed.optionLabel`),
  })),
)

function isRecordFeed(feed) {
  return (feed.resource || feedResources.record) === feedResources.record
}
</script>
