<template>
  <div>
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <div
      v-else-if="mode === 'list'"
    >
      <draggable
        :list.sync="attachments"
        :disabled="!enableOrder"
      >
        <b-row
          v-for="(a, index) in attachments"
          :key="a.attachmentID"
          no-gutters
          class="flex-nowrap item mb-2"
        >
          <b-col cols="auto">
            <div
              v-if="enableOrder"
              class="d-inline p-1 mr-2"
            >
              <font-awesome-icon
                :icon="['fas', 'bars']"
                class="handle text-secondary"
              />
            </div>
          </b-col>

          <b-col style="word-break:break-all;">
            <attachment-link :attachment="a" />

            <i18next
              path="general.label.attachmentFileInfo"
              tag="small"
              class="d-block text-muted"
            >
              <span>{{ size(a) }}</span>

              <span>{{ uploadedAt(a) }}</span>
            </i18next>
          </b-col>

          <b-col
            cols="auto"
            class="d-flex align-items-start"
          >
            <b-button
              v-if="a.download"
              :href="a.download"
              variant="outline-light"
              size="sm"
              class="border-0 text-primary px-2 ml-2"
            >
              <font-awesome-icon :icon="['fas', 'download']" />
            </b-button>

            <c-input-confirm
              v-if="enableDelete"
              show-icon
              class="ml-2"
              @confirmed="deleteAttachment(index)"
            />
          </b-col>
        </b-row>
      </draggable>
    </div>

    <div
      v-else
      class="d-flex align-items-center justify-content-around flex-wrap h-100"
    >
      <div
        v-for="a in attachments"
        :key="a.attachmentID"
        :class="{ 'h-100': attachments.length === 1, 'w-100': !canPreview(a) }"
        class="item mb-2"
      >
        <c-preview-inline
          v-if="canPreview(a)"
          :src="inlineUrl(a)"
          :meta="a.meta"
          :name="a.name"
          :alt="a.name"
          :preview-style="{ width: 'unset', ...inlineCustomStyles(a) }"
          :preview-class="[
            !previewOptions.clickToView ? 'disable-zoom-cursor' : '',

          ]"
          :labels="previewLabels"
          class="mb-1"
          @openPreview="openLightbox({ ...a, ...$event })"
        />

        <div
          v-if="!hideFileName"
          class="d-flex align-items-start justify-content-between"
          :class="{ 'w-100': canPreview(a) }"
        >
          <div
            :class="{ 'text-center': canPreview(a) }"
            class="text-wrap"
          >
            <attachment-link
              :attachment="a"
            />
          </div>

          <b-button
            v-if="a.download"
            :href="a.download"
            variant="outline-light"
            class="border-0 text-primary px-2"
          >
            <font-awesome-icon :icon="['fas', 'download']" />
          </b-button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import numeral from 'numeral'
import moment from 'moment'
import { compose, shared } from '@cortezaproject/corteza-js'
import AttachmentLink from './Link'
import draggable from 'vuedraggable'
import { url, components } from '@cortezaproject/corteza-vue'
const { CPreviewInline, canPreview } = components

export default {
  i18nOptions: {
    namespaces: 'preview',
  },

  components: {
    CPreviewInline,
    AttachmentLink,
    draggable,
  },

  props: {
    enableDelete: {
      type: Boolean,
    },

    enableOrder: {
      type: Boolean,
      default: false,
    },

    namespace: {
      type: compose.Namespace,
      required: true,
    },

    kind: {
      type: String,
      required: true,
    },

    mode: {
      type: String,
      required: true,
    },

    set: {
      type: Array,
      required: true,
    },

    hideFileName: {
      type: Boolean,
      default: false,
    },

    previewOptions: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      processing: false,

      attachments: [],
    }
  },

  computed: {
    inlineUrl () {
      return (a) => (this.ext(a) === 'pdf' ? a.download : a.url)
    },

    previewLabels () {
      return {
        loading: this.$t('pdf.loading'),
        firstPagePreview: this.$t('pdf.firstPagePreview'),
        pageLoadFailed: this.$t('pdf.pageLoadFailed'),
        pageLoading: this.$t('pdf.pageLoading'),
      }
    },

    canPreview () {
      return (a) => {
        const meta = a.meta || {}
        const type = (meta.preview || meta.original || {}).mimetype
        const src = this.inlineUrl(a)
        return canPreview({ type, src, name: a.name })
      }
    },

    baseURL () {
      return url.Make({ url: window.CortezaAPI + '/compose' })
    },
  },

  watch: {
    set: {
      immediate: true,
      handler (set) {
        // Handle attachments provided as objects
        const att = set.map(a => {
          if (typeof a === 'object') {
            return new shared.Attachment(a, this.baseURL)
          } else {
            return null
          }
        })

        // Handle attachmentsprovided as attachmentID
        const namespaceID = this.namespace.namespaceID

        this.processing = true

        Promise.all(Object.entries(set).map(([index, attachmentID]) => {
          if (typeof attachmentID === 'string') {
            return this.$ComposeAPI.attachmentRead({ kind: this.kind, attachmentID, namespaceID }).then(a => {
              att.splice(index, 1, new shared.Attachment(a, this.baseURL))
            })
          }
        }))
          .then(() => {
          // Filter out invalid/missing attachments
            const { clickToView = true, enableDownload = true } = this.previewOptions

            this.attachments = att
              .filter(a => !!a)
              .filter(a => typeof a === 'object')
              .map(a => {
                return {
                  ...a,
                  download: enableDownload ? a.download : undefined,
                  clickToView,
                }
              })
          })
          .finally(() => {
            this.processing = false
          })
      },
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    size (a) {
      return numeral(a.meta.original.size).format('0b')
    },

    uploadedAt (a) {
      return moment(a.updatedAt || a.createdAt).fromNow()
    },

    openLightbox (e) {
      if (!this.previewOptions.clickToView) return

      this.$root.$emit('showAttachmentsModal', e)
    },

    deleteAttachment (index) {
      this.attachments.splice(index, 1)
      this.$emit('update:set', this.attachments.map(a => a.attachmentID))
    },

    ext (a) {
      const { meta } = a
      switch (meta && meta.original ? meta.original.ext : null) {
        case 'odt':
        case 'doc':
        case 'docx':
          return 'word'
        case 'pdf':
          return 'pdf'
        case 'ppt':
        case 'pptx':
          return 'powerpoint'
        case 'zip':
        case 'rar':
          return 'archive'
        case 'xls':
        case 'xlsx':
        case 'csv':
          return 'excel'
        case 'mov':
        case 'mp3':
        case 'mp4':
          return 'video'
        case 'png':
        case 'jpg':
        case 'jpeg':
          return 'image'
        default: return 'alt'
      }
    },

    inlineCustomStyles (a) {
      const {
        width,
        height,
        borderRadius,
        backgroundColor,
      } = this.previewOptions
      let { maxWidth, maxHeight, margin } = this.previewOptions

      maxWidth = maxWidth || '100%'
      maxHeight = maxHeight || '100%'
      margin = margin || 'auto'

      if (this.ext(a) === 'image') {
        return {
          height,
          width,
          maxHeight,
          maxWidth,
          borderRadius,
          backgroundColor,
          margin,
        }
      }

      return {}
    },

    setDefaultValues () {
      this.processing = false
      this.attachments = []
    },
  },
}
</script>

<style lang="scss" scoped>
.handle {
  cursor: grab;
}

.item:hover {
  background-color: var(--gray-200);
}
</style>
