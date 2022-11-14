<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form-group>
      <slot name="uploadLabel">
        <label>{{ $t('import.uploadFileLabel') }}</label>
      </slot>

      <uploader
        class="uploader"
        :label="$t('import.uploadFilePlaceholder')"
        :endpoint="endpoint"
        :accepted-files="['application/zip']"
        @uploaded="onUploaded"
      />
    </b-form-group>
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

<style lang="scss" scoped>
.uploader {
  height: 130px;
}
</style>
