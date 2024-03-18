<template>
  <wrap
    v-bind="$props"
    :scrollable-body="false"
    card-class="tabs-base-block-container"
    header-class="border-0 border-white"
    v-on="$listeners"
  >
    <div
      v-if="!tabbedBlocks.length"
      class="d-flex h-100 align-items-center justify-content-center"
    >
      <p class="mb-0">
        {{ $t('noTabs') }}
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
        justified: block.options.style.justify === 'justify',
        pills: block.options.style.appearance === 'pills',
        tabs: block.options.style.appearance === 'tabs',
        small: block.options.style.appearance === 'small',
        vertical: block.options.style.orientation === 'vertical',
        end: block.options.style.position === 'end'
      }"
      no-fade
      class="h-100"
      :class="{ 'd-flex flex-column': block.options.style.orientation !== 'vertical' }"
    >
      <b-tab
        v-for="(tab, index) in tabbedBlocks"
        :key="`${getTabTitle(tab, index)}-${index}`"
        class="h-100"
        :title-item-class="getTitleItemClass(index)"
        :title-link-class="getTitleItemClass(index)"
        no-body
        :lazy="isTabLazy(tab)"
      >
        <template #title>
          <span>
            {{ getTabTitle(tab, index) }}
          </span>

          <div
            v-if="tab.block && editable"
            class="d-inline ml-3"
          >
            <div
              v-if="unsavedBlocks.has(tab.block.blockID !== '0' ? tab.block.blockID : tab.block.meta.tempID)"
              v-b-tooltip.noninteractive.hover="{ title: $t('unsavedChanges'), container: '#body' }"
              class="btn btn-sm border-0"
            >
              <font-awesome-icon
                :icon="['fas', 'exclamation-triangle']"
                class="text-warning"
              />
            </div>

            <b-button-group size="sm">
              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.edit'), container: '#body' }"
                variant="outline-light"
                class="text-primary border-0 toolbox-button"
                @click="editTabbedBlock(tab)"
              >
                <font-awesome-icon
                  :icon="['far', 'edit']"
                />
              </b-button>

              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.clone'), container: '#body' }"
                variant="outline-light"
                class="text-primary border-0 toolbox-button"
                @click="cloneTabbedBlock(tab)"
              >
                <font-awesome-icon
                  :icon="['far', 'clone']"
                />
              </b-button>

              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.copy'), container: '#body' }"
                variant="outline-light"
                class="text-primary border-0 toolbox-button"
                @click="copyTabbedBlock(tab)"
              >
                <font-awesome-icon
                  :icon="['far', 'copy']"
                />
              </b-button>
            </b-button-group>

            <c-input-confirm
              :tooltip="$t('tooltip.delete')"
              show-icon
              class="ml-1"
              @confirmed="deleteTab(index)"
            />
          </div>
        </template>

        <page-block-tab
          v-if="tab.block"
          v-bind="{ ...$attrs, ...$props, page, block: tab.block, blockIndex: index }"
          :record="record"
          :module="module"
          :magnified="magnified"
          header-class="border-0 border-white"
        />

        <div
          v-else-if="!tab.block"
          class="d-flex h-100 align-items-center justify-content-center"
        >
          <p class="mb-0">
            {{ $t('noBlock') }}
          </p>
        </div>
      </b-tab>
    </b-tabs>
  </wrap>
</template>

<script>
import base from './base'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { fetchID } from 'corteza-webapp-compose/src/lib/block'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  i18nOptions: {
    namespaces: 'block',
    keyPrefix: 'tabs',
  },

  name: 'TabBase',

  components: {
    PageBlockTab: () => import('corteza-webapp-compose/src/components/PageBlocks'),
  },

  extends: base,

  computed: {
    tabbedBlocks () {
      return this.block.options.tabs.reduce((acc, { blockID, title }) => {
        const unparsedBlock = blockID ? this.blocks.find(b => fetchID(b) === blockID) : undefined

        if (!unparsedBlock) {
          if (!blockID && title) {
            acc.push({ title })
          }

          return acc
        }

        let block = JSON.parse(JSON.stringify(unparsedBlock))

        // Blocks should display as Plain, to avoid card shadow/border
        block.style.wrap.kind = 'Plain'
        block.style.border.enabled = false
        block = compose.PageBlockMaker(block)

        acc.push({
          block,
          title,
        })

        return acc
      }, [])
    },

    contentClass () {
      return `overflow-hidden mh-100 ${this.block.options.style.orientation === 'vertical' ? 'd-block' : 'flex-fill'}`
    },

    navClass () {
      const { orientation } = this.block.options.style
      const style = orientation === 'vertical' ? 'px-3' : 'px-2'
      return `bg-transparent ${style}`
    },

    navWrapperClass () {
      const { orientation, position } = this.block.options.style
      let border = ''
      let style = 'bg-transparent mh-100'

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
    editTabbedBlock (tab) {
      const blockIndex = this.blocks.findIndex(block => fetchID(block) === fetchID(tab.block))
      if (blockIndex > -1) {
        this.$emit('edit-block', blockIndex)
      }
    },

    cloneTabbedBlock (tab) {
      const tabbedBlockIndex = this.blocks.findIndex(block => fetchID(block) === fetchID(tab.block))
      if (tabbedBlockIndex > -1) {
        this.$emit('clone-block', { tabbedBlockIndex, tabBlockIndex: this.blockIndex, title: tab.title })
      }
    },

    copyTabbedBlock (tab) {
      const tabbedBlockIndex = this.blocks.findIndex(block => fetchID(block) === fetchID(tab.block))
      if (tabbedBlockIndex > -1) {
        this.$emit('copy-block', tabbedBlockIndex)
      }
    },

    deleteTab (tabIndex) {
      this.$emit('delete-tab', { tabIndex, blockIndex: this.blockIndex })
    },

    getTitleItemClass (index) {
      const { justify, alignment } = this.block.options.style
      return `order-${index} text-truncate text-${alignment} ${justify !== 'none' ? 'flex-fill' : ''}`
    },

    getTabTitle ({ title = '', block = {} }, tabIndex) {
      const { title: blockTitle, kind } = block
      const interpolatedTitle = evaluatePrefilter(blockTitle, {
        record: this.record,
        user: this.$auth.user || {},
        recordID: (this.record || {}).recordID || NoID,
        ownerID: (this.record || {}).ownedBy || NoID,
        userID: (this.$auth.user || {}).userID || NoID,
      })

      title = evaluatePrefilter(title, {
        record: this.record,
        user: this.$auth.user || {},
        recordID: (this.record || {}).recordID || NoID,
        ownerID: (this.record || {}).ownedBy || NoID,
        userID: (this.$auth.user || {}).userID || NoID,
      })

      return title || interpolatedTitle || kind || `${this.$t('tab')} ${tabIndex + 1}`
    },

    isTabLazy ({ block = {} }) {
      const { kind } = block
      return ['Calendar', 'Metric', 'Geometry'].includes(kind)
    },
  },
}
</script>

<style lang="scss">
.tabs-base-block-container {
  .nav-pills {
    .active .toolbox-button {
      color: var(--white) !important;
    }
  }

  .tabs {
    .card {
      border-radius: 0;
    }

    .card-header {
      border-radius: 0;
    }
  }
}
</style>
