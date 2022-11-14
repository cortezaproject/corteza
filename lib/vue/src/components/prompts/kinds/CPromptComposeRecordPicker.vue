<template>
  <div>
    <p>
      {{ message }}
    </p>
    <div
      class="text-center m-2"
    >
      <vue-select
        v-model="value"
        :options="options"
        :loading="processing"
        append-to-body
        :calculate-position="calculatePosition"
        option-value="recordID"
        option-text="label"
        placeholder="Select record"
        :filterable="false"
        class="bg-white w-100 mb-3"
        @search="search"
      >
        <c-pagination
          v-if="showPagination"
          slot="list-footer"
          :has-prev-page="hasPrevPage"
          :has-next-page="hasNextPage"
          @prev="goToPage(false)"
          @next="goToPage(true)"
        />
      </vue-select>

      <b-button
        @click="$emit('submit', { value: encodeValue() })"
        :disabled="loading"
      >
        {{ pVal('buttonLabel', 'Submit') }}
      </b-button>
    </div>
  </div>
</template>
<script lang="js">
import base from './base.vue'
import CPagination from '../common/CPagination.vue'
import { pVal } from '../utils.ts'
import { VueSelect } from 'vue-select'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { createPopper } from '@popperjs/core'
import { debounce } from 'lodash'

export default {
  extends: base,
  name: 'c-prompt-compose-record-picker',

  components: {
    VueSelect,
    CPagination,
  },

  data () {
    return {
      processing: false,
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
    // moduleFields returns the available field names
    moduleFields () {
      if (!this.module) {
        return []
      }
      return this.module.fields.map(({ name }) => name)
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

  methods: {
    encodeValue () {
      if (!this.value) {
        return { '@type': 'Any', '@value': null }
      }

      return { '@type': 'ComposeRecord', '@value': this.value.record }
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

      if (moduleID && moduleID !== NoID) {
        // Determine what fields to use for searching
        // Default to label field
        let qf = this.pVal('queryFields').map(f => f['@value']).filter(f => !!f)
        if (!qf || qf.length === 0) {
          qf = [this.pVal('labelField')]
        }

        if (query.length > 0) {
          // Construct query
          query = qf.map(qf => {
            return `${qf} LIKE '%${query}%'`
          }).join(' OR ')
        }

        this.fetchPrefiltered({ namespaceID, moduleID, query, limit })
      }
    }, 300),

    calculatePosition (dropdownList, component, { width }) {
      /**
       * We need to explicitly define the dropdown width since
       * it is usually inherited from the parent with CSS.
       */
      dropdownList.style.width = width
      dropdownList.style['z-index'] = 10000

      /**
       * Here we position the dropdownList relative to the $refs.toggle Element.
       *
       * The 'offset' modifier aligns the dropdown so that the $refs.toggle and
       * the dropdownList overlap by 1 pixel.
       *
       * The 'toggleClass' modifier adds a 'drop-up' class to the Vue Select
       * wrapper so that we can set some styles for when the dropdown is placed
       * above.
       */
      const popper = createPopper(component.$refs.toggle, dropdownList, {
        placement: 'bottom',
        modifiers: [
          {
            name: 'offset',
            options: {
              offset: [0, -1],
            },
          },
          {
            name: 'toggleClass',
            enabled: true,
            phase: 'write',
            fn ({ state }) {
              component.$el.classList.toggle('drop-up', state.placement === 'top')
            },
          }],
      })

      /**
       * To prevent memory leaks Popper needs to be destroyed.
       * If you return function, it will be called just before dropdown is removed from DOM.
       */
      return () => popper.destroy()
    },


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

      this.$ComposeAPI.recordList({ ...q, query })
        .then(({ filter, set }) => {
          this.filter = { ...this.filter, ...filter }
          this.filter.nextPage = filter.nextPage
          this.filter.prevPage = filter.prevPage

          this.options = set.map(r => {
            const record = new compose.Record(this.module, r)
            const label = record.values[this.pVal('labelField')] || record.recordID

            return {
              recordID: record.recordID,
              label,
              record,
            }
          })

          return { filter, set }
        })
        .finally(() => {
          this.processing = false
        })
    },

    goToPage (next = true) {
      this.filter.pageCursor = next ? this.filter.nextPage : this.filter.prevPage
    },
  },
}
</script>
