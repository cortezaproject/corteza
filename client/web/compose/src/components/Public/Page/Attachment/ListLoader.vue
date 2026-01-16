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
        handle=".handle"
      >
        <b-row
          v-for="(a, index) in attachments"
          :key="a.attachmentID"
          no-gutters
          class="list-item flex-nowrap mb-1 rounded"
        >
          <b-col cols="auto">
            <font-awesome-icon
              v-if="enableOrder"
              :icon="['fas', 'bars']"
              class="handle text-secondary my-1 mr-3"
              style="padding-top: 0.05rem;"
            />
          </b-col>

          <b-col>
            <div class="d-flex flex-column flex-wrap align-items-start">
              <div
                class="d-flex align-items-start gap-1"
                style="word-break: break-all;"
              >
                <div style="margin-top: 0.1rem;">
                  <attachment-link :attachment="a">
                    {{ a.name }}
                  </attachment-link>
                </div>

                <div class="d-flex align-items-center gap">
                  <b-button
                    v-if="a.download"
                    :href="a.download"
                    variant="outline-extra-light"
                    size="sm"
                    class="download-button border-0"
                    @click.stop
                  >
                    <font-awesome-icon
                      :icon="['fas', 'download']"
                      class="text-secondary"
                    />
                  </b-button>

                  <c-input-confirm
                    v-if="enableDelete"
                    show-icon
                    class="delete-button"
                    @confirmed="deleteAttachment(index)"
                  />
                </div>
              </div>

              <i18next
                path="general.label.attachmentFileInfo"
                tag="small"
                class="d-block text-muted"
              >
                <span>{{ size(a) }}</span>

                <span>{{ uploadedAt(a) }}</span>
              </i18next>
            </div>
          </b-col>
        </b-row>
      </draggable>
    </div>

    <div
      v-else
      class="d-flex align-items-start justify-content-around gap-3 flex-wrap h-100"
    >
      <div
        v-for="a in attachments"
        :key="a.attachmentID"
        class="item-preview"
      >
        <c-preview-inline
          v-if="canPreview(a)"
          :src="inlineUrl(a)"
          :title="a.name"
          :meta="a.meta"
          :name="a.name"
          :alt="a.name"
          :preview-style="{ width: 'unset', ...inlineCustomStyles(a) }"
          :labels="previewLabels"
          @openPreview="openLightbox({ ...a, ...$event })"
        />

        <div
          class="d-flex align-items-start justify-content-center"
          :style="{ width: `calc(${inlineCustomStyles(a).width})` }"
        >
          <div
            v-if="!hideFileName"
            class="text-wrap filename-container text-center"
            :style="{ marginTop: '0.1rem' }"
          >
            <attachment-link :attachment="a" />
          </div>
        </div>
        <b-button
          v-if="a.download"
          :href="a.download"
          variant="extra-light"
          size="sm"
          class="preview-download-button border-0"
          @click.stop
        >
          <font-awesome-icon
            :icon="['fas', 'download']"
            class="text-secondary"
          />
        </b-button>
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
const { CPreviewInline, canPreview, getExtensionIconType } = components

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
      return (a) => a.url
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

          return Promise.resolve([])
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
      if (this.ext(e) === 'pdf') {
        window.open(e.url, '_blank')
      } else {
        this.$root.$emit('showAttachmentsModal', e)
      }
    },

    deleteAttachment (index) {
      this.attachments.splice(index, 1)
      this.$emit('update:set', this.attachments.map(a => a.attachmentID))
    },

    ext (a) {
      const { meta } = a
      const { original = {} } = meta || {}
      const { ext } = original || {}
      return getExtensionIconType(ext)
    },

    inlineCustomStyles (a) {
      const {
        borderRadius,
        backgroundColor,
      } = this.previewOptions
      let { width, height, maxWidth, maxHeight, margin } = this.previewOptions

      maxWidth = maxWidth || '100%'
      maxHeight = maxHeight || '100%'
      margin = margin || 'auto'

      if (this.ext(a) !== 'image') {
        width = width || '200px'
        height = height || 'auto'
      }

      return {
        width,
        height,
        maxWidth,
        maxHeight,
        borderRadius,
        backgroundColor,
        margin,
      }
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

.list-item {
  .download-button {
    opacity: 0;
    transition: opacity 0.2s;
  }

  &:hover .download-button {
    opacity: 1;
  }

  .delete-button {
    opacity: 0;
    transition: opacity 0.2s;
  }

  &:hover .delete-button {
    opacity: 1;
  }

  &:hover {
    background-color: var(--light);
  }
}

.item-preview {
  position: relative;
  .preview-download-button {
    position: absolute;
    top: 0;
    right: 0;
    opacity: 0;
    transition: opacity 0.2s;
  }

  &:hover .preview-download-button {
    opacity: 1;
  }

  .filename-container {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-word;
    max-width: 100%;

    &:hover {
      -webkit-line-clamp: unset;
      line-clamp: unset;
      overflow: visible;
    }
  }
}
</style>
