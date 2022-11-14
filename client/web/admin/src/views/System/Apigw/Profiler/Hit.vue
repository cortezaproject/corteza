<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('hit.title')"
    >
      <span
        class="text-nowrap"
      />
    </c-content-header>

    <c-profiler-hit-info
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      :hit="hit"
    />
  </b-container>
</template>

<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CProfilerHitInfo from 'corteza-webapp-admin/src/components/Apigw/Profiler/CProfilerHitInfo'
import { mapGetters } from 'vuex'

export default {
  components: {
    CProfilerHitInfo,
  },

  i18nOptions: {
    namespaces: [ 'system.apigw' ],
    keyPrefix: 'profiler',
  },

  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      hit: {},
      info: {
        processing: false,
        success: false,
      },
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('system/', 'apigw-route.create')
    },
  },

  watch: {
    hitID: {
      immediate: true,
      handler () {
        if (this.hitID) {
          this.fetchHit()
        } else {
          this.hit = {}
        }
      },
    },
  },

  mounted () {
    this.fetchHit()
  },

  methods: {
    fetchHit () {
      this.incLoader()

      this.$SystemAPI.apigwProfilerHit({ hitID: this.$route.params.hitID })
        .then(h => { this.hit = h })
        .catch(this.toastErrorHandler(this.$t('notification:queue.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },
  },
}
</script>
