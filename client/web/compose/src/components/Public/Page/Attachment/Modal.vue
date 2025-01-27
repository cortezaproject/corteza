<template>
  <c-preview-lightbox
    v-if="show"
    :src="(attachment || {}).document || (attachment || {}).src"
    :name="(attachment || {}).name"
    :alt="(attachment || {}).name"
    :labels="previewLabels"
    :meta="(attachment || {}).meta"
    @close="attachment=undefined"
  >
    <p
      slot="header.left"
      class="m-0"
    >
      {{ (attachment || {}).name }}
    </p>

    <a
      v-if="attachment.download"
      slot="header.right"
      :href="(attachment || {}).download"
    >
      {{ $t('general.label.download') }}
    </a>
  </c-preview-lightbox>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CPreviewLightbox } = components

export default {
  i18nOptions: {
    namespaces: 'preview',
  },

  components: {
    CPreviewLightbox,
  },

  data () {
    return {
      attachment: undefined,
    }
  },

  computed: {
    show: {
      get () {
        return !!this.attachment
      },

      set (show) {
        if (!show) {
          this.attachment = undefined
        }
      },
    },

    previewLabels () {
      return {
        loading: this.$t('pdf.loading'),
        downloadForAll: this.$t('pdf.downloadForAll'),
        pageLoadFailed: this.$t('pdf.pageLoadFailed'),
        pageLoading: this.$t('pdf.pageLoading'),
        noPages: this.$t('pdf.noPages'),
        clickToRetry: this.$t('pdf.clickToRetry'),
      }
    },
  },

  created () {
    window.addEventListener('keyup', this.onKeyUp)
    this.$root.$on('showAttachmentsModal', this.showAttachmentModal)
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    onKeyUp ({ key }) {
      if (key === 'Escape') {
        this.attachment = undefined
      }
    },

    showAttachmentModal ({ url, download, name, document = undefined, meta, enableDownload }) {
      this.attachment = {
        document,
        download,
        meta,
        src: url,
        name,
        caption: name,
        enableDownload,
      }
    },

    setDefaultValues () {
      this.attachment = undefined
    },

    destroyEvents () {
      window.removeEventListener('keyup', this.onKeyUp)
      this.$root.$off('showAttachmentsModal', this.showAttachmentModal)
    },
  },
}
</script>
