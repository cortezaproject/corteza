<script>
import Function from './Function'

export default {
  extends: Function,

  methods: {
    async getFunctionTypes () {
      return this.$AutomationAPI.functionList()
        .then(({ set }) => {
          this.functions = set.filter(({ kind = '' }) => kind === 'iterator').sort((a, b) => a.meta.short.localeCompare(b.meta.short))
        })
        .catch(this.defaultErrorHandler(this.$t('notification:failed-fetch-functions')))
    },
  },
}
</script>
