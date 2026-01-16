<template>
  <wrap
    v-bind="$props"
    :scrollable-body="true"
    v-on="$listeners"
    @refreshBlock="refresh"
  >
    <div
      class="d-flex flex-column align-items-center h-100 overflow-hidden"
    >
      <span
        v-if="revisionsDisabledOnModule"
        class="my-auto"
      >
        {{ $t('errors.disabled-on-module') }}
      </span>

      <b-spinner
        v-else-if="isProcessing"
        class="my-auto"
      />

      <b-button
        v-else-if="!preloadRevisions && !loadedRevisions"
        variant="light"
        class="my-auto"
        @click="refresh()"
      >
        {{ $t('show-revisions', { revision: record.revision }) }}
      </b-button>

      <template
        v-else
      >
        <b-table-lite
          :items="revisions"
          :fields="columns"
          sticky-header
          hover
          details-td-class="bg-light p-0"
          responsive
          class="flex-fill mh-100 mb-0 w-100 rounded"
        >
          <template #cell(revision)="row">
            {{ row.item.values.revision }}
          </template>

          <template #cell(operation)="row">
            {{ $t(`operations.${row.item.values.operation}`) }}
          </template>

          <template #cell(user)="row">
            <field-viewer
              value-only
              :field="mockRevisionModule.findField('userID')"
              :record="row.item"
              :module="mockRevisionModule"
              :namespace="namespace"
            />
          </template>

          <template #cell(timestamp)="row">
            <field-viewer
              value-only
              :field="mockRevisionModule.findField('timestamp')"
              :record="row.item"
              :module="mockRevisionModule"
              :namespace="namespace"
            />
          </template>

          <template
            #cell(adt)="row"
          >
            <b-button
              v-if="row.item.meta.changes.length > 0"
              variant="outline-extra-light"
              class="d-flex align-items-center m-0 border-0 text-primary ml-auto"
              @click="row.toggleDetails"
            >
              {{ $t(`show-changes`, { count: row.item.meta.changes.length }) }}
              <font-awesome-icon
                :icon="row.detailsShowing ? 'chevron-up' : 'chevron-down'"
                class="ml-2"
              />
            </b-button>
          </template>

          <template #row-details="row">
            <div class="overflow-hidden">
              <b-table-simple
                table-variant="light"
                class="mb-0"
              >
                <b-thead>
                  <b-tr class="text-primary">
                    <b-th>{{ $t('changes.columns.field.label') }}</b-th>
                    <b-th>{{ $t('changes.columns.old-value.label') }}</b-th>
                    <b-th>{{ $t('changes.columns.new-value.label') }}</b-th>
                  </b-tr>
                </b-thead>

                <b-tbody>
                  <b-tr
                    v-for="(change) in row.item.meta.changes"
                    :key="change.key"
                  >
                    <b-td>
                      {{ change.label }}
                    </b-td>

                    <b-td>
                      <field-viewer
                        v-if="row.item.meta.oldRecord"
                        value-only
                        :field="module.findField(change.key)"
                        :record="row.item.meta.oldRecord"
                        :module="module"
                        :namespace="namespace"
                      />

                      <span v-else>
                        -
                      </span>
                    </b-td>

                    <b-td>
                      <field-viewer
                        v-if="row.item.meta.newRecord"
                        value-only
                        :field="module.findField(change.key)"
                        :record="row.item.meta.newRecord"
                        :module="module"
                        :namespace="namespace"
                      />

                      <span v-else>
                        -
                      </span>
                    </b-td>
                  </b-tr>
                </b-tbody>
              </b-table-simple>
            </div>
          </template>
        </b-table-lite>

        <div
          v-if="!revisions.length"
          class="position-absolute text-center mt-5 d-print-none"
          style="left: 0; right: 0; bottom: calc(50% - 33px);"
        >
          <p class="mb-0 mx-2">
            {{ $t('errors.no-revisions') }}
          </p>
        </div>
      </template>
    </div>
  </wrap>
</template>
<script>
import base from './base'
import { compose, NoID } from '@cortezaproject/corteza-js'
import users from 'corteza-webapp-compose/src/mixins/users'
import records from 'corteza-webapp-compose/src/mixins/records'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'

export default {
  i18nOptions: {
    namespaces: 'block',
    keyPrefix: 'recordRevisions.viewer',
  },

  components: {
    FieldViewer,
  },

  extends: base,

  mixins: [
    users,
    records,
  ],

  data () {
    return {
      /**
       * Last error that occurred
       * while loading the revisions
       */
      error: null,

      /**
       * When true, the revisions are being loaded
       * controled from refresh() method
       */
      processing: false,

      /**
       * Flag for if user clicked on show revisions button
       */
      loadedRevisions: false,

      mockRevisionModule: undefined,

      /**
       * List of revisions when loaded
       */
      revisions: [],

      /**
       * Revisions table fields
       *
       * Please note that table utilizes row-details feature
       * where changes are displayed
       */
      columns: [
        {
          key: 'revision',
          label: '#',
          thClass: 'border-top-0',
          class: 'text-left',
        },
        {
          key: 'operation',
          label: this.$t('revisions.columns.operation.label'),
          thClass: 'border-top-0',
        },
        {
          key: 'user',
          label: this.$t('revisions.columns.user.label'),
          thClass: 'border-top-0 text-right',
          tdClass: 'text-right',
        },
        {
          key: 'timestamp',
          label: this.$t('revisions.columns.timestamp.label'),
          thClass: 'border-top-0 text-right',
          tdClass: 'text-right',
        },
        {
          key: 'adt',
          label: '',
          thClass: 'border-top-0',
          class: 'text-nowrap text-right',
        },
      ],
    }
  },

  computed: {
    showHeader () {
      return !!this.block.title || !!this.block.description
    },

    revisionsDisabledOnModule () {
      return this.module ? !this.module.config.recordRevisions.enabled : false
    },

    preloadRevisions () {
      return this.options.preload
    },
  },

  watch: {
    'record.recordID': {
      immediate: true,
      handler () {
        if (this.preloadRevisions) {
          this.refresh()
        }
      },
    },

    options: {
      deep: true,
      handler () {
        this.refresh()
      },
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
    this.destroyEvents()
  },

  created () {
    this.refreshBlock(this.refresh)
  },

  mounted () {
    this.createEvents()
  },

  methods: {
    createEvents () {
      this.$root.$on('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$on('refetch-records', this.refresh)
    },

    refreshOnRelatedRecordsUpdate ({ moduleID } = {}) {
      if (this.module.moduleID === moduleID) {
        this.refresh()
      }
    },

    async loadRevisions () {
      this.revisions = []

      if (this.revisionsDisabledOnModule) {
        return
      }

      if (!this.record || this.record.recordID === NoID) {
        return
      }

      this.processing = true

      const fields = [
        { name: 'revision', kind: 'Number' },
        { name: 'changeID', kind: 'String' },
        { name: 'userID', kind: 'User' },
        { name: 'timestamp', kind: 'DateTime' },
        { name: 'operation', kind: 'String' },
      ]

      this.mockRevisionModule = new compose.Module({
        ...this.module,
        fields,
      })

      return this.block.fetch(this.$ComposeAPI, this.record, this.options.sortDirection).then(set => {
        this.revisions = set.map(r => {
          // Build old and new record values from all changes
          let oldOwnedBy = NoID
          let newOwnedBy = NoID
          const oldValues = {}
          const newValues = {}

          r.changes.forEach(c => {
            if (c.old !== undefined) {
              if (c.key === 'ownedBy') {
                oldOwnedBy = c.old[0]
              } else {
                oldValues[c.key] = c.old
              }
            }

            if (c.new !== undefined) {
              if (c.key === 'ownedBy') {
                newOwnedBy = c.new[0]
              } else {
                newValues[c.key] = c.new
              }
            }
          })

          // Create old and new records for this revision
          // Start with the current record's system fields, then apply changes
          const oldRecord = Object.keys(oldValues).length > 0
            ? new compose.Record(this.module, {
              ownedBy: oldOwnedBy,
              values: oldValues,
            })
            : null

          const newRecord = Object.keys(newValues).length > 0
            ? new compose.Record(this.module, {
              ownedBy: newOwnedBy,
              values: newValues,
            })
            : null

          // Changes now contain the key and label
          const changes = r.changes.map(c => {
            const field = this.module.findField(c.key)
            return {
              key: c.key,
              label: field ? field.label || field.name : c.key,
            }
          })

          return new compose.Record(this.mockRevisionModule, {
            recordID: r.resource,
            values: {
              revision: r.revision,
              changeID: r.changeID,
              operation: r.operation,
              timestamp: r.timestamp,
              userID: r.userID,
              recordID: r.recordID,
            },
            meta: {
              changes,
              oldRecord,
              newRecord,
            },
          })
        })

        // Collect all old/new records for fetching users and related records
        const allRecords = []
        this.revisions.forEach(rev => {
          if (rev.meta.oldRecord) allRecords.push(rev.meta.oldRecord)
          if (rev.meta.newRecord) allRecords.push(rev.meta.newRecord)
        })

        return Promise.all([
          this.fetchUsers(this.module.fields, this.revisions),
          this.fetchRecords(this.namespace.namespaceID, this.module.fields, allRecords),
        ])
      })
        .finally(() => {
          setTimeout(() => {
            this.processing = false
          }, 300)
        })
    },

    refresh () {
      this.loadedRevisions = true
      this.loadRevisions()
    },

    setDefaultValues () {
      this.error = null
      this.processing = false
      this.loadedRevisions = false
      this.revisions = []
      this.columns = []
    },

    destroyEvents () {
      this.$root.$off('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$off('refetch-records', this.refresh)
    },
  },
}
</script>
