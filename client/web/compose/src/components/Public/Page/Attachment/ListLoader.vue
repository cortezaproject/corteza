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
      class="list"
    >
      <draggable
        :list.sync="attachments"
        :disabled="!enableOrder"
      >
        <div
          v-for="(a, index) in attachments"
          :key="a.attachmentID"
          class="item"
        >
          <b-row no-gutters>
            <b-col>
              <div
                v-if="enableOrder"
                class="mr-2 d-inline"
              >
                <font-awesome-icon
                  v-b-tooltip.hover
                  :icon="['fas', 'bars']"
                  :title="$t('general.tooltip.dragAndDrop')"
                  class="handle text-light"
                />
              </div>
              <attachment-link :attachment="a" />
              &nbsp;
              <i18next
                path="general.label.attachmentFileInfo"
                tag="label"
              >
                <span>{{ size(a) }}</span>
                <span>{{ uploadedAt(a) }}</span>
              </i18next>
            </b-col>
            <div class="col-sm-2 text-right my-auto">
              <a
                :href="a.download"
                class="px-0 btn text-primary mr-2"
              >
                <font-awesome-icon :icon="['fas', 'download']" />
              </a>
              <b-button
                v-if="enableDelete"
                variant="link"
                class="px-0"
                @click="deleteAttachment(index)"
              >
                <font-awesome-icon
                  :icon="['far', 'trash-alt']"
                  class="action text-danger"
                />
              </b-button>
            </div>
          </b-row>
        </div>
      </draggable>
    </div>

    <div
      v-else-if="mode === 'grid'"
      class="grid"
    >
      <div
        v-for="a in attachments"
        :key="a.attachmentID"
        class="p-2"
      >
        <attachment-link
          :attachment="a"
          class="d-block"
        >
          <font-awesome-icon
            :icon="['far', 'file-'+ext(a)]"
            class="text-dark float-left mr-2"
          />
          {{ a.name }}
        </attachment-link>
        <i18next
          path="general.label.attachmentFileInfo"
          tag="label"
        >
          <span>{{ size(a) }}</span>
          <span>{{ uploadedAt(a) }}</span>
        </i18next>
      </div>
    </div>

    <div
      v-else
      class="single gallery"
    >
      <div
        v-for="(a) in files"
        :key="a.attachmentID"
        class="my-2"
      >
        <c-preview-inline
          v-if="canPreview(a)"
          class="ml-0"
          :src="inlineUrl(a)"
          :meta="a.meta"
          :name="a.name"
          :alt="a.name"
          :preview-style="{ width: 'unset' }"
          :labels="previewLabels"
          @openPreview="openLightbox({ ...a, ...$event })"
        />

        <div v-else>
          <font-awesome-icon
            :icon="['far', 'file-'+ext(a)]"
            title="Open bookmarks"
          />
        </div>

        <div
          v-if="!hideFileName"
          class="m-1"
        >
          <attachment-link
            :attachment="a"
          />
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
  },

  data () {
    return {
      processing: false,

      attachments: [],
    }
  },

  computed: {
    inlineUrl () {
      return (a) => (this.ext(a) === 'pdf' ? a.download : a.previewUrl)
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

    files () {
      if (this.mode === 'single') {
        return this.attachments.slice(this.attachments.length - 1)
      } else {
        return this.attachments
      }
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
            this.attachments = att
              .filter(a => !!a)
              .filter(a => typeof a === 'object')
          })
          .finally(() => {
            this.processing = false
          })
      },
    },
  },

  methods: {
    size (a) {
      return numeral(a.meta.original.size).format('0b')
    },

    uploadedAt (a) {
      return moment(a.updatedAt || a.createdAt).fromNow()
    },

    openLightbox (e) {
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
  },
}
</script>
<style lang="scss" scoped>

.grid {
  .svg-inline--fa {
    font-size: 40px;
  }
}

.single {
  .svg-inline--fa {
    font-size: 40px;
  }

  img {
    cursor: pointer;
  }
}

.handle {
  cursor: grab;
}
</style>
