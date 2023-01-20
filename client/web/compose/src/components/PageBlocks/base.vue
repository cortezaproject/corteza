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

    editable: {
      type: Boolean,
      required: false,
      default: false,
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
     * If reloading data does not refresh the page block
     * You should attach :key="key" to it and increment it in the refreshFunction
     */
    refreshBlock (refreshFunction) {
      if (!this.allowsRefresh || this.refreshInterval) return

      const interval = setInterval(() => {
        refreshFunction()
      }, this.options.refreshRate * 1000)

      this.refreshInterval = interval
    },
  },
}
</script>
