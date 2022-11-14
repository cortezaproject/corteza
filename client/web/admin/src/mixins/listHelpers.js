import { debounce } from 'lodash'
import { mapActions } from 'vuex'

export default {
  data () {
    return {
      /**
       * Placeholder, if component does not define own filter
       */
      filter: {},

      paging: {
        limit: 10,
        pageCursor: undefined,
        prevPage: '',
        nextPage: '',
      },

      sorting: {},
    }
  },

  watch: {
    /**
     * When fullPath for this component changes, we most likely should update our
     * filters.
     * @todo make this.filter reactive
     */
    '$route.fullPath': {
      handler: function () {
        this.handleQueryParams()
      },
    },
  },

  created () {
    this.handleQueryParams(true)
  },

  methods: {
    ...mapActions({
      incLoader: 'ui/incLoader',
      decLoader: 'ui/decLoader',
    }),

    /**
     * Parses query params into list filtering params.
     * @param initial {Boolean} - used to determine it this is the initial fetch
     */
    handleQueryParams (initial = false) {
      // Paging
      const {
        limit = this.paging.limit,
        pageCursor = this.paging.pageCursor,
        prevPage = this.paging.prevCursor,
        nextPage = this.paging.nextCursor,
        ...r1
      } = this.$route.query

      /// To prevent extra list fetch, check if pageCursor is defined (not first page)
      const refresh = this.$route.query.pageCursor !== this.paging.pageCursor
      this.paging = { limit, pageCursor, prevPage, nextPage }

      // Sorting
      const { sortBy = this.sorting.sortBy, sortDesc = this.sorting.sortDesc, ...r2 } = r1

      // Reset pageCursor when sort changes, except on first fetch (so we use the pageCursor from url)
      if (!initial && (sortBy !== this.sorting.sortBy || sortDesc !== this.sorting.sortDesc)) {
        this.paging.pageCursor = ''
      }
      this.sorting = { sortBy, sortDesc: sortDesc === true || sortDesc === 'true' }

      // Filtering
      // make sure filter fields are of the right type
      for (const key in r2) {
        if (typeof this.filter[key] === 'boolean') {
          r2[key] = r2[key] === 'true'
        }
      }

      this.filter = { ...this.filter, ...r2 }

      // Only refresh if pageCursor actually changed,
      if (refresh) {
        this.$root.$emit('bv::refresh::table', 'resource-list')
      }
    },

    filterList: debounce(function () {
      // reset paging when filtering changes
      //
      // we want to prevent situations with page is preset to a number that
      // exceeds the number of pages of a filtered results
      this.paging.pageCursor = ''

      // notify b-table about the change
      //
      // this effectively calls items()/procListResults()
      this.$root.$emit('bv::refresh::table', 'resource-list')
    }, 300),

    /**
     * Encode list params
     * @returns {{perPage: *, page: *, sort: (string|*)}}
     */
    encodeListParams () {
      const { sortBy, sortDesc } = this.sorting

      const sort = sortBy ? `${sortBy} ${sortDesc ? 'DESC' : 'ASC'}` : undefined

      return {
        ...this.filter,
        ...this.paging,
        ...{ nextPage: undefined, prevPage: undefined },
        sort: this.paging.pageCursor ? undefined : sort,
      }
    },

    encodeRouteParams () {
      return {
        query: {
          ...this.paging,
          ...this.sorting,
          ...this.filter,
        },
      }
    },

    /**
     *
     * @param p {Promise}
     * @returns {Promise}
     */
    procListResults (p, updateQuery = true) {
      this.incLoader()

      // Push new router/params to cause URL change
      //
      // We want this because in case when user refreshes or shares URL
      // he needs to land on the same page with the same parameters
      if (updateQuery) {
        this.$router.push(this.encodeRouteParams())
      }

      return p.then(async ({ set, filter } = {}) => {
        this.paging.pageCursor = undefined
        this.paging.nextPage = filter.nextPage
        this.paging.prevPage = filter.prevPage

        return set
      }).catch(this.toastErrorHandler(this.$t('notification:list.load.error')))
        .finally(async () => {
          await new Promise(
            resolve => setTimeout(resolve, 300)
          )

          this.decLoader()
        })
    },
  },
}
