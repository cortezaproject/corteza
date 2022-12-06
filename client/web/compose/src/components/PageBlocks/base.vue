<script>
import { compose } from '@cortezaproject/corteza-js'
import Wrap from './Wrap'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    Wrap,
  },

  props: {
    boundingRect: {
      type: Object,
      required: false,
      default: undefined,
    },

    blockIndex: {
      type: Number,
      default: -1,
    },

    namespace: {
      type: compose.Namespace,
      required: true,
    },

    page: {
      type: compose.Page,
      required: true,
    },

    block: {
      type: compose.PageBlock,
      required: true,
    },

    module: {
      type: compose.Module,
      required: false,
      default: undefined,
    },

    record: {
      type: compose.Record,
      required: false,
      default: undefined,
    },

    mode: {
      type: String,
      required: false,
      default: '',
    },
  },

  data () {
    return {
      refreshInterval: null,
      key: 0,
    }
  },

  computed: {
    options: {
      get () {
        return this.block.options
      },

      set (options) {
        this.block.options = options
      },
    },

    allowsRefresh () {
      return this.options.refreshRate >= 5 && ['page', 'page.record'].includes(this.$route.name)
    },
  },

  beforeDestroy () {
    clearInterval(this.refreshInterval)
  },

  methods: {
    /**
     *
     * @param {*} refreshFunction
     * @param {*} forceRerender
     * Key is used to force a component to re-render
     * However, sometimes reloading data is enough
     * Attach :key="key" if you need to re-render a component
     *
     */
    refreshBlock (refreshFunction, forceRerender) {
      if (!this.allowsRefresh || this.refreshInterval) return

      const interval = setInterval(() => {
        refreshFunction()

        if (forceRerender) {
          this.key++
        }
      }, this.options.refreshRate * 1000)

      this.refreshInterval = interval
    },
  },
}
</script>
