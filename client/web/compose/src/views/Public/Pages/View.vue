<template>
  <div
    v-if="!!page"
    class="d-flex w-100 overflow-hidden"
  >
    <div
      class="flex-grow-1 overflow-auto d-flex px-2 w-100"
    >
      <router-view
        v-if="recordID || isRecordCreatePage"
        :namespace="namespace"
        :module="module"
        :page="page"
        class="flex-grow-1 overflow-auto d-flex flex-column"
      />

      <grid
        v-else
        :namespace="namespace"
        :module="module"
        :page="page"
      />
    </div>

    <attachment-modal />
  </div>
</template>
<script>
import Grid from 'corteza-webapp-compose/src/components/Public/Page/Grid'
import AttachmentModal from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Modal'
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  components: {
    Grid,
    AttachmentModal,
  },

  props: {
    namespace: { // via router-view
      type: compose.Namespace,
      required: true,
    },

    page: { // via route-view
      type: compose.Page,
      required: true,
    },

    // We're using recordID to check if we need to display router-view or grid component
    recordID: {
      type: String,
      default: '',
    },
  },

  computed: {
    isRecordCreatePage () {
      return this.$route.name === 'page.record.create'
    },

    module () {
      if (this.page.moduleID && this.page.moduleID !== NoID) {
        return this.$store.getters['module/getByID'](this.page.moduleID)
      }

      return undefined
    },
  },

  watch: {
    'page.title': {
      immediate: true,
      handler (title) {
        document.title = [title, this.namespace.name, this.$t('general:label.app-name.public')].filter(v => v).join(' | ')
      },
    },
  },
}
</script>
