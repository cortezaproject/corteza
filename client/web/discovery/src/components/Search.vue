<template>
  <b-container
    fluid
    class="h-100 mh-100 p-0"
  >
    <b-row
      no-gutters
    >
      <b-col
        md="12"
        :lg="map.show ? '7' : '12'"
        :xl="map.show ? '8' : '12'"
        class="results-container"
        :class="{ 'with-map': map.show }"
      >
        <b-form
          @submit.prevent="onQuerySubmit()"
        >
          <b-form-group class="px-3">
            <b-input-group
              size="lg"
            >
              <b-form-input
                ref="query"
                v-model="query"
                :placeholder="$t('input-placeholder')"
                autocomplete="off"
              />
              <b-input-group-append>
                <b-button
                  v-if="query"
                  variant="link"
                  class="clear-query position-absolute text-secondary border-0"
                  @click="clearQuery()"
                >
                  <font-awesome-icon
                    :icon="['fas', 'times']"
                  />
                </b-button>
                <b-button
                  variant="link"
                  class="bg-white"
                  style="border: 2px solid #E4E9EF;"
                  @click="onQuerySubmit()"
                >
                  <font-awesome-icon
                    :icon="['fas', 'search']"
                  />
                </b-button>
                <input
                  type="submit"
                  hidden
                >
              </b-input-group-append>
            </b-input-group>

            <div
              class="d-flex align-items-center justify-content-between px-1 mt-1 text-muted"
            >
              <span
                :class="{ 'discovering': $store.state.processing }"
              >
                {{ searchDescription }}
              </span>
              <span>
                Use <samp>"text"</samp> for exact match
              </span>
            </div>
          </b-form-group>
        </b-form>

        <b-row
          class="results w-100 m-0 mh-100 overflow-auto"
        >
          <div
            v-if="$store.state.processing || !total.actual"
            class="position-absolute d-flex align-items-center justify-content-center w-100 h-100"
            style="opacity: 0.8; z-index: 1; background-color: #F3F3F5;"
          >
            <h5
              class="mb-5"
            >
              <b-spinner
                v-if="$store.state.processing"
                variant="primary"
                class="p-4"
              />
              <span
                v-else-if="!total.actual"
              >
                No results
              </span>
            </h5>
          </div>
          <b-col
            v-for="(hit, i) in filteredHits"
            :key="i"
            sm="12"
            md="6"
            :lg="map.show ? '6': '4'"
            class="py-3"
          >
            <result
              :id="hit.value.recordID || hit.value.moduleID"
              :index="i"
              :hit="hit"
              :show-map="map.show"
              :class="{ 'border-primary border shadow': map.clickedMarker && [hit.value.recordID, hit.value.moduleID].includes(map.clickedMarker) }"
              @hover="map.hoverIndex = $event"
            />
          </b-col>
        </b-row>

        <div
          class="position-fixed map-button"
        >
          <b-button
            variant="warning"
            class="rounded-circle p-3"
            @click="toggleMap"
          >
            <font-awesome-icon
              :icon="['fas', 'map-marked-alt']"
              class="h3 mb-0"
            />
          </b-button>
        </div>
      </b-col>

      <b-col
        v-if="map.show"
        md="12"
        lg="5"
        xl="4"
      >
        <discovery-map
          :markers="map.markers"
          :hover-index="map.hoverIndex"
          @hover="markerHovered"
        />
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import Result from './Results'
import DiscoveryMap from './DiscoveryMap.vue'

export default {
  i18nOptions: {
    namespaces: 'search',
  },

  components: {
    Result,
    DiscoveryMap,
  },

  data () {
    return {
      query: '',

      hits: [],
      filteredHits: [],

      pagination: {
        limit: 50,
        from: 0,
        size: 0,
      },

      total: {
        all: 0,
        actual: 0,
      },

      initial: false,

      map: {
        show: false,
        markers: [],
        clickedMarker: undefined,
        hoverIndex: undefined,
      },
    }
  },

  computed: {
    searchDescription () {
      if (this.$store.state.processing) {
        return 'Discovering'
      }

      if (this.total.all > 0) {
        return `Showing ${this.total.actual} of ${this.total.all} results`
      }

      return ''
    },
  },

  watch: {
    '$store.state.types': {
      handler: function () {
        this.getFilteredData()
      },
    },

    '$store.state.modules': {
      handler () {
        if (this.initial) return
        this.pagination.size = this.pagination.limit
        this.getSearchData(this.query)
      },
    },

    '$store.state.namespaces': {
      handler () {
        if (this.initial) return
        this.pagination.size = this.pagination.limit
        this.getSearchData(this.query)
      },
    },
  },

  created () {
    this.initial = true

    const { query = '', modules = [], namespaces = [], size = 0 } = this.$route.query

    this.query = query
    this.pagination.size = size

    this.$store.commit('updateModules', Array.isArray(modules) ? modules : [modules])
    this.$store.commit('updateNamespaces', Array.isArray(namespaces) ? namespaces : [namespaces])

    this.getSearchData(this.query)

    setTimeout(() => {
      this.initial = false
    }, 1000)
  },

  mounted () {
    const listElm = document.querySelector('.results')
    listElm.addEventListener('scroll', e => {
      if (listElm.scrollTop + listElm.clientHeight >= listElm.scrollHeight - 10) {
        if (!this.$store.state.processing && this.total.actual < this.total.all) {
          this.getSearchData(this.query, true)
        }
      }
    })
  },

  methods: {
    getSearchData (query = '', append = false) {
      this.$store.commit('updateProcessing', true)

      if (!append) {
        this.map.markers = []
      }

      // Filters
      const modules = this.$store.state.modules
      const namespaces = this.$store.state.namespaces

      // Pagination
      if (append) {
        this.pagination.size += this.pagination.limit
      }

      const { size } = this.pagination

      this.updateRouteQuery({ query, modules, namespaces, size })

      this.$DiscoveryAPI.query({ query, modules, namespaces, size })
        .then((response = {}) => {
          if (response) {
            this.hits = (response.hits || [])

            this.total.all = response.total_results || 0

            this.getFilteredData()

            this.pagination = {
              ...this.pagination,
              from: response.from || 0,
              size: response.size || 0,
            }

            this.$store.commit('updateAggregations', response.aggregations)

            this.getMarkers()
          }
        }).catch(e => {
          this.toastErrorHandler(this.$t('notification:search.failed'))(e)
          this.hits = []
          this.filteredHits.splice(0, this.filteredHits.length)
        })
        .finally(() => {
          this.$store.commit('updateProcessing', false)
        })
    },

    getFilteredData () {
      let filteredHits = this.hits

      if (this.$store.state.types.length > 0 && this.hits.length) {
        filteredHits = this.hits.filter(hit => this.$store.state.types.includes(hit.type))
      }

      this.filteredHits.splice(0, this.filteredHits.length, ...filteredHits)
      this.total.actual = this.filteredHits.length
    },

    onQuerySubmit () {
      if (!this.$store.state.processing) {
        this.pagination.size = this.pagination.limit
        this.getSearchData(this.query)
      }
    },

    clearQuery () {
      this.query = ''
      this.$refs.query.focus()
    },

    getMarkers () {
      const markers = []

      this.filteredHits.forEach(({ type, value }) => {
        if (type === 'compose:record' && Array.isArray(value.values)) {
          const id = value.recordID
          value.values.forEach(({ value = [] }) => {
            const isGeometry = value && value.find(v => {
              return v.toString().includes('{"coordinates":[')
            })

            if (isGeometry) {
              value.forEach(coordinates => {
                coordinates = JSON.parse(coordinates || '{}').coordinates
                if (coordinates) {
                  markers.push({ id, coordinates })
                }
              })
            }
          })
        }
      })

      this.map.markers = markers
    },

    markerHovered (ID) {
      if (ID) {
        document.getElementById(ID).scrollIntoView({
          behavior: 'smooth',
          block: 'center',
        })
      }

      this.map.clickedMarker = ID
    },

    toggleMap () {
      this.map.show = !this.map.show
    },

    updateRouteQuery ({ query = undefined, modules = [], namespaces = [], size = 0 }) {
      if (JSON.stringify(this.$route.query) !== JSON.stringify({ query, modules, namespaces })) {
        this.$router.push({ query: { query: query || undefined, modules, namespaces, size } })
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.results-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 64px);
}

.results-container.with-map {
  height: calc(60vh - 64px);
}

@media (min-width: 992px) {
  .results-container {
    height: calc(100vh - 64px) !important;
  }
}

.results {
  flex: 1 1 auto;
}

.clear-query {
  z-index: 3 !important;
  right: 52px;
  top: 2px;
}

.map-button {
  bottom: 1rem;
  right: 1rem;
  z-index: 99999;
}

// https://stackoverflow.com/a/40991531/17926309
.discovering::after {
  display: inline-block;
  animation: discovering steps(1, end) 1s infinite;
  content: '';
}

@keyframes discovering {
  0% { content: ''; }
  25% { content: '.'; }
  50% { content: '..'; }
  75% { content: '...'; }
  100% { content: ''; }
}
</style>
