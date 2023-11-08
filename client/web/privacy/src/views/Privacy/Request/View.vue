<template>
  <b-container
    class="d-flex flex-column p-3"
  >
    <div
      v-if="processing.request"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <template v-else-if="request">
      <request-viewer
        :request="request"
        class="mb-3"
      />

      <request-comments
        :comments="comments"
        :processing="processing.comments"
        :sort="sort"
        @sort="sort = $event"
        @submit="submitComment"
      />
    </template>

    <portal to="editor-toolbar">
      <editor-toolbar
        :processing="processing.toolbar"
        :processing-confirm="processingReject"
        :back-link="{ name: 'request.list' }"
        :delete-show="isDC"
        :delete-disabled="!request || !isPending"
        :delete-label="$t('reject')"
        @delete="handleRequest('rejected')"
      >
        <c-input-confirm
          v-if="!isDC"
          :disabled="!request || !isPending"
          :processing="processingCancel"
          :text="$t('cancel')"
          variant="light"
          size="lg"
          size-confirm="lg"
          @confirmed="handleRequest('canceled')"
        />

        <c-input-confirm
          v-else
          :disabled="!request || !isPending"
          :processing="processingApprove"
          :text="$t('approve')"
          variant="primary"
          variant-ok="primary"
          size="lg"
          size-confirm="lg"
          class="ml-2"
          @confirmed="handleRequest('approved')"
        />
      </editor-toolbar>
    </portal>
  </b-container>
</template>

<script>
import EditorToolbar from 'corteza-webapp-privacy/src/components/Common/EditorToolbar'
import RequestViewer from 'corteza-webapp-privacy/src/components/Requests/Viewer'
import RequestComments from 'corteza-webapp-privacy/src/components/Requests/Comments'

export default {
  name: 'RequestView',

  i18nOptions: {
    namespaces: 'request',
    keyPrefix: 'view',
  },

  components: {
    EditorToolbar,
    RequestViewer,
    RequestComments,
  },

  props: {
    requestID: {
      type: String,
      required: false,
      default: '',
    },
  },

  data () {
    return {
      processing: {
        comments: false,
        request: false,
        toolbar: false,
      },

      processingApprove: false,
      processingReject: false,
      processingCancel: false,

      isDC: null,

      sort: 'createdAt DESC',

      request: undefined,
      comments: [],
    }
  },

  computed: {
    isPending () {
      return this.request.status === 'pending'
    },
  },

  watch: {
    requestID: {
      immediate: true,
      handler () {
        if (this.requestID) {
          this.fetchRequest()
          this.fetchComments()
        }
      },
    },

    sort: {
      handler () {
        if (this.requestID) {
          this.fetchComments()
        }
      },
    },
  },

  created () {
    this.checkIsDC()
  },

  methods: {
    checkIsDC () {
      this.$SystemAPI.roleList({ query: 'data-privacy-officer', memberID: this.$auth.user.userID })
        .then(({ set = [] }) => {
          this.isDC = !!set.length
        })
    },

    fetchRequest (requestID = this.requestID) {
      this.processing.request = true

      return this.$SystemAPI.dataPrivacyRequestRead({ requestID })
        .then(request => {
          this.request = request
        })
        .catch(this.toastErrorHandler(this.$t('notification:list.load.error')))
        .finally(() => {
          this.processing.request = false
        })
    },

    fetchComments (requestID = this.requestID) {
      this.processing.comments = true

      return this.$SystemAPI.dataPrivacyRequestCommentList({ requestID, sort: this.sort })
        .then(({ set }) => {
          this.comments = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:list.load.error')))
        .finally(() => {
          this.processing.comments = false
        })
    },

    handleRequest (status) {
      this.processing.toolbar = true

      if (status === 'approved') {
        this.processingApprove = true
      } else if (status === 'rejected') {
        this.processingReject = true
      } else {
        this.processingCancel = true
      }

      this.$SystemAPI.dataPrivacyRequestUpdateStatus({ requestID: this.requestID, status })
        .then(() => {
          this.$router.push({ name: 'request.list' })
        })
        .finally(() => {
          this.processing.toolbar = false

          if (status === 'approved') {
            this.processingApprove = false
          } else if (status === 'rejected') {
            this.processingReject = false
          } else {
            this.processingCancel = false
          }
        })
    },

    submitComment (comment) {
      this.processing.comments = true

      this.$SystemAPI.dataPrivacyRequestCommentCreate({ requestID: this.requestID, comment })
        .then(() => {
          return this.fetchComments()
        })
        .finally(() => {
          this.processing.comments = false
        })
    },
  },
}
</script>
