<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="title"
    >
      <span
        class="text-nowrap"
      >
        <b-button
          v-if="$Settings.get('apigw.profilerEnabled', false)"
          class="ml-2"
          variant="info"
          :to="{ name: 'system.apigw.profiler' }"
        >
          {{ $t('label') }}
        </b-button>
      </span>
    </c-content-header>

    <c-profiler-route-hits
      :route="$route.params.routeID"
    />
  </b-container>
</template>

<script>
import CProfilerRouteHits from 'corteza-webapp-admin/src/components/Apigw/Profiler/CProfilerRouteHits'

export default {
  components: {
    CProfilerRouteHits,
  },

  i18nOptions: {
    namespaces: ['system.apigw'],
    keyPrefix: 'profiler',
  },

  watch: {
    '$route.params.routeID': {
      immediate: true,
      handler () {
        this.title = `${this.$t('title')} - ${this.decodeRouteID(this.$route.params.routeID)}`
      },
    },
  },

  methods: {
    decodeRouteID (routeID) {
      return atob(routeID)
    },
  },
}
</script>
