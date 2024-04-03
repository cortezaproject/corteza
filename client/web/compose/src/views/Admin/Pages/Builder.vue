<template>
  <div
    v-if="page"
    id="page-builder"
    ref="pageBuilder"
    class="flex-grow-1 overflow-auto d-flex p-2 w-100"
    tabIndex="1"
  >
    <portal to="topbar-title">
      {{ title }}
    </portal>

    <portal to="topbar-tools">
      <c-input-select
        v-if="layout && layouts.length > 1"
        ref="layoutSelect"
        :value="layout.pageLayoutID"
        :options="layouts"
        :reduce="layout => layout.pageLayoutID"
        size="sm"
        style="min-width: 250px; max-width: 300px;"
        @input="setLayout"
      />

      <b-button-group
        v-if="page.canUpdatePage"
        size="sm"
        class="ml-2 text-nowrap"
      >
        <b-button
          variant="primary"
          class="d-flex align-items-center"
          :to="pageViewer"
        >
          {{ $t('navigation.viewPage') }}
          <font-awesome-icon
            :icon="['far', 'eye']"
            class="ml-2"
          />
        </b-button>

        <page-translator
          :page.sync="trPage"
          :page-layout.sync="layout"
          button-variant="primary"
          style="margin-left:2px;"
        />

        <b-button
          v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.edit.page'), container: '#body' }"
          variant="primary"
          :to="pageEditor"
          class="d-flex align-items-center"
          style="margin-left:2px;"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <div
      v-if="processingLayout"
      class="d-flex align-items-center justify-content-center w-100"
    >
      <b-spinner />
    </div>

    <grid
      v-else-if="layout"
      :blocks.sync="blocks"
      editable
      @item-updated="onBlockUpdated"
    >
      <template
        slot-scope="{ index, block, resizing }"
      >
        <div
          :data-test-id="`block-${block.kind}`"
          class="h-100"
        >
          <div
            class="toolbox border-0 p-2 m-0 text-light text-center"
            data-test-id="block-toolbox"
          >
            <div
              v-if="unsavedBlocks.has(block.blockID !== '0' ? block.blockID : block.meta.tempID)"
              v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.unsavedChanges'), container: '#body' }"
              class="btn border-0"
            >
              <font-awesome-icon
                :icon="['fas', 'exclamation-triangle']"
                class="text-warning"
              />
            </div>

            <b-button-group>
              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.edit.block'), container: '#body' }"
                data-test-id="button-edit"
                variant="outline-light"
                class="border-0"
                @click="editBlock(index)"
              >
                <font-awesome-icon
                  :icon="['far', 'edit']"
                />
              </b-button>

              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.clone.block'), container: '#body' }"
                variant="outline-light"
                class="border-0"
                @click="cloneBlock(index)"
              >
                <font-awesome-icon
                  :icon="['far', 'clone']"
                />
              </b-button>

              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.copy.block'), container: '#body' }"
                variant="outline-light"
                class="border-0"
                @click="copyBlock(index)"
              >
                <font-awesome-icon
                  :icon="['far', 'copy']"
                />
              </b-button>
            </b-button-group>

            <c-input-confirm
              :tooltip="$t('tooltip.delete.block')"
              show-icon
              link
              size="md"
              class="ml-1"
              @confirmed="deleteBlock(index)"
            />
          </div>

          <page-block
            v-bind="{
              ...$attrs,
              ...$props
            }"
            :page="page"
            :blocks="usedBlocks"
            :block-index="index"
            :block="block"
            :module="module"
            :record="record"
            :resizing="resizing"
            :unsaved-blocks="unsavedBlocks"
            editable
            class="p-2"
            @edit-block="editBlock"
            @clone-block="cloneTabbedBlock"
            @copy-block="copyBlock"
            @delete-tab="deleteTab"
          />
        </div>
      </template>
    </grid>

    <b-modal
      id="createBlockSelector"
      :title="$t('build.selectBlockTitle')"
      size="lg"
      scrollable
    >
      <new-block-selector
        :record-page="!!module"
        :existing-blocks="selectableExistingBlocks"
        style="max-height: 75vh;"
        @select="addBlock"
      />

      <template #modal-footer>
        {{ $t('block:selectBlockFootnote') }}
      </template>
    </b-modal>

    <b-modal
      scrollable
      :ok-title="$t('build.addBlock')"
      ok-variant="primary"
      :ok-disabled="blockEditorOkDisabled"
      cancel-variant="light"
      :cancel-title="$t('block.general.label.cancel')"
      size="xl"
      :visible="showCreator"
      body-class="p-0 border-top-0"
      header-class="p-3 pb-0 border-bottom-0"
      @ok="updateBlocks()"
      @hide="editor = undefined"
    >
      <template #modal-title>
        <div class="d-flex gap-1 align-items-center">
          <h5 class="mb-0">
            {{ $t('block.general.title') }}
          </h5>
          <font-awesome-icon
            v-if="isEditorBlockReferenced"
            v-b-tooltip.noninteractive.hover.right="{ title: $t('referencedBlock'), container: '#body' }"
            :icon="['fas', 'exclamation-circle']"
            class="text-warning"
          />
        </div>
      </template>

      <configurator
        v-if="showCreator"
        :namespace="namespace"
        :module="module"
        :page="page"
        :blocks="usedBlocks"
        :block.sync="editor.block"
        :record="record"
      />
    </b-modal>

    <b-modal
      scrollable
      size="xl"
      :visible="showEditor"
      body-class="p-0 border-top-0"
      footer-class="d-flex justify-content-between"
      header-class="p-3 pb-0 border-bottom-0"
      @hide="editor = undefined"
    >
      <template #modal-title>
        <div class="d-flex gap-1 align-items-center">
          <h5 class="mb-0">
            {{ $t('changeBlock') }}
          </h5>
          <font-awesome-icon
            v-if="isEditorBlockReferenced"
            v-b-tooltip.noninteractive.hover.right="{ title: $t('referencedBlock'), container: '#body' }"
            :icon="['fas', 'exclamation-circle']"
            class="text-warning"
          />
        </div>
      </template>

      <configurator
        v-if="showEditor"
        :namespace="namespace"
        :module="module"
        :page="page"
        :blocks="usedBlocks"
        :block.sync="editor.block"
        :block-index="editor.index"
        :record="record"
      />

      <template #modal-footer="{ cancel }">
        <c-input-confirm
          size="md"
          size-confirm="md"
          variant="danger"
          :text="$t('label.delete')"
          :tooltip="$t('label.delete')"
          @confirmed="deleteBlock(editor.index)"
        />

        <div>
          <b-button
            variant="light"
            class="mr-2"
            @click="cancel()"
          >
            {{ $t('label.cancel') }}
          </b-button>

          <b-button
            variant="primary"
            :title="$t('label.saveAndClose')"
            :disabled="blockEditorOkDisabled"
            @click="updateBlocks()"
          >
            {{ $t('label.saveAndClose') }}
          </b-button>
        </div>
      </template>
    </b-modal>

    <portal to="admin-toolbar">
      <editor-toolbar
        :hide-save="!page.canUpdatePage"
        :processing="processing"
        :processing-save="processingSave"
        :processing-save-and-close="processingSaveAndClose"
        :processing-delete="processingDelete"
        :processing-clone="processingClone"
        hide-clone
        @save="handleSaveLayout()"
        @delete="handleDeleteLayout()"
        @saveAndClose="handleSaveLayout({ closeOnSuccess: true })"
        @back="$router.push(previousPage || { name: 'admin.pages' })"
      >
        <b-button
          v-if="page.canUpdatePage"
          v-b-modal.createBlockSelector
          data-test-id="button-add-block"
          variant="light"
          size="lg"
        >
          + {{ $t('build.addBlock') }}
        </b-button>

        <template #saveAsCopy>
          <b-dropdown
            v-if="page.canUpdatePage"
            data-test-id="dropdown-saveAsCopy"
            :text="$t('general:label.saveAsCopy')"
            :disabled="processing"
            size="lg"
            variant="light"
          >
            <b-dropdown-item
              data-test-id="dropdown-item-saveAsCopy-ref"
              @click="handleCloneLayout({ ref: true })"
            >
              {{ $t('build.saveAsCopy.ref') }}
            </b-dropdown-item>
            <b-dropdown-item
              data-test-id="dropdown-item-saveAsCopy-noRef"
              @click="handleCloneLayout({ ref: false })"
            >
              {{ $t('build.saveAsCopy.noRef') }}
            </b-dropdown-item>
          </b-dropdown>
        </template>
      </editor-toolbar>
    </portal>

    <record-modal
      :namespace="namespace"
    />

    <magnification-modal
      :namespace="namespace"
    />
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import pages from 'corteza-webapp-compose/src/mixins/pages'
import NewBlockSelector from 'corteza-webapp-compose/src/components/Admin/Page/Builder/Selector'
import PageTranslator from 'corteza-webapp-compose/src/components/Admin/Page/PageTranslator'
import Grid from 'corteza-webapp-compose/src/components/Common/Grid'
import PageBlock from 'corteza-webapp-compose/src/components/PageBlocks'
import EditorToolbar from 'corteza-webapp-compose/src/components/Admin/EditorToolbar'
import { compose, NoID } from '@cortezaproject/corteza-js'
import Configurator from 'corteza-webapp-compose/src/components/PageBlocks/Configurator'
import RecordModal from 'corteza-webapp-compose/src/components/Public/Record/Modal'
import MagnificationModal from 'corteza-webapp-compose/src/components/Public/Page/Block/Modal'
import { fetchID } from 'corteza-webapp-compose/src/lib/block'
import { handle } from '@cortezaproject/corteza-vue'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  components: {
    Configurator,
    Grid,
    NewBlockSelector,
    PageBlock,
    EditorToolbar,
    PageTranslator,
    RecordModal,
    MagnificationModal,
  },

  mixins: [
    pages,
  ],

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedBlocks(next, to)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedBlocks(next)
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    pageID: {
      type: String,
      required: true,
    },
  },

  data () {
    return {
      title: '',

      processing: false,
      processingSave: false,
      processingSaveAndClose: false,
      processingClone: false,
      processingDelete: false,

      processingLayout: false,

      page: undefined,
      layout: undefined,
      layouts: [],

      blocks: [],

      editor: undefined,

      unsavedBlocks: new Set(),
    }
  },

  computed: {
    ...mapGetters({
      pages: 'page/set',
      getModuleByID: 'module/getByID',
      previousPage: 'ui/previousPage',
    }),

    trPage: {
      get () {
        if (!this.page) {
          return new compose.Page()
        }
        return this.page
      },
      set (v) {
        this.page = v
        this.updatePageSet(v)
      },
    },

    showEditor () {
      return this.editor && this.editor.index !== undefined
    },

    showCreator () {
      return this.editor && this.editor.index === undefined
    },

    module () {
      if (this.page && this.page.moduleID !== NoID) {
        return this.$store.getters['module/getByID'](this.page.moduleID)
      } else {
        return undefined
      }
    },

    /**
     * Create a dummy record object when we are editing a record page.
     * This enables compose:record triggers & Record page blocks
     */
    record () {
      if (this.module) {
        return new compose.Record({}, this.module)
      }
      return null
    },

    pageViewer () {
      const name = this.module ? 'page.record.create' : 'page'
      return { name, params: { pageID: this.pageID } }
    },

    pageEditor () {
      return { name: 'admin.pages.edit', params: { pageID: this.pageID } }
    },

    hasChildren () {
      return this.page ? this.pages.some(({ selfID }) => selfID === this.page.pageID) : false
    },

    hideDelete () {
      return this.hasChildren || !this.page.canDeletePage || !!this.page.deletedAt
    },

    selectableExistingBlocks () {
      return this.page.blocks.filter(({ blockID }) => !this.usedBlocks.some(b => b.blockID === blockID))
    },

    // Blocks used on page or tabbed
    usedBlocks () {
      const tabbedIDs = new Set()

      // If tab is not on layout include it
      this.blocks.forEach(block => {
        if (block.kind !== 'Tabs') return

        const { tabs = [] } = block.options
        tabs.forEach(tab => {
          if (this.blocks.some(({ blockID }) => blockID === tab.blockID)) return
          const { blockID } = this.page.blocks.find(({ blockID }) => blockID === tab.blockID) || {}
          if (blockID) {
            tabbedIDs.add(blockID)
          }
        })
      })

      return [
        ...this.blocks.filter(({ blockID }) => !tabbedIDs.has(blockID)),
        ...this.page.blocks.filter(({ blockID }) => tabbedIDs.has(blockID)),
      ]
    },

    // Set of blockIDs used on other layouts
    otherLayoutBlockIDs () {
      const set = new Set()

      return this.layouts.reduce((acc, { blocks, pageLayoutID }) => {
        if (pageLayoutID === this.layout.pageLayoutID) return acc

        blocks.forEach(({ blockID }) => acc.add(blockID))

        return acc
      }, set)
    },

    // Is block open in editor referenced on other layouts
    isEditorBlockReferenced () {
      const { block } = this.editor || {}
      if (!block || block.blockID === NoID) return

      return this.otherLayoutBlockIDs.has(this.editor.block.blockID)
    },

    blockEditorOkDisabled () {
      if (!this.editor) return true

      const { block } = this.editor

      if (!block) return true

      const { customCSSClass, customID } = block.meta

      return [handle.handleState(customID), handle.classState(customCSSClass)].includes(false)
    },
  },

  watch: {
    pageID: {
      immediate: true,
      handler (pageID) {
        this.processingLayout = true

        this.unsavedBlocks.clear()
        this.layouts = []
        this.layout = undefined

        const { namespaceID, name } = this.namespace
        this.findPageByID({ namespaceID, pageID, force: true }).then(page => {
          let { title = '', handle } = page
          title = title || handle
          this.title = `${this.$t('label.pageBuilder')} - "${title}"`
          document.title = [page.title, name, this.$t('general:label.app-name.private')].filter(v => v).join(' | ')
          this.page = page.clone()
          return this.fetchPageLayouts().then(() => {
            this.setLayout()
          })
        }).catch(() => {
          this.processingLayout = false
        })
      },
    },
  },

  mounted () {
    window.addEventListener('paste', this.pasteBlock)
    this.$root.$on('tab-editRequest', this.fulfilEditRequest)
    this.$root.$on('tab-createRequest', this.fulfilCreateRequest)
    this.$root.$on('tabChange', this.untabBlock)
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      findPageByID: 'page/findByID',
      updatePage: 'page/update',
      deletePage: 'page/delete',
      updatePageSet: 'page/updateSet',
      createPage: 'page/create',
      loadPages: 'page/load',
      findLayoutByID: 'pageLayout/findByID',
      findLayoutsByPageID: 'pageLayout/findByPageID',
      createPageLayout: 'pageLayout/create',
      updatePageLayout: 'pageLayout/update',
      deletePageLayout: 'pageLayout/delete',
    }),

    fulfilEditRequest (blockID) {
      // this ensures whatever changes in tabs is not lost before we lose its configurator
      // because we are reusing that modal component
      this.updateBlocks()

      const blockIndex = this.blocks.findIndex(block => fetchID(block) === blockID)
      if (blockIndex > -1) {
        this.editBlock(blockIndex)
      }
    },

    fulfilCreateRequest (block) {
      this.updateBlocks(block)
    },

    untabBlock (block) {
      const where = this.tabLocation(block)

      if (!where.length) return

      where.forEach(({ block, index }) => {
        const { tabs } = block.options
        tabs.splice(index, 1)
      })
    },

    tabLocation (tabbedBlock) {
      const where = []
      this.blocks.forEach((block, i) => {
        if (block.kind !== 'Tabs') return
        const { tabs } = block.options
        const index = tabs.findIndex(({ blockID }) => blockID === fetchID(tabbedBlock))
        where.push({ block, index })
      })
      return where
    },

    addBlock (block, index = undefined) {
      this.$bvModal.hide('createBlockSelector')
      this.calculateNewBlockPosition(block)
      this.editor = { index, block: compose.PageBlockMaker(block) }
    },

    editBlock (index = undefined) {
      this.$nextTick(() => {
        this.editor = { index, block: compose.PageBlockMaker(this.blocks[index]) }
      })
    },

    deleteBlock (index) {
      // If the deleted block is hidden, we need to remove it from the related tabs blocks if it is tabbed.
      if (this.blocks[index].meta.hidden) {
        this.blocks.forEach((block) => {
          if (block.kind !== 'Tabs' || !block.options.tabs.some(({ blockID }) => blockID === fetchID(this.blocks[index]))) return
          block.options.tabs = block.options.tabs.filter(({ blockID }) => blockID !== fetchID(this.blocks[index]))
        })
      }

      const block = this.blocks[index]

      this.blocks.splice(index, 1)

      if (block.blockID !== NoID) {
        this.unsavedBlocks.add(block.blockID)
      } else {
        this.unsavedBlocks.delete(block.meta.tempID)
      }

      if (block.kind === 'Tabs') {
        this.showUntabbedHiddenBlocks()
      }

      if (this.editor) this.editor = undefined
    },

    deleteTab ({ blockIndex, tabIndex }) {
      const { blockID } = this.blocks[blockIndex] || {}

      if (!blockID) return

      this.unsavedBlocks.add(blockID)
      this.blocks[blockIndex].options.tabs.splice(tabIndex, 1)

      this.showUntabbedHiddenBlocks()
    },

    // Changes meta.hidden property to false, for all blocks that are hidden but not in a tab
    showUntabbedHiddenBlocks () {
      const tabbedBlocks = new Set()

      this.blocks.forEach(block => {
        if (block.kind !== 'Tabs') return

        block.options.tabs.forEach(({ blockID }) => tabbedBlocks.add(blockID))
      })

      this.blocks.forEach((block, index) => {
        if (!block.meta.hidden || tabbedBlocks.has(fetchID(block))) return

        this.blocks[index].meta.hidden = false
        this.calculateNewBlockPosition(this.blocks[index])
      })

      tabbedBlocks.clear()
    },

    onBlockUpdated (index) {
      this.unsavedBlocks.add(fetchID(this.blocks[index]))
    },

    // When debugging this, make sure to remove the @hide event handle from the block editor/creator modals
    updateBlocks (block = this.editor.block) {
      block = compose.PageBlockMaker(block)

      const creatingTabbedBlock = this.editor.block.kind !== block.kind

      if (creatingTabbedBlock) {
        this.$root.$emit('builder-createRequestFulfilled', {
          blockID: fetchID(block),
          title: block.title,
        })
      }

      if (this.editor.index !== undefined && !creatingTabbedBlock) {
        const oldBlock = this.blocks[this.editor.index]

        if (oldBlock.meta.hidden === true && this.editor.block.meta.hidden === false) {
          this.untabBlock(this.editor.block)
          this.calculateNewBlockPosition(block)
        }

        this.blocks.splice(this.editor.index, 1, block)
        this.unsavedBlocks.add(fetchID(block))
      } else {
        this.blocks.push(block)
        this.unsavedBlocks.add(fetchID(block))
        this.scrollToBottom()
      }

      if (block.kind === 'Tabs') {
        block.options.tabs.forEach((tab) => {
          if (!tab.blockID) return
          let tabbedBlock = this.blocks.find(b => fetchID(b) === tab.blockID)

          if (!tabbedBlock) {
            tabbedBlock = this.page.blocks.find(({ blockID }) => blockID === tab.blockID)
            this.blocks.push(tabbedBlock)
          }

          tabbedBlock.meta.hidden = true
        })

        this.showUntabbedHiddenBlocks()
      }

      if (this.editor.block.kind === block.kind) {
        this.editor = undefined
      }
    },

    cloneBlock (index) {
      this.appendBlock(this.blocks[index].clone(), this.$t('notification:page.cloneSuccess'))
    },

    cloneTabbedBlock ({ tabbedBlockIndex, tabBlockIndex, title }) {
      const block = this.blocks[tabbedBlockIndex].clone()
      block.meta.hidden = true

      this.blocks[tabBlockIndex].options.tabs.push({
        blockID: fetchID(block),
        title,
      })

      this.blocks.push(block)
      this.unsavedBlocks.add(fetchID(block))
    },

    appendBlock (block, msg) {
      this.calculateNewBlockPosition(block)

      this.editor = { index: undefined, block }
      this.updateBlocks()

      if (!this.editor) {
        msg && this.toastSuccess(msg)
        return true
      } else {
        msg && this.toastErrorHandler(this.$t('notification:page.duplicateFailed'))
        return false
      }
    },

    isBlockUnsaved (block) {
      return this.unsavedBlocks.has(block.blockID)
    },

    calculateNewBlockPosition (block) {
      if (this.blocks.length) {
        // ensuring we append the block to the end of the page
        // eslint-disable-next-line
          const maxY = this.blocks.filter(({ meta }) => !meta.hidden).map((block) => block.xywh[1]).reduce((acc, val) => {
          return acc > val ? acc : val
        }, 0)
        block.xywh = [0, maxY + 2, 20, 15]
      }
    },

    async fetchPageLayouts () {
      const { namespaceID } = this.namespace

      return this.findLayoutsByPageID({ namespaceID, pageID: this.pageID }).then(layouts => {
        this.layouts = layouts.map(l => {
          l = new compose.PageLayout(l)
          l.label = l.meta.title || l.handle || l.pageLayoutID
          return l
        })
      })
    },

    checkRequiredRecordFields () {
      // Find all required fields
      const req = new Set(this.module.fields.filter(({ isRequired = false }) => isRequired).map(({ name }) => name))

      // Check if all required fields are there
      for (const b of this.usedBlocks) {
        if (b.kind !== 'Record') {
          continue
        }

        // If no fields are in Record block, means all fields are present(default), no need to check
        if (!b.options || !b.options.fields.length) {
          return true
        }

        for (const f of b.options.fields) {
          req.delete(f.name)
        }
      }

      // If required fields are satisfied, then the validation passes
      return !req.size
    },

    async handleSaveLayout ({ closeOnSuccess = false, previewOnSuccess = false, alert = true } = {}) {
      const { namespaceID } = this.namespace

      // Record blocks
      if (this.module && !this.checkRequiredRecordFields()) {
        this.toastErrorHandler(this.$t('notification:page.saveFailedRequired'))()
        return
      }

      // Inline record lists
      this.usedBlocks.forEach((b, index) => {
        if (b.kind === 'RecordList' && b.options.editable) {
          const recordListModule = this.getModuleByID(b.options.moduleID)
          const req = new Set(recordListModule.fields.filter(({ isRequired = false }) => isRequired).map(({ name }) => name))

          // Check if all required fields are there
          for (const f of b.options.editFields) {
            req.delete(f.name)
          }

          if (req.size) {
            this.toastErrorHandler(this.$t('notification:page.saveFailedRequired'))()
          }
        }
      })

      this.processing = true

      if (closeOnSuccess) {
        this.processingSaveAndClose = true
      } else {
        this.processingSave = true
      }

      return Promise.all([
        this.findPageByID({ ...this.page, force: true }),
        this.findLayoutByID({ ...this.layout }),
      ]).then(([page, layout]) => {
        const blocks = [
          ...page.blocks.filter(({ blockID }) => {
            // Check if block exists in any other layout, if not delete it permanently
            return !this.blocks.some(b => b.blockID === blockID) && this.layouts.some(({ pageLayoutID, blocks }) => pageLayoutID !== layout.pageLayoutID && blocks.some(b => b.blockID === blockID))
          }),
          ...this.blocks,
        ]

        return this.updatePage({ namespaceID, ...page, blocks })
          .then(this.updateTabbedBlockIDs)
          .then(async page => {
            const blocks = this.blocks.map(({ blockID, meta, xywh }) => {
              if (blockID === NoID) {
                blockID = (page.blocks.find(block => block.meta.tempID === meta.tempID) || {}).blockID
              }

              return { blockID, xywh, meta }
            })
            layout = await this.updatePageLayout({ ...layout, blocks })
            return { page, layout }
          })
      }).then(async ({ page, layout }) => {
        this.unsavedBlocks.clear()

        if (closeOnSuccess) {
          this.$router.push(this.previousPage || { name: 'admin.pages' })
          return
        }

        if (alert) {
          this.toastSuccess(this.$t('notification:page.page-layout.save.success'))
        }

        this.page = new compose.Page(page)

        await this.fetchPageLayouts()
        this.setLayout(layout.pageLayoutID, false)
      }).finally(() => {
        this.processing = false

        if (closeOnSuccess) {
          this.processingSaveAndClose = false
        } else {
          this.processingSave = false
        }
      }).catch(this.toastErrorHandler(this.$t('notification:page.page-layout.save.failed')))
    },

    async handleCloneLayout ({ ref = false }) {
      this.processing = true
      this.processingLayout = true
      this.processingClone = true

      const layout = {
        ...this.layout.clone(),
        handle: '',
        weight: this.layouts.length + 1,
      }

      layout.meta.title = `${this.$t('copyOf')}${layout.meta.title}`

      // If we are cloning a layout with references, we need to clone the blocks
      if (!ref) {
        const oldBlockIDs = {}
        layout.blocks = []

        // Sort based on if tab or not
        this.blocks = this.blocks.toSorted((a, b) => {
          if (a.kind === 'Tabs' && b.kind !== 'Tabs') {
            return 1 // Move 'Tabs' to the end
          } else if (a.kind !== 'Tabs' && b.kind === 'Tabs') {
            return -1 // Keep 'Tabs' before other elements
          } else {
            return 0 // No change in order
          }
        }).map(block => {
          const oldBlockID = block.blockID

          if (block.kind === 'Tabs') {
            block.options.tabs = block.options.tabs.map(tab => {
              tab.blockID = oldBlockIDs[tab.blockID]
              return tab
            })
          }

          block = block.clone()
          oldBlockIDs[oldBlockID] = block.meta.tempID

          return block
        })
      }

      this.createPageLayout(layout).then(layout => {
        this.layout = layout
        this.layouts.push({ ...layout, label: layout.meta.title || layout.handle || layout.pageLayoutID })
        return this.handleSaveLayout({ alert: false })
      }).then(() => {
        this.toastSuccess(this.$t('notification:page.page-layout.clone.success'))
      }).finally(() => {
        this.processing = false
        this.processingLayout = false
        this.processingClone = false
      }).catch(this.toastErrorHandler(this.$t('notification:page.page-layout.clone.failed')))
    },

    handleDeleteLayout () {
      this.processing = true
      this.processingDelete = true

      this.deletePageLayout({ ...this.layout }).then(() => {
        return this.fetchPageLayouts()
      }).then(() => {
        this.setLayout()
        this.toastSuccess(this.$t('notification:page.page-layout.delete.success'))
      }).finally(() => {
        this.processing = false
        this.processingDelete = false
      }).catch(this.toastErrorHandler(this.$t('notification:page.page-layout.delete.failed')))
    },

    /**
     * Validates block, returns true if there are no problems with it
     *
     * @param {compose.PageBlock} block
     * @returns {boolean}
     */
    isValid (block) {
      if (typeof block.validate === 'function') {
        return block.validate().length === 0
      }

      return true
    },

    async copyBlock (index) {
      const block = JSON.stringify(this.blocks[index].clone())

      // Change tabbed blockID to use tempID's since they are persisted on save
      if (block.kind === 'Tabs') {
        const { tabs = [] } = block.options

        block.options.tabs = tabs.map(b => {
          const { tempID } = (this.blocks.find(({ blockID }) => blockID === b.blockID) || {}).meta || {}
          b.blockID = tempID
          return b
        })
      }

      navigator.clipboard.writeText(block).then(() => {
        this.toastSuccess(this.$t('notification:page.copySuccess'))
        this.$refs.pageBuilder.focus()
      },
      (err) => {
        this.toastErrorHandler(this.$t('notification:page.copyFailed', { reason: err }))
      })
    },

    pasteBlock (event) {
      // ensuring page-builder is focused before pasting a block
      if (document.querySelector('#page-builder') === document.activeElement) {
        event.preventDefault()
        const paste = (event.clipboardData || window.clipboardData).getData('text')
        // Doing this to handle JSON parse error
        try {
          const block = compose.PageBlockMaker(JSON.parse(paste))
          const valid = this.isValid(block)

          if (valid) {
            this.appendBlock(block, this.$t('notification:page.pasteSuccess'))
          }
        } catch (error) {
          this.toastWarning(this.$t('notification:page.invalidBlock'))
        }
      }
    },

    // Trigger browser dialog on page leave to prevent unsaved changes
    checkUnsavedBlocks (next, to = { query: {} }) {
      // Check if additional query params will be appended to url
      const queryParams = Object.keys(to.query).filter(key => key !== 'layoutID').length > 0
      next(!this.unsavedBlocks.size || queryParams || window.confirm(this.$t('build.unsavedChanges')))
    },

    async setLayout (layoutID, processing = true) {
      const oldLayoutID = this.$route.query.layoutID

      // Cancelable redirect
      if (layoutID && oldLayoutID !== layoutID) {
        try {
          await this.$router.replace({ ...this.$route, query: { ...this.$route.query, layoutID } })
        } catch {
          this.$refs.layoutSelect.localValue = oldLayoutID
          return
        }
      }

      if (processing) {
        this.processingLayout = true
      }

      layoutID = layoutID || this.$route.query.layoutID

      if (layoutID) {
        this.layout = this.layouts.find(({ pageLayoutID }) => pageLayoutID === layoutID)
      }

      this.layout = this.layout || this.layouts[0]
      if (!this.layout) {
        this.toastWarning(this.$t('notification:page.page-layout.notFound.edit'))
        return this.$router.push(this.pageEditor)
      }

      // If no previous layout was set and it exists, replace the URL with the proper layoutID
      if (this.$route.query.layoutID !== this.layout.pageLayoutID) {
        this.$router.replace({ ...this.$route, query: { ...this.$route.query, layoutID: this.layout.pageLayoutID } })
      }

      this.unsavedBlocks.clear()

      const tempBlocks = []
      const { blocks = [] } = this.layout || {}

      blocks.forEach(({ blockID, xywh, meta = {} }) => {
        let block = this.page.blocks.find(b => b.blockID === blockID)

        if (block) {
          block.xywh = xywh
          block.meta.hidden = !!meta.hidden
          tempBlocks.push(block)

          if (block.kind === 'Tabs') {
            const { tabs = [] } = block.options
            tabs.forEach(tab => {
              if (blocks.some(b => b.blockID === tab.blockID)) return

              block = this.page.blocks.find(b => b.blockID === tab.blockID)

              if (block) {
                tempBlocks.push(block)
              }
            })
          }
        }
      })

      this.blocks = tempBlocks

      setTimeout(() => {
        this.processingLayout = false
      }, 400)
    },

    scrollToBottom () {
      const pageBuilderElement = document.getElementById('page-builder')

      this.$nextTick(() => {
        pageBuilderElement.scrollTo({
          top: pageBuilderElement.scrollHeight,
          behavior: 'smooth',
        })
      })
    },

    setDefaultValues () {
      this.title = ''
      this.processing = false
      this.processingSaveAndClose = false
      this.processingSave = false
      this.processingClone = false
      this.processingLayout = false
      this.page = undefined
      this.layout = undefined
      this.layouts = []
      this.blocks = []
      this.editor = undefined
      this.unsavedBlocks.clear()
    },

    destroyEvents () {
      window.removeEventListener('paste', this.pasteBlock)
      this.$root.$off('tab-editRequest', this.fulfilEditRequest)
      this.$root.$off('tab-createRequest', this.fulfilCreateRequest)
      this.$root.$off('tabChange', this.untabBlock)
    },
  },
}
</script>
<style lang="scss">
div.toolbox {
  position: absolute;
  background-color: var(--secondary);
  bottom: 0;
  left: 0;
  z-index: 1001;
  border-top-right-radius: 10px;
  opacity: 0.5;
  pointer-events: none;

  &:hover {
    opacity: 1;
  }

  & * {
    pointer-events: auto;
  }
}

[dir="rtl"] {
  div.toolbox {
    left: 0;
    right: auto;
  }
}
</style>
