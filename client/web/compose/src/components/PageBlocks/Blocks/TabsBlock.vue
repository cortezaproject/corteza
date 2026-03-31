<template>
  <PageBlock :block="block">
    <div
      v-if="!tabbedBlocks.length"
      class="flex items-center justify-center h-full p-3 text-muted-color italic"
    >
      {{ $t('block.tabs.noTabs') }}
    </div>

    <Tabs
      v-else
      :value="activeTab"
      :orientation="tabStyle.orientation"
      :show-navigators="tabStyle.orientation !== 'vertical'"
      class="h-full tabs-block flex"
      :class="tabsClasses"
      @update:value="onTabChange"
    >
      <TabList :class="tabListClasses" :pt="tabListPt">
        <Tab
          v-for="(tab, index) in tabbedBlocks"
          :key="`tab-header-${index}`"
          :value="index"
          :class="tabClasses"
        >
          {{ tab.title || `${$t('block.tabs.tab')} ${index + 1}` }}
        </Tab>
      </TabList>

      <TabPanels :class="tabPanelsClasses">
        <TabPanel
          v-for="(tab, index) in tabbedBlocks"
          :key="`tab-panel-${index}`"
          :value="index"
          class="h-full p-0"
        >
          <div v-if="tab.block && shouldRenderTab(tab, index)" class="h-full">
            <component
              :is="resolveBlock(tab.block.kind)"
              v-if="resolveBlock(tab.block.kind)"
              :block="tab.block"
              :blocks="blocks"
              :namespace="namespace"
              :page="page"
              :record="record"
            />
            <div v-else class="flex items-center justify-center h-full text-muted-color italic p-2">
              {{ tab.block.kind || $t('block.tabs.noBlock') }}
            </div>
          </div>
          <div v-else class="flex items-center justify-center h-full p-3 text-muted-color italic">
            {{ $t('block.tabs.noBlock') }}
          </div>
        </TabPanel>
      </TabPanels>
    </Tabs>
  </PageBlock>
</template>

<script setup>
import { computed, ref } from 'vue'
import PageBlock from './PageBlock.vue'
import { resolveBlock } from '../registry'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
  // All blocks on the page (needed to resolve tabbed block references)
  blocks: { type: Array, default: () => [] },
})

const activeTab = ref(0)
const visitedTabs = ref({ 0: true })

const tabStyle = computed(() => ({
  appearance: 'tabs',
  alignment: 'center',
  justify: 'justify',
  orientation: 'horizontal',
  position: 'start',
  ...(props.block.options?.style || {}),
}))

const tabsClasses = computed(() => [
  {
    'tabs-block--vertical': tabStyle.value.orientation === 'vertical',
    'tabs-block--vertical-end':
      tabStyle.value.orientation === 'vertical' && tabStyle.value.position === 'end',
  },
  tabStyle.value.orientation === 'vertical'
    ? tabStyle.value.position === 'end'
      ? 'flex-row-reverse'
      : 'flex-row'
    : tabStyle.value.position === 'end'
      ? 'flex-col-reverse'
      : 'flex-col',
])

const tabListClasses = computed(() => [
  'flex flex-nowrap',
  tabStyle.value.orientation === 'vertical'
    ? [
        'w-56 min-w-56',
        tabStyle.value.justify === 'justify'
          ? 'h-full overflow-hidden'
          : 'overflow-x-hidden overflow-y-auto',
      ]
    : 'overflow-x-auto',
  {
    'tabs-block__list--justify': tabStyle.value.justify === 'justify',
    'tabs-block__list--left':
      tabStyle.value.justify !== 'justify' && tabStyle.value.alignment === 'left',
    'tabs-block__list--center':
      tabStyle.value.justify !== 'justify' && tabStyle.value.alignment === 'center',
    'tabs-block__list--right':
      tabStyle.value.justify !== 'justify' && tabStyle.value.alignment === 'right',
  },
])

const tabListPt = computed(() => ({
  root: {
    class:
      tabStyle.value.orientation === 'vertical'
        ? tabStyle.value.justify === 'justify'
          ? 'block h-full'
          : 'block'
        : undefined,
  },
  content: {
    class:
      tabStyle.value.orientation === 'vertical'
        ? tabStyle.value.justify === 'justify'
          ? 'h-full overflow-hidden'
          : 'overflow-x-hidden overflow-y-auto'
        : 'overflow-x-auto overflow-y-hidden',
  },
  tabList: {
    style: {
      borderWidth: tabStyle.value.orientation === 'vertical'
        ? tabStyle.value.position === 'end'
          ? '0 0 0 1px'
          : '0 1px 0 0'
        : tabStyle.value.position === 'end'
          ? '1px 0 0 0'
          : '0 0 1px 0',
    },
    class: [
      'flex w-full flex-nowrap border-solid border-surface',
      tabStyle.value.orientation === 'vertical'
        ? ['!flex-col !items-stretch', tabStyle.value.justify === 'justify' ? 'h-full' : undefined]
        : 'flex-row',
      tabStyle.value.appearance === 'pills' ? 'gap-2 p-2' : undefined,
      {
        'justify-start':
          tabStyle.value.justify !== 'justify' && tabStyle.value.alignment === 'left',
        'justify-center':
          tabStyle.value.justify !== 'justify' && tabStyle.value.alignment === 'center',
        'justify-end': tabStyle.value.justify !== 'justify' && tabStyle.value.alignment === 'right',
      },
    ],
  },
  activeBar: {
    class:
      tabStyle.value.appearance === 'pills' || tabStyle.value.orientation === 'vertical'
        ? 'hidden'
        : undefined,
  },
}))

const tabClasses = computed(() => [
  'border-0 border-surface whitespace-nowrap',
  tabStyle.value.appearance === 'pills'
    ? ['border-0', 'border-b-0', 'border-t-0', 'border-x-0']
    : tabStyle.value.orientation === 'vertical'
      ? ['border-b', 'border-x-0']
      : tabStyle.value.position === 'end'
        ? ['border-t-0', 'border-b-0', 'border-x-0']
        : ['border-b', 'border-t-0', 'border-x-0'],
  {
    'flex-1': tabStyle.value.justify === 'justify',
    'w-full justify-start': tabStyle.value.orientation === 'vertical',
    'text-sm px-3 py-2': tabStyle.value.appearance === 'small',
    'rounded-md border-0 border-b-0 border-t-0': tabStyle.value.appearance === 'pills',
    'tabs-block__tab--pills': tabStyle.value.appearance === 'pills',
  },
])

const tabPanelsClasses = computed(() => 'h-full flex-1 overflow-hidden !p-0')

const tabbedBlocks = computed(() => {
  const tabs = props.block.options?.tabs || []
  return tabs.reduce((acc, { blockID, title, lazy = true }) => {
    let block = null

    if (blockID) {
      block = props.blocks.find(b => {
        const bid = b.blockID && b.blockID !== '0' ? b.blockID : b.meta?.tempID
        return bid === blockID
      })
    }

    if (!block && !title) return acc

    acc.push({ block: block || null, title: title || '', lazy })
    return acc
  }, [])
})

function onTabChange(index) {
  activeTab.value = index
  visitedTabs.value[index] = true
}

function shouldRenderTab(tab, index) {
  if (!tab.lazy) return true
  return !!visitedTabs.value[index]
}
</script>

<style scoped>
.tabs-block :deep(.p-tabpanel) {
  height: 100%;
}

.tabs-block--vertical :deep(.p-tablist-tab-list) {
  flex-direction: column;
  align-items: stretch;
}

.tabs-block--vertical :deep(.p-tab) {
  width: 100%;
  justify-content: flex-start;
}

.tabs-block :deep(.p-tab + .p-tab) {
  border-left-width: 0;
}

.tabs-block--vertical :deep(.p-tab.p-tab-active:not(.tabs-block__tab--pills)) {
  border-bottom-style: solid !important;
  border-bottom-width: 1px !important;
  border-bottom-color: var(--p-primary-color) !important;
}

.tabs-block :deep(.p-tab-active.tabs-block__tab--pills) {
  background: var(--p-primary-color);
  color: var(--p-primary-contrast-color);
}
</style>
