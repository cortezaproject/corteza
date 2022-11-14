<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>
    <template v-else-if="roModule && contentField">
      <section
        v-if="canAddRecord"
        class="d-flex flex-column px-3 py-2"
      >
        <b-form-input
          v-if="titleField"
          v-model="newRecord.title"
          class="mb-2"
          :placeholder="$t('comment.titleInput')"
        />
        <b-form-textarea
          v-model.trim="newRecord.content"
          :value="true"
          :placeholder="$t('comment.contentInput')"
        />
        <b-button
          variant="primary"
          class="ml-auto mt-2 mb-2"
          :disabled="!isValid"
          @click="createNewRecord()"
        >
          {{ $t('comment.submit') }}
        </b-button>
      </section>
      <div
        v-if="sortableRecords.length && canAddRecord"
        class="border w-100 mb-3"
      />
      <section v-if="sortableRecords.length">
        <b-list-group class="px-3 py-2">
          <b-list-group-item
            v-for="record in sortableRecords"
            :key="record.recordID"
            class="p-0 pb-3 border-0"
          >
            <div class="d-flex flex-wrap border p-2">
              <div class="text-primary">
                {{ getAuthor(record.ownedBy) }}
              </div>
              <span class="ml-auto text-muted">
                {{ getFormattedDate((record || {}).updatedAt || (record || {}).createdAt) }}
              </span>
            </div>
            <div class="border p-3 d-flex flex-column">
              <field-viewer
                v-if="titleField && titleField.canReadRecordValue"
                class="mb-3 text-muted font-weight-bold"
                :field="titleField"
                :record="record"
                :namespace="namespace"
                value-only
              />
              <template v-else-if="!options.titleField" />
              <i
                v-else
                class="text-secondary h6"
              >{{ $t('field.noPermission') }}</i>
              <field-viewer
                v-if="contentField.canReadRecordValue"
                :field="contentField"
                :record="record"
                :namespace="namespace"
                value-only
              />
              <i
                v-else
                class="text-secondary h6"
              >{{ $t('field.noPermission') }}</i>
            </div>
          </b-list-group-item>
        </b-list-group>
      </section>
    </template>
  </wrap>
</template>
<script>
import { mapGetters } from 'vuex'
import base from './base'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import users from 'corteza-webapp-compose/src/mixins/users'
import { compose, NoID, fmt } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldViewer,
  },

  extends: base,

  mixins: [
    users,
  ],

  data () {
    return {
      processing: false,

      filter: {
        sort: '',
        filter: '',
      },

      records: [],
      newRecord: {
        title: '',
        content: '',
      },
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      pages: 'page/set',
      findByID: 'user/findByID',
    }),

    roModule () {
      return this.getModuleByID(this.moduleID)
    },

    moduleID () {
      return this.options.moduleID
    },

    titleField () {
      const { titleField } = this.options

      if (!titleField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === titleField)
    },

    contentField () {
      const { contentField } = this.options

      if (!contentField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === contentField)
    },

    referenceField () {
      const { referenceField } = this.options

      if (!referenceField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === referenceField) || {}
    },

    canAddRecord () {
      return this.roModule && this.roModule.canCreateRecord
    },

    isValid () {
      return !!this.newRecord.content && !this.isNewRecord
    },

    isNewRecord () {
      if (this.record) {
        return this.record.recordID === NoID
      }
      return false
    },

    sortableRecords () {
      if (this.options.sortDirection === 'asc') {
        return [...this.records].sort((a, b) => a.createdAt - b.createdAt)
      } else {
        return [...this.records].sort((a, b) => b.createdAt - a.createdAt)
      }
    },

    reference () {
      if (this.record) {
        return this.record.recordID !== NoID ? this.record.recordID : ''
      }
      return NoID
    },
  },

  watch: {
    'record.recordID': {
      immediate: true,
      handler () {
        this.reloadRecords()
      },
    },
  },

  methods: {
    getFormattedDate (date) {
      return fmt.fullDateTime(date)
    },

    getAuthor (userID) {
      const user = this.findByID(userID) || {}
      return user.name || user.handle || user.email || ''
    },

    fetchUsers () {
      const userListID = this.records.map(r => {
        return r.ownedBy
      }).filter((x, i, r) => r.indexOf(x) === i)
      this.$store.dispatch('user/fetchUsers', userListID)
    },

    reloadRecords () {
      if (!this.options.moduleID) {
      // Make sure block is properly configured
        throw Error(this.$t('record.moduleOrPageNotSet'))
      }
      if (this.roModule && this.contentField) {
        this.processing = true
        this.fetchRecords(this.roModule, this.expandFilter())
          .then(rr => {
            this.records = rr
            this.fetchUsers()
          })
          .catch(e => {
            console.error(e)
          })
          .finally(() => {
            this.processing = false
          })
      }
    },

    createNewRecord () {
      // cannot create new record if content field is empty
      if (this.newRecord.content.length) {
        const record = {}
        record.values = []
        if (this.contentField) {
          record.values.push({
            name: this.contentField.name,
            value: this.newRecord.content,
          })
        }
        if (this.referenceField) {
          record.values.push({
            name: this.referenceField.name,
            value: this.reference,
          })
        }
        if (this.titleField) {
          record.values.push({
            name: this.titleField.name,
            value: this.newRecord.title,
          })
        }
        record.moduleID = this.options.moduleID
        record.namespaceID = this.roModule.namespaceID
        this.$ComposeAPI.recordCreate(record).then(rec => {
          // clean the input and reload data
          this.newRecord.title = ''
          this.newRecord.content = ''
          this.reloadRecords()
        })
          .catch(this.toastErrorHandler(this.$t('notification:record.createFailed')))
      }
    },

    expandFilter () {
      /* eslint-disable no-template-curly-in-string */
      if (!this.record) {
        // If there is no current record and we are using recordID/ownerID variable in (pre)filter
        // we should disable the block
        if ((this.options.filter || '').includes('${record')) {
          throw Error(this.$t('record.invalidRecordVar'))
        }

        if ((this.options.filter || '').includes('${ownerID}')) {
          throw Error(this.$t('record.invalidOwnerVar'))
        }
      }

      if (this.options.filter) {
        return evaluatePrefilter(this.options.filter, {
          record: this.record,
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).userID || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      }

      return ''
    },
    /**
     * Fetches group of records using configured options & module
     *
     * @param {Compose}           api Compose API client
     * @param {Module}            module Module to use for assembling API request & casting results
     * @param {String}            filter Filter records
     * @returns {Promise<Record[]>}
     */
    async fetchRecords (module, filter) {
      if (module.moduleID !== this.options.moduleID) {
        throw Error('Module incompatible, module mismatch')
      }
      if (this.referenceField) {
        if (filter.length) {
          filter += ' AND '
        }
        filter += `${this.referenceField.name} = '${this.reference}' `
      }
      const { positionField: sort } = this.options
      const { moduleID, namespaceID } = module

      const params = {
        namespaceID,
        moduleID,
        query: filter,
        sort,
      }

      return this.$ComposeAPI
        .recordList(params)
        .then(({ set }) => set.map(r => Object.freeze(new compose.Record(module, r))))
    },
  },
}
</script>

<style lang="scss" scoped>
.grab {
  cursor: grab;
}
</style>
