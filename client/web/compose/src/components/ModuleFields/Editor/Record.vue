<template>
  <b-form-group
    :label-cols-md="horizontal && '5'"
    :label-cols-xl="horizontal && '4'"
    :content-cols-md="horizontal && '7'"
    :content-cols-xl="horizontal && '8'"
    :class="formGroupStyleClasses"
  >
    <template
      #label
    >
      <div
        v-if="!valueOnly"
        class="d-flex align-items-center text-primary p-0"
      >
        <span
          :title="label"
          class="d-inline-block mw-100"
        >
          {{ label }}
        </span>

        <c-hint :tooltip="hint" />

        <slot name="tools" />
      </div>
      <div
        class="small text-muted"
        :class="{ 'mb-1': description }"
      >
        {{ description }}
      </div>
    </template>

    <multi
      v-if="field.isMulti"
      :value.sync="value"
      :errors="errors"
      :single-input="field.options.selectType !== 'each'"
      :show-list="field.options.selectType !== 'multiple'"
    >
      <template #single>
        <b-input-group class="d-flex w-100">
          <c-input-select
            v-if="field.options.selectType === 'multiple'"
            v-model="multipleSelected"
            :options="options"
            :get-option-key="getOptionKey"
            :get-option-label="getOptionLabel"
            :disabled="!module"
            :loading="processing"
            :clearable="false"
            :filterable="false"
            :searchable="searchable"
            :selectable="isSelectable"
            :placeholder="placeholder"
            multiple
            @search="search"
          >
            <pagination
              v-if="showPagination"
              slot="list-footer"
              :has-prev-page="hasPrevPage"
              :has-next-page="hasNextPage"
              @prev="goToPage(false)"
              @next="goToPage(true)"
            />
          </c-input-select>

          <c-input-select
            v-else
            ref="singleSelect"
            :options="options"
            :get-option-key="getOptionKey"
            :get-option-label="getOptionLabel"
            :disabled="!module"
            :loading="processing"
            :clearable="false"
            :filterable="false"
            :searchable="searchable"
            :selectable="isSelectable"
            :placeholder="placeholder"
            @input="selectChange($event)"
            @search="search"
          >
            <pagination
              v-if="showPagination"
              slot="list-footer"
              :has-prev-page="hasPrevPage"
              :has-next-page="hasNextPage"
              @prev="goToPage(false)"
              @next="goToPage(true)"
            />
          </c-input-select>

          <b-input-group-append v-if="canAddRecordThroughSelectField">
            <b-button
              v-b-tooltip.hover="{ title: $t('kind.record.tooltip.addRecord'), container: '#body' }"
              variant="light"
              class="d-flex align-items-center"
              @click="addRecordThroughRecordSelectField()"
            >
              <font-awesome-icon
                :icon="['fas', 'plus']"
                class="text-primary"
              />
            </b-button>
          </b-input-group-append>
        </b-input-group>
      </template>
      <template #default="ctx">
        <b-input-group
          v-if="field.options.selectType === 'each'"
          class="d-flex w-100"
        >
          <c-input-select
            :options="options"
            :get-option-key="getOptionKey"
            :get-option-label="getOptionLabel"
            :disabled="!module"
            :loading="processing"
            :clearable="false"
            :filterable="false"
            :searchable="searchable"
            :selectable="isSelectable"
            :placeholder="placeholder"
            :value="getRecord(ctx.index)"
            @input="setRecord($event, ctx.index)"
            @search="search"
          >
            <pagination
              v-if="showPagination"
              slot="list-footer"
              :has-prev-page="hasPrevPage"
              :has-next-page="hasNextPage"
              @prev="goToPage(false)"
              @next="goToPage(true)"
            />
          </c-input-select>

          <b-input-group-append v-if="canAddRecordThroughSelectField">
            <b-button
              v-b-tooltip.hover="{ title: $t('kind.record.tooltip.addRecord'), container: '#body' }"
              variant="light"
              class="d-flex align-items-center"
              @click="addRecordThroughRecordSelectField()"
            >
              <font-awesome-icon
                :icon="['fas', 'plus']"
                class="text-primary"
              />
            </b-button>
          </b-input-group-append>
        </b-input-group>
        <b-spinner
          v-else-if="resolving"
          variant="primary"
          small
        />

        <span v-else>
          {{ getOptionLabel(multipleSelected[ctx.index]) }}
        </span>
      </template>
    </multi>

    <template
      v-else
    >
      <b-input-group>
        <c-input-select
          v-model="selected"
          :options="options"
          :get-option-key="getOptionKey"
          :get-option-label="getOptionLabel"
          :disabled="!module"
          :loading="processing"
          :placeholder="placeholder"
          :filterable="false"
          :searchable="searchable"
          :selectable="isSelectable"
          @search="search"
        >
          <pagination
            v-if="showPagination"
            slot="list-footer"
            :has-prev-page="hasPrevPage"
            :has-next-page="hasNextPage"
            @prev="goToPage(false)"
            @next="goToPage(true)"
          />
        </c-input-select>

        <b-input-group-append v-if="canAddRecordThroughSelectField">
          <b-button
            v-b-tooltip.hover="{ title: $t('kind.record.tooltip.addRecord'), container: '#body' }"
            variant="light"
            class="d-flex align-items-center"
            @click="addRecordThroughRecordSelectField()"
          >
            <font-awesome-icon
              :icon="['fas', 'plus']"
              class="text-primary"
            />
          </b-button>
        </b-input-group-append>
      </b-input-group>

      <errors :errors="errors" />
    </template>
  </b-form-group>
</template>
<script>
import base from './base'
import { debounce } from 'lodash'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { mapActions, mapGetters } from 'vuex'
import { queryToFilter, evaluatePrefilter, isFieldInFilter } from 'corteza-webapp-compose/src/lib/record-filter'
import Pagination from '../Common/Pagination.vue'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    Pagination,
  },

  extends: base,

  data () {
    return {
      processing: false,
      resolving: false,

      query: '',

      records: [],

      recordValues: {},

      filter: {
        query: '',
        sort: '',
        limit: 10,
        pageCursor: '',
        prevPage: '',
        nextPage: '',
      },
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      findUserByID: 'user/findByID',
      findRecordsByIDs: 'record/findByIDs',
      pages: 'page/set',
    }),

    options () {
      return this.records
    },

    module () {
      if (this.field.options.moduleID !== NoID) {
        return this.getModuleByID(this.field.options.moduleID)
      } else {
        return undefined
      }
    },

    searchable () {
      return !this.field.options.recordLabelField
    },

    placeholder () {
      return this.searchable ? this.$t('kind.record.suggestionPlaceholder') : this.$t('kind.select.placeholder')
    },

    multipleSelected: {
      get () {
        return this.value
      },

      set (value) {
        this.value = value.map(v => {
          return typeof v === 'string' ? v : v.recordID
        })
      },
    },

    selected: {
      get () {
        return this.getRecord()
      },

      set (value) {
        this.setRecord(value)
      },
    },

    showPagination () {
      return this.hasPrevPage || this.hasNextPage
    },

    hasPrevPage () {
      return !!this.filter.prevPage
    },

    hasNextPage () {
      return !!this.filter.nextPage
    },

    canAddRecordThroughSelectField () {
      if (!this.extraOptions.recordSelectorShowAddRecordButton || this.module === undefined) return

      return !!this.getRecordSelectorPage().page.pageID && this.module.canCreateRecord
    },
  },

  watch: {
    'filter.pageCursor': {
      handler (pageCursor) {
        if (pageCursor) {
          this.fetchPrefiltered(this.filter)
        }
      },
    },
  },

  created () {
    this.loadLatest()
    if (this.value) {
      this.resolving = true
      this.formatRecordValues(this.value).finally(() => {
        this.resolving = false
      })
    }
  },

  beforeDestroy () {
    this.setDefaultValues()
    this.destroyEvents()
  },

  mounted () {
    this.createEvents()
  },

  methods: {
    ...mapActions({
      findModuleByID: 'module/findByID',
      resolveUsers: 'user/resolveUsers',
      resolveRecords: 'record/resolveRecords',
      updateRecords: 'record/updateRecords',
    }),

    createEvents () {
      this.$root.$on('record-field-change', this.refetchOnPrefilterValueChange)
      this.$root.$on('module-records-updated', this.refreshOnRelatedRecordsUpdate)
    },

    refetchOnPrefilterValueChange ({ fieldName }) {
      const { prefilter } = this.field.options

      if (isFieldInFilter(fieldName, prefilter)) {
        const namespaceID = this.namespace.namespaceID
        const moduleID = this.field.options.moduleID
        this.fetchPrefiltered({ namespaceID, moduleID })
      }
    },

    getRecord (index = undefined) {
      return index !== undefined ? this.value[index] : this.value
    },

    setRecord ({ recordID } = {}, index = undefined) {
      const crtValue = index !== undefined ? this.value[index] : this.value

      if (recordID !== crtValue) {
        if (recordID) {
          if (index !== undefined) {
            this.value.splice(index, 1, recordID)
          } else {
            this.value = recordID
          }
        } else {
          if (index !== undefined) {
            this.value.splice(index, 1)
          } else {
            this.value = undefined
          }
        }
      }
    },

    isSelectable ({ recordID } = {}) {
      if (!recordID) {
        return false
      }

      if (this.field.isMulti) {
        return !this.field.options.isUniqueMultiValue || !this.value.includes(recordID)
      } else {
        return this.value !== recordID
      }
    },

    search: debounce(function (query = '') {
      if (query !== this.query) {
        this.query = query
        this.filter.pageCursor = undefined
      }

      const { limit, pageCursor } = this.filter
      const namespaceID = this.namespace.namespaceID
      const moduleID = this.field.options.moduleID

      if (moduleID && moduleID !== NoID) {
        // Determine what fields to use for searching
        // Default to label field
        let qf = this.field.options.queryFields
        if (!qf || qf.length === 0) {
          qf = [this.field.options.labelField]
        }

        if (query.length > 0) {
          const fields = qf.map(f => this.module.fields.find(({ name }) => name === f))
          query = queryToFilter(query, '', fields)
        }

        this.fetchPrefiltered({ namespaceID, moduleID, query, sort: this.sortString(), limit, pageCursor })
      }
    }, 300),

    loadLatest () {
      const namespaceID = this.namespace.namespaceID
      const moduleID = this.field.options.moduleID
      const { limit } = this.filter
      if (moduleID && moduleID !== NoID) {
        this.fetchPrefiltered({ namespaceID, moduleID, limit })
      }
    },

    fetchPrefiltered (q = this.filter) {
      this.processing = true

      // Support prefilters
      let { query = '' } = q
      if (this.field.options.prefilter) {
        const pf = evaluatePrefilter(this.field.options.prefilter, {
          record: this.record,
          user: this.$auth.user || {},
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
        if (query) {
          query = `(${pf}) AND (${query})`
        } else {
          query = pf
        }
      }

      if (q.pageCursor) {
        q.sort = ''
      }

      this.$ComposeAPI.recordList({ ...q, query })
        .then(({ filter, set }) => {
          this.filter = { ...this.filter, ...filter }
          this.filter.nextPage = filter.nextPage
          this.filter.prevPage = filter.prevPage

          this.updateRecords(set)

          return this.formatRecordValues(set.map(({ recordID }) => recordID)).then(() => {
            this.records = set.map(r => new compose.Record(this.module, r))
          })
        }).finally(() => {
          this.processing = false
        })
    },

    sortString () {
      return [this.field.options.labelField].filter(f => !!f).join(', ')
    },

    async formatRecordValues (recordIDs) {
      recordIDs = Array.isArray(recordIDs) ? recordIDs : [recordIDs].filter(v => v) || []
      const { namespaceID = NoID } = this.namespace
      const { moduleID = NoID, labelField, recordLabelField } = this.field.options

      if (!recordIDs.length || [moduleID, namespaceID].includes(NoID) || !labelField) {
        return
      }

      return this.findModuleByID({ namespace: this.namespace, moduleID }).then(async module => {
        const relatedField = module.fields.find(({ name }) => name === labelField)
        let records = this.findRecordsByIDs(recordIDs).map(r => new compose.Record(module, r))
        const mappedIDs = {}

        if (relatedField.kind === 'Record' && recordLabelField) {
          const relatedModule = await this.findModuleByID({ namespaceID, moduleID: relatedField.options.moduleID })
          const relatedRecordIDs = new Set()

          records.forEach(r => {
            const recordValue = relatedField.isMulti ? r.values[relatedField.name] : [r.values[relatedField.name]]
            recordValue.forEach(rID => relatedRecordIDs.add(rID))
          })
          await this.resolveRecords({ namespaceID, moduleID: relatedModule.moduleID, recordIDs: [...relatedRecordIDs] })

          const relatedLabelField = relatedModule.fields.find(({ name }) => name === recordLabelField)

          for (let r of await this.findRecordsByIDs([...relatedRecordIDs])) {
            r = new compose.Record(relatedModule, r)
            let relatedRecordValue = relatedLabelField.isMulti ? r.values[relatedLabelField.name] : [r.values[relatedLabelField.name]]

            if (relatedLabelField.kind === 'User') {
              await this.resolveUsers(relatedRecordValue)
              relatedRecordValue = relatedRecordValue.map(v => relatedLabelField.formatter(this.findUserByID(v)))
            }

            mappedIDs[r.recordID] = relatedRecordValue.join(relatedLabelField.options.multiDelimiter)
            relatedRecordIDs.clear()
          }
        } else if (relatedField.kind === 'User') {
          this.processing = true

          const relatedUserIDs = new Set()
          records.forEach(r => {
            const recordValue = relatedField.isMulti ? r.values[relatedField.name] : [r.values[relatedField.name]]
            recordValue.forEach(uID => relatedUserIDs.add(uID))
          })

          await this.resolveUsers([...relatedUserIDs])
        } else if (records.length === 0) {
          await this.resolveRecords({ namespaceID, moduleID, recordIDs: [...recordIDs] })
          records = this.findRecordsByIDs(recordIDs).map(r => new compose.Record(module, r))
        }

        records.forEach(record => {
          let recordValue = relatedField.isMulti ? record.values[relatedField.name] : [record.values[relatedField.name]]

          if (relatedField.kind === 'User') {
            recordValue = recordValue.map(v => relatedField.formatter(this.findUserByID(v)))
          } else if (relatedField.kind === 'Record' && recordLabelField) {
            recordValue = recordValue.map(v => mappedIDs[v])
          }

          this.$set(this.recordValues, record.recordID, recordValue.join(relatedField.options.multiDelimiter))
        })
      })
    },

    addRecordThroughRecordSelectField () {
      const { page } = this.getRecordSelectorPage()

      if (page === undefined) return

      const { pageID } = page
      const { recordSelectorAddRecordDisplayOption } = this.extraOptions

      const route = {
        name: 'page.record.create',
        params: { pageID, edit: true },
      }

      if (recordSelectorAddRecordDisplayOption === 'modal') {
        this.$root.$emit('show-record-modal', {
          recordID: NoID,
          recordPageID: pageID,
          edit: true,
        })
      } else if (recordSelectorAddRecordDisplayOption === 'newTab') {
        window.open(this.$router.resolve(route).href)
      } else {
        this.$router.push(route)
      }
    },

    refreshOnRelatedRecordsUpdate () {
      const { page } = this.getRecordSelectorPage()
      if (page === undefined || this.module === undefined) return

      if (page.pageID !== this.$route.params.pageID) {
        this.loadLatest()
      }
    },

    getRecordSelectorPage () {
      const recordFieldModuleID = this.field.options.moduleID

      if (!recordFieldModuleID) return

      const recordFieldPage = this.pages.find(p => p.moduleID === recordFieldModuleID)

      if (!recordFieldPage) return

      return {
        page: recordFieldPage,
      }
    },

    selectChange ({ recordID } = {}) {
      if (!recordID) return

      this.value.push(recordID)

      // reset singleSelect value for better value presentation
      if (this.$refs.singleSelect) {
        this.$refs.singleSelect._data._value = undefined
      }
    },

    goToPage (next = true) {
      this.filter.pageCursor = next ? this.filter.nextPage : this.filter.prevPage
    },

    getOptionKey (value) {
      if (!value) return

      return typeof value === 'string' ? value : value.recordID
    },

    getOptionLabel (value) {
      if (!value) {
        return ''
      }

      const recordID = typeof value === 'string' ? value : value.recordID

      return this.recordValues[recordID] || recordID
    },

    destroyEvents () {
      this.$root.$off('record-field-change', this.refetchOnPrefilterValueChange)
      this.$root.$off('module-records-updated', this.refreshOnRelatedRecordsUpdate)
    },

    setDefaultValues () {
      this.processing = false
      this.resolving = false
      this.query = ''
      this.records = []
      this.recordValues = {}
      this.filter = {}
    },
  },
}
</script>
