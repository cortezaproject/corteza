<template>
  <wrap
    v-bind="$props"
    :scrollable-body="false"
    v-on="$listeners"
  >
    <div
      v-if="!options.tabs.length"
      class="d-flex h-100 align-items-center justify-content-center"
    >
      <p class="mb-0">
        {{ $t('tabs.noTabs') }}
      </p>
    </div>

    <b-tabs
      v-else
      card
      :nav-class="navClass"
      :nav-wrapper-class="navWrapperClass"
      :content-class="contentClass"
      v-bind="{
        align: block.options.style.alignment,
        fill: block.options.style.fillJustify === 'fill',
        justified: block.options.style.fillJustify === 'justify',
        pills: block.options.style.appearance === 'pills',
        tabs: block.options.style.appearance === 'tabs',
        small: block.options.style.appearance === 'small',
        vertical: block.options.style.orientation === 'vertical',
        end: block.options.style.position === 'end'
      }"
      class="h-100"
      :class="{ 'd-flex flex-column': block.options.style.orientation !== 'vertical' }"
    >
      <b-tab
        v-for="(tab, index) in tabbedBlocks"
        :key="`${getTabTitle(tab, index)}-${index}`"
        :title="getTabTitle(tab, index)"
        class="h-100"
        :title-item-class="getTitleItemClass(index)"
        :title-link-class="getTitleItemClass(index)"
        no-body
        :lazy="isTabLazy(tab)"
      >
        <page-block-tab
          v-if="tab.block"
          v-bind="{ ...$attrs, ...$props, page, block: tab.block, blockIndex: index, boundingRect: { xywh: block.xywh} }"
          :record="record"
          :module="module"
        />

        <div
          v-else
          class="d-flex h-100 align-items-center justify-content-center"
        >
          <p class="mb-0">
            {{ $t('tabs.noBlock') }}
          </p>
        </div>
      </b-tab>
    </b-tabs>
  </wrap>
</template>

<script>
import base from './base'
import { compose } from '@cortezaproject/corteza-js'
import { fetchID } from 'corteza-webapp-compose/src/lib/block'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'TabBase',

  components: {
    PageBlockTab: () => import('corteza-webapp-compose/src/components/PageBlocks'),
  },
  extends: base,

  computed: {
    tabbedBlocks () {
      return this.block.options.tabs.map(({ blockID, title }) => {
        let block = this.page.blocks.find(b => fetchID(b) === blockID)
        block = block ? compose.PageBlockMaker(block) : undefined

        // Blocks should display as Plain, to avoid card shadow/border
        if (block) {
          block.style.wrap.kind = 'Plain'
        }

        return {
          block: block ? compose.PageBlockMaker(block) : undefined,
          title,
        }
      })
    },

    contentClass () {
      return `overflow-hidden mh-100 ${this.block.options.style.orientation === 'vertical' ? 'd-block' : 'flex-fill'}`
    },

    navClass () {
      const { orientation } = this.block.options.style
      const style = orientation === 'vertical' ? 'px-3' : 'px-2'
      return `bg-white ${style}`
    },

    navWrapperClass () {
      const { orientation, position } = this.block.options.style
      let border = 'border-bottom'
      let style = 'bg-white mh-100'

      if (orientation === 'vertical') {
        border = position === 'end' ? 'border-left' : 'border-right'
        style = `${style} overflow-auto`
      } else if (position === 'end') {
        border = 'border-top'
      }

      return `${border} ${style}`
    },
  },

  methods: {
    getTitleItemClass (index) {
      const { fillJustify, alignment } = this.block.options.style
      return `order-${index} text-truncate text-${alignment} ${fillJustify !== 'none' ? 'flex-fill' : ''}`
    },

    getTabTitle ({ title, block = {} }, tabIndex) {
      const { title: blockTitle, kind } = block
      return title || blockTitle || kind || `${this.$t('tabs.tab')} ${tabIndex + 1}`
    },

    isTabLazy ({ block = {} }) {
      const { kind } = block
      return ['Calendar', 'Metric', 'Geometry'].includes(kind)
    },
  },
}
</script>
