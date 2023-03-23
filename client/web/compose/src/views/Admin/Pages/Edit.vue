<template>
  <div class="py-3">
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
            :icon="['fas', 'cogs']"
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
      <b-row no-gutters>
        <b-col>
          <b-card
            no-body
            class="shadow-sm"
          >
            <b-form
              class="px-4 py-3"
            >
              <b-row>
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
              </b-row>
            </b-form>
          </b-card>
        </b-col>
      </b-row>
    </b-container>

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
        @clone="handleClone()"
        @delete="handleDeletePage"
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
import { compose, NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  name: 'PageEdit',

  components: {
    EditorToolbar,
    PageTranslator,
    Uploader,
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
      page: new compose.Page(),

      showIconModal: false,
      attachments: [],
      selectedAttachmentID: '',
      linkUrl: '',

      processing: false,
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
      return this.pages.some(({ selfID }) => selfID === this.page.pageID)
    },

    disableSave () {
      return [this.titleState, this.handleState].includes(false)
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
        this.page.config.navItem.icon = icon
      },
    },

    isIconSet () {
      return !!this.selectedAttachmentID
    },

    pageIcon () {
      return this.icon.type === 'link' ? this.icon.src : this.makeAttachmentUrl(this.icon.src)
    },
  },

  watch: {
    pageID: {
      immediate: true,
      handler (pageID) {
        if (pageID) {
          this.findPageByID({ namespaceID: this.namespace.namespaceID, pageID }).then((page) => {
            this.page = page.clone()
            return this.fetchAttachments()
          }).catch(this.toastErrorHandler(this.$t('notification:page.loadFailed')))
        }
      },
    },
  },

  methods: {
    ...mapActions({
      findPageByID: 'page/findByID',
      updatePage: 'page/update',
      deletePage: 'page/delete',
      createPage: 'page/create',
      loadPages: 'page/load',
    }),

    async handleSave ({ closeOnSuccess = false } = {}) {
      /**
       * Pass a special tag alongside payload that
       * instructs store layer to add content-language header to the API request
       */
      const resourceTranslationLanguage = this.currentLanguage
      return this.saveIcon().then(icon => {
        this.page.config.navItem.icon = icon
        return this.updatePage({ namespaceID: this.namespace.namespaceID, ...this.page, resourceTranslationLanguage }).then((page) => {
          this.page = page.clone()
          this.toastSuccess(this.$t('notification:page.saved'))
          if (closeOnSuccess) {
            this.$router.push({ name: 'admin.pages' })
          }
        })
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
      this.processing = true

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
        .finally(() => {
          this.processing = false
        })
    },

    async saveIcon () {
      return this.$ComposeAPI.pageUpdateIcon({
        namespaceID: this.namespace.namespaceID,
        pageID: this.pageID,
        type: this.icon.type,
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

  },
}
</script>

<style lang="scss" scoped>
.selected-icon {
  outline: 2px solid $success;
}
</style>
