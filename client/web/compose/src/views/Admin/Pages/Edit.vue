<template>
  <div
    v-if="page"
    class="py-3"
  >
    <portal to="topbar-title">
      {{ $t('edit.edit') }}
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
          :to="{ name: 'admin.pages.builder' }"
        >
          {{ $t('label.pageBuilder') }}
          <font-awesome-icon
            :icon="['far', 'edit']"
            class="ml-2"
          />
        </b-button>

        <page-translator
          v-if="page"
          :page="page"
          style="margin-left:2px;"
        />

        <b-button
          variant="primary"
          :title="$t('tooltip.view')"
          :disabled="!pageViewer"
          :to="pageViewer"
          class="d-flex align-items-center"
          style="margin-left:2px;"
        >
          <font-awesome-icon
            :icon="['far', 'eye']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <b-container fluid="xl">
      <b-card
        no-body
        class="shadow-sm"
      >
        <b-form-row
          v-if="page"
          class="px-4 py-3"
        >
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="`${$t('newPlaceholder')} *`"
              label-class="text-primary"
            >
              <input
                id="id"
                v-model="page.pageID"
                required
                type="hidden"
              >
              <b-form-input
                v-model="page.title"
                data-test-id="input-title"
                required
                :state="titleState"
                class="mb-2"
              />
            </b-form-group>
          </b-col>
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('label.handle')"
              label-class="text-primary"
            >
              <b-form-input
                v-model="page.handle"
                data-test-id="input-handle"
                :state="handleState"
                class="mb-2"
                :placeholder="$t('block.general.placeholder.handle')"
              />
              <b-form-invalid-feedback :state="handleState">
                {{ $t('block.general.invalid-handle-characters') }}
              </b-form-invalid-feedback>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
          >
            <b-form-group
              :label="$t('label.description')"
              label-class="text-primary"
            >
              <b-form-textarea
                v-model="page.description"
                data-test-id="input-description"
                :placeholder="$t('edit.pageDescription')"
                rows="4"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              label-class="d-flex align-items-center text-primary"
            >
              <template #label>
                {{ $t('icon.page') }}
                <b-button
                  :title="$t('icon.configure')"
                  variant="outline-light"
                  class="d-flex align-items-center px-1 text-primary border-0 ml-1"
                  @click="openIconModal"
                >
                  <font-awesome-icon
                    :icon="['fas', 'cog']"
                  />
                </b-button>
              </template>

              <img
                v-if="icon.src"
                :src="pageIcon"
                width="auto"
                height="50"
              >

              <span v-else>
                {{ $t('icon.noIcon') }}
              </span>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('edit.otherOptions')"
              label-class="text-primary"
            >
              <b-form-checkbox
                v-if="!isRecordPage"
                v-model="page.visible"
                data-test-id="checkbox-page-visibility"
              >
                {{ $t('edit.visible') }}
              </b-form-checkbox>

              <b-form-checkbox
                v-model="page.config.navItem.expanded"
                data-test-id="checkbox-show-sub-pages-in-sidebar"
              >
                {{ $t('showSubPages') }}
              </b-form-checkbox>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
          >
            <hr>
            <b-form-group
              label="Layouts"
              label-class="text-primary"
            >
              <b-table-simple
                responsive="lg"
                borderless
                small
              >
                <b-thead>
                  <tr>
                    <th />

                    <th
                      class="text-primary"
                      style="width: 45%; min-width: 200px;"
                    >
                      Title *
                    </th>

                    <th
                      class="text-primary"
                      style="width: 45%; min-width: 200px;"
                    >
                      Handle
                    </th>

                    <th style="width: 80px;" />
                  </tr>
                </b-thead>

                <draggable
                  v-model="layouts"
                  handle=".handle"
                  tag="b-tbody"
                >
                  <tr
                    v-for="(layout, index) in layouts"
                    :key="index"
                  >
                    <b-td class="handle text-center align-middle pr-2">
                      <font-awesome-icon
                        :icon="['fas', 'bars']"
                        class="grab m-0 text-light p-0"
                      />
                    </b-td>

                    <b-td
                      class="align-middle"
                    >
                      <b-form-input
                        v-model="layout.meta.title"
                        :state="layoutTitleState(layout.meta.title)"
                        @input="layout.meta.updated = true"
                      />
                    </b-td>

                    <b-td
                      class="align-middle"
                    >
                      <b-input-group>
                        <b-form-input
                          v-model="layout.handle"
                          :state="layoutHandleState(layout.handle)"
                          @input="layout.meta.updated = true"
                        />

                        <b-input-group-append>
                          <b-button
                            variant="light"
                            class="d-flex align-items-center px-3"
                            @click="configureLayout(index)"
                          >
                            <font-awesome-icon
                              :icon="['fas', 'wrench']"
                            />
                          </b-button>

                          <b-button
                            variant="primary"
                            :disabled="layout.pageLayoutID === '0'"
                            class="d-flex align-items-center"
                            :to="{ name: 'admin.pages.builder', query: { layoutID: layout.pageLayoutID} }"
                          >
                            <font-awesome-icon
                              :icon="['far', 'edit']"
                            />
                          </b-button>
                        </b-input-group-append>
                      </b-input-group>
                    </b-td>

                    <td
                      class="text-center align-middle"
                    >
                      <c-input-confirm
                        :title="$t('tabs.tooltip.delete')"
                        class="ml-2"
                        @confirmed="deleteLayout(index)"
                      />
                    </td>
                  </tr>
                </draggable>
              </b-table-simple>
              <b-button
                variant="primary"
                @click="addLayout"
              >
                Add layout
              </b-button>
            </b-form-group>
          </b-col>
        </b-form-row>
      </b-card>
    </b-container>

    <b-modal
      v-if="layoutEditor.layout"
      :visible="!!layoutEditor.layout"
      title="Configure layout"
      :ok-title="$t('general:label.saveAndClose')"
      ok-variant="primary"
      cancel-variant="link"
      size="lg"
      @ok="updateLayout()"
      @cancel="layoutEditor.layout = undefined"
      @hide="layoutEditor.layout = undefined"
    >
      <b-form-group
        label="Condition"
        label-class="text-primary"
      >
        <b-input-group>
          <b-input-group-prepend>
            <b-button variant="dark">
              ƒ
            </b-button>
          </b-input-group-prepend>
          <b-form-input
            v-model="layoutEditor.layout.config.visibility.expression"
            placeholder="When will the layout be shown"
          />
        </b-input-group>
      </b-form-group>

      <b-form-group
        label="Roles"
        label-class="text-primary"
      >
        <vue-select
          v-model="currentLayoutRoles"
          :options="roles.options"
          :loading="roles.processing"
          placeholder="Pick roles that the layout will be shown to"
          :get-option-label="role => role.name"
          :reduce="role => role.roleID"
          :selectable="role => !currentLayoutRoles.includes(role.roleID)"
          append-to-body
          multiple
          class="bg-white"
        />
      </b-form-group>
    </b-modal>

    <b-modal
      v-model="showIconModal"
      :title="$t('icon.configure')"
      :ok-title="$t('label.saveAndClose')"
      size="lg"
      label-class="text-primary"
      cancel-variant="link"
      @close="closeIconModal"
      @ok="saveIconModal"
    >
      <b-form-group
        :label="$t('icon.upload')"
        label-class="text-primary"
        class="mb-0"
      >
        <uploader
          :endpoint="endpoint"
          :accepted-files="['image/*']"
          :param-name="'icon'"
          @uploaded="uploadAttachment"
        />

        <b-form-group
          :label="$t('url.label')"
          label-class="text-primary"
          class="my-2"
        >
          <b-input-group>
            <b-input
              v-model="linkUrl"
              :disabled="isIconSet"
            />
            <b-input-group-append>
              <b-button
                v-b-modal.logo
                :title="$t('tooltip.preview-link')"
                :disabled="!linkUrl"
                variant="light"
                rounded
                class="d-flex align-items-center btn-light"
              >
                <font-awesome-icon :icon="['fas', 'external-link-alt']" />
              </b-button>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>
      </b-form-group>

      <template v-if="attachments.length > 0">
        <hr>

        <b-form-group
          :label="$t('icon.list')"
          label-class="text-primary"
        >
          <div
            v-if="processing"
            class="d-flex align-items-center justify-content-center h-100"
          >
            <b-spinner />
          </div>

          <div
            v-else
            class="d-flex flex-wrap"
          >
            <img
              v-for="a in attachments"
              :key="a.attachmentID"
              :src="a.src"
              :alt="a.name"
              width="auto"
              height="50"
              :class="{ 'selected-icon': selectedAttachmentID === a.attachmentID }"
              class="rounded pointer mr-2"
              @click="toggleSelectedIcon(a.attachmentID)"
            >
          </div>
        </b-form-group>
      </template>
    </b-modal>

    <b-modal
      id="logo"
      hide-header
      hide-footer
      centered
      body-class="p-1"
    >
      <b-img
        :src="linkUrl"
        fluid-grow
      />
    </b-modal>

    <portal to="admin-toolbar">
      <editor-toolbar
        :back-link="{ name: 'admin.pages' }"
        :hide-delete="hideDelete"
        :hide-save="!page.canUpdatePage"
        :disable-save="disableSave"
        :processing="processing"
        @clone="handleClone()"
        @delete="handleDeletePage()"
        @save="handleSave()"
        @saveAndClose="handleSave({ closeOnSuccess: true })"
      >
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
import EditorToolbar from 'corteza-webapp-compose/src/components/Admin/EditorToolbar'
import PageTranslator from 'corteza-webapp-compose/src/components/Admin/Page/PageTranslator'
import pages from 'corteza-webapp-compose/src/mixins/pages'
import Uploader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Uploader'
import Draggable from 'vuedraggable'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'
import { VueSelect } from 'vue-select'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  name: 'PageEdit',

  components: {
    EditorToolbar,
    PageTranslator,
    Uploader,
    Draggable,
    VueSelect,
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

      page: new compose.Page(),

      showIconModal: false,
      attachments: [],
      selectedAttachmentID: '',
      linkUrl: '',

      layouts: [],

      layoutEditor: {
        index: undefined,
        layout: undefined,
      },

      deletedLayouts: new Set(),

      roles: {
        processing: false,
        options: [],
      },
    }
  },

  computed: {
    ...mapGetters({
      pages: 'page/set',
    }),

    titleState () {
      return this.page.title.length > 0 ? null : false
    },

    handleState () {
      return handle.handleState(this.page.handle)
    },

    pageViewer () {
      if (this.isRecordPage) {
        return undefined
      }
      const { pageID } = this.page
      return { name: 'page', params: { pageID } }
    },

    isRecordPage () {
      return this.page && this.page.moduleID !== NoID
    },

    hasChildren () {
      return this.page ? this.pages.some(({ selfID }) => selfID === this.page.pageID) : false
    },

    disableSave () {
      return [this.titleState, this.handleState].includes(false) || this.layouts.some(l => !l.meta.title || handle.handleState(l.handle) === false)
    },

    hideDelete () {
      return this.hasChildren || !this.page.canDeletePage || !!this.page.deletedAt
    },

    showDeleteDropdown () {
      return this.hasChildren && this.page.canDeletePage && !this.page.deletedAt
    },

    endpoint () {
      return this.$ComposeAPI.iconUploadEndpoint({
        namespaceID: this.namespace.namespaceID,
      })
    },

    icon: {
      get () {
        return this.page.config.navItem.icon || {}
      },

      set (icon) {
        this.$set(this.page.config.navItem, 'icon', icon)
      },
    },

    isIconSet () {
      return !!this.selectedAttachmentID
    },

    pageIcon () {
      if (!this.icon.src) {
        return
      }

      return this.icon.type === 'link' ? this.icon.src : this.makeAttachmentUrl(this.icon.src)
    },

    currentLayoutRoles: {
      get () {
        if (!this.layoutEditor.layout) {
          return []
        }

        return this.layoutEditor.layout.config.visibility.roles
      },

      set (roles) {
        this.$set(this.layoutEditor.layout.config.visibility, 'roles', roles)
      },
    },
  },

  watch: {
    pageID: {
      immediate: true,
      handler (pageID) {
        this.page = undefined
        this.layouts = []

        this.deletedLayouts = new Set()

        if (pageID) {
          this.processing = true

          const { namespaceID } = this.namespace
          this.findPageByID({ namespaceID, pageID }).then((page) => {
            this.page = page.clone()
            return this.fetchAttachments()
          }).then(this.fetchLayouts)
            .finally(() => {
              this.processing = false
            }).catch(this.toastErrorHandler(this.$t('notification:page.loadFailed')))
        }
      },
    },
  },

  created () {
    this.fetchRoles()
  },

  methods: {
    ...mapActions({
      findPageByID: 'page/findByID',
      updatePage: 'page/update',
      deletePage: 'page/delete',
      createPage: 'page/create',
      loadPages: 'page/load',
      findLayoutsByPageID: 'pageLayout/findByPageID',
      createPageLayout: 'pageLayout/create',
      updatePageLayout: 'pageLayout/update',
      deletePageLayout: 'pageLayout/delete',
    }),

    async fetchLayouts () {
      const { namespaceID } = this.namespace
      return this.findLayoutsByPageID({ namespaceID, pageID: this.pageID }).then(layouts => {
        this.layouts = layouts.map(layout => new compose.PageLayout(layout))
      })
    },

    async fetchRoles () {
      this.roles.processing = true

      this.$SystemAPI.roleList().then(({ set: roles = [] }) => {
        this.roles.options = roles.filter(({ meta }) => !(meta.context && meta.context.resourceTypes))
      }).finally(() => {
        this.roles.processing = false
      })
    },

    addLayout () {
      this.layouts.push(new compose.PageLayout({ namespaceID: this.namespace.namespaceID, pageID: this.pageID }))
    },

    updateLayout () {
      this.layoutEditor.layout.meta.updated = true
      this.layouts.splice(this.layoutEditor.index, 1, this.layoutEditor.layout)
      this.layoutEditor.index = undefined
      this.layoutEditor.layout = undefined
    },

    deleteLayout (index) {
      const { pageLayoutID } = this.layouts[index] || {}
      if (pageLayoutID !== NoID) {
        this.deletedLayouts.add(this.layouts[index])
      }

      this.layouts.splice(index, 1)
    },

    configureLayout (index) {
      this.layoutEditor.index = index
      this.layoutEditor.layout = { ...this.layouts[index] }
    },

    async handleSaveLayouts () {
      // Delete first so old deleted handles don't interfere with new identical ones
      return Promise.all([...this.deletedLayouts].map(this.deletePageLayout)).then(() => {
        return Promise.all(this.layouts.map(layout => {
          if (layout.pageLayoutID === NoID) {
            return this.createPageLayout(layout)
          } else if (layout.meta.updated) {
            return this.updatePageLayout(layout)
          }
        }))
      })
    },

    async handlePageLayoutReorder () {
      const { namespaceID } = this.namespace
      const pageIDs = this.layouts.map(({ pageLayoutID }) => pageLayoutID)

      return this.$ComposeAPI.pageLayoutReorder({ namespaceID, pageID: this.pageID, pageIDs }).then(() => {
        return this.$store.dispatch('pageLayout/load', { namespaceID, clear: true, force: true })
      })
    },

    handleSave ({ closeOnSuccess = false } = {}) {
      this.processing = true

      /**
       * Pass a special tag alongside payload that
       * instructs store layer to add content-language header to the API request
       */
      const resourceTranslationLanguage = this.currentLanguage
      const { namespaceID } = this.namespace

      return this.saveIcon().then(icon => {
        this.page.config.navItem.icon = icon
        return this.updatePage({ namespaceID, ...this.page, resourceTranslationLanguage })
      }).then(page => {
        this.page = page.clone()
        return this.handleSaveLayouts()
      }).then(this.handlePageLayoutReorder)
        .then(() => {
          this.fetchLayouts()
          this.deletedLayouts = new Set()

          this.toastSuccess(this.$t('notification:page.saved'))
          if (closeOnSuccess) {
            this.$router.push({ name: 'admin.pages' })
          }
        }).finally(() => {
          this.processing = false
        }).catch(this.toastErrorHandler(this.$t('notification:page.saveFailed')))
    },

    handleDeletePage (strategy = 'abort') {
      this.deletePage({ ...this.page, strategy }).then(() => {
        this.$router.push({ name: 'admin.pages' })
      }).catch(this.toastErrorHandler(this.$t('notification:page.deleteFailed')))
    },

    uploadAttachment () {
      this.fetchAttachments()
    },

    async fetchAttachments () {
      return this.$ComposeAPI.iconList({ sort: 'id DESC' })
        .then(({ set: attachments = [] }) => {
          const baseURL = this.$ComposeAPI.baseURL
          this.attachments = []

          if (attachments) {
            attachments.forEach(a => {
              const src = !a.url.includes(baseURL) ? this.makeAttachmentUrl(a.url) : a.url
              this.attachments.push({ ...a, src })
            })
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.iconFetchFailed')))
    },

    async saveIcon () {
      return this.$ComposeAPI.pageUpdateIcon({
        namespaceID: this.namespace.namespaceID,
        pageID: this.pageID,
        type: this.icon.type || 'link',
        source: this.icon.src,
      })
    },

    toggleSelectedIcon (attachmentID = '') {
      this.selectedAttachmentID = this.selectedAttachmentID === attachmentID ? '' : attachmentID
    },

    openIconModal () {
      this.linkUrl = this.icon.type === 'link' ? this.icon.src : ''
      this.selectedAttachmentID = (this.attachments.find(a => a.url === this.icon.src) || {}).attachmentID
      this.showIconModal = true
    },

    saveIconModal () {
      const type = this.selectedAttachmentID ? 'attachment' : 'link'

      let src = this.linkUrl
      if (this.selectedAttachmentID) {
        src = (this.attachments.find(({ attachmentID }) => attachmentID === this.selectedAttachmentID) || {}).url
      }

      this.icon = { type, src }
    },

    closeIconModal () {
      this.linkUrl = this.icon.type === 'link' ? this.icon.src : ''
      this.selectedAttachmentID = (this.attachments.find(a => a.url === this.icon.src) || {}).attachmentID
    },

    makeAttachmentUrl (src) {
      return `${this.$ComposeAPI.baseURL}${src}`
    },

    layoutTitleState (title) {
      return title ? null : false
    },

    layoutHandleState (layoutHandle) {
      return handle.handleState(layoutHandle)
    },
  },
}
</script>

<style lang="scss" scoped>
.selected-icon {
  outline: 2px solid $success;
}
</style>
