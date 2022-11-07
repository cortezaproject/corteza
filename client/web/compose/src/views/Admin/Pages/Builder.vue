<template>
  <div
    v-if="!!page"
    class="flex-grow-1 overflow-auto d-flex px-2 w-100"
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
    >
      <template
        slot-scope="{ boundingRect, block, index }"
      >
        <div
          class="h-100 editable-block"
          :class="{ 'bg-warning': !isValid(block) }"
        >
          <div
            class="toolbox border-0 p-2 pr-3 m-0 text-light text-center"
          >
            <b-button
              :title="$t('tooltip.edit.block')"
              variant="link"
              class="p-1 text-light"
              @click="editBlock(index)"
            >
              <font-awesome-icon
                :icon="['far', 'edit']"
              />
            </b-button>

            <c-input-confirm
              class="p-1"
              size="md"
              link
              @confirmed="deleteBlock(index)"
            />
          </div>

          <page-block
            v-bind="{ ...$attrs, ...$props, page, block, boundingRect, blockIndex: index }"
            :record="record"
            :module="module"
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
      @ok="updateBlocks"
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
      :ok-title="$t('label.saveAndClose')"
      ok-variant="primary"
      :cancel-title="$t('label.cancel')"
      cancel-variant="link"
      size="xl"
      :visible="showEditor"
      body-class="p-0 border-top-0"
      header-class="p-3 pb-0 border-bottom-0"
      @ok="updateBlocks"
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
            class="mr-1"
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
import NewBlockSelector from 'corteza-webapp-compose/src/components/Admin/Page/Builder/Selector'
import PageTranslator from 'corteza-webapp-compose/src/components/Admin/Page/PageTranslator'
import Grid from 'corteza-webapp-compose/src/components/Common/Grid'
import PageBlock from 'corteza-webapp-compose/src/components/PageBlocks'
import EditorToolbar from 'corteza-webapp-compose/src/components/Admin/EditorToolbar'
import { compose, NoID } from '@cortezaproject/corteza-js'
import Configurator from 'corteza-webapp-compose/src/components/PageBlocks/Configurator'

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

        if (pageID) {
          const { namespaceID, name } = this.namespace
          this.findPageByID({ namespaceID, pageID: this.pageID, force: true })
            .then(page => {
              document.title = [page.title, name, this.$t('general:label.app-name.private')].filter(v => v).join(' | ')

              this.page = page.clone()
            })
        }
      },
    },
  },

  methods: {
    ...mapActions({
      findPageByID: 'page/findByID',
      updatePage: 'page/update',
      deletePage: 'page/delete',
      updatePageSet: 'page/updateSet',
      createPage: 'page/create',
    }),

    addBlock (block, index = undefined) {
      this.$bvModal.hide('createBlockSelector')
      this.editor = { index, block: compose.PageBlockMaker(block) }
    },

    editBlock (index = undefined) {
      this.editor = { index, block: compose.PageBlockMaker(this.blocks[index]) }
    },

    deleteBlock (index) {
      this.blocks.splice(index, 1)
      this.page.blocks = this.blocks
    },

    updatePageBlockGrid (blocks) {
      this.blocks = blocks
    },

    updateBlocks () {
      const block = compose.PageBlockMaker(this.editor.block)
      this.page.blocks = this.blocks

      if (this.editor.index !== undefined) {
        this.page.blocks.splice(this.editor.index, 1, block)
      } else {
        this.page.blocks.push(block)
      }

      this.editor = undefined
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
            this.$root.$emit(`page-block:validate:${this.page.pageID}-${(this.record || {}).recordID || NoID}-${index}`, resolve)
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

    handleClone () {
      let page = this.page.clone()
      page = {
        ...page,
        pageID: NoID,
        title: this.$t('copyOf', { title: this.page.title }),
        handle: '',
      }

      const { namespaceID = NoID } = this.namespace
      this.createPage({ namespaceID, ...page })
        .then(({ pageID }) => {
          this.$router.push({ name: 'admin.pages.builder', params: { pageID } })
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.cloneFailed')))
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
  z-index: 1000;
  border-top-right-radius: 10px;
  opacity: 0.5;

  &:hover {
    opacity: 1;
  }
}

[dir="rtl"] {
  div.toolbox {
    left: 0;
    right: auto;
  }
}
</style>
