<template>
  <div
    v-if="!!page"
    id="page-builder"
    ref="pageBuilder"
    class="flex-grow-1 overflow-auto d-flex px-2 w-100"
    tabIndex="1"
  >
    <portal to="topbar-title">
      {{ title }}
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="page && page.canUpdatePage"
        size="sm"
        class="mr-1"
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
      :blocks="page.blocks"
      editable
      @change="updatePageBlockGrid"
      @item-updated="onBlockUpdated"
    >
      <template
        slot-scope="{ boundingRect, block, index }"
      >
        <div
          :data-test-id="`block-${block.kind}`"
          class="h-100 editable-block"
          :class="{ 'bg-warning': !isValid(block) }"
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
            <c-input-confirm
              class="ml-1"
              size="md"
              link
              @confirmed="deleteBlock(index)"
            />
          </div>

          <page-block
            v-bind="{ ...$attrs, ...$props, page, block, boundingRect, blockIndex: index, editable: true }"
            :record="record"
            :module="module"
            class="p-2"
          />
        </div>
      </template>
    </grid>

    <b-modal
      id="createBlockSelector"
      size="lg"
      scrollable
      hide-footer
      :title="$t('build.selectBlockTitle')"
    >
      <new-block-selector
        :record-page="!!module"
        @select="addBlock"
      />
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
      @ok="$event => updateBlocks()"
      @hide="editor = undefined"
    >
      <configurator
        v-if="showCreator"
        :namespace="namespace"
        :module="module"
        :page="page"
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
            size="md"
            :title="$t('label.cancel')"
            class="text-decoration-none"
            @click="cancel()"
          >
            {{ $t('label.cancel') }}
          </b-button>

          <b-button
            variant="primary"
            size="md"
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
        :hide-delete="hideDelete"
        :hide-save="!page.canUpdatePage"
        :disable-clone="disableClone"
        :clone-tooltip="cloneTooltip"
        @save="handleSave()"
        @delete="handleDeletePage"
        @saveAndClose="handleSave({ closeOnSuccess: true })"
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

        <template #delete>
          <b-dropdown
            v-if="showDeleteDropdown"
            data-test-id="dropdown-delete"
            size="lg"
            variant="danger"
            :text="$t('general:label.delete')"
          >
            <b-dropdown-item
              data-test-id="dropdown-item-delete-update-parent-of-sub-pages"
              @click="handleDeletePage('rebase')"
            >
              {{ $t('delete.rebase') }}
            </b-dropdown-item>
            <b-dropdown-item
              data-test-id="dropdown-item-delete-sub-pages"
              @click="handleDeletePage('cascade')"
            >
              {{ $t('delete.cascade') }}
            </b-dropdown-item>
          </b-dropdown>
        </template>
      </editor-toolbar>
    </portal>
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
import { fetchID } from 'corteza-webapp-compose/src/lib/tabs.js'

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
      editor: undefined,
      page: undefined,
      blocks: [],
      board: null,
      unsavedBlocks: new Set(),
      id: null,
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

    showDeleteDropdown () {
      return this.hasChildren && this.page.canDeletePage && !this.page.deletedAt
    },

    disableClone () {
      return !!this.module
    },

    cloneTooltip () {
      return this.disableClone ? this.$t('tooltip.saveAsCopy') : ''
    },

  },

  watch: {
    pageID: {
      immediate: true,
      handler (pageID) {
        this.page = undefined
        this.unsavedBlocks.clear()

        if (pageID) {
          const { namespaceID, name } = this.namespace
          this.findPageByID({ namespaceID, pageID, force: true })
            .then(page => {
              document.title = [page.title, name, this.$t('general:label.app-name.private')].filter(v => v).join(' | ')

              this.page = page.clone()
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
    }),

    fulfilEditRequest (blockID) {
      // this ensures whatever changes in tabs is not lost before we lose its configurator
      // because we are reusing that modal component
      this.updateBlocks()
      this.blocks.find((block, i) => fetchID(block) === blockID && this.editBlock(i))
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
      this.editor = { index, block: compose.PageBlockMaker(this.blocks[index]) }
    },

    deleteBlock (index) {
      /**
       * If the block is tabbed, we need to remove it from the tabs when it is deleted.
       * We find where it is tabbed then remove it from its tabs.
       * We don't bother checking if it is actually tabbed because you cannot delete a
       * manually tabed block that is not added to a Tabs block.
       */
      if (this.blocks[index].meta.tabbed) {
        const whereIsItTabbed = this.blocks.filter((block) => block.kind === 'Tabs' && block.options.tabs.some(({ blockID }) => blockID === fetchID(this.blocks[index])))
        whereIsItTabbed.forEach((block) => {
          block.options.tabs = block.options.tabs.filter(({ blockID }) => blockID !== fetchID(this.blocks[index]))
        })
      }

      /**
       * First we get all the tabs that are tabbed in other Tabs blocks.
       * The reduce eliminates duplicates as we only need reps of what  is tabbed
       * We check what needs to be freed by comparing allTabs to the tabs of the block
       * being deleted.
       */
      if (this.blocks[index].kind === 'Tabs') {
        const allTabs = this.blocks.filter((block) => block.kind === 'Tabs' && fetchID(block) !== fetchID(this.blocks[index]))
          .map(({ options }) => options.tabs).flat().reduce((unique, o) => {
            if (!unique.some(tab => tab.blockID === o.blockID)) {
              unique.push(o)
            }
            return unique
          }, []).map(({ blockID }) => blockID)

        const tobefreed = this.blocks[index].options.tabs.filter(({ blockID }) => !allTabs.includes(blockID))
        tobefreed.forEach(({ blockID }) => {
          this.blocks.find((b) => fetchID(b) === blockID).meta.tabbed = false
        })
      }

      this.blocks.splice(index, 1)
      this.page.blocks = this.blocks

      if (this.editor) this.editor = undefined
      this.unsavedBlocks.add(index)
    },

    updatePageBlockGrid (blocks) {
      this.blocks = blocks
    },

    onBlockUpdated (index) {
      this.unsavedBlocks.add(index)
    },

    generateUID () {
      return Math.random().toString(36).substring(2) + (new Date()).getTime().toString(36)
    },

    updateBlocks (block = this.editor.block) {
      block = compose.PageBlockMaker(block)
      this.page.blocks = this.blocks

      if (this.editor.index !== undefined) {
        const oldBlock = this.blocks[this.editor.index]

        if (oldBlock.meta.tabbed === true && this.editor.block.meta.tabbed === false) {
          this.untabBlock(this.editor.block)
        }

        if (this.editor.block.kind !== block.kind) {
          this.page.blocks.push(block)
          const tab = {
            blockID: fetchID(block),
            title: block.title,
          }
          this.$root.$emit('builder-createRequestFulfilled', { tab })
        } else {
          this.page.blocks.splice(this.editor.index, 1, block)
          this.unsavedBlocks.add(this.editor.index)
        }
      } else {
        this.page.blocks.push(block)
        this.unsavedBlocks.add(this.page.blocks.length - 1)
      }

      if (block.kind === 'Tabs') {
        block.options.tabs = block.options.tabs.filter(t => t.blockID !== null)
        block.options.tabs.forEach((tab) => {
          this.page.blocks.find(b => fetchID(b) === tab.blockID).meta.tabbed = true
        })
      }

      if (this.editor.block.kind === block.kind) {
        this.editor = undefined
      }
    },

    cloneBlock (index) {
      this.appendBlock({ ...this.blocks[index].clone() }, this.$t('notification:page.cloneSuccess'))
    },

    appendBlock (block, msg) {
      if (this.blocks.length) {
        // ensuring we append the block to the end of the page
        // eslint-disable-next-line
          const maxY = this.blocks.map((block) => block.xywh[1]).reduce((acc, val) => {
          return acc > val ? acc : val
        }, 0)
        block.xywh = [0, maxY + 2, 3, 3]
      }
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

    async handleSave ({ closeOnSuccess = false, previewOnSuccess = false } = {}) {
      const { namespaceID } = this.namespace

      // Record blocks
      if (this.module && !this.validateModuleFieldSelection(this.module, this.page)) {
        this.toastErrorHandler(this.$t('notification:page.saveFailedRequired'))()
        return
      }

      // Inline record lists
      const queue = []
      this.blocks.forEach((b, index) => {
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

      this.findPageByID({ namespaceID, pageID: this.pageID, force: true })
        .then(page => {
          // Merge changes
          const mergedPage = new compose.Page({ namespaceID, ...page, blocks: this.blocks })

          this.updatePage(mergedPage).then((page) => {
            // get the Tabs Block that still has tabs with tempIDs
            const tabsWithTempIDs = this.page.blocks.filter(b => b.kind === 'Tabs').filter(b => b.options.tabs.some(t => t.blockID.startsWith('tempID-')))

            if (!tabsWithTempIDs.length) return page

            page.blocks.forEach(b => {
              if (b.kind !== 'Tabs') return

              tabsWithTempIDs.find(t => t.meta.tempID === b.meta.tempID).options.tabs.forEach((t, j) => {
                if (!t.blockID.startsWith('tempID-')) return false

                // find a block with the same tempID that should be updated by now and get its blockID
                const updatedBlock = page.blocks.find(block => block.meta.tempID === t.blockID)

                if (!updatedBlock) return false

                const tab = {
                  // fetchID gets the blockID using the found block
                  blockID: fetchID(updatedBlock),
                  title: updatedBlock.title,
                }
                b.options.tabs.splice(j, 1, tab)
              })
            })
            return this.updatePage(page)
          }).then((page) => {
            this.unsavedBlocks.clear()
            this.toastSuccess(this.$t('notification:page.saved'))
            if (closeOnSuccess) {
              this.$router.push({ name: 'admin.pages' })
            } else if (previewOnSuccess) {
              this.$router.push({ name: 'page', params: { pageID: this.pageID } })
            }
            this.page = new compose.Page(page)
          }).catch(this.toastErrorHandler(this.$t('notification:page.saveFailed')))
        })
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

    handleDeletePage (strategy = 'abort') {
      this.deletePage({ ...this.page, strategy }).then(() => {
        this.$router.push({ name: 'admin.pages' })
      }).catch(this.toastErrorHandler(this.$t('notification:page.deleteFailed')))
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
      const parsedBlock = JSON.stringify(this.blocks[index].clone())
      navigator.clipboard.writeText(parsedBlock).then(() => {
        this.toastSuccess(this.$t('notification:page.copySuccess'))
        this.toastInfo(this.$t('notification:page.blockWaiting'))
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
