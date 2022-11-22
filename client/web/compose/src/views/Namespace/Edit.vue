<template>
  <div
    v-if="loaded"
    class="d-flex flex-column w-100 h-100"
  >
    <portal to="topbar-title">
      {{ pageTitle }}
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="isEdit"
        size="sm"
        class="mr-1"
      >
        <b-button
          data-test-id="button-visit-namespace"
          variant="primary"
          class="d-flex align-items-center"
          :to="openNamespace"
          :disabled="!namespaceEnabled"
        >
          {{ $t('visit') }}
        </b-button>
        <b-button
          v-if="namespace.canManageNamespace"
          :title="$t('configure')"
          data-test-id="button-visit-admin-panel"
          variant="primary"
          class="d-flex align-items-center"
          :to="{ name: 'admin.modules', params: { slug: namespace.slug } }"
          style="margin-left:2px;"
        >
          <font-awesome-icon
            :icon="['fas', 'cogs']"
          />
        </b-button>
        <namespace-translator
          v-if="namespace"
          :namespace="namespace"
          :disabled="isNew"
          style="margin-left:2px;"
        />
      </b-button-group>
    </portal>

    <div class="flex-grow-1 overflow-auto mb-2">
      <b-container
        fluid="xl"
        class="flex-grow-1"
      >
        <div
          v-if="isEdit"
          class="d-flex align-items-center mt-1 mb-2"
        >
          <b-btn
            data-test-id="button-export-namespace"
            variant="light"
            size="lg"
            class="ml-1"
            @click="exportNamespace"
          >
            {{ $t('export') }}
          </b-btn>

          <c-permissions-button
            v-if="namespace.canGrant"
            data-test-id="button-permissions"
            :title="namespace.name"
            :target="namespace.name"
            :resource="'corteza::compose:namespace/'+namespace.namespaceID"
            button-variant="light"
            :button-label="$t('label.permissions')"
            class="ml-1 btn-lg"
          />
        </div>

        <b-card
          body-class="p-3"
          footer-bg-variant="warning"
        >
          <b-form>
            <b-form-group :label="$t('name.label')">
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
                    button-variant="light"
                    :disabled="isNew"
                  />
                </b-input-group-append>
              </b-input-group>
            </b-form-group>
            <b-form-group
              :label="$t('slug.label')"
              :description="$t('slug.description')"
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

            <b-form-group :label="$t('subtitle.label')">
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
                    highlight-key="subtitle"
                    button-variant="light"
                    :disabled="isNew"
                  />
                </b-input-group-append>
              </b-input-group>
            </b-form-group>

            <b-form-group
              :label="$t('description.label')"
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
                    highlight-key="description"
                    button-variant="light"
                    :disabled="isNew"
                  />
                </b-input-group-append>
              </b-input-group>
            </b-form-group>
            <hr>

            <b-form-group
              :label="$t('sidebar.configure')"
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
      :back-link="{ name: 'namespace.manage' }"
      :hide-delete="hideDelete"
      :hide-clone="!isEdit"
      :hide-save="hideSave"
      :disable-save="disableSave"
      @delete="handleDelete"
      @save="handleSave()"
      @clone="$router.push({ name: 'namespace.clone', params: { namespaceID: namespace.namespaceID }})"
      @saveAndClose="handleSave({ closeOnSuccess: true })"
    />

    <b-modal
      id="logo"
      hide-header
      hide-footer
      centered
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
import { compose, NoID } from '@cortezaproject/corteza-js'
import { url, handle } from '@cortezaproject/corteza-vue'
import EditorToolbar from 'corteza-webapp-compose/src/components/Admin/EditorToolbar'
import NamespaceTranslator from 'corteza-webapp-compose/src/components/Namespaces/NamespaceTranslator'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'namespace',
  },

  components: {
    EditorToolbar,
    NamespaceTranslator,
  },

  data () {
    return {
      loaded: false,
      processing: false,

      namespace: new compose.Namespace({ enabled: true }),
      namespaceAssets: {
        logo: undefined,
        icon: undefined,
      },
      namespaceEnabled: false,

      application: undefined,
      isApplication: false,
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
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
      return this.$route.name === 'namespace.edit'
    },

    isClone () {
      return this.$route.name === 'namespace.clone'
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
      handler (namespaceID) {
        this.fetchNamespace(namespaceID)
      },
    },
  },

  methods: {
    async fetchNamespace () {
      this.processing = true

      const namespaceID = this.$route.params.namespaceID

      this.application = undefined
      this.isApplication = false

      if (namespaceID) {
        await this.$store.dispatch('namespace/findByID', { namespaceID })
          .then(ns => {
            this.namespaceEnabled = ns.enabled
            this.namespace = new compose.Namespace(ns)

            if (this.isClone) {
              this.namespace.name = `${ns.name} (${this.$t('cloneSuffix')})`
              this.namespace.slug = `${ns.slug}_${this.$t('cloneSuffix')}`
            }

            this.fetchApplication()
          })
      } else {
        this.namespace = new compose.Namespace({ enabled: true })
      }

      this.namespace.meta = {
        subtitle: '',
        description: '',
        hideSidebar: false,
        ...this.namespace.meta,
      }

      this.processing = false
      this.loaded = true
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
      this.$SystemAPI.applicationList({ name: this.namespace.slug })
        .then(({ set = [] }) => {
          if (set.length) {
            this.application = set[0]
            this.isApplication = this.application.enabled
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:namespace.application.fetchFailed')))
    },

    async handleSave ({ closeOnSuccess = false } = {}) {
      this.processing = true

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
        } catch (e) {
          const error = JSON.stringify(e) === '{}' ? '' : e
          this.toastErrorHandler(this.$t('notification:namespace.assetUploadFailed'))(error)
          this.processing = false
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
          await this.$store.dispatch('namespace/update', { ...payload, namespaceID }).then((ns) => {
            this.namespaceEnabled = ns.enabled
            this.namespace = new compose.Namespace(ns)

            this.toastSuccess(this.$t('notification:namespace.saved'))
          })
        } catch (e) {
          this.toastErrorHandler(this.$t('notification:namespace.saveFailed'))(e)
          this.processing = false
          return
        }
      } else if (this.isClone) {
        try {
          await this.$store.dispatch('namespace/clone', { namespaceID, name, slug, enabled, meta }).then((ns) => {
            this.namespace = new compose.Namespace(ns)
          })
        } catch (e) {
          this.toastErrorHandler(this.$t('notification:namespace.cloneFailed'))(e)
          this.processing = false
          return
        }
      } else {
        try {
          await this.$store.dispatch('namespace/create', payload).then((ns) => {
            this.namespaceEnabled = ns.enabled
            this.namespace = new compose.Namespace(ns)

            this.toastSuccess(this.$t('notification:namespace.saved'))
          })
        } catch (e) {
          this.toastErrorHandler(this.$t('notification:namespace.createFailed'))(e)
          this.processing = false
          return
        }
      }

      await this.handleApplicationSave()
        .catch(() => this.toastErrorHandler(this.$t('notification:namespace.createAppFailed')))

      this.processing = false

      if (closeOnSuccess) {
        this.$router.push({ name: 'namespace.manage' })
      } else if (!this.isEdit || this.isClone) {
        this.$router.push({ name: 'namespace.edit', params: { namespaceID: this.namespace.namespaceID } })
      }

      this.namespace.meta = {
        subtitle: '',
        description: '',
        hideSidebar: false,
        ...this.namespace.meta,
      }
    },

    handleDelete () {
      this.processing = true

      const { namespaceID } = this.namespace
      const { applicationID } = this.application || {}

      this.$store.dispatch('namespace/delete', { namespaceID })
        .catch(this.toastErrorHandler(this.$t('notification:namespace.deleteFailed')))
        .then(() => {
          if (applicationID) {
            return this.$SystemAPI.applicationDelete({ applicationID })
          }
        })
        .finally(() => {
          this.processing = false
          this.$router.push({ name: 'namespace.manage' })
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
          .then(app => { this.application = app })
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
          .then(app => { this.application = app })
          .catch(this.toastErrorHandler(this.$t('notification:namespace.application.createFailed')))
      }
    },

    async uploadAssets () {
      const rr = {}

      const rq = async (file) => {
        var formData = new FormData()
        formData.append('upload', file)

        const rsp = await this.$ComposeAPI.api().request({
          method: 'post',
          url: this.$ComposeAPI.namespaceUploadEndpoint(),
          data: formData,
        })
        if (rsp.data.error) {
          throw new Error(rsp.data.error)
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
  },
}
</script>
