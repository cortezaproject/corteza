<template>
  <b-spinner
    v-if="processing"
    variant="primary"
    small
  />

  <div
    v-else
    @click.stop
  >
    <span
      v-for="(v, index) of formattedValue"
      :key="v.id || index"
      :class="{ 'd-block': field.options.multiDelimiter === '\n', 'mt-1': field.options.multiDelimiter === '\n' && index !== 0 }"
    >
      <template v-if="disableClick">
        <field-viewer
          v-if="v.record"
          :field="labelField"
          :record="v.record"
          :namespace="namespace"
          disable-click
          value-only
        />
        <template v-else>
          {{ v.id }}
        </template>
        {{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}
      </template>

      <a
        v-else-if="['modal', 'newTab'].includes(extraOptions.recordSelectorDisplayOption)"
        :href="resolveHref(v.to)"
        :class="{ 'text-decoration-none default-cursor': !v.to}"
        @click="(e) => onRecordSelectorClick(e, v.to)"
      >
        <field-viewer
          v-if="v.record"
          :field="labelField"
          :record="v.record"
          :namespace="namespace"
          disable-click
          value-only
        />
        <template v-else>
          {{ v.id }}
        </template>
        {{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}
      </a>

      <router-link
        v-else
        :to="v.to"
        :class="{ 'text-decoration-none default-cursor': !v.to}"
      >
        <field-viewer
          v-if="v.record"
          :field="labelField"
          :record="v.record"
          :namespace="namespace"
          disable-click
          value-only
        />
        <template v-else>
          {{ v.id }}
        </template>
        {{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}
      </router-link>
    </span>
  </div>
</template>

<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import records from 'corteza-webapp-compose/src/mixins/records'
import users from 'corteza-webapp-compose/src/mixins/users'
import { mapGetters } from 'vuex'
import base from './base'

export default {
  components: {
    FieldViewer: () => import('corteza-webapp-compose/src/components/ModuleFields/Viewer'),
  },

  extends: base,

  mixins: [
    users,
    records,
  ],

  data () {
    return {
      processing: false,

      recordValues: {},
    }
  },

  computed: {
    ...mapGetters({
      pages: 'page/set',
      getModuleByID: 'module/getByID',
      findRecordByID: 'record/findByID',
      findRecordsByIDs: 'record/findByIDs',
    }),

    formattedValue () {
      const value = Array.isArray(this.value) ? this.value : [this.value]
      return value.map(recordID => {
        let record = this.findRecordByID(recordID, this.field.options.moduleID)

        if (record) {
          record = new compose.Record(this.recordModule, record)
        }

        return {
          id: recordID,
          to: this.linkToRecord(recordID),
          record,
        }
      })
    },

    recordPage () {
      return this.pages.find(p => p.moduleID === this.field.options.moduleID)
    },

    recordModule () {
      if (!this.field.options.moduleID) {
        return undefined
      }

      return this.getModuleByID(this.field.options.moduleID)
    },

    labelField () {
      if (!this.field.options.labelField || !this.recordModule) {
        return undefined
      }

      let labelField = this.recordModule.fields.find(({ name }) => name === this.field.options.labelField)

      if (!labelField) {
        return undefined
      }

      labelField = compose.ModuleFieldMaker(labelField)

      if (labelField.kind === 'Record' && this.field.options.recordLabelField) {
        labelField.options.labelField = this.field.options.recordLabelField
      }

      return labelField
    },

    relatedModule () {
      if (!this.labelField.options.moduleID) {
        return undefined
      }

      return this.getModuleByID(this.labelField.options.moduleID)
    },

    relatedLabelField () {
      if (!this.field.options.recordLabelField || !this.relatedModule) {
        return undefined
      }

      return this.relatedModule.fields.find(({ name }) => name === this.field.options.recordLabelField)
    },
  },

  mounted () {
    this.formatRecordValues(this.value)
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {

    linkToRecord (recordID) {
      if (!this.recordPage || !recordID) {
        return ''
      }

      return {
        name: 'page.record',
        params: {
          pageID: this.recordPage.pageID,
          recordID,
        },
      }
    },

    resolveHref (route) {
      if (!route) {
        return '#'
      }
      return this.$router.resolve(route).href
    },

    formatRecordValues (recordIDs) {
      recordIDs = Array.isArray(recordIDs) ? recordIDs : [recordIDs].filter(v => v) || []
      const { namespaceID = NoID } = this.namespace
      const { moduleID = NoID, recordLabelField } = this.field.options

      if (!recordIDs.length || [moduleID, namespaceID].includes(NoID) || !this.labelField || !this.recordModule) {
        return
      }

      const stopProcessing = () => {
        setTimeout(() => {
          this.processing = false
        }, 300)
      }

      const records = this.findRecordsByIDs(recordIDs, this.field.options.moduleID).map(r => new compose.Record(this.recordModule, r))

      if (this.labelField.kind === 'Record' && recordLabelField) {
        this.processing = true

        return this.fetchRecords(namespaceID, [this.labelField], records).finally(stopProcessing)
      } else if (this.labelField.kind === 'User') {
        this.processing = true

        return this.fetchUsers([this.labelField], records).finally(stopProcessing)
      }
    },

    onRecordSelectorClick (e, route) {
      // Let the browser handle ctrl/cmd/shift/middle-click natively so users
      // can open the record page in a new tab via standard browser conventions.
      if (e.ctrlKey || e.metaKey || e.shiftKey || e.button === 1) {
        return
      }

      e.preventDefault()

      if (!route) {
        return
      }

      if (this.extraOptions.recordSelectorDisplayOption === 'modal' || this.inModal) {
        if (route.params && route.params.recordID && route.params.pageID) {
          this.$root.$emit('show-record-modal', {
            recordID: route.params.recordID,
            recordPageID: route.params.pageID,
          })
        }
      } else if (this.extraOptions.recordSelectorDisplayOption === 'newTab') {
        const resolved = this.$router.resolve(route)
        if (resolved && resolved.href) {
          window.open(resolved.href, '_blank')
        }
      }
    },

    setDefaultValues () {
      this.processing = false
      this.recordValues = {}
    },
  },
}
</script>

<style lang="scss" scoped>
.default-cursor {
  cursor: default;
}
</style>
