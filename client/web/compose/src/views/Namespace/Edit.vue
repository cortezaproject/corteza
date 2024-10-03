<template>
  <div
    v-if="namespace"
    class="d-flex flex-column w-100 h-100"
  >
    <portal to="topbar-title">
      {{ pageTitle }}
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="isEdit"
        size="sm"
      >
        <b-button
          data-test-id="button-visit-namespace"
          variant="primary"
          class="d-flex align-items-center"
          :to="openNamespace"
          :disabled="!namespaceEnabled"
        >
          {{ $t('visit') }}
          <font-awesome-icon
            :icon="['far', 'eye']"
            class="ml-2"
          />
        </b-button>
        <b-button
          v-if="namespace.canManageNamespace"
          v-b-tooltip.noninteractive.hover="{ title: $t('configure'), container: '#body' }"
          data-test-id="button-visit-admin-panel"
          variant="primary"
          class="d-flex align-items-center"
          :to="{ name: 'admin.modules', params: { slug: namespace.slug } }"
          style="margin-left:2px;"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
          />
        </b-button>
        <namespace-translator
          v-if="namespace"
          :namespace="namespace"
          :disabled="isNew"
          button-variant="primary"
          style="margin-left:2px;"
        />
      </b-button-group>
    </portal>

    <div class="flex-grow-1 overflow-auto py-3">
      <b-container
        fluid="xl"
        class="flex-grow-1"
      >
        <b-card
          header-class="d-flex align-items-center gap-1 border-bottom"
          body-class="p-3"
          footer-bg-variant="warning"
        >
          <template
            v-if="isEdit"
            #header
          >
            <b-btn
              v-if="namespace.canExportNamespace"
              data-test-id="button-export-namespace"
              variant="light"
              size="lg"
              @click="exportNamespace"
            >
              {{ $t('export') }}
            </b-btn>

            <c-permissions-button
              v-if="namespace.canGrant"
              data-test-id="button-permissions"
              :title="namespace.name || namespace.slug || namespace.namespaceID"
              :target="namespace.name || namespace.slug || namespace.namespaceID"
              :resource="`corteza::compose:namespace/${namespace.namespaceID}`"
              :button-label="$t('label.permissions')"
              button-variant="light"
              class="btn-lg"
            />
          </template>

          <b-form>
            <b-form-group
              :label="$t('name.label')"
              label-class="text-primary"
            >
              <b-input-group>
                <b-form-input
                  id="ns-nm"
                  v-model="namespace.name"
                  data-test-id="input-name"
                  type="text"
                  required
                  :state="nameState"
                  :placeholder="$t('name.placeholder')"
                />
                <b-input-group-append>
                  <namespace-translator
                    :namespace="namespace"
                    highlight-key="name"
                    :disabled="isNew"
                  />
                </b-input-group-append>
              </b-input-group>
            </b-form-group>
            <b-form-group
              :label="$t('slug.label')"
              :description="$t('slug.description')"
              label-class="text-primary"
            >
              <b-form-input
                v-model="namespace.slug"
                data-test-id="input-slug"
                type="text"
                required
                :state="slugState"
                :placeholder="$t('slug.placeholder')"
              />
              <b-form-invalid-feedback :state="slugState">
                {{ $t('slug.invalid-handle-characters') }}
              </b-form-invalid-feedback>
            </b-form-group>
            <b-form-group>
              <b-form-checkbox
                v-model="namespace.enabled"
                data-test-id="checkbox-enable-namespace"
                class="mb-3"
              >
                {{ $t('enabled.label') }}
              </b-form-checkbox>
              <b-form-checkbox
                v-model="isApplication"
                data-test-id="checkbox-toggle-application"
                :disabled="!canToggleApplication"
              >
                {{ $t('application.label') }}
              </b-form-checkbox>
            </b-form-group>
            <hr>

            <b-form-group>
              <b-form-checkbox
                v-model="namespace.meta.logoEnabled"
                data-test-id="checkbox-show-logo"
              >
                {{ $t('logo.show') }}
              </b-form-checkbox>
            </b-form-group>

            <b-form-group
              v-if="namespace.meta.logoEnabled && isEdit"
            >
              <template #label>
                <div class="d-flex align-items-center">
                  {{ $t('logo.label') }}
                  <b-button
                    v-if="logoPreview"
                    v-b-modal.logo
                    data-test-id="button-logo-preview"
                    variant="link"
                    size="sm"
                    class="d-flex align-items-center border-0 p-0 ml-2"
                  >
                    <font-awesome-icon
                      :icon="['far', 'eye']"
                    />
                  </b-button>

                  <b-button
                    v-if="!!namespace.meta.logo"
                    data-test-id="button-logo-reset"
                    variant="light"
                    size="sm"
                    class="py-0 ml-2"
                    @click="resetLogo()"
                  >
                    {{ $t('logo.reset') }}
                  </b-button>
                </div>
              </template>

              <b-form-file
                v-model="namespaceAssets.logo"
                data-test-id="file-logo-upload"
                accept="image/*"
                :placeholder="$t('logo.placeholder')"
              />
            </b-form-group>

            <!-- <b-form-group>
              <template #label>
                <div class="d-flex align-items-center">
                  {{ $t('icon.label') }}
                  <b-button
                    v-if="namespace.meta.icon"
                    variant="primary"
                    class="py-0 ml-2"
                    v-b-modal.icon
                  >
                    Preview
                  </b-button>
                </div>
              </template>
              <b-form-file
                v-model="namespaceAssets.icon"
                accept="image/*"
                :placeholder="$t('icon.placeholder')"
              />
            </b-form-group> -->

            <b-form-group
              :label="$t('subtitle.label')"
              label-class="text-primary"
            >
              <b-input-group>
                <b-form-input
                  v-model="namespace.meta.subtitle"
                  data-test-id="input-subtitle"
                  type="text"
                  :placeholder="$t('subtitle.placeholder')"
                />
                <b-input-group-append>
                  <namespace-translator
                    :namespace="namespace"
                    highlight-key="meta.subtitle"
                    :disabled="isNew"
                  />
                </b-input-group-append>
              </b-input-group>
            </b-form-group>

            <b-form-group
              :label="$t('description.label')"
              label-class="text-primary"
              class="mb-3"
            >
              <b-input-group>
                <b-form-textarea
                  v-model="namespace.meta.description"
                  data-test-id="input-description"
                  rows="1"
                  :placeholder="$t('description.placeholder')"
                />
                <b-input-group-append>
                  <namespace-translator
                    :namespace="namespace"
                    highlight-key="meta.description"
                    :disabled="isNew"
                  />
                </b-input-group-append>
              </b-input-group>
            </b-form-group>
            <hr>

            <b-form-group
              :label="$t('sidebar.configure')"
              label-class="text-primary"
            >
              <b-form-checkbox
                v-model="namespace.meta.hideSidebar"
                data-test-id="checkbox-show-sidebar"
              >
                {{ $t('sidebar.hide') }}
              </b-form-checkbox>
            </b-form-group>
          </b-form>

          <template
            v-if="isClone"
            #footer
          >
            {{ $t('cloneWarning.wfInclusion') }}
          </template>
        </b-card>
      </b-container>
    </div>

    <editor-toolbar
      :processing="processing"
      :processing-save="processingSave"
      :processing-save-and-close="processingSaveAndClose"
      :processing-clone="processingClone"
      :processing-delete="processingDelete"
      :hide-delete="hideDelete"
      :hide-clone="!isEdit"
      :hide-save="hideSave"
      :disable-save="disableSave"
      @back="$router.go(-1)"
      @delete="handleDelete"
      @save="handleSave()"
      @clone="handleClone()"
      @saveAndClose="handleSave({ closeOnSuccess: true })"
    />

    <b-modal
      id="logo"
      hide-header
      hide-footer
      centered
      no-fade
      body-class="p-1"
    >
      <b-img
        v-if="logoPreview"
        :src="logoPreview"
        fluid-grow
      />
    </b-modal>

    <b-modal
      id="icon"
      hide-header
      hide-footer
      centered
      no-fade
      body-class="p-1"
    >
      <b-img
        :src="iconPreview"
        fluid-grow
      />
    </b-modal>
  </div>
</template>

<script>
import { isEqual } from 'lodash'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { url, handle } from '@cortezaproject/corteza-vue'
import EditorToolbar from 'corteza-webapp-compose/src/components/Admin/EditorToolbar'
import NamespaceTranslator from 'corteza-webapp-compose/src/components/Namespaces/NamespaceTranslator'
import { mapGetters, mapActions } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'namespace',
  },

  name: 'EditNamespace',

  components: {
    EditorToolbar,
    NamespaceTranslator,
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedNamespace(next)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedNamespace(next)
  },

  props: {
    isClone: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      processing: false,
      processingSave: false,
      processingSaveAndClose: false,
      processingClone: false,
      processingDelete: false,

      namespace: undefined,
      initialNamespaceState: undefined,

      namespaceAssets: {
        logo: undefined,
        icon: undefined,
      },

      namespaceAssetsInitialState: {
        logo: undefined,
        icon: undefined,
      },

      namespaceEnabled: false,

      application: undefined,
      isApplication: false,
      isApplicationInitialState: false,
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
      previousPage: 'ui/previousPage',
    }),

    canCreateApplication () {
      return this.can('system/', 'application.create')
    },

    isNew () {
      return this.namespace.namespaceID === NoID
    },

    pageTitle () {
      switch (true) {
        case this.isEdit:
          return this.$t('edit')
        case this.isClone:
          return this.$t('clone')

        default:
          return this.$t('create')
      }
    },

    watchKey () {
      return `${this.$route.params.namespaceID}|${this.$route.name}`
    },

    openNamespace () {
      return { name: 'pages', params: { slug: (this.namespace.slug || this.namespace.namespaceID) } }
    },

    isEdit () {
      return this.$route.name === 'namespace.edit' || this.isClone
    },

    logoPreview () {
      return this.namespace.meta.logo || this.$Settings.attachment('ui.mainLogo')
    },

    iconPreview () {
      return this.namespace.meta.icon || ''
    },

    nameState () {
      return this.namespace.name.length > 0 ? null : false
    },

    slugState () {
      return handle.handleState(this.namespace.slug)
    },

    canToggleApplication () {
      return this.canCreateApplication
    },

    disableSave () {
      return [this.nameState, this.slugState].includes(false)
    },

    hideSave () {
      return this.isEdit && !this.namespace.canUpdateNamespace
    },

    hideDelete () {
      return !this.isEdit || !!this.namespace.deletedAt || (this.isEdit && !this.namespace.canDeleteNamespace)
    },
  },

  watch: {
    watchKey: {
      immediate: true,
      handler () {
        this.fetchNamespace()
      },
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      updateNamespace: 'namespace/update',
      createNamespace: 'namespace/create',
      findNamespace: 'namespace/findByID',
      cloneNamespace: 'namespace/clone',
      deleteNamespace: 'namespace/delete',
    }),

    async fetchNamespace () {
      this.processing = true

      const namespaceID = this.$route.params.namespaceID

      this.namespace = undefined
      this.initialNamespaceState = undefined
      this.application = undefined
      this.isApplication = false
      this.isApplicationInitialState = this.isApplication

      if (namespaceID) {
        await this.findNamespace({ namespaceID })
          .then(ns => {
            this.namespaceEnabled = ns.enabled
            this.namespace = new compose.Namespace(ns)

            this.fetchApplication()
          })
      } else {
        this.namespace = new compose.Namespace({ enabled: true })
      }

      this.namespace.meta = {
        subtitle: '',
        description: '',
        hideSidebar: false,
        logoEnabled: null,
        ...this.namespace.meta,
      }

      this.initialNamespaceState = this.namespace.clone()
      this.namespaceAssetsInitialState = this.namespaceAssets

      this.processing = false
    },

    exportNamespace () {
      const params = {
        namespaceID: this.namespace.namespaceID,
        filename: encodeURIComponent(this.namespace.name.replace(/\./g, '-')),
      }

      const exportUrl = url.Make({
        url: `${this.$ComposeAPI.baseURL}${this.$ComposeAPI.namespaceExportEndpoint(params)}`,
        query: {
          jwt: this.$auth.accessToken,
        },
      })

      window.open(exportUrl)
    },

    fetchApplication () {
      const { namespaceID, slug } = this.namespace
      this.$SystemAPI.applicationList({ name: slug || namespaceID })
        .then(({ set = [] }) => {
          if (set.length) {
            this.application = set[0]
            this.isApplication = this.application.enabled
            this.isApplicationInitialState = this.isApplication
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:namespace.application.fetchFailed')))
    },

    async handleSave ({ closeOnSuccess = false } = {}) {
      const toggleProcessing = () => {
        this.processing = !this.processing

        if (closeOnSuccess) {
          this.processingSaveAndClose = !this.processingSaveAndClose
        } else {
          this.processingSave = !this.processingSave
        }
      }

      toggleProcessing()

      /**
       * Pass a special tag alongside payload that
       * instructs store layer to add content-language header to the API request
       */
      const resourceTranslationLanguage = this.currentLanguage
      let { namespaceID, name, slug, enabled, meta } = this.namespace
      let assets

      // Firstly handle any new namespace assets
      if (this.namespaceAssets.logo || this.namespaceAssets.icon) {
        try {
          assets = await this.uploadAssets()
          meta = { ...meta, ...assets }
          this.namespaceAssetsInitialState = this.namespaceAssets
        } catch (e) {
          this.toastErrorHandler(this.$t('notification:namespace.assetUploadFailed'))(e)
          toggleProcessing()
          return
        }
      }

      const payload = {
        name,
        slug,
        enabled,
        meta,
        resourceTranslationLanguage,
      }

      if (this.isEdit) {
        try {
          await this.updateNamespace({ ...payload, namespaceID }).then((ns) => {
            this.namespaceEnabled = ns.enabled
            this.namespace = new compose.Namespace(ns)

            this.toastSuccess(this.$t('notification:namespace.saved'))
          })
        } catch (e) {
          this.toastErrorHandler(this.$t('notification:namespace.saveFailed'))(e)
          toggleProcessing()
          return
        }
      } else {
        try {
          await this.createNamespace(payload).then((ns) => {
            this.namespaceEnabled = ns.enabled
            this.namespace = new compose.Namespace(ns)

            this.toastSuccess(this.$t('notification:namespace.saved'))
          })
        } catch (e) {
          this.toastErrorHandler(this.$t('notification:namespace.createFailed'))(e)
          toggleProcessing()
          return
        }
      }

      await this.handleApplicationSave()
        .catch(() => this.toastErrorHandler(this.$t('notification:namespace.createAppFailed')))

      this.initialNamespaceState = this.namespace.clone()
      this.isApplicationInitialState = this.isApplication

      toggleProcessing()

      if (closeOnSuccess) {
        this.$router.push(this.previousPage || { name: 'namespace.manage' })
      } else if (!this.isEdit || this.isClone) {
        this.$router.push({ name: 'namespace.edit', params: { namespaceID: this.namespace.namespaceID } })
      }
    },

    handleClone () {
      this.processingClone = true

      let { name, slug } = this.namespace

      name = `${name} (${this.$t('cloneSuffix')})`
      slug = slug ? `${slug}_${this.$t('cloneSuffix')}` : ''

      return this.cloneNamespace({ ...this.namespace, name, slug }).then(({ namespaceID }) => {
        this.$route.params.namespaceID = namespaceID
        this.toastSuccess(this.$t('notification:namespace.cloned'))
        this.$router.push({ name: 'namespace.edit', params: { namespaceID, isClone: true } })
      }).catch(e => {
        this.toastErrorHandler(this.$t('notification:namespace.cloneFailed'))(e)
      }).finally(() => {
        this.processingClone = false
      })
    },

    handleDelete () {
      this.processing = true
      this.processingDelete = true

      const { namespaceID } = this.namespace
      const { applicationID } = this.application || {}

      this.deleteNamespace({ namespaceID })
        .catch(this.toastErrorHandler(this.$t('notification:namespace.deleteFailed')))
        .then(() => {
          this.namespace.deletedAt = new Date()

          if (applicationID) {
            return this.$SystemAPI.applicationDelete({ applicationID })
          }
        })
        .then(() => {
          this.$router.push({ name: 'namespace.manage' })
          this.toastSuccess(this.$t('notification:namespace.deleted'))
        })
        .finally(() => {
          this.processing = false
          this.processingDelete = false
        })
    },

    async handleApplicationSave () {
      if (this.application) {
        this.application.name = this.namespace.slug || this.namespace.namespaceID
        this.application.unify.name = this.namespace.name
        this.application.unify.url = `/compose/ns/${this.application.name}/pages`

        let enabled = this.application.enabled
        if (this.isApplication && !this.application.enabled) {
          enabled = true
        } else if (!this.isApplication && this.application.enabled) {
          enabled = false
        }

        this.application.unify.listed = enabled

        // Assets
        // Don't take note of the ID, it will be different on the system side
        this.application.unify.icon = this.application.unify.icon || this.namespace.meta.icon
        this.application.unify.logo = this.application.unify.logo || this.namespace.meta.logo

        return this.$SystemAPI.applicationUpdate({ ...this.application, enabled })
          .then(app => {
            this.application = app
            this.isApplication = this.application.enabled
          })
          .catch(this.toastErrorHandler(this.$t('notification:namespace.application.saveFailed')))
      } else if (this.isApplication) {
        // If namespace not an application - create one and enable
        const application = {
          name: this.namespace.slug || this.namespace.namespaceID,
          enabled: true,
          unify: {
            name: this.namespace.name,
            listed: true,
            url: `compose/ns/${this.namespace.slug || this.namespace.namespaceID}/pages`,
            icon: this.namespace.meta.icon || this.$Settings.attachment('ui.iconLogo'),
            logo: this.namespace.meta.logo || this.$Settings.attachment('ui.mainLogo'),
          },
        }
        return this.$SystemAPI.applicationCreate({ ...application })
          .then(app => {
            this.application = app
            this.isApplication = this.application.enabled
          })
          .catch(this.toastErrorHandler(this.$t('notification:namespace.application.createFailed')))
      }
    },

    async uploadAssets () {
      const rr = {}

      const rq = async (file) => {
        const formData = new FormData()
        formData.append('upload', file)

        const rsp = await this.$ComposeAPI.api().request({
          method: 'post',
          url: this.$ComposeAPI.namespaceUploadEndpoint(),
          data: formData,
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })
        if (rsp.data.error) {
          throw new Error(rsp.data.error.message)
        }
        return rsp.data.response
      }

      const baseURL = this.$ComposeAPI.baseURL

      if (this.namespaceAssets.logo) {
        const rsp = await rq(this.namespaceAssets.logo)
        rr.logo = baseURL + rsp.url
        rr.logoID = rsp.attachmentID

        this.namespaceAssets.logo = undefined
      }

      if (this.namespaceAssets.icon) {
        const rsp = await rq(this.namespaceAssets.icon)
        rr.icon = baseURL + rsp.url
        rr.iconID = rsp.attachmentID

        this.namespaceAssets.icon = undefined
      }

      return rr
    },

    resetLogo () {
      this.namespace.meta.logo = undefined
      this.namespace.meta.logoID = undefined
    },

    checkUnsavedNamespace (next) {
      if (this.isNew || this.processingClone || this.namespace.deletedAt) {
        return next(true)
      }

      const namespaceState = !isEqual(this.namespace.clone(), this.initialNamespaceState.clone())
      const isApplicationState = !(this.isApplication === this.isApplicationInitialState)
      const namespaceAssetsState = !isEqual(this.namespaceAssets, this.namespaceAssetsInitialState)

      return next((namespaceState || isApplicationState || namespaceAssetsState) ? window.confirm(this.$t('manage.unsavedChanges')) : true)
    },

    setDefaultValues () {
      this.processing = false
      this.processingSaveAndClose = false
      this.processingSave = false
      this.processingClone = false
      this.namespace = undefined
      this.initialNamespaceState = undefined
      this.namespaceAssets = {}
      this.namespaceAssetsInitialState = {}
      this.namespaceEnabled = false
      this.application = undefined
      this.isApplication = false
      this.isApplicationInitialState = false
    },
  },
}
</script>
