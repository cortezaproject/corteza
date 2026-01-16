<template>
  <div class="d-flex flex-column gap-1">
    <p
      v-if="!!message"
      class="text-break"
      v-html="message"
    />

    <b-form-group
      :label="label"
      label-class="text-primary"
    >
      <c-input-select
        v-model="value"
        :options="options"
        :get-option-key="r => r.recordID"
        :loading="processing"
        append-to-body
        option-value="recordID"
        :placeholder="placeholder"
        :filterable="false"
        :reduce="r => r.recordID"
        class="w-100"
        @search="search"
      >
        <template #list-footer>
          <c-pagination
            v-if="showPagination"
            :has-prev-page="hasPrevPage"
            :has-next-page="hasNextPage"
            @prev="goToPage(false)"
            @next="goToPage(true)"
          />
        </template>
      </c-input-select>
    </b-form-group>

    <b-button
      :disabled="loading"
      variant="primary"
      class="ml-auto"
      @click="$emit('submit', { value: encodeValue() })"
    >
      {{ pVal('buttonLabel', 'Submit') }}
    </b-button>
  </div>
</template>

<script lang="js">
import base from './base.vue'
import CPagination from '../common/CPagination.vue'
import CInputSelect from '../../input/CInputSelect.vue'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { debounce } from 'lodash'
import axios from 'axios'

export default {
  name: 'CPromptComposeRecordPicker',

  components: {
    CInputSelect,
    CPagination,
  },

  extends: base,

  data () {
    return {
      processing: false,
      cancelRequest: null,
      query: '',
      filter: {
        query: '',
        sort: '',
        limit: 10,
        pageCursor: '',
        prevPage: '',
        nextPage: '',
      },

      namespaceID: NoID,
      module: undefined,

      options: [],
      value: undefined,
    }
  },

  computed: {
    labelField () {
      return this.module.fields.find(f => f.name === this.pVal('labelField'))
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

    placeholder () {
      return this.pVal('placeholder', 'Select a record')
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

  async created () {
    // Prep the data
    const module = this.pVal('module')
    const moduleType = this.pType('module')
    const namespace = this.pVal('namespace')
    const namespaceType = this.pType('namespace')

    // Resolve bits
    // namespace
    if (namespaceType === 'ID') {
      this.namespaceID = namespace
    } else if (namespaceType === 'ComposeNamespace') {
      this.namespaceID = namespace.namespaceID
    } else {
      // @ts-ignore
      const { set: nn } = await this.$ComposeAPI.namespaceList({ slug: namespace })
      if (!nn || nn.length !== 1) {
        throw new Error('namespace not resolved')
      }

      this.namespaceID = nn[0].namespaceID
    }

    // module; get the full thing as we need fields
    if (moduleType === 'ID') {
      this.module = await this.$ComposeAPI.moduleRead({ namespaceID: this.namespaceID, moduleID: module })
      if (!this.module) {
        throw new Error('module not resolved')
      }
    } else if (moduleType === 'ComposeModule') {
      this.module = module
    } else {
      // @ts-ignore
      const { set: nn } = await this.$ComposeAPI.moduleList({ handle: module, namespaceID: this.namespaceID })
      if (!nn || nn.length !== 1) {
        throw new Error('module not resolved')
      }

      this.module = nn[0]
    }

    // Preload
    this.loadLatest()
  },

  beforeUnmount () {
    this.setDefaultValues()
  },

  methods: {
    encodeValue () {
      if (!this.value) {
        return { '@type': 'Any', '@value': null }
      }

      const { record = {} } = this.options.find(({ recordID }) => recordID === this.value) || {}

      return { '@type': 'ComposeRecord', '@value': record }
    },

    loadLatest () {
      const namespaceID = this.namespaceID
      const moduleID = this.module.moduleID
      const { limit } = this.filter
      if (moduleID && moduleID !== NoID) {
        this.fetchPrefiltered({ namespaceID, moduleID, limit })
      }
    },

    search: debounce(function (query = '') {
      if (query !== this.query) {
        this.query = query
        this.filter.pageCursor = undefined
      }

      const { limit, pageCursor } = this.filter
      const namespaceID = this.namespaceID
      const moduleID = this.module.moduleID
      const queryFields = this.pVal('queryFields') || []

      if (moduleID && moduleID !== NoID) {
        // Determine what fields to use for searching
        // Default to label field
        let qf = queryFields.map(f => f['@value']).filter(f => !!f)

        if ((!qf || qf.length === 0) && this.pVal('labelField')) {
          qf = [this.pVal('labelField')]
        }

        if (query.length > 0) {
          // Construct query
          query = qf.map(qf => {
            return `${qf} LIKE '%${query}%'`
          }).join(' OR ')
        }

        const sort = qf.filter(f => !!f).join(', ')

        this.fetchPrefiltered({ namespaceID, moduleID, query, sort, limit })
      }
    }, 600),

    fetchPrefiltered (q) {
      this.processing = true

      // Prefilter...
      let { query = '' } = q
      if (this.pVal('prefilter')) {
        const pf = this.pVal('prefilter')
        if (query) {
          query = `(${pf}) AND (${query})`
        } else {
          query = pf
        }
      }

      if (this.cancelRequest) {
        this.cancelRequest()
        this.cancelRequest = null
      }

      const { response, cancel } = this.$ComposeAPI.recordListCancellable({ ...q, query })
      this.cancelRequest = cancel

      Promise.all([response(), new Promise(resolve => setTimeout(resolve, 300))])
        .then(([{ filter, set }]) => {
          this.filter = { ...this.filter, ...filter }
          this.filter.nextPage = filter.nextPage
          this.filter.prevPage = filter.prevPage

          this.options = set.map(r => {
            const record = new compose.Record(this.module, r)

            let label
            if (this.labelField) {
              label = this.labelField.isMulti ? record.values[this.pVal('labelField')].join(', ') : record.values[this.pVal('labelField')]
            }

            return {
              recordID: record.recordID,
              label: label || record.recordID,
              record,
            }
          })
          this.processing = false
          return { filter, set }
        })
        .catch((e) => {
          if (axios.isCancel(e)) return
          this.processing = false
          throw e
        })
    },

    goToPage (next = true) {
      this.filter.pageCursor = next ? this.filter.nextPage : this.filter.prevPage
    },

    setDefaultValues () {
      this.processing = false
      this.query = ''
      this.filter = {}
      this.namespaceID = NoID
      this.module = undefined
      this.options = []
      this.value = undefined
    },
  },
}
</script>
