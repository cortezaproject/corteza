<script>
import Function from './Function.vue'

export default {
  extends: Function,

  methods: {
    async getFunctionTypes () {
      return this.$AutomationAPI.functionList()
        .then(({ set = [] }) => {
          this.functions = set.filter(({ kind = '' }) => kind === 'iterator').sort((a, b) => a.meta.short.localeCompare(b.meta.short))
        })
        .catch(e => this.toast.add({ severity: 'error', summary: this.$t('notification.failed-fetch-functions'), detail: e?.message, life: 5000 }))
    },
  },
}
</script>
