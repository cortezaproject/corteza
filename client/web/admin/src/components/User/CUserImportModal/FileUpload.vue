<template>
  <div>
    <c-uploader
      :labels="{
        instructions: $t('import.uploadFilePlaceholder'),
      }"
      :endpoint="endpoint"
      :accepted-files="['application/zip']"
      @upload="onUploaded"
    />
  </div>
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
