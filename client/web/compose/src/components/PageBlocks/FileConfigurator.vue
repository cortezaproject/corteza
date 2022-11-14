<template>
  <b-tab :title="$t('kind.file.label')">
    <b-form-group
      horizontal
      :description="$t('kind.file.view.modeFootnote')"
      :label="$t('kind.file.view.modeLabel')"
    >
      <b-form-radio-group
        v-model="options.mode"
        buttons
        button-variant="outline-secondary"
        size="sm"
        name="buttons2"
        :options="modes"
      />
    </b-form-group>
    <b-form-checkbox
      v-model="options.hideFileName"
      :disabled="!enableFileNameHiding"
      class="mb-3"
    >
      {{ $t('kind.file.view.showName') }}
    </b-form-checkbox>
    <uploader
      :endpoint="endpoint"
      :max-filesize="$s('compose.Page.Attachments.MaxSize', 100)"
      :accepted-files="$s('compose.Page.Attachments.Mimetypes', ['*/*'])"
      @uploaded="appendAttachment"
    />
    <list-loader
      kind="page"
      enable-delete
      :namespace="namespace"
      :set.sync="options.attachments"
      mode="list"
      class="mt-2"
    />
  </b-tab>
</template>
<script>
import base from './base'
import Uploader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Uploader'
import ListLoader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/ListLoader'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  name: 'File',

  components: {
    Uploader,
    ListLoader,
  },

  extends: base,

  computed: {
    endpoint () {
      const { pageID } = this.page

      return this.$ComposeAPI.pageUploadEndpoint({
        namespaceID: this.namespace.namespaceID,
        pageID,
      })
    },

    modes () {
      return [
        { value: 'list', text: this.$t('kind.file.view.list') },
        { value: 'grid', text: this.$t('kind.file.view.grid') },
        { value: 'single', text: this.$t('kind.file.view.single') },
        { value: 'gallery', text: this.$t('kind.file.view.gallery') },
      ]
    },

    enableFileNameHiding () {
      return (this.options.mode === 'single') || (this.options.mode === 'gallery')
    },
  },

  methods: {
    appendAttachment ({ attachmentID } = {}) {
      this.options.attachments.push(attachmentID)
    },
  },
}
</script>
