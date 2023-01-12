<template>
  <b-modal
    id="record-modal"
    v-model="showModal"
    scrollable
    dialog-class="h-100 mw-90"
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
      v-if="showModal"
      :namespace="namespace"
      :page="page"
      :module="module"
      :record-i-d="recordID"
      :show-record-modal="showModal"
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
      showModal: false,
      recordID: undefined,
      module: undefined,
      page: undefined,
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      getPageByID: 'page/getByID',
    }),
  },

  watch: {
    '$route.query.recordID': {
      immediate: true,
      handler (recordID, oldRecordID) {
        const { recordPageID } = this.$route.query

        if (!recordID) {
          this.showModal = false
          return
        }

        if (this.showModal && (recordID !== oldRecordID)) {
          this.showModal = false

          setTimeout(() => {
            this.$router.push({
              query: {
                ...this.$route.query,
                recordID,
                recordPageID,
              },
            })
          }, 300)

          return
        }

        setTimeout(() => {
          this.loadModal({ recordID, recordPageID })
        }, 100)
      },
    },
  },

  created () {
    this.$root.$on('show-record-modal', ({ recordID, recordPageID }) => {
      this.$router.push({
        query: {
          ...this.$route.query,
          recordID,
          recordPageID,
        },
      })
    })
  },

  beforeDestroy () {
    this.$root.$off('show-record-modal')
  },

  methods: {
    loadModal ({ recordID, recordPageID }) {
      if (recordID && recordPageID) {
        this.recordID = recordID
        this.page = this.getPageByID(recordPageID)

        if (this.page) {
          this.module = this.getModuleByID(this.page.moduleID)
          this.showModal = true
        }
      }
    },

    hideModal () {
      this.$router.push({
        query: {
          ...this.$route.query,
          recordID: undefined,
          moduleID: undefined,
          recordPageID: undefined,
        },
      })
    },
  },
}

</script>

<style lang="scss" scoped>
.mw-90 {
  max-width: 90vw;
}
</style>
