<template>
  <b-container
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="title"
    >
      <b-button
        v-if="nodeID"
        variant="link"
        @click="generate.modal = true"
      >
        {{ $t('generateUri') }}
      </b-button>
    </c-content-header>

    <c-federation-editor-info
      :node="node"
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      @submit="onInfoSubmit"
      @delete="onDelete"
    />

    <b-modal
      v-model="generate.modal"
      hide-header
      hide-footer
      centered
      size="lg"
      body-class="px-5"
    >
      <div
        class="text-center px-5"
      >
        <font-awesome-icon
          size="7x"
          :icon="['fas', 'share-alt']"
          class="text-light mb-2"
        />
        <h2>
          {{ $t('generate.description') }}
        </h2>
      </div>

      <b-input-group
        size="xl"
        class="mt-5"
      >
        <b-form-input
          v-model="generate.email"
          type="email"
          placeholder="email@example.com"
        />
        <b-input-group-append>
          <c-button-submit
            :disabled="!generate.url || !generate.email"
            :processing="generate.processing"
            :success="generate.success"
            :text="$t('generate.sendEmail')"
            @submit="sendEmail()"
          />
        </b-input-group-append>
      </b-input-group>

      <div
        class="mt-3"
      >
        <p>
          {{ $t('generate.subject') }} <strong>{{ $t('generate.invitation') }}</strong>
        </p>

        <p
          class="mt-4"
        >
          {{ $t('generate.hello') }}
        </p>

        <p>
          {{ $t('generate.body', { userLabel }) }}
        </p>

        <p
          class="text-center text-break"
        >
          <i>
            {{ generate.url || $t('generate.notGenerated') }}
          </i>
        </p>

        <p>
          {{ $t('generate.kindRegards') }}
        </p>
      </div>

      <hr
        class="my-3"
      >

      <b-button
        variant="link"
        size="sm"
        :to="{}"
        class="p-1"
        @click="copyUrl()"
      >
        <font-awesome-icon
          :icon="['far', 'copy']"
          class="text-secondary pointer"
        />
      </b-button>
      <span
        class="text-break"
      >
        {{ generate.url || $t('generate.notGenerated') }}
      </span>
    </b-modal>
  </b-container>
</template>

<script>
import { mapGetters } from 'vuex'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CFederationEditorInfo from 'corteza-webapp-admin/src/components/Federation/CFederationEditorInfo'
import { cloneDeep, isEqual } from 'lodash'

export default {
  i18nOptions: {
    namespaces: 'federation.nodes',
    keyPrefix: 'editor',
  },

  components: {
    CFederationEditorInfo,
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    nodeID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      node: {},
      initialNodeState: {},

      // Processing and success flags for each form
      info: {
        processing: false,
        success: false,
      },

      generate: {
        modal: false,
        processing: false,
        success: false,
        email: '',
        url: '',
      },
    }
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('federation/', 'node.create')
    },

    userLabel () {
      return this.$auth.user.name || this.$auth.user.email
    },

    title () {
      return this.nodeID ? this.$t('title.edit') : this.$t('title.create')
    },
  },

  watch: {
    nodeID: {
      immediate: true,
      handler () {
        if (this.nodeID) {
          this.fetchNode()
          this.fetchGeneratedUrl()
        } else {
          this.node = {
            name: '',
            baseURL: '',
            contact: '',
          }

          this.initialNodeState = {
            name: '',
            baseURL: '',
            contact: '',
          }
        }
      },
    },
  },

  methods: {
    fetchNode () {
      this.incLoader()

      this.$FederationAPI.nodeRead({ nodeID: this.nodeID })
        .then(node => {
          this.node = node // new federation.Node(node)
          this.initialNodeState = cloneDeep(node)
        })
        .catch(this.toastErrorHandler(this.$t('notification:federation.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    fetchGeneratedUrl () {
      this.incLoader()

      this.$FederationAPI.nodeGenerateUri({ nodeID: this.nodeID })
        .then(url => {
          this.generate.url = url
        })
        .catch(this.toastErrorHandler(this.$t('notification:federation.url.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    /**
     * Handles node info submit event, calls node update or create API endpoint
     * and handles response & errors
     *
     * @param node {Object}
     */
    onInfoSubmit (node) {
      this.info.processing = true

      const payload = { ...node }

      if (payload.nodeID) {
        // On update, reset the node obj
        this.$FederationAPI.nodeUpdate(payload)
          .then(node => {
            this.node = node
            this.initialNodeState = cloneDeep(node)

            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:federation.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:federation.update.error')))
          .finally(() => {
            this.info.processing = false
          })
      } else {
        // On creation, redirect to edit page
        this.$FederationAPI.nodeCreate(payload)
          .then(({ nodeID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:federation.create.success'))

            this.$router.push({ name: 'federation.nodes.edit', params: { nodeID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:federation.create.error')))
          .finally(() => {
            this.info.processing = false
          })
      }
    },

    /**
     * Handles node delete event, calls node delete API endpoint
     * and handles response & errors
     */
    onDelete () {
      this.incLoader()

      if (this.node.deletedAt) {
        this.$FederationAPI.nodeUndelete({ nodeID: this.nodeID })
          .then(() => {
            this.fetchNode()

            this.toastSuccess(this.$t('notification:federation.undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:federation.undelete.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$FederationAPI.nodeDelete({ nodeID: this.nodeID })
          .then(() => {
            this.fetchNode()

            this.node.deletedAt = new Date()

            this.toastSuccess(this.$t('notification:federation.delete.success'))
            this.$router.push({ name: 'federation.nodes' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:federation.delete.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    async sendEmail () {
      this.generate.processing = true

      const html = `
        <p class="mt-4">Hello,</p>
        <p>${this.userLabel} is sending you an invitation for Corteza Federated Network.</p>
        <p>To start sharing data between organizations, go to the admin panel of your Corteza application, click on &ldquo;Federation&rdquo; and select &ldquo;Pair Federation Network&rdquo; on top right corner.<br />Copy the link below and await confirmation from another administrator.</p>
        <blockquote>
        <p class="text-center text-break"><em>${this.generate.url}</em></p>
        </blockquote>
        <p>Kind regards, Corteza team.</p>
      `

      const values = {
        to: [this.generate.email],
        subject: this.$t('generate.invitation'),
        content: {
          html,
        },
      }
      await this.$ComposeAPI.notificationEmailSend(values)
        .then(f => {
          this.generate.email = ''

          this.animateSuccess('generate')
          this.toastSuccess(this.$t('notification:workflow.email.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:workflow.email.error')))
        .finally(() => {
          this.generate.processing = false
        })
    },

    copyUrl () {
      navigator.clipboard.writeText(this.generate.url)
    },

    checkUnsavedChanges (next, to) {
      const isNewPage = this.$route.path.includes('/new') && to.name.includes('edit')

      if (isNewPage || this.node.deletedAt) {
        next(true)
      } else if (!to.name.includes('edit')) {
        next(!isEqual(this.node, this.initialNodeState) ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },
  },
}
</script>
