import { debounce } from 'lodash'
import { mapActions } from 'vuex'

export default {
  data () {
    return {
      /**
       * Placeholder, if component does not define own filter
       */
      filter: {},

      pagination: {
        limit: 10,
        pageCursor: undefined,
        prevPage: '',
        nextPage: '',
        total: 0,
        page: 1,
      },

      // Used to save query when fetching for total if we came to the page with a pageCursor in URL
      tempQuery: undefined,

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
      handler () {
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
      // Pagination
      let {
        limit = this.pagination.limit,
        pageCursor = this.pagination.pageCursor,
        prevPage = this.pagination.prevCursor,
        nextPage = this.pagination.nextCursor,
        total = this.pagination.total,
        page = this.pagination.page,
        ...r1
      } = this.$route.query

      limit = parseInt(limit)
      total = parseInt(total)
      page = parseInt(page)

      // If we came to the page with a pageCursor in the URL
      if (initial && pageCursor) {
        this.tempQuery = this.$route.query
        // Fetch replace query to trigger fetch of total number of items
        this.$router.replace({ query: { ...this.$route.query, limit: 1, pageCursor: undefined } })
        return
      }

      /// To prevent extra list fetch, check if pageCursor is defined (not first page)
      const refresh = this.$route.query.pageCursor !== this.pagination.pageCursor
      this.pagination = { limit, pageCursor, prevPage, nextPage, total, page }

      // Sorting
      let { sortBy = this.sorting.sortBy, sortDesc = this.sorting.sortDesc, ...r2 } = r1

      sortDesc = sortDesc === true || sortDesc === 'true'

      // Reset pageCursor when sort changes, except on first fetch (so we use the pageCursor from url)
      if (!initial && (sortBy !== this.sorting.sortBy || sortDesc !== this.sorting.sortDesc)) {
        this.pagination.pageCursor = ''
        this.pagination.page = 1
      }
      this.sorting = { sortBy, sortDesc }

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
      // reset pagination when filtering changes
      //
      // we want to prevent situations with page is preset to a number that
      // exceeds the number of pages of a filtered results
      this.pagination.pageCursor = ''
      this.pagination.page = 1

      // notify b-table about the change
      //
      // this effectively calls items()/procListResults()
      this.$root.$emit('bv::refresh::table', 'resource-list')
    }, 300),

    encodeListParams () {
      const { sortBy, sortDesc } = this.sorting
      const { limit, pageCursor } = this.pagination

      const sort = sortBy ? `${sortBy} ${sortDesc ? 'DESC' : 'ASC'}` : undefined

      return {
        limit,
        sort: pageCursor ? undefined : sort,
        ...this.filter,
        pageCursor,
        incTotal: !pageCursor || this.tempQuery,
      }
    },

    encodeRouteParams () {
      const { limit, pageCursor, page } = this.pagination

      return {
        query: {
          limit,
          ...this.sorting,
          ...this.filter,
          page,
          pageCursor,
        },
      }
    },

    /**
     *
     * @param p {Promise}
     * @returns {Promise}
     */
    procListResults (p, updateQuery = true) {
      // Push new router/params to cause URL change
      //
      // We want this because in case when user refreshes or shares URL
      // he needs to land on the same page with the same parameters
      if (updateQuery && !this.tempQuery) {
        this.$router.replace(this.encodeRouteParams())
      }

      return p.then(async ({ set, filter } = {}) => {
        if (filter.incTotal) {
          this.pagination.total = filter.total
        }

        // This was a fetch of total number of items
        if (this.tempQuery) {
          const query = this.tempQuery
          this.tempQuery = undefined

          this.$router.replace({ query })
          return []
        }

        this.pagination.pageCursor = undefined
        this.pagination.nextPage = filter.nextPage
        this.pagination.prevPage = filter.prevPage

        return set
      }).catch(this.toastErrorHandler(this.$t('notification:list.load.error')))
        .finally(async () => {
          await new Promise(resolve => setTimeout(resolve, 300))
        })
    },

    genericRowClass (item) {
      return { 'text-secondary': item && !!item.deletedAt }
    },
  },
}
