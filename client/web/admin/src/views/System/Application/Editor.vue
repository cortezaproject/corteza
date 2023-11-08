<template>
  <b-container
    v-if="application"
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="title"
    >
      <b-button
        v-if="applicationID && canCreate"
        data-test-id="button-new-application"
        variant="primary"
        :to="{ name: 'system.application.new' }"
      >
        {{ $t('new') }}
      </b-button>

      <c-permissions-button
        v-if="applicationID && canGrant"
        :title="application.name || applicationID"
        :target="application.name || applicationID"
        :resource="`corteza::system:application/${applicationID}`"
      >
        <font-awesome-icon :icon="['fas', 'lock']" />
        {{ $t('permissions') }}
      </c-permissions-button>
    </c-content-header>

    <c-application-editor-info
      :application="application"
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      @submit="onInfoSubmit"
      @delete="onDelete"
    />

    <c-application-editor-unify
      v-if="applicationID && application.unify && application.applicationID"
      class="mt-3"
      :unify="application.unify"
      :application="application"
      :can-pin="canPin"
      :processing="unify.processing"
      :success="unify.success"
      @change-detected="unifyAssetStateChange = true"
      @submit="onUnifySubmit"
    />
  </b-container>
</template>
<script>
import { isEqual } from 'lodash'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CApplicationEditorInfo from 'corteza-webapp-admin/src/components/Application/CApplicationEditorInfo'
import CApplicationEditorUnify from 'corteza-webapp-admin/src/components/Application/CApplicationEditorUnify'
import { system } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'

export default {
  components: {
    CApplicationEditorInfo,
    CApplicationEditorUnify,
  },

  i18nOptions: {
    namespaces: 'system.applications',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    applicationID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      application: undefined,
      initialApplicationState: undefined,

      info: {
        processing: false,
        success: false,
      },

      unify: {
        processing: false,
        success: false,
      },

      unifyAssetStateChange: false,
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('system/', 'application.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },

    canPin () {
      return this.can('system/', 'pin')
    },

    title () {
      return this.applicationID ? this.$t('title.edit') : this.$t('title.create')
    },
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  watch: {
    applicationID: {
      immediate: true,
      handler () {
        if (this.applicationID) {
          this.fetchApplication()
        } else {
          this.application = new system.Application()

          this.initialApplicationState = this.application.clone()
        }
      },
    },
  },

  methods: {
    fetchApplication () {
      this.incLoader()

      this.$SystemAPI.applicationRead({ applicationID: this.applicationID, incFlags: 1 })
        .then((application = {}) => {
          if (!application.unify) {
            application.unify = {
              listed: true,
              pinned: false,
              name: this.application.name,
              config: '',
              icon: '',
              logo: '',
              url: '',
            }
          }

          application.unify.pinned = (application.flags || []).includes('pinned')
          application.unify.name = application.unify.name ? application.unify.name : application.name

          this.application = new system.Application(application)
          this.initialApplicationState = this.application.clone()
        },)
        .catch(this.toastErrorHandler(this.$t('notification:application.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onInfoSubmit (application) {
      this.info.processing = true

      if (this.applicationID) {
        application = {
          ...application,
          unify: this.initialApplicationState.unify,
        }

        this.$SystemAPI.applicationUpdate(application)
          .then(application => {
            this.initialApplicationState = new system.Application({
              ...application,
              unify: this.initialApplicationState.unify,
            })

            this.application = new system.Application({
              ...application,
              unify: this.application.unify,
            })

            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:application.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:application.update.error')))
          .finally(() => {
            this.info.processing = false
          })
      } else {
        this.$SystemAPI.applicationCreate(application)
          .then(({ applicationID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:application.create.success'))

            this.$router.push({ name: 'system.application.edit', params: { applicationID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:application.create.error')))
          .finally(() => {
            this.info.processing = false
          })
      }
    },

    async onUnifySubmit ({ unify, unifyAssets }) {
      this.unify.processing = true

      // Firstly handle any new application assets
      if (unifyAssets.logo || unifyAssets.icon) {
        try {
          const assets = await this.uploadAssets(unifyAssets)
          unify = { ...unify, ...assets }
        } catch (e) {
          this.toastErrorHandler(this.$t('notification:application.assetsUpload.error'))(e)
          this.unify.processing = false
          return
        }
      }

      if (this.applicationID) {
        const flagPayload = {
          applicationID: this.applicationID,
          flag: 'pinned',
          ownedBy: '0',
        }

        if (unify.pinned) {
          await this.$SystemAPI.applicationFlagCreate(flagPayload)
            .catch(() => {})
        } else {
          await this.$SystemAPI.applicationFlagDelete(flagPayload)
            .catch(() => {})
        }

        return this.$SystemAPI.applicationUpdate({ ...this.initialApplicationState, unify })
          .then(() => {
            this.application = new system.Application({ ...this.application, unify })
            this.initialApplicationState = new system.Application({
              ...this.initialApplicationState,
              unify,
            })

            this.unifyAssetStateChange = false

            this.animateSuccess('unify')
            this.toastSuccess(this.$t('notification:application.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:application.update.error')))
          .finally(() => {
            this.unify.processing = false
          })
      }
    },

    async uploadAssets (assets) {
      const rr = {}

      const rq = async (file) => {
        var formData = new FormData()
        formData.append('upload', file)

        const rsp = await this.$SystemAPI.api().request({
          method: 'post',
          url: this.$SystemAPI.applicationUploadEndpoint(),
          data: formData,
        })
        if (rsp.data.error) {
          throw new Error(rsp.data.error)
        }
        return rsp.data.response
      }

      const baseURL = this.$SystemAPI.baseURL

      if (assets.logo) {
        const rsp = await rq(assets.logo)
        rr.logo = baseURL + rsp.url
        rr.logoID = rsp.attachmentID

        assets.logo = undefined
      }

      if (assets.icon) {
        const rsp = await rq(assets.icon)
        rr.icon = baseURL + rsp.url
        rr.iconID = rsp.attachmentID

        assets.icon = undefined
      }

      return rr
    },

    onDelete () {
      this.incLoader()

      if (this.application.deletedAt) {
        this.$SystemAPI.applicationUndelete({ applicationID: this.applicationID })
          .then(() => {
            this.fetchApplication()

            this.toastSuccess(this.$t('notification:application.undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:application.undelete.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI.applicationDelete({ applicationID: this.applicationID })
          .then(() => {
            this.fetchApplication()

            this.application.deletedAt = new Date()

            this.toastSuccess(this.$t('notification:application.delete.success'))
            this.$router.push({ name: 'system.application' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:application.delete.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    checkUnsavedChanges (next, to) {
      const isNewPage = this.$route.path.includes('/new') && to.name.includes('edit')

      if (isNewPage || this.application.deletedAt) {
        next(true)
      } else if (!to.name.includes('edit')) {
        next(!isEqual(this.application, this.initialApplicationState) || this.unifyAssetStateChange ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },
  },
}
</script>
