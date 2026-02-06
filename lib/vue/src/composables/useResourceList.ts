import axios from 'axios'
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

interface PaginationState {
  limit: number
  pageCursor?: string
  prevPage: string
  nextPage: string
  total: number
  page: number
}

interface SortingState {
  sortBy?: string
  sortDesc?: boolean
}

interface FilterState {
  [key: string]: any
}

interface ListParams {
  limit: number
  sort?: string
  pageCursor?: string
  incTotal: boolean
  [key: string]: any
}

interface CancellableResponse<T = any> {
  response: () => Promise<T>
  cancel: () => void
}

export function useResourceList<T = any>(
  apiFn: (params: ListParams) => CancellableResponse,
  options: {
    filter?: FilterState
    sorting?: SortingState
    pagination?: Partial<PaginationState>
  } = {},
) {
  const route = useRoute()
  const router = useRouter()

  // Reactive state
  const filter = reactive<FilterState>({ ...options.filter })
  const sorting = reactive<SortingState>({ ...options.sorting })
  const pagination = reactive<PaginationState>({
    limit: 100,
    pageCursor: undefined,
    prevPage: '',
    nextPage: '',
    total: 0,
    page: 1,
    ...options.pagination,
  })

  // Loading and error states - start with loading true to prevent flash of empty state
  const loading = ref(true)
  const error = ref<Error | null>(null)

  // Items state
  const items = ref<T[]>([])

  // Temporary query for handling pageCursor from URL
  let tempQuery: any = undefined

  // Abortable requests
  const abortableRequests: (() => void)[] = []

  // Parse query params into state
  const handleQueryParams = (initial = false) => {
    // Pagination
    let {
      limit = pagination.limit,
      pageCursor = pagination.pageCursor,
      prevPage = pagination.prevPage,
      nextPage = pagination.nextPage,
      total = pagination.total,
      page = pagination.page,
      ...r1
    } = route.query

    limit = parseInt(String(limit))
    total = parseInt(String(total))
    page = parseInt(String(page))

    // If we came to the page with a pageCursor in the URL
    if (initial && pageCursor) {
      tempQuery = route.query
      // Fetch replace query to trigger fetch of total number of items
      router.replace({ query: { ...route.query, limit: 1, pageCursor: undefined } })
      return
    }

    /// To prevent extra list fetch, check if pageCursor is defined (not first page)
    const refresh = String(route.query.pageCursor) !== String(pagination.pageCursor)
    Object.assign(pagination, { limit, pageCursor, prevPage, nextPage, total, page })

    // Sorting
    let { sortBy = sorting.sortBy, sortDesc = sorting.sortDesc, ...r2 } = r1

    sortDesc = sortDesc === true || sortDesc === 'true'

    // Reset pageCursor when sort changes, except on first fetch (so we use the pageCursor from url)
    if (!initial && (sortBy !== sorting.sortBy || sortDesc !== sorting.sortDesc)) {
      pagination.pageCursor = ''
      pagination.page = 1
    }
    Object.assign(sorting, { sortBy, sortDesc })

    // Filtering
    // make sure filter fields are of the right type
    const processedFilter: Record<string, any> = {}
    for (const key in r2) {
      const value = r2[key]
      if (typeof filter[key] === 'boolean') {
        processedFilter[key] = value === 'true'
      } else {
        processedFilter[key] = value
      }
    }

    Object.assign(filter, processedFilter)

    // Only refresh if pageCursor actually changed
    if (refresh) {
      fetchItems()
    }
  }

  // Encode list params for API call
  const encodeListParams = (): ListParams => {
    let { sortBy, sortDesc } = sorting
    const { limit, pageCursor } = pagination

    if (sortBy === 'changedAt') {
      sortBy = 'coalesce(deletedAt, updatedAt, createdAt)'
    }

    const sort = sortBy ? `${sortBy} ${sortDesc ? 'DESC' : 'ASC'}` : undefined

    return {
      limit,
      sort: pageCursor ? undefined : sort,
      ...filter,
      pageCursor,
      incTotal: !pageCursor || !!tempQuery,
    }
  }

  // Encode route params for URL update
  const encodeRouteParams = () => {
    const { limit, pageCursor, page } = pagination

    return {
      path: route.path,
      query: {
        limit: limit.toString(),
        sortBy: sorting.sortBy || '',
        sortDesc: sorting.sortDesc?.toString() || 'false',
        ...Object.fromEntries(
          Object.entries(filter).map(([key, value]) => [
            key,
            typeof value === 'boolean' ? value.toString() : value?.toString() || '',
          ]),
        ),
        page: page.toString(),
        pageCursor: pageCursor || '',
      },
    }
  }

  // Fetch items from API
  const fetchItems = async (updateQuery = true) => {
    abortRequests()

    loading.value = true
    error.value = null

    try {
      const { response, cancel } = apiFn(encodeListParams())
      abortableRequests.push(cancel)

      // Push new router/params to cause URL change
      if (updateQuery && !tempQuery) {
        router.replace(encodeRouteParams())
      }

      const result = await response()

      if (result.filter?.incTotal) {
        pagination.total = result.filter.total
      }

      // This was a fetch of total number of items
      if (tempQuery) {
        const query = tempQuery
        tempQuery = undefined
        router.replace({ query })
        items.value = []
        return
      }

      pagination.pageCursor = undefined
      pagination.nextPage = result.filter?.nextPage || ''
      pagination.prevPage = result.filter?.prevPage || ''

      // Small delay to prevent rapid loading states
      await new Promise(resolve => setTimeout(resolve, 300))

      items.value = result.set || []
    } catch (err) {
      if (!axios.isCancel(err)) {
        error.value = err as Error
        console.error('Failed to fetch items:', err)
      }
    } finally {
      loading.value = false
    }
  }

  // Filter list (reset pagination)
  const filterList = () => {
    // reset pagination when filtering changes
    pagination.pageCursor = ''
    pagination.page = 1

    abortRequests()
    fetchItems()
  }

  // Handle sort changes
  const handleSort = (event: { sortField?: string; sortOrder?: number }) => {
    const { sortField, sortOrder } = event

    if (sortField) {
      // If sorting by the same field, toggle direction
      if (sorting.sortBy === sortField) {
        sorting.sortDesc = !sorting.sortDesc
      } else {
        // New field, default to ascending
        sorting.sortBy = sortField
        sorting.sortDesc = sortOrder === -1
      }

      filterList()
    }
  }

  // Abort all pending requests
  const abortRequests = () => {
    abortableRequests.forEach(cancel => cancel())
    abortableRequests.length = 0
  }

  const handleRowClick = ({ data = {} } = {}) => {
    const { namespaceID, slug } = data as { namespaceID?: string; slug?: string }
    router.push({ name: 'namespace.edit', params: { slug: slug || namespaceID } })
  }

  // Watch for route changes
  watch(
    () => route.fullPath,
    () => handleQueryParams(),
    { immediate: false },
  )

  // Initialize on mount
  onMounted(() => {
    handleQueryParams(true)

    fetchItems()
  })

  // Cleanup on unmount
  onBeforeUnmount(() => {
    abortRequests()
  })

  return {
    // State
    items: computed(() => items.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    filter,
    sorting,
    pagination,

    // Methods
    fetchItems,
    filterList,
    handleSort,
    abortRequests,
    handleRowClick,
  }
}
