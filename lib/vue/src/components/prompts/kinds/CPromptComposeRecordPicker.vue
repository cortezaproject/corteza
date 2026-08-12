<template>
  <div class="flex flex-col gap-4">
    <p v-if="message" class="m-0 whitespace-pre-wrap" v-html="message" />

    <label v-if="label" class="text-sm font-medium text-color">{{ label }}</label>

    <Message v-if="resolveError" severity="error" :closable="false" class="text-sm">
      {{ resolveError }}
    </Message>

    <Select
      v-else
      v-model="selectedRecordID"
      :options="options"
      option-label="label"
      option-value="recordID"
      :placeholder="placeholder"
      class="w-full"
      filter
      :loading="processing"
      :disabled="loading"
      empty-message=""
      append-to="self"
      @filter="onFilter"
    >
      <template #footer>
        <div v-if="showPagination" class="flex gap-1 p-1">
          <Button
            class="flex-1"
            severity="secondary"
            size="small"
            icon="pi pi-chevron-left"
            :disabled="!hasPrevPage || processing"
            @click="goToPage(false)"
          />
          <Button
            class="flex-1"
            severity="secondary"
            size="small"
            icon="pi pi-chevron-right"
            :disabled="!hasNextPage || processing"
            @click="goToPage(true)"
          />
        </div>
      </template>
    </Select>

    <div class="flex justify-end">
      <Button
        :disabled="loading || !selectedRecordID"
        :label="pVal('buttonLabel', tF('general.label.submit', 'Submit'))"
        @click="$emit('submit', { value: encodeValue() })"
      />
    </div>
  </div>
</template>

<script>
import { compose, NoID } from '@planetcrust/human-js'
import { debounce } from 'lodash-es'
import base from './base.vue'

const PAGE_LIMIT = 10
const SEARCH_DEBOUNCE_MS = 600

export default {
  name: 'CPromptComposeRecordPicker',
  extends: base,
  emits: ['submit'],
  data() {
    return {
      processing: false,
      query: '',
      filter: {
        query: '',
        sort: '',
        limit: PAGE_LIMIT,
        pageCursor: '',
        prevPage: '',
        nextPage: '',
      },
      namespaceID: NoID,
      module: undefined,
      options: [],
      selectedRecordID: undefined,
      cancelRequest: null,
      resolveError: '',
    }
  },
  computed: {
    placeholder() {
      return this.pVal(
        'placeholder',
        this.tF('prompt.record-picker.placeholder', 'Select a record'),
      )
    },
    labelField() {
      if (!this.module) return undefined
      return this.module.fields?.find(f => f.name === this.pVal('labelField'))
    },
    showPagination() {
      return this.hasPrevPage || this.hasNextPage
    },
    hasPrevPage() {
      return !!this.filter.prevPage
    },
    hasNextPage() {
      return !!this.filter.nextPage
    },
  },
  watch: {
    'filter.pageCursor': {
      handler(pageCursor) {
        if (pageCursor !== undefined) {
          this.fetchPrefiltered({
            namespaceID: this.namespaceID,
            moduleID: this.module?.moduleID,
            query: this.filter.query,
            sort: this.filter.sort,
            limit: this.filter.limit,
            pageCursor,
          })
        }
      },
    },
  },
  async created() {
    this.search = debounce(this.runSearch, SEARCH_DEBOUNCE_MS)

    try {
      await this.resolveNamespace()
      await this.resolveModule()
      this.loadLatest()
    } catch (e) {
      // Surface the failure inline so the user understands why the picker is
      // empty and can close the prompt instead of waiting on a stuck spinner.
      const reason = e?.message || 'Unknown error'
      this.resolveError = this.tF(
        'prompt.record-picker.resolve-failed',
        `Could not load records: ${reason}. Check that the namespace and module exist and that you have access.`,
      )
      // eslint-disable-next-line no-console
      console.warn('[CPromptComposeRecordPicker] resolve failed:', e)
    }
  },
  beforeUnmount() {
    if (this.cancelRequest) {
      this.cancelRequest()
      this.cancelRequest = null
    }
    if (this.search?.cancel) {
      this.search.cancel()
    }
  },
  methods: {
    async resolveNamespace() {
      const namespace = this.pVal('namespace')
      const namespaceType = this.pType('namespace')

      if (namespaceType === 'ID') {
        this.namespaceID = namespace
      } else if (namespaceType === 'ComposeNamespace') {
        this.namespaceID = namespace?.namespaceID
      } else {
        const { set = [] } = await this.$ComposeAPI.namespaceList({ slug: namespace })
        if (set.length !== 1) throw new Error('namespace not resolved')
        this.namespaceID = set[0].namespaceID
      }
    },

    async resolveModule() {
      const module = this.pVal('module')
      const moduleType = this.pType('module')

      if (moduleType === 'ID') {
        this.module = await this.$ComposeAPI.moduleRead({
          namespaceID: this.namespaceID,
          moduleID: module,
        })
      } else if (moduleType === 'ComposeModule') {
        this.module = module
      } else {
        const { set = [] } = await this.$ComposeAPI.moduleList({
          handle: module,
          namespaceID: this.namespaceID,
        })
        if (set.length !== 1) throw new Error('module not resolved')
        this.module = set[0]
      }
    },

    loadLatest() {
      const namespaceID = this.namespaceID
      const moduleID = this.module?.moduleID
      if (moduleID && moduleID !== NoID) {
        this.fetchPrefiltered({ namespaceID, moduleID, limit: this.filter.limit })
      }
    },

    onFilter(event) {
      this.search(event?.value ?? '')
    },

    runSearch(query = '') {
      if (query !== this.query) {
        this.query = query
        this.filter.pageCursor = ''
      }

      const moduleID = this.module?.moduleID
      if (!moduleID || moduleID === NoID) return

      const queryFields = this.pVal('queryFields') || []
      let qf = (Array.isArray(queryFields) ? queryFields : [])
        .map(f => (f && typeof f === 'object' && '@value' in f ? f['@value'] : f))
        .filter(f => !!f)

      if (qf.length === 0 && this.pVal('labelField')) {
        qf = [this.pVal('labelField')]
      }

      let composedQuery = ''
      if (query.length > 0 && qf.length > 0) {
        const escaped = query.replace(/'/g, "''")
        composedQuery = qf.map(f => `${f} LIKE '%${escaped}%'`).join(' OR ')
      }
      const sort = qf.filter(f => !!f).join(', ')

      this.fetchPrefiltered({
        namespaceID: this.namespaceID,
        moduleID,
        query: composedQuery,
        sort,
        limit: this.filter.limit,
      })
    },

    fetchPrefiltered(params) {
      this.processing = true

      let { query = '' } = params
      const prefilter = this.pVal('prefilter')
      if (prefilter) {
        query = query ? `(${prefilter}) AND (${query})` : prefilter
      }

      if (this.cancelRequest) {
        this.cancelRequest()
        this.cancelRequest = null
      }

      const { response, cancel } = this.$ComposeAPI.recordListCancellable({ ...params, query })
      this.cancelRequest = cancel

      response()
        .then(({ filter, set }) => {
          this.filter = { ...this.filter, ...filter, query: filter.query || '' }
          this.options = set.map(raw => {
            const record = new compose.Record(this.module, raw)
            let label
            if (this.labelField) {
              const v = record.values?.[this.labelField.name]
              label = this.labelField.isMulti && Array.isArray(v) ? v.join(', ') : v
            }
            return {
              recordID: record.recordID,
              label: label || record.recordID,
              record,
            }
          })
        })
        .catch(e => {
          if (e?.name === 'CanceledError' || e?.message === 'canceled') return
          // network error — keep stale options; the picker remains usable
        })
        .finally(() => {
          this.processing = false
        })
    },

    goToPage(next = true) {
      this.filter.pageCursor = next ? this.filter.nextPage : this.filter.prevPage
    },

    encodeValue() {
      const entry = this.options.find(({ recordID }) => recordID === this.selectedRecordID)
      if (!entry) return { '@type': 'Any', '@value': null }
      return { '@type': 'ComposeRecord', '@value': entry.record }
    },
  },
}
</script>
