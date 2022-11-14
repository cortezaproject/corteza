<template>
  <b-container
    class="d-flex flex-column p-3"
  >
    <request-editor
      v-if="kind"
      :kind="kind"
      @submit="onSubmit"
    />
  </b-container>
</template>

<script>
import RequestEditor from 'corteza-webapp-privacy/src/components/Requests/Editor'

export default {
  name: 'RequestView',

  i18nOptions: {
    namespaces: 'request',
    keyPrefix: 'create',
  },

  components: {
    RequestEditor,
  },

  props: {
    kind: {
      type: String,
      required: true,
    },
  },

  data () {
    return {
      processing: false,
    }
  },

  methods: {
    onSubmit ({ kind, payload }) {
      this.processing = true

      payload = [payload]

      return this.$SystemAPI.dataPrivacyRequestCreate({ kind, payload })
        .then(({ requestID, kind } = {}) => {
          this.$router.push({ name: 'request.view', params: { requestID, kind } })
        })
        .catch(this.toastErrorHandler(this.$t('notification:list.load.error')))
        .finally(() => {
          this.processing = false
        })
    },
  },
}
</script>
