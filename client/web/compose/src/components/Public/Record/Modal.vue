<template>
  <b-modal
    id="record-modal"
    v-model="showRecordModal"
    scrollable
    body-class="p-0"
    footer-class="p-0"
    size="xl"
    @hidden="hideModal"
  >
    <template #modal-title>
      <portal-target
        name="record-modal-header"
      />
    </template>

    <view-record
      :namespace="namespace"
      :page="page"
      :module="module"
      :record-i-d="recordID"
      :show-record-modal="showRecordModal"
    />

    <template #modal-footer>
      <portal-target
        name="record-modal-footer"
        class="w-100 m-0"
      />
    </template>
  </b-modal>
</template>

<script>
import record from 'corteza-webapp-compose/src/mixins/record'
import { compose } from '@cortezaproject/corteza-js'
import ViewRecord from 'corteza-webapp-compose/src/views/Public/Pages/Records/View'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'RecordModal',

  components: {
    ViewRecord,
  },

  mixins: [
    record,
  ],

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },
  },

  data () {
    return {
      showRecordModal: false,
      recordID: null,
      module: null,
      page: null,
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      getPageByID: 'page/getByID',
    }),
  },

  created () {
    this.$root.$on('show-record-modal', this.showModal)
    this.showModal(this.$route.query)
  },

  beforeDestroy () {
    this.$root.$off('show-record-modal')
  },

  methods: {
    showModal ({ recordID, moduleID, recordPageID }) {
      if (recordID && moduleID && recordPageID) {
        this.recordID = recordID
        this.module = this.getModuleByID(moduleID)
        this.page = this.getPageByID(recordPageID)

        this.showRecordModal = true

        // persist the modal in the url
        this.$router.replace({
          query: {
            recordID,
            moduleID,
            recordPageID: this.page.pageID,
          },
        })
      }
    },

    hideModal () {
      this.$router.replace({ query: { } })
      this.recordID = null
      this.module = null
      this.page = null
    },
  },
}

</script>

<style>
#record-modal .modal-dialog {
  height: 100%;
  max-width: 90vw;
}
</style>
