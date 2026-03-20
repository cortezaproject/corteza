<template>
  <PageBlock :block="block">
    <div v-if="!tabbedBlocks.length" class="flex items-center justify-center h-full p-3 text-muted-color italic">
      {{ $t('block.tabs.noTabs') }}
    </div>

    <Tabs
      v-else
      :value="activeTab"
      class="h-full tabs-block"
      @update:value="onTabChange"
    >
      <TabList>
        <Tab
          v-for="(tab, index) in tabbedBlocks"
          :key="`tab-header-${index}`"
          :value="index"
        >
          {{ tab.title || `${$t('block.tabs.tab')} ${index + 1}` }}
        </Tab>
      </TabList>

      <TabPanels class="h-full">
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

<style>
.tabs-block .p-tabpanels {
  padding: 0 !important;
  height: 100%;
}
.tabs-block .p-tabpanel {
  height: 100%;
}
</style>
