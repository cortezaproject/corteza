<template>
  <b-card>
    <uploader
      class="uploader"
      :label="$t('import.uploadFilePlaceholder')"
      :endpoint="endpoint"
      :accepted-files="['application/zip']"
      @uploaded="onUploaded"
    />
  </b-card>
</template>

<script>
import Uploader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Uploader'

export default {
  i18nOptions: {
    namespaces: 'namespace',
  },

  components: {
    Uploader,
  },

  data () {
    return {
      session: null,
      sessionFile: null,
    }
  },

  computed: {
    endpoint () {
      return this.$ComposeAPI.namespaceImportInitEndpoint({})
    },

    canContinue () {
      return !!this.session
    },
  },

  methods: {
    onUploaded (e, f) {
      this.session = e
      this.sessionFile = f
      this.fileUploaded()
    },

    fileUploaded () {
      this.$emit('fileUploaded', {
        ...this.session || {},
      })
    },
  },
}
</script>
