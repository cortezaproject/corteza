<template>
  <div
    v-if="page"
    id="page-builder"
    ref="pageBuilder"
    class="flex-grow-1 overflow-auto d-flex px-2 w-100"
    tabIndex="1"
  >
    <portal to="topbar-title">
      {{ title }}
    </portal>

    <portal to="topbar-tools">
      <b-form-select
        v-if="layout"
        ref="layoutSelect"
        size="sm"
        :value="layout.pageLayoutID"
        value-field="pageLayoutID"
        text-field="label"
        :options="layouts"
        style="width: 300px;"
        @change="switchLayout"
      />

      <b-button-group
        v-if="page.canUpdatePage"
        size="sm"
        class="ml-2 text-nowrap"
      >
        <b-button
          variant="primary"
          class="d-flex align-items-center"
          :disabled="!pageViewer"
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
          style="margin-left:2px;"
        />
        <b-button
          variant="primary"
          :title="$t('tooltip.edit.page')"
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

    <grid
      :blocks.sync="blocks"
      editable
      @item-updated="onBlockUpdated"
    >
      <template
        slot-scope="{ boundingRect, block, index }"
      >
        <div
          :data-test-id="`block-${block.kind}`"
          class="h-100 editable-block"
        >
          <div
            class="toolbox border-0 p-2 m-0 text-light text-center"
            data-test-id="block-toolbox"
          >
            <div
              v-if="unsavedBlocks.has(index)"
              :title="$t('tooltip.unsavedChanges')"
              class="btn border-0"
            >
              <font-awesome-icon
                :icon="['fas', 'exclamation-triangle']"
                class="text-warning"
              />
            </div>

            <b-button-group>
              <b-button
                data-test-id="button-edit"
                :title="$t('tooltip.edit.block')"
                variant="outline-light"
                class="border-0"
                @click="editBlock(index)"
              >
                <font-awesome-icon
                  :icon="['far', 'edit']"
                />
              </b-button>

              <b-button
                :title="$t('tooltip.clone.block')"
                variant="outline-light"
                class="border-0"
                @click="cloneBlock(index)"
              >
                <font-awesome-icon
                  :icon="['far', 'clone']"
                />
              </b-button>

              <b-button
                :title="$t('tooltip.copy.block')"
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
              :title="$t('tooltip.delete.block')"
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
            :bounding-rect="boundingRect"
            editable
            class="p-2"
          />
        </div>
      </template>
    </grid>

    <b-modal
      id="createBlockSelector"
      size="lg"
      scrollable
      :title="$t('build.selectBlockTitle')"
    >
      <new-block-selector
        :record-page="!!module"
        :existing-blocks="selectableExistingBlocks"
        @select="addBlock"
      />

      <template #modal-footer>
        {{ $t('block:selectBlockFootnote') }}
      </template>
    </b-modal>

    <b-modal
      :title="$t('block.general.title')"
      :ok-title="$t('build.addBlock')"
      ok-variant="primary"
      cancel-variant="link"
      :cancel-title="$t('block.general.label.cancel')"
      size="xl"
      :visible="showCreator"
      body-class="p-0 border-top-0"
      header-class="p-3 pb-0 border-bottom-0"
      @ok="updateBlocks()"
      @hide="editor = undefined"
    >
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
      :title="$t('changeBlock')"
      size="xl"
      :visible="showEditor"
      body-class="p-0 border-top-0"
      footer-class="d-flex justify-content-between"
      header-class="p-3 pb-0 border-bottom-0"
      @hide="editor = undefined"
    >
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
          :title="$t('label.delete')"
          :borderless="false"
          @confirmed="deleteBlock(editor.index)"
        >
          {{ $t('label.delete') }}
        </c-input-confirm>

        <div>
          <b-button
            variant="link"
            :title="$t('label.cancel')"
            class="text-decoration-none"
            @click="cancel()"
          >
            {{ $t('label.cancel') }}
          </b-button>

          <b-button
            variant="primary"
            :title="$t('label.saveAndClose')"
            @click="updateBlocks()"
          >
            {{ $t('label.saveAndClose') }}
          </b-button>
        </div>
      </template>
    </b-modal>

    <portal to="admin-toolbar">
      <editor-toolbar
        :back-link="{name: 'admin.pages'}"
        :hide-save="!page.canUpdatePage"
        :disable-clone="disableClone"
        :disable-save="processing"
        :clone-tooltip="cloneTooltip"
        @save="handleSaveLayout()"
        @delete="handleDeleteLayout"
        @saveAndClose="handleSaveLayout({ closeOnSuccess: true })"
        @clone="handleClone()"
      >
        <b-button
          v-if="page.canUpdatePage"
          v-b-modal.createBlockSelector
          data-test-id="button-add-block"
          variant="light"
          size="lg"
          class="mr-1 float-right"
        >
          + {{ $t('build.addBlock') }}
        </b-button>
      </editor-toolbar>
    </portal>

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
import MagnificationModal from 'corteza-webapp-compose/src/components/Public/Page/Block/Modal'
import { fetchID } from 'corteza-webapp-compose/src/lib/block'

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
    MagnificationModal,
  },

  mixins: [
    pages,
  ],

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
      processing: false,

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

    title () {
      let { title = '', handle } = this.page || {}
      title = title || handle

      return this.$t('label.pageBuilder') + ' - ' + (title ? `"${title}"` : this.$t('label.noHandle'))
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
      if (this.module) {
        return undefined
      }

      return { name: 'page', params: { pageID: this.pageID } }
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

    disableClone () {
      return !!this.module
    },

    cloneTooltip () {
      return this.disableClone ? this.$t('tooltip.saveAsCopy') : ''
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
          const { blockID } = this.page.blocks.find(({ blockID }) => blockID === tab.blockID) || {}
          if (blockID) {
            tabbedIDs.add(blockID)
          }
        })
      })

      return [
        ...this.page.blocks.filter(({ blockID }) => tabbedIDs.has(blockID)),
        ...this.blocks.filter(({ blockID }) => !tabbedIDs.has(blockID)),
      ]
    },
  },

  watch: {
    pageID: {
      immediate: true,
      handler (pageID) {
        this.page = undefined
        this.layout = undefined
        this.layouts = []
        this.unsavedBlocks.clear()

        if (pageID) {
          const { namespaceID, name } = this.namespace
          this.findPageByID({ namespaceID, pageID, force: true }).then(page => {
            document.title = [page.title, name, this.$t('general:label.app-name.private')].filter(v => v).join(' | ')
            this.page = page.clone()
            return this.fetchPageLayouts()
          }).then(() => {
            this.setLayout()
          })
        }
      },
    },
  },

  created () {
    this.$root.$on('tab-editRequest', this.fulfilEditRequest)
    this.$root.$on('tab-createRequest', this.fulfilCreateRequest)
    this.$root.$on('tabChange', this.untabBlock)
  },

  mounted () {
    window.addEventListener('paste', this.pasteBlock)
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedBlocks(next)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedBlocks(next)
  },

  destroyed () {
    window.removeEventListener('paste', this.pasteBlock)
    this.$root.$off('tab-editRequest', this.fulfilEditRequest)
    this.$root.$off('tab-createRequest', this.fulfilCreateRequest)
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

      // Changes meta.hidden property to false, for all blocks that were tabbed only in the deleted block
      if (this.blocks[index].kind === 'Tabs') {
        const tabbedBlocks = this.blocks.filter((block) => block.kind === 'Tabs' && fetchID(block) !== fetchID(this.blocks[index]))
          .map(({ options }) => options.tabs).flat().reduce((unique, o) => {
            if (!unique.some(tab => tab.blockID === o.blockID)) {
              unique.push(o)
            }
            return unique
          }, []).map(({ blockID }) => blockID)

        this.blocks[index].options.tabs.forEach(({ blockID }) => {
          if (tabbedBlocks.includes(blockID)) return
          const index = this.blocks.findIndex((b) => fetchID(b) === blockID)
          if (index === -1) return
          this.blocks[index].meta.hidden = false
          this.calculateNewBlockPosition(this.blocks[index])
        })
      }

      this.blocks.splice(index, 1)

      if (this.editor) this.editor = undefined
      this.unsavedBlocks.add(index)
    },

    onBlockUpdated (index) {
      this.unsavedBlocks.add(index)
    },

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
        this.unsavedBlocks.add(this.editor.index)
      } else {
        this.blocks.push(block)
        this.unsavedBlocks.add(this.blocks.length - 1)
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
      }

      if (this.editor.block.kind === block.kind) {
        this.editor = undefined
      }
    },

    cloneBlock (index) {
      this.appendBlock(this.blocks[index].clone(), this.$t('notification:page.cloneSuccess'))
    },

    appendBlock (block, msg) {
      this.calculateNewBlockPosition(block)

      this.editor = { index: undefined, block: compose.PageBlockMaker(block) }
      this.updateBlocks()

      if (!this.editor) {
        msg && this.toastSuccess(msg)
        return true
      } else {
        msg && this.toastErrorHandler(this.$t('notification:page.duplicateFailed'))
        return false
      }
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

    async handleSaveLayout ({ closeOnSuccess = false, previewOnSuccess = false } = {}) {
      const { namespaceID } = this.namespace

      // Record blocks
      if (this.module && !this.validateModuleFieldSelection(this.module, this.page)) {
        this.toastErrorHandler(this.$t('notification:page.saveFailedRequired'))()
        return
      }

      // Inline record lists
      const queue = []
      this.usedBlocks.forEach((b, index) => {
        if (b.kind === 'RecordList' && b.options.editable) {
          const p = new Promise((resolve) => {
            const recordListUniqueID = [this.page.pageID, (this.record || {}).recordID || NoID, b.blockID].map(v => v || NoID).join('-')
            this.$root.$emit(`page-block:validate:${recordListUniqueID}`, resolve)
          })

          queue.push(p)
        }
      })

      const validated = await Promise.all(queue)
      if (validated.find(({ valid }) => !valid)) {
        this.toastErrorHandler(this.$t('notification:page.saveFailedRequired'))()
        return
      }

      this.processing = true

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
      }).then(({ page, layout }) => {
        this.page = new compose.Page(page)
        this.layout = new compose.PageLayout(layout)
        this.fetchPageLayouts()
        this.$route.query.layoutID = layout.pageLayoutID
        this.unsavedBlocks.clear()
        this.toastSuccess(this.$t('notification:page.page-layout.save.success'))

        if (closeOnSuccess) {
          this.$router.push({ name: 'admin.pages' })
        } else if (previewOnSuccess) {
          this.$router.push({ name: 'page', params: { pageID: this.pageID } })
        }
      }).finally(() => {
        this.processing = false
      }).catch(this.toastErrorHandler(this.$t('notification:page.page-layout.save.failed')))
    },

    validateModuleFieldSelection (module, page) {
      // Find all required fields
      const req = new Set(module.fields.filter(({ isRequired = false }) => isRequired).map(({ name }) => name))

      // Check if all required fields are there
      for (const b of page.blocks) {
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

    handleDeleteLayout () {
      this.processing = true

      this.deletePageLayout({ ...this.layout }).then(() => {
        this.setLayout()
        this.toastSuccess(this.$t('notification:page.page-layout.delete.success'))
      }).finally(() => {
        this.processing = false
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
          const block = JSON.parse(paste)
          const valid = this.isValid(block)

          if (valid) {
            this.appendBlock(block, this.$t('notification:page.pasteSuccess'))
          }
        } catch (error) {
          this.toastWarning(this.$t('notification:page.invalidBlock'))
          console.log(error)
        }
      }
    },

    // Trigger browser dialog on page leave to prevent unsaved changes
    checkUnsavedBlocks (next) {
      next(!this.unsavedBlocks.size || window.confirm(this.$t('build.unsavedChanges')))
    },

    async setLayout () {
      if (this.$route.query.layoutID) {
        const { namespaceID } = this.namespace
        this.layout = await this.findLayoutByID({ namespaceID, pageID: this.pageID, pageLayoutID: this.$route.query.layoutID })
      }

      this.layout = this.layout || this.layouts[0]
      if (!this.layout) {
        this.toastWarning('No layout, create one to edit it')
        return this.$router.push(this.pageEditor)
      }

      if (this.$route.query.layoutID !== this.layout.pageLayoutID) {
        this.$router.replace({ ...this.$route, query: { ...this.$route.query, layoutID: this.layout.pageLayoutID } })
      }

      this.unsavedBlocks.clear()

      const tempBlocks = []
      const { blocks = [] } = this.layout || {}

      blocks.forEach(({ blockID, xywh }) => {
        let block = this.page.blocks.find(b => b.blockID === blockID)
        block.xywh = xywh
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
      })

      this.blocks = tempBlocks
    },

    switchLayout (layoutID) {
      this.$router.push({ ...this.$route, query: { ...this.$route.query, layoutID } }).then(() => {
        this.setLayout()
      }).catch(() => {
        // Change layout value of select back to previous one if redirect was canceled
        this.$refs.layoutSelect.localValue = this.layout.pageLayoutID
      })
    },
  },
}
</script>
<style lang="scss">
div.toolbox {
  position: absolute;
  background-color: $dark;
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
