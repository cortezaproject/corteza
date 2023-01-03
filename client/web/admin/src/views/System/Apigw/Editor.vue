<template>
  <b-container class="py-3">
    <c-content-header :title="$t('title')">
      <span class="text-nowrap">
        <b-button
          v-if="routeID && canCreate"
          data-test-id="button-add"
          variant="primary"
          :to="{ name: 'system.apigw.new' }"
        >
          {{ $t("new") }}
        </b-button>
        <c-permissions-button
          v-if="routeID && canGrant"
          :title="route.endpoint || routeID"
          :target="route.endpoint || routeID"
          :resource="`corteza::system:apigw-route/${routeID}`"
          button-variant="light"
          class="ml-2"
        >
          <font-awesome-icon :icon="['fas', 'lock']" />
          {{ $t("permissions") }}
        </c-permissions-button>
      </span>
    </c-content-header>

    <c-route-editor-info
      :route="route"
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      @submit="onInfoSubmit"
      @delete="onInfoDelete"
    />

    <c-filters-stepper
      v-if="routeID"
      ref="stepper"
      :fetching="stepper.fetching"
      :processing="stepper.processing"
      :success="stepper.success"
      :filters.sync="filters"
      :available-filters="availableFilters"
      :steps="steps"
      @submit="onFiltersSubmit"
    />

    <c-profiler-route-hits
      v-if="routeID && showProfiler"
      :route="routeEndpoint"
      class="mt-3"
    />
  </b-container>
</template>
<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CRouteEditorInfo from 'corteza-webapp-admin/src/components/Apigw/CRouteEditorInfo'
import CFiltersStepper from 'corteza-webapp-admin/src/components/Apigw/CFiltersStepper'
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import CProfilerRouteHits from 'corteza-webapp-admin/src/components/Apigw/Profiler/CProfilerRouteHits'

export default {
  components: {
    CRouteEditorInfo,
    CFiltersStepper,
    CProfilerRouteHits,
  },

  i18nOptions: {
    namespaces: ['system.apigw'],
    keyPrefix: 'editor',
  },

  mixins: [editorHelpers],

  props: {
    routeID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      route: {},
      routeEndpoint: undefined,

      info: {
        processing: false,
        success: false,
      },
      stepper: {
        fetching: false,
        processing: false,
        success: false,
      },

      filters: [],
      availableFilters: [],
      steps: [],
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('system/', 'apigw-route.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },

    showProfiler () {
      return this.$Settings.get('apigw.profiler.enabled', false) && (this.$Settings.get('apigw.profiler.global', false) || this.filters.some(({ ref, enabled = false }) => ref === 'profiler' && enabled))
    },
  },

  watch: {
    routeID: {
      immediate: true,
      handler () {
        this.routeEndpoint = undefined

        if (this.routeID) {
          this.fetchSteps()
          this.fetchRoute()
          this.fetchFilters()
        } else {
          this.route = {
            method: 'GET',
          }
        }
      },
    },
  },
  methods: {
    fetchRoute () {
      this.incLoader()

      this.$SystemAPI.apigwRouteRead({ routeID: this.routeID, incFlags: 1 })
        .then((api) => {
          this.route = api
          this.routeEndpoint = btoa(api.endpoint)
        })
        .catch(this.toastErrorHandler(this.$t('notification:gateway.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onInfoSubmit (route) {
      this.info.processing = true
      if (this.routeID) {
        this.$SystemAPI
          .apigwRouteUpdate(route)
          .then(() => {
            this.fetchRoute()

            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:gateway.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:gateway.update.error')))
          .finally(() => {
            this.info.processing = false
          })
      } else {
        this.$SystemAPI
          .apigwRouteCreate(route)
          .then(({ routeID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:gateway.create.success'))

            this.$router.push({
              name: 'system.apigw.edit',
              params: { routeID },
            })
          })
          .catch(this.toastErrorHandler(this.$t('notification:gateway.create.error')))
          .finally(() => {
            this.info.processing = false
          })
      }
    },

    onInfoDelete () {
      this.incLoader()

      if (this.route.deletedAt) {
        this.$SystemAPI
          .apigwRouteUndelete({ routeID: this.routeID })
          .then(() => {
            this.fetchRoute()

            this.toastSuccess(this.$t('notification:gateway.undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:gateway.undelete.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI
          .apigwRouteDelete({ routeID: this.routeID })
          .then(() => {
            this.fetchRoute()

            this.toastSuccess(this.$t('notification:gateway.delete.success'))
            this.$router.push({ name: 'system.apigw' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:gateway.delete.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    onFiltersSubmit () {
      if (this.routeID) {
        this.stepper.processing = true

        Promise.all(this.filters.map(filter => {
          if (filter.created || filter.updated || filter.deleted) {
            filter.params = this.encodeParams(filter.params)
            filter.weight = filter.weight.toString()

            if (filter.filterID && filter.filterID !== NoID) {
              return filter.deleted ? this.deleteFilter(filter) : this.updateFilter(filter)
            } else {
              return filter.deleted ? undefined : this.createFilter(filter)
            }
          }
        })).then(async () => {
          await this.fetchFilters()

          this.animateSuccess('stepper')
          this.toastSuccess(this.$t('notification:gateway.filter.update.success'))
        })
          .catch(this.toastErrorHandler(this.$t('notification:gateway.filter.update.error')))
          .finally(() => {
            this.stepper.processing = false
          })
      }
    },

    createFilter (filter) {
      return this.$SystemAPI.apigwFilterCreate({ ...filter, routeID: this.routeID })
    },

    updateFilter (filter) {
      return this.$SystemAPI.apigwFilterUpdate({ ...filter, routeID: this.routeID })
    },

    deleteFilter ({ filterID = '' }) {
      if (filterID) {
        return this.$SystemAPI.apigwFilterDelete({ filterID })
      }
    },

    fetchFilters () {
      this.incLoader()
      this.stepper.fetching = true

      this.$SystemAPI.apigwFilterList({ routeID: this.routeID })
        .then(({ set = [] }) => {
          return this.setRouteFilters(set)
        })
        .catch(this.toastErrorHandler(this.$t('notification:gateway.filter.fetch.error')))
        .finally(() => {
          this.decLoader()
          this.stepper.fetching = false
        })
    },

    setRouteFilters (routeFilters = []) {
      return this.fetchAllAvailableFilters().then(() => {
        this.filters = (routeFilters || []).map(filter => {
          const f = { ...this.availableFilters.find((af) => af.ref === filter.ref) }
          f.params = this.decodeParams(f, { ...filter.params })
          f.weight = parseInt(filter.weight)
          f.filterID = filter.filterID
          f.enabled = !!filter.enabled
          return { ...f }
        })
      })
    },

    decodeParams (filter = {}, values = {}) {
      const { params = [] } = filter
      return params.map(({ label, type }) => {
        return {
          label,
          type,
          value: values[label],
        }
      })
    },

    encodeParams (params = []) {
      return params.reduce((result, p) => {
        result[p.label] = p.value
        return result
      }, {})
    },

    fetchAllAvailableFilters () {
      this.incLoader()

      return this.$SystemAPI.apigwFilterDefFilter()
        .then((api) => {
          this.availableFilters = api.map((f) => {
            return { ...f, ref: f.name, enabled: true, options: { checked: false } }
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:gateway.filter.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    fetchSteps () {
      this.steps = ['prefilter', 'processer', 'postfilter']
    },
  },
}
</script>
