<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form-group>
      <c-uploader
        class="uploader"
        :labels="{
          instructions: $t('import.uploadFilePlaceholder'),
        }"
        :endpoint="endpoint"
        :accepted-files="['application/zip']"
        @upload="onUploaded"
      />
    </b-form-group>
  </b-card>
</template>

<script>
import CUploader from 'corteza-webapp-admin/src/components/CUploader'

export default {
  i18nOptions: {
    namespaces: 'system.users',
  },

  name: 'CFileUpload',

  components: {
    CUploader,
  },

  data () {
    return {
      session: null,
      sessionFile: null,
    }
  },

  computed: {
    endpoint () {
      return this.$SystemAPI.userImportEndpoint({})
    },

    canContinue () {
      return !!this.session
    },
  },

  methods: {
    onUploaded () {
      this.$emit('imported', {})
    },
  },
}
</script>

<style lang="scss" scoped>
.uploader {
  height: 130px;
}
</style>
