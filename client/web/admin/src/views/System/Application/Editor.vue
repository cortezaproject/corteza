<template>
  <b-container
    v-if="application"
    class="py-3"
  >
    <c-content-header
      :title="title"
    >
      <span
        class="text-nowrap"
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
          button-variant="light"
          class="ml-2"
        >
          <font-awesome-icon :icon="['fas', 'lock']" />
          {{ $t('permissions') }}
        </c-permissions-button>
      </span>
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
      v-if="application.unify && application.applicationID"
      class="mt-3"
      :unify="application.unify"
      :application="application"
      :can-pin="canPin"
      :processing="unify.processing"
      :success="unify.success"
      @submit="onUnifySubmit"
    />
  </b-container>
</template>
<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CApplicationEditorInfo from 'corteza-webapp-admin/src/components/Application/CApplicationEditorInfo'
import CApplicationEditorUnify from 'corteza-webapp-admin/src/components/Application/CApplicationEditorUnify'
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

      info: {
        processing: false,
        success: false,
      },
      unify: {
        processing: false,
        success: false,
      },
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

  watch: {
    applicationID: {
      immediate: true,
      handler () {
        if (this.applicationID) {
          this.fetchApplication()
        } else {
          this.application = {}
        }
      },
    },
  },

  methods: {
    fetchApplication () {
      this.incLoader()

      this.$SystemAPI.applicationRead({ applicationID: this.applicationID, incFlags: 1 })
        .then(this.prepare)
        .catch(this.toastErrorHandler(this.$t('notification:application.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onInfoSubmit (application) {
      this.info.processing = true

      if (this.applicationID) {
        this.$SystemAPI.applicationUpdate(application)
          .then(() => {
            this.fetchApplication()

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

        return this.$SystemAPI.applicationUpdate({ ...this.application, unify })
          .then(() => {
            this.fetchApplication()

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

            this.toastSuccess(this.$t('notification:application.delete.success'))
            this.$router.push({ name: 'system.application' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:application.delete.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    prepare (application = {}) {
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

      this.application = application
    },
  },
}
</script>
